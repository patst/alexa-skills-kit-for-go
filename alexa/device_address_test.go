package alexa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
