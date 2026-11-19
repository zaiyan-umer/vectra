---
name: Vectra Engineered Precision
colors:
  surface: '#f8f9ff'
  surface-dim: '#cbdbf5'
  surface-bright: '#f8f9ff'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#eff4ff'
  surface-container: '#e5eeff'
  surface-container-high: '#dce9ff'
  surface-container-highest: '#d3e4fe'
  on-surface: '#0b1c30'
  on-surface-variant: '#434655'
  inverse-surface: '#213145'
  inverse-on-surface: '#eaf1ff'
  outline: '#737686'
  outline-variant: '#c3c6d7'
  surface-tint: '#0053db'
  primary: '#004ac6'
  on-primary: '#ffffff'
  primary-container: '#2563eb'
  on-primary-container: '#eeefff'
  inverse-primary: '#b4c5ff'
  secondary: '#565e74'
  on-secondary: '#ffffff'
  secondary-container: '#dae2fd'
  on-secondary-container: '#5c647a'
  tertiary: '#006242'
  on-tertiary: '#ffffff'
  tertiary-container: '#007d55'
  on-tertiary-container: '#bdffdb'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#dbe1ff'
  primary-fixed-dim: '#b4c5ff'
  on-primary-fixed: '#00174b'
  on-primary-fixed-variant: '#003ea8'
  secondary-fixed: '#dae2fd'
  secondary-fixed-dim: '#bec6e0'
  on-secondary-fixed: '#131b2e'
  on-secondary-fixed-variant: '#3f465c'
  tertiary-fixed: '#6ffbbe'
  tertiary-fixed-dim: '#4edea3'
  on-tertiary-fixed: '#002113'
  on-tertiary-fixed-variant: '#005236'
  background: '#f8f9ff'
  on-background: '#0b1c30'
  surface-variant: '#d3e4fe'
typography:
  headline-xl:
    fontFamily: Inter
    fontSize: 2rem
    fontWeight: '600'
    lineHeight: 2.5rem
    letterSpacing: -0.025em
  headline-xl-mobile:
    fontFamily: Inter
    fontSize: 1.5rem
    fontWeight: '600'
    lineHeight: 2rem
    letterSpacing: -0.02em
  headline-lg:
    fontFamily: Inter
    fontSize: 1.5rem
    fontWeight: '600'
    lineHeight: 2rem
    letterSpacing: -0.02em
  headline-md:
    fontFamily: Inter
    fontSize: 1.125rem
    fontWeight: '600'
    lineHeight: 1.5rem
    letterSpacing: -0.015em
  headline-sm:
    fontFamily: Inter
    fontSize: 0.875rem
    fontWeight: '600'
    lineHeight: 1.25rem
    letterSpacing: -0.005em
  body-lg:
    fontFamily: Inter
    fontSize: 1rem
    fontWeight: '400'
    lineHeight: 1.5rem
  body-md:
    fontFamily: Inter
    fontSize: 0.875rem
    fontWeight: '400'
    lineHeight: 1.25rem
  body-sm:
    fontFamily: Inter
    fontSize: 0.75rem
    fontWeight: '400'
    lineHeight: 1rem
  code-lg:
    fontFamily: JetBrains Mono
    fontSize: 1.125rem
    fontWeight: '500'
    lineHeight: 1.5rem
    letterSpacing: -0.02em
  code-md:
    fontFamily: JetBrains Mono
    fontSize: 0.875rem
    fontWeight: '400'
    lineHeight: 1.25rem
    letterSpacing: -0.01em
  code-sm:
    fontFamily: JetBrains Mono
    fontSize: 0.75rem
    fontWeight: '400'
    lineHeight: 1rem
    letterSpacing: '0'
  label-caps:
    fontFamily: JetBrains Mono
    fontSize: 0.6875rem
    fontWeight: '600'
    lineHeight: 0.875rem
    letterSpacing: 0.05em
rounded:
  sm: 0.125rem
  DEFAULT: 0.25rem
  md: 0.375rem
  lg: 0.5rem
  xl: 0.75rem
  full: 9999px
spacing:
  gutter: 1rem
  gutter-desktop: 1.5rem
  margin: 1rem
  margin-desktop: 2rem
  space-xs: 0.25rem
  space-sm: 0.5rem
  space-md: 0.75rem
  space-lg: 1rem
  space-xl: 1.5rem
---

## Brand & Style

The design system establishes a high-performance, developer-first operational atmosphere tailored for cloud-infrastructure architects, Site Reliability Engineers (SREs), and network performance engineers. The interface prioritizes density, immediate telemetry visibility, and deterministic interaction over decorative fluff.

### Visual Style & Methodology
- **Aesthetic Movement:** Technical Precision Minimalism. Stripped of non-functional ornament, relying on razor-sharp borders, exact alignment grids, and systematic spatial rhythms.
- **Atmospheric Intent:** Evokes the disciplined confidence of mission-critical aircraft flight displays and high-frequency trading terminals. The user should feel in control of an ultra-responsive, zero-latency instrument.
- **Structural Tenet:** Content precedes decor. Every pixel, badge, and grid line corresponds to operational metrics, throughput status, or distribution histograms.

## Colors

The color system delivers high legibility and unambiguous state signaling under high information density.

### Palette Architecture
- **Base Canvas & Cards:** Default light backdrop grounded in slate-50 (`#f8fafc`). Elevated operational planes and cards use pure white (`#ffffff`) framed by hairline dividers (`#e2e8f0`).
- **Primary Kinetic Accent:** Electric Blue (`#2563eb`, active hover `#1d4ed8`) reserved strictly for user-driven mutations, execution triggers (e.g., "Run Benchmark"), and primary navigation states.
- **Secondary Structural Neutral:** Deep Slate (`#0f172a` down to `#1e293b`) for dominant typography, metric value readouts, and active tab indicators.
- **Semantic Triad:**
  - **Healthy / Optimal:** Emerald (`#10b981`) paired with tint background (`#ecfdf5`) for successful HTTP 2xx distributions, latency p99 thresholds met, and idle node health.
  - **Degraded / Running:** Amber (`#f59e0b`) paired with tint background (`#fffbeb`) for active test runs, p95 buffer pressure, and non-fatal socket retransmits.
  - **Critical / Error:** Rose Red (`#ef4444`) paired with tint background (`#fef2f2`) for HTTP 5xx failures, connection timeouts, and dropped connections.

## Typography

The type scale integrates the UI neutrality of Inter alongside the monospaced mechanical precision of JetBrains Mono.

### Hierarchical Roles
- **Inter (Display & Body):** Employs tight negative tracking on headlines (`-0.025em`) to consolidate visual weight. Body text maintains natural tracking for comfortable scan-reading across dense table rows and configuration specs.
- **JetBrains Mono (Data & Code):** Mandated for all dynamic values, latency figures (e.g., `12.4ms`), RPS throughput counters, IP targets, and raw curl payloads. Enables columnar visual scanning without digit jitter via tabular figures.
- **Label-Caps (`0.6875rem` / `11px`):** Utilized exclusively for chart axis markers, column headers, protocol chips, and micro metadata tags. Always uppercase.

## Layout & Spacing

A 12-column grid layout provides deterministic structure for real-time telemetry panels, timeline charts, and parameter input panels.

### Layout Rules
- **Desktop (>= 1280px):** 12-column fluid grid locked within a maximum 1600px wrapper. Metrics and distribution summaries utilize 3-column spans; dynamic waterfall charts and response distributions span 8 to 12 columns. Fixed sidebar width at 240px.
- **Tablet (768px - 1279px):** Reflows into an 8-column layout. Metric cards collapse to 4-column spans (2x2 grid). Gutter tightens to `1rem`.
- **Mobile (< 768px):** Single-column stacked stream. Complex charts switch to summary stat rows with progressive disclosure drawers. Margins scale to `1rem`.
- **Spatial Cadence:** Standard component padding is strictly bound to `0.75rem` (`space-md`) and `1rem` (`space-lg`) to enforce a compact, tool-like footprint that optimizes vertical viewport capacity.

## Elevation & Depth

Visual hierarchy in this design system avoids heavy shadows and deep blurs in favor of crisp surface tiers and hairline stroke boundaries.

### Depth Stratification
- **Ground Floor (Canvas):** `#f8fafc` — host surface for primary layout scaffolding.
- **Tier 1 (Surface Cards & Toolbars):** Pure `#ffffff` framed with a continuous `1px` border of `#e2e8f0`. Shadow is minimal: `0 1px 2px 0 rgba(15, 23, 42, 0.04)`.
- **Tier 2 (Popovers, Filter Menus, Code Drawers):** `#ffffff` layered over a crisp border of `#cbd5e1` with a focused elevation drop: `0 4px 6px -1px rgba(15, 23, 42, 0.08), 0 2px 4px -2px rgba(15, 23, 42, 0.04)`.
- **Tier 3 (Modals & Command Palette):** Centered floating layers with an ambient, cool-gray tinted drop: `0 20px 25px -5px rgba(15, 23, 42, 0.1), 0 8px 10px -6px rgba(15, 23, 42, 0.04)`. Accompanied by a subtle backdrop scrim of `#0f172a` at 20% opacity.

## Shapes

The design system maintains a structured, low-radius contour profile (`roundedness: 1`), communicating architectural precision rather than playful casualness.

### Corner Radius System
- **Base (0.25rem / 4px):** Standard interactive elements including input fields, standard buttons, badge containers, and code tags.
- **Large (0.5rem / 8px):** Metric surface cards, log panels, visual charts, and modal frames.
- **Extra-Large (0.75rem / 12px):** Top-level floating command interfaces.
- **Zero Radius (0px):** Code blocks within tables, split data gutters, and raw timeline graphs to retain straight edge-to-edge alignment.

## Components

### Buttons & Trigger Controls
- **Primary:** Background `#2563eb`, text `#ffffff`, border `1px solid transparent`. Hover `#1d4ed8`, active `#1e40af`. Focus ring: `2px` offset with `#93c5fd`.
- **Secondary:** Background `#ffffff`, text `#0f172a`, border `1px solid #e2e8f0`. Hover `#f8fafc`, active `#f1f5f9`.
- **Destructive:** Background `#ffffff`, text `#ef4444`, border `1px solid #fecaca`. Hover `#fef2f2`.
- **Size Specification:** Height is compact (32px default, 28px dense table actions), font Inter 13px/weight 500.

### Status Badges & Chips
- Status indicators pair a subtle tinted fill with a hard-coded semantic colored dot:
  - **Healthy:** Background `#ecfdf5`, border `#a7f3d0`, text `#065f46`, inner pulse dot `#10b981`.
  - **Running/Warm:** Background `#fffbeb`, border `#fde68a`, text `#92400e`, inner pulse dot `#f59e0b`.
  - **Failed:** Background `#fef2f2`, border `#fecaca`, text `#991b1b`, inner dot `#ef4444`.
- Typography is strictly `label-caps` in JetBrains Mono.

### Data Tables & Log Lists
- Header cells: height 36px, background `#f8fafc`, typography `label-caps` in `#64748b`, bottom border `1px solid #e2e8f0`.
- Data rows: height 44px (collapsible to 36px under dense mode), background `#ffffff`, hover row background `#f8fafc`.
- Divider: hairline `1px solid #f1f5f9`. Numbers and endpoints styled using `code-sm`.

### Input Fields & Selects
- Base: height 34px, background `#ffffff`, border `1px solid #cbd5e1`, radius `0.25rem`. Text: Inter 13px (`#0f172a`), placeholder `#94a3b8`.
- Focus state: border `#2563eb`, box-shadow `0 0 0 1px #2563eb`.
- Monospace inputs (e.g., target URI, regex patterns): JetBrains Mono 12px with a fixed prefix tag (e.g., `https://`).

### Telemetry Cards
- Surface: `#ffffff`, border `1px solid #e2e8f0`, padding `1rem`.
- Layout: Top header displaying metric name in uppercase JetBrains Mono (`label-caps`), secondary change delta indicator (+/- %), and large primary value in `headline-lg` (interlocked with `code-lg`). Embedded sparklines span 100% of card bottom edge without padding.