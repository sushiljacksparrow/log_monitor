package queryservice

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

func TestAppendQueries(t *testing.T) {
	tests := []struct {
		name        string
		field       string
		value       *string
		expectedLen int
	}{
		{"nil value", "service", nil, 0},
		{"empty value", "service", strPtr(""), 0},
		{"valid value", "service", strPtr("auth-service"), 1},
	}

	testCases := []struct {
		name string
		fn   func(*[]types.Query, string, *string)
	}{
		{"Term", appendTermQuery},
		{"Match", appendMatchQuery},
	}

	for _, tc := range testCases {
		for _, tt := range tests {
			tt := tt

			t.Run(tc.name+"/"+tt.name, func(t *testing.T) {
				filters := make([]types.Query, 0)

				tc.fn(&filters, tt.field, tt.value)

				if len(filters) != tt.expectedLen {
					t.Errorf("%s expected %d filters, got %d",
						tc.name, tt.expectedLen, len(filters))
				}
			})
		}
	}
}

func strPtr(s string) *string {
	return &s
}
