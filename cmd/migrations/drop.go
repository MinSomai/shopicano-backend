package migration

import (
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/spf13/cobra"
)

var MigDropCmd = &cobra.Command{
	Use:   "drop",
	Short: "drop drops database tables",
	Run:   drop,
}

func drop(cmd *cobra.Command, args []string) {
	tx := app.DB().Begin()

	var tables []core.Table
	tables = append(tables, &models.AdditionalChargeOfProduct{}, &models.AdditionalCharge{})
	tables = append(tables, &models.OrderedProduct{}, &models.Order{})
	tables = append(tables, &models.Category{}, &models.Collection{}, &models.Product{}, &models.ProductOfCollection{})
	tables = append(tables, &models.ShippingMethod{}, &models.PaymentMethod{}, &models.Settings{})
	tables = append(tables, &models.StorePermission{}, &models.Store{}, &models.Staff{})
	tables = append(tables, &models.UserPermission{}, &models.User{}, &models.Session{})

	for _, t := range tables {
		if err := tx.DropTableIfExists(t).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Log().Errorln(err)
	}

	log.Log().Infoln("Migration drop completed")
}
