package data

import "time"

type Metic interface {
	// Record 매트릭 수집 처리
	Record(metricName string, value float64)
	// Flush 매트릭 저장 처리
	Flush() error
}

type DataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     string    `json:"value"`
}
