package ledger

import (
	"strings"
	"testing"
)

func makeCheckpointBook(t *testing.T) *Book {
	t.Helper()
	b := NewBook("cp-book")
	acc, _ := NewAccount("a1", "Alpha", AccountTypeAsset)
	_ = b.AddAccount(acc)
	acc2, _ := NewAccount("a2", "Beta", AccountTypeLiability)
	_ = b.AddAccount(acc2)
	tx, _ := NewTransaction("tx1", []Entry{
		{AccountID: "a1", Amount: 500, Currency: "USD", Type: EntryTypeDebit},
		{AccountID: "a2", Amount: 500, Currency: "USD", Type: EntryTypeCredit},
	})
	_ = b.Post(tx)
	return b
}

func TestCheckpointIndex_Capture_And_Get(t *testing.T) {
	b := makeCheckpointBook(t)
	ci := NewCheckpointIndex()
	if err := ci.Capture("v1", b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cp, ok := ci.Get("v1")
	if !ok {
		t.Fatal("expected checkpoint v1 to exist")
	}
	if cp.Name != "v1" {
		t.Errorf("expected name v1, got %s", cp.Name)
	}
}

func TestCheckpointIndex_Capture_NilBook(t *testing.T) {
	ci := NewCheckpointIndex()
	if err := ci.Capture("v1", nil); err == nil {
		t.Fatal("expected error for nil book")
	}
}

func TestCheckpointIndex_Capture_EmptyName(t *testing.T) {
	b := makeCheckpointBook(t)
	ci := NewCheckpointIndex()
	if err := ci.Capture("", b); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCheckpointIndex_Get_NotFound(t *testing.T) {
	ci := NewCheckpointIndex()
	_, ok := ci.Get("missing")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestCheckpointIndex_Names(t *testing.T) {
	b := makeCheckpointBook(t)
	ci := NewCheckpointIndex()
	_ = ci.Capture("v1", b)
	_ = ci.Capture("v2", b)
	names := ci.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestCheckpointIndex_DiffCheckpoints(t *testing.T) {
	b := makeCheckpointBook(t)
	ci := NewCheckpointIndex()
	_ = ci.Capture("before", b)

	tx2, _ := NewTransaction("tx2", []Entry{
		{AccountID: "a1", Amount: 200, Currency: "USD", Type: EntryTypeDebit},
		{AccountID: "a2", Amount: 200, Currency: "USD", Type: EntryTypeCredit},
	})
	_ = b.Post(tx2)
	_ = ci.Capture("after", b)

	delta, err := ci.DiffCheckpoints("before", "after")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if delta["a1"]["USD"] != 200 {
		t.Errorf("expected delta 200 for a1/USD, got %.2f", delta["a1"]["USD"])
	}
}

func TestCheckpointIndex_DiffCheckpoints_NotFound(t *testing.T) {
	ci := NewCheckpointIndex()
	_, err := ci.DiffCheckpoints("x", "y")
	if err == nil {
		t.Fatal("expected error for missing checkpoints")
	}
}

func TestExportCheckpointCSV_Valid(t *testing.T) {
	b := makeCheckpointBook(t)
	ci := NewCheckpointIndex()
	_ = ci.Capture("v1", b)
	cp, _ := ci.Get("v1")
	out, err := ExportCheckpointCSV(cp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "checkpoint,created_at,account_id,currency,balance") {
		t.Error("expected CSV header")
	}
	if !strings.Contains(out, "a1") {
		t.Error("expected account a1 in output")
	}
}

func TestExportCheckpointCSV_Nil(t *testing.T) {
	_, err := ExportCheckpointCSV(nil)
	if err == nil {
		t.Fatal("expected error for nil checkpoint")
	}
}
