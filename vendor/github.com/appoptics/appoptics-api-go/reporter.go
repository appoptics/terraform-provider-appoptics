package appoptics

import (
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	outputMeasurementsIntervalSeconds = 15
	outputMeasurementsInterval        = outputMeasurementsIntervalSeconds * time.Second
	maxRetries                        = 3
	maxMeasurementsPerBatch           = 1000
)

// Reporter provides a way to persist data from a set collection of Aggregators and Counters at a regular interval
type Reporter struct {
	measurementSet   *MeasurementSet
	measurementsComm MeasurementsCommunicator
	prefix           string

	batchChan             chan *MeasurementsBatch
	measurementSetReports chan *MeasurementSetReport

	globalTags map[string]string
}

// NewReporter returns a reporter for a given MeasurementSet, providing a way to sync metric information
// to AppOptics for a collection of running metrics.
func NewReporter(measurementSet *MeasurementSet, communicator MeasurementsCommunicator, prefix string) *Reporter {
	r := &Reporter{
		measurementSet:        measurementSet,
		measurementsComm:      communicator,
		prefix:                prefix,
		batchChan:             make(chan *MeasurementsBatch, 100),
		measurementSetReports: make(chan *MeasurementSetReport, 1000),
	}
	r.initGlobalTags()
	return r
}

// Start kicks off two goroutines that help batch and report metrics measurements to AppOptics.
func (r *Reporter) Start() {
	go r.postMeasurementBatches()
	go r.flushReportsForever()
}

func (r *Reporter) initGlobalTags() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "na"
	}
	r.globalTags = map[string]string{
		"hostname": hostname + os.Getenv("HOST_SUFFIX"),
	}
}

func (r *Reporter) postMeasurementBatches() {
	for batch := range r.batchChan {
		tryCount := 0
		for {
			log.Debug("Uploading AppOptics measurements batch", "time", time.Unix(batch.Time, 0), "numMeasurements", len(batch.Measurements), "globalTags", r.globalTags)
			_, err := r.measurementsComm.Create(batch)
			if err == nil {
				break
			}
			tryCount++
			aborting := tryCount == maxRetries
			log.Error("Error uploading AppOptics measurements batch", "err", err, "tryCount", tryCount, "aborting", aborting)
			if aborting {
				break
			}
		}
	}
}

func (r *Reporter) flushReport(report *MeasurementSetReport) {
	batchTimeUnixSecs := (time.Now().Unix() / outputMeasurementsIntervalSeconds) * outputMeasurementsIntervalSeconds

	var batch *MeasurementsBatch
	resetBatch := func() {
		batch = &MeasurementsBatch{
			Time:   batchTimeUnixSecs,
			Period: outputMeasurementsIntervalSeconds,
		}
	}
	flushBatch := func() {
		r.batchChan <- batch
	}
	addMeasurement := func(measurement Measurement) {
		batch.Measurements = append(batch.Measurements, measurement)
		// AppOptics API docs advise sending very large numbers of metrics in multiple HTTP requests; so we'll flush
		// batches of 500 measurements at a time.
		if len(batch.Measurements) >= maxMeasurementsPerBatch {
			flushBatch()
			resetBatch()
		}
	}
	resetBatch()
	report.Counts["num_measurements"] = int64(len(report.Counts)) + int64(len(report.Aggregators)) + 1
	for key, value := range report.Counts {
		metricName, tags := parseMeasurementKey(key)
		m := Measurement{
			Name: r.prefix + regexpIllegalNameChars.ReplaceAllString(metricName, "_"),
			Tags: r.mergeGlobalTags(tags),
		}
		if value != 0 {
			m.Value = float64(value)
		}
		addMeasurement(m)
	}
	// TODO: refactor to use Aggregator methods
	for key, agg := range report.Aggregators {
		metricName, tags := parseMeasurementKey(key)
		m := Measurement{
			Name: r.prefix + regexpIllegalNameChars.ReplaceAllString(metricName, "_"),
			Tags: r.mergeGlobalTags(tags),
		}
		if agg.Sum != 0 {
			m.Sum = agg.Sum
		}
		if agg.Count != 0 {
			m.Count = agg.Count
		}
		if agg.Min != 0 {
			m.Min = agg.Min
		}
		if agg.Max != 0 {
			m.Max = agg.Max
		}
		if agg.Last != 0 {
			m.Last = agg.Last
		}
		addMeasurement(m)
	}
	if len(batch.Measurements) > 0 {
		flushBatch()
	}
}

func (r *Reporter) flushReportsForever() {
	// Sleep for a random duration between 0 and outputMeasurementsInterval in order to randomize the counters output cycle.
	time.Sleep(time.Duration(rand.Int63n(int64(outputMeasurementsInterval))))
	report := r.measurementSet.Reset()
	r.flushReport(report)
	// After the initial random sleep, start a regular interval timer. This will output measurements at a consistent time
	// modulo outputMeasurementsInterval.
	for range time.Tick(outputMeasurementsInterval) {
		report := r.measurementSet.Reset()
		r.flushReport(report)
	}
}

func (r *Reporter) mergeGlobalTags(tags map[string]string) map[string]string {
	if tags == nil {
		return r.globalTags
	}

	if r.globalTags == nil {
		return tags
	}

	for k, v := range r.globalTags {
		if _, ok := tags[v]; !ok {
			tags[k] = v
		}
	}

	return tags
}
