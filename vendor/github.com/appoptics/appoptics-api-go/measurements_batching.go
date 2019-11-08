package appoptics

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// MeasurementsBatch is a collection of Measurements persisted to the API at the same time.
// It can optionally have tags that are applied to all contained Measurements.
type MeasurementsBatch struct {
	// Measurements is the collection of timeseries entries being sent to the server
	Measurements []Measurement `json:"measurements,omitempty"`
	// Period is a slice of time measured in seconds, used in service-side aggregation
	Period int64 `json:"period,omitempty"`
	// Time is a Unix epoch timestamp used to align a group of Measurements on a time boundary
	Time int64 `json:"time"`
	// Tags are key-value identifiers that will be applied to all Measurements in the batch
	Tags *map[string]string `json:"tags,omitempty"`
}

// BatchPersister implements persistence to AppOptics and enforces error limits
type BatchPersister struct {
	// mc is the MeasurementsCommunicator used to talk to the AppOptics API
	mc MeasurementsCommunicator
	// errors holds the errors received in attempting to persist to AppOptics
	errors []error
	// errorLimit is the number of persistence errors that will be tolerated
	errorLimit int
	// prepChan is a channel of Measurements slices
	prepChan chan []Measurement
	// batchChan is used to create MeasurementsBatches for persistence to AppOptics
	batchChan chan *MeasurementsBatch
	// stopBatchingChan is used to cease persisting MeasurementsBatches to AppOptics
	stopBatchingChan chan bool
	// stopPersistingChan is used to cease persisting MeasurementsBatches to AppOptics
	stopPersistingChan chan bool
	// stopErrorChan is used to cease the error checking for/select
	stopErrorChan chan bool
	// errorChan is used to tally errors that occur in batching/persisting
	errorChan chan error
	// maximumPushInterval is the max time (in milliseconds) to wait before pushing a batch whether its length is equal
	// to the MeasurementPostMaxBatchSize or not
	maximumPushInterval int
	// sendStats is a flag for whether to persist to AppOptics or simply print messages to stdout
	sendStats bool
}

// NewBatchPersister sets up a new instance of batched persistence capabilities using the provided MeasurementsCommunicator
func NewBatchPersister(mc MeasurementsCommunicator, sendStats bool) *BatchPersister {
	return &BatchPersister{
		mc:                  mc,
		errorLimit:          DefaultPersistenceErrorLimit,
		prepChan:            make(chan []Measurement),
		batchChan:           make(chan *MeasurementsBatch),
		stopBatchingChan:    make(chan bool),
		stopErrorChan:       make(chan bool),
		stopPersistingChan:  make(chan bool),
		errorChan:           make(chan error),
		errors:              []error{},
		maximumPushInterval: 2000,
		sendStats:           sendStats,
	}
}

func NewMeasurementsBatch(m []Measurement, tags *map[string]string) *MeasurementsBatch {
	return &MeasurementsBatch{
		Measurements: m,
		Tags:         tags,
		Time:         time.Now().UTC().Unix(),
	}
}

// MeasurementsSink gives calling code write-only access to the Measurements prep channel
func (bp *BatchPersister) MeasurementsSink() chan<- []Measurement {
	return bp.prepChan
}

// MeasurementsStopBatchingChannel gives calling code write-only access to the Measurements batching control channel
func (bp *BatchPersister) MeasurementsStopBatchingChannel() chan<- bool {
	return bp.stopBatchingChan
}

// MeasurementsErrorChannel gives calling code write-only access to the Measurements error channel
func (bp *BatchPersister) MeasurementsErrorChannel() chan<- error {
	return bp.errorChan
}

// MaxiumumPushIntervalMilliseconds returns the number of milliseconds the system will wait before pushing any
// accumulated Measurements to AppOptics
func (bp *BatchPersister) MaximumPushInterval() int {
	return bp.maximumPushInterval
}

// SetMaximumPushInterval sets the number of milliseconds the system will wait before pushing any accumulated
// Measurements to AppOptics
func (bp *BatchPersister) SetMaximumPushInterval(ms int) {
	bp.maximumPushInterval = ms
}

// batchMeasurements reads slices of Measurements off a channel and packages them into batches conforming to the
// limitations imposed by the API. If Measurements are arriving slowly, collected Measurements will be pushed on an
// interval defined by maximumPushIntervalMilliseconds
func (bp *BatchPersister) batchMeasurements() {
	var currentMeasurements = []Measurement{}
	ticker := time.NewTicker(time.Millisecond * time.Duration(bp.maximumPushInterval))
LOOP:
	for {
		select {
		case receivedMeasurements := <-bp.prepChan:
			currentMeasurements = append(currentMeasurements, receivedMeasurements...)
			if len(currentMeasurements) >= MeasurementPostMaxBatchSize {
				bp.batchChan <- &MeasurementsBatch{Measurements: currentMeasurements[:MeasurementPostMaxBatchSize]}
				currentMeasurements = currentMeasurements[MeasurementPostMaxBatchSize:]
			}
		case <-ticker.C:
			if len(currentMeasurements) > 0 {
				pushBatch := &MeasurementsBatch{}
				if len(currentMeasurements) >= MeasurementPostMaxBatchSize {
					pushBatch.Measurements = currentMeasurements[:MeasurementPostMaxBatchSize]
					bp.batchChan <- pushBatch
					currentMeasurements = currentMeasurements[MeasurementPostMaxBatchSize:]
				} else {
					pushBatch.Measurements = currentMeasurements
					bp.batchChan <- pushBatch
					currentMeasurements = []Measurement{}
				}
			}
		case <-bp.stopBatchingChan:
			ticker.Stop()
			if len(currentMeasurements) > 0 {
				if len(bp.errors) < bp.errorLimit {
					bp.batchChan <- &MeasurementsBatch{Measurements: currentMeasurements[:MeasurementPostMaxBatchSize]}
				}
			}
			close(bp.batchChan)
			bp.stopPersistingChan <- true
			bp.stopErrorChan <- true
			break LOOP
		}
	}
}

// BatchAndPersistMeasurementsForever continually packages up Measurements from the channel returned by MeasurementSink()
// and persists them to AppOptics.
func (bp *BatchPersister) BatchAndPersistMeasurementsForever() {
	go bp.batchMeasurements()
	go bp.persistBatches()
	go bp.managePersistenceErrors()
}

// persistBatches reads maximal slices of Measurements off a channel and persists them to the remote AppOptics
// API. Errors are placed on the error channel.
func (bp *BatchPersister) persistBatches() {
	ticker := time.NewTicker(time.Millisecond * 500)
LOOP:
	for {
		select {
		case <-ticker.C:
			batch := <-bp.batchChan
			if batch != nil {
				err := bp.persistBatch(batch)
				if err != nil {
					bp.errorChan <- err
				}
			}
		case <-bp.stopPersistingChan:
			if len(bp.errors) > bp.errorLimit {
				batch := <-bp.batchChan
				if batch != nil {
					bp.persistBatch(batch)
				}
			}
			ticker.Stop()
			break LOOP
		}
	}
}

// managePersistenceErrors tracks errors on the provided channel and sends a stop signal if the ErrorLimit is reached.
func (bp *BatchPersister) managePersistenceErrors() {
LOOP:
	for {
		select {
		case err := <-bp.errorChan:
			bp.errors = append(bp.errors, err)
			if len(bp.errors) == bp.errorLimit {
				bp.stopBatchingChan <- true
				break LOOP
			}
		case <-bp.stopErrorChan:
			break LOOP
		}
	}
}

// persistBatch sends to the remote AppOptics endpoint unless config.SendStats() returns false, when it prints to stdout
func (bp *BatchPersister) persistBatch(batch *MeasurementsBatch) error {
	if bp.sendStats {
		// TODO: make this conditional upon log level
		log.Printf("persisting %d Measurements to AppOptics\n", len(batch.Measurements))
		resp, err := bp.mc.Create(batch)
		if resp == nil {
			fmt.Println("response is nil")
			return err
		}
		// TODO: make this conditional upon log level
		dumpResponse(resp)
	} else {
		// TODO: make this more verbose upon log level
		log.Printf("received %d Measurements for persistence\n", len(batch.Measurements))
		//printMeasurements(batch.Measurements)
	}
	return nil
}
