package vtb24

import (
	"encoding/json"
	"testing"
	"time"
)

func TestRateValue_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		JSON  string
		Value float64
	}{
		{`"0"`, 0},
		{`"12"`, 12},
		{`"12,34"`, 12.34},
		{`"12.34"`, 12.34},
	}

	for _, tt := range tests {
		t.Run(tt.JSON, func(t *testing.T) {
			var v RateValue
			if err := json.Unmarshal([]byte(tt.JSON), &v); err != nil {
				t.Fatal(err)
			}
			if v2 := RateValue(tt.Value); v != v2 {
				t.Errorf("want %s, got %s", v2, v)
			}
		})
	}
}

func TestRateTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		JSON string
		Time time.Time
	}{
		{
			`"/Date(1506453186593)/"`,
			time.Date(2017, time.September, 26, 19, 13, 6, 593, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.JSON, func(t *testing.T) {
			var v RateTime
			if err := json.Unmarshal([]byte(tt.JSON), &v); err != nil {
				t.Fatal(err)
			}
			if v := time.Time(v); tt.Time != v {
				t.Errorf("want %s, got %s", tt.Time, v)
			}
		})
	}
}
