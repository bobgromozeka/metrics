package metrics

import (
	"strconv"
)

type Gauge = float64
type Counter = int64

type RequestPayload struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

var validNames = map[string]struct{}{
	GaugeType:   {},
	CounterType: {},
}

func ParseCounter(value string) (Counter, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func ParseGauge(value string) (Gauge, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func IsValidValue(metricsType string, value string) bool {
	isValid := false

	switch metricsType {
	case CounterType:
		_, err := ParseCounter(value)
		isValid = err == nil
	case GaugeType:
		_, err := ParseGauge(value)
		isValid = err == nil
	}

	return isValid
}

func IsValidType(metricsType string) bool {
	_, ok := validNames[metricsType]
	return ok
}
