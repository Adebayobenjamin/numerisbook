package router

import (
	"time"

	controller_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/controllers/interfaces"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewApplicationRouter(
	uploadController controller_interfaces.InvoiceController,
) *gin.Engine {
	router := gin.Default()

	// Configure CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Authorization", "Content-Type", "Accept"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	router.Use(gin.Recovery())

	router.GET("/ping", PingHandler())
	// Group routes (api/v1/uploads
	apiRoutes := router.Group("/api/v1")
	NewInvoiceRouter(uploadController, apiRoutes)

	return router

}

func PingHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
