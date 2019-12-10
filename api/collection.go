package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/middlewares"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"net/http"
	"strconv"
)

func RegisterCollectionRoutes(g *echo.Group) {
	func(*echo.Group) {
		g.Use(middlewares.MightBeStoreStaffWithStoreActivation)
		g.GET("/", listCollections)
	}(g)

	func(*echo.Group) {
		// Private endpoints only
		g.Use(middlewares.IsStoreStaffWithStoreActivation)
		g.POST("/", createCollection)
		g.DELETE("/:id/", deleteCollection)
		g.PUT("/:id/", updateCollection)
	}(g)
}

func createCollection(ctx echo.Context) error {
	storeID := ctx.Get(utils.StoreID).(string)

	c, err := validators.ValidateCreateCollection(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.CollectionCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	c.StoreID = storeID

	cu := data.NewCollectionRepository()
	if err := cu.CreateCollection(c); err != nil {
		msg, ok := errors.IsDuplicateKeyError(err)
		if ok {
			resp.Title = msg
			resp.Status = http.StatusConflict
			resp.Code = errors.CollectionAlreadyExists
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
	resp.Data = c
	return resp.ServerJSON(ctx)
}

func updateCollection(ctx echo.Context) error {
	return nil
}

func deleteCollection(ctx echo.Context) error {
	storeID := ctx.Get(utils.StoreID).(string)
	collectionID := ctx.Param("id")

	resp := core.Response{}

	cu := data.NewCollectionRepository()
	if err := cu.DeleteCollection(storeID, collectionID); err != nil {
		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Collection not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.CollectionNotFound
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

func listCollections(ctx echo.Context) error {
	pageQ := ctx.Request().URL.Query().Get("page")
	limitQ := ctx.Request().URL.Query().Get("limit")
	query := ctx.Request().URL.Query().Get("query")

	var err error

	page, err := strconv.ParseInt(pageQ, 10, 64)
	if err != nil {
		page = 1
	}
	limit, err := strconv.ParseInt(limitQ, 10, 64)
	if err != nil {
		limit = 10
	}

	resp := core.Response{}

	var collections []models.Collection

	if query == "" {
		collections, err = fetchCollections(ctx, page, limit, utils.IsStoreStaff(ctx))
	} else {
		collections, err = searchCollections(ctx, query, page, limit, utils.IsStoreStaff(ctx))
	}

	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = collections
	return resp.ServerJSON(ctx)
}

func searchCollections(ctx echo.Context, query string, page, limit int64, isPublic bool) ([]models.Collection, error) {
	from := (page - 1) * limit
	cu := data.NewCollectionRepository()
	if isPublic {
		return cu.SearchCollections(query, int(from), int(limit))
	}
	return cu.SearchCollectionsWithStore(ctx.Get(utils.StoreID).(string), query, int(from), int(limit))
}

func fetchCollections(ctx echo.Context, page, limit int64, isPublic bool) ([]models.Collection, error) {
	from := (page - 1) * limit
	cu := data.NewCollectionRepository()
	if isPublic {
		return cu.ListCollections(int(from), int(limit))
	}
	return cu.ListCollectionsWithStore(ctx.Get(utils.StoreID).(string), int(from), int(limit))
}
