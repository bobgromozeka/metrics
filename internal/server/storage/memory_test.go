package storage

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestMemStorage_AddCounter(t *testing.T) {
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name   string
		fields Metrics
		args   args
		want   int64
	}{
		{
			name: "Adds counter to existing metrics",
			fields: Metrics{
				Gauge:   map[string]float64{},
				Counter: CounterMetrics{"a": 5},
			},
			args: args{
				name:  "a",
				value: 10,
			},
			want: 15,
		},
		{
			name: "Adds counter to non-existing metrics",
			fields: Metrics{
				Gauge:   map[string]float64{},
				Counter: CounterMetrics{"a": 5},
			},
			args: args{
				name:  "b",
				value: 10,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				s := NewMemory()
				s.AddCounters(ctx, tt.fields.Counter)
				s.SetGauges(ctx, tt.fields.Gauge)

				if got, _ := s.AddCounter(ctx, tt.args.name, tt.args.value); got != tt.want {
					t.Errorf("AddCounter() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestMemStorage_GetAllCounterMetrics(t *testing.T) {
	tests := []struct {
		name   string
		fields Metrics
		want   CounterMetrics
	}{
		{
			name: "Can get all metrics",
			fields: Metrics{
				Gauge:   GaugeMetrics{"a": 1.11},
				Counter: CounterMetrics{"b": 123},
			},
			want: CounterMetrics{"b": 123},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				s := NewMemory()
				s.AddCounters(ctx, tt.fields.Counter)
				s.SetGauges(ctx, tt.fields.Gauge)

				if got, _ := s.GetAllCounterMetrics(ctx); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetAllCounterMetrics() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestMemStorage_GetAllGaugeMetrics(t *testing.T) {
	tests := []struct {
		name   string
		fields Metrics
		want   GaugeMetrics
	}{
		{
			name: "Can get all metrics",
			fields: Metrics{
				Gauge:   GaugeMetrics{"a": 1.11},
				Counter: CounterMetrics{"b": 123},
			},
			want: GaugeMetrics{"a": 1.11},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				s := NewMemory()
				s.AddCounters(ctx, tt.fields.Counter)
				s.SetGauges(ctx, tt.fields.Gauge)

				if got, _ := s.GetAllGaugeMetrics(ctx); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetAllGaugeMetrics() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestMemStorage_GetCounterMetrics(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  Metrics
		args    args
		wantV   int64
		wantErr bool
	}{
		{
			name: "Can get metrics when exists",
			fields: Metrics{
				Gauge:   GaugeMetrics{},
				Counter: CounterMetrics{"a": 1234},
			},
			args: args{
				name: "a",
			},
			wantV:   1234,
			wantErr: false,
		},
		{
			name: "Can't get metrics when exists",
			fields: Metrics{
				Gauge:   GaugeMetrics{},
				Counter: CounterMetrics{"a": 1234},
			},
			args: args{
				name: "b",
			},
			wantV:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				s := NewMemory()
				s.AddCounters(ctx, tt.fields.Counter)
				s.SetGauges(ctx, tt.fields.Gauge)

				gotV, err := s.GetCounterMetrics(ctx, tt.args.name)
				if gotV != tt.wantV {
					t.Errorf("GetCounterMetrics() gotV = %v, want %v", gotV, tt.wantV)
				}
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

func TestMemStorage_GetGaugeMetrics(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  Metrics
		args    args
		wantV   float64
		wantErr bool
	}{
		{
			name: "Can get metrics when exists",
			fields: Metrics{
				Gauge:   GaugeMetrics{"a": 1234.123},
				Counter: CounterMetrics{},
			},
			args: args{
				name: "a",
			},
			wantV:   1234.123,
			wantErr: false,
		},
		{
			name: "Can't get metrics when exists",
			fields: Metrics{
				Gauge:   GaugeMetrics{"a": 1234.123},
				Counter: CounterMetrics{},
			},
			args: args{
				name: "b",
			},
			wantV:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				s := NewMemory()
				s.AddCounters(ctx, tt.fields.Counter)
				s.SetGauges(ctx, tt.fields.Gauge)

				gotV, err := s.GetGaugeMetrics(ctx, tt.args.name)
				if gotV != tt.wantV {
					t.Errorf("GetGaugeMetrics() gotV = %v, want %v", gotV, tt.wantV)
				}
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name   string
		fields Metrics
		args   args
		want   float64
	}{
		{
			name: "Sets gauge to existing metrics",
			fields: Metrics{
				Gauge:   GaugeMetrics{"a": 1.11},
				Counter: CounterMetrics{"a": 5},
			},
			args: args{
				name:  "a",
				value: 2.22,
			},
			want: 2.22,
		},
		{
			name: "Sets gauge to non-existing metrics",
			fields: Metrics{
				Gauge:   GaugeMetrics{},
				Counter: CounterMetrics{"a": 5},
			},
			args: args{
				name:  "b",
				value: 10.111,
			},
			want: 10.111,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				s := NewMemory()
				s.AddCounters(ctx, tt.fields.Counter)
				s.SetGauges(ctx, tt.fields.Gauge)

				if got, _ := s.SetGauge(ctx, tt.args.name, tt.args.value); got != tt.want {
					t.Errorf("SetGauge() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
