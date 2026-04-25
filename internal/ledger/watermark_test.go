package ledger

import (
	"testing"
	"time"
)

func TestWatermarkIndex_Observe_And_Get(t *testing.T) {
	wi := NewWatermarkIndex()
	now := time.Now()

	if err := wi.Observe("acc1", "USD", 100.0, now); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	high := wi.Get("acc1", "USD", WatermarkHigh)
	if high == nil {
		t.Fatal("expected high watermark, got nil")
	}
	if high.Amount != 100.0 {
		t.Errorf("expected 100.0, got %f", high.Amount)
	}

	low := wi.Get("acc1", "USD", WatermarkLow)
	if low == nil {
		t.Fatal("expected low watermark, got nil")
	}
	if low.Amount != 100.0 {
		t.Errorf("expected 100.0, got %f", low.Amount)
	}
}

func TestWatermarkIndex_HighUpdates(t *testing.T) {
	wi := NewWatermarkIndex()
	t0 := time.Now()
	t1 := t0.Add(time.Hour)

	_ = wi.Observe("acc1", "USD", 50.0, t0)
	_ = wi.Observe("acc1", "USD", 200.0, t1)

	high := wi.Get("acc1", "USD", WatermarkHigh)
	if high.Amount != 200.0 {
		t.Errorf("expected high 200.0, got %f", high.Amount)
	}
	low := wi.Get("acc1", "USD", WatermarkLow)
	if low.Amount != 50.0 {
		t.Errorf("expected low 50.0, got %f", low.Amount)
	}
}

func TestWatermarkIndex_LowUpdates(t *testing.T) {
	wi := NewWatermarkIndex()
	t0 := time.Now()

	_ = wi.Observe("acc2", "EUR", 300.0, t0)
	_ = wi.Observe("acc2", "EUR", -10.0, t0.Add(time.Minute))

	low := wi.Get("acc2", "EUR", WatermarkLow)
	if low.Amount != -10.0 {
		t.Errorf("expected low -10.0, got %f", low.Amount)
	}
}

func TestWatermarkIndex_Get_NotFound(t *testing.T) {
	wi := NewWatermarkIndex()
	if wm := wi.Get("missing", "USD", WatermarkHigh); wm != nil {
		t.Errorf("expected nil, got %+v", wm)
	}
}

func TestWatermarkIndex_Observe_EmptyAccountID(t *testing.T) {
	wi := NewWatermarkIndex()
	if err := wi.Observe("", "USD", 10.0, time.Now()); err == nil {
		t.Error("expected error for empty accountID")
	}
}

func TestWatermarkIndex_Observe_EmptyCurrency(t *testing.T) {
	wi := NewWatermarkIndex()
	if err := wi.Observe("acc1", "", 10.0, time.Now()); err == nil {
		t.Error("expected error for empty currency")
	}
}

func TestWatermarkIndex_All(t *testing.T) {
	wi := NewWatermarkIndex()
	_ = wi.Observe("acc1", "USD", 100.0, time.Now())
	_ = wi.Observe("acc2", "GBP", 50.0, time.Now())

	all := wi.All()
	// 2 accounts × 2 kinds = 4 watermarks
	if len(all) != 4 {
		t.Errorf("expected 4 watermarks, got %d", len(all))
	}
}
