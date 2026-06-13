package db

import (
	"database/sql"
	"fmt"
	"time"

	"qfis/internal/models"

	_ "modernc.org/sqlite"
)

var conn *sql.DB

func Init(path string) error {
	var err error
	conn, err = sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	conn.Exec("PRAGMA journal_mode=WAL")
	conn.Exec("PRAGMA foreign_keys=ON")

	if err := migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

func migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS merchants (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		mcc TEXT NOT NULL DEFAULT '',
		npwp_valid INTEGER NOT NULL DEFAULT 0,
		flagged INTEGER NOT NULL DEFAULT 0,
		x REAL NOT NULL DEFAULT 0,
		y REAL NOT NULL DEFAULT 0,
		created_at INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS reports (
		id TEXT PRIMARY KEY,
		merchant_id TEXT NOT NULL,
		note TEXT DEFAULT '',
		score REAL NOT NULL DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'queued',
		created_at INTEGER NOT NULL,
		FOREIGN KEY (merchant_id) REFERENCES merchants(id)
	);

	CREATE TABLE IF NOT EXISTS rules_cache (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		merchant_id TEXT NOT NULL,
		rule_name TEXT NOT NULL,
		score REAL NOT NULL DEFAULT 0,
		checked_at INTEGER NOT NULL
	);
	`
	_, err := conn.Exec(schema)
	return err
}

func Close() {
	if conn != nil {
		conn.Close()
	}
}

func GetMerchant(id string) (*models.Merchant, error) {
	row := conn.QueryRow(
		`SELECT id, name, mcc, npwp_valid, flagged, x, y FROM merchants WHERE id = ?`,
		id,
	)
	m := &models.Merchant{}
	err := row.Scan(&m.ID, &m.Name, &m.MCC, &m.NPWP, &m.Flagged, &m.X, &m.Y)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func SearchMerchants(q string) ([]models.Merchant, error) {
	like := "%" + q + "%"
	rows, err := conn.Query(
		`SELECT id, name, mcc, npwp_valid, flagged, x, y FROM merchants WHERE id LIKE ? OR name LIKE ? LIMIT 20`,
		like, like,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchants []models.Merchant
	for rows.Next() {
		var m models.Merchant
		if err := rows.Scan(&m.ID, &m.Name, &m.MCC, &m.NPWP, &m.Flagged, &m.X, &m.Y); err != nil {
			return nil, err
		}
		merchants = append(merchants, m)
	}
	return merchants, nil
}

func GetAllMerchants() ([]models.Merchant, error) {
	rows, err := conn.Query(
		`SELECT id, name, mcc, npwp_valid, flagged, x, y FROM merchants`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchants []models.Merchant
	for rows.Next() {
		var m models.Merchant
		if err := rows.Scan(&m.ID, &m.Name, &m.MCC, &m.NPWP, &m.Flagged, &m.X, &m.Y); err != nil {
			return nil, err
		}
		merchants = append(merchants, m)
	}
	return merchants, nil
}

func UpsertMerchant(m *models.Merchant) error {
	_, err := conn.Exec(
		`INSERT INTO merchants (id, name, mcc, npwp_valid, flagged, x, y, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET name=excluded.name, flagged=excluded.flagged`,
		m.ID, m.Name, m.MCC, btoi(m.NPWP), btoi(m.Flagged), m.X, m.Y, time.Now().Unix(),
	)
	return err
}

func InsertReport(r *models.Report) error {
	_, err := conn.Exec(
		`INSERT INTO reports (id, merchant_id, note, score, status, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		r.ID, r.MerchantID, r.Note, r.Score, r.Status, time.Now().Unix(),
	)
	return err
}

func GetReports(limit int) ([]models.Report, error) {
	rows, err := conn.Query(
		`SELECT r.id, r.merchant_id, COALESCE(m.name,''), r.score, r.status, r.created_at
		 FROM reports r LEFT JOIN merchants m ON r.merchant_id = m.id
		 ORDER BY r.created_at DESC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var r models.Report
		var ts int64
		if err := rows.Scan(&r.ID, &r.MerchantID, &r.MerchantName, &r.Score, &r.Status, &ts); err != nil {
			return nil, err
		}
		ago := time.Since(time.Unix(ts, 0))
		r.Time = fmtAgo(ago)
		r.Status = riskLabel(r.Score)
		reports = append(reports, r)
	}
	return reports, nil
}

func GetDashboardStats() (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}

	conn.QueryRow(`SELECT COUNT(*) FROM merchants WHERE flagged=1`).Scan(&stats.TotalFlagged)

	conn.QueryRow(`SELECT COUNT(*) FROM merchants WHERE flagged=1`).Scan(&stats.HighRiskCount)

	todayStart := time.Now().Truncate(24 * time.Hour).Unix()
	conn.QueryRow(`SELECT COUNT(*) FROM reports WHERE created_at >= ?`, todayStart).Scan(&stats.ReportsToday)

	conn.QueryRow(`SELECT COUNT(*) FROM reports WHERE status='queued'`).Scan(&stats.PendingOJK)

	return stats, nil
}

func SaveRuleCache(merchantID, ruleName string, score float64) error {
	_, err := conn.Exec(
		`INSERT INTO rules_cache (merchant_id, rule_name, score, checked_at) VALUES (?, ?, ?, ?)`,
		merchantID, ruleName, score, time.Now().Unix(),
	)
	return err
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func fmtAgo(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "baru saja"
	case d < time.Hour:
		m := int(d.Minutes())
		return fmt.Sprintf("%d menit lalu", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		return fmt.Sprintf("%d jam lalu", h)
	default:
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%d hari lalu", days)
	}
}

func riskLabel(score float64) string {
	if score >= 70 {
		return "MERAH"
	}
	if score >= 30 {
		return "KUNING"
	}
	return "HIJAU"
}
