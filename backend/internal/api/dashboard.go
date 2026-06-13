package api

import (
	"net/http"

	"qfis/internal/db"
	"qfis/internal/engine"
	"qfis/internal/models"

	"github.com/gin-gonic/gin"
)

func GetDashboardStats(c *gin.Context) {
	stats, err := db.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func GetMerchantNetwork(c *gin.Context) {
	merchants, err := db.GetAllMerchants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if merchants == nil {
		merchants = []models.Merchant{}
	}

	for i, m := range merchants {
		resp, _ := engine.CalculateRisk(m.ID, m.Name, m.MCC)
		merchants[i].RiskScore = resp.RiskScore
		merchants[i].RiskLabel = resp.RiskLabel
	}

	c.JSON(http.StatusOK, merchants)
}
