package main

import (
	"fmt"
	"log"
	"os"

	"qfis/internal/api"
	"qfis/internal/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	dbPath := "qfis.db"
	if path := os.Getenv("QFIS_DB_PATH"); path != "" {
		dbPath = path
	}

	if err := db.Init(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	if err := db.Seed(); err != nil {
		log.Printf("Warning: seed failed (may already be seeded): %v", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/merchant/check", api.CheckMerchant)
		v1.GET("/merchant/search", api.SearchMerchant)
		v1.GET("/merchant/:id", api.GetMerchantByID)

		v1.POST("/report/submit", api.SubmitReport)
		v1.GET("/report/list", api.ListReports)
		v1.GET("/report/:id", api.GetReport)

		v1.GET("/dashboard/stats", api.GetDashboardStats)
		v1.GET("/dashboard/network", api.GetMerchantNetwork)
	}

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("QFIS backend running on http://localhost:%s\n", port)
	r.Run(":" + port)
}
