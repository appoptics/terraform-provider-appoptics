package appoptics

import "fmt"

type Service struct {
	ID       int               `json:"id,omitempty"`
	Type     string            `json:"type,omitempty"`
	Settings map[string]string `json:"settings,omitempty"`
	Title    string            `json:"title,omitempty"`
}

type ServicesCommunicator interface {
	List() (*ListServicesResponse, error)
	Retrieve(int) (*Service, error)
	Create(*Service) (*Service, error)
	Update(*Service) error
	Delete(int) error
}

type ListServicesResponse struct {
	Query    QueryInfo  `json:"query,omitempty"`
	Services []*Service `json:"services,omitempty"`
}

type ServicesService struct {
	client *Client
}

func NewServiceService(c *Client) *ServicesService {
	return &ServicesService{c}
}

// List retrieves all Services
func (ss *ServicesService) List() (*ListServicesResponse, error) {
	req, err := ss.client.NewRequest("GET", "services", nil)
	if err != nil {
		return nil, err
	}

	servicesResponse := &ListServicesResponse{}

	_, err = ss.client.Do(req, &servicesResponse)

	if err != nil {
		return nil, err
	}

	return servicesResponse, nil
}

// Retrieve returns the Service identified by the parameter
func (ss *ServicesService) Retrieve(id int) (*Service, error) {
	service := &Service{}
	path := fmt.Sprintf("services/%d", id)
	req, err := ss.client.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	_, err = ss.client.Do(req, service)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// Create creates the Service
func (ss *ServicesService) Create(s *Service) (*Service, error) {
	req, err := ss.client.NewRequest("POST", "services", s)
	if err != nil {
		return nil, err
	}

	createdService := &Service{}

	_, err = ss.client.Do(req, createdService)
	if err != nil {
		return nil, err
	}

	return createdService, nil
}

// Update updates the Service
func (ss *ServicesService) Update(s *Service) error {
	path := fmt.Sprintf("services/%d", s.ID)
	req, err := ss.client.NewRequest("PUT", path, s)
	if err != nil {
		return err
	}

	_, err = ss.client.Do(req, nil)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the Service
func (ss *ServicesService) Delete(id int) error {
	path := fmt.Sprintf("services/%d", id)
	req, err := ss.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = ss.client.Do(req, nil)

	return err
}
