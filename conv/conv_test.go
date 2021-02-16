package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt2Letter(t *testing.T) {
	tests := []struct {
		in     int
		expect string
	}{
		{1, "A"},
		{25, "Y"},
		{26, "Z"},
		{27, "AA"},
		{100, "CV"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expect, Int2Letter(tt.in))
	}
}
