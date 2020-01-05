package models

type Customer struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	ProfilePicture    string `json:"profile_picture"`
	Phone             string `json:"phone"`
	IsEmailVerified   bool   `json:"is_email_verified"`
	IsPhoneVerified   bool   `json:"is_phone_verified"`
	StoreID           string `json:"store_id"`
	NumberOfPurchases int    `json:"number_of_purchases"`
}
