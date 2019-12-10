package data

import (
	"fmt"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"strings"
	"time"
)

type ProductRepositoryImpl struct {
}

var productRepository ProductRepository

func NewProductRepository() ProductRepository {
	if productRepository == nil {
		productRepository = &ProductRepositoryImpl{}
	}

	return productRepository
}

func (pu *ProductRepositoryImpl) CreateProduct(req *validators.ReqProductCreate) (*models.Product, error) {
	tx := app.DB().Begin()

	images := ""
	for _, i := range req.AdditionalImages {
		if strings.TrimSpace(i) == "" {
			continue
		}

		if images != "" {
			images += ","
		}
		images += strings.TrimSpace(i)
	}

	p := models.Product{
		ID:                  utils.NewUUID(),
		StoreID:             req.StoreID,
		Price:               req.Price,
		Stock:               req.Stock,
		Name:                req.Name,
		IsShippable:         req.IsShippable,
		CategoryID:          req.CategoryID,
		IsPublished:         req.IsPublished,
		IsDigital:           req.IsDigital,
		AdditionalImages:    images,
		SKU:                 req.SKU,
		Unit:                req.Unit,
		DigitalDownloadLink: req.DigitalDownloadLink,
		Image:               req.Image,
		Description:         req.Description,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}

	if err := tx.Table(p.TableName()).Create(&p).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, c := range req.CollectionsToAdd {
		cop := models.ProductOfCollection{
			StoreID:      req.StoreID,
			ProductID:    p.ID,
			CollectionID: c,
		}
		if err := tx.Table(cop.TableName()).Create(&cop).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, ac := range req.AdditionalChargesToAdd {
		acp := models.AdditionalChargeOfProduct{
			ProductID:          p.ID,
			AdditionalChargeID: ac,
		}
		if err := tx.Table(acp.TableName()).
			Create(&acp).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (pu *ProductRepositoryImpl) UpdateProduct(productID string, req *validators.ReqProductUpdate) (*models.Product, error) {
	tx := app.DB().Begin()

	images := ""
	for _, i := range req.AdditionalImages {
		if strings.TrimSpace(i) == "" {
			continue
		}

		if images != "" {
			images += ","
		}
		images += strings.TrimSpace(i)
	}

	p := models.Product{}
	if err := tx.Table(p.TableName()).Where("id = ?", productID).First(&p).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	p.ID = productID
	p.StoreID = req.StoreID
	p.Price = req.Price
	p.Stock = req.Stock
	p.Name = req.Name
	p.IsShippable = req.IsShippable
	p.CategoryID = req.CategoryID
	p.IsPublished = req.IsPublished
	p.IsDigital = req.IsDigital
	p.AdditionalImages = images
	p.SKU = req.SKU
	p.Unit = req.Unit
	p.DigitalDownloadLink = req.DigitalDownloadLink
	p.Image = req.Image
	p.Description = req.Description
	p.UpdatedAt = time.Now().UTC()

	if err := tx.Table(p.TableName()).Where("id = ?", productID).Save(&p).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, c := range req.CollectionsToAdd {
		cop := models.ProductOfCollection{
			StoreID:      req.StoreID,
			ProductID:    p.ID,
			CollectionID: c,
		}
		if err := tx.Table(cop.TableName()).FirstOrCreate(&cop).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, c := range req.CollectionsToRemove {
		cop := models.ProductOfCollection{}
		if err := tx.Table(cop.TableName()).
			Where("collection_id = ? AND store_id = ? AND product_id = ?", c, p.StoreID, p.ID).
			Delete(&cop).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, ac := range req.AdditionalChargesToAdd {
		acp := models.AdditionalChargeOfProduct{
			ProductID:          p.ID,
			AdditionalChargeID: ac,
		}
		if err := tx.Table(acp.TableName()).
			FirstOrCreate(&acp).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, ac := range req.AdditionalChargesToRemove {
		acp := models.AdditionalChargeOfProduct{}
		if err := tx.Table(acp.TableName()).
			Where("product_id = ? AND additional_charge_id = ?", p.ID, ac).
			Delete(&acp).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (pu *ProductRepositoryImpl) ListProducts(from, limit int) ([]models.ProductDetails, error) {
	db := app.DB()
	var ps []models.ProductDetails
	p := models.Product{}
	if err := db.Table(p.TableName()).
		Select("products.id, products.stock, products.name, products.price, products.description, products.is_published, products.is_shippable, products.is_digital, c.id AS category_id, c.name AS category_name, products.image, products.created_at, products.updated_at").
		Joins("LEFT JOIN categories AS c ON products.category_id = c.id").
		Where("products.is_published = ?", true).
		Offset(from).Limit(limit).
		Order("products.created_at DESC").Find(&ps).Error; err != nil {
		return nil, err
	}
	return ps, nil
}

func (pu *ProductRepositoryImpl) ListProductsWithStore(storeID string, from, limit int) ([]models.ProductDetails, error) {
	db := app.DB()
	var ps []models.ProductDetails
	p := models.Product{}
	if err := db.Table(p.TableName()).
		Select("products.id, products.stock, products.sku, products.name, products.price, products.description, products.is_published, products.is_shippable, products.is_digital, c.id AS category_id, c.name AS category_name, products.image, products.created_at, products.updated_at").
		Joins("LEFT JOIN categories AS c ON products.category_id = c.id").
		Where("products.store_id = ?", storeID).
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&ps).Error; err != nil {
		return nil, err
	}
	return ps, nil
}

func (pu *ProductRepositoryImpl) SearchProducts(query string, from, limit int) ([]models.ProductDetails, error) {
	db := app.DB()
	var ps []models.ProductDetails
	p := models.Product{}
	if err := db.Table(p.TableName()).
		Select("products.id, products.name, products.stock, products.price, products.description, products.is_published, products.is_shippable, products.is_digital, c.id AS category_id, c.name AS category_name, products.image, products.created_at, products.updated_at").
		Joins("LEFT JOIN categories AS c ON products.category_id = c.id").
		Where("products.is_published = ? AND (LOWER(products.name) LIKE ? OR LOWER(c.name) LIKE ?)", true, "%"+strings.ToLower(query)+"%", "%"+strings.ToLower(query)+"%").
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&ps).Error; err != nil {
		return nil, err
	}
	return ps, nil
}

func (pu *ProductRepositoryImpl) SearchProductsWithStore(storeID, query string, from, limit int) ([]models.ProductDetails, error) {
	db := app.DB()
	var ps []models.ProductDetails
	p := models.Product{}
	if err := db.Table(p.TableName()).
		Select("products.id, products.name, products.stock, products.price, products.description, products.is_published, products.is_shippable, products.is_digital, c.id AS category_id, c.name AS category_name, products.image, products.created_at, products.updated_at").
		Joins("LEFT JOIN categories AS c ON products.category_id = c.id").
		Where("products.store_id = ? AND (LOWER(products.name) LIKE ? OR LOWER(c.name) LIKE ?)", storeID, "%"+strings.ToLower(query)+"%", "%"+strings.ToLower(query)+"%").
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&ps).Error; err != nil {
		return nil, err
	}
	return ps, nil
}

func (pu *ProductRepositoryImpl) DeleteProduct(storeID, productID string) error {
	tx := app.DB().Begin()

	poc := models.ProductOfCollection{}
	if err := tx.Table(poc.TableName()).Where("store_id = ? AND product_id = ?", storeID, productID).Delete(&poc).Error; err != nil {
		tx.Rollback()
		return err
	}

	p := models.Product{}
	if err := tx.Table(p.TableName()).
		Where("store_id = ? AND id = ?", storeID, productID).
		Delete(&p).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (pu *ProductRepositoryImpl) GetProduct(productID string) (*models.ProductDetails, error) {
	db := app.DB()
	p := models.Product{}
	ps := models.ProductDetails{}
	cat := models.Category{}
	store := models.Store{}

	if err := db.Table(fmt.Sprintf("%s", p.TableName())).
		Select("products.id, s.id AS store_id, s.name AS store_name, products.digital_download_link, products.price, products.unit, products.stock, products.sku, products.additional_images, products.name, products.description, products.is_published, products.is_shippable, products.is_digital, c.id AS category_id, c.name AS category_name, products.image, products.created_at, products.updated_at").
		Joins(fmt.Sprintf("LEFT JOIN %s AS c ON products.category_id = c.id", cat.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS s ON products.store_id = s.id", store.TableName())).
		Where("products.id = ? AND products.is_published = ?", productID, true).
		First(&ps).Error; err != nil {
		return nil, err
	}

	poc := models.ProductOfCollection{}
	c := models.Collection{}
	var collections []models.CollectionDetails
	if err := db.Table(fmt.Sprintf("%s AS poc", poc.TableName())).
		Select("c.id, c.name, c.description").
		Joins(fmt.Sprintf("JOIN %s AS c ON poc.collection_id = c.id", c.TableName())).
		Where("poc.product_id = ?", productID).
		Scan(&collections).Error; err != nil {
		return nil, err
	}

	acp := models.AdditionalChargeOfProduct{}
	ac := models.AdditionalCharge{}
	var additionalCharges []models.AdditionalChargeDetails
	if err := db.Table(fmt.Sprintf("%s AS acp", acp.TableName())).
		Select("ac.id, ac.name, ac.charge_type, ac.amount, ac.amount_type, ac.amount_max, ac.amount_min").
		Joins(fmt.Sprintf("JOIN %s AS ac ON acp.additional_charge_id = ac.id", ac.TableName())).
		Where("acp.product_id = ?", productID).
		Scan(&additionalCharges).Error; err != nil {
		return nil, err
	}

	ps.Collections = collections
	ps.AdditionalCharges = additionalCharges
	return &ps, nil
}

func (pu *ProductRepositoryImpl) GetProductWithStore(storeID, productID string) (*models.ProductDetails, error) {
	db := app.DB()
	p := models.Product{}
	ps := models.ProductDetails{}
	cat := models.Category{}
	store := models.Store{}

	if err := db.Table(fmt.Sprintf("%s", p.TableName())).
		Select("products.id, s.id AS store_id, s.name AS store_name, products.digital_download_link, products.price, products.unit, products.stock, products.sku, products.additional_images, products.name, products.description, products.is_published, products.is_shippable, products.is_digital, c.id AS category_id, c.name AS category_name, products.image, products.created_at, products.updated_at").
		Joins(fmt.Sprintf("LEFT JOIN %s AS c ON products.category_id = c.id", cat.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS s ON products.store_id = s.id", store.TableName())).
		Where("products.id = ? AND products.store_id = ?", productID, storeID).
		First(&ps).Error; err != nil {

		return nil, err
	}

	poc := models.ProductOfCollection{}
	c := models.Collection{}
	var collections []models.CollectionDetails
	if err := db.Table(fmt.Sprintf("%s AS poc", poc.TableName())).
		Select("c.id, c.name, c.description").
		Joins(fmt.Sprintf("JOIN %s AS c ON poc.collection_id = c.id", c.TableName())).
		Where("poc.store_id = ? AND poc.product_id = ?", storeID, productID).
		Scan(&collections).Error; err != nil {
		return nil, err
	}

	acp := models.AdditionalChargeOfProduct{}
	ac := models.AdditionalCharge{}
	var additionalCharges []models.AdditionalChargeDetails
	if err := db.Table(fmt.Sprintf("%s AS acp", acp.TableName())).
		Select("ac.id, ac.name, ac.charge_type, ac.amount, ac.amount_type, ac.amount_max, ac.amount_min").
		Joins(fmt.Sprintf("JOIN %s AS ac ON acp.additional_charge_id = ac.id", ac.TableName())).
		Where("acp.product_id = ?", productID).
		Scan(&additionalCharges).Error; err != nil {
		return nil, err
	}

	ps.Collections = collections
	ps.AdditionalCharges = additionalCharges
	return &ps, nil
}

func (pu *ProductRepositoryImpl) GetProductForOrder(storeID, productID string, quantity int) (*models.Product, error) {
	tx := app.DB().Begin()
	p := models.Product{}

	if err := tx.Table(p.TableName()).
		Where("products.id = ? AND products.store_id = ? AND products.quantity - ? >= 0", productID, storeID, quantity).
		First(&p).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Table(p.TableName()).
		Where("products.id = ? AND products.store_id = ? AND products.quantity - ? >= 0", productID, storeID, quantity).
		UpdateColumn("products.quantity = products.quantity - ?", quantity).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &p, nil
}
