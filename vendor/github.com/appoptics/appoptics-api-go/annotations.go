package appoptics

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// AnnotationStream is a group of AnnotationEvents with a common name, representing a timeseries of occurrences of similar events
type AnnotationStream struct {
	Name        string                         `json:"name"`
	DisplayName string                         `json:"display_name,omitempty"`
	Events      []map[string][]AnnotationEvent `json:"events,omitempty"` // keys are Source names
}

// AnnotationEvent is the main data structure for the Annotations API
type AnnotationEvent struct {
	ID          int              `json:"id"`
	Title       string           `json:"title"`
	Source      string           `json:"source,omitempty"`
	Description string           `json:"description,omitempty"`
	Links       []AnnotationLink `json:"links,omitempty"`
	StartTime   int64            `json:"start_time,omitempty"`
	EndTime     int64            `json:"end_time,omitempty"`
}

// AnnotationLink represents the Link metadata for on the AnnotationEvent
type AnnotationLink struct {
	Rel   string `json:"rel"`
	Href  string `json:"href"`
	Label string `json:"label,omitempty"`
}

type RetrieveAnnotationsRequest struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Sources   []string
}

type ListAnnotationsResponse struct {
	AnnotationStreams []*AnnotationStream `json:"annotations"`
	Query             QueryInfo           `json:"query"`
}

// AnnotationsCommunicator provides an interface to the Annotations API from AppOptics
type AnnotationsCommunicator interface {
	List(*string) (*ListAnnotationsResponse, error)
	Retrieve(*RetrieveAnnotationsRequest) (*AnnotationStream, error)
	RetrieveEvent(string, int) (*AnnotationEvent, error)
	Create(*AnnotationEvent, string) (*AnnotationEvent, error)
	UpdateStream(string, string) error
	UpdateEvent(string, int, *AnnotationLink) (*AnnotationLink, error)
	Delete(string) error
}

type AnnotationsService struct {
	client *Client
}

func NewAnnotationsService(c *Client) *AnnotationsService {
	return &AnnotationsService{c}
}

// List retrieves paginated AnnotationEvents for all streams with name LIKE argument string
func (as *AnnotationsService) List(streamNameSearch *string) (*ListAnnotationsResponse, error) {
	var (
		path        string
		annotations *ListAnnotationsResponse
	)

	if streamNameSearch != nil {
		path = fmt.Sprintf("annotations?name=%s", *streamNameSearch)
	} else {
		path = "annotations"
	}

	req, err := as.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	_, err = as.client.Do(req, &annotations)

	if err != nil {
		return nil, err
	}

	return annotations, nil
}

// Retrieve fetches all AnnotationEvents matching the provided sources
func (as *AnnotationsService) Retrieve(retReq *RetrieveAnnotationsRequest) (*AnnotationStream, error) {
	stream := &AnnotationStream{}
	path := fmt.Sprintf("annotations/%s", retReq.Name)
	req, err := as.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = retReq.queryString(req)

	_, err = as.client.Do(req, stream)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (retReq *RetrieveAnnotationsRequest) queryString(req *http.Request) string {
	q := req.URL.Query()

	if !retReq.StartTime.IsZero() {
		q.Add("start_time", strconv.FormatInt(retReq.StartTime.Unix(), 10))
	}

	if !retReq.EndTime.IsZero() {
		q.Add("end_time", strconv.FormatInt(retReq.EndTime.Unix(), 10))
	}

	if len(retReq.Sources) > 0 {
		for _, source := range retReq.Sources {
			q.Add("sources[]", source)
		}
	}
	return q.Encode()
}

// RetrieveEvent returns a single event identified by an integer ID from a given stream
func (as *AnnotationsService) RetrieveEvent(streamName string, id int) (*AnnotationEvent, error) {
	event := &AnnotationEvent{}
	path := fmt.Sprintf("annotations/%s/%d", streamName, id)
	req, err := as.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	_, err = as.client.Do(req, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// Create makes an AnnotationEvent on the stream with the given name
func (as *AnnotationsService) Create(event *AnnotationEvent, streamName string) (*AnnotationEvent, error) {
	path := fmt.Sprintf("annotations/%s", streamName)
	req, err := as.client.NewRequest("POST", path, event)
	if err != nil {
		return nil, err
	}

	createdEvent := &AnnotationEvent{}

	_, err = as.client.Do(req, createdEvent)

	if err != nil {
		return nil, err
	}

	return createdEvent, nil
}

// UpdateStream updates the display name of the stream
func (as *AnnotationsService) UpdateStream(streamName, displayName string) error {
	path := fmt.Sprintf("annotations/%s", streamName)
	jsonTemplate := `{"display_name": %s}`
	req, err := as.client.NewRequest("POST", path, fmt.Sprintf(jsonTemplate, displayName))
	if err != nil {
		return err
	}

	_, err = as.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// UpdateEvent adds a link to an annotation Event
func (as *AnnotationsService) UpdateEvent(streamName string, id int, link *AnnotationLink) (*AnnotationLink, error) {
	newLink := &AnnotationLink{}
	path := fmt.Sprintf("annotations/%s/%d/links", streamName, id)
	req, err := as.client.NewRequest("POST", path, link)
	if err != nil {
		return nil, err
	}

	_, err = as.client.Do(req, newLink)

	if err != nil {
		return nil, err
	}

	return newLink, nil
}

// Delete deletes the annotation stream matching the provided name
func (as *AnnotationsService) Delete(streamName string) error {
	path := fmt.Sprintf("annotations/%s", streamName)
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
