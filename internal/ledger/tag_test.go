package ledger

import (
	"testing"
)

func TestTagIndex_Add_And_Lookup(t *testing.T) {
	ti := NewTagIndex()
	tag := Tag{Key: "region", Value: "us-east"}

	if err := ti.Add("tx-001", tag); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ti.Add("tx-002", tag); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results := ti.Lookup(tag)
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestTagIndex_Lookup_NotFound(t *testing.T) {
	ti := NewTagIndex()
	results := ti.Lookup(Tag{Key: "env", Value: "prod"})
	if len(results) != 0 {
		t.Errorf("expected empty slice, got %v", results)
	}
}

func TestTagIndex_Add_EmptyKey(t *testing.T) {
	ti := NewTagIndex()
	err := ti.Add("tx-001", Tag{Key: "", Value: "val"})
	if err == nil {
		t.Error("expected error for empty tag key, got nil")
	}
}

func TestTagIndex_Add_EmptyTxID(t *testing.T) {
	ti := NewTagIndex()
	err := ti.Add("", Tag{Key: "env", Value: "prod"})
	if err == nil {
		t.Error("expected error for empty transaction ID, got nil")
	}
}

func TestTagIndex_Keys(t *testing.T) {
	ti := NewTagIndex()
	_ = ti.Add("tx-001", Tag{Key: "env", Value: "prod"})
	_ = ti.Add("tx-002", Tag{Key: "region", Value: "eu-west"})

	keys := ti.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "env:prod" {
		t.Errorf("unexpected first key: %s", keys[0])
	}
	if keys[1] != "region:eu-west" {
		t.Errorf("unexpected second key: %s", keys[1])
	}
}

func TestTagIndex_Remove(t *testing.T) {
	ti := NewTagIndex()
	tag := Tag{Key: "env", Value: "staging"}
	_ = ti.Add("tx-001", tag)

	ti.Remove(tag)
	results := ti.Lookup(tag)
	if len(results) != 0 {
		t.Errorf("expected empty results after remove, got %v", results)
	}
}

func TestTagIndex_Lookup_ReturnsSortedIDs(t *testing.T) {
	ti := NewTagIndex()
	tag := Tag{Key: "team", Value: "payments"}
	_ = ti.Add("tx-zzz", tag)
	_ = ti.Add("tx-aaa", tag)
	_ = ti.Add("tx-mmm", tag)

	results := ti.Lookup(tag)
	if results[0] != "tx-aaa" || results[1] != "tx-mmm" || results[2] != "tx-zzz" {
		t.Errorf("results not sorted: %v", results)
	}
}
