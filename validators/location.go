package validators

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/errors"
	"strconv"
)

type reqUpdateLocation struct {
	IsPublished bool `json:"is_published"`
}

func ValidateUpdateLocation(ctx echo.Context, single bool) (*reqUpdateLocation, error) {
	body := reqUpdateLocation{}

	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}

	ve := errors.ValidationError{}

	if single {
		locationIDQ := ctx.Param("location_id")
		_, err := strconv.ParseInt(locationIDQ, 10, 64)
		if err != nil {
			ve.Add("location_id", "is invalid")
		}
	}

	ok, err := govalidator.ValidateStruct(&body)
	if ok {
		return &body, nil
	}

	if !ok {
		for k, v := range govalidator.ErrorsByField(err) {
			ve.Add(k, v)
		}
	}

	if len(ve) == 0 {
		return &body, nil
	}
	return nil, &ve
}
