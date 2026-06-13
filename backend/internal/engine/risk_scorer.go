package engine

import (
	"fmt"
	"math"

	"qfis/internal/db"
	"qfis/internal/models"
)

func CalculateRisk(merchantID, merchantName, mcc string) (*models.CheckResponse, error) {
	var err error
	resp := &models.CheckResponse{
		Flags: []string{},
	}

	merchant, dbErr := db.GetMerchant(merchantID)
	if dbErr == nil {
		merchantName = merchant.Name
		mcc = merchant.MCC
	}

	npwpScore, _ := CheckNPWP(merchantID)
	if npwpScore > 0 {
		resp.Flags = append(resp.Flags, "NPWP tidak terdaftar atau tidak valid")
	}

	_, nameFlags := CheckMerchantName(merchantName)
	resp.Flags = append(resp.Flags, nameFlags...)

	mccScore := CheckMCC(mcc)
	if mccScore > 0 {
		resp.Flags = append(resp.Flags, fmt.Sprintf("MCC code %s termasuk kategori berisiko tinggi", mcc))
	}

	mlScore := ClassifyMerchant(merchantName)

	reports := 0
	npwpValid := false
	if merchant != nil {
		npwpValid = merchant.NPWP
	}
	reportsList, _ := db.GetReports(100)
	for _, r := range reportsList {
		if r.MerchantID == merchantID {
			reports++
		}
	}

	anomalyScore, anomalyFlags := DetectAnomaly(merchantName, reports, npwpValid)
	resp.Flags = append(resp.Flags, anomalyFlags...)

	var communityScore float64
	if reports > 0 {
		communityScore = math.Min(float64(reports)*0.05, 1.0)
	}

	finalScore := (0.30*npwpScore +
		0.25*mlScore +
		0.20*anomalyScore +
		0.15*communityScore +
		0.10*mccScore) * 100

	finalScore = math.Min(math.Max(finalScore, 0), 100)

	resp.RiskScore = math.Round(finalScore*10) / 10

	switch {
	case resp.RiskScore >= 70:
		resp.RiskLabel = "MERAH"
		resp.Recommendation = "Segera laporkan ke OJK untuk investigasi lebih lanjut"
	case resp.RiskScore >= 30:
		resp.RiskLabel = "KUNING"
		resp.Recommendation = "Masukkan dalam watch list, pantau aktivitas selanjutnya"
	default:
		resp.RiskLabel = "HIJAU"
		resp.Recommendation = "Merchant terlihat legitimate, tidak ada indikasi mencurigakan"
	}

	return resp, err
}
