package alexa

import ()

type GadgetControllerSetLightDirective struct {
	Type          string           `json:"type,omitempty"`
	Version       int              `json:"version"`
	TargetGadgets []string         `json:"targetGadgets,omitempty"`
	Parameters    GadgetParameters `json:"parameters"`
}
type GadgetTriggerEventType string

const (
	BUTTON_DOWN = "buttonDown"
	BUTTON_UP   = "buttonUp"
	NONE        = "none"
)

type GadgetParameters struct {
	TriggerEvent       GadgetTriggerEventType `json:"triggerEvent"`
	TriggerEventTimeMs int                    `json:"triggerEventTimeMs"`
	Animations         []GadgetAnimation      `json:"animations"`
}

type GadgetAnimation struct {
	Repeat       int                   `json:"repeat"`
	TargetLights []int                 `json:"targetLights"`
	Sequence     []GadgetAnimationStep `json:"sequence"`
}

type GadgetAnimationStep struct {
	DurationMs int    `json:"durationMs"`
	Color      string `json:"color"`
	Blend      bool   `json:"blend"`
}

type SystemExceptionEncounteredRequest struct {
	CommonRequest
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
	Cause struct {
		RequestId string `json:"requestId"`
	} `json:"cause"`
}
