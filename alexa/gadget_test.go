package alexa

import (
	"encoding/json"
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

func TestAddSetLightDirective(t *testing.T) {
	var response Response
	gadgetAnimation := GadgetAnimation{
		Repeat:       1,
		TargetLights: []string{"myTarget"},
		Sequence: []GadgetAnimationStep{
			GadgetAnimationStep{
				Blend:      true,
				Color:      "red",
				DurationMs: 1000,
			},
		},
	}
	directive := response.AddGadgetControllerSetLightDirective([]string{"target"}, "none", 100, []GadgetAnimation{
		gadgetAnimation,
	})

	bytes, err := json.Marshal(directive)
	assert.NoError(t, err)
	var unmarshalledDirective GadgetControllerSetLightDirective
	json.Unmarshal(bytes, &unmarshalledDirective)

	assert.Equal(t, directive, &unmarshalledDirective)
}
