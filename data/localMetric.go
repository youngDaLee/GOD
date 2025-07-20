package data

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type JsonMetric struct {
	mu      *sync.RWMutex
	metrics map[string][]DataPoint
	path    string
}

func NewJsonMetric(path string) *JsonMetric {
	return &JsonMetric{
		metrics: make(map[string][]DataPoint),
		path:    path,
	}
}

func (m *JsonMetric) Record(metricName string, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.metrics[metricName]; !ok {
		m.metrics[metricName] = []DataPoint{}
	}

	m.metrics[metricName] = append(m.metrics[metricName], DataPoint{
		Timestamp: time.Now(),
		Value:     value,
	})
}

func (m *JsonMetric) Flush() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	file, err := os.Create(m.path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(m.metrics)
}
