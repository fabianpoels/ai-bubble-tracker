package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Datapoint struct {
	bun.BaseModel `bun:"table:datapoints,alias:d"`

	Date string `bun:"date,pk,notnull" json:"date"`

	// Your bubble indicator metrics
	BubbleIndex float64 `bun:"bubble_index" json:"bubble_index"`
	// Add your specific metrics here, for example:
	// ValuationScore      float64   `bun:"valuation_score" json:"valuation_score"`
	// CapexRevenueRatio   float64   `bun:"capex_revenue_ratio" json:"capex_revenue_ratio"`
	// SentimentIndex      float64   `bun:"sentiment_index" json:"sentiment_index"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

func DateToString(t time.Time) string {
	return t.Format("20060102")
}

func StringToDate(s string) (time.Time, error) {
	return time.Parse("20060102", s)
}
