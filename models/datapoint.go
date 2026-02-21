package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Datapoint struct {
	bun.BaseModel `bun:"table:datapoints,alias:d"`

	Date int `bun:"date,pk,notnull" json:"date"`

	ShillerPE            *float64 `bun:"shiller_pe" json:"shiller_pe"`
	SP500ForwardPE       *float64 `bun:"sp500_forward_pe" json:"sp500_forward_pe"`
	NvdaPriceToSales     *float64 `bun:"nvda_price_to_sales" json:"nvda_price_to_sales"`
	MarketConcentration  *float64 `bun:"market_concentration" json:"market_concentration"`
	BigTechCapex         *float64 `bun:"big_tech_capex" json:"big_tech_capex"`
	NvdaDataCenterRev    *float64 `bun:"nvda_data_center_rev" json:"nvda_data_center_rev"`
	CapexToRevenueRatio  *float64 `bun:"capex_to_revenue_ratio" json:"capex_to_revenue_ratio"`
	Vix                  *float64 `bun:"vix" json:"vix"`
	GoogleTrendsAIBubble *float64 `bun:"google_trends_ai_bubble" json:"google_trends_ai_bubble"`
	InsiderNetSelling    *float64 `bun:"insider_net_selling" json:"insider_net_selling"`
	FedFundsRate         *float64 `bun:"fed_funds_rate" json:"fed_funds_rate"`
	TenYearYield         *float64 `bun:"ten_year_yield" json:"ten_year_yield"`

	ValuationScore *float64 `bun:"valuation_score" json:"valuation_score"`
	CapexScore     *float64 `bun:"capex_score" json:"capex_score"`
	SentimentScore *float64 `bun:"sentiment_score" json:"sentiment_score"`
	MacroScore     *float64 `bun:"macro_score" json:"macro_score"`
	BubbleIndex    *float64 `bun:"bubble_index" json:"bubble_index"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

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
