package ledger

import (
	"fmt"
	"time"
)

// Snapshot captures the state of all account balances at a point in time.
type Snapshot struct {
	TakenAt  time.Time
	BookName string
	Balances map[string]map[string]float64 // accountID -> currency -> balance
}

// TakeSnapshot records the current balance of every account in the book.
func TakeSnapshot(b *Book) (*Snapshot, error) {
	if b == nil {
		return nil, fmt.Errorf("snapshot: book must not be nil")
	}

	snap := &Snapshot{
		TakenAt:  time.Now().UTC(),
		BookName: b.Name,
		Balances: make(map[string]map[string]float64),
	}

	for id, acc := range b.accounts {
		currencies := acc.Currencies()
		if len(currencies) == 0 {
			snap.Balances[id] = map[string]float64{}
			continue
		}
		snap.Balances[id] = make(map[string]float64, len(currencies))
		for _, cur := range currencies {
			snap.Balances[id][cur] = acc.Balance(cur)
		}
	}

	return snap, nil
}

// DiffSnapshots compares two snapshots and returns accounts whose balances
// changed between them. Only accounts present in at least one snapshot are
// included in the diff.
func DiffSnapshots(before, after *Snapshot) map[string]map[string][2]float64 {
	diff := make(map[string]map[string][2]float64)

	allIDs := make(map[string]struct{})
	for id := range before.Balances {
		allIDs[id] = struct{}{}
	}
	for id := range after.Balances {
		allIDs[id] = struct{}{}
	}

	for id := range allIDs {
		beforeCurs := before.Balances[id]
		afterCurs := after.Balances[id]

		allCurs := make(map[string]struct{})
		for c := range beforeCurs {
			allCurs[c] = struct{}{}
		}
		for c := range afterCurs {
			allCurs[c] = struct{}{}
		}

		for c := range allCurs {
			b := beforeCurs[c]
			a := afterCurs[c]
			if b != a {
				if diff[id] == nil {
					diff[id] = make(map[string][2]float64)
				}
				diff[id][c] = [2]float64{b, a}
			}
		}
	}

	return diff
}
