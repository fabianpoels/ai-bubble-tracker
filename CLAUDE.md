# AI Bubble Tracker — Project Briefing for Claude Code

This file gives you full context on the project so you can contribute effectively without needing prior conversation history.

---

## What This Project Is

A web application that displays a single animated value: **how close we are to an AI market bubble burst**. The indicator is a composite score derived from publicly available numerical data across five weighted pillars. The goal is a minimal, credible, transparent tool — not a trading platform.

There are two types of signals:
- **Pre-burst indicators** — measure bubble proximity (the main score)
- **Post-burst indicators** — confirm when a burst has actually occurred

---

## Core Concept: The Five Pillars

The bubble proximity score is a weighted composite:

| Pillar | Weight | What It Measures |
|--------|--------|-----------------|
| Valuation Disconnect | 30% | Gap between AI stock valuations and fundamentals |
| CapEx-to-Revenue Chasm | 25% | How much is being spent vs. what's being earned |
| Enterprise Reality Check | 20% | Actual ROI realization, project success rates |
| Market Sentiment & Behavior | 15% | Retail interest, volatility, insider activity |
| External Shock Vulnerability | 10% | Macro fragility (rates, dollar, spreads) |

**Known bias to be aware of:** This framework overweights market signals because they're easy to obtain, while potentially underweighting harder-to-measure factors like energy constraints, competitive commoditization, and real usage metrics. Transparency about this bias is a design goal.

---

## Indicator Roadmap

Phases are organized by **data collection complexity**, not by importance.

### Phase 1 — Direct API or stable programmatic fetch (17 indicators)
All fetched from official APIs or well-maintained data services. No scraping.

| # | Indicator | Source | Frequency |
|---|-----------|--------|-----------|
| 1 | Shiller PE Ratio (CAPE) | FRED API | Monthly |
| 2 | Federal Funds Rate | FRED API | Daily |
| 3 | 10-Year Treasury Yield | FRED API | Daily |
| 4 | Corporate Bond Spread | FRED API (ICE BofA series) | Daily |
| 5 | VIX | Yahoo Finance `^VIX` | Daily |
| 6 | NVDA Price-to-Sales | Yahoo Finance | Daily |
| 7 | S&P 500 Forward P/E | Yahoo Finance | Daily |
| 8 | Individual AI Stock P/S (MSFT, GOOGL, META, PLTR, AVGO, ORCL) | Yahoo Finance | Daily |
| 9 | Market Concentration Index (top 7 tech / S&P 500 market cap) | Yahoo Finance (calculated) | Daily |
| 10 | NVDA Data Center Revenue | Yahoo Finance | Quarterly |
| 11 | Big Tech Aggregate CapEx (MSFT, GOOGL, META, AMZN, AAPL) | Yahoo Finance | Quarterly |
| 12 | CapEx-to-Revenue Ratio | Derived from #10 and #11 | Quarterly |
| 13 | Put/Call Ratio | Yahoo Finance / CBOE | Daily |
| 14 | US Dollar Index (DXY) | Yahoo Finance `DX-Y.NYB` | Daily |
| 15 | Bitcoin Price | Yahoo Finance `BTC-USD` | Daily |
| 16 | Google Trends "AI Bubble" | pytrends | Weekly |
| 17 | Google Trends "AI Stock" | pytrends | Weekly |

### Phase 2 — Scraping, relatively stable (3 indicators)
Structured HTML sources that are unlikely to change frequently.

| # | Indicator | Source | Frequency |
|---|-----------|--------|-----------|
| 18 | Insider Net Selling | openinsider.com | Weekly |
| 19 | Tech Layoff Count | layoffs.fyi | Weekly |
| 20 | ChatGPT Traffic Rank | SimilarWeb public pages | Monthly |

### Phase 3 — Manual or unreliable (2 indicators)
Low-frequency or no machine-readable source available.

| # | Indicator | Source | Frequency |
|---|-----------|--------|-----------|
| 21 | H1B AI Engineer Salaries | DOL bulk data (annual, heavy parsing) | Annual |
| 22 | OpenAI API Pricing | Manual monitoring | Quarterly |

---

## Tech Stack

| Component | Choice | Why |
|-----------|--------|-----|
| Language | Go 1.26 | Performance, simplicity |
| HTTP Framework | Gin | Lightweight, fast |
| Database | PostgreSQL | Relational, reliable |
| DB Driver | pgx | Modern, PostgreSQL-native, actively maintained |
| ORM | Bun | Lightweight, SQL-first, better performance than GORM |
| Cache | Redis | Standard caching layer |
| Dev Environment | Devcontainer | Reproducible, Docker-based |
| Base Image | `golang:1.26.0-bookworm` (official) | Minimal, no Microsoft devcontainer overhead |

---

## Architecture

### Entry Point Pattern
The application uses a **CLI-based task routing system** via Go flags, not separate binaries:

```bash
go run main.go -task server       # Start the API server
go run main.go -task db-create    # Create database tables
go run main.go -task db-drop      # Drop tables (with confirmation prompt)
go run main.go -task fetch-<name> # Fetch data from a specific external API
```

Each external data source gets its own `-task fetch-X` entry point.

### Database Migrations
Handled **programmatically through Go code**, not CLI migration tools (no golang-migrate binary). Bun ORM's `CreateTable` / `DropTable` is used directly in task handlers.

### Key Data Model

One row per day — a wide table that holds all raw indicator values, pillar scores, and the final computed bubble index. This simplifies time-series queries and gap detection (NULLs are meaningful: they indicate data not yet fetched or unavailable for that date).

```go
type Datapoint struct {
    bun.BaseModel `bun:"table:datapoints,alias:d"`

    Date int `bun:"date,pk,notnull" json:"date"` // YYYYMMDD — e.g. 20260215

    // Pillar 1: Valuation Disconnect (30%)
    ShillerPE           *float64 `bun:"shiller_pe" json:"shiller_pe"`
    SP500ForwardPE      *float64 `bun:"sp500_forward_pe" json:"sp500_forward_pe"`
    NvdaPriceToSales    *float64 `bun:"nvda_price_to_sales" json:"nvda_price_to_sales"`
    MarketConcentration *float64 `bun:"market_concentration" json:"market_concentration"`

    // Pillar 2: CapEx-to-Revenue Chasm (25%)
    BigTechCapex        *float64 `bun:"big_tech_capex" json:"big_tech_capex"`
    NvdaDataCenterRev   *float64 `bun:"nvda_data_center_rev" json:"nvda_data_center_rev"`
    CapexToRevenueRatio *float64 `bun:"capex_to_revenue_ratio" json:"capex_to_revenue_ratio"`

    // Pillar 3: Enterprise Reality Check (20%)
    // Phase 2+ indicators — columns to be added as implemented

    // Pillar 4: Market Sentiment & Behavior (15%)
    Vix                  *float64 `bun:"vix" json:"vix"`
    GoogleTrendsAIBubble *float64 `bun:"google_trends_ai_bubble" json:"google_trends_ai_bubble"`
    InsiderNetSelling    *float64 `bun:"insider_net_selling" json:"insider_net_selling"`

    // Pillar 5: External Shock Vulnerability (10%)
    FedFundsRate  *float64 `bun:"fed_funds_rate" json:"fed_funds_rate"`
    TenYearYield  *float64 `bun:"ten_year_yield" json:"ten_year_yield"`

    // Computed scores (populated by aggregation task)
    ValuationScore *float64 `bun:"valuation_score" json:"valuation_score"`
    CapexScore     *float64 `bun:"capex_score" json:"capex_score"`
    SentimentScore *float64 `bun:"sentiment_score" json:"sentiment_score"`
    MacroScore     *float64 `bun:"macro_score" json:"macro_score"`
    BubbleIndex    *float64 `bun:"bubble_index" json:"bubble_index"`

    CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
    UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}
```

**Why pointer types (`*float64`):** A regular `float64` cannot distinguish between "value is zero" and "value was never fetched". Pointers map to true PostgreSQL NULLs, making gap detection reliable.

**Why integer dates as PK:** Better index performance than strings, no timezone ambiguity, trivially sortable and comparable, human-readable.

**Why a wide table over EAV (narrow) table:** Time-series queries are trivial (one row = one day's full snapshot). Gap detection is straightforward (query for NULLs in specific columns). The alternative — one row per date+indicator — is more flexible for schema changes but makes every query a pivot. Since the indicator set is planned and phased, schema migrations on a low-traffic internal table are an acceptable trade-off.

### Upsert Pattern

Each fetcher task only touches its own columns, leaving others untouched:

```go
_, err := db.GetBun().NewInsert().
    Model(&datapoint).
    On("CONFLICT (date) DO UPDATE").
    Set("vix = EXCLUDED.vix, updated_at = current_timestamp").
    Exec(ctx)
```

Never overwrite unrelated columns — each fetcher is scoped to its own `SET` clause.

### Historical Data

Fetchers should backfill historical data on first run, up to the depth available per source (e.g. Shiller PE back to 1871, VIX back to 1990, NVDA data center revenue ~5 years). Historical depth varies per indicator — NULLs for pre-existence periods are expected and correct.

### Helper Functions

```go
func DateToInt(t time.Time) int {
    year, month, day := t.Date()
    return year*10000 + int(month)*100 + day
}

func IntToDate(i int) time.Time {
    year := i / 10000
    month := (i % 10000) / 100
    day := i % 100
    return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
```

### Package Structure
```
.
├── .devcontainer/
│   ├── devcontainer.json   # Uses common-utils feature for user management
│   └── docker-compose.yml
├── cmd/ or main.go         # Entry point with flag-based task routing
├── db/                     # Bun connection management (GetBun(), Close())
├── models/                 # Bun model structs
├── tasks/                  # One file per task (db_create, db_drop, fetch_X)
├── server/                 # Gin router and handlers
└── .env                    # Local config (not committed)
```

---

## Development Principles

These are non-negotiable preferences — always apply them:

**Minimalism over convenience.** If two solutions solve the same problem, choose the simpler one. Challenge redundant dependencies before adding them. Every tool in the stack should earn its place.

**Understand the trade-off.** When recommending a library or approach, explain *why* — what is being gained and what is being given up. Don't just recommend the most popular option.

**No redundant tooling.** The project uses `gopls` via editor integration for code analysis. Don't suggest adding `golangci-lint` or similar tools unless there's a specific, concrete gap. Don't install things that are already covered elsewhere.

**SQL-first ORM.** Bun was chosen *specifically* because it's SQL-first. Don't suggest GORM or other magic-heavy ORMs. Write explicit queries when needed.

**Handle migrations in code.** Don't suggest adding golang-migrate or other CLI migration tools. Table management happens through Go task functions using Bun's schema API.

**Clean devcontainer setup.** The devcontainer uses the official `golang:1.26.0-bookworm` image with the `common-utils` devcontainer feature for non-root user management. Don't suggest switching to Microsoft's devcontainer Go image.

**Environment variables via `.env`.** All config comes from a `.env` file loaded with `godotenv`. No hardcoded config values.

---

## Data Collection Approach

Each external data source is a separate task (`-task fetch-X`). Data is stored as `Datapoint` records with the indicator name as part of the composite PK.

**Primary free API sources:**
- **FRED** (`api.stlouisfed.org`) — macro data, Shiller PE, Fed Funds Rate
- **Yahoo Finance** (`query1.finance.yahoo.com`) — stock prices, P/S ratios, VIX, CapEx
- **Google Trends** — sentiment via pytrends (or equivalent Go client)
- **SEC EDGAR** (`sec.gov/cgi-bin/browse-edgar`) — CapEx from filings
- **openinsider.com** — insider trading activity

Collection frequencies vary: daily for market data, weekly for trends, quarterly for earnings/CapEx.

---

## What Has Been Built So Far

- [x] Devcontainer setup (golang:1.26.0-bookworm + postgres + redis)
- [x] Go project initialized with Gin, pgx, Bun, godotenv
- [x] Database connection management (`db.GetBun()`, `db.Close()`)
- [x] `Datapoint` model with integer date composite PK
- [x] CLI task router in `main.go` (`-task` flag)
- [x] `db-create` and `db-drop` tasks
- [ ] First data fetch tasks (Phase 1 indicators)
- [ ] Gin API endpoints serving indicator data
- [ ] Scoring/aggregation logic
- [ ] Frontend

---

## What To Avoid

- Don't add libraries without a clear reason
- Don't use GORM — Bun is the ORM
- Don't suggest CLI migration tools — use Bun's schema API
- Don't use Microsoft's devcontainer Go image — use the official `golang:1.26.0-bookworm`
- Don't hardcode credentials or config values
- Don't add golangci-lint unless there's a specific unmet need
- Don't over-engineer early phases — get Phase 1 working cleanly first