package storage

import (
	"testing"
)

func TestStorage(t *testing.T) {
	tests := []struct {
		name    string
		storage *Store
		value   int64
	}{
		{
			name:    "test set counter",
			storage: NewStore(),
			value:   5, // errors.New("float"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage.SetCounter("test", tt.value)
		})
	}
}
