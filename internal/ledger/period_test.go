package ledger

import (
	"strings"
	"testing"
	"time"
)

func TestGeneratePeriods_Monthly(t *testing.T) {
	from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	periods, err := GeneratePeriods(from, PeriodMonthly, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(periods) != 3 {
		t.Fatalf("expected 3 periods, got %d", len(periods))
	}
	if periods[0].Label != "2024-01" {
		t.Errorf("expected 2024-01, got %s", periods[0].Label)
	}
	if periods[2].Label != "2024-03" {
		t.Errorf("expected 2024-03, got %s", periods[2].Label)
	}
}

func TestGeneratePeriods_Daily(t *testing.T) {
	from := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	periods, err := GeneratePeriods(from, PeriodDaily, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if periods[0].Label != "2024-06-01" {
		t.Errorf("unexpected label: %s", periods[0].Label)
	}
	if periods[1].Label != "2024-06-02" {
		t.Errorf("unexpected label: %s", periods[1].Label)
	}
}

func TestGeneratePeriods_InvalidN(t *testing.T) {
	_, err := GeneratePeriods(time.Now(), PeriodMonthly, 0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestPeriod_Contains(t *testing.T) {
	p := Period{
		Start: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
	}
	inside := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	outside := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	if !p.Contains(inside) {
		t.Error("expected inside to be contained")
	}
	if p.Contains(outside) {
		t.Error("expected outside not to be contained")
	}
}

func TestSummariseByPeriod_Valid(t *testing.T) {
	book := makeBook()
	periods, _ := GeneratePeriods(time.Now().AddDate(0, -1, 0), PeriodMonthly, 2)
	summaries, err := SummariseByPeriod(book, periods, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(summaries) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(summaries))
	}
}

func TestSummariseByPeriod_NilBook(t *testing.T) {
	periods, _ := GeneratePeriods(time.Now(), PeriodMonthly, 1)
	_, err := SummariseByPeriod(nil, periods, "")
	if err == nil {
		t.Fatal("expected error for nil book")
	}
}

func TestExportPeriodSummaryCSV_Valid(t *testing.T) {
	periods, _ := GeneratePeriods(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), PeriodMonthly, 1)
	summaries := []PeriodSummary{
		{Period: periods[0], Count: 2, Total: map[string]float64{"USD": 150.00}},
	}
	csv := ExportPeriodSummaryCSV(summaries)
	if !strings.Contains(csv, "period,currency,count,total") {
		t.Error("missing CSV header")
	}
	if !strings.Contains(csv, "2024-01") {
		t.Error("missing period label in CSV")
	}
	if !strings.Contains(csv, "150.00") {
		t.Error("missing total in CSV")
	}
}

func TestExportPeriodSummaryCSV_Empty(t *testing.T) {
	periods, _ := GeneratePeriods(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), PeriodMonthly, 1)
	summaries := []PeriodSummary{{Period: periods[0], Count: 0, Total: map[string]float64{}}}
	csv := ExportPeriodSummaryCSV(summaries)
	if !strings.Contains(csv, "0.00") {
		t.Error("expected zero total in empty period CSV")
	}
}
