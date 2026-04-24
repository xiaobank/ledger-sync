package ledger

import (
	"strings"
	"testing"
)

func TestExportJournalCSV_Valid(t *testing.T) {
	b := makeJournalBook(t)
	j, _ := BuildJournal(b)

	csv, err := ExportJournalCSV(j)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(csv, "date,tx_id") {
		t.Error("expected CSV header")
	}
	if !strings.Contains(csv, "acc-a") {
		t.Error("expected acc-a in CSV output")
	}
	if !strings.Contains(csv, "acc-b") {
		t.Error("expected acc-b in CSV output")
	}
}

func TestExportJournalCSV_NilJournal(t *testing.T) {
	_, err := ExportJournalCSV(nil)
	if err == nil {
		t.Fatal("expected error for nil journal")
	}
}

func TestExportJournalCSV_Empty(t *testing.T) {
	j := &Journal{}
	csv, err := ExportJournalCSV(j)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(csv), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}

func TestSummariseJournal_Valid(t *testing.T) {
	b := makeJournalBook(t)
	j, _ := BuildJournal(b)

	summaries := SummariseJournal(j)
	if len(summaries) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(summaries))
	}
	for _, s := range summaries {
		if s.Date == "" {
			t.Error("expected non-empty date in summary")
		}
		if s.TxID == "" {
			t.Error("expected non-empty tx_id in summary")
		}
	}
}

func TestSummariseJournal_Nil(t *testing.T) {
	result := SummariseJournal(nil)
	if result != nil {
		t.Error("expected nil for nil journal")
	}
}
