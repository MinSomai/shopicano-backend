package helpers

type ProductStats struct {
	ID            string `json:"id" sql:"id"`
	Name          string `json:"name" sql:"name"`
	Stock         int    `json:"stock" sql:"stock"`
	Price         int    `json:"price" sql:"price"`
	Image         string `json:"image" sql:"image"`
	Description   string `json:"description" sql:"description"`
	NumberOfSells int    `json:"number_of_sells" sql:"number_of_sells"`
}

type CategoryStats struct {
	ID            string `json:"id" sql:"id"`
	Name          string `json:"name" sql:"name"`
	Image         string `json:"image" sql:"image"`
	Description   string `json:"description" sql:"description"`
	NumberOfSells int    `json:"number_of_sells" sql:"number_of_sells"`
}
