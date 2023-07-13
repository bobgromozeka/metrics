package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type PersistentStorage struct {
	Storage
	persistencePath string
	syncPersisting  bool
}

func NewPersistenceStorage(s Storage, config PersistenceSettings) Storage {
	ps := &PersistentStorage{
		Storage:         s,
		persistencePath: config.Path,
	}

	if config.Path != "" {
		if config.Restore {
			data, restoreErr := restoreFrom(config.Path)
			if restoreErr != nil {
				log.Println("Could not restore data from file: ", restoreErr)
			} else {
				if cErr := ps.AddCounters(context.Background(), data.Counter); cErr != nil {
					log.Println("Could not populate counters with restored data:", cErr)
				}

				if gErr := ps.SetGauges(context.Background(), data.Gauge); gErr != nil {
					log.Println("Could not populate gauges with restored data:", gErr)
				}
			}
		}

		if config.Interval == 0 {
			ps.syncPersisting = true
		} else {
			go func() {
				ticker := time.Tick(time.Second * time.Duration(config.Interval))
				for range ticker {
					ps.persist()
				}
			}()
		}
	}

	return ps
}

func (s *PersistentStorage) SetGauge(ctx context.Context, name string, value float64) (float64, error) {
	res, err := s.Storage.SetGauge(ctx, name, value)

	if s.syncPersisting {
		s.persist()
	}

	return res, err
}

func (s *PersistentStorage) AddCounter(ctx context.Context, name string, value int64) (int64, error) {
	res, err := s.Storage.AddCounter(ctx, name, value)

	if s.syncPersisting {
		s.persist()
	}

	return res, err
}

func (s *PersistentStorage) persist() {
	ctx := context.Background()
	gm, gErr := s.GetAllGaugeMetrics(ctx)
	cm, cErr := s.GetAllCounterMetrics(ctx)
	if gErr != nil || cErr != nil {
		log.Println("Could not get one of metrics with error: ", cErr, gErr)
		return
	}

	err := persistToPath(
		s.persistencePath, Metrics{
			Gauge:   gm,
			Counter: cm,
		},
	)
	if err != nil {
		fmt.Println("Error during syncing storage data: ", err)
	}
}

func persistToPath(path string, data Metrics) error {
	jsonData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		return jsonErr
	}

	if writeErr := os.WriteFile(path, jsonData, 0666); writeErr != nil {
		return writeErr
	}

	return nil
}

func restoreFrom(filepath string) (Metrics, error) {
	data := Metrics{}

	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		return data, err
	}

	unmarshalErr := json.Unmarshal(jsonData, &data)
	if unmarshalErr != nil {
		return data, unmarshalErr
	}

	return data, nil
}
