package engine

import (
	"math"
	"strings"

	"qfis/internal/db"
)

var suspiciousKeywords = []string{
	"pulsa", "token", "reload", "top up", "topup", "voucher",
	"gaming", "game", "chip", "poker", "depo", "saldo",
	"refill", "isi ulang", "mobile legend", "ml", "free fire",
	"chip", "slot", "jackpot", "casino", "togel",
}

var legitKeywords = []string{
	"warung", "toko", "bakso", "nasi", "ayam", "soto",
	"pecel", "bengkel", "laundry", "salon", "minimarket",
	"martabak", "foto", "buku", "elektronik", "sembako",
	"kopi", "teh", "roti", "cemilan", "tahu",
}

func CheckNPWP(merchantID string) (float64, error) {
	merchant, err := db.GetMerchant(merchantID)
	if err != nil {
		return 0.5, nil
	}

	if !merchant.NPWP {
		db.SaveRuleCache(merchantID, "npwp_check", 1.0)
		return 1.0, nil
	}
	db.SaveRuleCache(merchantID, "npwp_check", 0.0)
	return 0.0, nil
}

func CheckMerchantName(name string) (float64, []string) {
	nameLower := strings.ToLower(name)
	var flags []string

	maxKeywordScore := 0.0
	for _, kw := range suspiciousKeywords {
		if strings.Contains(nameLower, kw) {
			flags = append(flags, "Nama mengandung keyword mencurigakan: "+kw)
			keywordScore := 0.6 + (float64(len(kw)) / 30.0)
			if keywordScore > maxKeywordScore {
				maxKeywordScore = keywordScore
			}
		}
	}

	for _, kw := range legitKeywords {
		if strings.Contains(nameLower, kw) {
			if len(flags) > 0 {
				maxKeywordScore = math.Max(0, maxKeywordScore-0.3)
			} else {
				maxKeywordScore = math.Max(maxKeywordScore, 0.0)
			}
			break
		}
	}

	nameEntropy := calcNameEntropy(name)
	entropyScore := 0.0
	if nameEntropy > 4.0 {
		entropyScore = math.Min((nameEntropy-4.0)/2.0, 1.0)
		flags = append(flags, "Nama merchant memiliki entropy tinggi (tidak wajar)")
	}

	finalScore := math.Max(maxKeywordScore, entropyScore)
	return math.Min(finalScore, 1.0), flags
}

func calcNameEntropy(name string) float64 {
	freq := make(map[rune]int)
	for _, c := range name {
		freq[c]++
	}
	entropy := 0.0
	length := len(name)
	for _, count := range freq {
		if count == 0 {
			continue
		}
		p := float64(count) / float64(length)
		entropy -= p * math.Log2(p)
	}
	return entropy
}

func CheckMCC(mcc string) float64 {
	highRiskMCCs := map[string]float64{
		"5999": 0.7, // Miscellaneous & Specialty Retail
		"7993": 0.9, // Gambling Transactions
		"7995": 0.8, // Betting
		"7801": 0.6, // Government Lottery
		"7802": 0.6, // Government Lottery
	}
	score, exists := highRiskMCCs[mcc]
	if !exists {
		return 0.0
	}
	return score
}
