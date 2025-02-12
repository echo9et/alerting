package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/echo9et/alerting/internal/entities"
	"github.com/go-chi/chi/v5"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type UnknowType struct {
	Message string
}

func (e *UnknowType) Error() string {
	return "Unknow Type"
}

var supportMetrics = map[string]func(entities.ManagerValues, string, string) error{
	Gauge:   handlerGauge,
	Counter: handlerCounters,
}

func handlerCounters(s entities.ManagerValues, name, sValue string) error {
	iValue, err := strconv.ParseInt(sValue, 10, 64)
	if err != nil {
		return err
	}
	s.SetCounter(name, iValue)
	return nil
}

func handlerGauge(s entities.ManagerValues, name, sValue string) error {
	fValue, err := strconv.ParseFloat(sValue, 64)
	if err != nil {
		return err
	}
	s.SetGauge(name, fValue)
	return nil
}

func WriteMetric(w http.ResponseWriter, r *http.Request, s entities.Storage) error {
	handlerMetric, ok := supportMetrics[chi.URLParam(r, "type")]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	name, value := chi.URLParam(r, "name"), chi.URLParam(r, "value")
	err := handlerMetric(s, name, value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func WriteMetricJSON(w http.ResponseWriter, r *http.Request, s entities.Storage) error {
	var mj entities.MetricsJSON
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf.Bytes(), &mj); err != nil {
		return err
	}

	if err = saveMetricsJSON(s, mj); err != nil {
		return err
	}

	out, err := json.Marshal(mj)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
	return nil
}

func WriteMetricsJSON(w http.ResponseWriter, r *http.Request, s entities.Storage) error {
	var metricsJSON = make([]entities.MetricsJSON, 0)
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		fmt.Println(time.Now(), "ERROR: READFROM", err)
		return err
	}

	if err = json.Unmarshal(buf.Bytes(), &metricsJSON); err != nil {
		fmt.Println(time.Now(), "ERROR: UMARHAL", err)
		return err
	}

	return s.SetMetrics(metricsJSON)
}

func ReadMetricJSON(w http.ResponseWriter, r *http.Request, s entities.Storage) error {
	var mj entities.MetricsJSON
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf.Bytes(), &mj); err != nil {
		return err
	}
	if err = getMetricsJSON(s, &mj); err != nil {
		return err
	}

	out, err := json.Marshal(mj)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)

	return nil
}

func getMetricsJSON(s entities.Storage, mj *entities.MetricsJSON) error {
	switch mj.MType {
	case Counter:
		value, status := s.GetCounter(mj.ID)
		if !status {
			return errors.New("counter not found")
		}
		iValue, _ := strconv.ParseInt(value, 10, 64)
		mj.Delta = &iValue
	case Gauge:
		value, status := s.GetGauge(mj.ID)
		if !status {
			return errors.New("gauge not found")
		}
		dValue, _ := strconv.ParseFloat(value, 64)
		mj.Value = &dValue
	}

	return nil
}

func saveMetricsJSON(s entities.Storage, mj entities.MetricsJSON) error {
	switch mj.MType {
	case Counter:
		s.SetCounter(mj.ID, *mj.Delta)
	case Gauge:
		s.SetGauge(mj.ID, *mj.Value)
	default:
		return &UnknowType{}
	}

	return nil
}
