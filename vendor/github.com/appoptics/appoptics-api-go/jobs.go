package appoptics

import "fmt"

// Job is the representation of a task happening in the AppOptics cloud
type Job struct {
	ID       int                 `json:"id"`
	State    string              `json:"state"`
	Progress float64             `json:"progress,omitempty"`
	Output   string              `json:"output,omitempty"`
	Errors   map[string][]string `json:"errors,omitempty"`
}

type JobsCommunicator interface {
	Retrieve(int) (*Job, error)
}

type JobsService struct {
	client *Client
}

func NewJobsService(c *Client) *JobsService {
	return &JobsService{c}
}

// Retrieve gets the Job identified by the provided ID
func (js *JobsService) Retrieve(id int) (*Job, error) {
	path := fmt.Sprintf("jobs/%d", id)
	req, err := js.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	job := &Job{}

	_, err = js.client.Do(req, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}
