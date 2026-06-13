package engine

import "testing"

func TestClassifyMerchant_Suspicious(t *testing.T) {
	cases := []struct {
		name      string
		minScore  float64
	}{
		{"Pulsa Murah 24Jam", 0.4},
		{"Token Digital Express", 0.3},
		{"Chip Poker Online", 0.3},
		{"Depo Gaming Cepat", 0.3},
	}

	for _, c := range cases {
		score := ClassifyMerchant(c.name)
		if score < c.minScore {
			t.Errorf("%s: expected score >= %f, got %f", c.name, c.minScore, score)
		}
	}
}

func TestClassifyMerchant_Legit(t *testing.T) {
	cases := []struct {
		name     string
		maxScore float64
	}{
		{"Warung Kopi Barokah", 0.3},
		{"Alfamart Gontor", 0.3},
		{"Bakso Mas Joko", 0.3},
		{"Toko Sembako Ibu", 0.3},
	}

	for _, c := range cases {
		score := ClassifyMerchant(c.name)
		if score > c.maxScore {
			t.Errorf("%s: expected score <= %f, got %f", c.name, c.maxScore, score)
		}
	}
}

func TestClassifyMerchant_EdgeCases(t *testing.T) {
	cases := []struct {
		name      string
		minScore  float64
		maxScore  float64
	}{
		{"", 0, 0.5},           // empty
		{"a", 0.1, 0.6},        // single word = suspicious pattern
		{"Toko", 0, 0.3},       // single legit-ish word
	}

	for _, c := range cases {
		score := ClassifyMerchant(c.name)
		if score < c.minScore || score > c.maxScore {
			t.Errorf("%s: expected score in [%f,%f], got %f", c.name, c.minScore, c.maxScore, score)
		}
	}
}

func TestClassifyMerchant_Bounds(t *testing.T) {
	score := ClassifyMerchant("Pulsa Murah 24Jam Express Instan Digital Pro VIP")
	if score > 1.0 {
		t.Errorf("expected score <= 1.0, got %f", score)
	}
	if score < 0 {
		t.Errorf("expected score >= 0, got %f", score)
	}
}
