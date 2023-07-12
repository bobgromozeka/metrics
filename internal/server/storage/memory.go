package storage

import (
	"context"

	"github.com/bobgromozeka/metrics/internal/metrics"
)

type MemStorage struct {
	Metrics Metrics
}

type Metrics struct {
	Gauge   GaugeMetrics
	Counter CounterMetrics
}

func NewMemory() Storage {
	return &MemStorage{
		Metrics: Metrics{
			Gauge:   GaugeMetrics{},
			Counter: CounterMetrics{},
		},
	}
}

func (s *MemStorage) SetMetrics(ctx context.Context, m Metrics) {
	s.Metrics = m
}

func (s *MemStorage) GetMetricsByType(ctx context.Context, mtype string, name string) (any, error) {
	switch mtype {
	case metrics.GaugeType:
		return s.GetGaugeMetrics(ctx, name)
	case metrics.CounterType:
		return s.GetCounterMetrics(ctx, name)
	default:
		return nil, ErrWrongMetrics
	}
}

func (s *MemStorage) GetAllGaugeMetrics(ctx context.Context) (GaugeMetrics, error) {
	return s.Metrics.Gauge, nil
}

func (s *MemStorage) GetAllCounterMetrics(ctx context.Context) (CounterMetrics, error) {
	return s.Metrics.Counter, nil
}

func (s *MemStorage) GetGaugeMetrics(ctx context.Context, name string) (float64, error) {
	gm, _ := s.GetAllGaugeMetrics(ctx)
	v, ok := gm[name]
	if !ok {
		return v, ErrNotFound
	}
	return v, nil
}

func (s *MemStorage) GetCounterMetrics(ctx context.Context, name string) (int64, error) {
	cm, _ := s.GetAllCounterMetrics(ctx)
	v, ok := cm[name]
	if !ok {
		return v, ErrNotFound
	}
	return v, nil
}

func (s *MemStorage) AddCounter(ctx context.Context, name string, value int64) (int64, error) {
	if _, ok := s.Metrics.Counter[name]; !ok {
		s.Metrics.Counter[name] = 0
	}

	s.Metrics.Counter[name] += value

	return s.Metrics.Counter[name], nil
}

func (s *MemStorage) SetGauge(ctx context.Context, name string, value float64) (float64, error) {
	if _, ok := s.Metrics.Gauge[name]; !ok {
		s.Metrics.Gauge[name] = 0
	}

	s.Metrics.Gauge[name] = value

	return value, nil
}

func (s *MemStorage) UpdateMetricsByType(ctx context.Context, metricsType string, name string, value string) (any, error) {
	switch metricsType {
	case metrics.CounterType:
		parsedValue, err := metrics.ParseCounter(value)
		if err != nil {
			return false, err
		}
		return s.AddCounter(ctx, name, parsedValue)
	case metrics.GaugeType:
		parsedValue, err := metrics.ParseGauge(value)
		if err != nil {
			return false, err
		}
		return s.SetGauge(ctx, name, parsedValue)
	default:
		return false, nil
	}
}
