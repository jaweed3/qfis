# QFIS — QRIS Fraud Intelligence System

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev)
[![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)](https://react.dev)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

**Target:** Fingerprint merchant QRIS ilegal dan laporkan ke OJK/BI secara otomatis. Bukan blokir website — ini serang **alur dana**-nya.

> Anti-Judol Hackathon Project — Solo Dev Architecture

---

## Fitur

- **Risk Engine** — Rule-based + heuristic scoring (0–100) per merchant QRIS
- **Network Graph** — Animated canvas visualization koneksi antar merchant mencurigakan
- **Merchant Checker** — Cek risiko merchant QRIS instan via ID atau nama
- **Crowdsource Reports** — Community intelligence layer, laporan terverifikasi
- **OJK-ready Reports** — Format laporan standar untuk regulator (coming soon)
- **Klein Blue Design System** — Forensic intelligence UI dengan Playfair + JetBrains Mono

## Quick Start

### Opsi 1 — Binary (recommended untuk server)

```bash
# Build (butuh Go 1.26+ & Node 22+)
./build.sh

# Atau download binary dari releases
scp backend/qfis-linux user@server:~/qfis/
ssh user@server
cd ~/qfis
mkdir -p data
GIN_MODE=release PORT=8080 ./qfis-linux
# → http://server:8080
```

### Opsi 2 — Docker

```bash
cp .env.example .env
# Edit API_KEY jika perlu
docker compose up -d
# → http://localhost:8080
```

### Opsi 3 — Development

```bash
# Terminal 1 — Backend
cd backend
go run .    # → :8080

# Terminal 2 — Frontend (hot reload)
cd frontend
npm install
npm run dev  # → :5173 (proxy ke :8080)
```

---

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/merchant/check` | Cek risiko merchant |
| GET | `/api/v1/merchant/search?q=` | Cari merchant |
| GET | `/api/v1/merchant/:id` | Detail merchant |
| POST | `/api/v1/report/submit` | Laporkan merchant |
| GET | `/api/v1/report/list` | Daftar laporan |
| GET | `/api/v1/report/:id` | Detail laporan |
| GET | `/api/v1/dashboard/stats` | Statistik dashboard |
| GET | `/api/v1/dashboard/network` | Data network graph |

### Contoh: Cek Merchant

```bash
curl -X POST http://localhost:8080/api/v1/merchant/check \
  -H "Content-Type: application/json" \
  -d '{"qris_merchant_id": "QRIS-9942A", "merchant_name": "Pulsa Murah 24Jam"}'
```

Response:
```json
{
  "risk_score": 72.4,
  "risk_label": "MERAH",
  "flags": [
    "NPWP tidak terdaftar",
    "Nama mengandung keyword mencurigakan: pulsa"
  ],
  "recommendation": "Segera laporkan ke OJK untuk investigasi lebih lanjut"
}
```

### Auth (opsional)

Set env `API_KEY=secret123` — semua request ke `/api/*` butuh header:
```
X-API-Key: secret123
```

---

## Risk Scoring Formula

| Komponen | Bobot | Sumber |
|----------|-------|--------|
| NPWP Invalid | 30% | OSS Checker |
| Keyword Match | 25% | Classifier |
| Anomaly Behavior | 20% | Isolation Forest |
| Community Report | 15% | Crowdsource |
| Domain Blocklist | 10% | Komdigi |

**Threshold:** `<30 HIJAU` · `30–70 KUNING` · `≥70 MERAH`

---

## Project Structure

```
qfis/
├── backend/                  # Go backend
│   ├── main.go               # Entrypoint (API + embedded frontend)
│   ├── internal/
│   │   ├── api/              # HTTP handlers
│   │   ├── engine/           # Risk engine (rule, classifier, anomaly, scorer)
│   │   ├── db/               # SQLite + seeder
│   │   └── models/           # Data types
│   └── frontend/dist/        # Embedded React build
├── frontend/                 # React frontend
│   └── src/
│       ├── components/       # UI components
│       ├── pages/            # Dashboard, Check, Reports
│       └── lib/api.js        # API client
├── docs/                     # Documentation
│   ├── architecture.md       # System architecture
│   ├── design.md             # Klein Blue design system
│   └── frontend.md           # Frontend component spec
├── Dockerfile                # Multi-stage Docker build
├── docker-compose.yml        # One-command deploy
├── build.sh / build.bat      # Build scripts
└── .env.example              # Environment config template
```

---

## Tech Stack

| Layer | Tech |
|-------|------|
| Backend | Go + Gin + SQLite |
| Frontend | React 18 + Vite + Tailwind CSS v4 |
| Fonts | Playfair Display + JetBrains Mono + Inter |
| Embed | `embed.FS` — single binary deployment |
| ML | Rule-based + heuristics (scikit-learn planned) |

---

## Security

- **GIN_MODE=release** by default
- **CORS** configurable via `CORS_ORIGINS` env
- **API Key** optional auth via `X-API-Key` header
- **SQL Injection** prevented — parameterized queries
- **Path Traversal** prevented — embed.FS sanitizes paths
- **Docker** runs as non-root user
- **Rate Limiting** — behind nginx recommended

> Untuk production: pasang di belakang nginx/caddy untuk HTTPS + rate limiting.

---

## Deployment Checklist

- [ ] Set `API_KEY` untuk proteksi API
- [ ] Set `CORS_ORIGINS` ke domain frontend
- [ ] Gunakan reverse proxy (nginx) untuk HTTPS
- [ ] Backup `data/qfis.db` secara periodik
- [ ] Monitor logs: `docker compose logs -f`

---

## License

MIT — built for hackathon purposes. Data sources used are all publicly available. System designed to assist OJK/BI/PPATK regulatory reporting, not to replace official investigation procedures.
