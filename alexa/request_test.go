package alexa

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var skill = Skill{}

func TestLaunchRequest(t *testing.T) {
	launchRequest, _ := ioutil.ReadFile("../resources/launch_request.json")
	var r RequestEnvelope
	err := json.Unmarshal(launchRequest, &r)
	if err != nil {
		t.Fatal("Error occurred", err)
		t.FailNow()
	}
	skill.OnLaunch = func(request *LaunchRequest, response *ResponseEnvelope) {
		assert.Equal(t, "string", request.Locale, "Locale does not match")
		assert.Equal(t, "amzn1.echo-api.request.0000000-0000-0000-0000-00000000000", request.RequestID, "RequestID does not match")
		assert.Equal(t, "2015-05-13T12:34:56Z", request.Timestamp, "Timestamp does not match")
		assert.Equal(t, "LaunchRequest", request.Type, "Type does not match")
		response.Response.SetOutputSpeech("output")
	}
	response, err := r.handleRequest(&skill)

	assert.Equal(t, "<speak> output </speak>", response.Response.OutputSpeech.Ssml, "OutputSpeech does not match")
	assert.Equal(t, "SSML", response.Response.OutputSpeech.Type, "OutputSpeech Type does not match")

	if err != nil {
		t.Fatal("Error occurred", err)
		t.FailNow()
	}
}

func TestIntentRequest(t *testing.T) {
	launchRequest, _ := ioutil.ReadFile("../resources/intent_request.json")
	var r RequestEnvelope
	err := json.Unmarshal(launchRequest, &r)
	if err != nil {
		t.Fatal("Error occurred", err)
		t.FailNow()
	}
	skill.OnIntent = func(request *IntentRequest, response *ResponseEnvelope) {
		assert.Equal(t, "COMPLETED", request.DialogState, "DialogState does not match")
		assert.Equal(t, "NONE", request.Intent.ConfirmationStatus, "ConfirmationStatus does not match")
		assert.Equal(t, "GetZodiacHoroscopeIntent", request.Intent.Name, "Name does not match")
		assert.Equal(t, "ZodiacSign", request.Intent.Slots["ZodiacSign"].Name, "Name does not match")
		assert.Equal(t, "virgo", request.Intent.Slots["ZodiacSign"].Value, "Value does not match")
		assert.Equal(t, "NONE", request.Intent.Slots["ZodiacSign"].ConfirmationStatus, "ConfirmationStatus does not match")

	}
	_, err = r.handleRequest(&skill)

	if err != nil {
		t.Fatal("Error occurred", err)
		t.FailNow()
	}
}

func TestSessionEndedRequest(t *testing.T) {
	launchRequest, _ := ioutil.ReadFile("../resources/session_ended_request.json")
	var r RequestEnvelope
	err := json.Unmarshal(launchRequest, &r)
	if err != nil {
		t.Fatal("Error occurred", err)
		t.FailNow()
	}
	skill.OnSessionEnded = func(request *SessionEndedRequest, response *ResponseEnvelope) {
		assert.Equal(t, "SessionEndedRequest", request.Type, "Type does not match")
		assert.Equal(t, "USER_INITIATED", request.Reason, "Reason does not match")
	}
	_, err = r.handleRequest(&skill)

	if err != nil {
		t.Fatal("Error occurred", err)
		t.FailNow()
	}
}

func TestSessionAttributes(t *testing.T) {
	launchRequest, _ := ioutil.ReadFile("../resources/intent_request.json")
	var r RequestEnvelope
	err := json.Unmarshal(launchRequest, &r)
	if err != nil {
		t.Fatal("Error occurred", err)
		t.FailNow()
	}
	skill.OnIntent = func(request *IntentRequest, response *ResponseEnvelope) {
		assert.Equal(t, false, request.Session.New, "Session.New does not match")
		assert.Equal(t, "amzn1.echo-api.session.0000000-0000-0000-0000-00000000000", request.Session.SessionID, "Session.SessionID does not match")
		assert.Equal(t, "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe", request.Session.Application.ApplicationID, "Session.ApplicationID does not match")
		assert.Equal(t, "amzn1.account.AM3B00000000000000000000000", request.Session.User.UserID, "Session.UserID does not match")
		assert.Equal(t, "", request.Session.User.AccessToken, "")
		assert.Equal(t, true, request.Session.Attributes["supportedHoroscopePeriods"].(map[string]interface{})["daily"], "Session attribute daily does not match")
		// Add an session attribute
		request.Session.Attributes["newProp"] = "newPropValue"
	}
	response, err := r.handleRequest(&skill)

	assert.Equal(t, "newPropValue", response.SessionAttributes["newProp"], "Session attribute newProp does not match")

	if err != nil {
		t.Fatal("Error occurred", err)
		t.FailNow()
	}
}

func TestContextAttributes(t *testing.T) {
	launchRequest, _ := ioutil.ReadFile("../resources/intent_request.json")
	var r RequestEnvelope
	err := json.Unmarshal(launchRequest, &r)
	if err != nil {
		t.Fatal("Error occurred", err)
		t.FailNow()
	}
	skill.OnIntent = func(request *IntentRequest, response *ResponseEnvelope) {
		assert.Equal(t, "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe", request.Context.System.Application.ApplicationID, "")
		assert.Equal(t, "amzn1.account.AM3B00000000000000000000000", request.Context.System.User.UserID, "")
		assert.Equal(t, 0, request.Context.AudioPlayer.OffsetInMilliseconds, "")
		assert.Equal(t, "IDLE", request.Context.AudioPlayer.PlayerActivity, "")
	}
	_, err = r.handleRequest(&skill)

	if err != nil {
		t.Fatal("Error occurred", err)
		t.FailNow()
	}
}
