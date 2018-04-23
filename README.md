## `alexa-skills-kit-for-go`: Create Amazon Alexa Go skills

[![GoDoc][1]][2]
[![GoCard][3]][4]
[![Build Status][5]][6]
[![codecov][7]][8]
[![License][9]][10]

[1]: https://godoc.org/github.com/patst/alexa-skills-kit-for-go?status.svg
[2]: https://godoc.org/github.com/patst/alexa-skills-kit-for-go
[3]: https://goreportcard.com/badge/patst/alexa-skills-kit-for-go
[4]: https://goreportcard.com/report/patst/alexa-skills-kit-for-go
[5]: https://travis-ci.org/patst/alexa-skills-kit-for-go.svg?branch=master
[6]: https://travis-ci.org/patst/alexa-skills-kit-for-go
[7]: https://codecov.io/gh/patst/alexa-skills-kit-for-go/branch/master/graph/badge.svg
[8]: https://codecov.io/gh/patst/alexa-skills-kit-for-go
[9]: https://img.shields.io/badge/License-Apache%202.0-blue.svg
[10]: https://github.com/patst/alexa-skills-kit-for-go/blob/master/LICENSE

This library is for building Amazon Alexa custom skills with Go. Unfortunately Amazon does only provide Alexa SDKs for building skills with Java and NodeJS.

Supported features:

* Amazon Alexa Custom skills in general
* Custom Skill as webservice backend ([AWS - Host a Custom Skill as Web Service](https://developer.amazon.com/docs/custom-skills/host-a-custom-skill-as-a-web-service.html))
* Custom Skill as Lambda function ([AWS - Host a Custom Skill as an AWS Lambda Function](https://developer.amazon.com/docs/custom-skills/host-a-custom-skill-as-an-aws-lambda-function.html))
* GameEngine and Gadget Skill API ([AWS - Understand Gadgets Skill API](https://developer.amazon.com/docs/gadget-skills/understand-gadgets-skill-api.html))
* Dialog Interface ([AWS - Dialog Interface Reference](https://developer.amazon.com/docs/custom-skills/dialog-interface-reference.html))
* Display Interface ([AWS - Display Interface Reference](https://developer.amazon.com/docs/custom-skills/display-interface-reference.html))
* AudioPlayer Interface ([AWS - AudioPlayer Interface Reference](https://developer.amazon.com/docs/custom-skills/audioplayer-interface-reference.html))
* Device Address Service ([AWS - Enhance you skill with customer address information](https://developer.amazon.com/docs/custom-skills/device-address-api.html))
* SessionStorage - store data in session attribute

There is a excellent API description what attributes must be included in responses and how to use the different interfaces in the [AWS Request and Response JSON reference](https://developer.amazon.com/docs/custom-skills/request-and-response-json-reference.html)

## Functionality

Amazon has an exhaustive description of the Alexa API including all the request and response objects. This library mapped the objects to structs and provides a way to interact with the Alexa Voice Service.

The connection is either via Web Service for self hosted custom skills or AWS Lambda Function for skill hosted at the AWS Lambda Function service.
The second option is a easier deployment way because a developer has not to think about valid HTTPS certificates.

## Examples

The examples show the usage of the library for different Alexa controllers. Every example resides in its own folder and does not reference other examples.

Usage:

* Copy the Go file(s) from one of the examples folders (e.g. `examples/helloworld/main.go`) in a Go project
* Run `go get ./...`
* Run `go build`
* Execute the binary!

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

### Hello buttons lambda

Example how to use the new Echo Gadget Buttons. The example is converted from the Amazon NodeJS [Hello buttons tutorial](https://github.com/alexa/skill-sample-nodejs-buttons-hellobuttons).

The example deploys as a lambda function. The necessary steps to create the skill are described in the above tutorial.

## Next Steps

* Get the Amazon NodeJS hello buttons example working
* Provide the so far not supported features (Services, AudioPlayer Interface, Dialog Interface)

## Author

Patrick Steinig ([@patst87](http://twitter.com/patst87))
