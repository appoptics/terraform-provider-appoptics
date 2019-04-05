package appoptics

import "sync"

// An Aggregator uses the "summary fields" measurement feature in AppOptics to aggregate multiple values
// into a single measurement, storing a count/sum/min/max/last that can be periodically sent as a single aggregate measurement.
// It can be either updated by passing sequential values to UpdateValue or by passing an Aggregator
// to Update, e.g. s.Update(Aggregator{Sum:100,Count:10,Min:5,Max:15})
type Aggregator struct {
	Count int64
	Sum   float64
	Min   float64
	Max   float64
	Last  float64
}

// UpdateValue sets the most recently observed value for this Aggregator, updating sum/count/min/max/last
// accordingly.
func (g *Aggregator) UpdateValue(val float64) {
	if g.Count == 0 {
		g.Min = val
		g.Max = val
	} else {
		if val < g.Min {
			g.Min = val
		}
		if val > g.Max {
			g.Max = val
		}
	}
	g.Count++
	g.Sum += val
	g.Last = val
}

// Update merges another Aggregator into this Aggregator, merging sum/count/min/max/last accordingly. It can
// be used to facilitate efficient input of many data points into an Aggregator in one call, and it can
// also be used to merge two different Aggregators (for example, workers can each maintain their own
// and periodically merge them).
func (g *Aggregator) Update(other Aggregator) {
	if g.Count == 0 {
		g.Count = other.Count
		g.Sum = other.Sum
		g.Min = other.Min
		g.Max = other.Max
		g.Last = other.Last
	} else {
		g.Count += other.Count
		g.Sum += other.Sum
		if other.Min < g.Min {
			g.Min = other.Min
		}
		if other.Max > g.Max {
			g.Max = other.Max
		}
		g.Last = other.Last
	}
}

// SynchronizedAggregator augments an Aggregator with a mutex to allow concurrent access from multiple
// goroutines.
type SynchronizedAggregator struct {
	Aggregator
	m sync.Mutex
}

// UpdateValue is a concurrent-safe wrapper around Aggregator.UpdateValue
func (g *SynchronizedAggregator) UpdateValue(val float64) {
	g.m.Lock()
	defer g.m.Unlock()
	g.Aggregator.UpdateValue(val)
}

// Update is a concurrent-safe wrapper around Aggregator.Update
func (g *SynchronizedAggregator) Update(other Aggregator) {
	g.m.Lock()
	defer g.m.Unlock()
	g.Aggregator.Update(other)
}

// Reset returns a copy the current Aggregator state and resets it back to its zero state.
func (g *SynchronizedAggregator) Reset() Aggregator {
	g.m.Lock()
	defer g.m.Unlock()
	current := g.Aggregator
	g.Aggregator = Aggregator{}
	return current
}
