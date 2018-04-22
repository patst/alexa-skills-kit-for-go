package alexa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDialogDirectives(t *testing.T) {
	skill := Skill{
		ApplicationID:  "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
		SkipValidation: true,
	}
	skillHandler := skill.GetHTTPSkillHandler()

	skill.OnLaunch = func(request *LaunchRequest, responseEnvelope *ResponseEnvelope) {
		responseEnvelope.Response.AddDialogConfirmIntentDirective()
		responseEnvelope.Response.AddDialogConfirmSlotDirective("slot1")
		responseEnvelope.Response.AddDialogDelegateDirective()
		responseEnvelope.Response.AddDialogElicitSlotDirective("slot2")
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
	respBytes, _ := ioutil.ReadAll(responseWriter.Body)
	var resp map[string]interface{}
	json.Unmarshal(respBytes, &resp)
	fmt.Println("response: ", string(respBytes))
	directives := resp["response"].(map[string]interface{})["directives"].([]interface{})
	assert.Equal(t, 4, len(directives))
	assert.Equal(t, "Dialog.ConfirmIntent", directives[0].(map[string]interface{})["type"])
	assert.Equal(t, "Dialog.ConfirmSlot", directives[1].(map[string]interface{})["type"])
	assert.Equal(t, "slot1", directives[1].(map[string]interface{})["slotToConfirm"])
	assert.Equal(t, "Dialog.Delegate", directives[2].(map[string]interface{})["type"])
	assert.Equal(t, "Dialog.ElicitSlot", directives[3].(map[string]interface{})["type"])
	assert.Equal(t, "slot2", directives[3].(map[string]interface{})["slotToElicit"])
}
