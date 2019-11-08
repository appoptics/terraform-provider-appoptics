package appoptics

import "net/http"

type MockMeasurementsService struct {
	OnCreate func(batch *MeasurementsBatch) (*http.Response, error)
}

func (m *MockMeasurementsService) Create(batch *MeasurementsBatch) (*http.Response, error) {
	return m.OnCreate(batch)
}
