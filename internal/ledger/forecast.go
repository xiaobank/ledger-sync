package ledger

import (
	"fmt"
	"time"
)

// ForecastEntry represents a projected balance for an account at a future date.
type ForecastEntry struct {
	AccountID string
	Currency  string
	Projected float64
	At        time.Time
}

// ForecastResult holds the full set of projected entries for a forecast run.
type ForecastResult struct {
	GeneratedAt time.Time
	Entries     []ForecastEntry
}

// ForecastOptions controls how the forecast is computed.
type ForecastOptions struct {
	// Periods is the number of future periods to project.
	Periods int
	// PeriodDuration is the length of each period (e.g. 30 * 24 * time.Hour).
	PeriodDuration time.Duration
	// GrowthRate is a per-period multiplier applied to each account balance (e.g. 0.02 for 2%).
	GrowthRate float64
}

// Forecast projects account balances forward based on current balances and a
// simple compound growth model. It returns an error if the book is nil or
// options are invalid.
func Forecast(b *Book, opts ForecastOptions) (*ForecastResult, error) {
	if b == nil {
		return nil, fmt.Errorf("forecast: book must not be nil")
	}
	if opts.Periods <= 0 {
		return nil, fmt.Errorf("forecast: periods must be greater than zero")
	}
	if opts.PeriodDuration <= 0 {
		return nil, fmt.Errorf("forecast: period duration must be positive")
	}

	now := time.Now().UTC()
	result := &ForecastResult{GeneratedAt: now}

	b.mu.RLock()
	accounts := make([]*Account, 0, len(b.accounts))
	for _, a := range b.accounts {
		accounts = append(accounts, a)
	}
	b.mu.RUnlock()

	for _, acc := range accounts {
		for currency, balance := range acc.balances {
			projected := balance
			for p := 1; p <= opts.Periods; p++ {
				projected = projected * (1.0 + opts.GrowthRate)
				at := now.Add(time.Duration(p) * opts.PeriodDuration)
				result.Entries = append(result.Entries, ForecastEntry{
					AccountID: acc.ID,
					Currency:  currency,
					Projected: projected,
					At:        at,
				})
			}
		}
	}

	return result, nil
}
