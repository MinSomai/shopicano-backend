package api

import (
	"fmt"
	"gopkg.in/h2non/baloo.v3"
	"net/http"
	"testing"
)

var url = "http://localhost:9119"
var token = "ODdjNmQzMTMtOTM0YS00YzRlLThhZmEtNTZhMmMzZjE2ZjgzXzE1NzYzNDYyODRfMjAxOS0xMi0xNCAxNzo1ODowNC42MzI0MSArMDAwMCBVVEM="

//{
//"name": "NMD_R1 SHOES",
//"description": "NMD_R1 SHOES A MODERN NMD TRAINER WITH A SNUG KNIT UPPER. Run with it. These adidas NMD_R1 Shoes are a little technical and a lot street smart. Their midsole plugs flash back to the '80s, but the knit upper, full-length cushioned midsole and EVA inserts are 100 percent fashion forward.",
//"is_published": true,
//"is_shippable": true,
//"is_digital": false,
//"price": 130,
//"sku": "RS-R101",
//"stock": 10,
//"unit": "item",
//"image": "products/d90c343e-6780-415e-a49a-773675fd79eb-c2hvZV9hZGRpZGFz.webp",
//"additional_images": [
//"products/d90c343e-6780-415e-a49a-773675fd79eb-c2hvZV9hZGRpZGFz.webp",
//"products/d90c343e-6780-415e-a49a-773675fd79eb-c2hvZV9hZGRpZGFz.webp"
//],
//"category_id": "9fd756a9-cb01-439c-9e16-8c5de7d0365c"
//}

func TestCreateProduct(t *testing.T) {
	pld := struct {
		Name             string   `json:"name"`
		Description      string   `json:"description"`
		IsPublished      bool     `json:"is_published"`
		IsShippable      bool     `json:"is_shippable"`
		IsDigital        bool     `json:"is_digital"`
		Price            int      `json:"price"`
		SKU              string   `json:"sku"`
		Stock            string   `json:"stock"`
		Unit             string   `json:"unit"`
		Image            string   `json:"image"`
		AdditionalImages []string `json:"additional_images"`
		CategoryID       string   `json:"category_id"`
	}{}

	test := baloo.New(url)
	test.Post("/v1/products").
		JSON(pld).
		AddHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Expect(t).
		Status(http.StatusCreated).
		End()
}
