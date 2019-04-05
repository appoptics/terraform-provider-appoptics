package appoptics

import (
	"fmt"
)

type Alert struct {
	ID           int                    `json:"id,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Active       bool                   `json:"active,omitempty"`
	RearmSeconds int                    `json:"rearm_seconds,omitempty"`
	Conditions   []*AlertCondition      `json:"conditions,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
	Services     []*Service             `json:"services,omitempty"` // correspond to IDs of Service objects
	CreatedAt    int                    `json:"created_at,omitempty"`
	UpdatedAt    int                    `json:"updated_at,omitempty"`
}

type AlertCondition struct {
	ID              int     `json:"id,omitempty"`
	Type            string  `json:"type,omitempty"`
	MetricName      string  `json:"metric_name,omitempty"`
	Threshold       float64 `json:"threshold,omitempty"`
	SummaryFunction string  `json:"summary_function,omitempty"`
	Duration        int     `json:"duration,omitempty"`
	Tags            []*Tag  `json:"tags,omitempty"`
}

type AlertStatus struct {
	Alert  Alert  `json:"alert,omitempty"`
	Status string `json:"status,omitempty"`
}

type AlertsResponse struct {
	Query  QueryInfo `json:"query,omitempty"`
	Alerts []*Alert  `json:"alerts,omitempty"`
}

type AlertsService struct {
	client *Client
}

type AlertsCommunicator interface {
	List() (*AlertsResponse, error)
	Retrieve(int) (*Alert, error)
	Create(*Alert) (*Alert, error)
	Update(*Alert) error
	AssociateToService(int, int) error
	DisassociateFromService(alertId, serviceId int) error
	Delete(int) error
	Status(int) (*AlertStatus, error)
}

func NewAlertsService(c *Client) *AlertsService {
	return &AlertsService{c}
}

// List retrieves all Alerts
func (as *AlertsService) List() (*AlertsResponse, error) {
	req, err := as.client.NewRequest("GET", "alerts", nil)
	if err != nil {
		return nil, err
	}

	alertsResponse := &AlertsResponse{}

	_, err = as.client.Do(req, &alertsResponse)

	if err != nil {
		return nil, err
	}

	return alertsResponse, nil
}

// Retrieve returns the Alert identified by the parameter
func (as *AlertsService) Retrieve(id int) (*Alert, error) {
	alert := &Alert{}
	path := fmt.Sprintf("alerts/%d", id)
	req, err := as.client.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	_, err = as.client.Do(req, alert)
	if err != nil {
		return nil, err
	}

	return alert, nil
}

// Create creates the Alert
func (as *AlertsService) Create(a *Alert) (*Alert, error) {
	req, err := as.client.NewRequest("POST", "alerts", a)
	if err != nil {
		return nil, err
	}

	createdAlert := &Alert{}

	_, err = as.client.Do(req, createdAlert)
	if err != nil {
		return nil, err
	}

	return createdAlert, nil
}

// Update updates the Alert
func (as *AlertsService) Update(a *Alert) error {
	path := fmt.Sprintf("alerts/%d", a.ID)
	req, err := as.client.NewRequest("PUT", path, a)
	if err != nil {
		return err
	}
	_, err = as.client.Do(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// AssociateToService updates the Alert to allow assign it to the Service identified
func (as *AlertsService) AssociateToService(alertId, serviceId int) error {
	path := fmt.Sprintf("alerts/%d/services", alertId)
	bodyStruct := struct {
		ID int `json:"service"`
	}{serviceId}
	req, err := as.client.NewRequest("POST", path, bodyStruct)

	if err != nil {
		return err
	}

	_, err = as.client.Do(req, nil)

	if err != nil {
		return err
	}
	return nil
}

// DisassociateFromService updates the Alert to remove the Service identified
func (as *AlertsService) DisassociateFromService(alertId, serviceId int) error {
	path := fmt.Sprintf("alerts/%d/services/%d", alertId, serviceId)
	req, err := as.client.NewRequest("DELETE", path, nil)

	if err != nil {
		return err
	}

	_, err = as.client.Do(req, nil)

	if err != nil {
		return err
	}
	return nil
}

// Delete deletes the Alert
func (as *AlertsService) Delete(id int) error {
	path := fmt.Sprintf("alerts/%d", id)
	req, err := as.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = as.client.Do(req, nil)

	return err
}

// Status returns the Alert's status
func (as *AlertsService) Status(id int) (*AlertStatus, error) {
	path := fmt.Sprintf("alerts/%d/status", id)
	req, err := as.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	alertStatus := &AlertStatus{}

	_, err = as.client.Do(req, alertStatus)

	if err != nil {
		return nil, err
	}

	return alertStatus, nil
}
