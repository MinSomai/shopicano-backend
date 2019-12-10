package data

import "github.com/shopicano/shopicano-backend/models"

type CollectionRepository interface {
	CreateCollection(c *models.Collection) error
	ListCollections(from, limit int) ([]models.Collection, error)
	SearchCollections(query string, from, limit int) ([]models.Collection, error)
	ListCollectionsWithStore(storeID string, from, limit int) ([]models.Collection, error)
	SearchCollectionsWithStore(storeID, query string, from, limit int) ([]models.Collection, error)
	DeleteCollection(storeID, collectionID string) error
}
