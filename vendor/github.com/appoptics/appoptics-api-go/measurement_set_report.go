package appoptics

type MeasurementSetReport struct {
	Counts      map[string]int64
	Aggregators map[string]Aggregator
}

func NewMeasurementSetReport() *MeasurementSetReport {
	return &MeasurementSetReport{
		Counts:      map[string]int64{},
		Aggregators: map[string]Aggregator{},
	}
}
