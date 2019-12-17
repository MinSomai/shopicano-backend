package migration

import (
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/spf13/cobra"
	"strings"
)

var MigAutoCmd = &cobra.Command{
	Use:   "auto",
	Short: "auto alter database tables if required",
	Run:   auto,
}

func auto(cmd *cobra.Command, args []string) {
	tx := app.DB().Begin()

	var tables []core.Table
	tables = append(tables, &models.Address{})
	tables = append(tables, &models.UserPermission{}, &models.User{}, &models.Session{})
	tables = append(tables, &models.StorePermission{}, &models.Store{}, &models.Staff{})
	tables = append(tables, &models.ShippingMethod{}, &models.PaymentMethod{}, &models.Settings{})
	tables = append(tables, &models.Category{}, &models.Collection{}, &models.Product{}, &models.ProductOfCollection{})
	tables = append(tables, &models.Order{}, &models.OrderedItem{})
	tables = append(tables, &models.AdditionalCharge{}, &models.AdditionalChargeOfProduct{})

	for _, t := range tables {
		if err := tx.AutoMigrate(t).Error; err != nil {
			tx.Rollback()
			log.Log().Errorln(err)
			return
		}
	}

	var tForeignKeys []core.Model
	tForeignKeys = append(tForeignKeys, &models.Address{}, &models.Category{}, &models.Collection{})
	tForeignKeys = append(tForeignKeys, &models.Order{}, &models.OrderedItem{})
	tForeignKeys = append(tForeignKeys, &models.Product{}, &models.ProductOfCollection{})
	tForeignKeys = append(tForeignKeys, &models.Settings{}, &models.Store{}, &models.Staff{})
	tForeignKeys = append(tForeignKeys, &models.User{}, &models.Session{})

	for _, t := range tForeignKeys {
		for _, fks := range t.ForeignKeys() {
			fk := strings.Split(fks, ";")
			if err := tx.Model(t).AddForeignKey(fk[0], fk[1], fk[2], fk[3]).Error; err != nil {
				tx.Rollback()
				log.Log().Errorln(err)
				return
			}
		}
	}

	var views []core.View
	views = append(views, &models.StoreUserProfile{})
	views = append(views, &models.OrderDetailsView{})
	views = append(views, &models.OrderedItemView{})

	for _, v := range views {
		if err := v.CreateView(tx); err != nil {
			tx.Rollback()
			log.Log().Errorln(err)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Log().Errorln(err)
		return
	}

	log.Log().Infoln("Migration auto completed")
}
