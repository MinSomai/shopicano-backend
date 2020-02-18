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
	SendSignUpVerificationEmailTaskName = "send_sign_up_verification_email"
)

func SendSignUpVerificationEmailFn(userID string) error {
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

	t := utils.NewToken()
	u.VerificationToken = &t

	if err := userDao.Update(db, u); err != nil {
		db.Rollback()

		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	verificationUrl := fmt.Sprintf("%s/v1/email-verification?uid=%s&token=%s",
		config.App().BackendUrl, userID, *u.VerificationToken)

	params := map[string]interface{}{
		"platformName":    settings.Name,
		"platformWebsite": settings.Website,
		"verificationUrl": verificationUrl,
		"userName":        u.Name,
	}

	body, err := templates.GenerateActivateAccountEmailHTML(params)
	if err != nil {
		db.Rollback()

		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	if err := services.SendSignUpVerificationEmail(u.Email, "Please verify your account", body); err != nil {
		db.Rollback()

		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	if err := db.Commit().Error; err != nil {
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	return nil
}
