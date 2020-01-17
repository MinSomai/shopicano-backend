package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
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

func (su *StoreRepositoryImpl) GetStoreUserProfile(db *gorm.DB, userID string) (*models.StoreUserProfile, error) {
	sup := models.StoreUserProfile{}
	if err := db.Table(sup.TableName()).Where("user_id = ?", userID).First(&sup).Error; err != nil {
		return nil, err
	}
	return &sup, nil
}

func (su *StoreRepositoryImpl) CreateStore(db *gorm.DB, s *models.Store) error {
	if err := db.Table(s.TableName()).Create(s).Error; err != nil {
		return err
	}
	return nil
}

func (su *StoreRepositoryImpl) AddStoreStuff(db *gorm.DB, st *models.Staff) error {
	if err := db.Table(st.TableName()).Create(&st).Error; err != nil {
		return err
	}
	return nil
}

func (su *StoreRepositoryImpl) UpdateStoreStuffPermission(db *gorm.DB, storeID, userID, permissionID string) error {
	st := models.Staff{}

	if err := db.Table(st.TableName()).Select("permission_id").Update(map[string]interface{}{
		"permission_id": permissionID,
	}).Where("store_id = ? AND user_id = ?", storeID, userID).Error; err != nil {
		return err
	}
	return nil
}

func (su *StoreRepositoryImpl) DeleteStoreStuffPermission(db *gorm.DB, storeID, userID string) error {
	st := models.Staff{}

	if err := db.Table(st.TableName()).Delete(&st, "store_id = ? AND user_id = ?", storeID, userID).Error; err != nil {
		return err
	}
	return nil
}

func (su *StoreRepositoryImpl) FindStoreByID(db *gorm.DB, ID string) (*models.Store, error) {
	s := models.Store{}
	if err := db.Table(s.TableName()).Find(&s, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &s, nil
}
