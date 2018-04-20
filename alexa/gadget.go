package alexa

import (
	"fmt"
	"strconv"
)

// GadgetControllerSetLightDirective sends a command to animate the LEDs of connected Echo Buttons. The following example shows the general form of the directive.
type GadgetControllerSetLightDirective struct {
	Type          string           `json:"type,omitempty"`
	Version       int              `json:"version"`
	TargetGadgets []string         `json:"targetGadgets,omitempty"`
	Parameters    GadgetParameters `json:"parameters"`
}

// GadgetTriggerEventType describes the action that triggers the animation.
type GadgetTriggerEventType string

const (
	// BUTTON_DOWN GadgetTriggerEventType value for button down event
	buttonDown = "buttonDown"
	// BUTTON_UP GadgetTriggerEventType value for button up event
	buttonUp = "buttonUp"
	none     = "none"
)

// GadgetParameters contains instructions on how to animate the buttons.
type GadgetParameters struct {
	TriggerEvent       GadgetTriggerEventType `json:"triggerEvent"`
	TriggerEventTimeMs int                    `json:"triggerEventTimeMs"`
	Animations         []GadgetAnimation      `json:"animations"`
}

// GadgetAnimation contains a sequence of instructions to be performed in a specific order, along with the number of times to play the overall animation.
type GadgetAnimation struct {
	Repeat       int                   `json:"repeat"`
	TargetLights []string              `json:"targetLights"`
	Sequence     []GadgetAnimationStep `json:"sequence"`
}

// GadgetAnimationStep a step to render in a animation.
type GadgetAnimationStep struct {
	DurationMs int    `json:"durationMs"`
	Color      string `json:"color"`
	Blend      bool   `json:"blend"`
}

// RgbToHex converts single rgb values to a hex string representation.
func RgbToHex(r, g, b int) string {
	return fmt.Sprintf("%02s%02s%02s", strconv.FormatInt(int64(r), 16), strconv.FormatInt(int64(g), 16), strconv.FormatInt(int64(b), 16))
}
