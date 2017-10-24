package librato

import (
	"fmt"
	"net/http"
)

// AnnotationsService handles communication with the Librato API methods related to
// Annotations.
type AnnotationsService struct {
	client *Client
}

// Annotation represents a Librato Annotation.
type Annotation struct {
	Name        *string          `json:"name"`
	Title       *string          `json:"title"`
	Source      *string          `json:"source,omitempty"`
	Description *string          `json:"description,omitempty"`
	Links       []AnnotationLink `json:"links,omitempty"`
	StartTime   *uint            `json:"start_time,omitempty"`
	EndTime     *uint            `json:"end_time,omitempty"`
}

// AnnotationLink represents a Librato Annotation Link.
type AnnotationLink struct {
	Label *string `json:"label,omitempty"`
	Rel   *string `json:"rel"`
	Href  *string `json:"href,omitempty"`
}

// Create an Annotation
//
// Librato API docs: https://www.librato.com/docs/api/?shell#create-an-Annotation
func (a *AnnotationsService) Create(annotation *Annotation) (*Annotation, *http.Response, error) {
	u := fmt.Sprintf("annotations/%s", *annotation.Name)
	req, err := a.client.NewRequest("POST", u, annotation)
	if err != nil {
		return nil, nil, err
	}

	an := new(Annotation)
	resp, err := a.client.Do(req, an)
	if err != nil {
		return nil, resp, err
	}

	return an, resp, err
}
