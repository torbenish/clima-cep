package temperature

import "testing"

func TestConversions(t *testing.T) {
	tests := []struct {
		name string
		c, f, k float64
	}{
		{"exemplo", 28.5, 83.3, 301.5},
		{"zero", 0, 32, 273},
		{"negativo", -10, 14, 263},
		{"ebulicao", 100, 212, 373},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CelsiusToFahrenheit(tt.c); got != tt.f {
				t.Errorf("F: got %v, want %v", got, tt.f)
			}
			if got := CelsiusToKelvin(tt.c); got != tt.k {
				t.Errorf("K: got %v, want %v", got, tt.k)
			}
		})
	}
}