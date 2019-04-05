package appoptics

import (
	"fmt"
	"time"
)

// Snapshot represents a portrait of a Chart at a specific point in time
type Snapshot struct {
	Href      string                   `json:"href,omitempty"`
	JobHref   string                   `json:"job_href,omitempty"`
	ImageHref string                   `json:"image_href,omitempty"`
	Duration  int                      `json:"duration,omitempty"`
	EndTime   time.Time                `json:"end_time,omitempty"`
	CreatedAt time.Time                `json:"created_at,omitempty"`
	UpdatedAt time.Time                `json:"updated_at,omitempty"`
	Subject   map[string]SnapshotChart `json:"subject"`
}

// SnapshotChart contains the metadata for the chart requested in the Snapshot
type SnapshotChart struct {
	ID      int      `json:"id"`
	Sources []string `json:"sources"`
	Type    string   `json:"type"`
}

type SnapshotsCommunicator interface {
	Create(*Snapshot) (*Snapshot, error)
	Retrieve(int) (*Snapshot, error)
}

type SnapshotsService struct {
	client *Client
}

func NewSnapshotsService(c *Client) *SnapshotsService {
	return &SnapshotsService{c}
}

// Create requests the creation of a new Snapshot for later retrieval
func (ss *SnapshotsService) Create(s *Snapshot) (*Snapshot, error) {
	path := fmt.Sprintf("snapshots")
	req, err := ss.client.NewRequest("POST", path, s)

	if err != nil {
		return nil, err
	}

	newSnapshot := &Snapshot{}

	_, err = ss.client.Do(req, newSnapshot)
	if err != nil {
		return nil, err
	}

	return newSnapshot, nil
}

// Retrieve fetches data about a Snapshot, including a fully qualified URL for fetching the image asset itself
func (ss *SnapshotsService) Retrieve(id int) (*Snapshot, error) {
	path := fmt.Sprintf("snapshots/%d", id)
	req, err := ss.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	snapshot := &Snapshot{}
	_, err = ss.client.Do(req, snapshot)

	if err != nil {
		return nil, err
	}

	return snapshot, nil
}
