package alexa

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRGBToHex(t *testing.T) {
	values := []string{"552200", "ff0000"}

	for i, element := range values {
		r, g, b := convert(element)
		hex := RgbToHex(r, g, b)
		assert.Equal(t, values[i], hex)
	}
}

func convert(input string) (int, int, int) {
	if strings.HasPrefix(input, "#") {
		input = strings.Replace(input, "#", "", 1)
	}

	if len(input) == 3 {
		input = fmt.Sprintf("%c%c%c%c%c%c", input[0], input[0], input[1], input[1], input[2], input[2])
	}

	d, _ := hex.DecodeString(input)

	return int(d[0]), int(d[1]), int(d[2])
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
