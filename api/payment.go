package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/repositories"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"io/ioutil"
	"net/http"
)

func RegisterPaymentRoutes(g *echo.Group) {
	g.GET("/success/", func(ctx echo.Context) error {
		b, _ := ioutil.ReadAll(ctx.Request().Body)
		log.Log().Infoln(string(b))
		return nil
	})
	g.GET("/failure/", func(ctx echo.Context) error {
		b, _ := ioutil.ReadAll(ctx.Request().Body)
		log.Log().Infoln(string(b))
		return nil
	})
}

func onPaymentSuccess(ctx echo.Context) error {
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

	ou := repositories.NewOrderRepository()
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
