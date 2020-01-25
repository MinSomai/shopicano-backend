package data

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/helpers"
	"github.com/shopicano/shopicano-backend/models"
	"strings"
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

func (pu *ProductRepositoryImpl) Create(db *gorm.DB, p *models.Product) error {
	if err := db.Model(p).Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (pu *ProductRepositoryImpl) Update(db *gorm.DB, p *models.Product) error {
	if err := db.Table(p.TableName()).
		Select("name, description, is_published, category_id, sku, stock, unit, price, additional_images, image, is_shippable, is_digital, digital_download_link, updated_at").
		Where("id = ? AND store_id = ?", p.ID, p.StoreID).
		Updates(map[string]interface{}{
			"name":                  p.Name,
			"description":           p.Description,
			"is_published":          p.IsPublished,
			"category_id":           p.CategoryID,
			"sku":                   p.SKU,
			"stock":                 p.Stock,
			"unit":                  p.Unit,
			"price":                 p.Price,
			"additional_images":     p.AdditionalImages,
			"image":                 p.Image,
			"is_shippable":          p.IsShippable,
			"is_digital":            p.IsDigital,
			"digital_download_link": p.DigitalDownloadLink,
			"updated_at":            p.UpdatedAt,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (pu *ProductRepositoryImpl) IncreaseDownloadCounter(db *gorm.DB, p *models.Product) error {
	if err := db.Table(p.TableName()).
		Where("id = ? AND store_id = ?", p.ID, p.StoreID).
		Update("download_counter", gorm.Expr("download_counter + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func (pu *ProductRepositoryImpl) List(db *gorm.DB, from, limit int) ([]models.ProductDetails, error) {
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

func (pu *ProductRepositoryImpl) ListAsStoreStuff(db *gorm.DB, storeID string, from, limit int) ([]models.ProductDetailsInternal, error) {
	var ps []models.ProductDetailsInternal
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

func (pu *ProductRepositoryImpl) Search(db *gorm.DB, query string, from, limit int) ([]models.ProductDetails, error) {
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

func (pu *ProductRepositoryImpl) SearchAsStoreStuff(db *gorm.DB, storeID, query string, from, limit int) ([]models.ProductDetailsInternal, error) {
	var ps []models.ProductDetailsInternal
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

func (pu *ProductRepositoryImpl) Delete(db *gorm.DB, storeID, productID string) error {
	p := models.Product{}
	if err := db.Table(p.TableName()).
		Where("store_id = ? AND id = ?", storeID, productID).
		Delete(&p).Error; err != nil {
		return err
	}
	return nil
}

func (pu *ProductRepositoryImpl) Get(db *gorm.DB, productID string) (*models.Product, error) {
	p := models.Product{}
	if err := db.Table(fmt.Sprintf("%s", p.TableName())).
		Where("products.id = ?", productID).
		First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (pu *ProductRepositoryImpl) GetAsStoreStuff(db *gorm.DB, storeID, productID string) (*models.Product, error) {
	p := models.Product{}
	if err := db.Table(fmt.Sprintf("%s", p.TableName())).
		Where("products.id = ? AND products.store_id = ?", productID, storeID).
		First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (pu *ProductRepositoryImpl) GetDetails(db *gorm.DB, productID string) (*models.ProductDetails, error) {
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

	c := models.Collection{}
	cop := models.CollectionOfProduct{}
	var collections []models.Collection
	if err := db.Table(fmt.Sprintf("%s AS cop", cop.TableName())).
		Select("c.id, c.name, c.description").
		Joins(fmt.Sprintf("JOIN %s AS c ON cop.collection_id = c.id", c.TableName())).
		Where("cop.product_id = ?", productID).
		Scan(&collections).Error; err != nil {
		return nil, err
	}

	a := models.ProductAttribute{}
	var attributes []models.ProductAttribute
	if err := db.Table(fmt.Sprintf("%s AS pa", a.TableName())).
		Select("pa.key AS key, pa.value AS value").
		Where("pa.product_id = ?", productID).
		Scan(&attributes).Error; err != nil {
		return nil, err
	}

	ps.Collections = collections
	ps.Attributes = attributes
	return &ps, nil
}

func (pu *ProductRepositoryImpl) GetDetailsAsStoreStuff(db *gorm.DB, storeID, productID string) (*models.ProductDetailsInternal, error) {
	p := models.Product{}
	ps := models.ProductDetailsInternal{}
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

	cop := models.CollectionOfProduct{}
	c := models.Collection{}
	var collections []models.Collection
	if err := db.Table(fmt.Sprintf("%s AS cop", cop.TableName())).
		Select("c.id, c.name, c.description").
		Joins(fmt.Sprintf("JOIN %s AS c ON cop.collection_id = c.id", c.TableName())).
		Where("c.store_id = ? AND cop.product_id = ?", storeID, productID).
		Scan(&collections).Error; err != nil {
		return nil, err
	}

	a := models.ProductAttribute{}
	var attributes []models.ProductAttribute
	if err := db.Table(fmt.Sprintf("%s AS pa", a.TableName())).
		Select("pa.key AS key, pa.value AS value").
		Where("pa.product_id = ?", productID).
		Scan(&attributes).Error; err != nil {
		return nil, err
	}
	ps.Attributes = attributes

	ps.Collections = collections
	return &ps, nil
}

func (pu *ProductRepositoryImpl) GetForOrder(db *gorm.DB, productID string, quantity int) (*models.Product, error) {
	p := models.Product{}

	if err := db.Table(p.TableName()).
		Where("id = ? AND (stock - ? >= 0 OR is_digital)", productID, quantity).
		First(&p).Error; err != nil {
		return nil, err
	}

	if p.IsDigital {
		return &p, nil
	}

	if err := db.Table(p.TableName()).
		Where("id = ? AND stock - ? >= 0", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (pu *ProductRepositoryImpl) Stats(db *gorm.DB, from, limit int) ([]helpers.ProductStats, error) {
	var stats []helpers.ProductStats

	p := models.Product{}
	oi := models.OrderedItem{}

	if err := db.Table(fmt.Sprintf("%s AS p", p.TableName())).
		Select("p.id AS id, p.name AS name, p.stock AS stock, p.price AS price, p.image AS image, p.description AS description, COALESCE(SUM(oi.quantity), 0) AS number_of_sells").
		Joins(fmt.Sprintf("LEFT JOIN %s AS oi ON p.id = oi.product_id", oi.TableName())).
		Group("p.id, p.name, p.stock, p.price, p.image, p.description").
		Order("number_of_sells DESC").
		Offset(from).
		Limit(limit).
		Find(&stats).Error; err != nil {
		return nil, err
	}

	if stats == nil {
		stats = []helpers.ProductStats{}
	}

	return stats, nil
}

func (pu *ProductRepositoryImpl) StatsAsStoreStuff(db *gorm.DB, storeID string, from, limit int) ([]helpers.ProductStats, error) {
	var stats []helpers.ProductStats

	p := models.Product{}
	oi := models.OrderedItem{}

	if err := db.Table(fmt.Sprintf("%s AS p", p.TableName())).
		Select("p.id AS id, p.name AS name, p.stock AS stock, p.price AS price, p.image AS image, p.description AS description, COALESCE(SUM(oi.quantity), 0) AS number_of_sells").
		Joins(fmt.Sprintf("LEFT JOIN %s AS oi ON p.id = oi.product_id", oi.TableName())).
		Group("p.id, p.name, p.stock, p.price, p.image, p.description").
		Order("number_of_sells DESC").
		Offset(from).
		Limit(limit).
		Find(&stats, "p.store_id = ?", storeID).Error; err != nil {
		return nil, err
	}

	if stats == nil {
		stats = []helpers.ProductStats{}
	}

	return stats, nil
}

func (pu *ProductRepositoryImpl) AddAttribute(db *gorm.DB, v *models.ProductAttribute) error {
	if err := db.Table(v.TableName()).Create(v).Error; err != nil {
		return err
	}
	return nil
}

func (pu *ProductRepositoryImpl) RemoveAttribute(db *gorm.DB, productID, attributeKey string) error {
	v := models.ProductAttribute{}
	if err := db.Table(v.TableName()).Delete(&v, "product_id = ? AND key = ?", productID, attributeKey).Error; err != nil {
		return err
	}
	return nil
}
