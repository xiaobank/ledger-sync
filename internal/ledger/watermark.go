package ledger

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// WatermarkKind distinguishes high-water from low-water marks.
type WatermarkKind string

const (
	WatermarkHigh WatermarkKind = "high"
	WatermarkLow  WatermarkKind = "low"
)

// Watermark records the peak or trough balance seen for an account.
type Watermark struct {
	AccountID string
	Kind      WatermarkKind
	Amount    float64
	Currency  string
	RecordedAt time.Time
}

// WatermarkIndex tracks high and low watermarks per account/currency pair.
type WatermarkIndex struct {
	mu    sync.RWMutex
	marks map[string]*Watermark // key: accountID+":"+currency+":"+kind
}

// NewWatermarkIndex creates an empty WatermarkIndex.
func NewWatermarkIndex() *WatermarkIndex {
	return &WatermarkIndex{
		marks: make(map[string]*Watermark),
	}
}

func watermarkKey(accountID, currency string, kind WatermarkKind) string {
	return fmt.Sprintf("%s:%s:%s", accountID, currency, kind)
}

// Observe updates the watermarks for the given account/currency with the
// supplied balance value, recording the timestamp if a new extreme is found.
func (w *WatermarkIndex) Observe(accountID, currency string, balance float64, at time.Time) error {
	if accountID == "" {
		return errors.New("watermark: accountID must not be empty")
	}
	if currency == "" {
		return errors.New("watermark: currency must not be empty")
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	for _, kind := range []WatermarkKind{WatermarkHigh, WatermarkLow} {
		key := watermarkKey(accountID, currency, kind)
		existing, ok := w.marks[key]
		if !ok {
			w.marks[key] = &Watermark{
				AccountID:  accountID,
				Kind:       kind,
				Amount:     balance,
				Currency:   currency,
				RecordedAt: at,
			}
			continue
		}
		if kind == WatermarkHigh && balance > existing.Amount {
			existing.Amount = balance
			existing.RecordedAt = at
		}
		if kind == WatermarkLow && balance < existing.Amount {
			existing.Amount = balance
			existing.RecordedAt = at
		}
	}
	return nil
}

// Get returns the watermark for the given account, currency, and kind.
// Returns nil if no observation has been recorded yet.
func (w *WatermarkIndex) Get(accountID, currency string, kind WatermarkKind) *Watermark {
	w.mu.RLock()
	defer w.mu.RUnlock()
	key := watermarkKey(accountID, currency, kind)
	wm, ok := w.marks[key]
	if !ok {
		return nil
	}
	copy := *wm
	return &copy
}

// All returns a slice of all recorded watermarks.
func (w *WatermarkIndex) All() []Watermark {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]Watermark, 0, len(w.marks))
	for _, wm := range w.marks {
		out = append(out, *wm)
	}
	return out
}
