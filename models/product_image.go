package models

import "fmt"

type ProductImage struct {
	ProductID string `json:"product_id" gorm:"column:product_id;primary_key"`
	ImagePath string `json:"image_path" gorm:"column:image_path;primary_key"`
}

func (pi *ProductImage) TableName() string {
	return "product_images"
}

func (pi *ProductImage) ForeignKeys() []string {
	p := Product{}

	return []string{
		fmt.Sprintf("product_id;%s(id);RESTRICT;RESTRICT", p.TableName()),
	}
}
