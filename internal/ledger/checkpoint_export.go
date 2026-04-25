package ledger

import (
	"fmt"
	"sort"
	"strings"
)

// CheckpointSummary holds a human-readable summary row for a checkpoint.
type CheckpointSummary struct {
	Name      string
	CreatedAt string
	AccountID string
	Currency  string
	Balance   float64
}

// ExportCheckpointCSV renders a checkpoint's balances as CSV.
func ExportCheckpointCSV(cp *Checkpoint) (string, error) {
	if cp == nil {
		return "", fmt.Errorf("checkpoint must not be nil")
	}

	var sb strings.Builder
	sb.WriteString("checkpoint,created_at,account_id,currency,balance\n")

	accIDs := make([]string, 0, len(cp.Balances))
	for id := range cp.Balances {
		accIDs = append(accIDs, id)
	}
	sort.Strings(accIDs)

	for _, accID := range accIDs {
		currencies := cp.Balances[accID]
		curs := make([]string, 0, len(currencies))
		for c := range currencies {
			curs = append(curs, c)
		}
		sort.Strings(curs)
		for _, cur := range curs {
			sb.WriteString(fmt.Sprintf("%s,%s,%s,%s,%.2f\n",
				cp.Name,
				cp.CreatedAt.Format("2006-01-02T15:04:05Z"),
				accID,
				cur,
				cp.Balances[accID][cur],
			))
		}
	}
	return sb.String(), nil
}

// ExportCheckpointDiffCSV renders a balance delta map as CSV.
func ExportCheckpointDiffCSV(from, to string, delta map[string]map[string]float64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("from,to,account_id,currency,delta\n"))

	accIDs := make([]string, 0, len(delta))
	for id := range delta {
		accIDs = append(accIDs, id)
	}
	sort.Strings(accIDs)

	for _, accID := range accIDs {
		currencies := delta[accID]
		curs := make([]string, 0, len(currencies))
		for c := range currencies {
			curs = append(curs, c)
		}
		sort.Strings(curs)
		for _, cur := range curs {
			sb.WriteString(fmt.Sprintf("%s,%s,%s,%s,%.2f\n",
				from, to, accID, cur, delta[accID][cur],
			))
		}
	}
	return sb.String()
}
