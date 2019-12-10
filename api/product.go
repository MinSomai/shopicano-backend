package api

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/middlewares"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"net/http"
	"strconv"
)

func RegisterProductRoutes(g *echo.Group) {
	g.GET("/", listProducts)
	g.GET("/search/", searchProducts)
	g.GET("/:id/", getProduct)

	g.Use(middlewares.IsStoreStaffWithStoreActivation)
	g.POST("/", createProduct)
	g.GET("/with_store/", listProductsWithStore)
	g.GET("/search/with_store/", searchProductsWithStore)
	g.DELETE("/:id/", deleteProduct)
	g.PUT("/:id/", updateProduct)
	g.GET("/:id/with_store/", getProductWithStore)
}

func createProduct(ctx echo.Context) error {
	storeID := ctx.Get(utils.StoreID).(string)

	req, err := validators.ValidateCreateProduct(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.ProductCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	req.StoreID = storeID
	if req.CategoryID != nil && *req.CategoryID == "" {
		req.CategoryID = nil
	}

	pu := data.NewProductRepository()
	p, err := pu.CreateProduct(req)
	if err != nil {
		msg, ok := errors.IsDuplicateKeyError(err)
		if ok {
			resp.Title = msg
			resp.Status = http.StatusConflict
			resp.Code = errors.ProductAlreadyExists
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
	resp.Title = "Product created"
	resp.Data = p
	return resp.ServerJSON(ctx)
}

func updateProduct(ctx echo.Context) error {
	storeID := ctx.Get(utils.StoreID).(string)
	productID := ctx.Param("id")

	req, err := validators.ValidateUpdateProduct(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.ProductCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	req.StoreID = storeID
	if req.CategoryID != nil && *req.CategoryID == "" {
		req.CategoryID = nil
	}

	pu := data.NewProductRepository()
	p, err := pu.UpdateProduct(productID, req)
	if err != nil {
		msg, ok := errors.IsDuplicateKeyError(err)
		if ok {
			resp.Title = msg
			resp.Status = http.StatusConflict
			resp.Code = errors.ProductAlreadyExists
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		if err == gorm.ErrRecordNotFound {
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

	resp.Status = http.StatusOK
	resp.Title = "Product updated"
	resp.Data = p
	return resp.ServerJSON(ctx)
}

func deleteProduct(ctx echo.Context) error {
	storeID := ctx.Get(utils.StoreID).(string)
	productID := ctx.Param("id")

	resp := core.Response{}

	pu := data.NewProductRepository()
	if err := pu.DeleteProduct(storeID, productID); err != nil {
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

	resp.Status = http.StatusNoContent
	return resp.ServerJSON(ctx)
}

func getProductWithStore(ctx echo.Context) error {
	storeID := ctx.Get(utils.StoreID).(string)
	productID := ctx.Param("id")

	resp := core.Response{}

	pu := data.NewProductRepository()
	p, err := pu.GetProductWithStore(storeID, productID)
	if err != nil {
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

	resp.Status = http.StatusOK
	resp.Data = p
	return resp.ServerJSON(ctx)
}

func getProduct(ctx echo.Context) error {
	productID := ctx.Param("id")

	resp := core.Response{}

	pu := data.NewProductRepository()
	p, err := pu.GetProduct(productID)
	if err != nil {
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

	resp.Status = http.StatusOK
	resp.Data = p
	return resp.ServerJSON(ctx)
}

func searchProducts(ctx echo.Context) error {
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

	from := (page - 1) * limit
	pu := data.NewProductRepository()
	products, err := pu.SearchProducts(query, int(from), int(limit))
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = products
	return resp.ServerJSON(ctx)
}

func listProducts(ctx echo.Context) error {
	pageQ := ctx.Request().URL.Query().Get("page")
	limitQ := ctx.Request().URL.Query().Get("limit")

	page, err := strconv.ParseInt(pageQ, 10, 64)
	if err != nil {
		page = 1
	}
	limit, err := strconv.ParseInt(limitQ, 10, 64)
	if err != nil {
		limit = 10
	}

	resp := core.Response{}

	from := (page - 1) * limit
	pu := data.NewProductRepository()
	products, err := pu.ListProducts(int(from), int(limit))
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = products
	return resp.ServerJSON(ctx)
}

func searchProductsWithStore(ctx echo.Context) error {
	storeID := ctx.Get(utils.StoreID).(string)
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

	from := (page - 1) * limit
	pu := data.NewProductRepository()
	products, err := pu.SearchProductsWithStore(storeID, query, int(from), int(limit))
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = products
	return resp.ServerJSON(ctx)
}

func listProductsWithStore(ctx echo.Context) error {
	storeID := ctx.Get(utils.StoreID).(string)
	pageQ := ctx.Request().URL.Query().Get("page")
	limitQ := ctx.Request().URL.Query().Get("limit")

	page, err := strconv.ParseInt(pageQ, 10, 64)
	if err != nil {
		page = 1
	}
	limit, err := strconv.ParseInt(limitQ, 10, 64)
	if err != nil {
		limit = 10
	}

	resp := core.Response{}

	from := (page - 1) * limit
	pu := data.NewProductRepository()
	products, err := pu.ListProductsWithStore(storeID, int(from), int(limit))
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = products
	return resp.ServerJSON(ctx)
}
