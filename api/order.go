package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/middlewares"
	"github.com/shopicano/shopicano-backend/models"
	payment_gateways "github.com/shopicano/shopicano-backend/payment-gateways"
	"github.com/shopicano/shopicano-backend/queue"
	"github.com/shopicano/shopicano-backend/services"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"github.com/shopicano/shopicano-backend/values"
	"net/http"
	"strconv"
	"time"
)

func RegisterOrderRoutes(g *echo.Group) {
	g.POST("/:order_id/pay/", payOrder)
	g.GET("/:order_id/pay/", payOrder)

	func(g echo.Group) {
		g.Use(middlewares.MustBeUserOrStoreStaffAndStoreActive)
		g.GET("/", listOrders)
		g.GET("/:order_id/", getOrder)
	}(*g)

	func(g echo.Group) {
		g.Use(middlewares.AuthUser)
		g.POST("/", createOrder)
		g.POST("/:order_id/nonce/", generatePayNonce)
		g.POST("/:order_id/review/", createReview)
	}(*g)

	func(g echo.Group) {
		g.Use(middlewares.AuthUserWithQueryToken)
		g.GET("/:order_id/products/:product_id/download/", downloadProductAsUser)
		g.GET("/:order_id/nonce/", generatePayNonce)
	}(*g)

	func(g echo.Group) {
		g.Use(middlewares.IsStoreStaffAndStoreActive)
		g.POST("/internal/", createOrder)
		g.PATCH("/internal/:order_id/", createOrder)
		g.PUT("/internal/items/:order_id/", createOrder)
		g.DELETE("/internal/items/:order_id/", createOrder)
	}(*g)
}

func createOrder(ctx echo.Context) error {
	userID := ctx.Get(utils.UserID).(string)

	pld, err := validators.ValidateCreateOrder(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.OrderDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	pld.UserID = userID
	return createNewOrder(ctx, pld)
}

func createNewOrder(ctx echo.Context, pld *validators.ReqOrderCreate) error {
	resp := core.Response{}

	db := app.DB().Begin()

	o := models.Order{}
	o.ID = utils.NewUUID()
	o.Hash = utils.NewShortUUID()
	o.UserID = pld.UserID
	o.ShippingAddressID = pld.ShippingAddressID
	o.BillingAddressID = pld.BillingAddressID
	o.PaymentMethodID = pld.PaymentMethodID
	o.ShippingMethodID = pld.ShippingMethodID
	o.Status = models.OrderPending
	o.PaymentStatus = models.PaymentPending

	pu := data.NewProductRepository()
	ou := data.NewOrderRepository()
	au := data.NewAdminRepository()
	cu := data.NewCouponRepository()

	pm, err := au.GetPaymentMethod(db, o.PaymentMethodID)
	if err != nil {
		db.Rollback()

		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Payment method not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.PaymentMethodNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	var sm *models.ShippingMethod

	if o.ShippingMethodID != nil {
		sm, err = au.GetShippingMethod(db, *o.ShippingMethodID)
		if err != nil {
			db.Rollback()

			if errors.IsRecordNotFoundError(err) {
				resp.Title = "Shipping method not found"
				resp.Status = http.StatusNotFound
				resp.Code = errors.ShippingMethodNotFound
				resp.Errors = err
				return resp.ServerJSON(ctx)
			}

			resp.Title = "Database query failed"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.DatabaseQueryFailed
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}
	}

	hasDigitalProducts := false
	hasNonDigitalProducts := false

	isAllDigitalProduct := true

	var availableItems []*models.OrderedItem

	var storeID *string

	for _, v := range pld.Items {
		item, err := pu.GetForOrder(db, v.ID, v.Quantity)
		if err != nil {
			db.Rollback()

			if errors.IsRecordNotFoundError(err) {
				resp.Title = "Product unavailable"
				resp.Status = http.StatusNotFound
				resp.Code = errors.ProductUnavailable
				resp.Errors = err
				return resp.ServerJSON(ctx)
			}

			resp.Title = "Database query failed"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.DatabaseQueryFailed
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		if item.IsDigital {
			v.Quantity = 1
		}

		if storeID == nil {
			storeID = &item.StoreID
			o.StoreID = *storeID
		} else {
			if *storeID != item.StoreID {
				db.Rollback()

				resp.Title = "All products must be from same store"
				resp.Status = http.StatusBadRequest
				resp.Code = errors.AllProductsMustBeFromSameStore
				return resp.ServerJSON(ctx)
			}
		}

		if !hasDigitalProducts {
			hasDigitalProducts = item.IsDigital
		}
		if !hasNonDigitalProducts {
			hasNonDigitalProducts = !item.IsDigital
		}

		if isAllDigitalProduct {
			isAllDigitalProduct = item.IsDigital
		}

		oi := &models.OrderedItem{
			OrderID:     o.ID,
			ProductID:   item.ID,
			Quantity:    v.Quantity,
			Price:       item.Price,
			ProductCost: item.ProductCost,
		}
		oi.SubTotal = v.Quantity * item.Price

		availableItems = append(availableItems, oi)

		o.SubTotal += oi.SubTotal
	}

	if hasDigitalProducts && hasNonDigitalProducts {
		db.Rollback()

		resp.Title = "Cart must have all digital or all non-digital products"
		resp.Status = http.StatusBadRequest
		resp.Code = errors.CartMustHaveAllDigitalOrAllNonDigitalProducts
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	o.IsAllDigitalProducts = isAllDigitalProduct

	if !isAllDigitalProduct && sm == nil {
		db.Rollback()

		resp.Title = "Shipping method required"
		resp.Status = http.StatusBadRequest
		resp.Code = errors.ShippingMethodNotFound
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if !isAllDigitalProduct {
		o.ShippingCharge = sm.CalculateDeliveryCharge(0)
	}

	o.GrandTotal = o.SubTotal + o.ShippingCharge

	var couponID *string

	if pld.CouponCode != nil {
		coupon, err := cu.GetByCode(db, *storeID, *pld.CouponCode)
		if err != nil {
			db.Rollback()

			if errors.IsRecordNotFoundError(err) {
				resp.Title = "coupon not found"
				resp.Status = http.StatusNotFound
				resp.Code = errors.CouponNotFound
				resp.Errors = err
				return resp.ServerJSON(ctx)
			}

			resp.Title = "Database query failed"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.DatabaseQueryFailed
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		if !coupon.IsValid() {
			db.Rollback()

			resp.Title = "Coupon is invalid"
			resp.Status = http.StatusBadRequest
			resp.Code = errors.InvalidCoupon
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		if coupon.IsUserSpecific {
			ok, err := cu.HasUser(db, *storeID, coupon.ID, pld.UserID)
			if err != nil {
				db.Rollback()

				resp.Title = "Database query failed"
				resp.Status = http.StatusInternalServerError
				resp.Code = errors.DatabaseQueryFailed
				resp.Errors = err
				return resp.ServerJSON(ctx)
			}

			if !ok {
				db.Rollback()

				resp.Title = "Coupon not applicable for the user"
				resp.Status = http.StatusNotFound
				resp.Code = errors.CouponNotFound
				resp.Errors = err
				return resp.ServerJSON(ctx)
			}
		}

		previousUsage, err := cu.GetUsage(db, coupon.ID, o.UserID)
		if err != nil {
			db.Rollback()

			resp.Title = "Failed to get coupon usage"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.DatabaseQueryFailed
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		if coupon.MaxUsage != 0 && previousUsage >= coupon.MaxUsage {
			db.Rollback()

			resp.Title = "Coupon usage exceed"
			resp.Status = http.StatusBadRequest
			resp.Code = errors.InvalidCoupon
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		discount := 0
		switch coupon.DiscountType {
		case models.ProductDiscount:
			discount = coupon.CalculateDiscount(o.SubTotal)
		case models.ShippingDiscount:
			discount = coupon.CalculateDiscount(o.ShippingCharge)
		case models.TotalDiscount:
			discount = coupon.CalculateDiscount(o.GrandTotal)
		default:
			discount = 0
		}

		o.DiscountedAmount = discount
		couponID = &coupon.ID
	}

	pgName := payment_gateways.GetActivePaymentGateway().GetName()
	o.PaymentProcessingFee = pm.CalculateProcessingFee(o.GrandTotal)
	o.PaymentGateway = &pgName

	o.GrandTotal += o.PaymentProcessingFee
	o.OriginalGrandTotal = o.GrandTotal
	o.GrandTotal -= o.DiscountedAmount

	err = ou.Create(db, &o)
	if err != nil {
		db.Rollback()

		log.Log().Errorln(err)

		if errors.IsPreparedError(err) {
			resp.Title = "Invalid request"
			resp.Status = http.StatusBadRequest
			resp.Code = errors.InvalidRequest
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if couponID != nil {
		couponUsage := models.CouponUsage{
			CouponID: *couponID,
			UserID:   o.UserID,
			OrderID:  o.ID,
		}
		if err := cu.AddUsage(db, &couponUsage); err != nil {
			db.Rollback()

			resp.Title = "Failed to add coupon usage"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.DatabaseQueryFailed
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}
	}

	for _, v := range availableItems {
		if err := ou.AddOrderedItem(db, v); err != nil {
			db.Rollback()

			resp.Title = "Database query failed"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.DatabaseQueryFailed
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}
	}

	ol := models.OrderLog{
		ID:        utils.NewUUID(),
		OrderID:   o.ID,
		Action:    string(o.Status),
		Details:   "Order has been created",
		CreatedAt: time.Now(),
	}
	if err := ou.CreateLog(db, &ol); err != nil {
		db.Rollback()

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if isAllDigitalProduct {
		o.Status = models.OrderDelivered

		err = ou.UpdateStatus(db, &o)
		if err != nil {
			db.Rollback()

			log.Log().Errorln(err)

			if errors.IsPreparedError(err) {
				resp.Title = "Invalid request"
				resp.Status = http.StatusBadRequest
				resp.Code = errors.InvalidRequest
				resp.Errors = err
				return resp.ServerJSON(ctx)
			}

			resp.Title = "Database query failed"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.DatabaseQueryFailed
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		ol := models.OrderLog{
			ID:        utils.NewUUID(),
			OrderID:   o.ID,
			Action:    string(o.Status),
			Details:   "Order has been delivered",
			CreatedAt: time.Now(),
		}
		if err := ou.CreateLog(db, &ol); err != nil {
			db.Rollback()

			resp.Title = "Database query failed"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.DatabaseQueryFailed
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}
	}

	if o.GrandTotal == 0 {
		o.PaymentStatus = models.PaymentCompleted

		err = ou.UpdatePaymentStatus(db, &o)
		if err != nil {
			db.Rollback()

			log.Log().Errorln(err)

			if errors.IsPreparedError(err) {
				resp.Title = "Invalid request"
				resp.Status = http.StatusBadRequest
				resp.Code = errors.InvalidRequest
				resp.Errors = err
				return resp.ServerJSON(ctx)
			}

			resp.Title = "Database query failed"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.DatabaseQueryFailed
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		ol := models.OrderLog{
			ID:        utils.NewUUID(),
			OrderID:   o.ID,
			Action:    string(o.PaymentStatus),
			Details:   "Payment has been completed",
			CreatedAt: time.Now(),
		}
		if err := ou.CreateLog(db, &ol); err != nil {
			db.Rollback()

			resp.Title = "Database query failed"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.DatabaseQueryFailed
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}
	}

	m, err := ou.GetDetails(db, o.ID)
	if err != nil {
		db.Rollback()

		resp.Title = "Order not found"
		resp.Status = http.StatusNotFound
		resp.Code = errors.OrderNotFound
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if err := queue.SendOrderDetailsEmail(o.ID); err != nil {
		db.Rollback()

		resp.Title = "Failed to queue send order details"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.FailedToEnqueueTask
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if err := db.Commit().Error; err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusCreated
	resp.Data = m
	return resp.ServerJSON(ctx)
}

func createOrderWithStore(ctx echo.Context) error {
	//storeID := ctx.Get(utils.StoreID).(string)
	//
	//c, err := validateCreateCollection(ctx)
	//
	//resp := core.Response{}
	//
	//if err != nil {
	//	resp.Title = "Invalid data"
	//	resp.Status = http.StatusUnprocessableEntity
	//	resp.Code = errors.CollectionCreationDataInvalid
	//	resp.Errors = err
	//	return resp.ServerJSON(ctx)
	//}
	//
	//c.StoreID = storeID
	//
	//cu := NewCollectionRepository()
	//if err := cu.CreateCollection(c); err != nil {
	//	msg, ok := errors.IsDuplicateKeyError(err)
	//	if ok {
	//		resp.Title = msg
	//		resp.Status = http.StatusConflict
	//		resp.Code = errors.CollectionAlreadyExists
	//		resp.Errors = err
	//		return resp.ServerJSON(ctx)
	//	}
	//
	//	resp.Title = "Database query failed"
	//	resp.Status = http.StatusInternalServerError
	//	resp.Code = errors.DatabaseQueryFailed
	//	resp.Errors = err
	//	return resp.ServerJSON(ctx)
	//}
	//
	//resp.Status = http.StatusCreated
	//resp.Data = c
	//return resp.ServerJSON(ctx)
	return nil
}

func getOrder(ctx echo.Context) error {
	orderID := ctx.Param("order_id")

	resp := core.Response{}

	db := app.DB()

	var r interface{}
	var err error

	ou := data.NewOrderRepository()

	if utils.IsStoreStaff(ctx) {
		r, err = ou.GetDetailsAsStoreStuff(db, utils.GetStoreID(ctx), orderID)
	} else {
		r, err = ou.GetDetailsAsUser(db, utils.GetUserID(ctx), orderID)
	}

	if err != nil {
		resp.Title = "Order not found"
		resp.Status = http.StatusNotFound
		resp.Code = errors.OrderNotFound
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = r
	return resp.ServerJSON(ctx)
}

func listOrders(ctx echo.Context) error {
	pageQ := ctx.Request().URL.Query().Get("page")
	limitQ := ctx.Request().URL.Query().Get("limit")
	query := ctx.Request().URL.Query().Get("query")

	page, err := strconv.ParseInt(pageQ, 10, 64)
	if err != nil {
		page = 1
	}
	limit, err := strconv.ParseInt(limitQ, 10, 64)
	if err != nil {
		limit = 10
	}

	resp := core.Response{}

	var r interface{}

	if query == "" {
		r, err = fetchOrders(ctx, page, limit, !utils.IsStoreStaff(ctx))
	} else {
		r, err = searchOrders(ctx, query, page, limit, !utils.IsStoreStaff(ctx))
	}

	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = r
	return resp.ServerJSON(ctx)
}

func fetchOrders(ctx echo.Context, page, limit int64, isPublic bool) ([]models.OrderDetailsViewExternal, error) {
	db := app.DB()
	from := (page - 1) * limit
	ou := data.NewOrderRepository()

	log.Log().Infoln("UserID : ", ctx.Get(utils.UserID).(string))
	log.Log().Infoln("Offset : ", from)
	log.Log().Infoln("Limit : ", limit)
	log.Log().Infoln("IsPublic : ", isPublic)

	if isPublic {
		return ou.List(db, ctx.Get(utils.UserID).(string), int(from), int(limit))
	}

	log.Log().Infoln("StoreID : ", ctx.Get(utils.StoreID).(string))
	return ou.ListAsStoreStuff(db, ctx.Get(utils.StoreID).(string), int(from), int(limit))
}

func searchOrders(ctx echo.Context, query string, page, limit int64, isPublic bool) ([]models.OrderDetailsView, error) {
	db := app.DB()
	from := (page - 1) * limit
	ou := data.NewOrderRepository()

	if isPublic {
		return ou.Search(db, query, ctx.Get(utils.UserID).(string), int(from), int(limit))
	}
	return ou.SearchAsStoreStuff(db, query, ctx.Get(utils.StoreID).(string), int(from), int(limit))
}

func downloadProductAsUser(ctx echo.Context) error {
	orderID := ctx.Param("order_id")
	productID := ctx.Param("product_id")

	resp := core.Response{}

	db := app.DB()

	ou := data.NewOrderRepository()

	o, err := ou.GetDetailsAsUser(db, utils.GetUserID(ctx), orderID)
	if err != nil {
		log.Log().Errorln(err)

		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Order not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.OrderNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	_, err = ou.GetOrderedItem(db, o.ID, productID)
	if err != nil {
		log.Log().Errorln(err)

		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Product not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.ProductNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	pu := data.NewProductRepository()
	m, err := pu.Get(db, productID)
	if err != nil {
		log.Log().Errorln(err)

		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Product not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.ProductNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if !(o.Status == models.OrderConfirmed && o.PaymentStatus == models.PaymentCompleted) {
		resp.Title = "Unauthorized to download the product"
		resp.Status = http.StatusForbidden
		resp.Code = errors.UserScopeUnauthorized
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	f, err := services.ServeAsStreamFromMinio(fmt.Sprintf("%s/%s", values.ReservedBucketName, m.DigitalDownloadLink))
	if err != nil {
		log.Log().Errorln(err)

		resp.Title = "Minio service failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.MinioServiceFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	err = pu.IncreaseDownloadCounter(db, m)
	if err != nil {
		log.Log().Errorln(err)

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	return resp.ServeStreamFromMinio(ctx, f)
}
