package models

type Customer struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	ProfilePicture    string `json:"profile_picture"`
	Phone             string `json:"phone"`
	IsEmailVerified   bool
	IsPhoneVerified   bool
	StoreID           string `json:"store_id"`
	NumberOfPurchases int    `json:"number_of_purchases"`
}
