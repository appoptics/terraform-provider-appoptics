package appoptics

type TaggedMeasurementSet struct {
	*MeasurementSet
	tags map[string]interface{}
}

// Tags returns the tags map
func (s *TaggedMeasurementSet) Tags() map[string]interface{} {
	return s.tags
}

// SetTags sets the value of the tags map
func (s *TaggedMeasurementSet) SetTags(tags map[string]interface{}) {
	s.tags = tags
}

// GetCounter returns a SynchronizedCounter assigned to the specified key with tags, creating a new one
// if necessary.
func (s *TaggedMeasurementSet) GetCounter(key string) *SynchronizedCounter {
	return s.MeasurementSet.GetCounter(MetricWithTags(key, s.tags))
}

// GetAggregator returns a SynchronizedAggregator assigned to the specified key with tags, creating a new one
// if necessary.
func (s *TaggedMeasurementSet) GetAggregator(key string) *SynchronizedAggregator {
	return s.MeasurementSet.GetAggregator(MetricWithTags(key, s.tags))
}

// Incr is a convenience function to get the specified Counter and call Incr on it. See Counter.Incr.
func (s *TaggedMeasurementSet) Incr(key string) {
	s.GetCounter(key).Incr()
}

// Add is a convenience function to get the specified Counter and call Add on it. See Counter.Add.
func (s *TaggedMeasurementSet) Add(key string, delta int64) {
	s.GetCounter(key).Add(delta)
}

// AddInt is a convenience function to get the specified Counter and call AddInt on it. See
// Counter.AddInt.
func (s *TaggedMeasurementSet) AddInt(key string, delta int) {
	s.GetCounter(key).AddInt(delta)
}

// UpdateAggregatorValue is a convenience to get the specified Aggregator and call UpdateValue on it.
// See Aggregator.UpdateValue.
func (s *TaggedMeasurementSet) UpdateAggregatorValue(key string, val float64) {
	s.GetAggregator(key).UpdateValue(val)
}

// UpdateAggregator is a convenience to get the specified Aggregator and call Update on it. See Aggregator.Update.
func (s *TaggedMeasurementSet) UpdateAggregator(key string, other Aggregator) {
	s.GetAggregator(key).Update(other)
}

// Merge takes a MeasurementSetReport and merges all of it Counters and Aggregators into this MeasurementSet.
// This in turn calls Counter.Add for each Counter in the report, and Aggregator.Update for each Aggregator in
// the report. Any keys that do not exist in this MeasurementSet will be created.
func (s *TaggedMeasurementSet) Merge(report *MeasurementSetReport) {
	for key, value := range report.Counts {
		s.GetCounter(key).Add(value)
	}
	for key, agg := range report.Aggregators {
		s.GetAggregator(key).Update(agg)
	}
}
