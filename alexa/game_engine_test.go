package alexa

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGameEngineInputEvent(t *testing.T) {
	skill := Skill{
		ApplicationID: "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
	}
	skillHandler := skill.GetSkillHandler()

	skill.OnGameEngineEvent = func(requestEnvelope *RequestEnvelope, request *GameEngineInputHandlerEventRequest, response *OutgoingResponse) {
		assertEqual(t, "GameEngine.InputHandlerEvent", request.Type, "")
		assertEqual(t, "amzn1.echo-api.request.463eaf71-0206-412e-b7dd-164936862994", request.OriginatingRequestID, "")
		assertEqual(t, "button_down_event", request.Events[0].Name, "")
		assertEqual(t, "down", request.Events[0].InputEvents[0].Action, "")
		assertEqual(t, "press", request.Events[0].InputEvents[0].Feature, "")
		assertEqual(t, "0000FF", request.Events[0].InputEvents[0].Color, "")
		assertEqual(t, "amzn1.ask.gadget.05RPH7PJG9C61DHI4QR0RLOQOHKUMULN8NS600CRDU8UGIM96405THTNT0283R6JJTBOND2Q9LK4MD84880C4U6J4AUHU4689FF3TTBITEACDA8V8B8E5MFRDOUM247V8GUVJKA09O1CBVSHK6LAD2J0BLV607IH03U4A13S9MS9OUO02EKIS", request.Events[0].InputEvents[0].GadgetID, "")
	}

	launchRequestReader, err := os.Open("../resources/gameengine_inputhandlerevent_request.json")
	if err != nil {
		t.Error("Error reading input file", err)
	}

	httpRequest := httptest.NewRequest("POST", "/url?dev=true", launchRequestReader)
	responseWriter := httptest.NewRecorder()
	skillHandler.ServeHTTP(responseWriter, httpRequest)
	if responseWriter.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			responseWriter.Code, http.StatusOK)
	}
}

func TestRegisterForEvents(t *testing.T) {
	skill := Skill{
		ApplicationID: "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
	}
	skillHandler := skill.GetSkillHandler()

	skill.OnLaunch = func(requestEnvelope *RequestEnvelope, request *LaunchRequest, response *OutgoingResponse) {
		response.Response.SetOutputSpeech("outputspeech")
		response.Response.SetReprompt("reprompt")
		response.Response.SimpleCard("TestCard", "TestCardText")

		directive := NewGameEngineStartInputDirective()
		directive.Timeout = 30000
		directive.Recognizers = map[string]interface{}{
			"button_down_recognizer": GameEnginePatternRecognizer{
				Type:   "match",
				Fuzzy:  false,
				Anchor: "end",
				Pattern: []GameEnginePattern{
					{
						Action: "down",
					},
				},
			},
			"button_up_recognizer": GameEnginePatternRecognizer{
				Type:   "match",
				Fuzzy:  false,
				Anchor: "end",
				Pattern: []GameEnginePattern{
					{
						Action: "up",
					},
				},
			},
		}
		directive.Events = map[string]GameEngineRegistrationEvent{
			"button_down_event": {
				Meets:                 []string{"button_down_recognizer"},
				Reports:               "matches",
				ShouldEndInputHandler: false,
			},
			"button_up_event": {
				Meets:                 []string{"button_up_recognizer"},
				Reports:               "matches",
				ShouldEndInputHandler: false,
			},
			"timeout": {
				Meets:                 []string{"timed out"},
				Reports:               "history",
				ShouldEndInputHandler: true,
			},
		}
		response.Response.Directives = append(response.Response.Directives, directive)
		response.Response.ShouldEndSession = false

	}
	launchRequestReader, err := os.Open("../resources/launch_request.json")
	if err != nil {
		t.Error("Error reading input file", err)
	}

	httpRequest := httptest.NewRequest("POST", "/url?dev=true", launchRequestReader)
	responseWriter := httptest.NewRecorder()
	skillHandler.ServeHTTP(responseWriter, httpRequest)
	if responseWriter.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			responseWriter.Code, http.StatusOK)
	}
}
