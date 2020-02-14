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

func (su *StoreRepositoryImpl) UpdateStoreStuffPermission(db *gorm.DB, staff *models.Staff) error {
	st := models.Staff{}

	if err := db.Table(st.TableName()).Select("permission_id").Update(map[string]interface{}{
		"permission_id": staff.PermissionID,
	}).Where("store_id = ? AND user_id = ? AND is_creator = ?", staff.StoreID, staff.UserID, false).Error; err != nil {
		return err
	}
	return nil
}

func (su *StoreRepositoryImpl) DeleteStoreStuffPermission(db *gorm.DB, storeID, userID string) error {
	st := models.Staff{}

	if err := db.Table(st.TableName()).Delete(&st, "store_id = ? AND user_id = ? AND is_creator = ?", storeID, userID, false).Error; err != nil {
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

func (su *StoreRepositoryImpl) IsAlreadyStaff(db *gorm.DB, userID string) (bool, error) {
	staff := models.Staff{}

	var count int

	if err := db.Table(staff.TableName()).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (su *StoreRepositoryImpl) ListStaffs(db *gorm.DB, storeID string, from, limit int) ([]models.StoreUserProfile, error) {
	var staffs []models.StoreUserProfile
	sup := models.StoreUserProfile{}
	if err := db.Table(sup.TableName()).
		Where("id = ?", storeID).
		Offset(from).
		Limit(limit).
		Find(&staffs).Error; err != nil {
		return nil, err
	}
	return staffs, nil
}

func (su *StoreRepositoryImpl) SearchStaffs(db *gorm.DB, storeID, query string, from, limit int) ([]models.StoreUserProfile, error) {
	var staffs []models.StoreUserProfile
	sup := models.StoreUserProfile{}
	if err := db.Table(sup.TableName()).
		Where("id = ? AND (user_email LIKE ? OR user_phone LIKE ?)", storeID, "%"+query+"%", "%"+query+"%").
		Offset(from).
		Limit(limit).
		Find(&staffs).Error; err != nil {
		return nil, err
	}
	return staffs, nil
}
