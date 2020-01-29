package validators

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/errors"
)

type ReqPaymentMethodCreate struct {
	Name             string `json:"name" valid:"required"`
	IsPublished      bool   `json:"is_published"`
	IsFlat           bool   `json:"is_flat"`
	ProcessingFee    int    `json:"processing_fee"`
	MinProcessingFee int    `json:"min_processing_fee"`
	MaxProcessingFee int    `json:"max_processing_fee"`
	IsOfflinePayment bool   `json:"is_offline_payment"`
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
	IsFlat                  bool   `json:"is_flat"`
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
