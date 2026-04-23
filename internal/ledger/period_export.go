package ledger

import (
	"fmt"
	"strings"
)

// PeriodSummary holds aggregated transaction data for a single period.
type PeriodSummary struct {
	Period Period
	Count  int
	Total  map[string]float64 // currency -> net amount
}

// SummariseByPeriod groups transactions from book into the supplied periods
// for the given accountID. Pass an empty accountID to include all accounts.
func SummariseByPeriod(book *Book, periods []Period, accountID string) ([]PeriodSummary, error) {
	if book == nil {
		return nil, fmt.Errorf("period_export: book is nil")
	}
	if len(periods) == 0 {
		return nil, fmt.Errorf("period_export: no periods supplied")
	}

	summaries := make([]PeriodSummary, len(periods))
	for i, p := range periods {
		summaries[i] = PeriodSummary{
			Period: p,
			Total:  make(map[string]float64),
		}
	}

	for _, tx := range book.Transactions() {
		for i, p := range periods {
			if !p.Contains(tx.Timestamp) {
				continue
			}
			for _, leg := range tx.Legs {
				if accountID != "" && leg.AccountID != accountID {
					continue
				}
				summaries[i].Total[leg.Currency] += leg.Amount
				summaries[i].Count++
			}
		}
	}
	return summaries, nil
}

// ExportPeriodSummaryCSV renders a slice of PeriodSummary as CSV text.
func ExportPeriodSummaryCSV(summaries []PeriodSummary) string {
	var sb strings.Builder
	sb.WriteString("period,currency,count,total\n")
	for _, s := range summaries {
		if len(s.Total) == 0 {
			sb.WriteString(fmt.Sprintf("%s,,0,0.00\n", s.Period.Label))
			continue
		}
		for currency, total := range s.Total {
			sb.WriteString(fmt.Sprintf("%s,%s,%d,%.2f\n",
				s.Period.Label, currency, s.Count, total))
		}
	}
	return sb.String()
}
