package ledger

import (
	"strings"
	"testing"
)

func TestAuditLog_Record_And_Events(t *testing.T) {
	log := NewAuditLog()
	log.Record(AuditEventPost, "svc-payments", "txn-001")
	log.Record(AuditEventAddAccount, "svc-accounts", "acc-cash")

	events := log.Events()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].EventType != AuditEventPost {
		t.Errorf("expected POST, got %s", events[0].EventType)
	}
	if events[1].ActorID != "svc-accounts" {
		t.Errorf("expected actor svc-accounts, got %s", events[1].ActorID)
	}
}

func TestAuditLog_Events_ReturnsCopy(t *testing.T) {
	log := NewAuditLog()
	log.Record(AuditEventSnapshot, "svc-audit", "snap-1")

	events := log.Events()
	events[0].ActorID = "tampered"

	original := log.Events()
	if original[0].ActorID == "tampered" {
		t.Error("Events() should return a copy, not a reference")
	}
}

func TestAuditLog_FilterByType(t *testing.T) {
	log := NewAuditLog()
	log.Record(AuditEventPost, "svc-a", "txn-1")
	log.Record(AuditEventReconcile, "svc-b", "recon-1")
	log.Record(AuditEventPost, "svc-a", "txn-2")

	posts := log.FilterByType(AuditEventPost)
	if len(posts) != 2 {
		t.Errorf("expected 2 POST events, got %d", len(posts))
	}

	reconciles := log.FilterByType(AuditEventReconcile)
	if len(reconciles) != 1 {
		t.Errorf("expected 1 RECONCILE event, got %d", len(reconciles))
	}
}

func TestAuditLog_FilterByActor(t *testing.T) {
	log := NewAuditLog()
	log.Record(AuditEventPost, "svc-payments", "txn-10")
	log.Record(AuditEventPost, "svc-billing", "txn-11")
	log.Record(AuditEventAddAccount, "svc-payments", "acc-revenue")

	paymentEvents := log.FilterByActor("svc-payments")
	if len(paymentEvents) != 2 {
		t.Errorf("expected 2 events for svc-payments, got %d", len(paymentEvents))
	}
}

func TestAuditLog_String_Empty(t *testing.T) {
	log := NewAuditLog()
	out := log.String()
	if !strings.Contains(out, "no events") {
		t.Errorf("expected 'no events' in empty log string, got: %s", out)
	}
}

func TestAuditLog_String_WithEvents(t *testing.T) {
	log := NewAuditLog()
	log.Record(AuditEventPost, "svc-x", "txn-abc")
	out := log.String()
	if !strings.Contains(out, "POST") {
		t.Errorf("expected POST in log string, got: %s", out)
	}
	if !strings.Contains(out, "svc-x") {
		t.Errorf("expected actor svc-x in log string, got: %s", out)
	}
}
