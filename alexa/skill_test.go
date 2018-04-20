package alexa

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWrongApplicationId(t *testing.T) {
	wrongAppID := "wrong app id"

	launchRequestReader, err := os.Open("../resources/launch_request.json")
	if err != nil {
		t.Error("Error reading input file", err)
	}

	bodyBytes, _ := ioutil.ReadAll(launchRequestReader)
	var reqEnvelope RequestEnvelope
	json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&reqEnvelope)

	reqEnvelope.Request.(map[string]interface{})["timestamp"] = time.Now().Format("2006-01-02T15:04:05Z")
	err = reqEnvelope.isRequestValid(wrongAppID)
	assert.Error(t, err)
}
