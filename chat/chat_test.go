package chat

import "testing"

var FormatValueTests = []struct {
	Value float64
	S     string
}{
	{100.00, "100"},
	{100.10, "100.10"},
	{100.01, "100.01"},
	{100.001, "100"},
}

func TestFormatValue(t *testing.T) {
	for _, tt := range FormatValueTests {
		if s := FormatValue(tt.Value); s != tt.S {
			t.Errorf("%v: want %q, got %q", tt.Value, tt.S, s)
		}
	}
}
