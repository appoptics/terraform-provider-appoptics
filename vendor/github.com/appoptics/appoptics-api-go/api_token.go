package appoptics

import "fmt"

type ApiToken struct {
	ID     *int    `json:"id,omitempty"`
	Active *bool   `json:"active,omitempty"`
	Name   *string `json:"name,omitempty"`
	Role   *string `json:"role,omitempty"`
	Token  *string `json:"token,omitempty"`
}

type ApiTokensResponse struct {
	Query     QueryInfo   `json:"query,omitempty"`
	ApiTokens []*ApiToken `json:"api_tokens,omitempty"`
}

type ApiTokensCommunicator interface {
	List() (*ApiTokensResponse, error)
	Retrieve(string) (*ApiTokensResponse, error)
	Create(*ApiToken) (*ApiToken, error)
	Update(*ApiToken) (*ApiToken, error)
	Delete(int) error
}

type ApiTokensService struct {
	client *Client
}

func NewApiTokensService(c *Client) *ApiTokensService {
	return &ApiTokensService{c}
}

// List retrieves all ApiTokens
func (ts *ApiTokensService) List() (*ApiTokensResponse, error) {
	req, err := ts.client.NewRequest("GET", "api_tokens", nil)
	if err != nil {
		return nil, err
	}

	apiResponse := &ApiTokensResponse{}

	_, err = ts.client.Do(req, &apiResponse)

	if err != nil {
		return nil, err
	}

	return apiResponse, nil
}

// Retrieve returns the ApiToken identified by the parameter
func (ts *ApiTokensService) Retrieve(name string) (*ApiTokensResponse, error) {
	tokenResponse := &ApiTokensResponse{}
	path := fmt.Sprintf("api_tokens/%s", name)
	req, err := ts.client.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	_, err = ts.client.Do(req, tokenResponse)
	if err != nil {
		return nil, err
	}

	return tokenResponse, nil
}

// Create creates the ApiToken
func (ts *ApiTokensService) Create(at *ApiToken) (*ApiToken, error) {
	req, err := ts.client.NewRequest("POST", "api_tokens", at)
	if err != nil {
		return nil, err
	}

	createdToken := &ApiToken{}

	_, err = ts.client.Do(req, createdToken)
	if err != nil {
		return nil, err
	}

	return createdToken, nil
}

// Update updates the ApiToken
func (ts *ApiTokensService) Update(at *ApiToken) (*ApiToken, error) {
	path := fmt.Sprintf("api_tokens/%d", at.ID)
	req, err := ts.client.NewRequest("PUT", path, at)
	if err != nil {
		return nil, err
	}

	updatedToken := &ApiToken{}

	_, err = ts.client.Do(req, updatedToken)

	if err != nil {
		return nil, err
	}

	return updatedToken, nil
}

// Delete deletes the ApiToken
func (ts *ApiTokensService) Delete(id int) error {
	path := fmt.Sprintf("api_tokens/%d", id)
	req, err := ts.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = ts.client.Do(req, nil)

	return err
}
