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

func TestDeviceAddressServiceFullAddress(t *testing.T) {
	var tokenOk, tokenNotOk string = "tokenOk", "tokenNotOk"
	// Start mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "Bearer "+tokenOk {
			var deviceAddress DeviceAddress
			deviceAddress.CountryCode = authHeader[len("Bearer "):]
			bytes, marshallErr := json.Marshal(deviceAddress)

			if marshallErr != nil {
				fmt.Println("Error marshalling ", marshallErr)
				t.Fail()
			}
			w.Write(bytes)
		} else if authHeader == "Bearer "+tokenNotOk {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
	}))
	defer ts.Close()

	var deviceAddressService = GetDeviceAddressService()

	system := &System{
		APIAccessToken: tokenOk,
		APIEndpoint:    ts.URL,
		Device: Device{
			DeviceID: "deviceID",
		},
	}

	//Test okay case
	addr, err := deviceAddressService.GetFullAddress(system)
	assert.NoError(t, err)
	assert.NotNil(t, addr)
	assert.Equal(t, tokenOk, addr.CountryCode)

	// Invalid access token
	system.APIAccessToken = tokenNotOk
	addr, err = deviceAddressService.GetFullAddress(system)
	assert.Error(t, err)
	assert.True(t, deviceAddressService.IsNotAuthorizedError(err))

	// Some other error
	system.APIAccessToken = "random token"
	addr, err = deviceAddressService.GetFullAddress(system)
	assert.Error(t, err)
	assert.False(t, deviceAddressService.IsNotAuthorizedError(err))

	//Wrong url
	system.APIAccessToken = tokenOk
	system.APIEndpoint = "http://wrong"
	addr, err = deviceAddressService.GetFullAddress(system)
	assert.Error(t, err)
	assert.False(t, deviceAddressService.IsNotAuthorizedError(err))

}

func TestDeviceAddressServiceCountryAndPostalCode(t *testing.T) {
	var tokenOk, tokenNotOk string = "tokenOk", "tokenNotOk"
	// Start mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "Bearer "+tokenOk {
			var deviceAddress DeviceAddress
			deviceAddress.CountryCode = authHeader[len("Bearer "):]
			bytes, marshallErr := json.Marshal(deviceAddress)

			if marshallErr != nil {
				fmt.Println("Error marshalling ", marshallErr)
				t.Fail()
			}
			w.Write(bytes)
		} else if authHeader == "Bearer "+tokenNotOk {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
	}))
	defer ts.Close()

	var deviceAddressService = GetDeviceAddressService()

	system := &System{
		APIAccessToken: tokenOk,
		APIEndpoint:    ts.URL,
		Device: Device{
			DeviceID: "deviceID",
		},
	}

	//Test okay case
	addr, err := deviceAddressService.GetCountryAndPostalCode(system)
	assert.NoError(t, err)
	assert.NotNil(t, addr)
	assert.Equal(t, tokenOk, addr.CountryCode)

	// Invalid access token
	system.APIAccessToken = tokenNotOk
	addr, err = deviceAddressService.GetCountryAndPostalCode(system)
	assert.Error(t, err)
	assert.True(t, deviceAddressService.IsNotAuthorizedError(err))

	// Some other error
	system.APIAccessToken = "random token"
	addr, err = deviceAddressService.GetCountryAndPostalCode(system)
	assert.Error(t, err)
	assert.False(t, deviceAddressService.IsNotAuthorizedError(err))

	//Wrong url
	system.APIAccessToken = tokenOk
	system.APIEndpoint = "http://wrong"
	addr, err = deviceAddressService.GetCountryAndPostalCode(system)
	assert.Error(t, err)
	assert.False(t, deviceAddressService.IsNotAuthorizedError(err))

}

func TestPermissionsConsentCard(t *testing.T) {
	skill := Skill{
		ApplicationID:  "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
		SkipValidation: true,
	}
	skillHandler := skill.GetHTTPSkillHandler()

	skill.OnLaunch = func(request *LaunchRequest, response *ResponseEnvelope) {
		response.Response.SetAskForPermissionsConsentCard("Consent", "Please allow address info", []string{"read::alexa:device:all:address"})
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
	card := resp["response"].(map[string]interface{})["card"].(map[string]interface{})
	assert.Equal(t, "AskForPermissionsConsent", card["type"])
}
