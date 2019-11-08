package appoptics

import (
	"fmt"
)

type Chart struct {
	ID           int         `json:"id,omitempty"`
	Name         string      `json:"name,omitempty"`
	Type         string      `json:"type,omitempty"`
	Streams      []Stream    `json:"streams,omitempty"`
	Min          float64     `json:"min,omitempty"`
	Max          float64     `json:"max,omitempty"`
	Label        string      `json:"label,omitempty"`
	RelatedSpace int         `json:"related_space,omitempty"`
	Thresholds   []Threshold `json:"threshold,omitempty"`
}

type Stream struct {
	ID                 int    `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	Metric             string `json:"metric,omitempty"`
	Composite          string `json:"composite,omitempty"`
	Type               string `json:"type,omitempty"`
	Tags               []Tag  `json:"tags,omitempty"`
	GroupFunction      string `json:"group_function,omitempty"` // valid: average, sum, min, max
	GroupBy            string `json:"group_by,omitempty"`
	SummaryFunction    string `json:"summary_function,omitempty"`    // valid: average, sum, min, max, count
	DownsampleFunction string `json:"downsample_function,omitempty"` // valid: average, min, max, sum, count
	Color              string `json:"color,omitempty"`
	UnitsShort         string `json:"units_short,omitempty"`
	UnitsLong          string `json:"units_long,omitempty"`
	TransformFunction  string `json:"transform_function,omitempty"`
	Period             int    `json:"period,omitempty"`
	Min                int    `json:"min,omitempty"`
	Max                int    `json:"max,omitempty"`
}

type Threshold struct {
	Operator string  `json:"operator,omitempty"`
	Value    float64 `json:"value,omitempty"`
	Type     string  `json:"type,omitempty"`
}

type ChartsCommunicator interface {
	List(int) ([]*Chart, error)
	Retrieve(int, int) (*Chart, error)
	Create(*Chart, int) (*Chart, error)
	Update(*Chart, int) (*Chart, error)
	Delete(int, int) error
}

type ChartsService struct {
	client *Client
}

func NewChartsService(c *Client) *ChartsService {
	return &ChartsService{c}
}

// List retrieves the Charts for the provided Space ID
func (cs *ChartsService) List(spaceId int) ([]*Chart, error) {
	path := fmt.Sprintf("spaces/%d/charts", spaceId)
	charts := []*Chart{}
	req, err := cs.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	_, err = cs.client.Do(req, &charts)

	if err != nil {
		return nil, err
	}

	return charts, nil
}

// Retrieve finds the Chart identified by the provided parameters
func (cs *ChartsService) Retrieve(chartId, spaceId int) (*Chart, error) {
	chart := &Chart{}
	path := fmt.Sprintf("spaces/%d/charts/%d", spaceId, chartId)
	req, err := cs.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	_, err = cs.client.Do(req, chart)
	if err != nil {
		return nil, err
	}

	return chart, nil
}

// Create creates the Chart in the Space
func (cs *ChartsService) Create(chart *Chart, spaceId int) (*Chart, error) {
	path := fmt.Sprintf("spaces/%d/charts", spaceId)
	req, err := cs.client.NewRequest("POST", path, chart)
	if err != nil {
		return nil, err
	}

	createdChart := &Chart{}

	_, err = cs.client.Do(req, createdChart)
	if err != nil {
		return nil, err
	}

	return createdChart, nil
}

// Update takes a Chart representing requested changes to an existing Chart on the server
// and returns the altered Chart from the server.
func (cs *ChartsService) Update(existingChart *Chart, spaceId int) (*Chart, error) {
	path := fmt.Sprintf("spaces/%d/charts/%d", spaceId, existingChart.ID)
	req, err := cs.client.NewRequest("PUT", path, existingChart)
	if err != nil {
		return nil, err
	}

	updatedChart := &Chart{}
	_, err = cs.client.Do(req, updatedChart)

	if err != nil {
		return nil, err
	}

	return updatedChart, nil

}

// Delete deletes the Chart from the Space
func (cs *ChartsService) Delete(chartId, spaceId int) error {
	path := fmt.Sprintf("spaces/%d/charts/%d", spaceId, chartId)
	req, err := cs.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = cs.client.Do(req, nil)

	return err
}
