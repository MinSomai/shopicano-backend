package api

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/middlewares"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"net/http"
	"strconv"
	"strings"
)

func RegisterLocationRoutes(publicEndpoints, platformEndpoints *echo.Group) {
	locationsPublicPath := publicEndpoints.Group("/locations")
	locationsPlatformPath := platformEndpoints.Group("/locations")

	func(g echo.Group) {
		//g.Use(middlewares.MustBeUserOrStoreStaffAndStoreActive)
		g.GET("/", listLocations)
	}(*locationsPublicPath)

	func(g echo.Group) {
		g.Use(middlewares.IsPlatformManager)
		g.PATCH("/:location_id/", updateLocation)
		g.PATCH("/", enableLocation)
	}(*locationsPlatformPath)
}

func listLocations(ctx echo.Context) error {
	resp := core.Response{}

	name := ctx.QueryParam("name")
	//locationTypeQ := ctx.QueryParam("type")
	var locationType models.LocationType
	//if locationTypeQ == "city" {
	//	locationType = models.LocationTypeCity
	//} else if locationTypeQ == "state" {
	//	locationType = models.LocationTypeState
	//} else {
	//	locationType = models.LocationTypeCountry
	//}

	locationType = models.LocationTypeCountry

	parentIDQ := ctx.QueryParam("parent_id")
	parentID, err := strconv.ParseInt(parentIDQ, 10, 64)
	if locationType != models.LocationTypeCountry && err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusBadRequest
		resp.Code = errors.InvalidRequest
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	var locations []models.Location

	db := app.DB()

	if utils.IsPlatformAdmin(ctx) {
		locations, err = listLocationsForAdmin(name, locationType, parentID, db)
	} else {
		locations, err = listLocationsForUser(name, locationType, parentID, db)
	}

	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = locations
	return resp.ServerJSON(ctx)
}

func listLocationsForAdmin(name string, locationType models.LocationType, parentID int64, db *gorm.DB) ([]models.Location, error) {
	locDao := data.NewLocationRepository()

	where := "(type = ? AND parent_id = ?)"
	var args []interface{}
	args = append(args, locationType)
	args = append(args, parentID)

	if name != "" {
		where += " OR LOWER(name) LIKE ?"
		args = append(args, "%"+strings.ToLower(name)+"%")
	}

	return locDao.List(db, where, args)
}

func listLocationsForUser(name string, locationType models.LocationType, parentID int64, db *gorm.DB) ([]models.Location, error) {
	locDao := data.NewLocationRepository()

	where := "(type = ? AND parent_id = ? AND is_published = ?)"
	var args []interface{}
	args = append(args, locationType)
	args = append(args, parentID)
	args = append(args, 1)

	if name != "" {
		where += " OR LOWER(name) LIKE ?"
		args = append(args, "%"+strings.ToLower(name)+"%")
	}

	return locDao.List(db, where, args)
}

func updateLocation(ctx echo.Context) error {
	resp := core.Response{}

	req, err := validators.ValidateUpdateLocation(ctx, true)
	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusBadRequest
		resp.Code = errors.InvalidRequest
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	locationIDQ := ctx.Param("location_id")
	locationID, _ := strconv.ParseInt(locationIDQ, 10, 64)

	toggle := int64(0)
	if req.IsPublished {
		toggle = 1
	}

	db := app.DB()

	locDao := data.NewLocationRepository()
	err = locDao.UpdateByID(db, locationID, toggle)
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	return resp.ServerJSON(ctx)
}

func enableLocation(ctx echo.Context) error {
	resp := core.Response{}

	req, err := validators.ValidateUpdateLocation(ctx, false)
	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusBadRequest
		resp.Code = errors.InvalidRequest
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	toggle := int64(0)
	if req.IsPublished {
		toggle = 1
	}

	db := app.DB()

	locDao := data.NewLocationRepository()
	err = locDao.UpdateAll(db, toggle)
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	return resp.ServerJSON(ctx)
}
