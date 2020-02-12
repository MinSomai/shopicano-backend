package tasks

import (
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/services"
	"github.com/shopicano/shopicano-backend/utils"
	"time"
)

const (
	SendSignUpVerificationEmailTaskName = "send_sign_up_verification_email"
)

func SendSignUpVerificationEmailFn(userID string) error {
	db := app.DB().Begin()

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

	if err := services.SendSignUpVerificationEmail(u.Name, u.Email, u.ID, *u.VerificationToken); err != nil {
		db.Rollback()

		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	if err := db.Commit().Error; err != nil {
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	return nil
}
