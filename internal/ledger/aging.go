package ledger

import (
	"errors"
	"sort"
	"time"
)

// AgingBucket represents a time-based grouping of outstanding balances.
type AgingBucket struct {
	Label   string
	MinDays int
	MaxDays int // -1 means unbounded
	Total   map[string]float64 // currency -> amount
}

// AgingReport holds buckets of aged balances for a given account.
type AgingReport struct {
	AccountID string
	AsOf      time.Time
	Buckets   []AgingBucket
}

var defaultBuckets = []struct {
	label   string
	min     int
	max     int
}{
	{"Current", 0, 30},
	{"31-60 days", 31, 60},
	{"61-90 days", 61, 90},
	{"91+ days", 91, -1},
}

// AgingAnalysis computes an aging report for the given account by examining
// posted transactions in the book relative to asOf.
func AgingAnalysis(b *Book, accountID string, asOf time.Time) (*AgingReport, error) {
	if b == nil {
		return nil, errors.New("aging: book is nil")
	}
	if accountID == "" {
		return nil, errors.New("aging: accountID is empty")
	}

	buckets := make([]AgingBucket, len(defaultBuckets))
	for i, d := range defaultBuckets {
		buckets[i] = AgingBucket{
			Label:   d.label,
			MinDays: d.min,
			MaxDays: d.max,
			Total:   make(map[string]float64),
		}
	}

	txns := b.Transactions()
	sort.Slice(txns, func(i, j int) bool {
		return txns[i].Date.Before(txns[j].Date)
	})

	for _, tx := range txns {
		if tx.Date.After(asOf) {
			continue
		}
		age := int(asOf.Sub(tx.Date).Hours() / 24)
		for _, leg := range tx.Legs {
			if leg.AccountID != accountID {
				continue
			}
			for i, d := range defaultBuckets {
				if age >= d.min && (d.max == -1 || age <= d.max) {
					buckets[i].Total[leg.Currency] += leg.Amount
					break
				}
			}
		}
	}

	return &AgingReport{
		AccountID: accountID,
		AsOf:      asOf,
		Buckets:   buckets,
	}, nil
}
