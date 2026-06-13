# QFIS Frontend Design Document
> QRIS Fraud Intelligence System — Visual & Component Spec

---

## 1. Design Concept

**"Forensic Network"** — bukan dashboard fintech biasa, bukan gov site membosankan.

Visual language-nya diambil dari dunia **fraud investigation tools** (OSINT platforms, network analysis) — dark, data-dense, authoritative. Setiap elemen harus communicate: *ini serius, ini akurat, ini bisa dipercaya regulator.*

**Signature element:** Live animated network graph di hero — setiap node adalah merchant QRIS, warna = risk score, edge = koneksi antar merchant yang dicurigai. Bukan static chart. Ini yang pertama orang lihat, dan ini yang mereka inget.

---

## 2. Design Tokens

### Color Palette
```
--void:       #0A0E1A   → background utama (deep navy, bukan pure black)
--surface:    #0F1F3D   → card, panel background
--elevated:   #1A3A6B   → border, divider, elevated element
--cyan:       #00D4FF   → primary accent — data, live indicator, link
--red:        #FF3B5C   → high risk, alert, danger
--amber:      #F59E0B   → medium risk, warning, pending
--green:      #10B981   → low risk, safe, success
--muted:      #8892A4   → secondary text, placeholder, label
--primary:    #E2E8F0   → body text
```

### Typography
```
Display/Data: JetBrains Mono — semua angka, ID, kode, label kapital
Body:         Inter        — semua prose, deskripsi, UI copy
```

**Aturan type:**
- Semua angka → JetBrains Mono, minimum 600 weight
- Semua label uppercase → JetBrains Mono, letter-spacing: 2px
- Prose/help text → Inter 400, color: --muted
- Score/metric besar → JetBrains Mono 700, ukuran 24–32px

### Border Radius
```
Card/Panel:   8px
Button:       6px
Badge:        4px
Input:        6px
```

---

## 3. Component Library

### 3.1 RiskBadge
```
Props: score (0–100)
Variants:
  MERAH  (≥70) → color: #FF3B5C, bg: #FF3B5C18, border: #FF3B5C44
  KUNING (30–69) → color: #F59E0B, bg: #F59E0B18, border: #F59E0B44
  HIJAU  (<30) → color: #10B981, bg: #10B98118, border: #10B98144

Font: JetBrains Mono 700, 11px, letter-spacing: 1px
Padding: 2px 8px
```

### 3.2 StatCard
```
Layout:
  ┌──────────────────┐
  │ 1,847            │  ← JetBrains Mono 28px, color=accent
  │                  │
  │ TOTAL TERFLAG    │  ← Inter 11px uppercase, color=muted
  └──────────────────┘

Background: --surface
Border: 1px solid --elevated
Padding: 16px 20px
```

### 3.3 NetworkGraph (Canvas)
```
Behavior:
  - Node drift: slow, randomized velocity ±0.3px/frame
  - Edge: dashed line antara flagged nodes dalam radius 180px
  - Edge animation: lineDashOffset moving (menandakan alur dana)
  - Node pulse: high-risk nodes punya expanding ring animation
  - Click: select node → tampil detail di panel kanan
  - Hover cursor: crosshair

Visual layer order (bottom to top):
  1. Background grid (cyan, opacity 4%)
  2. Edges (red dashed, animated)
  3. Node glow (radial gradient)
  4. Node core (filled circle)
  5. Selection ring (cyan, 2px)
  6. Node label (JetBrains Mono 10px)
```

### 3.4 MerchantDetail Panel
```
Triggered by: click node di NetworkGraph
Content:
  - Merchant name (JetBrains Mono 14px bold)
  - Merchant ID (JetBrains Mono 11px muted)
  - Risk score bar (gradient: green→red, animated width)
  - Score number large (JetBrains Mono 22px, color=riskColor)
  - Metadata grid: MCC, NPWP, Laporan, Status
  - CTA button "LAPORKAN KE OJK" — hanya muncul jika score ≥ 70

Empty state:
  - Dashed circle icon
  - Help text "Klik node pada graph"
  - Opacity: 40%
```

### 3.5 MerchantChecker
```
Flow:
  Input → [SCAN] button → loading 900ms → result card

Input: JetBrains Mono 13px, background --void
Button: cyan when idle, muted when loading

Result card:
  - Border color: riskColor(score)
  - 3 metadata items: Risk Score, NPWP, Laporan
  - Warning banner jika score ≥ 70
```

### 3.6 ReportForm
```
Fields: Merchant ID (required), Catatan (optional textarea)
Submit: amber-themed (warna warning, bukan primary action)
Success state: green confirmation dengan Report ID
Auto-reset: 3 detik setelah submit
```

---

## 4. Page Layout

### Header (sticky, 56px)
```
LEFT:   Logo dot + "QFIS" + tagline
CENTER: Tab navigation (DASHBOARD / CEK MERCHANT / LAPORAN)
RIGHT:  Live indicator (green dot pulse + "LIVE")

Background: #0A0E1Aee + backdrop-filter: blur(8px)
Border-bottom: 1px solid --elevated
```

### Dashboard Tab
```
┌─────────────────────────────────────────────────────┐
│  [Stat]  [Stat]  [Stat]  [Stat]                    │
├───────────────────────────┬─────────────────────────┤
│                           │                         │
│   Network Graph           │   Merchant Detail       │
│   (canvas, 420px high)    │   Panel                 │
│                           │                         │
│   [Legend bar]            │                         │
├───────────────────────────┴─────────────────────────┤
│   Recent Reports Table                              │
└─────────────────────────────────────────────────────┘
```

### Cek Merchant Tab
```
Max-width: 600px, centered
- Title + description
- MerchantChecker component
- Example IDs grid
```

### Laporan Tab
```
2-column grid, max-width 900px:
LEFT:  Antrian laporan ke OJK (list view)
RIGHT: ReportForm + disclaimer
```

---

## 5. Motion Design

| Element | Animation | Rationale |
|---------|-----------|-----------|
| Network nodes | Slow drift ±0.3px/frame | Simulates live data |
| Edge dash | lineDashOffset scroll | Menggambarkan alur dana mengalir |
| High-risk pulse ring | sin wave expand/contract | Draws attention ke threat |
| Risk score bar | width transition 0.6s ease | Satisfying reveal |
| Live counter | +1 tiap ~4 detik (random) | Sense of urgency |
| Tab switch | Instant — no slide | Responsiveness > aesthetics |

**Rule:** Satu "wow" moment (network graph). Semua animasi lain subtle. Jangan AI-generated vibe.

---

## 6. Responsiveness

```
Desktop (>1024px): Full layout, 2-column graph+detail
Tablet (768–1024px): Stack graph dan detail vertically
Mobile (<768px): Single column, graph height 280px, tabs jadi scroll
```

---

## 7. Copy Style Guide

- Semua label UI: **KAPITAL, letter-spacing tinggi** (forensic, bukan playful)
- Error messages: Langsung, tanpa apologi. "NPWP tidak terdaftar" bukan "Sepertinya NPWP belum ada"
- CTA: Aktif + spesifik. "LAPORKAN KE OJK" bukan "Submit"
- Empty states: Invitation to act. "Klik node untuk lihat detail" bukan "Belum ada data"
- Risk labels: MERAH/KUNING/HIJAU (familiar, Indonesian, tidak perlu translate)

---

## 8. File Structure

```
src/
├── App.jsx                  ← Single file untuk hackathon speed
│
├── Kalau mau split nanti:
├── components/
│   ├── NetworkGraph.jsx     ← Canvas animation
│   ├── RiskBadge.jsx
│   ├── StatCard.jsx
│   ├── MerchantDetail.jsx
│   ├── MerchantChecker.jsx
│   └── ReportForm.jsx
├── pages/
│   ├── Dashboard.jsx
│   ├── CheckMerchant.jsx
│   └── Reports.jsx
└── constants/
    └── theme.js             ← Design tokens
```

---

## 9. Dependencies

```json
{
  "react": "^18",
  "react-dom": "^18"
}
```

**Zero external UI library.** Semua komponen custom CSS-in-JS (inline styles). Ini bukan malas — ini intentional:
- Bundle size kecil
- Tidak ada Tailwind conflict
- Design token 100% controlled
- Tidak keliatan shadcn/MUI template

Google Fonts via `@import` di dalam style tag: `JetBrains Mono + Inter`.

---

## 10. Demo Preparation Checklist

- [ ] Mock data 8 merchant (campuran flagged + aman)
- [ ] Network graph klik-able dan responsive
- [ ] Merchant checker flow end-to-end (input → scan → result)
- [ ] Live counter naik selama demo
- [ ] Report submit → success state
- [ ] Mobile view tidak broken
- [ ] Semua risk badge warna benar

---

*Design language: Forensic Intelligence. Bukan startup. Bukan gov. Di antara keduanya — yang bikin regulator percaya dan juri terkesan.*
