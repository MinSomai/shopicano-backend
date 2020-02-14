package models

type Summary struct {
	Time        string `json:"time"`
	TotalOrders int    `json:"total_orders"`
	Earnings    int    `json:"earnings"`
	Expenses    int    `json:"expenses"`
	Profits     int    `json:"profits"`
	Discounts   int    `json:"discounts"`
	Customers   int    `json:"customers"`
}
