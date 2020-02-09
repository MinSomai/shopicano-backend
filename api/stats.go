package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/middlewares"
	"github.com/shopicano/shopicano-backend/utils"
	"net/http"
	"time"
)

func RegisterStatsRoutes(g *echo.Group) {
	func(g *echo.Group) {
		g.Use(middlewares.MightBeStoreStaffAndStoreActive)
		g.GET("/products/", productStats)
		g.GET("/categories/", categoryStats)
		//g.GET("/collections/", collectionStats)
		//g.GET("/stores/", storeStats)
	}(g)

	func(g *echo.Group) {
		g.Use(middlewares.IsStoreStaffAndStoreActive)
		g.GET("/orders/", orderStats)
		//g.GET("/collections/", collectionStats)
		//g.GET("/stores/", storeStats)
	}(g)
}

func productStats(ctx echo.Context) error {
	resp := core.Response{}

	db := app.DB()
	pu := data.NewProductRepository()

	var res interface{}
	var err error

	isPublic := !utils.IsStoreStaff(ctx)
	if !isPublic {
		res, err = pu.Stats(db, utils.GetStoreID(ctx), 0, 25)
	}

	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = res
	return resp.ServerJSON(ctx)
}

func categoryStats(ctx echo.Context) error {
	resp := core.Response{}

	db := app.DB()
	cu := data.NewCategoryRepository()

	var res interface{}
	var err error

	isPublic := !utils.IsStoreStaff(ctx)
	if isPublic {
		res, err = cu.Stats(db, 0, 25)
	} else {
		res, err = cu.StatsAsStoreStuff(db, utils.GetStoreID(ctx), 0, 25)
	}

	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = res
	return resp.ServerJSON(ctx)
}

func orderStats(ctx echo.Context) error {
	resp := core.Response{}

	db := app.DB()
	ou := data.NewOrderRepository()

	timeFrames := map[time.Time]time.Time{}

	timeline := ctx.QueryParam("timeline")
	switch timeline {
	case "w":
		now := time.Now()
		prev := now.Add(time.Hour * -24)

		for len(timeFrames) < 7 {
			timeFrames[prev] = now

			now = prev
			prev = now.Add(time.Hour * -24)
		}
	case "m":
		now := time.Now()
		prev := now.Add(time.Hour * 7 * -24)

		for len(timeFrames) < 5 {
			timeFrames[prev] = now

			now = prev
			prev = now.Add(time.Hour * 7 * -24)
		}
	case "y":
		now := time.Now()
		prev := now.Add(time.Hour * 30 * -24)

		for len(timeFrames) < 12 {
			timeFrames[prev] = now

			now = prev
			prev = now.Add(time.Hour * 30 * -24)
		}
	}

	var res []interface{}

	for k, v := range timeFrames {
		count, err := ou.CountByTimeAsStoreStuff(db, utils.GetStoreID(ctx), k, v)
		if err != nil {
			resp.Title = "Database query failed"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.DatabaseQueryFailed
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		res = append(res, struct {
			ID        int
			StartTime string `json:"start_time"`
			EndTime   string `json:"end_time"`
			Count     int    `json:"count"`
		}{
			StartTime: k.Format(utils.DateFormat),
			EndTime:   v.Format(utils.DateFormat),
			Count:     count,
		})
	}

	resp.Status = http.StatusOK
	resp.Data = res
	return resp.ServerJSON(ctx)
}
