package alexa

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// DeviceAddressService provides methods to read customer adress information data.
type DeviceAddressService interface {
	// GetCountyAndPostelCode gets the country and postal code associated with a device specified by deviceId in the syste struct.
	GetCountryAndPostalCode(system *System) (*DeviceShortAddress, error)
	// GetFullAddress gets the full address associated with the device specified by deviceId in the system struct.
	GetFullAddress(system *System) (*DeviceAddress, error)
	//IsNotAuthorizedError return true if it is a not authorized error
	IsNotAuthorizedError(err error) bool
}

// DeviceShortAddress contains the customers country and postal code.
type DeviceShortAddress struct {
	CountryCode string `json:"countryCode"`
	PostalCode  string `json:"postalCode"`
}

// DeviceAddress contains all the customers address informations.
type DeviceAddress struct {
	DeviceShortAddress
	StateOrRegion    string `json:"stateOrRegion"`
	City             string `json:"city"`
	AddressLine1     string `json:"addressLine1"`
	AddressLine2     string `json:"addressLine2"`
	AddressLine3     string `json:"addressLine3"`
	DistrictOrCounty string `json:"districtOrCounty"`
}

type deviceAddressService struct{}

var deviceAddressServiceInstance = &deviceAddressService{}

var errorForbidden = errors.New("The authentication token is invalid or doesn't have access to the resource")

func (s *deviceAddressService) executeAlexaCall(url, accessToken string, targetObj interface{}) error {
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error executing http request for URL "+url, err)
		return err
	}

	if resp.StatusCode == http.StatusForbidden {
		return errorForbidden
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected StatusCode %d, X-Amnz-RequestId=%s", resp.StatusCode, resp.Header.Get("X-Amzn-RequestId"))
	}

	defer resp.Body.Close()
	respBytes, _ := ioutil.ReadAll(resp.Body)
	log.Println("Service response: ", string(respBytes))
	err = json.Unmarshal(respBytes, targetObj)
	return err
}

func (s *deviceAddressService) GetCountryAndPostalCode(system *System) (*DeviceShortAddress, error) {
	url := fmt.Sprintf("%s/v1/devices/%s/settings/address/countryAndPostalCode", system.APIEndpoint, system.Device.DeviceID)
	var shortAddr DeviceShortAddress
	err := s.executeAlexaCall(url, system.APIAccessToken, &shortAddr)

	return &shortAddr, err
}

func (s *deviceAddressService) GetFullAddress(system *System) (*DeviceAddress, error) {
	url := fmt.Sprintf("%s/v1/devices/%s/settings/address", system.APIEndpoint, system.Device.DeviceID)
	var address DeviceAddress
	err := s.executeAlexaCall(url, system.APIAccessToken, &address)

	return &address, err
}

func (s *deviceAddressService) IsNotAuthorizedError(err error) bool {
	return err == errorForbidden
}
