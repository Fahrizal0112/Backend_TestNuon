package main

import (
	"log"
	"os"
	"transaction-api/config"
	"transaction-api/models"
	"transaction-api/routes"
	"transaction-api/utils"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	log.Println("Connecting to database...")
	config.ConnectDatabase()


	log.Println("Running auto migration...")
	if err := config.DB.AutoMigrate(&models.Transaction{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed")

	csvFile := "data.csv"
	if _, err := os.Stat(csvFile); err == nil {
		log.Printf("Found CSV file: %s, loading data...", csvFile)
		if err := utils.LoadCSVData(csvFile); err != nil {
			log.Printf("Error loading CSV data: %v", err)
		} else {
			log.Println("CSV data loaded successfully")
		}
	} else {
		log.Printf("CSV file not found: %s", csvFile)
	}
	
	var count int64
	config.DB.Model(&models.Transaction{}).Count(&count)
	log.Printf("Total transactions in database: %d", count)

	r := routes.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
