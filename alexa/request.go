package alexa

import (
	"encoding/json"
	"log"
	"time"
)

// RequestEnvelope is the deserialized http post request sent by alexa.
type RequestEnvelope struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	// one of the request structs
	Request interface{} `json:"request"`
	Context Context     `json:"context"`
}

// Session object contained in standard request types like LaunchRequest, IntentRequest, SessionEndedRequest and GameEngine interface.
type Session struct {
	New         bool                   `json:"new"`
	SessionID   string                 `json:"sessionId"`
	Attributes  map[string]interface{} `json:"attributes"`
	Application Application            `json:"application"`
	User        User                   `json:"user"`
}

// Application object with the applications unique id.
type Application struct {
	ApplicationID string `json:"applicationId"`
}

// User contains the userId and access token if existent.
type User struct {
	UserID      string `json:"userId"`
	AccessToken string `json:"accessToken,omitempty"`
}

type Context struct {
	System      System      `json:"system"`
	AudioPlayer AudioPlayer `json:"audioPlayer"`
}

type System struct {
	ApiAccessToken string      `json:"apiAccessToken"`
	ApiEndpoint    string      `json:"apiEndpoint"`
	Application    Application `json:"application"`
	Device         Device      `json:"device"`
	User           User        `json:"user"`
}

type Device struct {
	DeviceId            string                 `json:"deviceId"`
	SupportedInterfaces map[string]interface{} `json:"supportedInterfaces"`
}

type AudioPlayer struct {
	Token                string `json:"token,omitempty"`
	OffsetInMilliseconds int    `json:"offsetInMilliseconds,omitempty"`
	PlayerActivity       string `json:"playerActivity"`
}

type CommonRequest struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	Locale    string `json:"locale"`
}

type LaunchRequest struct {
	CommonRequest
}

type IntentRequest struct {
	CommonRequest
	Intent      Intent `json:"intent,omitempty"`
	DialogState string `json:"dialogState,omitempty"`
}

// Intent provided in Intent requests
type Intent struct {
	Name               string          `json:"name,omitempty"`
	Slots              map[string]Slot `json:"slots,omitempty"`
	ConfirmationStatus string          `json:"confirmationStatus,omitempty"`
}

// Slot is provided in Intents
type Slot struct {
	Name               string      `json:"name"`
	Value              string      `json:"value"`
	ConfirmationStatus string      `json:"confirmationStatus,omitempty"`
	Resolutions        interface{} `json:"resolutions"`
}

type SessionEndedRequest struct {
	CommonRequest
	Reason string      `json:"reason,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

type AudioPlayerRequest struct {
	CommonRequest
	// tbd
}

// GetTypedRequest provides the request object mapped to the given struct
func (request *RequestEnvelope) GetTypedRequest(requestObj interface{}) error {
	data, _ := json.Marshal(request.Request)
	return json.Unmarshal(data, &requestObj)
}

// VerifyTimestamp checks if the the timestamp is not older than 30 seconds
func (request *RequestEnvelope) VerifyTimestamp() bool {
	timestampStr := request.Request.(map[string]interface{})["timestamp"].(string)
	requestTimestamp, err := time.Parse("2006-01-02T15:04:05Z", timestampStr)
	if err != nil {
		log.Fatalln("Error parsing request timestamp with value ", timestampStr)
	}
	if time.Since(requestTimestamp) < time.Duration(30)*time.Second {
		return true
	}
	return false
}
