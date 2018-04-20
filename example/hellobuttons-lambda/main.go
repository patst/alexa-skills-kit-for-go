package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/patst/alexa-skills-kit-for-go/alexa"
)

func main() {
	skill := alexa.Skill{
		ApplicationID:     "FILL WITH SKILL ID", // Echo App ID from Amazon Dashboard
		OnIntent:          intentDispatchHandler,
		OnLaunch:          launchRequestHandler,
		OnSessionEnded:    sessionEndedRequestHandler,
		OnSystemException: systemExceptionHandler,
		OnGameEngineEvent: gameEngineInputEventHandler,
	}
	skillHandler := skill.GetLambdaSkillHandler()

	lambda.Start(skillHandler)
}

func intentDispatchHandler(request *alexa.IntentRequest, responseEnvelope *alexa.ResponseEnvelope) {
	switch request.Intent.Name {
	case "AMAZON.StopIntent":
		stopIntentHandler(request, responseEnvelope)
	case "AMAZON.CancelIntent":
		cancelIntentHandler(request, responseEnvelope)
	case "AMAZON.HelpIntent":
		helpIntentHandler(request, responseEnvelope)
	case "Unhandled":
	default:
		log.Println("Unknown intent!", request.Intent.Name)
		unhandledIntentHandler(request, responseEnvelope)
	}
}

func gameEngineInputEventHandler(request *alexa.GameEngineInputHandlerEventRequest, responseEnvelope *alexa.ResponseEnvelope) {
	jsonR, _ := json.Marshal(request)
	log.Println("Received GameEngine event ", string(jsonR))

	if request.Type == "GameEngine.InputHandlerEvent" {
		for _, event := range request.Events {
			switch event.Name {
			case "button_down_event":
				gadgetID := event.InputEvents[0].GadgetID
				if responseEnvelope.SessionAttributes[gadgetID+"_initialized"] == nil {
					//This is a new button
					buttonCountVal := responseEnvelope.SessionAttributes["buttonCount"]
					buttonCount := int(buttonCountVal.(float64)) + 1
					responseEnvelope.SessionAttributes["buttonCount"] = buttonCount
					responseEnvelope.SessionAttributes[gadgetID+"_initialized"] = true

					/*
					   This is a new button, as in new to our understanding.
					   Because this button may have just woken up, it may not have
					   received the initial animations during the launch intent.
					   We'll resend the animations here, but instead of the empty array
					   broadcast above, we'll send the animations ONLY to this buttonId.
					*/
					responseEnvelope.Response.AddDirective(buildButtonIdleAnimationDirective([]string{gadgetID}, breathAnimationRed))
					responseEnvelope.Response.AddDirective(buildButtonDownAnimationDirective([]string{gadgetID}))
					responseEnvelope.Response.AddDirective(buildButtonUpAnimationDirective([]string{gadgetID}))

					// Say something when we first encounter a button.
					responseEnvelope.Response.SetOutputSpeech("Hello button " + strconv.Itoa(buttonCount) + ".")
				}
			case "button_up_event":
				gadgetID := event.InputEvents[0].GadgetID

				newAnimationIndex := 1
				//  On releasing the button, we'll replace the 'none' animation with a new color from a set of animations.
				if animationIndex, ok := responseEnvelope.SessionAttributes[gadgetID]; ok {
					//Gadget does alreay exist increase index by one
					newAnimationIndex = int(animationIndex.(float64)) + 1
					if newAnimationIndex >= len(animations) {
						newAnimationIndex = 0
					}
				}
				responseEnvelope.SessionAttributes[gadgetID] = newAnimationIndex

				responseEnvelope.Response.AddDirective(buildButtonIdleAnimationDirective([]string{gadgetID}, animations[newAnimationIndex]))
			case "timeout":
				if buttonCount, ok := request.Session.Attributes["buttonCount"]; ok {
					responseEnvelope.Response.SetOutputSpeech(fmt.Sprintf("Thank you for using the Gadgets Test Skill. I counted %d buttons. Goodbye.", int(buttonCount.(float64))))
				} else {
					responseEnvelope.Response.SetOutputSpeech("I didn't detect any buttons.  You must have at least one Echo Button to use this skill. Goodbye.")
				}
				responseEnvelope.Response.ShouldEndSession = true
			}
		}
	}
}

func sessionEndedRequestHandler(request *alexa.SessionEndedRequest, response *alexa.ResponseEnvelope) {
	log.Printf("Session with SessionID %v ended \n", request.Session.SessionID)
}

func launchRequestHandler(request *alexa.LaunchRequest, responseEnvelope *alexa.ResponseEnvelope) {
	log.Println("Launch request received. SessionID: ", request.Session.SessionID)

	directive := responseEnvelope.Response.AddGameEngineStartInputDirective(30000)
	buttonDownRecognizer := directive.AddPatternRecognizer("button_down_recognizer")
	buttonDownRecognizer.Fuzzy = false
	buttonDownRecognizer.Anchor = "end"
	buttonDownRecognizer.AddPattern(nil, nil, "down")

	buttonUpRecognizer := directive.AddPatternRecognizer("button_up_recognizer")
	buttonUpRecognizer.Fuzzy = false
	buttonUpRecognizer.Anchor = "end"
	buttonUpRecognizer.AddPattern(nil, nil, "up")

	buttonDownEvent := directive.AddEvent("button_down_event", false, []string{"button_down_recognizer"})
	buttonDownEvent.Reports = "matches"

	buttonUpEvent := directive.AddEvent("button_up_event", false, []string{"button_up_recognizer"})
	buttonUpEvent.Reports = "matches"

	timeoutEvent := directive.AddEvent("timeout", true, []string{"timed out"})
	timeoutEvent.Reports = "history"

	// Preserve the originatingRequestId.  We'll use this to stop the InputHandler later. See the Note at https://developer.amazon.com/docs/gadget-skills/receive-echo-button-events.html#start
	responseEnvelope.SessionAttributes["inputHandler_originatingRequestId"] = request.RequestID

	// Start keeping track of some state
	responseEnvelope.SessionAttributes["buttonCount"] = 0

	/*
	   If the buttons are awake before the Skill starts, the Skill can send
	   animations to all of the buttons by targeting the empty array [].
	*/

	// Build the breathing animation that will play immediately.
	responseEnvelope.Response.AddDirective(buildButtonIdleAnimationDirective([]string{}, breathAnimationRed))

	// Build the 'button down' animation for when the button is pressed.
	responseEnvelope.Response.AddDirective(buildButtonDownAnimationDirective([]string{}))

	// Build the 'button up' animation for when the button is released.
	responseEnvelope.Response.AddDirective(buildButtonUpAnimationDirective([]string{}))

	responseEnvelope.Response.SetOutputSpeech("Welcome to the Gadgets Test Skill. Press your Echo Buttons to change the colors of the lights. <audio src='https://s3.amazonaws.com/ask-soundlibrary/foley/amzn_sfx_rhythmic_ticking_30s_01.mp3'/>")

	json, err := json.Marshal(responseEnvelope)
	if err != nil {
		log.Fatal("Error serializing response", err)
	}
	log.Println("--> Response: ", string(json))
}

func systemExceptionHandler(request *alexa.SystemExceptionEncounteredRequest, responseEnvelope *alexa.ResponseEnvelope) {
	marshalled, _ := json.Marshal(request)
	log.Printf("System Exception encountered. RequestId %v . Request JSON: %v \n", request.RequestID, string(marshalled))
}

// Intent handling start

func helpIntentHandler(request *alexa.IntentRequest, responseEnvelope *alexa.ResponseEnvelope) {
	responseEnvelope.Response.SetOutputSpeech("Welcome to the Gadgets Test Skill. Press your Echo Buttons to change the lights. <audio src='https://s3.amazonaws.com/ask-soundlibrary/foley/amzn_sfx_rhythmic_ticking_30s_01.mp3'/>")
}

func stopIntentHandler(request *alexa.IntentRequest, responseEnvelope *alexa.ResponseEnvelope) {
	responseEnvelope.Response.SetOutputSpeech("Thank you for using the Gadgets Test Skill.  Goodbye.")
	if originatingRequestID, ok := request.Session.Attributes["inputHandler_originatingRequestId"]; ok {
		responseEnvelope.Response.AddGameEngineStopInputHandlerDirective(originatingRequestID.(string))
	}
	responseEnvelope.Response.ShouldEndSession = true
}

func cancelIntentHandler(request *alexa.IntentRequest, responseEnvelope *alexa.ResponseEnvelope) {
	responseEnvelope.Response.SetOutputSpeech("Thank you for using the Gadgets Test Skill.  Goodbye.")
	if originatingRequestID, ok := request.Session.Attributes["inputHandler_originatingRequestId"]; ok {
		responseEnvelope.Response.AddGameEngineStopInputHandlerDirective(originatingRequestID.(string))
	}
	responseEnvelope.Response.ShouldEndSession = true
}

func unhandledIntentHandler(request *alexa.IntentRequest, responseEnvelope *alexa.ResponseEnvelope) {
	responseEnvelope.Response.SetOutputSpeech("Sorry, I didn't get that.  Please press your Echo Buttons to change the color of the lights. <audio src='https://s3.amazonaws.com/ask-soundlibrary/foley/amzn_sfx_rhythmic_ticking_30s_01.mp3'/>")
}

// Intent handling stop
