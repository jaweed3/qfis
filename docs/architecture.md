# QRIS Fraud Intelligence System (QFIS)
> Anti-Judol Hackathon Project — Solo Dev Architecture

---

## 1. Problem Statement

Judol bandar gak mati karena website-nya diblokir — mereka pindah server dalam jam. Yang bikin mereka hidup adalah **alur dana**. Sejak 2025, deposit judol beralih masif ke QRIS karena:
- Transfer instan, sulit dilacak manual
- Merchant ID bisa dibuat bulk dengan identitas palsu
- "Smurfing" — ribuan transaksi kecil untuk obscure aliran dana besar

**Target sistem ini bukan pemain judol. Target adalah fingerprint merchant QRIS ilegal dan laporkan ke OJK/BI secara otomatis.**

---

## 2. Threat Model yang Diserang

```
[Pemain Judol]
      |
      | deposit via QRIS (< Rp 500rb, berulang, jam malam)
      v
[Merchant QRIS Shell] ← TARGET DETEKSI
      |
      | agregasi dana (smurfing)
      v
[Rekening Aggregator] ← TARGET DETEKSI
      |
      | transfer ke luar negeri / crypto
      v
[Bandar / Server Luar]
```

**Pola yang bisa dideteksi tanpa akses data bank:**
1. Merchant QRIS tanpa SIUP/NPWP valid (data publik OSS)
2. Merchant name mengandung pola evasion ("pulsa", "token", "digital goods")
3. Velocity pattern dari laporan publik / crowdsource
4. Domain website terkait QRIS yang match judol keyword

---

## 3. System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    DATA INGESTION LAYER                  │
│                                                         │
│  [OSS/NIB Scraper]  [Komdigi Blocklist]  [Crowdsource]  │
│         |                  |                   |        │
│         └──────────────────┴───────────────────┘        │
│                            |                            │
│                    [ETL Pipeline]                       │
│                    (Python + Pandas)                    │
└─────────────────────────────┬───────────────────────────┘
                              |
┌─────────────────────────────▼───────────────────────────┐
│                  INTELLIGENCE ENGINE                     │
│                                                         │
│  ┌─────────────────┐    ┌──────────────────────────┐   │
│  │  Rule Engine    │    │   ML Risk Scorer         │   │
│  │                 │    │                          │   │
│  │ • NPWP check    │    │ • IndoBERT merchant name │   │
│  │ • Velocity rule │    │   classifier             │   │
│  │ • Time pattern  │    │ • Anomaly detection      │   │
│  │ • Geo cluster   │    │   (Isolation Forest)     │   │
│  └────────┬────────┘    └────────────┬─────────────┘   │
│           └────────────┬─────────────┘                  │
│                        |                               │
│              [Risk Score Aggregator]                    │
│              Score 0–100 per merchant ID               │
└─────────────────────────┬───────────────────────────────┘
                          |
┌─────────────────────────▼───────────────────────────────┐
│                   OUTPUT LAYER                           │
│                                                         │
│  [Dashboard Web]    [REST API]    [Auto-Report PDF]     │
│  (FastAPI + React)  (public)      (ke OJK/BI format)   │
└─────────────────────────────────────────────────────────┘
```

---

## 4. Tech Stack — Solo Dev Friendly

| Layer | Tech | Alasan |
|-------|------|--------|
| Backend | **FastAPI** (Python) | Kamu udah familiar, async, cepat |
| Database | **PostgreSQL** (Supabase) | Udah punya akun dari CatatTernak |
| ML | **scikit-learn** + HuggingFace | Isolation Forest + IndoBERT-lite |
| Scraping | **Playwright** + **httpx** | Headless scraping OSS Indonesia |
| Frontend | **React** + Tailwind | Dashboard minimal, bisa pakai template |
| Deployment | **Railway** atau **Fly.io** | Free tier, satu command deploy |
| Report | **WeasyPrint** | Generate PDF laporan otomatis |

**Total estimasi setup:** 2–3 jam kalau environment bersih.

---

## 5. Data Sources (Semua Publik, No Auth Required)

### 5.1 Data Primer
```
1. OSS Indonesia (oss.go.id)
   - Cek NIB/NPWP merchant
   - Endpoint: pencarian badan usaha

2. Komdigi Blocklist
   - https://trustpositif.kominfo.go.id/
   - Domain judol yang sudah diblokir → ground truth label

3. PPATK Press Release
   - Pola transaksi mencurigakan yang dipublikasi
   - Manual extraction → seed rules

4. BI QRIS Merchant Registry
   - Merchant name + kategori MCC (Merchant Category Code)
   - Bisa minta akses publik atau scrape dari sandbox BI
```

### 5.2 Data Crowdsource (Fitur Unik)
```
User bisa submit:
- QRIS merchant ID yang dicurigai
- Screenshot bukti (OCR → extract merchant name)
- Upvote/downvote komunitas

→ Community intelligence layer yang self-improving
```

---

## 6. ML Model Detail

### 6.1 Merchant Name Classifier
```python
# Fine-tune IndoBERT pada task binary classification
# Label: LEGIT vs SUSPICIOUS

# Training data (bootstrap):
# LEGIT: merchant name dari Tokopedia, Gojek, dll (publik)
# SUSPICIOUS: nama merchant dari kasus PPATK yang dipublikasi
#             + keyword list: "pulsa murah", "token", "reload"

from transformers import AutoTokenizer, AutoModelForSequenceClassification

model_name = "indolem/indobert-base-uncased"
# Fine-tune dengan ~500 contoh → cukup untuk hackathon
```

### 6.2 Behavioral Anomaly (Rule-based dulu, ML nanti)
```python
# Isolation Forest pada fitur:
features = {
    "hour_of_report": int,          # jam crowdsource dilaporkan
    "merchant_age_days": int,        # umur QRIS merchant (dari registry)
    "name_entropy": float,           # keanehan nama merchant
    "mcc_gambling_proximity": float, # cosine distance ke kata judol
    "npwp_valid": bool,              # hasil cek OSS
    "domain_blocklisted": bool,       # ada di trustpositif
}
```

---

## 7. Risk Scoring Formula

```
Risk Score = (
    0.30 * npwp_invalid_score +
    0.25 * ml_classifier_score +
    0.20 * behavioral_anomaly_score +
    0.15 * community_report_score +
    0.10 * domain_blocklist_score
) * 100

Threshold:
- Score < 30  → HIJAU (legitimate)
- Score 30–70 → KUNING (watch list)
- Score > 70  → MERAH (flag untuk laporan OJK)
```

---

## 8. API Endpoints

```
POST /api/v1/merchant/check
  body: { qris_merchant_id: str, merchant_name: str }
  return: { risk_score: float, flags: [], recommendation: str }

POST /api/v1/report/submit
  body: { merchant_id: str, evidence_url: str, reporter_note: str }
  return: { report_id: str, status: "queued" }

GET  /api/v1/dashboard/stats
  return: { total_flagged: int, high_risk_count: int, reports_today: int }

GET  /api/v1/report/export/{report_id}
  return: PDF (format laporan OJK)
```

---

## 9. Sprint Plan — Solo Dev, Hackathon Timeline

### Asumsi: 48–72 jam hackathon

```
JAM 0–6: FOUNDATION
├── Setup repo, FastAPI skeleton, Supabase schema
├── Scraper Komdigi blocklist (ground truth)
└── Rule engine basic (NPWP check via OSS)

JAM 6–18: INTELLIGENCE ENGINE  
├── IndoBERT merchant classifier (fine-tune minimal)
├── Isolation Forest behavioral anomaly
├── Risk score aggregator
└── Unit test semua komponen

JAM 18–30: API + DASHBOARD
├── Semua endpoint FastAPI
├── React dashboard (pakai shadcn/ui biar cepat)
├── Crowdsource submit form
└── Risk score visualization

JAM 30–42: REPORT GENERATION + DEMO DATA
├── PDF auto-report generator (format OJK)
├── Seed database dengan ~50 merchant contoh
├── Demo flow end-to-end
└── Stress test

JAM 42–48: POLISH + PITCH
├── Landing page / README
├── Recorded demo video
└── Slide pitch (5 menit)
```

---

## 10. Differentiator vs Tim Lain

| Solusi Umum (Tim Lain) | QFIS (Kamu) |
|------------------------|-------------|
| Blokir website judol | Fingerprint alur DANA-nya |
| Edukasi literasi digital | Attack supply chain finansial |
| Deteksi konten/keyword | Deteksi merchant QRIS ilegal |
| Reactive (setelah viral) | Proactive (sebelum laporan polisi) |
| Output: notifikasi | Output: laporan siap submit OJK |

---

## 11. Limitation Jujur (Untuk Pitch)

1. **Tanpa akses data transaksi BI/OJK** — sistem ini inference dari data publik + crowdsource, bukan raw transaction. Akurasi terbatas.
2. **False positive risk** — merchant QRIS legit dengan nama mirip bisa kena flag. Threshold harus konservatif.
3. **Adversarial evasion** — bandar bisa ganti nama merchant tiap hari. Sistem butuh continuous learning.

**Cara pitch-nya:** Frame sebagai "early warning intelligence layer" yang augment kemampuan PPATK/OJK, bukan replace investigasi manual mereka.

---

## 12. Repo Structure

```
qfis/
├── backend/
│   ├── main.py                 # FastAPI entrypoint
│   ├── api/
│   │   ├── merchant.py
│   │   ├── report.py
│   │   └── dashboard.py
│   ├── engine/
│   │   ├── rule_engine.py      # Rule-based checks
│   │   ├── ml_classifier.py    # IndoBERT wrapper
│   │   ├── anomaly.py          # Isolation Forest
│   │   └── risk_scorer.py      # Aggregator
│   ├── scrapers/
│   │   ├── komdigi.py
│   │   └── oss_checker.py
│   └── reports/
│       └── pdf_generator.py
├── frontend/
│   ├── src/
│   │   ├── pages/
│   │   │   ├── Dashboard.jsx
│   │   │   ├── CheckMerchant.jsx
│   │   │   └── Report.jsx
│   │   └── components/
│   │       ├── RiskBadge.jsx
│   │       └── MerchantCard.jsx
├── ml/
│   ├── train_classifier.py
│   ├── data/
│   │   ├── legit_merchants.csv
│   │   └── suspicious_merchants.csv
│   └── models/                 # saved .pkl / .bin
├── docker-compose.yml
└── README.md
```

---

## 13. Nama Project Rekomendasi

- **QFIS** — QRIS Fraud Intelligence System *(profesional, pitch ke regulator)*
- **JaringBandar** *(kontekstual, memorable untuk juri)*  
- **SIGAP-QRIS** *(akronim, kesan government-ready)*

---

*Generated for hackathon purposes. Data sources used are all publicly available. System designed to assist OJK/BI/PPATK regulatory reporting, not to replace official investigation procedures.*
