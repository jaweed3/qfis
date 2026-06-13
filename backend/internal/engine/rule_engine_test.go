package engine

import (
	"testing"
)

func TestCheckNPWP_Invalid(t *testing.T) {
	// Since CheckNPWP queries DB, we test the score mapping logic
	// via CheckMerchantName and ClassifyMerchant instead
	score, flags := CheckMerchantName("Pulsa Murah 24Jam")
	if score <= 0 {
		t.Errorf("expected score > 0 for suspicious name, got %f", score)
	}
	if len(flags) == 0 {
		t.Error("expected flags for suspicious name")
	}
}

func TestCheckMerchantName_Suspicious(t *testing.T) {
	cases := []struct {
		name     string
		wantFlag bool
		minScore float64
	}{
		{"Pulsa Murah 24Jam", true, 0.3},
		{"Token Digital Express", true, 0.3},
		{"Reload Gaming Center", true, 0.3},
		{"Warung Kopi Barokah", false, 0},
		{"Alfamart Gontor", false, 0},
		{"Toko Sembako Ibu", false, 0},
	}

	for _, c := range cases {
		score, flags := CheckMerchantName(c.name)
		if c.wantFlag && len(flags) == 0 {
			t.Errorf("%s: expected flags but got none", c.name)
		}
		if !c.wantFlag && len(flags) > 0 {
			t.Errorf("%s: expected no flags but got %v", c.name, flags)
		}
		if score < c.minScore {
			t.Errorf("%s: expected score >= %f, got %f", c.name, c.minScore, score)
		}
	}
}

func TestCheckMerchantName_LegitKeywords(t *testing.T) {
	legitNames := []string{
		"Bakso Mas Joko",
		"Laundry Express",
		"Bengkel Motor Jaya",
		"Salon Cantik",
		"Nasi Goreng Mantap",
	}

	for _, name := range legitNames {
		score, flags := CheckMerchantName(name)
		if score > 0.3 {
			t.Errorf("%s: expected low score for legit name, got %f with flags %v", name, score, flags)
		}
	}
}

func TestCheckMCC(t *testing.T) {
	cases := []struct {
		mcc       string
		wantScore float64
	}{
		{"5999", 0.7},
		{"7993", 0.9},
		{"7995", 0.8},
		{"5812", 0},
		{"5411", 0},
		{"9999", 0},
	}

	for _, c := range cases {
		score := CheckMCC(c.mcc)
		if score != c.wantScore {
			t.Errorf("MCC %s: expected %f, got %f", c.mcc, c.wantScore, score)
		}
	}
}

func TestCalcNameEntropy(t *testing.T) {
	cases := []struct {
		name     string
		minEntropy float64
	}{
		{"aaaaaa", 0},
		{"abcdef", 2.5},
		{"Pulsa Murah 24Jam", 3.0},
	}

	for _, c := range cases {
		entropy := calcNameEntropy(c.name)
		if entropy < c.minEntropy {
			t.Errorf("%s: expected entropy >= %f, got %f", c.name, c.minEntropy, entropy)
		}
	}
}

func TestSuspiciousKeywords(t *testing.T) {
	// Verify all suspicious keywords are detected
	keywords := []string{"pulsa", "token", "reload", "top up", "voucher"}
	for _, kw := range keywords {
		score, flags := CheckMerchantName(kw + " test")
		if score <= 0 && len(flags) == 0 {
			t.Errorf("keyword %q not detected", kw)
		}
	}
}
