package ledger

import (
	"strings"
	"testing"
)

func TestRuleIndex_Add_Valid(t *testing.T) {
	ri := NewRuleIndex()
	err := ri.Add(Rule{ID: "r1", Type: RuleTypeMinBalance, AccountID: "acc1", Currency: "USD", Threshold: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r, ok := ri.Get("r1")
	if !ok || r.AccountID != "acc1" {
		t.Error("rule not stored correctly")
	}
}

func TestRuleIndex_Add_Invalid(t *testing.T) {
	ri := NewRuleIndex()
	cases := []Rule{
		{ID: "", Type: RuleTypeMinBalance, AccountID: "acc1", Currency: "USD"},
		{ID: "r1", Type: RuleTypeMinBalance, AccountID: "", Currency: "USD"},
		{ID: "r2", Type: RuleTypeMinBalance, AccountID: "acc1", Currency: ""},
		{ID: "r3", Type: RuleTypeRequireTag, AccountID: "acc1", TagKey: ""},
		{ID: "r4", Type: "unknown", AccountID: "acc1"},
	}
	for _, c := range cases {
		if err := ri.Add(c); err == nil {
			t.Errorf("expected error for rule %+v", c)
		}
	}
}

func TestRuleIndex_Add_Duplicate(t *testing.T) {
	ri := NewRuleIndex()
	r := Rule{ID: "r1", Type: RuleTypeMaxBalance, AccountID: "acc1", Currency: "USD", Threshold: 500}
	_ = ri.Add(r)
	if err := ri.Add(r); err == nil {
		t.Error("expected duplicate error")
	}
}

func TestCheckRules_NoViolation(t *testing.T) {
	b, ri, _ := makeRuleBook(t)
	_ = ri.Add(Rule{ID: "r1", Type: RuleTypeMinBalance, AccountID: "cash", Currency: "USD", Threshold: 0})
	v, err := CheckRules(b, ri, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Errorf("expected no violations, got %d", len(v))
	}
}

func TestCheckRules_MinBalanceViolation(t *testing.T) {
	b, ri, _ := makeRuleBook(t)
	_ = ri.Add(Rule{ID: "r1", Type: RuleTypeMinBalance, AccountID: "cash", Currency: "USD", Threshold: 999999})
	v, err := CheckRules(b, ri, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 1 {
		t.Errorf("expected 1 violation, got %d", len(v))
	}
}

func TestCheckRules_NilBook(t *testing.T) {
	ri := NewRuleIndex()
	_, err := CheckRules(nil, ri, nil)
	if err == nil {
		t.Error("expected error for nil book")
	}
}

func TestCheckRules_NilIndex(t *testing.T) {
	b, _, _ := makeRuleBook(t)
	_, err := CheckRules(b, nil, nil)
	if err == nil {
		t.Error("expected error for nil index")
	}
}

func TestExportRulesCSV(t *testing.T) {
	_, ri, _ := makeRuleBook(t)
	_ = ri.Add(Rule{ID: "r1", Type: RuleTypeMaxBalance, AccountID: "cash", Currency: "USD", Threshold: 10000})
	csv, err := ExportRulesCSV(ri)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(csv, "r1") {
		t.Error("CSV missing rule ID")
	}
}

// makeRuleBook is a helper that creates a minimal book and rule index.
func makeRuleBook(t *testing.T) (*Book, *RuleIndex, *TagIndex) {
	t.Helper()
	b := NewBook("test-book")
	acc, _ := NewAccount("cash", "Cash", AccountTypeAsset)
	_ = b.AddAccount(acc)
	tx, _ := NewTransaction("tx1", []Entry{
		{AccountID: "cash", Currency: "USD", Amount: 500, Type: EntryTypeDebit},
		{AccountID: "cash", Currency: "USD", Amount: 500, Type: EntryTypeCredit},
	})
	_ = b.Post(tx)
	return b, NewRuleIndex(), NewTagIndex()
}
