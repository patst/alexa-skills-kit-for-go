package alexa

type GameEngineStartInputDirective struct {
	Type                 string                                 `json:"type,omitempty"`
	Timeout              int                                    `json:"timeout"`
	MaximumHistoryLength int                                    `json:"maximumHistoryLength,omitempty"`
	Proxies              []interface{}                          `json:"proxies,omitempty"`
	Recognizers          interface{}                            `json:"recognizers"`
	Events               map[string]GameEngineRegistrationEvent `json:"events"`
}

// This recognizer is true when all of the specified events have occurred in the specified order.
type GameEnginePatternRecognizer struct {
	// Must be match
	Type      string              `json:"type"`
	Anchor    string              `json:"anchor,omitempty"`
	Fuzzy     bool                `json:"fuzzy"`
	GadgetIds []string            `json:"gadgetIds,omitempty"`
	Actions   []interface{}       `json:"actions,omitempty"`
	Pattern   []GameEnginePattern `json:"pattern"`
}

type GameEnginePattern struct {
	GadgetIds []string `json:"gadgetIds,omitempty"`
	Colors    []string `json:"colors,omitempty"`
	Action    string   `json:"action,omitempty"`
}

// The deviation recognizer returns true when another specified recognizer reports that the player has deviated from its expected pattern.
type GameEngineDeviationRecognizer struct {
	// Must be deviation
	Type       string `json:"type"`
	Recognizer string `json:"recognizer"`
}

// This recognizer consults another recognizer for the degree of completion, and is true if that degree is above the specified threshold. The completion parameter is specified as a decimal percentage.
type GameEngineProgressRecognizer struct {
	// Must be progress
	Type       string `json:"type"`
	Recognizer string `json:"recognizer"`
}

// The events object is where you define the conditions that must be met for your skill to be notified of Echo Button input. You must define at least one event.
type GameEngineRegistrationEvent struct {
	Meets []string `json:"meets"`
	Fails []string `json:"fails,omitempty"`
	// Possible values: history, matches
	Reports                 string `json:"reports,omitempty"`
	ShouldEndInputHandler   bool   `json:"shouldEndInputHandler"`
	MaximumInvocations      int    `json:"maximumInvocations,omitempty"`
	TriggerTimeMilliseconds int    `json:"triggerTimeMilliseconds,omitempty"`
}

// A list of events sent from the Input Handler. Each event that you specify will be sent only once to your skill as it becomes true. Note that in any InputHandlerEvent request one or more events may have become true at the same time.
type GameEngineInputEvent struct {
	Name        string `json:"name"`
	InputEvents []struct {
		GadgetId  string `json:"gadgetId"`
		Timestamp string `json:"timestamp"`
		Action    string `json:"action"`
		Color     string `json:"color"`
		Feature   string `json:"feature"`
	} `json:"inputEvents"`
}

// GameEngineInputHandlerEventRequest is send by GameEngine to notify your skill about Echo Button events
type GameEngineInputHandlerEventRequest struct {
	CommonRequest
	// From GamEngine.InputHandlerEvent
	OriginatingRequestId string                 `json:"originatingRequestId"`
	Events               []GameEngineInputEvent `json:"events"`
}

func NewGameEngineStartInputDirective() *GameEngineStartInputDirective {
	return &GameEngineStartInputDirective{
		Type: "GameEngine.StartInputHandler",
	}
}
