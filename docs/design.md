# QFIS Design System
## Klein Blue Authority — Frontend Specification

---

## 1. Concept

**Reference:** Hermes Agent (nousresearch.com) — tapi bukan clone. Ini *interpretation*.

Hermes pakai Klein Blue sebagai statement: *ini bukan SaaS, ini movement*. QFIS pakai DNA yang sama tapi dengan layer intelligence — fraud forensics yang punya moral authority. Kombinasi Serif baroque + monospace utility bukan kebetulan: serif = institusional, trustworthy, berat; mono = data, precise, tidak bisa dibantah. Clash keduanya adalah the whole point.

**Satu kalimat:** *Kalau OJK punya war room, ini tampilannya.*

---

## 2. Design Tokens

### 2.1 Color

```css
:root {
  /* Core */
  --qfis-blue:       #2B2BFF;   /* Klein Blue — primary brand */
  --qfis-blue-dark:  #1A1AE8;   /* hover state */
  --qfis-blue-void:  #0A0A2E;   /* deep blue — section bg */
  --qfis-white:      #FFFFFF;
  --qfis-off-white:  #F0EDFF;   /* tinted white — body text on blue */

  /* Semantic — Risk */
  --risk-red:        #FF2D55;   /* MERAH ≥70 */
  --risk-red-bg:     rgba(255, 45, 85, 0.12);
  --risk-amber:      #FF9F0A;   /* KUNING 30–69 */
  --risk-amber-bg:   rgba(255, 159, 10, 0.12);
  --risk-green:      #30D158;   /* HIJAU <30 */
  --risk-green-bg:   rgba(48, 209, 88, 0.12);

  /* UI Surface (non-blue sections) */
  --surface-0:       #07070F;   /* deepest bg */
  --surface-1:       #0D0D1F;   /* card bg */
  --surface-2:       #141428;   /* elevated */
  --surface-border:  rgba(43, 43, 255, 0.2);  /* blue-tinted border */
  --surface-border-strong: rgba(43, 43, 255, 0.5);

  /* Text */
  --text-primary:    #FFFFFF;
  --text-secondary:  rgba(255, 255, 255, 0.55);
  --text-muted:      rgba(255, 255, 255, 0.25);
  --text-on-blue:    #FFFFFF;
  --text-label:      rgba(255, 255, 255, 0.4); /* uppercase labels */
}
```

**Aturan penggunaan:**
- `--qfis-blue` hanya untuk: hero section, CTA primary, active state, logo
- Jangan pake blue untuk card background di dashboard — itu territory hero doang
- Semantic risk color tidak pernah dimodifikasi opacity-nya kecuali untuk bg (`-bg` variant)

---

### 2.2 Typography

```css
/* Import */
@import url('https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,700;0,900;1,700&family=JetBrains+Mono:wght@400;600;700&family=Inter:wght@400;500&display=swap');

:root {
  --font-display: 'Playfair Display', Georgia, serif;
  --font-mono:    'JetBrains Mono', 'Courier New', monospace;
  --font-body:    'Inter', -apple-system, sans-serif;
}
```

**Type Scale:**

| Token | Family | Size | Weight | Use |
|-------|--------|------|--------|-----|
| `--t-hero` | Playfair Display | 72–96px | 900 | Hero headline (mobile: 48px) |
| `--t-display` | Playfair Display | 40px | 900 | Section headline |
| `--t-heading` | Playfair Display | 24px | 700 | Card heading, panel title |
| `--t-label` | JetBrains Mono | 11px | 700 | ALL CAPS labels, eyebrow, nav |
| `--t-data` | JetBrains Mono | 28–40px | 700 | Big numbers, score display |
| `--t-data-sm` | JetBrains Mono | 14–16px | 600 | Merchant ID, code, metadata |
| `--t-body` | Inter | 15px | 400 | Prose, descriptions |
| `--t-caption` | Inter | 12px | 400 | Help text, footnote |

**Typography Rules:**
1. Hero: Playfair Display 900, ALL CAPS, tight tracking (`letter-spacing: -0.02em`)
2. All UI labels: JetBrains Mono 700, ALL CAPS, `letter-spacing: 0.15em`
3. All numbers/IDs/codes: JetBrains Mono — tidak ada pengecualian
4. Prose body: Inter 400 — hanya untuk deskripsi panjang, bukan UI
5. Jangan mix Playfair dan Inter dalam satu component — pilih satu

---

### 2.3 Spacing

```
4px   — micro gap (badge padding)
8px   — tight (icon gap, inline item)
12px  — component internal
16px  — card padding standard
20px  — card padding comfortable
24px  — section gap, grid gap
40px  — section padding
80px  — hero vertical padding
```

### 2.4 Border Radius

```
0px   — hero elements, full-bleed sections, terminal blocks
4px   — badge, tag, small pill
6px   — button
8px   — card, input, modal
```

**Rule:** Semakin "data/technical" sebuah element, semakin mendekati 0px radius-nya. Semakin "UI/friendly" element-nya, pakai 6–8px.

---

## 3. Layout System

### 3.1 Hero Section (Klein Blue Full-Bleed)

```
┌──────────────────────────────────────────────────────────┐
│  #2B2BFF background, 100vw                              │
│                                                          │
│  NAV: [QFIS ●] ─────────── [DOCS] [API] [INSTALL →]    │
│                                                          │
│  ┌────────────────────┐  ┌──────────────────────────┐   │
│  │ OPEN SOURCE · OJK  │  │                          │   │
│  │                    │  │  [Engraving/SVG motif]   │   │
│  │ THE FRAUD          │  │  Scales, blindfold,      │   │
│  │ INTELLIGENCE       │  │  network lines           │   │
│  │ SYSTEM             │  │                          │   │
│  │                    │  └──────────────────────────┘   │
│  │ [SCAN MERCHANT →]  │                                  │
│  └────────────────────┘                                  │
│                                                          │
│  ── 1,847 TERFLAG  ── 312 RISIKO TINGGI  ── 18 → OJK ──  │
└──────────────────────────────────────────────────────────┘
```

**Proporsi:** Text LEFT 45%, Illustration RIGHT 55% (persis Hermes ratio)
**Stat bar:** Full-width, border-top 1px solid rgba(255,255,255,0.15), monospace

### 3.2 Dashboard Grid (Dark Sections)

```
┌──────────────────────────────────────────────────────────┐
│  surface-0 background                                    │
│                                                          │
│  ┌───────────────────────────────┐  ┌─────────────────┐ │
│  │  Network Graph                │  │ Merchant Detail │ │
│  │  canvas, border: surface-border│  │ Panel           │ │
│  │  420px                        │  │                 │ │
│  └───────────────────────────────┘  └─────────────────┘ │
│                                                          │
│  ┌──────────────────────────────────────────────────────┐│
│  │  Recent Reports Table                               ││
│  └──────────────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────────┘
```

---

## 4. Component Specifications

### 4.1 Navigation Bar

```
Height: 56px
Background: rgba(43, 43, 255, 0.95) + backdrop-filter: blur(12px)
Border-bottom: 1px solid rgba(255,255,255,0.1)
Position: sticky top:0, z-index:100

LEFT:   Logo — "QFIS" (Playfair Display 700, 18px, white)
        Dot indicator — 8px circle, white, pulsing
CENTER: Links — JetBrains Mono 700 11px, letter-spacing 0.15em
        ALL CAPS: DASHBOARD · CEK MERCHANT · LAPORAN
        Active: white. Inactive: rgba(255,255,255,0.45)
RIGHT:  "INSTALL →" — outlined pill button
        border: 1px solid rgba(255,255,255,0.4)
        hover: background white, color blue
```

### 4.2 Hero Headline

```css
.hero-headline {
  font-family: var(--font-display);
  font-size: clamp(48px, 7vw, 96px);
  font-weight: 900;
  color: white;
  line-height: 0.92;
  letter-spacing: -0.02em;
  text-transform: uppercase;
}
```

**Copy terkunci:** Jangan ganti ke sentence case. Ini intentional typographic choice.

### 4.3 Eyebrow Label

```css
.eyebrow {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  color: rgba(255,255,255,0.55);
  letter-spacing: 0.2em;
  text-transform: uppercase;
  margin-bottom: 16px;
}
/* Separator dot: "OPEN SOURCE · MIT LICENSE" */
/* Use · (U+00B7) bukan dash */
```

### 4.4 CTA Button

```css
/* Primary — on blue background */
.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: white;
  color: #2B2BFF;
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.15em;
  text-transform: uppercase;
  padding: 12px 24px;
  border-radius: 2px;   /* almost sharp — intentional */
  border: none;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}
.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.3);
}

/* Secondary — on dark background */
.btn-secondary {
  background: transparent;
  color: white;
  border: 1px solid var(--surface-border-strong);
  /* same font/size/padding */
}
.btn-secondary:hover {
  border-color: rgba(43,43,255,0.8);
  background: rgba(43,43,255,0.1);
}

/* Danger — flag/report action */
.btn-danger {
  background: var(--risk-red-bg);
  color: var(--risk-red);
  border: 1px solid rgba(255,45,85,0.3);
}
```

### 4.5 Risk Badge

```css
.risk-badge {
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.15em;
  text-transform: uppercase;
  padding: 3px 8px;
  border-radius: 4px;
}

/* Variants */
.risk-badge.merah  { color: var(--risk-red);   background: var(--risk-red-bg);   border: 1px solid rgba(255,45,85,0.25); }
.risk-badge.kuning { color: var(--risk-amber);  background: var(--risk-amber-bg); border: 1px solid rgba(255,159,10,0.25); }
.risk-badge.hijau  { color: var(--risk-green);  background: var(--risk-green-bg); border: 1px solid rgba(48,209,88,0.25); }
```

### 4.6 Stat Card

```
Layout:
┌─────────────────────────┐
│ 1,847                   │  ← JetBrains Mono 700, 36px, accent color
│                         │
│ TOTAL TERFLAG           │  ← JetBrains Mono 700, 10px, --text-muted
└─────────────────────────┘

Background: var(--surface-1)
Border: 1px solid var(--surface-border)
Border-radius: 8px
Padding: 20px

Number color mapping:
  Total/neutral → --qfis-blue (tetap branded)
  Risiko tinggi → --risk-red
  Pending/watch → --risk-amber
  Resolved/safe → --risk-green
```

### 4.7 Card (Dashboard Panel)

```css
.card {
  background: var(--surface-1);
  border: 1px solid var(--surface-border);
  border-radius: 8px;
  overflow: hidden;
}

.card-header {
  padding: 12px 20px;
  border-bottom: 1px solid var(--surface-border);
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 700;
  color: var(--text-label);
  letter-spacing: 0.2em;
  text-transform: uppercase;
}

.card-body {
  padding: 20px;
}
```

### 4.8 Input Field

```css
.input {
  background: var(--surface-0);
  border: 1px solid var(--surface-border);
  border-radius: 6px;
  padding: 11px 16px;
  color: var(--text-primary);
  font-family: var(--font-mono);
  font-size: 13px;
  outline: none;
  transition: border-color 0.15s;
  width: 100%;
}
.input::placeholder { color: var(--text-muted); }
.input:focus        { border-color: rgba(43,43,255,0.6); }
```

### 4.9 Table Row (Reports)

```
┌───────────────────────────────────────────────────────┐
│ RPT-001  │ Pulsa Murah 24Jam    │ 2 menit lalu │ MERAH│
│ [mono]   │ [Inter 13px]         │ [mono muted] │[badge]│
└───────────────────────────────────────────────────────┘

Row height: 52px
Row border-bottom: 1px solid rgba(255,255,255,0.04)
Row hover: background rgba(43,43,255,0.05)
Padding: 0 20px
```

---

## 5. Signature Element — Engraving Motif

Ini yang bikin QFIS beda dari semua dashboard hackathon.

**Reference:** Hermes Agent pakai engraving dewa Yunani. QFIS pakai engraving **Themis** (dewi keadilan, timbangan + pedang) atau motif **network web** bergaya etching abad 19.

**Implementasi:**
```
Option 1 (Rekomendasi): SVG engraving-style illustration
  - Garis tipis (stroke 0.5–1px), white on blue
  - Dense crosshatching untuk shadow
  - Subject: scales of justice + network nodes
  - Size: 400–500px, positioned absolute kanan hero

Option 2 (Cepat): CSS-only geometric
  - Konsentris circle dengan radial lines
  - Opacity 15–25%
  - Pure CSS, zero asset dependency

Option 3 (Paling Cepat): Background pattern
  - repeating-linear-gradient grid
  - opacity 6%, white on blue
  - Subtle texture, bukan focal point
```

**Untuk hackathon dengan time constraint: pakai Option 2 (CSS geometric) + focus ke konten.**

---

## 6. Motion Design

### 6.1 Prinsip
- **Purposeful, not decorative.** Setiap animasi harus punya alasan data/UX, bukan aesthetic semata
- **Fast entry, subtle loop.** Elemen masuk cepat (150–300ms), loop animation lambat dan subtle
- **No bounce, no spring** — ini bukan consumer app. `ease` atau `cubic-bezier(0.16,1,0.3,1)`

### 6.2 Spesifikasi per Element

| Element | Animation | Duration | Easing |
|---------|-----------|----------|--------|
| Page load | Fade up (y: 20px → 0, opacity 0 → 1) | 400ms | ease-out |
| Stat number | Count-up on mount | 1200ms | ease-out |
| Risk score bar | Width 0% → actual% | 600ms | cubic-bezier(0.16,1,0.3,1) |
| Network node pulse | Scale 1 → 1.15 → 1, opacity ping | 2s infinite | ease-in-out |
| Network edge dash | lineDashOffset scroll | continuous | linear |
| Live counter | +1 dengan fade | 300ms | ease |
| Hover button | translateY -1px | 150ms | ease |
| Tab switch | Opacity 0 → 1 | 200ms | ease |
| Scan loading | Dot pulse 3x | 900ms | ease-in-out |

### 6.3 Live Pulse Indicator (Nav)

```css
@keyframes live-pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50%       { opacity: 0.4; transform: scale(0.85); }
}
.live-dot {
  width: 6px; height: 6px;
  background: #30D158;
  border-radius: 50%;
  animation: live-pulse 2s ease-in-out infinite;
}
```

---

## 7. Engraving Illustration Spec (SVG)

Kalau ada waktu build ilustrasi, ini specnya:

```
Subject:  Timbangan keadilan (Themis/Justitia) dengan:
          - Network node di kedua piringan
          - Garis radial dari center (seperti Hermes angel)
          - Merchant QRIS di satu sisi, OJK logo di sisi lain

Style:    Pen-and-ink engraving
          - Semua stroke putih (#FFFFFF)
          - Opacity varies: 0.6 (detail) → 0.15 (background hatch)
          - Stroke width: 0.5px (hatch) → 1.5px (outline)
          - No fill shapes, hanya lines

Size:     480 × 520px viewBox
Color:    White on transparent (drop in blue bg)

Fallback: Jika tidak sempat, pakai SVG geometric:
```

```svg
<!-- Fallback: Concentric rings dengan radial lines -->
<svg viewBox="0 0 400 400" xmlns="http://www.w3.org/2000/svg">
  <!-- Concentric circles -->
  <circle cx="200" cy="200" r="60"  fill="none" stroke="white" stroke-width="0.5" opacity="0.2"/>
  <circle cx="200" cy="200" r="100" fill="none" stroke="white" stroke-width="0.5" opacity="0.15"/>
  <circle cx="200" cy="200" r="150" fill="none" stroke="white" stroke-width="0.5" opacity="0.1"/>
  <circle cx="200" cy="200" r="190" fill="none" stroke="white" stroke-width="0.5" opacity="0.07"/>
  <!-- Radial lines (setiap 15°) -->
  <!-- ... generate 24 lines dari center ke tepi -->
  <!-- Node dots di intersection -->
</svg>
```

---

## 8. Copy System

### 8.1 Voice
- **Bukan startup** → tidak ada "We help you...", "Supercharge your...", "Finally..."
- **Bukan pemerintah** → tidak ada "Dalam rangka...", "Berdasarkan Peraturan..."
- **Tone:** Intelligence analyst briefing. Singkat. Faktual. Berbobot.

### 8.2 Pola Headline
```
✓ "THE FRAUD INTELLIGENCE SYSTEM"     ← Playfair, caps, declarative
✓ "BANDAR TERDETEKSI"                 ← langsung ke point
✗ "Kami membantu Anda mendeteksi..."  ← too soft
✗ "Smart QRIS Fraud Detection"        ← too startup
```

### 8.3 UI Label Conventions
```
Section headers:  SEMUA KAPITAL, letter-spacing 0.2em, JetBrains Mono
CTA buttons:      KAPITAL + arrow → (bukan "Submit" atau "Click here")
Badge labels:     MERAH / KUNING / HIJAU (bukan "HIGH" / "MEDIUM" / "LOW")
Report ID format: RPT-XXXX (zero-padded 4 digit)
Merchant ID:      QRIS-XXXX (selalu prefix QRIS)
Scores:           "94/100" (bukan "94%" atau "Score: 94")
```

### 8.4 Error & Empty States
```
Empty network:    "Tidak ada merchant terflag hari ini."
                  (bukan "No data found" atau "Ups, belum ada data!")
Scan not found:   "Merchant tidak ada dalam database. Tambah ke watchlist?"
API error:        "Koneksi terputus. Coba lagi atau hubungi admin."
```

---

## 9. Responsive Breakpoints

```css
/* Mobile first */
@media (max-width: 768px) {
  .hero-headline { font-size: 48px; }
  .hero-layout   { flex-direction: column; }
  .hero-illustration { display: none; } /* prioritas konten */
  .dashboard-grid { grid-template-columns: 1fr; }
  .stat-row { grid-template-columns: 1fr 1fr; }
  .network-graph { height: 280px; }
}

@media (min-width: 769px) and (max-width: 1024px) {
  .hero-headline { font-size: 64px; }
  .dashboard-grid { grid-template-columns: 1fr; }
}

@media (min-width: 1025px) {
  /* Full desktop layout */
  .hero-layout { display: grid; grid-template-columns: 45% 55%; }
  .dashboard-grid { grid-template-columns: 1fr 300px; }
}
```

---

## 10. Anti-Patterns — Jangan Lakukan

```
✗ Card dengan drop-shadow — ini bukan Material Design
✗ Gradient button (linear-gradient) — flat only
✗ Terlalu banyak rounded corner (>8px) — ini bukan fintech consumer app
✗ Emoji di UI — tidak ada satu pun emoji di luar demo data
✗ Animasi bounce/spring — terlalu playful
✗ Warna biru yang terlalu terang (#00AAFF) — bukan Klein Blue
✗ Font weight campuran dalam satu line — pilih satu weight per element
✗ Prose teks dalam JetBrains Mono — mono hanya untuk data/label/code
✗ Blue background untuk card di dashboard — blue hanya untuk hero
✗ Lebih dari 3 warna dalam satu komponen
```

---

## 11. File Structure

```
src/
├── styles/
│   ├── tokens.css        ← semua CSS variables dari section 2
│   ├── typography.css    ← font import + type classes
│   ├── components.css    ← semua component styles
│   └── animations.css    ← keyframes
├── components/
│   ├── Hero.jsx          ← blue section, full-bleed
│   ├── Nav.jsx           ← sticky, translucent blue
│   ├── NetworkGraph.jsx  ← canvas animation
│   ├── StatCard.jsx
│   ├── RiskBadge.jsx
│   ├── MerchantDetail.jsx
│   ├── MerchantChecker.jsx
│   ├── ReportForm.jsx
│   └── Illustration.jsx  ← SVG engraving motif
├── pages/
│   ├── Dashboard.jsx
│   ├── CheckMerchant.jsx
│   └── Reports.jsx
└── App.jsx
```

---

## 12. Hackathon Priority Queue

Waktu terbatas? Build dalam urutan ini:

```
PRIORITY 1 — WAJIB (tanpa ini: tidak presentable)
  ├── Hero section blue full-bleed + headline
  ├── Nav sticky translucent
  ├── Stat cards (4 buah)
  └── Risk badge component

PRIORITY 2 — PENTING (ini yang bikin WOW)
  ├── Network graph canvas (animated)
  ├── Merchant detail panel
  └── Merchant checker (input → result)

PRIORITY 3 — POLISH (kalau ada sisa waktu)
  ├── SVG illustration/motif di hero
  ├── Count-up animation untuk stats
  ├── Report form + success state
  └── Mobile responsive

PRIORITY 4 — NICE TO HAVE
  ├── Page load animation
  ├── PDF export button (fake OK untuk demo)
  └── Dark/light mode toggle
```

---

*QFIS Design System v1.0 — Klein Blue Authority*
*For hackathon use. Build fast, ship sharp.*
