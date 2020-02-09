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

func ValidateLogin(ctx echo.Context) (string, string, error) {
	ul := struct {
		Email    string `json:"email" valid:"required,stringlength(3|100)"`
		Password string `json:"password"`
	}{}

	if err := ctx.Bind(&ul); err != nil {
		return "", "", err
	}

	ok, err := govalidator.ValidateStruct(&ul)
	if ok {
		return ul.Email, ul.Password, nil
	}

	ve := errors.ValidationError{}

	for k, v := range govalidator.ErrorsByField(err) {
		ve.Add(k, v)
	}

	return "", "", &ve
}

func ValidateRegister(ctx echo.Context) (*models.User, error) {
	ur := struct {
		Name           string  `json:"name" valid:"required,stringlength(3|100)"`
		Email          string  `json:"email" valid:"required,email"`
		ProfilePicture *string `json:"profile_picture"`
		Phone          *string `json:"phone"`
		Password       string  `json:"password" valid:"required,stringlength(8|100)"`
	}{}

	if err := ctx.Bind(&ur); err != nil {
		return nil, err
	}

	ok, err := govalidator.ValidateStruct(&ur)
	if ok {
		return &models.User{
			ID:             utils.NewUUID(),
			Name:           ur.Name,
			Email:          ur.Email,
			Password:       ur.Password,
			Phone:          ur.Phone,
			ProfilePicture: ur.ProfilePicture,
			Status:         models.UserRegistered,
			PermissionID:   values.UserGroupID,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		}, nil
	}

	ve := errors.ValidationError{}

	for k, v := range govalidator.ErrorsByField(err) {
		ve.Add(k, v)
	}

	return nil, &ve
}

type reqUserUpdate struct {
	Name           *string `json:"name"`
	Email          *string `json:"email"`
	ProfilePicture *string `json:"profile_picture"`
	Phone          *string `json:"phone"`
	Password       *string `json:"password"`
}

func ValidateUserUpdate(ctx echo.Context) (*reqUserUpdate, error) {
	body := reqUserUpdate{}

	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

type reqUserUpdateStatus struct {
	NewStatus models.UserStatus `json:"new_status" valid:"required"`
}

func ValidateUserUpdateStatus(ctx echo.Context) (*reqUserUpdateStatus, error) {
	body := reqUserUpdateStatus{}

	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}

	ok, err := govalidator.ValidateStruct(&body)
	if ok {
		return &body, nil
	}

	ve := errors.ValidationError{}

	for k, v := range govalidator.ErrorsByField(err) {
		ve.Add(k, v)
	}

	if !body.NewStatus.IsValid() {
		ve.Add("new_status", "is invalid")
	}

	return nil, &ve
}

type reqUserUpdatePermission struct {
	NewPermissionID string `json:"new_permission_id" valid:"required"`
}

func ValidateUserUpdatePermission(ctx echo.Context) (*reqUserUpdatePermission, error) {
	body := reqUserUpdatePermission{}

	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}

	ok, err := govalidator.ValidateStruct(&body)
	if ok {
		return &body, nil
	}

	ve := errors.ValidationError{}

	for k, v := range govalidator.ErrorsByField(err) {
		ve.Add(k, v)
	}

	return nil, &ve
}
