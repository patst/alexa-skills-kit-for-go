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

// Context object provides your skill with information about the current state of the Alexa service and device at the time the request is sent to your service.
type Context struct {
	System      System      `json:"system"`
	AudioPlayer AudioPlayer `json:"audioPlayer"`
}

// System object that provides information about the current state of the Alexa service and the device interacting with your skill.
type System struct {
	APIAccessToken string      `json:"apiAccessToken"`
	APIEndpoint    string      `json:"apiEndpoint"`
	Application    Application `json:"application"`
	Device         Device      `json:"device"`
	User           User        `json:"user"`
}

// Device object providing information about the device used to send the request.
type Device struct {
	DeviceID            string                 `json:"deviceId"`
	SupportedInterfaces map[string]interface{} `json:"supportedInterfaces"`
}

// AudioPlayer object providing the current state for the AudioPlayer interface.
type AudioPlayer struct {
	Token                string `json:"token,omitempty"`
	OffsetInMilliseconds int    `json:"offsetInMilliseconds,omitempty"`
	PlayerActivity       string `json:"playerActivity"`
}

// requestEnvelopeDataProvider provides a way to set Context and Session metadata for common requests.
type requestEnvelopeDataProvider interface {
	setContext(ctx *Context)
	setSession(session *Session)
}

// CommonRequest contains the attributes all alexa requests have in common.
type CommonRequest struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	Locale    string `json:"locale"`
	// Set manually from request envelope
	Session *Session
	Context *Context
}

// LaunchRequest send by Alexa if a skill is started.
type LaunchRequest struct {
	CommonRequest
}

// IntentRequest is send if a intent is invoked.
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

// SessionEndedRequest if a skill is stopped or cancelled.
type SessionEndedRequest struct {
	CommonRequest
	Reason string      `json:"reason,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

// AudioPlayerRequest for input events of audio player interface.
type AudioPlayerRequest struct {
	CommonRequest
	// tbd
}

// GetTypedRequest provides the request object mapped to the given struct
func (requestEnvelope *RequestEnvelope) GetTypedRequest(requestObj interface{}) error {
	data, _ := json.Marshal(requestEnvelope.Request)
	requestObj.(requestEnvelopeDataProvider).setContext(&requestEnvelope.Context)
	requestObj.(requestEnvelopeDataProvider).setSession(&requestEnvelope.Session)
	return json.Unmarshal(data, &requestObj)
}

func (cr *CommonRequest) setContext(ctx *Context) {
	cr.Context = ctx
}
func (cr *CommonRequest) setSession(session *Session) {
	cr.Session = session
}

// VerifyTimestamp checks if the the timestamp is not older than 30 seconds
func (requestEnvelope *RequestEnvelope) VerifyTimestamp() bool {
	timestampStr := requestEnvelope.Request.(map[string]interface{})["timestamp"].(string)
	requestTimestamp, err := time.Parse("2006-01-02T15:04:05Z", timestampStr)
	if err != nil {
		log.Fatalln("Error parsing request timestamp with value ", timestampStr)
	}
	if time.Since(requestTimestamp) < time.Duration(30)*time.Second {
		return true
	}
	return false
}
