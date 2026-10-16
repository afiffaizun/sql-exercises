package main

import (
	"log"
	"os"

	"crud-produk/config"
	"crud-produk/handlers"
	"crud-produk/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDatabase()

	if err := db.AutoMigrate(&models.Produk{}); err != nil {
		log.Fatalf("gagal migrasi database: %v", err)
	}
	config.Seed(db)

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(corsConfig))

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
		api.GET("/stats", handlers.GetStats)

		produk := api.Group("/produk")
		{
			produk.GET("", handlers.GetProduk)
			produk.GET("/:id", handlers.GetProdukByID)
			produk.POST("", handlers.CreateProduk)
			produk.PUT("/:id", handlers.UpdateProduk)
			produk.DELETE("/:id", handlers.DeleteProduk)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server berjalan di port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("gagal menjalankan server: %v", err)
	}
}
