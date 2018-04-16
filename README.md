## `alexa-skills-kit-for-go`: Create Amazon Alexa Go skills 

[![Go Report Card](https://goreportcard.com/badge/patst/alexa-skills-kit-for-go)](https://goreportcard.com/report/patst/alexa-skills-kit-for-go) [![Build Status](https://travis-ci.org/patst/alexa-skills-kit-for-go.svg?branch=master)](https://travis-ci.org/patst/alexa-skills-kit-for-go) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/patst/alexa-skills-kit-for-go/blob/master/LICENSE)

Unfortunately Amazon does only provide a Alexa SDK for building skills with Java and NodeJS. 

As a Go learner I tried to implement a Alexa Skill for the new gadget buttons and failed because the Java SDK was not up-do-date and I don't like Javascript.
There is another great project by mikeflynn on Github called [go-alexa](https://github.com/mikeflynn/go-alexa). To understand Go better and learn I tried to make another implementation providing the GameEngine and Gadget Controller interfaces.

## Examples

I try to create some examples how to use the library. The examples will show how to use the different Alexa controllers.

### Hello world

There is a [hello world example](example/helloworld) which reimplemented the [Amazon hello world example](https://github.com/alexa/alexa-skills-kit-sdk-for-java/tree/2.0.x/samples).

## Next Steps

- Get the Amazon NodeJS hello buttons example working
- Provide a Lambda connector for getting the freedom to deploy the skill either as Amazon Lambda or self hosted HTTPS application.
- Improve test coverage and code quality (Criticism welcome, I am still a Go newbie!)
- The `*request` classes are as `interface{}` type in the RequestEnvelope and have to be serialized again to get a concrete type. That seems ugly.

## Author

Patrick Steinig ([@patst87](http://twitter.com/patst87))