package ledger

import (
	"errors"
	"time"
)

// ProjectionEntry holds a single projected balance for an account at a point in time.
type ProjectionEntry struct {
	AccountID string
	At        time.Time
	Balance   map[string]float64 // currency -> projected balance
}

// ProjectionResult contains all projected entries produced by RunProjection.
type ProjectionResult struct {
	Entries []ProjectionEntry
}

// ProjectionConfig controls how a balance projection is computed.
type ProjectionConfig struct {
	AccountID    string
	Currency     string
	StepDuration time.Duration // e.g. 24*time.Hour for daily
	Steps        int           // number of steps to project
	GrowthRate   float64       // fractional growth per step, e.g. 0.01 = 1%
}

// RunProjection projects the future balance of a single account by applying
// a constant growth rate over a series of time steps.
func RunProjection(b *Book, cfg ProjectionConfig) (*ProjectionResult, error) {
	if b == nil {
		return nil, errors.New("projection: book must not be nil")
	}
	if cfg.AccountID == "" {
		return nil, errors.New("projection: account ID must not be empty")
	}
	if cfg.Currency == "" {
		return nil, errors.New("projection: currency must not be empty")
	}
	if cfg.StepDuration <= 0 {
		return nil, errors.New("projection: step duration must be positive")
	}
	if cfg.Steps <= 0 {
		return nil, errors.New("projection: steps must be positive")
	}

	acc, err := b.GetAccount(cfg.AccountID)
	if err != nil {
		return nil, errors.New("projection: account not found: " + cfg.AccountID)
	}

	currentBalance := acc.Balance(cfg.Currency)
	now := time.Now().UTC()

	result := &ProjectionResult{}
	for i := 1; i <= cfg.Steps; i++ {
		currentBalance = currentBalance * (1 + cfg.GrowthRate)
		entry := ProjectionEntry{
			AccountID: cfg.AccountID,
			At:        now.Add(time.Duration(i) * cfg.StepDuration),
			Balance:   map[string]float64{cfg.Currency: currentBalance},
		}
		result.Entries = append(result.Entries, entry)
	}
	return result, nil
}
