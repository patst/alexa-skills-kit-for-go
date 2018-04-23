package alexa

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudioPlayerDirectives(t *testing.T) {
	skill := Skill{
		ApplicationID:  "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
		SkipValidation: true,
	}
	skillHandler := skill.GetHTTPSkillHandler()

	skill.OnAudioPlayerState = func(request *AudioPlayerRequest, response *ResponseEnvelope) {
		response.Response.AddAudioPlayerStopDirective()
		response.Response.AddAudioPlayerClearQueueDirective("CLEAR_ENQUEUED")
		play := response.Response.AddAudioPlayerPlayDirective("play")
		play.SetAudioItemStream("url", "token", "previous", 0)
		play.SetAudioItemMetadata("title", "subtitle")
	}

	launchRequestReader, err := os.Open("../resources/audioplayer_request.json")
	if err != nil {
		t.Error("Error reading input file", err)
	}

	httpRequest := httptest.NewRequest("POST", "/", launchRequestReader)
	responseWriter := httptest.NewRecorder()
	skillHandler.ServeHTTP(responseWriter, httpRequest)
	if responseWriter.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			responseWriter.Code, http.StatusOK)
	}
}

func TestInvalidClearQueueMode(t *testing.T) {
	var resp Response
	var buf bytes.Buffer

	//Redirect log output
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	//Invalid clearBehavior value
	d := resp.AddAudioPlayerClearQueueDirective("invalidBehavior")
	assert.NotNil(t, d)
	assert.Contains(t, string(buf.Bytes()), invalidClearBehaviorStr)
	buf.Reset()

	//Valid clearBehavior value
	d2 := resp.AddAudioPlayerClearQueueDirective("CLEAR_ALL")
	assert.NotNil(t, d2)
	assert.NotContains(t, string(buf.Bytes()), invalidClearBehaviorStr)
}
