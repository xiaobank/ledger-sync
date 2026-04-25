package ledger

import (
	"fmt"
	"sync"
	"time"
)

// Checkpoint records a named point-in-time balance state for a book.
type Checkpoint struct {
	Name      string
	CreatedAt time.Time
	Balances  map[string]map[string]float64 // accountID -> currency -> balance
}

// CheckpointIndex stores named checkpoints for a book.
type CheckpointIndex struct {
	mu          sync.RWMutex
	checkpoints map[string]*Checkpoint
}

// NewCheckpointIndex creates an empty CheckpointIndex.
func NewCheckpointIndex() *CheckpointIndex {
	return &CheckpointIndex{
		checkpoints: make(map[string]*Checkpoint),
	}
}

// Capture records the current balances of all accounts in the book under the given name.
func (ci *CheckpointIndex) Capture(name string, b *Book) error {
	if name == "" {
		return fmt.Errorf("checkpoint name must not be empty")
	}
	if b == nil {
		return fmt.Errorf("book must not be nil")
	}

	balances := make(map[string]map[string]float64)
	for id, acc := range b.accounts {
		balances[id] = acc.balances()
	}

	cp := &Checkpoint{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Balances:  balances,
	}

	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.checkpoints[name] = cp
	return nil
}

// Get retrieves a checkpoint by name.
func (ci *CheckpointIndex) Get(name string) (*Checkpoint, bool) {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	cp, ok := ci.checkpoints[name]
	return cp, ok
}

// Names returns all stored checkpoint names.
func (ci *CheckpointIndex) Names() []string {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	names := make([]string, 0, len(ci.checkpoints))
	for n := range ci.checkpoints {
		names = append(names, n)
	}
	return names
}

// DiffCheckpoints returns balance deltas between two named checkpoints.
// Positive delta means balance increased from 'from' to 'to'.
func (ci *CheckpointIndex) DiffCheckpoints(from, to string) (map[string]map[string]float64, error) {
	ci.mu.RLock()
	defer ci.mu.RUnlock()

	a, ok := ci.checkpoints[from]
	if !ok {
		return nil, fmt.Errorf("checkpoint %q not found", from)
	}
	b, ok := ci.checkpoints[to]
	if !ok {
		return nil, fmt.Errorf("checkpoint %q not found", to)
	}

	delta := make(map[string]map[string]float64)
	for accID, currencies := range b.Balances {
		delta[accID] = make(map[string]float64)
		for cur, bal := range currencies {
			prev := 0.0
			if a.Balances[accID] != nil {
				prev = a.Balances[accID][cur]
			}
			delta[accID][cur] = bal - prev
		}
	}
	return delta, nil
}
