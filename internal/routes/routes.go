package routes

import (
	"aevum-emporium-be/internal/controllers"
	"aevum-emporium-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Auth Routes
	authGroup := router.Group("/auth/user")
	{
		authGroup.POST("/signup", controllers.SignUp())
		authGroup.POST("/login", controllers.Login())
	}

	// Product Routes
	productGroup := router.Group("/product")
	{
		productGroup.POST("/add", middleware.AuthMiddleware(), controllers.AddProduct())
		productGroup.PUT("/:product_id", middleware.AuthMiddleware(), controllers.UpdateProduct())
		productGroup.DELETE("/:product_id", middleware.AuthMiddleware(), controllers.DeleteProduct())

		productGroup.GET("/", controllers.GetProducts())
		productGroup.GET("/:product_id", controllers.GetProductByID())
		productGroup.GET("/search", controllers.SearchProductByCategoryOrName())
	}

	// Order Routes
	orderGroup := router.Group("/orders")
	{
		orderGroup.POST("/place", middleware.AuthMiddleware(), controllers.PlaceOrder())
		orderGroup.GET("/", middleware.AuthMiddleware(), controllers.GetOrders())
		orderGroup.PUT("/:order_id/status", middleware.AuthMiddleware(), controllers.UpdateOrder())
		orderGroup.DELETE("/:order_id", middleware.AuthMiddleware(), controllers.CancelOrder())
	}

	// Cart Routes
	cartGroup := router.Group("/cart")
	{
		cartGroup.POST("/", middleware.AuthMiddleware(), controllers.AddToCart())
		cartGroup.GET("/", middleware.AuthMiddleware(), controllers.ViewCart())
		cartGroup.DELETE("/:product_id", middleware.AuthMiddleware(), controllers.RemoveFromCart())
		cartGroup.DELETE("/clear", middleware.AuthMiddleware(), controllers.ClearCart())

	}

	// Address Routes
	addressGroup := router.Group("/address")
	{
		addressGroup.POST("/add", middleware.AuthMiddleware(), controllers.AddAddress())
		addressGroup.PUT("/home", middleware.AuthMiddleware(), controllers.EditHomeAddress())
		addressGroup.PUT("/work", middleware.AuthMiddleware(), controllers.EditWorkAddress())
		addressGroup.DELETE("/delete", middleware.AuthMiddleware(), controllers.DeleteAddress())
	}

	// Wishlist Routes
	wishlistGroup := router.Group("/wishlist")
	{
		wishlistGroup.POST("/", middleware.AuthMiddleware(), controllers.AddWishlist())
		wishlistGroup.GET("/:user_id", middleware.AuthMiddleware(), controllers.GetWishlistByUser())
		wishlistGroup.DELETE("/:wishlist_id", middleware.AuthMiddleware(), controllers.DeleteWishlist())
	}

	// Review Routes
	reviewGroup := router.Group("/review")
	{
		reviewGroup.POST("/add", middleware.AuthMiddleware(), controllers.AddReview())
		reviewGroup.GET("/:product_id", controllers.GetReviewsByProduct())
		reviewGroup.DELETE("/:review_id", middleware.AuthMiddleware(), controllers.DeleteReview())
	}

}
