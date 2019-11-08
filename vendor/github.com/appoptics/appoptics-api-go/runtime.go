package appoptics

import (
	"runtime"
	"time"
)

const (
	runtimeRecordInterval = 10 * time.Second
)

func RecordRuntimeMetrics(m *MeasurementSet) {
	go recordRuntimeMetrics(m)
}

func recordRuntimeMetrics(m *MeasurementSet) {
	var (
		memStats       = &runtime.MemStats{}
		lastSampleTime = time.Now()
		lastPauseNs    uint64
		lastNumGC      uint64
	)

	for {
		runtime.ReadMemStats(memStats)

		now := time.Now()

		m.UpdateAggregatorValue("go.goroutines", float64(runtime.NumGoroutine()))
		m.UpdateAggregatorValue("go.memory.allocated", float64(memStats.Alloc))
		m.UpdateAggregatorValue("go.memory.mallocs", float64(memStats.Mallocs))
		m.UpdateAggregatorValue("go.memory.frees", float64(memStats.Frees))
		m.UpdateAggregatorValue("go.memory.gc.total_pause", float64(memStats.PauseTotalNs))
		m.UpdateAggregatorValue("go.memory.gc.heap", float64(memStats.HeapAlloc))
		m.UpdateAggregatorValue("go.memory.gc.stack", float64(memStats.StackInuse))

		if lastPauseNs > 0 {
			pauseSinceLastSample := memStats.PauseTotalNs - lastPauseNs
			m.UpdateAggregatorValue("go.memory.gc.pause_per_second", float64(pauseSinceLastSample)/runtimeRecordInterval.Seconds())
		}
		lastPauseNs = memStats.PauseTotalNs

		countGC := int(uint64(memStats.NumGC) - lastNumGC)
		if lastNumGC > 0 {
			diff := float64(countGC)
			diffTime := now.Sub(lastSampleTime).Seconds()
			m.UpdateAggregatorValue("go.memory.gc.gc_per_second", diff/diffTime)
		}

		if countGC > 0 {
			if countGC > 256 {
				countGC = 256
			}

			for i := 0; i < countGC; i++ {
				idx := int((memStats.NumGC-uint32(i))+255) % 256
				pause := time.Duration(memStats.PauseNs[idx])
				m.UpdateAggregatorValue("go.memory.gc.pause", float64(pause))
			}
		}

		lastNumGC = uint64(memStats.NumGC)
		lastSampleTime = now

		time.Sleep(runtimeRecordInterval)
	}
}
