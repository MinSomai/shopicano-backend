package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/middlewares"
)

func RegisterSaleRoutes(g *echo.Group) {
	func(g *echo.Group) {
		// Private endpoints only
		g.Use(middlewares.IsStoreStaffWithStoreActivation)
		g.POST("/", createProduct)
		g.PATCH("/:product_id/", updateProduct)
		g.DELETE("/:product_id/", deleteProduct)
		g.PUT("/:product_id/attributes/", addProductAttribute)
		g.DELETE("/:product_id/attributes/:attribute_key/", deleteProductAttribute)
	}(g)

	func(g *echo.Group) {
		// Private endpoints only
		g.Use(middlewares.AuthUser)
		g.POST("/", createProduct)
		g.PATCH("/:product_id/check", updateProduct)
		g.DELETE("/:product_id/", deleteProduct)
		g.PUT("/:product_id/attributes/", addProductAttribute)
		g.DELETE("/:product_id/attributes/:attribute_key/", deleteProductAttribute)
	}(g)
}
