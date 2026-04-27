package ledger

import (
	"testing"
	"time"
)

func makeAgingBook(t *testing.T) *Book {
	t.Helper()
	b := makeBook(t)

	acc, _ := NewAccount("acc-ar", "Accounts Receivable", AccountTypeAsset)
	acc2, _ := NewAccount("acc-rev", "Revenue", AccountTypeRevenue)
	_ = b.AddAccount(acc)
	_ = b.AddAccount(acc2)
	return b
}

func TestAgingAnalysis_Valid(t *testing.T) {
	b := makeAgingBook(t)
	asOf := time.Now()

	tx, _ := NewTransaction("tx-1", asOf.AddDate(0, 0, -10), []Leg{
		{AccountID: "acc-ar", Amount: 100, Currency: "USD", Type: EntryTypeDebit},
		{AccountID: "acc-rev", Amount: 100, Currency: "USD", Type: EntryTypeCredit},
	})
	_ = b.Post(tx)

	tx2, _ := NewTransaction("tx-2", asOf.AddDate(0, 0, -50), []Leg{
		{AccountID: "acc-ar", Amount: 200, Currency: "USD", Type: EntryTypeDebit},
		{AccountID: "acc-rev", Amount: 200, Currency: "USD", Type: EntryTypeCredit},
	})
	_ = b.Post(tx2)

	report, err := AgingAnalysis(b, "acc-ar", asOf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.AccountID != "acc-ar" {
		t.Errorf("expected acc-ar, got %s", report.AccountID)
	}
	if report.Buckets[0].Total["USD"] != 100 {
		t.Errorf("expected 100 in Current bucket, got %v", report.Buckets[0].Total["USD"])
	}
	if report.Buckets[1].Total["USD"] != 200 {
		t.Errorf("expected 200 in 31-60 bucket, got %v", report.Buckets[1].Total["USD"])
	}
}

func TestAgingAnalysis_NilBook(t *testing.T) {
	_, err := AgingAnalysis(nil, "acc-ar", time.Now())
	if err == nil {
		t.Fatal("expected error for nil book")
	}
}

func TestAgingAnalysis_EmptyAccountID(t *testing.T) {
	b := makeAgingBook(t)
	_, err := AgingAnalysis(b, "", time.Now())
	if err == nil {
		t.Fatal("expected error for empty accountID")
	}
}

func TestAgingAnalysis_FutureTxExcluded(t *testing.T) {
	b := makeAgingBook(t)
	asOf := time.Now()

	tx, _ := NewTransaction("tx-future", asOf.AddDate(0, 0, 5), []Leg{
		{AccountID: "acc-ar", Amount: 500, Currency: "USD", Type: EntryTypeDebit},
		{AccountID: "acc-rev", Amount: 500, Currency: "USD", Type: EntryTypeCredit},
	})
	_ = b.Post(tx)

	report, err := AgingAnalysis(b, "acc-ar", asOf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, bucket := range report.Buckets {
		if bucket.Total["USD"] != 0 {
			t.Errorf("expected 0 in bucket %s, got %v", bucket.Label, bucket.Total["USD"])
		}
	}
}
