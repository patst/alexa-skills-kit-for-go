package alexa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	responseBytes, _ := ioutil.ReadAll(responseWriter.Body)
	fmt.Println("Response:", string(responseBytes))
	var response map[string]interface{}
	json.Unmarshal(responseBytes, &response)
	// The response must not contain image object in the metadata (optional)
	directives := response["response"].(map[string]interface{})["directives"].([]interface{})
	playDirective := directives[2].(map[string]interface{})
	require.Equal(t, "AudioPlayer.Play", playDirective["type"])
	require.Equal(t, "play", playDirective["playBehavior"])
	audioItem := playDirective["audioItem"].(map[string]interface{})
	stream := audioItem["stream"].(map[string]interface{})
	require.Equal(t, "url", stream["url"])
	require.Equal(t, "previous", stream["expectedPreviousToken"])
	require.Equal(t, "token", stream["token"])
	metadata := audioItem["metadata"].(map[string]interface{})
	require.Equal(t, "title", metadata["title"])
	require.Equal(t, "subtitle", metadata["subtitle"])
	require.Nil(t, metadata["art"])
	require.Nil(t, metadata["backgroundImage"])
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

func TestAudioPlayerDirectiveWithImage(t *testing.T) {
	skill := Skill{
		ApplicationID:  "amzn1.echo-sdk-ams.app.000000-d0ed-0000-ad00-000000d00ebe",
		SkipValidation: true,
	}
	skillHandler := skill.GetHTTPSkillHandler()

	skill.OnAudioPlayerState = func(request *AudioPlayerRequest, response *ResponseEnvelope) {
		play := response.Response.AddAudioPlayerPlayDirective("play")
		metadata := play.SetAudioItemMetadata("title", "subtitle")
		artImage := metadata.SetArtImage("artImage")
		artImage.AddImageSource("1", "url1", 2, 3)
		backgroundImage := metadata.SetBackgroundImage("backgroundImage")
		backgroundImage.AddImageSource("1", "url2", 2, 3)
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
	responseBytes, _ := ioutil.ReadAll(responseWriter.Body)
	fmt.Println("Response:", string(responseBytes))
	var response map[string]interface{}
	json.Unmarshal(responseBytes, &response)
	// The response must not contain image object in the metadata (optional)
	directives := response["response"].(map[string]interface{})["directives"].([]interface{})
	playDirective := directives[0].(map[string]interface{})
	require.Equal(t, "AudioPlayer.Play", playDirective["type"])
	metadata := playDirective["audioItem"].(map[string]interface{})["metadata"].(map[string]interface{})
	require.Equal(t, "title", metadata["title"])
	require.Equal(t, "subtitle", metadata["subtitle"])
	require.NotNil(t, metadata["art"])
	art := metadata["art"].(map[string]interface{})
	require.Equal(t, "artImage", art["contentDescription"])
	require.Equal(t, 1, len(art["sources"].([]interface{})))
	require.Equal(t, "1", art["sources"].([]interface{})[0].(map[string]interface{})["size"])
	require.Equal(t, "url1", art["sources"].([]interface{})[0].(map[string]interface{})["url"])
	require.Equal(t, 2.0, art["sources"].([]interface{})[0].(map[string]interface{})["heightPixels"])
	require.Equal(t, 3.0, art["sources"].([]interface{})[0].(map[string]interface{})["widthPixels"])

	require.NotNil(t, metadata["backgroundImage"])
	background := metadata["backgroundImage"].(map[string]interface{})
	require.Equal(t, "backgroundImage", background["contentDescription"])
	require.Equal(t, 1, len(background["sources"].([]interface{})))
	require.Equal(t, "1", background["sources"].([]interface{})[0].(map[string]interface{})["size"])
	require.Equal(t, "url2", background["sources"].([]interface{})[0].(map[string]interface{})["url"])
	require.Equal(t, 2.0, background["sources"].([]interface{})[0].(map[string]interface{})["heightPixels"])
	require.Equal(t, 3.0, background["sources"].([]interface{})[0].(map[string]interface{})["widthPixels"])
}
