package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/middlewares"
	"github.com/shopicano/shopicano-backend/repository"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"net/http"
)

func RegisterOrderRoutes(g *echo.Group) {
	func(g echo.Group) {
		g.Use(middlewares.AuthUser)
		g.POST("/", createOrder)
		g.GET("/", listOrders)
		g.GET("/search/", searchOrders)
		g.GET("/:order_id", getOrder)
	}(*g)

	func(g echo.Group) {
		g.Use(middlewares.IsStoreStaffWithStoreActivation)
		g.POST("/", createOrder)
		g.GET("/with-store", listOrdersWithStore)
		g.GET("/search/with-store", searchOrdersWithStore)
		g.GET("/:order_id/with-store", getOrderWithStore)
	}(*g)
}

func createOrder(ctx echo.Context) error {
	userID := ctx.Get(utils.UserID).(string)

	o, err := validators.ValidateCreateOrder(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.CollectionCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	o.UserID = userID

	ou := repository.NewOrderRepository()
	m, err := ou.CreateOrder(o)
	if err != nil {
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

func getOrderWithStore(ctx echo.Context) error {
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

func searchOrders(ctx echo.Context) error {
	//pageQ := ctx.Request().URL.Query().Get("page")
	//limitQ := ctx.Request().URL.Query().Get("limit")
	//query := ctx.Request().URL.Query().Get("query")
	//
	//page, err := strconv.ParseInt(pageQ, 10, 64)
	//if err != nil {
	//	page = 1
	//}
	//limit, err := strconv.ParseInt(limitQ, 10, 64)
	//if err != nil {
	//	limit = 10
	//}
	//
	//resp := core.Response{}
	//
	//from := (page - 1) * limit
	//cu := NewCollectionRepository()
	//collections, err := cu.SearchCollections(query, int(from), int(limit))
	//if err != nil {
	//	resp.Title = "Database query failed"
	//	resp.Status = http.StatusInternalServerError
	//	resp.Code = errors.DatabaseQueryFailed
	//	resp.Errors = err
	//	return resp.ServerJSON(ctx)
	//}
	//
	//resp.Status = http.StatusOK
	//resp.Data = collections
	//return resp.ServerJSON(ctx)
	return nil
}

func searchOrdersWithStore(ctx echo.Context) error {
	//pageQ := ctx.Request().URL.Query().Get("page")
	//limitQ := ctx.Request().URL.Query().Get("limit")
	//query := ctx.Request().URL.Query().Get("query")
	//
	//page, err := strconv.ParseInt(pageQ, 10, 64)
	//if err != nil {
	//	page = 1
	//}
	//limit, err := strconv.ParseInt(limitQ, 10, 64)
	//if err != nil {
	//	limit = 10
	//}
	//
	//resp := core.Response{}
	//
	//from := (page - 1) * limit
	//cu := NewCollectionRepository()
	//collections, err := cu.SearchCollections(query, int(from), int(limit))
	//if err != nil {
	//	resp.Title = "Database query failed"
	//	resp.Status = http.StatusInternalServerError
	//	resp.Code = errors.DatabaseQueryFailed
	//	resp.Errors = err
	//	return resp.ServerJSON(ctx)
	//}
	//
	//resp.Status = http.StatusOK
	//resp.Data = collections
	//return resp.ServerJSON(ctx)
	return nil
}

func listOrders(ctx echo.Context) error {
	//pageQ := ctx.Request().URL.Query().Get("page")
	//limitQ := ctx.Request().URL.Query().Get("limit")
	//
	//page, err := strconv.ParseInt(pageQ, 10, 64)
	//if err != nil {
	//	page = 1
	//}
	//limit, err := strconv.ParseInt(limitQ, 10, 64)
	//if err != nil {
	//	limit = 10
	//}
	//
	//resp := core.Response{}
	//
	//from := (page - 1) * limit
	//cu := NewCollectionRepository()
	//collections, err := cu.ListCollections(int(from), int(limit))
	//if err != nil {
	//	resp.Title = "Database query failed"
	//	resp.Status = http.StatusInternalServerError
	//	resp.Code = errors.DatabaseQueryFailed
	//	resp.Errors = err
	//	return resp.ServerJSON(ctx)
	//}
	//
	//resp.Status = http.StatusOK
	//resp.Data = collections
	//return resp.ServerJSON(ctx)
	return nil
}

func listOrdersWithStore(ctx echo.Context) error {
	//pageQ := ctx.Request().URL.Query().Get("page")
	//limitQ := ctx.Request().URL.Query().Get("limit")
	//
	//page, err := strconv.ParseInt(pageQ, 10, 64)
	//if err != nil {
	//	page = 1
	//}
	//limit, err := strconv.ParseInt(limitQ, 10, 64)
	//if err != nil {
	//	limit = 10
	//}
	//
	//resp := core.Response{}
	//
	//from := (page - 1) * limit
	//cu := NewCollectionRepository()
	//collections, err := cu.ListCollections(int(from), int(limit))
	//if err != nil {
	//	resp.Title = "Database query failed"
	//	resp.Status = http.StatusInternalServerError
	//	resp.Code = errors.DatabaseQueryFailed
	//	resp.Errors = err
	//	return resp.ServerJSON(ctx)
	//}
	//
	//resp.Status = http.StatusOK
	//resp.Data = collections
	//return resp.ServerJSON(ctx)
	return nil
}
