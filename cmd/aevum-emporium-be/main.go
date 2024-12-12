package main

import (
	"aevum-emporium-be/internal/controllers"
	"aevum-emporium-be/internal/datasource"
	"aevum-emporium-be/internal/routes"

	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app := controllers.NewApplication(datasource.ProductData(datasource.Client, "Products"), datasource.UserData(datasource.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	// router.Use(middleware.Authentication())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	// router.GET("/listcart", controllers.GetItemFromCart())
	// router.POST("/addaddress", controllers.AddAddress())
	// router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	// router.PUT("/editworkaddress", controllers.EditWorkAddress())
	// router.GET("/deleteaddresses", controllers.DeleteAddress())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())
	log.Fatal(router.Run(":" + port))
}
