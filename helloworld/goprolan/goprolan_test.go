package goprolan

import "testing"

func Test_GoProLan(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"case1"},
		{"case2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GoProLan()
		})
	}
}
