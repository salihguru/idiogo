package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JsonbArray[T any] []*T

func (j *JsonbArray[T]) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &j)
	case string:
		if v != "" {
			return j.Scan([]byte(v))
		}
	default:
		return fmt.Errorf("invalid type: %T", v)
	}
	return nil
}

func (j JsonbArray[T]) Value() (driver.Value, error) {
	data, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func JsonbObjScan(value interface{}, obj interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, obj)
	case string:
		if v != "" {
			return JsonbObjScan([]byte(v), obj)
		}
	default:
		return errors.New("not supported")
	}
	return nil
}

func JsonbObjValue(obj interface{}) (driver.Value, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

// JsonbMap represents a JSONB map for PostgreSQL
type JsonbMap map[string]interface{}

func (j *JsonbMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JsonbMap)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		if len(v) == 0 {
			*j = make(JsonbMap)
			return nil
		}
		return json.Unmarshal(v, j)
	case string:
		if v == "" {
			*j = make(JsonbMap)
			return nil
		}
		return j.Scan([]byte(v))
	default:
		return fmt.Errorf("invalid type: %T", v)
	}
}

func (j JsonbMap) Value() (driver.Value, error) {
	if j == nil || len(j) == 0 {
		return "{}", nil
	}

	data, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

// JsonbMapFloat64 represents a JSONB map with string keys and float64 values for PostgreSQL
type JsonbMapFloat64 map[string]float64

func (j *JsonbMapFloat64) Scan(value interface{}) error {
	if value == nil {
		*j = make(JsonbMapFloat64)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		if len(v) == 0 {
			*j = make(JsonbMapFloat64)
			return nil
		}
		return json.Unmarshal(v, j)
	case string:
		if v == "" {
			*j = make(JsonbMapFloat64)
			return nil
		}
		return j.Scan([]byte(v))
	default:
		return fmt.Errorf("invalid type: %T", v)
	}
}

func (j JsonbMapFloat64) Value() (driver.Value, error) {
	if j == nil || len(j) == 0 {
		return "{}", nil
	}

	data, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}
