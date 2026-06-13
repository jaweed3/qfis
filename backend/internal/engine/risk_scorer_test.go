package engine

import (
	"os"
	"testing"

	"qfis/internal/db"
)

func TestCalculateRisk_ScoreRange(t *testing.T) {
	// Setup in-memory DB
	tmpDB := os.TempDir() + "qfis_test.db"
	defer os.Remove(tmpDB)

	if err := db.Init(tmpDB); err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer db.Close()

	if err := db.Seed(); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	cases := []struct {
		id   string
		name string
		mcc  string
		min  float64
		max  float64
	}{
		{"QRIS-9999", "Pulsa Murah 24Jam", "5999", 0, 100},
		{"QRIS-9998", "Warung Kopi Barokah", "5812", 0, 100},
		{"QRIS-9997", "", "", 0, 100},
	}

	for _, c := range cases {
		resp, err := CalculateRisk(c.id, c.name, c.mcc)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", c.name, err)
			continue
		}
		if resp.RiskScore < c.min || resp.RiskScore > c.max {
			t.Errorf("%s: score %f out of range [%f,%f]", c.name, resp.RiskScore, c.min, c.max)
		}
		if resp.RiskLabel == "" {
			t.Errorf("%s: missing risk label", c.name)
		}
		if resp.Recommendation == "" {
			t.Errorf("%s: missing recommendation", c.name)
		}
	}
}

func TestCalculateRisk_SuspiciousVsLegit(t *testing.T) {
	tmpDB := os.TempDir() + "qfis_test2.db"
	defer os.Remove(tmpDB)

	db.Init(tmpDB)
	defer db.Close()
	db.Seed()

	suspiciousResp, _ := CalculateRisk("TEST-001", "Pulsa Murah 24 Jam Toko Chip Poker", "7993")
	legitResp, _ := CalculateRisk("TEST-002", "Warung Kopi Barokah Sejahtera", "5812")

	if suspiciousResp.RiskScore <= legitResp.RiskScore {
		t.Errorf("expected suspicious score (%f) > legit score (%f)",
			suspiciousResp.RiskScore, legitResp.RiskScore)
	}
}

func TestCalculateRisk_Flags(t *testing.T) {
	tmpDB := os.TempDir() + "qfis_test3.db"
	defer os.Remove(tmpDB)

	db.Init(tmpDB)
	defer db.Close()
	db.Seed()

	resp, _ := CalculateRisk("TEST-FLAG", "Chip Poker Slot Online", "7993")
	if len(resp.Flags) == 0 {
		t.Error("expected flags for suspicious merchant")
	}
}
