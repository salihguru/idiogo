package entity

import (
	"testing"
)

func TestJsonbMap_Value(t *testing.T) {
	tests := []struct {
		name     string
		input    JsonbMap
		expected string
	}{
		{
			name:     "nil map",
			input:    nil,
			expected: "{}",
		},
		{
			name:     "empty map", 
			input:    make(JsonbMap),
			expected: "{}",
		},
		{
			name:     "map with data",
			input:    JsonbMap{"key": "value", "number": 42},
			expected: `{"key":"value","number":42}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.input.Value()
			if err != nil {
				t.Errorf("JsonbMap.Value() error = %v", err)
				return
			}
			
			result := value.(string)
			if tt.name == "map with data" {
				// For map with data, just check it contains both keys
				if !contains(result, "key") || !contains(result, "number") {
					t.Errorf("JsonbMap.Value() = %v, want to contain both keys", result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("JsonbMap.Value() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestJsonbMap_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected JsonbMap
		wantErr  bool
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: make(JsonbMap),
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: make(JsonbMap),
			wantErr:  false,
		},
		{
			name:     "empty bytes",
			input:    []byte{},
			expected: make(JsonbMap),
			wantErr:  false,
		},
		{
			name:     "valid json string",
			input:    `{"key":"value","number":42}`,
			expected: JsonbMap{"key": "value", "number": float64(42)}, // JSON numbers are float64
			wantErr:  false,
		},
		{
			name:     "valid json bytes",
			input:    []byte(`{"key":"value"}`),
			expected: JsonbMap{"key": "value"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result JsonbMap
			err := result.Scan(tt.input)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonbMap.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if len(result) != len(tt.expected) {
					t.Errorf("JsonbMap.Scan() length = %v, want %v", len(result), len(tt.expected))
					return
				}
				
				for k, v := range tt.expected {
					if result[k] != v {
						t.Errorf("JsonbMap.Scan() key %v = %v, want %v", k, result[k], v)
					}
				}
			}
		})
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && func() bool {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	}()
}