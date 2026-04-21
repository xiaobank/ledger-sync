package ledger

import (
	"strings"
	"testing"
)

func makeFXRegistry(t *testing.T) *CurrencyRegistry {
	t.Helper()
	reg, err := NewCurrencyRegistry()
	if err != nil {
		t.Fatalf("NewCurrencyRegistry: %v", err)
	}
	if err := reg.Register("USD", "US Dollar", 2); err != nil {
		t.Fatalf("register USD: %v", err)
	}
	if err := reg.Register("EUR", "Euro", 2); err != nil {
		t.Fatalf("register EUR: %v", err)
	}
	if err := reg.SetRate("USD", 1.0); err != nil {
		t.Fatalf("set USD rate: %v", err)
	}
	if err := reg.SetRate("EUR", 0.92); err != nil {
		t.Fatalf("set EUR rate: %v", err)
	}
	return reg
}

func TestFXConverter_NilRegistry(t *testing.T) {
	_, err := NewFXConverter(nil)
	if err == nil {
		t.Fatal("expected error for nil registry")
	}
}

func TestFXConverter_SameCurrency(t *testing.T) {
	reg := makeFXRegistry(t)
	fx, _ := NewFXConverter(reg)
	result, conv, err := fx.Convert(100.0, "USD", "USD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 100.0 {
		t.Errorf("expected 100.0, got %.4f", result)
	}
	if conv.Rate != 1.0 {
		t.Errorf("expected rate 1.0, got %.4f", conv.Rate)
	}
}

func TestFXConverter_CrossCurrency(t *testing.T) {
	reg := makeFXRegistry(t)
	fx, _ := NewFXConverter(reg)
	result, conv, err := fx.Convert(100.0, "USD", "EUR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := 100.0 * (0.92 / 1.0)
	if result < expected-0.001 || result > expected+0.001 {
		t.Errorf("expected ~%.4f, got %.4f", expected, result)
	}
	if conv.FromCurrency != "USD" || conv.ToCurrency != "EUR" {
		t.Errorf("unexpected currencies in conversion record")
	}
}

func TestFXConverter_UnknownCurrency(t *testing.T) {
	reg := makeFXRegistry(t)
	fx, _ := NewFXConverter(reg)
	_, _, err := fx.Convert(50.0, "USD", "GBP")
	if err == nil {
		t.Fatal("expected error for unknown currency GBP")
	}
}

func TestFXConverter_NegativeAmount(t *testing.T) {
	reg := makeFXRegistry(t)
	fx, _ := NewFXConverter(reg)
	_, _, err := fx.Convert(-10.0, "USD", "EUR")
	if err == nil {
		t.Fatal("expected error for negative amount")
	}
}

func TestExportFXConversionsCSV_Empty(t *testing.T) {
	csv := ExportFXConversionsCSV(nil)
	if !strings.HasPrefix(csv, "from,to,rate,applied_at") {
		t.Errorf("unexpected header: %q", csv)
	}
	lines := strings.Split(strings.TrimSpace(csv), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line (header only), got %d", len(lines))
	}
}

func TestExportFXConversionsCSV_WithData(t *testing.T) {
	reg := makeFXRegistry(t)
	fx, _ := NewFXConverter(reg)
	_, conv, err := fx.Convert(200.0, "USD", "EUR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	csv := ExportFXConversionsCSV([]FXConversion{conv})
	if !strings.Contains(csv, "USD") || !strings.Contains(csv, "EUR") {
		t.Errorf("CSV missing currency info: %q", csv)
	}
	lines := strings.Split(strings.TrimSpace(csv), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}
