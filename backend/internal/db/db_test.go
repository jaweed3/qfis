package db

import (
	"os"
	"testing"

	"qfis/internal/models"
)

func TestInit(t *testing.T) {
	tmpDB := os.TempDir() + "qfis_db_test.db"
	defer os.Remove(tmpDB)

	if err := Init(tmpDB); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Close()
}

func TestSeed(t *testing.T) {
	tmpDB := os.TempDir() + "qfis_seed_test.db"
	defer os.Remove(tmpDB)

	if err := Init(tmpDB); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Close()

	if err := Seed(); err != nil {
		t.Fatalf("Seed failed: %v", err)
	}

	merchants, err := GetAllMerchants()
	if err != nil {
		t.Fatalf("GetAllMerchants failed: %v", err)
	}
	if len(merchants) == 0 {
		t.Fatal("expected seeded merchants")
	}

	reports, err := GetReports(100)
	if err != nil {
		t.Fatalf("GetReports failed: %v", err)
	}
	if len(reports) == 0 {
		t.Fatal("expected seeded reports")
	}
}

func TestGetMerchant(t *testing.T) {
	tmpDB := os.TempDir() + "qfis_get_test.db"
	defer os.Remove(tmpDB)

	Init(tmpDB)
	defer Close()
	Seed()

	merchants, _ := GetAllMerchants()
	if len(merchants) == 0 {
		t.Fatal("no merchants seeded")
	}

	m, err := GetMerchant(merchants[0].ID)
	if err != nil {
		t.Fatalf("GetMerchant(%s) failed: %v", merchants[0].ID, err)
	}
	if m.ID != merchants[0].ID {
		t.Errorf("expected ID %s, got %s", merchants[0].ID, m.ID)
	}
	if m.Name == "" {
		t.Error("expected merchant name")
	}
}

func TestSearchMerchants(t *testing.T) {
	tmpDB := os.TempDir() + "qfis_search_test.db"
	defer os.Remove(tmpDB)

	Init(tmpDB)
	defer Close()
	Seed()

	results, err := SearchMerchants("pulsa")
	if err != nil {
		t.Fatalf("SearchMerchants failed: %v", err)
	}
	if len(results) == 0 {
		t.Error("expected at least 1 result for 'pulsa'")
	}

	results, err = SearchMerchants("ZZZZNONEXISTENT")
	if err != nil {
		t.Fatalf("SearchMerchants failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for nonexistent query, got %d", len(results))
	}
}

func TestDashboardStats(t *testing.T) {
	tmpDB := os.TempDir() + "qfis_stats_test.db"
	defer os.Remove(tmpDB)

	Init(tmpDB)
	defer Close()
	Seed()

	stats, err := GetDashboardStats()
	if err != nil {
		t.Fatalf("GetDashboardStats failed: %v", err)
	}
	if stats.TotalFlagged <= 0 {
		t.Errorf("expected TotalFlagged > 0, got %d", stats.TotalFlagged)
	}
	if stats.ReportsToday <= 0 {
		t.Errorf("expected ReportsToday > 0, got %d", stats.ReportsToday)
	}
}

func TestInsertReport(t *testing.T) {
	tmpDB := os.TempDir() + "qfis_report_test.db"
	defer os.Remove(tmpDB)

	Init(tmpDB)
	defer Close()
	Seed()

	merchants, _ := GetAllMerchants()
	if len(merchants) == 0 {
		t.Fatal("no merchants")
	}

	report := &models.Report{
		ID:         "RPT-TEST-001",
		MerchantID: merchants[0].ID,
		Note:       "Test report",
		Score:      85.5,
		Status:     "queued",
	}
	if err := InsertReport(report); err != nil {
		t.Fatalf("InsertReport failed: %v", err)
	}

	reports, _ := GetReports(100)
	found := false
	for _, r := range reports {
		if r.ID == "RPT-TEST-001" {
			found = true
			if r.Score != 85.5 {
				t.Errorf("expected score 85.5, got %f", r.Score)
			}
			break
		}
	}
	if !found {
		t.Error("inserted report not found")
	}
}

func TestUpsertMerchant(t *testing.T) {
	tmpDB := os.TempDir() + "qfis_upsert_test.db"
	defer os.Remove(tmpDB)

	Init(tmpDB)
	defer Close()

	m := &models.Merchant{
		ID:      "QRIS-TEST",
		Name:    "Test Merchant",
		MCC:     "5812",
		NPWP:    true,
		Flagged: false,
		X:       0.5,
		Y:       0.5,
	}
	if err := UpsertMerchant(m); err != nil {
		t.Fatalf("UpsertMerchant failed: %v", err)
	}

	got, err := GetMerchant("QRIS-TEST")
	if err != nil {
		t.Fatalf("GetMerchant failed: %v", err)
	}
	if got.Name != "Test Merchant" {
		t.Errorf("expected 'Test Merchant', got '%s'", got.Name)
	}
}
