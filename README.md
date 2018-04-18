## `alexa-skills-kit-for-go`: Create Amazon Alexa Go skills 

[![Go Report Card](https://goreportcard.com/badge/patst/alexa-skills-kit-for-go)](https://goreportcard.com/report/patst/alexa-skills-kit-for-go) [![Build Status](https://travis-ci.org/patst/alexa-skills-kit-for-go.svg?branch=master)](https://travis-ci.org/patst/alexa-skills-kit-for-go) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/patst/alexa-skills-kit-for-go/blob/master/LICENSE) [![GoDoc](https://godoc.org/github.com/patst/alexa-skills-kit-for-go?status.svg)](https://godoc.org/github.com/patst/alexa-skills-kit-for-go)

This library is for building Amazon Alexa custom skills with Golang. Unfortunately Amazon does only provide Alexa SDKs for building skills with Java and NodeJS.

Supported features:

*Amazon Alexa Custom skills in general
*Custom Skill as webservice backend ([AWS - Host a Custom Skill as Web Service](https://developer.amazon.com/docs/custom-skills/host-a-custom-skill-as-a-web-service.html))
*Custom Skill as Lambda function ([AWS - Host a Custom Skill as an AWS Lambda Function](https://developer.amazon.com/docs/custom-skills/host-a-custom-skill-as-an-aws-lambda-function.html))
*GameEngine and Gadget Skill API ([AWS - Understand Gadgets Skill API](https://developer.amazon.com/docs/gadget-skills/understand-gadgets-skill-api.html))
*SessionStorage - store data in session attribute

Not supported so far:
*AudioPlayer interface
*Dialog Interface
*Device Address Service
*List Management Service
*Asynchronous Directive Service

There is a excellent API description what attributes must be included in responses and how to use the different interfaces in the [AWS Request and Response JSON reference](https://developer.amazon.com/docs/custom-skills/request-and-response-json-reference.html)

## Examples

The examples show the usage of the library for different Alexa controllers. Every example resides in its own folder and does not reference other examples.

Usage:

*Copy the Go file(s) from one of the examples folders (e.g. `examples/helloworld/main.go`) in a Go project
*Run `go get ./...`
*Run `go build`
*Execute the binary!

The AWS Alexa skill setup is described in the original tutorials.
The ApplicationID from the Alexa skill must be configured in the skill definition.

You need a HTTPS endpoint configured to run this skill or deploy the skill as a lambda method (coming soon).

### Hello world

There is a [hello world example](example/helloworld.go) which reimplemented the [Amazon hello world example](https://github.com/alexa/alexa-skills-kit-sdk-for-java/tree/2.0.x/samples).
You can see how to implement different Alexa intents and use a SessionEnded and LaunchRequest handler.

### Hello world lambda

Same examples as 'hello world' but it can be deployed as Amazon Lambda expression.

If you are using windows like I am it may be tricky to build the zip file with the executable file. The AWS lambda library contains a tool called `build-lambda-zip` for building a .zip file with a executable binary. (See [Building your function](https://github.com/aws/aws-lambda-go))

The Lambda Function can be deployed using the AWS CLI, AWS CloudFormation or the AWS Web console. For details see [deploying Lambda Functions](https://docs.aws.amazon.com/lambda/latest/dg/deploying-lambda-apps.html).

## Next Steps

*Get the Amazon NodeJS hello buttons example working
*Provide the so far not supported features (Services, AudioPlayer Interface, Dialog Interface)

## Author

Patrick Steinig ([@patst87](http://twitter.com/patst87))
