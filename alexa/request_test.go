package alexa

import (
	"encoding/json"
	"fmt"
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
		assertEqual(t, "string", request.Locale, "Locale does not match")
		assertEqual(t, "amzn1.echo-api.request.0000000-0000-0000-0000-00000000000", request.RequestID, "RequestID does not match")
		assertEqual(t, "2015-05-13T12:34:56Z", request.Timestamp, "Timestamp does not match")
		assertEqual(t, "LaunchRequest", request.Type, "Type does not match")
		response.Response.SetOutputSpeech("output")
	}
	response, err := handleRequest(&r, &skill)

	assertEqual(t, "output", response.Response.OutputSpeech.Text, "OutputSpeech does not match")
	assertEqual(t, "PlainText", response.Response.OutputSpeech.Type, "OutputSpeech Type does not match")

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
		assertEqual(t, "COMPLETED", request.DialogState, "DialogState does not match")
		assertEqual(t, "NONE", request.Intent.ConfirmationStatus, "ConfirmationStatus does not match")
		assertEqual(t, "GetZodiacHoroscopeIntent", request.Intent.Name, "Name does not match")
		assertEqual(t, "ZodiacSign", request.Intent.Slots["ZodiacSign"].Name, "Name does not match")
		assertEqual(t, "virgo", request.Intent.Slots["ZodiacSign"].Value, "Value does not match")
		assertEqual(t, "NONE", request.Intent.Slots["ZodiacSign"].ConfirmationStatus, "ConfirmationStatus does not match")

	}
	_, err = handleRequest(&r, &skill)

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
		assertEqual(t, "SessionEndedRequest", request.Type, "Type does not match")
		assertEqual(t, "USER_INITIATED", request.Reason, "Reason does not match")
	}
	_, err = handleRequest(&r, &skill)

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
		assertEqual(t, false, request.Session.New, "Session.New does not match")
		assertEqual(t, "amzn1.echo-api.session.0000000-0000-0000-0000-00000000000", request.Session.SessionID, "Session.SessionID does not match")
		assertEqual(t, "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe", request.Session.Application.ApplicationID, "Session.ApplicationID does not match")
		assertEqual(t, "amzn1.account.AM3B00000000000000000000000", request.Session.User.UserID, "Session.UserID does not match")
		assertEqual(t, "", request.Session.User.AccessToken, "")
		assertEqual(t, true, request.Session.Attributes["supportedHoroscopePeriods"].(map[string]interface{})["daily"], "Session attribute daily does not match")
		// Add an session attribute
		request.Session.Attributes["newProp"] = "newPropValue"
	}
	response, err := handleRequest(&r, &skill)

	assertEqual(t, "newPropValue", response.SessionAttributes["newProp"], "Session attribute newProp does not match")

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
		assertEqual(t, "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe", request.Context.System.Application.ApplicationID, "")
		assertEqual(t, "amzn1.account.AM3B00000000000000000000000", request.Context.System.User.UserID, "")
		assertEqual(t, 0, request.Context.AudioPlayer.OffsetInMilliseconds, "")
		assertEqual(t, "IDLE", request.Context.AudioPlayer.PlayerActivity, "")
	}
	_, err = handleRequest(&r, &skill)

	if err != nil {
		t.Fatal("Error occurred", err)
		t.FailNow()
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}
