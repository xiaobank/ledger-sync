package ledger

import (
	"fmt"
	"sort"
	"time"
)

// RollupPeriod defines the granularity of a rollup aggregation.
type RollupPeriod string

const (
	RollupDaily   RollupPeriod = "daily"
	RollupMonthly RollupPeriod = "monthly"
)

// RollupBucket holds aggregated transaction totals for a time bucket.
type RollupBucket struct {
	Label    string
	Period   RollupPeriod
	Currency string
	Total    int64 // in minor units
	Count    int
}

// RollupResult is the ordered set of buckets produced by RollupTransactions.
type RollupResult struct {
	AccountID string
	Period    RollupPeriod
	Buckets   []RollupBucket
}

// bucketLabel returns a string key for the given time and period.
func bucketLabel(t time.Time, p RollupPeriod) string {
	switch p {
	case RollupMonthly:
		return t.Format("2006-01")
	default:
		return t.Format("2006-01-02")
	}
}

// RollupTransactions aggregates posted transactions for a given account
// into time-bucketed totals grouped by currency.
func RollupTransactions(b *Book, accountID string, period RollupPeriod) (*RollupResult, error) {
	if b == nil {
		return nil, fmt.Errorf("rollup: book must not be nil")
	}
	if accountID == "" {
		return nil, fmt.Errorf("rollup: accountID must not be empty")
	}
	if period != RollupDaily && period != RollupMonthly {
		return nil, fmt.Errorf("rollup: unknown period %q", period)
	}

	type key struct {
		label    string
		currency string
	}
	agg := make(map[key]*RollupBucket)

	for _, tx := range b.Transactions() {
		for _, leg := range tx.Legs {
			if leg.AccountID != accountID {
				continue
			}
			k := key{label: bucketLabel(tx.PostedAt, period), currency: leg.Currency}
			if agg[k] == nil {
				agg[k] = &RollupBucket{
					Label:    k.label,
					Period:   period,
					Currency: leg.Currency,
				}
			}
			agg[k].Total += leg.Amount
			agg[k].Count++
		}
	}

	buckets := make([]RollupBucket, 0, len(agg))
	for _, v := range agg {
		buckets = append(buckets, *v)
	}
	sort.Slice(buckets, func(i, j int) bool {
		if buckets[i].Label == buckets[j].Label {
			return buckets[i].Currency < buckets[j].Currency
		}
		return buckets[i].Label < buckets[j].Label
	})

	return &RollupResult{
		AccountID: accountID,
		Period:    period,
		Buckets:   buckets,
	}, nil
}
