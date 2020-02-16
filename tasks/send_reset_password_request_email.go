package tasks

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/config"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/services"
	"github.com/shopicano/shopicano-backend/templates"
	"github.com/shopicano/shopicano-backend/utils"
	"time"
)

const (
	SendResetPasswordEmailTaskName             = "send_reset_password_email"
	SendResetPasswordConfirmationEmailTaskName = "send_reset_password_confirmation_email"
)

func SendResetPasswordEmailFn(userID string) error {
	db := app.DB().Begin()

	adminDao := data.NewAdminRepository()
	settings, err := adminDao.GetSettings(db)
	if err != nil {
		db.Rollback()

		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	userDao := data.NewUserRepository()
	u, err := userDao.Get(db, userID)
	if err != nil {
		db.Rollback()

		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	token := utils.NewToken()
	tokenAt := time.Now().Add(time.Hour * 24).UTC()

	u.ResetPasswordToken = &token
	u.ResetPasswordTokenGeneratedAt = &tokenAt

	if err := userDao.Update(db, u); err != nil {
		db.Rollback()

		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	address := fmt.Sprintf("%s, %s, %s - %s", settings.Address, settings.City, settings.Country, settings.Postcode)

	content, err := templates.GenerateResetPasswordEmailHTML(map[string]interface{}{
		"resetPasswordUrl": fmt.Sprintf("%s/#/reset-password?token=%s&email=%s", config.App().FrontStoreUrl, *u.ResetPasswordToken, u.Email),
		"shopicanoAddress": address,
		"shopicanoPhone":   settings.Phone,
		"shopicanoEmail":   settings.Email,
		"shopicanoWebsite": settings.Website,
	})
	if err != nil {
		db.Rollback()

		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	if err := services.SendEmail("Reset Password Requested", u.Email, content); err != nil {
		db.Rollback()

		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	if err := db.Commit().Error; err != nil {
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}
	return nil
}

func SendResetPasswordConfirmationEmailFn(userID string) error {
	db := app.DB()

	adminDao := data.NewAdminRepository()
	settings, err := adminDao.GetSettings(db)
	if err != nil {
		db.Rollback()

		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	userDao := data.NewUserRepository()
	u, err := userDao.Get(db, userID)
	if err != nil {
		db.Rollback()

		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	address := fmt.Sprintf("%s, %s, %s - %s", settings.Address, settings.City, settings.Country, settings.Postcode)

	content, err := templates.GenerateResetPasswordConfirmationEmailHTML(map[string]interface{}{
		"shopicanoAddress": address,
		"shopicanoPhone":   settings.Phone,
		"shopicanoEmail":   settings.Email,
		"shopicanoWebsite": settings.Website,
	})
	if err != nil {
		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	if err := services.SendEmail("Your password changed", u.Email, content); err != nil {
		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}
	return nil
}
