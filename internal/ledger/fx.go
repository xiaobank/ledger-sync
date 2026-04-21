package ledger

import (
	"errors"
	"fmt"
	"time"
)

// FXConversion represents a currency conversion applied to a transaction.
type FXConversion struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
	AppliedAt    time.Time
}

// FXConverter uses a CurrencyRegistry to convert amounts between currencies.
type FXConverter struct {
	registry *CurrencyRegistry
}

// NewFXConverter creates a new FXConverter backed by the given registry.
func NewFXConverter(registry *CurrencyRegistry) (*FXConverter, error) {
	if registry == nil {
		return nil, errors.New("fx: registry must not be nil")
	}
	return &FXConverter{registry: registry}, nil
}

// Convert converts an amount from one currency to another using the registry rates.
// Both currencies must be known and a rate must be set for toCurrency.
func (fx *FXConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, FXConversion, error) {
	if amount < 0 {
		return 0, FXConversion{}, errors.New("fx: amount must be non-negative")
	}
	if !fx.registry.IsKnown(fromCurrency) {
		return 0, FXConversion{}, fmt.Errorf("fx: unknown source currency %q", fromCurrency)
	}
	if !fx.registry.IsKnown(toCurrency) {
		return 0, FXConversion{}, fmt.Errorf("fx: unknown target currency %q", toCurrency)
	}
	if fromCurrency == toCurrency {
		conv := FXConversion{
			FromCurrency: fromCurrency,
			ToCurrency:   toCurrency,
			Rate:         1.0,
			AppliedAt:    time.Now().UTC(),
		}
		return amount, conv, nil
	}
	rate, err := fx.registry.GetRate(toCurrency)
	if err != nil {
		return 0, FXConversion{}, fmt.Errorf("fx: no rate available for %q: %w", toCurrency, err)
	}
	baseRate, err := fx.registry.GetRate(fromCurrency)
	if err != nil {
		return 0, FXConversion{}, fmt.Errorf("fx: no rate available for %q: %w", fromCurrency, err)
	}
	effectiveRate := rate / baseRate
	converted := amount * effectiveRate
	conv := FXConversion{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         effectiveRate,
		AppliedAt:    time.Now().UTC(),
	}
	return converted, conv, nil
}

// ConvertEntry converts all legs of a Transaction from their currency to a target currency,
// returning a new Transaction with updated amounts and a record of each conversion.
func (fx *FXConverter) ConvertEntry(tx Transaction, targetCurrency string) (Transaction, []FXConversion, error) {
	conversions := make([]FXConversion, 0, len(tx.Legs))
	newLegs := make([]Leg, len(tx.Legs))
	for i, leg := range tx.Legs {
		converted, conv, err := fx.Convert(leg.Amount, leg.Currency, targetCurrency)
		if err != nil {
			return Transaction{}, nil, fmt.Errorf("fx: leg %d conversion failed: %w", i, err)
		}
		newLegs[i] = Leg{
			AccountID: leg.AccountID,
			Type:      leg.Type,
			Amount:    converted,
			Currency:  targetCurrency,
		}
		conversions = append(conversions, conv)
	}
	newTx := Transaction{
		ID:          tx.ID,
		Description: tx.Description,
		Timestamp:   tx.Timestamp,
		Legs:        newLegs,
	}
	return newTx, conversions, nil
}
