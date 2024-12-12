package routes

import (
	"aevum-emporium-be/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Auth Routes
	authGroup := router.Group("/auth/users")
	{
		authGroup.POST("/signup", controllers.SignUp())
		authGroup.POST("/login", controllers.Login())
	}

	// Product Routes
	productGroup := router.Group("/product")
	{
		productGroup.POST("/add", controllers.AddProduct())
		productGroup.GET("/", controllers.GetProducts())
		productGroup.GET("/:id", controllers.GetProductByID())
		productGroup.GET("/search", controllers.SearchProductByCategoryOrName())
		productGroup.PUT("/:id", controllers.UpdateProduct())
		productGroup.DELETE("/:id", controllers.DeleteProduct())
	}

	// Order Routes
	orderGroup := router.Group("/order")
	{
		orderGroup.POST("/add", controllers.PlaceOrder())
		orderGroup.GET("/", controllers.GetOrders())
	}

	// Cart Routes
	cartGroup := router.Group("/cart")
	{
		cartGroup.POST("/add", controllers.AddToCart())
		cartGroup.GET("/", controllers.ViewCart())
		cartGroup.POST("/remove", controllers.RemoveFromCart())
		cartGroup.DELETE("/clear", controllers.ClearCart())
	}

	// Address Routes
	addressGroup := router.Group("/address")
	{
		addressGroup.POST("/add", controllers.AddAddress())
		addressGroup.PUT("/edit", controllers.EditHomeAddress())
		addressGroup.PUT("/edit", controllers.EditWorkAddress())
		addressGroup.DELETE("/:id", controllers.DeleteAddress())
	}
}
