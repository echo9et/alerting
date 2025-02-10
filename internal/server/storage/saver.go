package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/echo9et/alerting/internal/entities"
)

type Saver struct {
	Store         entities.Storage
	filename      string
	isRestore     bool
	storeInterval time.Duration
}

func NewSaver(storage entities.Storage, filename string, isRestore bool, duration time.Duration) (*Saver, error) {
	saver := &Saver{
		Store:         storage,
		filename:      filename,
		isRestore:     isRestore,
		storeInterval: duration,
	}

	if isRestore {
		if err := saver.restoreData(); err != nil {
			return nil, err
		}
	}

	if err := saver.saveData(); err != nil {
		panic(err)
	}

	if saver.storeInterval != 0 {
		ticker := time.NewTicker(saver.storeInterval)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					saver.saveData()
				case <-quit:
					ticker.Stop()
				}
			}
		}()
	}

	return saver, nil
}

func (s *Saver) GetCounter(name string) (string, bool) {
	return s.Store.GetCounter(name)
}

func (s *Saver) SetCounter(name string, iValue int64) {
	s.Store.SetCounter(name, iValue)
	if s.storeInterval == 0 {
		s.saveData()
	}
}

func (s *Saver) GetGauge(name string) (string, bool) {
	return s.Store.GetGauge(name)
}

func (s *Saver) SetGauge(name string, fValue float64) {
	s.Store.SetGauge(name, fValue)
	if s.storeInterval == 0 {
		s.saveData()
	}
}

func (s *Saver) AllMetrics() map[string]string {
	return s.Store.AllMetrics()
}

func (s *Saver) AllMetricsJSON() []entities.MetricsJSON {
	return s.Store.AllMetricsJSON()
}

func (s *Saver) restoreData() error {
	file, err := os.OpenFile(s.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	data, err := reader.ReadBytes('\n')

	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	metricsJSON := make([]entities.MetricsJSON, 0)
	err = json.Unmarshal(data, &metricsJSON)
	if err != nil {
		return err
	}

	for _, metric := range metricsJSON {
		switch metric.MType {
		case entities.Counter:
			s.Store.SetCounter(metric.ID, *metric.Delta)
		case entities.Gauge:
			s.Store.SetGauge(metric.ID, *metric.Value)
		default:
			fmt.Println("Не удалось прочитать тип данных при восстановление данных")
		}
	}

	// fmt.Println("READ METRICS", metricsJSON)
	return nil

}

func (s *Saver) saveData() error {

	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	metricsJSON := s.Store.AllMetricsJSON()
	data, err := json.Marshal(&metricsJSON)
	if err != nil {
		return err
	}

	if _, err := writer.Write(data); err != nil {
		return err
	}

	if err := writer.WriteByte('\n'); err != nil {
		return err
	}

	return writer.Flush()
}

func (s *Saver) Ping() bool {
	return s.Store.Ping()
}

func (s *Saver) SetMetrics(m []entities.MetricsJSON) error {
	return s.Store.SetMetrics(m)
}
