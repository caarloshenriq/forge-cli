package test

import "testing"

func sum(a, b int) int {
	return a + b
}

func TestSum(t *testing.T) {
	result := sum(2, 3)
	if result != 5 {
		t.Errorf("Espect 5, but get %d", result)
	}
}