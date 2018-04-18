package alexa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
)

type LambdaSkillHandler interface {
	Invoke(ctx context.Context, event interface{}) (interface{}, error)
}

type LambdaHandler func(ctx context.Context, event interface{}) (interface{}, error)

// GetLambdaSkillHandler provides a handler which can be used in a lambda function.
func (skill *Skill) GetLambdaSkillHandler() LambdaHandler {
	return func(ctx context.Context, event interface{}) (interface{}, error) {
		bodyBytes, err := json.Marshal(event)

		if err != nil {
			log.Println("Error reading request body. ", err)
			return nil, err
		}

		var requestEnvelope *RequestEnvelope
		err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&requestEnvelope)

		if err != nil {
			log.Println("Error decoding request body. ", err)
			return nil, err
		}

		if !isRequestValid(requestEnvelope, skill.ApplicationID, false, nil) {
			return nil, errors.New("Request is invalid")
		}

		response, err := handleRequest(requestEnvelope, skill)

		if err != nil {
			log.Println("Bad request.", err)
			return nil, err
		}

		return response, nil
	}
}
