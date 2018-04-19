package alexa

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestCertificateValidation(t *testing.T) {
	skill := Skill{
		ApplicationID: "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
	}
	skillHandler := skill.GetHTTPSkillHandler()

	launchRequestReader, err := os.Open("../resources/launch_request.json")
	if err != nil {
		t.Error("Error reading input file", err)
	}

	skill.SkipValidation = true
	httpRequest := httptest.NewRequest("POST", "/", launchRequestReader)
	responseWriter := httptest.NewRecorder()
	skillHandler.ServeHTTP(responseWriter, httpRequest)
	if responseWriter.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			responseWriter.Code, http.StatusOK)
	}

	// Testing without dev mode must throw an error due to a invalid certificate
	skill.SkipValidation = false
	httpRequest = httptest.NewRequest("POST", "/", launchRequestReader)
	responseWriter = httptest.NewRecorder()
	skillHandler.ServeHTTP(responseWriter, httpRequest)
	if responseWriter.Code != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			responseWriter.Code, http.StatusUnauthorized)
	}
}

func TestWrongApplicationId(t *testing.T) {
	wrongAppID := "wrong app id"

	launchRequestReader, err := os.Open("../resources/launch_request.json")
	if err != nil {
		t.Error("Error reading input file", err)
	}

	bodyBytes, _ := ioutil.ReadAll(launchRequestReader)
	var reqEnvelope RequestEnvelope
	json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&reqEnvelope)

	reqEnvelope.Request.(map[string]interface{})["timestamp"] = time.Now().Format("2006-01-02T15:04:05Z")
	err = reqEnvelope.isRequestValid(wrongAppID)
	assert.Error(t, err)
}
