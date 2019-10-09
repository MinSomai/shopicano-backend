package repositories

import (
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/validators"
)

type OrderRepository interface {
	CreateOrder(v *validators.ReqOrderCreate) (*models.OrderDetails, error)
}
