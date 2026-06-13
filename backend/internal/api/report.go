package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"qfis/internal/db"
	"qfis/internal/engine"
	"qfis/internal/models"

	"github.com/gin-gonic/gin"
)

func SubmitReport(c *gin.Context) {
	var req models.ReportSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	merchant, _ := db.GetMerchant(req.MerchantID)
	merchantName := req.MerchantID
	if merchant != nil {
		merchantName = merchant.Name
	}

	score, _ := engine.CalculateRisk(req.MerchantID, merchantName, "")
	if merchant != nil {
		_ = merchant
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	reportID := fmt.Sprintf("RPT-%04d", rng.Intn(10000))

	report := &models.Report{
		ID:         reportID,
		MerchantID: req.MerchantID,
		MerchantName: merchantName,
		Note:       req.ReporterNote,
		Score:      score.RiskScore,
		Status:     "queued",
		CreatedAt:  time.Now().Unix(),
	}

	if err := db.InsertReport(report); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auto-flag if high risk
	if score.RiskScore >= 70 && merchant != nil {
		merchant.Flagged = true
		db.UpsertMerchant(merchant)
	}

	c.JSON(http.StatusOK, models.ReportSubmitResponse{
		ReportID: reportID,
		Status:   "queued",
	})
}

func ListReports(c *gin.Context) {
	reports, err := db.GetReports(50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if reports == nil {
		reports = []models.Report{}
	}
	c.JSON(http.StatusOK, reports)
}

func GetReport(c *gin.Context) {
	id := c.Param("id")
	reports, err := db.GetReports(100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, r := range reports {
		if r.ID == id {
			c.JSON(http.StatusOK, r)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
}
