package ledger

import (
	"testing"
)

func TestCurrencyRegistry_Register_Valid(t *testing.T) {
	reg := NewCurrencyRegistry()
	if err := reg.Register("USD"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reg.IsKnown("USD") {
		t.Error("expected USD to be known after registration")
	}
}

func TestCurrencyRegistry_Register_Invalid(t *testing.T) {
	reg := NewCurrencyRegistry()
	if err := reg.Register("US"); err == nil {
		t.Error("expected error for 2-char currency code")
	}
	if err := reg.Register("USDD"); err == nil {
		t.Error("expected error for 4-char currency code")
	}
}

func TestCurrencyRegistry_IsKnown_CaseInsensitive(t *testing.T) {
	reg := NewCurrencyRegistry()
	_ = reg.Register("eur")
	if !reg.IsKnown("EUR") {
		t.Error("expected EUR to be known regardless of registration case")
	}
}

func TestCurrencyRegistry_SetRate_Valid(t *testing.T) {
	reg := NewCurrencyRegistry()
	_ = reg.Register("USD")
	_ = reg.Register("EUR")
	if err := reg.SetRate("USD", "EUR", 0.92); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCurrencyRegistry_SetRate_InvalidRate(t *testing.T) {
	reg := NewCurrencyRegistry()
	_ = reg.Register("USD")
	_ = reg.Register("EUR")
	if err := reg.SetRate("USD", "EUR", 0); err == nil {
		t.Error("expected error for zero rate")
	}
	if err := reg.SetRate("USD", "EUR", -1.5); err == nil {
		t.Error("expected error for negative rate")
	}
}

func TestCurrencyRegistry_SetRate_UnknownCurrency(t *testing.T) {
	reg := NewCurrencyRegistry()
	_ = reg.Register("USD")
	if err := reg.SetRate("USD", "GBP", 0.78); err == nil {
		t.Error("expected error for unknown target currency")
	}
}

func TestCurrencyRegistry_Convert_SameCurrency(t *testing.T) {
	reg := NewCurrencyRegistry()
	_ = reg.Register("USD")
	result, err := reg.Convert(1000, "USD", "USD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 1000 {
		t.Errorf("expected 1000, got %d", result)
	}
}

func TestCurrencyRegistry_Convert_WithRate(t *testing.T) {
	reg := NewCurrencyRegistry()
	_ = reg.Register("USD")
	_ = reg.Register("EUR")
	_ = reg.SetRate("USD", "EUR", 0.5)
	result, err := reg.Convert(200, "USD", "EUR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 100 {
		t.Errorf("expected 100, got %d", result)
	}
}

func TestCurrencyRegistry_Convert_MissingRate(t *testing.T) {
	reg := NewCurrencyRegistry()
	_ = reg.Register("USD")
	_ = reg.Register("JPY")
	_, err := reg.Convert(500, "USD", "JPY")
	if err == nil {
		t.Error("expected error when no rate is set")
	}
}
