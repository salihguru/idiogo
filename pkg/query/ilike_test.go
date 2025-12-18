package query

import (
	"testing"
)

func TestILikeMulti(t *testing.T) {
	tests := []struct {
		name       string
		fields     []string
		value      string
		wantKey    string
		wantValues []interface{}
		wantSkip   bool
	}{
		{
			name:       "Multiple fields with value",
			fields:     []string{"title", "description"},
			value:      "test",
			wantKey:    "title ILIKE ? OR description ILIKE ?",
			wantValues: []interface{}{"%test%", "%test%"},
			wantSkip:   false,
		},
		{
			name:       "Single field with value",
			fields:     []string{"title"},
			value:      "search",
			wantKey:    "title ILIKE ?",
			wantValues: []interface{}{"%search%"},
			wantSkip:   false,
		},
		{
			name:       "Empty value should skip",
			fields:     []string{"title", "description"},
			value:      "",
			wantKey:    "",
			wantValues: []interface{}{},
			wantSkip:   true,
		},
		{
			name:       "Empty fields should skip",
			fields:     []string{},
			value:      "test",
			wantKey:    "",
			wantValues: []interface{}{},
			wantSkip:   true,
		},
		{
			name:       "Three fields with value",
			fields:     []string{"title", "description", "content"},
			value:      "golang",
			wantKey:    "title ILIKE ? OR description ILIKE ? OR content ILIKE ?",
			wantValues: []interface{}{"%golang%", "%golang%", "%golang%"},
			wantSkip:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ILikeMulti(tt.fields, tt.value)

			if result.Key != tt.wantKey {
				t.Errorf("ILikeMulti() Key = %v, want %v", result.Key, tt.wantKey)
			}

			if result.Skip != tt.wantSkip {
				t.Errorf("ILikeMulti() Skip = %v, want %v", result.Skip, tt.wantSkip)
			}

			if len(result.Values) != len(tt.wantValues) {
				t.Errorf("ILikeMulti() Values length = %v, want %v", len(result.Values), len(tt.wantValues))
				return
			}

			for i, val := range result.Values {
				if val != tt.wantValues[i] {
					t.Errorf("ILikeMulti() Values[%d] = %v, want %v", i, val, tt.wantValues[i])
				}
			}
		})
	}
}

func TestILikeMultiValues(t *testing.T) {
	tests := []struct {
		name       string
		fields     []string
		values     []string
		wantKey    string
		wantValues []interface{}
		wantSkip   bool
	}{
		{
			name:       "Multiple fields with multiple values",
			fields:     []string{"title", "description"},
			values:     []string{"test", "search"},
			wantKey:    "title ILIKE ? OR title ILIKE ? OR description ILIKE ? OR description ILIKE ?",
			wantValues: []interface{}{"%test%", "%search%", "%test%", "%search%"},
			wantSkip:   false,
		},
		{
			name:       "Single field with multiple values",
			fields:     []string{"title"},
			values:     []string{"test", "search"},
			wantKey:    "title ILIKE ? OR title ILIKE ?",
			wantValues: []interface{}{"%test%", "%search%"},
			wantSkip:   false,
		},
		{
			name:       "Empty values should skip",
			fields:     []string{"title", "description"},
			values:     []string{},
			wantKey:    "",
			wantValues: []interface{}{},
			wantSkip:   true,
		},
		{
			name:       "Values with empty strings should be filtered out",
			fields:     []string{"title"},
			values:     []string{"test", "", "search"},
			wantKey:    "title ILIKE ? OR title ILIKE ?",
			wantValues: []interface{}{"%test%", "%search%"},
			wantSkip:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ILikeMultiValues(tt.fields, tt.values)

			if result.Key != tt.wantKey {
				t.Errorf("ILikeMultiValues() Key = %v, want %v", result.Key, tt.wantKey)
			}

			if result.Skip != tt.wantSkip {
				t.Errorf("ILikeMultiValues() Skip = %v, want %v", result.Skip, tt.wantSkip)
			}

			if len(result.Values) != len(tt.wantValues) {
				t.Errorf("ILikeMultiValues() Values length = %v, want %v", len(result.Values), len(tt.wantValues))
				return
			}

			for i, val := range result.Values {
				if val != tt.wantValues[i] {
					t.Errorf("ILikeMultiValues() Values[%d] = %v, want %v", i, val, tt.wantValues[i])
				}
			}
		})
	}
}
