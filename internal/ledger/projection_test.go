package ledger

import (
	"strings"
	"testing"
	"time"
)

func makeProjectionBook(t *testing.T) *Book {
	t.Helper()
	b := NewBook()
	acc, err := NewAccount("acc1", "Savings", AccountTypeAsset)
	if err != nil {
		t.Fatalf("NewAccount: %v", err)
	}
	if err := b.AddAccount(acc); err != nil {
		t.Fatalf("AddAccount: %v", err)
	}
	return b
}

func TestRunProjection_Valid(t *testing.T) {
	b := makeProjectionBook(t)
	cfg := ProjectionConfig{
		AccountID:    "acc1",
		Currency:     "USD",
		StepDuration: 24 * time.Hour,
		Steps:        5,
		GrowthRate:   0.01,
	}
	r, err := RunProjection(b, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Entries) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(r.Entries))
	}
	// Each step grows by 1%; verify monotonically increasing from 0
	for i := 1; i < len(r.Entries); i++ {
		prev := r.Entries[i-1].Balance["USD"]
		curr := r.Entries[i].Balance["USD"]
		if curr <= prev {
			t.Errorf("step %d: expected balance to grow, got %.4f <= %.4f", i, curr, prev)
		}
	}
}

func TestRunProjection_NilBook(t *testing.T) {
	_, err := RunProjection(nil, ProjectionConfig{AccountID: "acc1", Currency: "USD", StepDuration: time.Hour, Steps: 1})
	if err == nil {
		t.Fatal("expected error for nil book")
	}
}

func TestRunProjection_InvalidSteps(t *testing.T) {
	b := makeProjectionBook(t)
	cfg := ProjectionConfig{AccountID: "acc1", Currency: "USD", StepDuration: time.Hour, Steps: 0}
	_, err := RunProjection(b, cfg)
	if err == nil {
		t.Fatal("expected error for zero steps")
	}
}

func TestRunProjection_UnknownAccount(t *testing.T) {
	b := makeProjectionBook(t)
	cfg := ProjectionConfig{AccountID: "ghost", Currency: "USD", StepDuration: time.Hour, Steps: 3}
	_, err := RunProjection(b, cfg)
	if err == nil {
		t.Fatal("expected error for unknown account")
	}
}

func TestExportProjectionCSV_Valid(t *testing.T) {
	b := makeProjectionBook(t)
	cfg := ProjectionConfig{
		AccountID:    "acc1",
		Currency:     "USD",
		StepDuration: 24 * time.Hour,
		Steps:        3,
		GrowthRate:   0.05,
	}
	r, _ := RunProjection(b, cfg)
	csv, err := ExportProjectionCSV(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(csv, "account_id") {
		t.Error("expected CSV header")
	}
	if !strings.Contains(csv, "acc1") {
		t.Error("expected account ID in CSV")
	}
}

func TestExportProjectionCSV_Nil(t *testing.T) {
	_, err := ExportProjectionCSV(nil)
	if err == nil {
		t.Fatal("expected error for nil result")
	}
}

func TestSummariseProjection_Valid(t *testing.T) {
	b := makeProjectionBook(t)
	cfg := ProjectionConfig{
		AccountID:    "acc1",
		Currency:     "USD",
		StepDuration: 24 * time.Hour,
		Steps:        4,
		GrowthRate:   0.02,
	}
	r, _ := RunProjection(b, cfg)
	s := SummariseProjection(r, "USD")
	if s == nil {
		t.Fatal("expected non-nil summary")
	}
	if s.Steps != 4 {
		t.Errorf("expected 4 steps, got %d", s.Steps)
	}
	if s.Peak < s.Trough {
		t.Error("peak should be >= trough")
	}
}

func TestSummariseProjection_Nil(t *testing.T) {
	s := SummariseProjection(nil, "USD")
	if s != nil {
		t.Fatal("expected nil summary for nil result")
	}
}
