package alexa

import (
	"errors"
	"strings"
)

// Skill configures the different Handlers for skill execution.
type Skill struct {
	ApplicationID string
	// SkipValidation skips any request validation (TEST ONLY!)
	SkipValidation bool
	// Verbose enables request and response logging
	Verbose                  bool
	OnLaunch                 func(*LaunchRequest, *ResponseEnvelope)
	OnIntent                 func(*IntentRequest, *ResponseEnvelope)
	OnSessionEnded           func(*SessionEndedRequest, *ResponseEnvelope)
	OnAudioPlayerState       func(*AudioPlayerRequest, *ResponseEnvelope)
	OnAudioPlayerFailedState func(*AudioPlayerPlaybackFailedRequest, *ResponseEnvelope)
	OnSystemException        func(*SystemExceptionEncounteredRequest, *ResponseEnvelope)
	OnGameEngineEvent        func(*GameEngineInputHandlerEventRequest, *ResponseEnvelope)
}

// GetDeviceAddressService provides an instance of the device address service to query a customers address information.
func GetDeviceAddressService() DeviceAddressService {
	return deviceAddressServiceInstance
}

func (requestEnvelope *RequestEnvelope) handleRequest(skill *Skill) (*ResponseEnvelope, error) {
	//Read the type for this request to do the correct routing
	var commonRequest CommonRequest
	err := requestEnvelope.getTypedRequest(&commonRequest)
	if err != nil {
		return nil, err
	}

	requestType := commonRequest.Type

	// Create response and map the session attributes from the request
	response := newResponseEnvelope(requestEnvelope.Session.Attributes)

	// Request handling
	if requestType == "LaunchRequest" {

		if skill.OnLaunch != nil {
			//Map to the correct type
			var request LaunchRequest
			// Create concrete types
			requestEnvelope.getTypedRequest(&request)
			skill.OnLaunch(&request, response)
		}
	} else if requestType == "IntentRequest" {
		if skill.OnIntent != nil {
			var request IntentRequest
			// Create concrete types
			requestEnvelope.getTypedRequest(&request)
			skill.OnIntent(&request, response)
		}
	} else if requestType == "SessionEndedRequest" {
		if skill.OnSessionEnded != nil {
			var request SessionEndedRequest
			// Create concrete types
			requestEnvelope.getTypedRequest(&request)
			skill.OnSessionEnded(&request, response)
		}
	} else if strings.HasPrefix(requestType, "AudioPlayer.") {
		if skill.OnAudioPlayerState != nil {
			// Create concrete types
			if requestType == "AudioPlayer.PlaybackFailed" {
				var request AudioPlayerPlaybackFailedRequest
				requestEnvelope.getTypedRequest(&request)
				skill.OnAudioPlayerFailedState(&request, response)
			} else {
				var request AudioPlayerRequest
				requestEnvelope.getTypedRequest(&request)
				skill.OnAudioPlayerState(&request, response)
			}
		}
	} else if strings.HasPrefix(requestType, "GameEngine.") {
		if skill.OnGameEngineEvent != nil {
			var request GameEngineInputHandlerEventRequest
			// Create concrete types
			requestEnvelope.getTypedRequest(&request)
			skill.OnGameEngineEvent(&request, response)
		}
	} else if requestType == "System.ExceptionEncountered" {
		if skill.OnSystemException != nil {
			var request SystemExceptionEncounteredRequest
			// Create concrete types
			requestEnvelope.getTypedRequest(&request)
			skill.OnSystemException(&request, response)
		}
	} else {
		return nil, errors.New("Invalid request type: " + requestType)
	}
	return response, nil
}
