package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bobgromozeka/metrics/internal/metrics"

	"github.com/jackc/pgx/v5"
)

type Execer interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type DBStorage struct {
	*pgx.Conn
}

func NewDB(db *pgx.Conn) Storage {
	return &DBStorage{
		db,
	}
}

func (s *DBStorage) GetMetricsByType(ctx context.Context, mtype string, name string) (any, error) {
	switch mtype {
	case metrics.GaugeType:
		return s.GetGaugeMetrics(ctx, name)
	case metrics.CounterType:
		return s.GetCounterMetrics(ctx, name)
	default:
		return nil, ErrWrongMetrics
	}
}

func (s *DBStorage) GetAllGaugeMetrics(ctx context.Context) (GaugeMetrics, error) {
	gm := GaugeMetrics{}

	rows, rowsErr := s.Conn.Query(ctx, `select name, value from gauges`)
	if rowsErr != nil {
		return gm, rowsErr
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value metrics.Gauge

		scanErr := rows.Scan(&name, &value)
		if scanErr != nil {
			return gm, scanErr
		}

		gm[name] = value
	}

	if rows.Err() != nil {
		return gm, rows.Err()
	}

	return gm, nil
}

func (s *DBStorage) GetAllCounterMetrics(ctx context.Context) (CounterMetrics, error) {
	cm := CounterMetrics{}

	rows, rowsErr := s.Conn.Query(ctx, `select name, value from counters`)
	if rowsErr != nil {
		return cm, rowsErr
	}

	for rows.Next() {
		var name string
		var value metrics.Counter

		scanErr := rows.Scan(&name, &value)
		if scanErr != nil {
			return cm, scanErr
		}

		cm[name] = value
	}

	if rows.Err() != nil {
		return cm, rows.Err()
	}

	return cm, nil
}

func (s *DBStorage) GetGaugeMetrics(ctx context.Context, name string) (float64, error) {
	row := s.Conn.QueryRow(ctx, `select value from gauges where name = $1`, name)

	var val float64

	err := row.Scan(&val)
	if err != nil && err == pgx.ErrNoRows {
		return val, ErrNotFound
	}

	return val, nil
}

func (s *DBStorage) GetCounterMetrics(ctx context.Context, name string) (int64, error) {
	row := s.Conn.QueryRow(ctx, `select value from counters where name = $1`, name)

	var val int64

	err := row.Scan(&val)
	if err != nil && err == pgx.ErrNoRows {
		return val, ErrNotFound
	}

	return val, nil
}

func (s *DBStorage) AddCounter(ctx context.Context, name string, value int64) (int64, error) {
	err := addCounter(ctx, s.Conn, name, value)
	if err != nil {
		return 0, err
	}

	return s.GetCounterMetrics(ctx, name)
}

func (s *DBStorage) SetGauge(ctx context.Context, name string, value float64) (float64, error) {
	_, err := s.Conn.Exec(
		ctx,
		`insert into gauges (name, value) values($1, $2) on conflict (name) do update
		set value = $2`,
		name, value,
	)
	if err != nil {
		return 0, err
	}

	return s.GetGaugeMetrics(ctx, name)
}

func (s *DBStorage) UpdateMetricsByType(ctx context.Context, metricsType string, name string, value string) (any, error) {
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

func Bootstrap(db *pgx.Conn) error {
	ctx := context.Background()
	tx, txErr := db.Begin(ctx)
	if txErr != nil {
		return txErr
	}
	defer tx.Rollback(ctx)

	_, gaugeErr := tx.Exec(
		context.Background(),
		`create table if not exists gauges(
    			name varchar(255) not null,
    			value double precision,
    			primary key (name)
    			)`,
	)
	if gaugeErr != nil {
		return gaugeErr
	}

	_, counterErr := tx.Exec(
		context.Background(),
		`create table if not exists counters(
    			name varchar(255) not null,
    			value bigint,
    			primary key (name)
    			)`,
	)
	if counterErr != nil {
		return counterErr
	}

	return tx.Commit(ctx)
}

func (s *DBStorage) AddCounters(ctx context.Context, data CounterMetrics) error {
	tx, txErr := s.Conn.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer tx.Rollback(ctx)

	for key, val := range data {
		upsertErr := addCounter(ctx, tx, key, val)

		if upsertErr != nil {
			return upsertErr
		}
	}

	return tx.Commit(ctx)
}

func (s *DBStorage) SetGauges(ctx context.Context, data GaugeMetrics) error {
	tx, txErr := s.Conn.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer tx.Rollback(ctx)

	for key, val := range data {
		upsertErr := setGauge(ctx, tx, key, val)

		if upsertErr != nil {
			return upsertErr
		}
	}

	return tx.Commit(ctx)
}

func addCounter(ctx context.Context, conn Execer, name string, value int64) error {
	_, err := conn.Exec(
		ctx,
		`insert into counters (name, value) values($1, $2) on conflict (name) do update  
			set value = (counters.value + $2)`,
		name, value,
	)

	return err
}

func setGauge(ctx context.Context, conn Execer, name string, value float64) error {
	_, err := conn.Exec(
		ctx,
		`insert into gauges (name, value) values($1, $2) on conflict (name) do update  
			set value = $2`,
		name, value,
	)

	return err
}
