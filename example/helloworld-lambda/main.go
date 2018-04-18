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

func intentDispatchHandler(request *alexa.IntentRequest, response *alexa.OutgoingResponse) {
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
func helloWorldIntentHandler(request *alexa.IntentRequest, outgoingResponse *alexa.OutgoingResponse) {
	speechText := "Hello world"
	outgoingResponse.Response.SetOutputSpeech(speechText)
	outgoingResponse.Response.SimpleCard("HelloWorld", speechText)
}

func cancelAndStopIntentHandler(request *alexa.IntentRequest, outgoingResponse *alexa.OutgoingResponse) {
	speechText := "Goodbye"
	outgoingResponse.Response.SetOutputSpeech(speechText)
	outgoingResponse.Response.SimpleCard("HelloWorld", speechText)
}

func helpIntentHandler(request *alexa.IntentRequest, outgoingResponse *alexa.OutgoingResponse) {
	speechText := "You can say hello to me!"
	outgoingResponse.Response.SetOutputSpeech(speechText)
	outgoingResponse.Response.SetReprompt(speechText)
	outgoingResponse.Response.SimpleCard("HelloWorld", speechText)
}

func launchRequestHandler(request *alexa.LaunchRequest, outgoingResponse *alexa.OutgoingResponse) {
	speechText := "Welcome to the Alexa Skills Kit, you can say hello"
	outgoingResponse.Response.SetOutputSpeech(speechText)
	outgoingResponse.Response.SetReprompt(speechText)
	outgoingResponse.Response.SimpleCard("HelloWorld", speechText)
}

func sessionEndedRequestHandler(request *alexa.SessionEndedRequest, outgoingResponse *alexa.OutgoingResponse) {
	// cleanup stuff
}
