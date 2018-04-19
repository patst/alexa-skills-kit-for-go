package alexa

import (
	"testing"

	hex "github.com/dlion/hex2rgb"
	"github.com/stretchr/testify/assert"
)

func TestRGBToHex(t *testing.T) {
	values := []string{"552200", "ff0000"}

	for i, element := range values {
		r, g, b := hex.Convert(element)
		hex := RgbToHex(r, g, b)
		assert.Equal(t, values[i], hex)
	}
}
