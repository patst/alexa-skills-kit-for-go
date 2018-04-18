## `alexa-skills-kit-for-go`: Create Amazon Alexa Go skills 

[![Go Report Card](https://goreportcard.com/badge/patst/alexa-skills-kit-for-go)](https://goreportcard.com/report/patst/alexa-skills-kit-for-go) [![Build Status](https://travis-ci.org/patst/alexa-skills-kit-for-go.svg?branch=master)](https://travis-ci.org/patst/alexa-skills-kit-for-go) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/patst/alexa-skills-kit-for-go/blob/master/LICENSE) [![GoDoc](https://godoc.org/github.com/patst/alexa-skills-kit-for-go?status.svg)](https://godoc.org/github.com/patst/alexa-skills-kit-for-go)

Unfortunately Amazon does only provide a Alexa SDK for building skills with Java and NodeJS.

I tried to implement a Alexa Skill for the new gadget buttons and failed because the Java SDK was not up-do-date and I don't like Javascript.
There is another great project by mikeflynn on Github called [go-alexa](https://github.com/mikeflynn/go-alexa). To understand and learn Go I tried to make another implementation also providing the GameEngine and Gadget Controller interfaces.

## Examples

The examples show the usage of the different Alexa controllers.

Just run the Go file (e.g. helloworld.go). The AWS Alexa skill setup is described in the original tutorials.
The ApplicationID from the Alexa skill must be configured in the skill definition.

You need a HTTPS endpoint configured to run this skill or deploy the skill as a lambda method (coming soon).

### Hello world

There is a [hello world example](example/helloworld.go) which reimplemented the [Amazon hello world example](https://github.com/alexa/alexa-skills-kit-sdk-for-java/tree/2.0.x/samples).
You can see how to implement different Alexa intents and use a SessionEnded and LaunchRequest handler.

## Next Steps

- Get the Amazon NodeJS hello buttons example working
- Provide a Lambda connector for getting the freedom to deploy the skill either as Amazon Lambda or self hosted HTTPS application.
- Improve test coverage (Criticism welcome, I am still a Go newbie!)
- AudioController interface is a stub only at the moment. I do not have access to a Alexa AudioPlayer device at the moment.

## Author

Patrick Steinig ([@patst87](http://twitter.com/patst87))
