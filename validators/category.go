package validators

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/utils"
	"time"
)

func ValidateCreateCategory(ctx echo.Context) (*models.Category, error) {
	pld := struct {
		Name        string `json:"name" valid:"required,stringlength(1|20)"`
		Description string `json:"description" valid:"required,stringlength(1|50)"`
		Image       string `json:"image"`
		IsPublished bool   `json:"is_published"`
	}{}

	if err := ctx.Bind(&pld); err != nil {
		return nil, err
	}

	ok, err := govalidator.ValidateStruct(&pld)
	if ok {
		return &models.Category{
			ID:          utils.NewUUID(),
			Name:        pld.Name,
			Description: pld.Description,
			Image:       pld.Image,
			IsPublished: pld.IsPublished,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}, nil
	}

	ve := errors.ValidationError{}

	for k, v := range govalidator.ErrorsByField(err) {
		ve.Add(k, v)
	}

	return nil, &ve
}
