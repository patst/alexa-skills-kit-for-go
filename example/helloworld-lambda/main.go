package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/patst/alexa-skills-kit-for-go/alexa"
)

var skill alexa.Skill

func main() {
	skill = alexa.Skill{
		ApplicationID:  "FILL WITH SKILL ID", // Echo App ID from Amazon Dashboard
		OnIntent:       intentDispatchHandler,
		OnLaunch:       launchRequestHandler,
		OnSessionEnded: sessionEndedRequestHandler,
	}
	skillHandler := skill.GetLambdaSkillHandler()

	lambda.Start(skillHandler)
}

func intentDispatchHandler(request *alexa.IntentRequest, response *alexa.ResponseEnvelope) {
	switch request.Intent.Name {
	case "HelloWorldIntent":
		helloWorldIntentHandler(request, response)
	case "AMAZON.StopIntent":
		cancelAndStopIntentHandler(request, response)
	case "AMAZON.CancelIntent":
		cancelAndStopIntentHandler(request, response)
	case "AMAZON.HelpIntent":
		helpIntentHandler(request, response)
	default:
		log.Println("Unknown intent!", request.Intent.Name)
	}
}

// HelloWorldIntent
func helloWorldIntentHandler(request *alexa.IntentRequest, responseEnvelope *alexa.ResponseEnvelope) {
	speechText := "Hello world"
	responseEnvelope.Response.SetOutputSpeech(speechText)
	responseEnvelope.Response.SetSimpleCard("HelloWorld", speechText)
}

func cancelAndStopIntentHandler(request *alexa.IntentRequest, responseEnvelope *alexa.ResponseEnvelope) {
	speechText := "Goodbye"
	responseEnvelope.Response.SetOutputSpeech(speechText)
	responseEnvelope.Response.SetSimpleCard("HelloWorld", speechText)
}

func helpIntentHandler(request *alexa.IntentRequest, responseEnvelope *alexa.ResponseEnvelope) {
	speechText := "You can say hello to me!"
	responseEnvelope.Response.SetOutputSpeech(speechText)
	responseEnvelope.Response.SetReprompt(speechText)
	responseEnvelope.Response.SetSimpleCard("HelloWorld", speechText)
}

func launchRequestHandler(request *alexa.LaunchRequest, responseEnvelope *alexa.ResponseEnvelope) {
	speechText := "Welcome to the Alexa Skills Kit, you can say hello"
	responseEnvelope.Response.SetOutputSpeech(speechText)
	responseEnvelope.Response.SetReprompt(speechText)
	responseEnvelope.Response.SetSimpleCard("HelloWorld", speechText)
}

func sessionEndedRequestHandler(request *alexa.SessionEndedRequest, responseEnvelope *alexa.ResponseEnvelope) {
	// cleanup stuff
}
