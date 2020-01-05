package api

import (
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
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"net/http"
	"strconv"
)

func RegisterOrderRoutes(g *echo.Group) {
	g.POST("/:order_id/pay/", payOrder)
	g.GET("/:order_id/pay/", payOrder)

	func(g echo.Group) {
		g.Use(middlewares.MustBeUserOrStoreStaffWithStoreActivation)
		g.GET("/", listOrders)
		g.GET("/:order_id/", getOrder)
	}(*g)

	func(g echo.Group) {
		g.Use(middlewares.AuthUser)
		g.POST("/", createOrder)
		g.POST("/:order_id/nonce/", generatePayNonce)
	}(*g)

	func(g echo.Group) {
		g.Use(middlewares.IsStoreStaffWithStoreActivation)
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
	o.StoreID = pld.StoreID
	o.ShippingAddressID = pld.ShippingAddressID
	o.BillingAddressID = pld.BillingAddressID
	o.PaymentMethodID = pld.PaymentMethodID
	o.ShippingMethodID = pld.ShippingMethodID
	o.IsPaid = false

	pu := data.NewProductRepository()
	ou := data.NewOrderRepository()
	au := data.NewAdminRepository()

	pm, err := au.GetPaymentMethod(o.PaymentMethodID)
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
		sm, err = au.GetShippingMethod(*o.ShippingMethodID)
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

	var availableItems []*models.OrderedItem

	for _, v := range pld.Items {
		item, err := pu.GetForOrder(db, o.StoreID, v.ID, v.Quantity)
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

		oi := &models.OrderedItem{
			OrderID:   o.ID,
			ProductID: item.ID,
			Quantity:  v.Quantity,
			Price:     item.Price,
			TotalVat:  0,
			TotalTax:  0,
		}
		oi.SubTotal = v.Quantity * item.Price

		availableItems = append(availableItems, oi)

		o.SubTotal += oi.SubTotal
		o.TotalTax += oi.TotalTax
		o.TotalVat += oi.TotalVat
	}

	if o.ShippingMethodID != nil {
		o.ShippingCharge = sm.CalculateDeliveryCharge(0)
	}

	pgName := payment_gateways.GetActivePaymentGateway().GetName()

	o.GrandTotal = o.SubTotal + o.TotalTax + o.TotalVat + o.ShippingCharge
	o.PaymentProcessingFee = pm.CalculateProcessingFee(o.GrandTotal)
	o.PaymentGateway = &pgName
	o.Status = models.Pending

	err = ou.Create(db, &o)
	if err != nil {
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
