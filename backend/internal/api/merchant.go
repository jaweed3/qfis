package api

import (
	"net/http"

	"qfis/internal/db"
	"qfis/internal/engine"
	"qfis/internal/models"

	"github.com/gin-gonic/gin"
)

func CheckMerchant(c *gin.Context) {
	var req models.CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	score, err := engine.CalculateRisk(req.MerchantID, req.MerchantName, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, score)
}

func SearchMerchant(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	merchants, err := db.SearchMerchants(q)
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

func GetMerchantByID(c *gin.Context) {
	id := c.Param("id")
	merchant, err := db.GetMerchant(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "merchant not found"})
		return
	}

	resp, _ := engine.CalculateRisk(merchant.ID, merchant.Name, merchant.MCC)
	merchant.RiskScore = resp.RiskScore
	merchant.RiskLabel = resp.RiskLabel

	c.JSON(http.StatusOK, merchant)
}
