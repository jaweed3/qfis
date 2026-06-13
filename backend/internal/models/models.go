package models

type Merchant struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	RiskScore   float64 `json:"risk_score"`
	RiskLabel   string  `json:"risk_label"`
	Reports     int     `json:"reports"`
	MCC         string  `json:"mcc"`
	NPWP        bool    `json:"npwp"`
	Flagged     bool    `json:"flagged"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
}

type Report struct {
	ID         string  `json:"id"`
	MerchantID string  `json:"merchant_id"`
	MerchantName string `json:"merchant_name"`
	Score      float64 `json:"score"`
	Status     string  `json:"status"`
	Time       string  `json:"time"`
	Note       string  `json:"note,omitempty"`
	CreatedAt  int64   `json:"created_at"`
}

type CheckRequest struct {
	MerchantID   string `json:"qris_merchant_id" binding:"required"`
	MerchantName string `json:"merchant_name"`
}

type CheckResponse struct {
	RiskScore      float64  `json:"risk_score"`
	RiskLabel      string   `json:"risk_label"`
	Flags          []string `json:"flags"`
	Recommendation string   `json:"recommendation"`
}

type ReportSubmitRequest struct {
	MerchantID   string `json:"merchant_id" binding:"required"`
	EvidenceURL  string `json:"evidence_url"`
	ReporterNote string `json:"reporter_note"`
}

type ReportSubmitResponse struct {
	ReportID string `json:"report_id"`
	Status   string `json:"status"`
}

type DashboardStats struct {
	TotalFlagged   int `json:"total_flagged"`
	HighRiskCount  int `json:"high_risk_count"`
	ReportsToday   int `json:"reports_today"`
	PendingOJK     int `json:"pending_ojk"`
}
