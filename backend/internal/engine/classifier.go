package engine

import (
	"math"
	"strings"
)

var suspiciousPatterns = []string{
	"24", "jam", "murah", "cepat", "instan", "express",
	"online", "digital", "pro", "vip", "premium",
}

var legitimizingWords = []string{
	"barokah", "sejahtera", "makmur", "jaya", "mantap",
	"juara", "fresh", "enak", "indonesia",
}

func ClassifyMerchant(name string) float64 {
	nameLower := strings.ToLower(name)
	score := 0.0

	keywordScore, _ := CheckMerchantName(name)
	score += keywordScore * 0.5

	for _, p := range suspiciousPatterns {
		if strings.Contains(nameLower, p) {
			score += 0.1
		}
	}

	for _, w := range legitimizingWords {
		if strings.Contains(nameLower, w) {
			score = math.Max(0, score-0.15)
		}
	}

	words := strings.Fields(nameLower)
	if len(words) <= 1 {
		score += 0.2
	}
	if len(words) > 5 {
		score += 0.1
	}

	return math.Min(score, 1.0)
}
