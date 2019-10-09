package repositories

import (
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/values"
)

type StoreRepositoryImpl struct {
}

var storeRepository StoreRepository

func NewStoreRepository() StoreRepository {
	if storeRepository == nil {
		storeRepository = &StoreRepositoryImpl{}
	}
	return storeRepository
}

func (su *StoreRepositoryImpl) GetStoreUserProfile(userID string) (*models.StoreUserProfile, error) {
	db := app.DB()

	sup := models.StoreUserProfile{}
	if err := db.Table(sup.TableName()).Where("user_id = ?", userID).First(&sup).Error; err != nil {
		return nil, err
	}
	return &sup, nil
}

func (su *StoreRepositoryImpl) CreateStore(s *models.Store, userID string) error {
	tx := app.DB().Begin()

	if err := tx.Table(s.TableName()).Create(s).Error; err != nil {
		tx.Rollback()
		return err
	}

	st := models.Staff{
		UserID:       userID,
		StoreID:      s.ID,
		PermissionID: values.AdminGroupID,
	}

	if err := tx.Table(st.TableName()).Create(&st).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (su *StoreRepositoryImpl) FindStoreByID(ID string) (*models.Store, error) {
	db := app.DB()

	s := models.Store{}
	if err := db.Table(s.TableName()).Find(&s, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &s, nil
}
