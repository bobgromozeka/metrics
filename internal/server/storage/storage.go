package storage

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("metrics not found")
var ErrWrongMetrics = errors.New("wrong metrics type")

type GaugeMetrics map[string]float64
type CounterMetrics map[string]int64

type PersistenceSettings struct {
	Path     string
	Interval uint
	Restore  bool
}

type Storage interface {
	AddCounter(context.Context, string, int64) (int64, error)
	SetGauge(context.Context, string, float64) (float64, error)
	UpdateMetricsByType(context.Context, string, string, string) (any, error)
	GetAllGaugeMetrics(context.Context) (GaugeMetrics, error)
	GetAllCounterMetrics(context.Context) (CounterMetrics, error)
	GetGaugeMetrics(context.Context, string) (float64, error)
	GetCounterMetrics(context.Context, string) (int64, error)
	GetMetricsByType(context.Context, string, string) (any, error)
	SetMetrics(context.Context, Metrics)
}
