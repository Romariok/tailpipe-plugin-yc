package billing_log_api

import (
	"testing"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
)

func newTestRange(granularity time.Duration) (collection_state.DirectionalTimeRange, time.Time, time.Time) {
	now := time.Now().UTC().Truncate(granularity)
	from := now.Add(-24 * time.Hour)
	to := now
	return collection_state.DirectionalTimeRange{
		LowerBoundary:   from,
		UpperBoundary:   to,
		CollectionOrder: collection_state.CollectionOrderChronological,
	}, from, to
}

func TestBillingLogCollectionState_Init(t *testing.T) {
	gran := time.Hour
	tr, from, _ := newTestRange(gran)

	cs := NewBillingLogCollectionState().(*BillingLogCollectionState)
	cs.Init(tr, gran)

	if !cs.IsEmpty() {
		t.Fatalf("expected state to be empty after Init")
	}
	if got := cs.GetFromTime(); !got.Equal(from) {
		t.Fatalf("GetFromTime mismatch: got %v, want %v", got, from)
	}
	// after Init, the upper boundary equals the lower boundary (nothing collected yet)
	if got := cs.GetToTime(); !got.Equal(from) {
		t.Fatalf("GetToTime mismatch after Init: got %v, want %v", got, from)
	}
}

func TestBillingLogCollectionState_ShouldCollect_OnCollected(t *testing.T) {
	gran := time.Hour
	tr, from, _ := newTestRange(gran)

	cs := NewBillingLogCollectionState().(*BillingLogCollectionState)
	cs.Init(tr, gran)

	// event within the window, mid-hour
	ts1 := from.Add(30 * time.Minute)
	id1 := "acct|cloud|folder|svc|sku|res|2024-01-01|2024-01-01T10:30:00Z"

	if ok := cs.ShouldCollect(id1, ts1); !ok {
		t.Fatalf("ShouldCollect(ts1) = false, want true")
	}
	if err := cs.OnCollected(id1, ts1); err != nil {
		t.Fatalf("OnCollected(ts1) error: %v", err)
	}
	// upper boundary should remain at the hour start (truncate)
	if got := cs.GetToTime(); !got.Equal(from) {
		t.Fatalf("GetToTime after first OnCollected mismatch: got %v, want %v", got, from)
	}
	// duplicate of the same object within the same hour slot — do not collect
	if ok := cs.ShouldCollect(id1, ts1); ok {
		t.Fatalf("ShouldCollect(duplicate ts1) = true, want false")
	}

	// event in the next hour
	ts2 := from.Add(90 * time.Minute) // +1.5 hours
	id2 := "acct|cloud|folder|svc|sku|res|2024-01-01|2024-01-01T11:30:00Z"
	if ok := cs.ShouldCollect(id2, ts2); !ok {
		t.Fatalf("ShouldCollect(ts2) = false, want true")
	}
	if err := cs.OnCollected(id2, ts2); err != nil {
		t.Fatalf("OnCollected(ts2) error: %v", err)
	}
	// upper boundary should move to the next hour
	wantTo := from.Add(1 * time.Hour)
	if got := cs.GetToTime(); !got.Equal(wantTo) {
		t.Fatalf("GetToTime after second OnCollected mismatch: got %v, want %v", got, wantTo)
	}
}

func TestBillingLogCollectionState_ShouldCollect_OutOfRange(t *testing.T) {
	gran := time.Hour
	tr, from, to := newTestRange(gran)

	cs := NewBillingLogCollectionState().(*BillingLogCollectionState)
	cs.Init(tr, gran)

	before := from.Add(-1 * time.Minute)
	if ok := cs.ShouldCollect("x", before); ok {
		t.Fatalf("ShouldCollect(before from) = true, want false")
	}
	atTo := to // upper boundary is exclusive
	if ok := cs.ShouldCollect("y", atTo); ok {
		t.Fatalf("ShouldCollect(at to) = true, want false")
	}
}

func TestBillingLogCollectionState_OnCollectionComplete(t *testing.T) {
	gran := time.Hour
	tr, from, to := newTestRange(gran)

	cs := NewBillingLogCollectionState().(*BillingLogCollectionState)
	cs.Init(tr, gran)

	ts1 := from.Add(30 * time.Minute)
	if ok := cs.ShouldCollect("a", ts1); !ok {
		t.Fatalf("ShouldCollect(ts1) = false, want true")
	}
	if err := cs.OnCollected("a", ts1); err != nil {
		t.Fatalf("OnCollected(ts1) error: %v", err)
	}
	ts2 := from.Add(75 * time.Minute)
	if ok := cs.ShouldCollect("b", ts2); !ok {
		t.Fatalf("ShouldCollect(ts2) = false, want true")
	}
	if err := cs.OnCollected("b", ts2); err != nil {
		t.Fatalf("OnCollected(ts2) error: %v", err)
	}

	// completing collection should advance upper boundary to the end of the window (truncated to granularity)
	if err := cs.OnCollectionComplete(); err != nil {
		t.Fatalf("OnCollectionComplete error: %v", err)
	}
	if got := cs.GetToTime(); !got.Equal(to) {
		t.Fatalf("GetToTime after OnCollectionComplete mismatch: got %v, want %v", got, to)
	}
}

func TestBillingLogCollectionState_Validate(t *testing.T) {
	gran := time.Hour
	tr, _, _ := newTestRange(gran)

	cs := NewBillingLogCollectionState().(*BillingLogCollectionState)
	cs.Init(tr, gran)

	if err := cs.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}
