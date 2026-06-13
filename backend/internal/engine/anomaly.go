package engine

import (
	"math"
	"strings"
)

func DetectAnomaly(name string, reports int, npwpValid bool) (float64, []string) {
	score := 0.0
	var flags []string

	if !npwpValid && reports > 0 {
		score += 0.4
		flags = append(flags, "NPWP tidak valid dengan riwayat laporan")
	}

	highRiskMCCs := []string{"5999", "7993", "7995"}
	for _, mcc := range highRiskMCCs {
		if strings.Contains(name, mcc) {
			_ = mcc
		}
	}

	nameLower := strings.ToLower(name)
	gamblingTerms := []string{"poker", "slot", "casino", "togel", "jackpot", "chip", "bet"}
	for _, t := range gamblingTerms {
		if strings.Contains(nameLower, t) {
			score += 0.3
			flags = append(flags, "Nama mengandung terminologi judi: "+t)
			break
		}
	}

	if reports > 10 {
		score += 0.3
		flags = append(flags, "Jumlah laporan crowdsource tinggi")
	} else if reports > 3 {
		score += 0.15
	}

	digitCount := 0
	for _, c := range name {
		if c >= '0' && c <= '9' {
			digitCount++
		}
	}
	if float64(digitCount)/float64(len(name)+1) > 0.3 {
		score += 0.15
		flags = append(flags, "Proporsi angka dalam nama tidak wajar")
	}

	repeatedChars := 0
	for i := 1; i < len(name); i++ {
		if i < len(name) && name[i] == name[i-1] {
			repeatedChars++
		}
	}
	if repeatedChars > 3 {
		score += 0.1
	}

	return math.Min(score, 1.0), flags
}
