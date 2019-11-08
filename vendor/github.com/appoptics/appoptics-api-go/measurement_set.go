package appoptics

import (
	"sync"

	"golang.org/x/net/context"
)

type ctxMarker struct{}

// DefaultSink is a convenience instance of MeasurementSet that can be used to centrally aggregate
// measurements for an entire process.
var (
	DefaultSink  = NewMeasurementSet()
	ctxMarkerKey = &ctxMarker{}
)

// MeasurementSet represents a map of SynchronizedCounters and SynchronizedAggregators. All functions
// of MeasurementSet are safe for concurrent use.
type MeasurementSet struct {
	counters         map[string]*SynchronizedCounter
	aggregators      map[string]*SynchronizedAggregator
	countersMutex    sync.RWMutex
	aggregatorsMutex sync.RWMutex
}

// NewMeasurementSet returns a new empty MeasurementSet
func NewMeasurementSet() *MeasurementSet {
	return &MeasurementSet{
		counters:    map[string]*SynchronizedCounter{},
		aggregators: map[string]*SynchronizedAggregator{},
	}
}

// GetCounter returns a SynchronizedCounter assigned to the specified key, creating a new one
// if necessary.
func (s *MeasurementSet) GetCounter(key string) *SynchronizedCounter {
	s.countersMutex.RLock()
	counter, ok := s.counters[key]
	s.countersMutex.RUnlock()
	if !ok {
		s.countersMutex.Lock()
		counter, ok = s.counters[key]
		if !ok {
			counter = NewCounter()
			s.counters[key] = counter
		}
		s.countersMutex.Unlock()
	}
	return counter
}

// GetAggregator returns a SynchronizedAggregator assigned to the specified key, creating a new one
// if necessary.
func (s *MeasurementSet) GetAggregator(key string) *SynchronizedAggregator {
	s.aggregatorsMutex.RLock()
	agg, ok := s.aggregators[key]
	s.aggregatorsMutex.RUnlock()
	if !ok {
		s.aggregatorsMutex.Lock()
		agg, ok = s.aggregators[key]
		if !ok {
			agg = &SynchronizedAggregator{}
			s.aggregators[key] = agg
		}
		s.aggregatorsMutex.Unlock()
	}
	return agg
}

// Incr is a convenience function to get the specified Counter and call Incr on it. See Counter.Incr.
func (s *MeasurementSet) Incr(key string) {
	s.GetCounter(key).Incr()
}

// Add is a convenience function to get the specified Counter and call Add on it. See Counter.Add.
func (s *MeasurementSet) Add(key string, delta int64) {
	s.GetCounter(key).Add(delta)
}

// AddInt is a convenience function to get the specified Counter and call AddInt on it. See
// Counter.AddInt.
func (s *MeasurementSet) AddInt(key string, delta int) {
	s.GetCounter(key).AddInt(delta)
}

// UpdateAggregatorValue is a convenience to get the specified Aggregator and call UpdateValue on it.
// See Aggregator.UpdateValue.
func (s *MeasurementSet) UpdateAggregatorValue(key string, val float64) {
	s.GetAggregator(key).UpdateValue(val)
}

// UpdateAggregator is a convenience to get the specified Aggregator and call Update on it. See Aggregator.Update.
func (s *MeasurementSet) UpdateAggregator(key string, other Aggregator) {
	s.GetAggregator(key).Update(other)
}

// Merge takes a MeasurementSetReport and merges all of it Counters and Aggregators into this MeasurementSet.
// This in turn calls Counter.Add for each Counter in the report, and Aggregator.Update for each Aggregator in
// the report. Any keys that do not exist in this MeasurementSet will be created.
func (s *MeasurementSet) Merge(report *MeasurementSetReport) {
	for key, value := range report.Counts {
		s.GetCounter(key).Add(value)
	}
	for key, agg := range report.Aggregators {
		s.GetAggregator(key).Update(agg)
	}
}

// Reset generates a MeasurementSetReport with a copy of the state of each of the non-zero Counters and
// Aggregators in this MeasurementSet. Counters with a value of 0 and Aggregators with a count of 0 are omitted.
// All Counters and Aggregators are reset to the zero/nil state but are never removed from this
// MeasurementSet, so they can continue be used indefinitely.
func (s *MeasurementSet) Reset() *MeasurementSetReport {
	report := NewMeasurementSetReport()
	s.countersMutex.Lock()
	for key, counter := range s.counters {
		val := counter.Reset()
		if val != 0 {
			report.Counts[key] = val
		}
	}
	s.countersMutex.Unlock()
	s.aggregatorsMutex.Lock()
	for key, syncAggregator := range s.aggregators {
		agg := syncAggregator.Reset()
		if agg.Count != 0 {
			report.Aggregators[key] = agg
		}
	}
	s.aggregatorsMutex.Unlock()
	return report
}

// ContextWithMeasurementSet wraps the specified context with a MeasurementSet.
// XXX TODO: add convenience methods to read that MeasurementSet and manipulate Counters/Aggregators on it.
func ContextWithMeasurementSet(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, NewMeasurementSet())
}
