package alexa

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCertificateValidation(t *testing.T) {
	skill := Skill{
		ApplicationId: "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
	}
	skillHandler := skill.GetSkillHandler()

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

	// Testing without dev mode must throw an error due to a invalid certificate
	httpRequest = httptest.NewRequest("POST", "/url?dev=false", launchRequestReader)
	responseWriter = httptest.NewRecorder()
	skillHandler.ServeHTTP(responseWriter, httpRequest)
	if responseWriter.Code != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			responseWriter.Code, http.StatusUnauthorized)
	}
}

func TestWrongApplicationId(t *testing.T) {
	skill := Skill{
		ApplicationId: "wrong app id",
	}
	skillHandler := skill.GetSkillHandler()

	launchRequestReader, err := os.Open("../resources/launch_request.json")
	if err != nil {
		t.Error("Error reading input file", err)
	}

	httpRequest := httptest.NewRequest("POST", "/url?dev=true", launchRequestReader)
	responseWriter := httptest.NewRecorder()
	skillHandler.ServeHTTP(responseWriter, httpRequest)
	if responseWriter.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			responseWriter.Code, http.StatusBadRequest)
	}
}
