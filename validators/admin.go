package validators

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/errors"
)

type ReqPaymentMethodCreate struct {
	Name        string `json:"name" valid:"required"`
	IsPublished bool   `json:"is_published"`
}

func ValidateCreatePaymentMethod(ctx echo.Context) (*ReqPaymentMethodCreate, error) {
	pld := ReqPaymentMethodCreate{}
	if err := ctx.Bind(&pld); err != nil {
		return nil, err
	}

	ok, err := govalidator.ValidateStruct(&pld)
	if ok {
		return &pld, nil
	}

	ve := errors.ValidationError{}

	for k, v := range govalidator.ErrorsByField(err) {
		ve.Add(k, v)
	}

	return nil, &ve
}

type ReqShippingMethodCreate struct {
	Name                    string `json:"name" valid:"required"`
	ApproximateDeliveryTime int    `json:"approximate_delivery_time" valid:"required"`
	DeliveryCharge          int    `json:"delivery_charge" valid:"required"`
	IsPublished             bool   `json:"is_published"`
	WeightUnit              string `json:"weight_unit" valid:"required"`
}

func ValidateCreateShippingMethod(ctx echo.Context) (*ReqShippingMethodCreate, error) {
	pld := ReqShippingMethodCreate{}
	if err := ctx.Bind(&pld); err != nil {
		return nil, err
	}

	ok, err := govalidator.ValidateStruct(&pld)
	if ok {
		return &pld, nil
	}

	ve := errors.ValidationError{}

	for k, v := range govalidator.ErrorsByField(err) {
		ve.Add(k, v)
	}

	return nil, &ve
}

type ReqAdditionalChargeCreate struct {
	Name        string `json:"name" valid:"required"`
	ChargeType  string `json:"charge_type" valid:"required"`
	Amount      int    `json:"amount" valid:"required,range(1|1000000)"`
	AmountType  string `json:"amount_type" valid:"required"`
	AmountMax   int    `json:"amount_max" valid:"range(0|1000000)"`
	AmountMin   int    `json:"amount_min" valid:"range(0|1000000)"`
	IsPublished bool   `json:"is_published"`
}

func ValidateCreateAdditionalCharge(ctx echo.Context) (*ReqAdditionalChargeCreate, error) {
	pld := ReqAdditionalChargeCreate{}
	if err := ctx.Bind(&pld); err != nil {
		return nil, err
	}

	ok, err := govalidator.ValidateStruct(&pld)
	if ok {
		return &pld, nil
	}

	ve := errors.ValidationError{}

	for k, v := range govalidator.ErrorsByField(err) {
		ve.Add(k, v)
	}

	return nil, &ve
}
