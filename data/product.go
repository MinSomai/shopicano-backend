package data

import (
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/validators"
)

type ProductRepository interface {
	CreateProduct(req *validators.ReqProductCreate) (*models.Product, error)
	UpdateProduct(productID string, req *validators.ReqProductUpdate) (*models.Product, error)
	ListProducts(from, limit int) ([]models.ProductDetails, error)
	SearchProducts(query string, from, limit int) ([]models.ProductDetails, error)
	ListProductsWithStore(storeID string, from, limit int) ([]models.ProductDetails, error)
	SearchProductsWithStore(storeID, query string, from, limit int) ([]models.ProductDetails, error)
	DeleteProduct(storeID, productID string) error
	GetProduct(productID string) (*models.ProductDetails, error)
	GetProductWithStore(storeID, productID string) (*models.ProductDetails, error)
	GetProductForOrder(storeID, productID string, quantity int) (*models.Product, error)
}
