package alexa

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLambdaCall(t *testing.T) {
	skill := Skill{
		ApplicationID: "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
		OnLaunch: func(req *LaunchRequest, res *ResponseEnvelope) {
			res.Response.SetSimpleCard("title", "test")
		},
		SkipValidation: true,
	}
	skillHandler := skill.GetLambdaSkillHandler()

	launchRequestReader, err := os.Open("../resources/lambda_launch_request.json")
	if err != nil {
		t.Error("Error reading input file", err)
	}

	var event map[string]interface{}
	json.NewDecoder(launchRequestReader).Decode(&event)

	result, err := skillHandler(context.TODO(), event)

	assert.NoError(t, err)

	// result is a Outgoing response object
	responseEnvelope := result.(*ResponseEnvelope)
	assert.Equal(t, &Card{Type: "Simple", Title: "title", Content: "test"}, responseEnvelope.Response.Card)
}
