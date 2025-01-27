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
	MemStorage    *MemStorage
	filename      string
	isRestore     bool
	storeInterval time.Duration
}

func NewSaver(filename string, isRestore bool, duration time.Duration) (*Saver, error) {
	saver := &Saver{
		MemStorage:    NewMemStorage(),
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

	return saver, nil
}

func (s *Saver) GetCounter(name string) (string, bool) {
	return s.MemStorage.GetCounter(name)
}

func (s *Saver) SetCounter(name string, iValue int64) {
	s.MemStorage.SetCounter(name, iValue)
}

func (s *Saver) GetGauge(name string) (string, bool) {
	return s.MemStorage.GetGauge(name)
}

func (s *Saver) SetGauge(name string, fValue float64) {
	s.MemStorage.SetGauge(name, fValue)
}

func (s *Saver) AllMetrics() map[string]string {
	return s.MemStorage.AllMetrics()
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

	metrics := entities.DataJson{Data: make([]entities.MetricsJSON, 0)}
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		return err
	}

	fmt.Println("READ METRICS", metrics)
	return nil

}

func (s *Saver) saveData() error {

	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	jsonMetrics := entities.DataJson{Data: make([]entities.MetricsJSON, 0)}
	metric := entities.MetricsJSON{ID: "fdf", MType: "counter"}
	jsonMetrics.Data = append(jsonMetrics.Data, metric)
	data, err := json.Marshal(&jsonMetrics)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	if err := writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	fmt.Println("WRITE METRICS", jsonMetrics)
	return writer.Flush()
}
