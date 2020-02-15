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
	"net/http"
	"strconv"
	"strings"
)

func RegisterLocationRoutes(g *echo.Group) {
	func(g echo.Group) {
		g.Use(middlewares.MustBeUserOrStoreStaffAndStoreActive)
		g.GET("/", listLocations)
	}(*g)
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
