package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/patst/alexa-skills-kit-for-go/alexa"
)

var skill alexa.Skill

func main() {
	router := mux.NewRouter()
	skill = alexa.Skill{
		ApplicationID:  "FILL WITH SKILL ID", // Echo App ID from Amazon Dashboard
		OnIntent:       intentDispatchHandler,
		OnLaunch:       launchRequestHandler,
		OnSessionEnded: sessionEndedRequestHandler,
	}
	skillHandler := skill.GetHTTPSkillHandler()

	router.Handle("/echo/api/trueorfalse", skillHandler).Methods("POST")

	srv := &http.Server{
		Handler: router,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Starting webserver")
	// If no NGINX proxy is used this must be ListenAndServeTLS because Alexa requires a HTTPS connection for all requests!
	log.Fatal(srv.ListenAndServe())
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
	responseEnvelope.Response.SimpleCard("HelloWorld", speechText)
}

func cancelAndStopIntentHandler(request *alexa.IntentRequest, responseEnvelope *alexa.ResponseEnvelope) {
	speechText := "Goodbye"
	responseEnvelope.Response.SetOutputSpeech(speechText)
	responseEnvelope.Response.SimpleCard("HelloWorld", speechText)
}

func helpIntentHandler(request *alexa.IntentRequest, responseEnvelope *alexa.ResponseEnvelope) {
	speechText := "You can say hello to me!"
	responseEnvelope.Response.SetOutputSpeech(speechText)
	responseEnvelope.Response.SetReprompt(speechText)
	responseEnvelope.Response.SimpleCard("HelloWorld", speechText)
}

func launchRequestHandler(request *alexa.LaunchRequest, responseEnvelope *alexa.ResponseEnvelope) {
	speechText := "Welcome to the Alexa Skills Kit, you can say hello"
	responseEnvelope.Response.SetOutputSpeech(speechText)
	responseEnvelope.Response.SetReprompt(speechText)
	responseEnvelope.Response.SimpleCard("HelloWorld", speechText)
}

func sessionEndedRequestHandler(request *alexa.SessionEndedRequest, responseEnvelope *alexa.ResponseEnvelope) {
	// cleanup stuff
}
