package appoptics

import (
	"math/rand"
	"time"
)

type MultiReporter struct {
	measurementSet *MeasurementSet
	reporters      []*Reporter
}

func NewMultiReporter(m *MeasurementSet, reporters []*Reporter) *MultiReporter {
	return &MultiReporter{measurementSet: m, reporters: reporters}
}

func (m *MultiReporter) Start() {
	for _, r := range m.reporters {
		go r.postMeasurementBatches()
	}

	go m.flushReportsForever()
}

func (m *MultiReporter) flushReport(report *MeasurementSetReport) {
	for _, r := range m.reporters {
		r.flushReport(report)
	}
}

func (m *MultiReporter) flushReportsForever() {
	// Sleep for a random duration between 0 and outputMeasurementsInterval in order to randomize the counters output cycle.
	time.Sleep(time.Duration(rand.Int63n(int64(outputMeasurementsInterval))))
	m.flushReport(m.measurementSet.Reset())

	// After the initial random sleep, start a regular interval timer. This will output measurements at a consistent time
	// modulo outputMeasurementsInterval.
	for range time.Tick(outputMeasurementsInterval) {
		m.flushReport(m.measurementSet.Reset())
	}
}
