package ledger

import (
	"strings"
	"testing"
	"time"
)

func makeForecastBook(t *testing.T) *Book {
	t.Helper()
	b := makePopulatedBook(t)
	return b
}

func TestForecast_Valid(t *testing.T) {
	b := makeForecastBook(t)
	opts := ForecastOptions{
		Periods:        3,
		PeriodDuration: 30 * 24 * time.Hour,
		GrowthRate:     0.05,
	}

	r, err := Forecast(b, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil result")
	}
	if len(r.Entries) == 0 {
		t.Fatal("expected at least one forecast entry")
	}
}

func TestForecast_NilBook(t *testing.T) {
	_, err := Forecast(nil, ForecastOptions{Periods: 1, PeriodDuration: time.Hour})
	if err == nil {
		t.Fatal("expected error for nil book")
	}
}

func TestForecast_InvalidPeriods(t *testing.T) {
	b := makeForecastBook(t)
	_, err := Forecast(b, ForecastOptions{Periods: 0, PeriodDuration: time.Hour})
	if err == nil {
		t.Fatal("expected error for zero periods")
	}
}

func TestForecast_InvalidDuration(t *testing.T) {
	b := makeForecastBook(t)
	_, err := Forecast(b, ForecastOptions{Periods: 1, PeriodDuration: 0})
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestExportForecastCSV_Valid(t *testing.T) {
	b := makeForecastBook(t)
	r, err := Forecast(b, ForecastOptions{
		Periods:        2,
		PeriodDuration: 30 * 24 * time.Hour,
		GrowthRate:     0.0,
	})
	if err != nil {
		t.Fatalf("forecast error: %v", err)
	}

	csv, err := ExportForecastCSV(r)
	if err != nil {
		t.Fatalf("export error: %v", err)
	}
	if !strings.HasPrefix(csv, "account_id,currency,projected,period_date") {
		t.Errorf("missing CSV header, got: %s", csv[:40])
	}
}

func TestExportForecastCSV_NilResult(t *testing.T) {
	_, err := ExportForecastCSV(nil)
	if err == nil {
		t.Fatal("expected error for nil result")
	}
}

func TestForecastSummary_Empty(t *testing.T) {
	out := ForecastSummary(nil)
	if !strings.Contains(out, "No forecast data") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestForecastSummary_Valid(t *testing.T) {
	b := makeForecastBook(t)
	r, _ := Forecast(b, ForecastOptions{
		Periods:        1,
		PeriodDuration: 30 * 24 * time.Hour,
		GrowthRate:     0.1,
	})
	out := ForecastSummary(r)
	if !strings.Contains(out, "Forecast generated at:") {
		t.Errorf("expected header line in summary, got: %s", out)
	}
}
