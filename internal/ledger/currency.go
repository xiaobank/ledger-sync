package ledger

import (
	"errors"
	"fmt"
	"strings"
)

// CurrencyCode represents an ISO 4217 currency code.
type CurrencyCode string

// ExchangeRate holds a conversion rate between two currencies.
type ExchangeRate struct {
	From CurrencyCode
	To   CurrencyCode
	Rate float64
}

// CurrencyRegistry stores known currencies and exchange rates.
type CurrencyRegistry struct {
	currencies map[CurrencyCode]bool
	rates      map[string]float64
}

// NewCurrencyRegistry creates an empty CurrencyRegistry.
func NewCurrencyRegistry() *CurrencyRegistry {
	return &CurrencyRegistry{
		currencies: make(map[CurrencyCode]bool),
		rates:      make(map[string]float64),
	}
}

// Register adds a currency code to the registry.
func (r *CurrencyRegistry) Register(code CurrencyCode) error {
	c := CurrencyCode(strings.ToUpper(string(code)))
	if len(c) != 3 {
		return fmt.Errorf("invalid currency code %q: must be 3 characters", code)
	}
	r.currencies[c] = true
	return nil
}

// IsKnown reports whether the given currency code has been registered.
func (r *CurrencyRegistry) IsKnown(code CurrencyCode) bool {
	return r.currencies[CurrencyCode(strings.ToUpper(string(code)))]
}

// SetRate stores an exchange rate from one currency to another.
func (r *CurrencyRegistry) SetRate(from, to CurrencyCode, rate float64) error {
	if rate <= 0 {
		return errors.New("exchange rate must be positive")
	}
	if !r.IsKnown(from) {
		return fmt.Errorf("unknown currency: %s", from)
	}
	if !r.IsKnown(to) {
		return fmt.Errorf("unknown currency: %s", to)
	}
	r.rates[rateKey(from, to)] = rate
	return nil
}

// Convert converts an amount from one currency to another using stored rates.
func (r *CurrencyRegistry) Convert(amount int64, from, to CurrencyCode) (int64, error) {
	if from == to {
		return amount, nil
	}
	rate, ok := r.rates[rateKey(from, to)]
	if !ok {
		return 0, fmt.Errorf("no exchange rate from %s to %s", from, to)
	}
	return int64(float64(amount) * rate), nil
}

func rateKey(from, to CurrencyCode) string {
	return fmt.Sprintf("%s->%s", strings.ToUpper(string(from)), strings.ToUpper(string(to)))
}
