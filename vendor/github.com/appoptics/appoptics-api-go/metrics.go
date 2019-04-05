package appoptics

import "fmt"

// Metric is an AppOptics Metric
type Metric struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Type        string           `json:"type"`
	Period      int              `json:"period,omitempty"`
	DisplayName string           `json:"display_name,omitempty"`
	Composite   string           `json:"composite,omitempty"`
	Attributes  MetricAttributes `json:"attributes,omitempty"`
}

// MetricAttributes are key/value pairs of metadata about the Metric
type MetricAttributes struct {
	Color             string      `json:"color,omitempty"`
	DisplayMax        interface{} `json:"display_max,omitempty"`
	DisplayMin        interface{} `json:"display_min,omitempty"`
	DisplayUnitsLong  string      `json:"display_units_long,omitempty"`
	DisplayUnitsShort string      `json:"display_units_short,omitempty"`
	DisplayStacked    bool        `json:"display_stacked,omitempty"`
	CreatedByUA       string      `json:"created_by_ua,omitempty"`
	GapDetection      bool        `json:"gap_detection,omitempty"`
	Aggregate         bool        `json:"aggregate,omitempty"`
}

type MetricsResponse struct {
	Query   QueryInfo `json:"query,omitempty"`
	Metrics []*Metric `json:"metrics,omitempty"`
}

type MetricsService struct {
	client *Client
}

// MetricUpdatePayload will apply the state represented by Attributes to the Metrics identified by Names
type MetricUpdatePayload struct {
	Names      []string         `json:"names"`
	Attributes MetricAttributes `json:"attributes"`
}

type MetricsCommunicator interface {
	List() (*MetricsResponse, error)
	Retrieve(string) (*Metric, error)
	Create(*Metric) (*Metric, error)
	Update(string, *Metric) error
	Delete(string) error
}

func NewMetricsService(c *Client) *MetricsService {
	return &MetricsService{c}
}

// List lists the Metrics in the organization identified by the AppOptics token
func (ms *MetricsService) List() (*MetricsResponse, error) {
	req, err := ms.client.NewRequest("GET", "metrics", nil)
	if err != nil {
		return nil, err
	}

	metricsResponse := &MetricsResponse{}

	_, err = ms.client.Do(req, &metricsResponse)

	if err != nil {
		return nil, err
	}

	return metricsResponse, nil
}

// Retrieve fetches the Metric identified by the given name
func (ms *MetricsService) Retrieve(name string) (*Metric, error) {
	metric := &Metric{}
	path := fmt.Sprintf("metrics/%s", name)
	req, err := ms.client.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	_, err = ms.client.Do(req, metric)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

// Create creates the Metric in the organization identified by the AppOptics token
func (ms *MetricsService) Create(m *Metric) (*Metric, error) {
	path := fmt.Sprintf("metrics/%s", m.Name)
	req, err := ms.client.NewRequest("PUT", path, m)
	if err != nil {
		return nil, err
	}

	createdMetric := &Metric{}

	_, err = ms.client.Do(req, createdMetric)
	if err != nil {
		return nil, err
	}

	return createdMetric, nil
}

// Update updates the Metric with the given name, setting it to match the Metric pointer argument
func (ms *MetricsService) Update(originalName string, m *Metric) error {
	path := fmt.Sprintf("metrics/%s", originalName)
	req, err := ms.client.NewRequest("PUT", path, m)

	if err != nil {
		return err
	}

	_, err = ms.client.Do(req, nil)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the Metric matching the name argument
func (ms *MetricsService) Delete(name string) error {
	path := fmt.Sprintf("metrics/%s", name)
	req, err := ms.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = ms.client.Do(req, nil)

	return err
}
