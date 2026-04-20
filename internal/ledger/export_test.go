package ledger

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
)

func TestExportTransactions_CSV(t *testing.T) {
	b := makeBook(t)

	var buf bytes.Buffer
	err := ExportTransactions(b, &buf, ExportOptions{Format: FormatCSV})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("invalid csv output: %v", err)
	}

	// header + 2 entries from makeBook's single balanced transaction
	if len(records) < 3 {
		t.Fatalf("expected at least 3 rows (header + 2 entries), got %d", len(records))
	}

	header := records[0]
	expectedCols := []string{"transaction_id", "posted_at", "entry_type", "account_id", "amount", "currency", "description"}
	for i, col := range expectedCols {
		if header[i] != col {
			t.Errorf("header[%d] = %q, want %q", i, header[i], col)
		}
	}
}

func TestExportTransactions_FilterByAccount(t *testing.T) {
	b := makeBook(t)

	var buf bytes.Buffer
	err := ExportTransactions(b, &buf, ExportOptions{Format: FormatCSV, AccountID: "acc-debit"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("invalid csv output: %v", err)
	}

	// header + exactly 1 entry for acc-debit
	if len(records) != 2 {
		t.Fatalf("expected 2 rows (header + 1 entry), got %d", len(records))
	}
	if records[1][3] != "acc-debit" {
		t.Errorf("account_id = %q, want %q", records[1][3], "acc-debit")
	}
}

func TestExportTransactions_NilBook(t *testing.T) {
	err := ExportTransactions(nil, &bytes.Buffer{}, ExportOptions{})
	if err == nil {
		t.Fatal("expected error for nil book")
	}
}

func TestExportTransactions_UnsupportedFormat(t *testing.T) {
	b := makeBook(t)
	err := ExportTransactions(b, &bytes.Buffer{}, ExportOptions{Format: "xml"})
	if err == nil || !strings.Contains(err.Error(), "unsupported format") {
		t.Fatalf("expected unsupported format error, got %v", err)
	}
}

func TestFormatAmount(t *testing.T) {
	cases := []struct {
		minor int64
		want  string
	}{
		{1050, "10.50"},
		{100, "1.00"},
		{5, "0.05"},
		{-250, "-2.50"},
		{0, "0.00"},
	}
	for _, tc := range cases {
		got := formatAmount(tc.minor)
		if got != tc.want {
			t.Errorf("formatAmount(%d) = %q, want %q", tc.minor, got, tc.want)
		}
	}
}
