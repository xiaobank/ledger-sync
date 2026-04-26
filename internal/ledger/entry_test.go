package ledger

import (
	"testing"
)

func TestNewTransaction_Balanced(t *testing.T) {
	entries := []Entry{
		{Account: "cash", Type: Debit, AmountCents: 1000, Currency: "USD"},
		{Account: "revenue", Type: Credit, AmountCents: 1000, Currency: "USD"},
	}
	tx, err := NewTransaction("ref-001", entries)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if tx.ID == "" {
		t.Error("expected transaction ID to be set")
	}
	if tx.Reference != "ref-001" {
		t.Errorf("expected reference ref-001, got %s", tx.Reference)
	}
	for _, e := range tx.Entries {
		if e.TransactionID != tx.ID {
			t.Errorf("entry transaction_id mismatch: %s", e.TransactionID)
		}
		if e.ID == "" {
			t.Error("expected entry ID to be set")
		}
	}
}

func TestNewTransaction_Unbalanced(t *testing.T) {
	entries := []Entry{
		{Account: "cash", Type: Debit, AmountCents: 1000, Currency: "USD"},
		{Account: "revenue", Type: Credit, AmountCents: 500, Currency: "USD"},
	}
	_, err := NewTransaction("ref-002", entries)
	if err == nil {
		t.Fatal("expected error for unbalanced transaction")
	}
}

func TestNewTransaction_InvalidAmount(t *testing.T) {
	entries := []Entry{
		{Account: "cash", Type: Debit, AmountCents: -100, Currency: "USD"},
		{Account: "revenue", Type: Credit, AmountCents: -100, Currency: "USD"},
	}
	_, err := NewTransaction("ref-003", entries)
	if err == nil {
		t.Fatal("expected error for non-positive amount")
	}
}

func TestNewTransaction_MultiCurrencyBalanced(t *testing.T) {
	entries := []Entry{
		{Account: "cash_usd", Type: Debit, AmountCents: 500, Currency: "USD"},
		{Account: "revenue_usd", Type: Credit, AmountCents: 500, Currency: "USD"},
		{Account: "cash_eur", Type: Debit, AmountCents: 300, Currency: "EUR"},
		{Account: "revenue_eur", Type: Credit, AmountCents: 300, Currency: "EUR"},
	}
	_, err := NewTransaction("ref-004", entries)
	if err != nil {
		t.Fatalf("expected no error for multi-currency balanced tx, got: %v", err)
	}
}

func TestNewTransaction_InvalidType(t *testing.T) {
	entries := []Entry{
		{Account: "cash", Type: "invalid", AmountCents: 100, Currency: "USD"},
	}
	_, err := NewTransaction("ref-005", entries)
	if err == nil {
		t.Fatal("expected error for invalid entry type")
	}
}

func TestNewTransaction_EmptyEntries(t *testing.T) {
	_, err := NewTransaction("ref-006", []Entry{})
	if err == nil {
		t.Fatal("expected error for empty entries slice")
	}
}

func TestNewTransaction_EmptyReference(t *testing.T) {
	entries := []Entry{
		{Account: "cash", Type: Debit, AmountCents: 1000, Currency: "USD"},
		{Account: "revenue", Type: Credit, AmountCents: 1000, Currency: "USD"},
	}
	_, err := NewTransaction("", entries)
	if err == nil {
		t.Fatal("expected error for empty reference")
	}
}
