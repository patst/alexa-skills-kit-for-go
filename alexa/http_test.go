package alexa

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCertificateValidation(t *testing.T) {
	skill := Skill{
		ApplicationID:  "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
		SkipValidation: true,
	}
	skillHandler := skill.GetHTTPSkillHandler()

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

func TestInvalidBody(t *testing.T) {
	skill := Skill{
		ApplicationID:  "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
		Verbose:        true,
		SkipValidation: true,
	}
	skillHandler := skill.GetHTTPSkillHandler()

	launchRequestReader, err := os.Open("../resources/invalid_request.json")
	if err != nil {
		t.Error("Error reading input file", err)
	}

	httpRequest := httptest.NewRequest("POST", "/", launchRequestReader)
	responseWriter := httptest.NewRecorder()
	skillHandler.ServeHTTP(responseWriter, httpRequest)
	if responseWriter.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			responseWriter.Code, http.StatusBadRequest)
	}

	//try nil body
	httpRequest = httptest.NewRequest("POST", "/", nil)
	responseWriter = httptest.NewRecorder()
	skillHandler.ServeHTTP(responseWriter, httpRequest)
	if responseWriter.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			responseWriter.Code, http.StatusBadRequest)
	}
}

func TestErrorInHandler(t *testing.T) {
	skill := Skill{
		ApplicationID:  "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
		Verbose:        true,
		SkipValidation: true,
	}
	skillHandler := skill.GetHTTPSkillHandler()

	launchRequestReader, err := os.Open("../resources/invalid_request_type.json")
	if err != nil {
		t.Error("Error reading input file", err)
	}

	httpRequest := httptest.NewRequest("POST", "/", launchRequestReader)
	responseWriter := httptest.NewRecorder()
	skillHandler.ServeHTTP(responseWriter, httpRequest)
	if responseWriter.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			responseWriter.Code, http.StatusBadRequest)
	}
}
