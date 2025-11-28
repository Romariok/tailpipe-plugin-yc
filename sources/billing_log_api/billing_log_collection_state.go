package billing_log_api

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
)

type BillingLogCollectionState struct {
	State                       *collection_state.TimeRangeCollectionState `json:"state"`
	currentDirectionalTimeRange *collection_state.DirectionalTimeRange
	Granularity                 time.Duration `json:"granularity,omitempty"`
}

func NewBillingLogCollectionState() collection_state.CollectionState {
	return &BillingLogCollectionState{
		State: collection_state.NewTimeRangeCollectionState().(*collection_state.TimeRangeCollectionState),
	}
}

func (s *BillingLogCollectionState) Init(timeRange collection_state.DirectionalTimeRange, granularity time.Duration) {
	s.Granularity = granularity
	if s.State == nil {
		s.State = collection_state.NewTimeRangeCollectionState().(*collection_state.TimeRangeCollectionState)
	}
	s.State.Init(timeRange, granularity)
	s.currentDirectionalTimeRange = &timeRange
}

func (s *BillingLogCollectionState) IsEmpty() bool {
	if s.State == nil {
		return true
	}
	return s.State.IsEmpty()
}

func (s *BillingLogCollectionState) ShouldCollect(id string, timestamp time.Time) bool {
	if s.State == nil {
		return true
	}
	return s.State.ShouldCollect(id, timestamp)
}

func (s *BillingLogCollectionState) OnCollected(id string, timestamp time.Time) error {
	if s.State == nil {
		return nil
	}
	return s.State.OnCollected(id, timestamp)
}

func (s *BillingLogCollectionState) GetFromTime() time.Time {
	if s.State == nil {
		return time.Time{}
	}
	return s.State.GetFromTime()
}

func (s *BillingLogCollectionState) GetToTime() time.Time {
	if s.State == nil {
		return time.Time{}
	}
	return s.State.GetToTime()
}

func (s *BillingLogCollectionState) OnCollectionComplete() error {
	if s.State == nil {
		return nil
	}
	return s.State.OnCollectionComplete()
}

func (s *BillingLogCollectionState) MigrateFromLegacyState(bytes []byte) error {
	return nil
}

func (s *BillingLogCollectionState) Validate() error {
	if s.State == nil {
		return nil
	}
	return s.State.Validate()
}

func (s *BillingLogCollectionState) Clear(timeRange collection_state.DirectionalTimeRange) {
	if s.State == nil {
		return
	}
	s.State.Clear(timeRange)
}
