package routes

import (
	"transaction-api/controllers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.MaxMultipartMemory = 8 << 20

	v1 := r.Group("/api/v1")
	{
		v1.POST("/transactions", controllers.CreateTransaction)
		v1.GET("/transactions", controllers.GetAllTransactions)
		v1.GET("/transactions/:id", controllers.GetTransactionByID)

		v1.POST("/transactions/upload", controllers.UploadCSV)
		v1.DELETE("/transactions/clear", controllers.ClearAllTransactions)
	}
	return r
}