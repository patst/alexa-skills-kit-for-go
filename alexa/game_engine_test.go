package alexa

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGameEngineInputEvent(t *testing.T) {
	skill := Skill{
		ApplicationID:  "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
		SkipValidation: true,
	}
	skillHandler := skill.GetHTTPSkillHandler()

	skill.OnGameEngineEvent = func(request *GameEngineInputHandlerEventRequest, response *ResponseEnvelope) {
		assert.Equal(t, "GameEngine.InputHandlerEvent", request.Type, "")
		assert.Equal(t, "amzn1.echo-api.request.463eaf71-0206-412e-b7dd-164936862994", request.OriginatingRequestID, "")
		assert.Equal(t, "button_down_event", request.Events[0].Name, "")
		assert.Equal(t, "down", request.Events[0].InputEvents[0].Action, "")
		assert.Equal(t, "press", request.Events[0].InputEvents[0].Feature, "")
		assert.Equal(t, "0000FF", request.Events[0].InputEvents[0].Color, "")
		assert.Equal(t, "amzn1.ask.gadget.05RPH7PJG9C61DHI4QR0RLOQOHKUMULN8NS600CRDU8UGIM96405THTNT0283R6JJTBOND2Q9LK4MD84880C4U6J4AUHU4689FF3TTBITEACDA8V8B8E5MFRDOUM247V8GUVJKA09O1CBVSHK6LAD2J0BLV607IH03U4A13S9MS9OUO02EKIS", request.Events[0].InputEvents[0].GadgetID, "")
	}

	launchRequestReader, err := os.Open("../resources/gameengine_inputhandlerevent_request.json")
	if err != nil {
		t.Error("Error reading input file", err)
	}

	httpRequest := httptest.NewRequest("POST", "/", launchRequestReader)
	responseWriter := httptest.NewRecorder()
	skillHandler.ServeHTTP(responseWriter, httpRequest)
	if responseWriter.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			responseWriter.Code, http.StatusOK)
	}
}

func TestRegisterForEvents(t *testing.T) {
	skill := Skill{
		ApplicationID:  "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
		SkipValidation: true,
	}
	skillHandler := skill.GetHTTPSkillHandler()

	skill.OnLaunch = func(request *LaunchRequest, response *ResponseEnvelope) {
		response.Response.SetOutputSpeech("outputspeech")
		response.Response.SetReprompt("reprompt")
		response.Response.SetSimpleCard("TestCard", "TestCardText")

		directive := NewGameEngineStartInputDirective(30000)
		buttonDownRecognizer := directive.AddPatternRecognizer("button_down_recognizer")
		buttonDownRecognizer.Fuzzy = false
		buttonDownRecognizer.Anchor = "end"
		buttonDownRecognizer.AddPattern(nil, nil, "down")

		buttonUpRecognizer := directive.AddPatternRecognizer("button_up_recognizer")
		buttonUpRecognizer.Fuzzy = false
		buttonUpRecognizer.Anchor = "end"
		buttonUpRecognizer.AddPattern(nil, nil, "up")

		buttonDownEvent := directive.AddEvent("button_down_event", false, []string{"button_down_recognizer"})
		buttonDownEvent.Reports = "matches"

		buttonUpEvent := directive.AddEvent("button_up_event", false, []string{"button_up_recognizer"})
		buttonUpEvent.Reports = "matches"

		timeoutEvent := directive.AddEvent("timeout", true, []string{"timed out"})
		timeoutEvent.Reports = "history"

		response.Response.AddDirective(directive)
		response.Response.ShouldEndSession = false

	}
	launchRequestReader, err := os.Open("../resources/launch_request.json")
	if err != nil {
		t.Error("Error reading input file", err)
	}

	httpRequest := httptest.NewRequest("POST", "/", launchRequestReader)
	responseWriter := httptest.NewRecorder()
	skillHandler.ServeHTTP(responseWriter, httpRequest)
	if responseWriter.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			responseWriter.Code, http.StatusOK)
	}
}
