package engine

import "testing"

func TestDetectAnomaly(t *testing.T) {
	cases := []struct {
		name      string
		reports   int
		npwpValid bool
		minFlags  int
		minScore  float64
	}{
		{"Pulsa Murah 24Jam", 15, false, 1, 0.4},
		{"Chip Poker Online", 5, false, 1, 0.3},
		{"Warung Kopi Barokah", 0, true, 0, 0},
		{"Toko Sembako Ibu", 0, true, 0, 0},
	}

	for _, c := range cases {
		score, flags := DetectAnomaly(c.name, c.reports, c.npwpValid)
		if score < c.minScore {
			t.Errorf("%s: expected score >= %f, got %f with %d flags", c.name, c.minScore, score, len(flags))
		}
		if len(flags) < c.minFlags {
			t.Errorf("%s: expected >= %d flags, got %d", c.name, c.minFlags, len(flags))
		}
	}
}

func TestDetectAnomaly_GamblingTerms(t *testing.T) {
	gamblingNames := []string{
		"Poker Online Terpercaya",
		"Slot Gacor Maxwin",
		"Casino Online",
		"Togel Singapore",
		"Jackpot Besar",
	}

	for _, name := range gamblingNames {
		score, flags := DetectAnomaly(name, 0, false)
		if len(flags) == 0 {
			t.Errorf("%s: expected gambling term flag", name)
		}
		if score <= 0 {
			t.Errorf("%s: expected score > 0, got %f", name, score)
		}
	}
}

func TestDetectAnomaly_DigitRatio(t *testing.T) {
	score, flags := DetectAnomaly("T0k3n D1g1t4l 12345", 0, false)
	if score <= 0 {
		t.Errorf("expected score > 0 for name with many digits, got %f", score)
	}
	_ = flags
}

func TestDetectAnomaly_ScoreBounds(t *testing.T) {
	score, _ := DetectAnomaly("Normal Merchant Name", 0, true)
	if score > 0.3 {
		t.Errorf("expected low score for normal merchant, got %f", score)
	}

	score, _ = DetectAnomaly("Chip Poker Slot Togel Casino", 100, false)
	if score > 1.0 {
		t.Errorf("expected score capped at 1.0, got %f", score)
	}
	if score <= 0 {
		t.Errorf("expected score > 0 for high-risk merchant, got %f", score)
	}
}
