package appoptics

import (
	"fmt"
	"log"
	"net/http"
)

const (
	// AggregationKey is the key in the Measurement attributes used to tell the AppOptics system to aggregate values
	AggregationKey = "aggregate"
)

// Measurement wraps the corresponding API construct: https://docs.appoptics.com/api/#measurements
// Each Measurement represents a single timeseries value for an associated Metric. If AppOptics receives a Measurement
// with a Name field that doesn't correspond to an existing Metric, a new Metric will be created.
type Measurement struct {
	Name       string                 `json:"name"`
	Tags       map[string]string      `json:"tags,omitempty"`
	Value      interface{}            `json:"value,omitempty"`
	Time       int64                  `json:"time,omitempty"`
	Count      interface{}            `json:"count,omitempty"`
	Sum        interface{}            `json:"sum,omitempty"`
	Min        interface{}            `json:"min,omitempty"`
	Max        interface{}            `json:"max,omitempty"`
	Last       interface{}            `json:"last,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// MeasurementsCommunicator defines an interface for communicating with the Measurements portion of the AppOptics API
type MeasurementsCommunicator interface {
	Create(*MeasurementsBatch) (*http.Response, error)
}

// MeasurementsService implements MeasurementsCommunicator
type MeasurementsService struct {
	client *Client
}

func NewMeasurementsService(c *Client) *MeasurementsService {
	return &MeasurementsService{c}
}

// NewMeasurement returns a Measurement with the given name and an empty attributes map
func NewMeasurement(name string) Measurement {
	attrs := make(map[string]interface{})
	return Measurement{
		Name:       name,
		Attributes: attrs,
	}
}

// Create persists the given MeasurementCollection to AppOptics
func (ms *MeasurementsService) Create(batch *MeasurementsBatch) (*http.Response, error) {
	req, err := ms.client.NewRequest("POST", "measurements", batch)

	if err != nil {
		log.Println("error creating request:", err)
		return nil, err
	}
	return ms.client.Do(req, nil)
}

// printMeasurements pretty-prints the supplied measurements to stdout
func printMeasurements(data []Measurement) {
	for _, measurement := range data {
		fmt.Printf("\nMetric name: '%s' \n", measurement.Name)
		fmt.Printf("\t value: %d \n", measurement.Value)
		fmt.Printf("\t\tTags: ")
		for k, v := range measurement.Tags {
			fmt.Printf("\n\t\t\t%s: %s", k, v)
		}
	}
}
