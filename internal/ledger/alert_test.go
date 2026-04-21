package ledger

import (
	"strings"
	"testing"
)

func makeAlertBook(t *testing.T) *Book {
	t.Helper()
	b := makeBook(t)
	return b
}

func TestAlertEngine_NoViolation(t *testing.T) {
	b := makeAlertBook(t)
	e := NewAlertEngine()
	_ = e.AddRule(AlertRule{
		AccountID: "cash",
		Currency:  "USD",
		Threshold: -1_000_00,
		Below:     true,
		Severity:  SeverityCritical,
		Message:   "cash overdrawn",
	})
	alerts, err := e.Evaluate(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestAlertEngine_Triggered(t *testing.T) {
	b := makeAlertBook(t)
	e := NewAlertEngine()
	// cash balance is 0 after makeBook; alert if balance < 1 (i.e. always triggers here)
	_ = e.AddRule(AlertRule{
		AccountID: "cash",
		Currency:  "USD",
		Threshold: 1,
		Below:     true,
		Severity:  SeverityWarning,
		Message:   "low cash",
	})
	alerts, err := e.Evaluate(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Rule.Severity != SeverityWarning {
		t.Errorf("expected WARNING severity, got %s", alerts[0].Rule.Severity)
	}
}

func TestAlertEngine_NilBook(t *testing.T) {
	e := NewAlertEngine()
	_, err := e.Evaluate(nil)
	if err == nil {
		t.Fatal("expected error for nil book")
	}
}

func TestAlertEngine_AddRule_Invalid(t *testing.T) {
	e := NewAlertEngine()
	if err := e.AddRule(AlertRule{AccountID: "", Currency: "USD", Severity: SeverityInfo}); err == nil {
		t.Error("expected error for empty accountID")
	}
	if err := e.AddRule(AlertRule{AccountID: "cash", Currency: "", Severity: SeverityInfo}); err == nil {
		t.Error("expected error for empty currency")
	}
	if err := e.AddRule(AlertRule{AccountID: "cash", Currency: "USD", Severity: ""}); err == nil {
		t.Error("expected error for empty severity")
	}
}

func TestAlertEngine_UnknownAccount(t *testing.T) {
	b := makeAlertBook(t)
	e := NewAlertEngine()
	_ = e.AddRule(AlertRule{
		AccountID: "nonexistent",
		Currency:  "USD",
		Threshold: 0,
		Below:     true,
		Severity:  SeverityInfo,
		Message:   "ghost account",
	})
	alerts, err := e.Evaluate(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts for unknown account, got %d", len(alerts))
	}
}

func TestExportAlertsCSV_Empty(t *testing.T) {
	out := ExportAlertsCSV(nil)
	if !strings.HasPrefix(out, "triggered_at") {
		t.Errorf("expected CSV header, got: %s", out)
	}
}

func TestExportAlertsCSV_WithAlerts(t *testing.T) {
	b := makeAlertBook(t)
	e := NewAlertEngine()
	_ = e.AddRule(AlertRule{
		AccountID: "cash",
		Currency:  "USD",
		Threshold: 1,
		Below:     true,
		Severity:  SeverityCritical,
		Message:   "test alert",
	})
	alerts, _ := e.Evaluate(b)
	out := ExportAlertsCSV(alerts)
	if !strings.Contains(out, "CRITICAL") {
		t.Errorf("expected CRITICAL in CSV output, got: %s", out)
	}
	if !strings.Contains(out, "cash") {
		t.Errorf("expected account id in CSV output, got: %s", out)
	}
}
