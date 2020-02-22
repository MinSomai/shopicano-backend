package validators

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/values"
	"time"
)

func ValidateCreateStore(ctx echo.Context) (*models.Store, error) {
	pld := struct {
		Name        string `json:"name" valid:"required,stringlength(1|100)"`
		Address     string `json:"address" valid:"required,stringlength(1|100)"`
		City        string `json:"city" valid:"required,stringlength(1|30)"`
		Country     string `json:"country" valid:"required,stringlength(1|30)"`
		Postcode    string `json:"postcode" valid:"required,stringlength(1|100)"`
		Email       string `json:"email" valid:"required,email"`
		Phone       string `json:"phone" valid:"required,stringlength(1|20)"`
		Description string `json:"description" valid:"required,stringlength(1|1000)"`
		Image       string `json:"image"`
	}{}

	if err := ctx.Bind(&pld); err != nil {
		return nil, err
	}

	ok, err := govalidator.ValidateStruct(&pld)
	if ok {
		return &models.Store{
			ID:                       utils.NewUUID(),
			Name:                     pld.Name,
			Email:                    pld.Email,
			Phone:                    pld.Phone,
			Status:                   models.StoreRegistered,
			Address:                  pld.Address,
			Description:              pld.Description,
			IsOrderCreationEnabled:   false,
			IsProductCreationEnabled: false,
			Postcode:                 pld.Postcode,
			City:                     pld.City,
			Country:                  pld.Country,
			CreatedAt:                time.Now().UTC(),
			UpdatedAt:                time.Now().UTC(),
		}, nil
	}

	ve := errors.ValidationError{}

	for k, v := range govalidator.ErrorsByField(err) {
		ve.Add(k, v)
	}

	return nil, &ve
}

func ValidateCreateStoreStaff(ctx echo.Context) (*string, *string, error) {
	pld := struct {
		Email        string `json:"email" valid:"required,email"`
		PermissionID string `json:"permission_id" valid:"required"`
	}{}

	if err := ctx.Bind(&pld); err != nil {
		return nil, nil, err
	}

	ve := errors.ValidationError{}

	ok, err := govalidator.ValidateStruct(&pld)
	if ok {
		if pld.PermissionID == values.AdminGroupID || pld.PermissionID == values.ManagerGroupID {
			return &pld.Email, &pld.PermissionID, nil
		}

		ve.Add("permission_id", "is invalid")
	}

	for k, v := range govalidator.ErrorsByField(err) {
		ve.Add(k, v)
	}

	return nil, nil, &ve
}

func ValidateUpdateStoreStaff(ctx echo.Context) (*string, *string, error) {
	pld := struct {
		Email        string `json:"user_id" valid:"required"`
		PermissionID string `json:"permission_id" valid:"required"`
	}{}

	if err := ctx.Bind(&pld); err != nil {
		return nil, nil, err
	}

	ve := errors.ValidationError{}

	ok, err := govalidator.ValidateStruct(&pld)
	if ok {
		if pld.PermissionID == values.AdminGroupID || pld.PermissionID == values.ManagerGroupID {
			return &pld.Email, &pld.PermissionID, nil
		}

		ve.Add("permission_id", "is invalid")
	}

	for k, v := range govalidator.ErrorsByField(err) {
		ve.Add(k, v)
	}

	return nil, nil, &ve
}

func ValidateUpdateStoreStatus(ctx echo.Context) (*models.StoreStatus, *int64, error) {
	pld := struct {
		Status         *models.StoreStatus `json:"status"`
		CommissionRate *int64              `json:"commission_rate"`
	}{}

	if err := ctx.Bind(&pld); err != nil {
		return nil, nil, err
	}

	ve := errors.ValidationError{}

	if pld.Status != nil && !pld.Status.IsValid() {
		ve.Add("status", "is invalid")
	}
	if pld.CommissionRate != nil && (*pld.CommissionRate < 0 || *pld.CommissionRate > 100) {
		ve.Add("commission_rate", "is invalid")
	}

	if len(ve) > 0 {
		return nil, nil, &ve
	}

	return pld.Status, pld.CommissionRate, nil
}
