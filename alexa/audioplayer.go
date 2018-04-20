package alexa

import (
	"log"
)

// AudioPlayerRequest represents an incoming request from the Audioplayer Interface. It does not have a session context.
// Response to such a request must be a AudioPlayerDirective or empty
type AudioPlayerRequest struct {
	CommonRequest
	Token                string `json:"token"`
	OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
}

// AudioPlayerPlaybackFailedRequest is sent when Alexa encounters an error when attempting to play a stream.
type AudioPlayerPlaybackFailedRequest struct {
	AudioPlayerRequest
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
	CurrentPlaybackState struct {
		Token                string `json:"token"`
		OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
		PlayerActivity       string `json:"playerActivity"`
	} `json:"currentPlaybackState"`
}

// AudioPlayerPlayDirective sends Alexa a command to stream the audio file identified by the specified audioItem. Use the playBehavior parameter to determine whether the stream begins playing immediately, or is added to the queue.
// shouldEndSession should be set to false otherwise playback will pause immediately
type AudioPlayerPlayDirective struct {
	Type         string `json:"type"`
	PlayBehavior string `json:"playBehavior"`
	AudioItem    struct {
		Stream struct {
			URL                   string `json:"url"`
			Token                 string `json:"token"`
			ExpectedPreviousToken string `json:"expectedPreviousToken,omitempty"`
			OffsetInMilliseconds  int    `json:"offsetInMilliseconds"`
		} `json:"stream"`
		Metadata struct {
			Title           string             `json:"title,omitempty"`
			Subtitle        string             `json:"subtitle,omitempty"`
			Art             DisplayImageObject `json:"art,omitempty"`
			BackgroundImage DisplayImageObject `json:"backgroundImage,omitempty"`
		} `json:"metadata,omitempty"`
	} `json:"audioItem"`
}

// AudioPlayerStopDirective stopts the current audio playback
type AudioPlayerStopDirective struct {
	Type string `json:"type"`
}

// AudioPlayerClearQueueDirective clears the audio playback queue. You can set this directive to clear the queue without stopping the currently playing stream, or clear the queue and stop any currently playing stream.
type AudioPlayerClearQueueDirective struct {
	Type          string `json:"type"`
	ClearBehavior string `json:"clearBehavior"`
}

// AddAudioPlayerPlayDirective creates a new play directive for AudioPlayer interfaces.
func (r *Response) AddAudioPlayerPlayDirective(playBehavior string) *AudioPlayerPlayDirective {
	playDirective := &AudioPlayerPlayDirective{
		Type:         "AudioPlayer.Play",
		PlayBehavior: playBehavior,
	}
	r.AddDirective(playDirective)
	return playDirective
}

// SetAudioItemStream sets the stream attributes for the audio item associated with the play directive.
func (d *AudioPlayerPlayDirective) SetAudioItemStream(url, token, expectedPreviousToken string, offsetInMilliseconds int) {
	d.AudioItem.Stream.URL = url
	d.AudioItem.Stream.Token = token
	d.AudioItem.Stream.ExpectedPreviousToken = expectedPreviousToken
	d.AudioItem.Stream.OffsetInMilliseconds = offsetInMilliseconds
}

// SetAudioItemMetadata sets the metadata attributes for the audio item associated with the play directive.
func (d *AudioPlayerPlayDirective) SetAudioItemMetadata(title, subtitle string) {
	d.AudioItem.Metadata.Title = title
	d.AudioItem.Metadata.Subtitle = subtitle
	// TODO set image
}

// AddAudioPlayerStopDirective creates a new stop directive for AudioPlayer interface.
func (r *Response) AddAudioPlayerStopDirective() *AudioPlayerStopDirective {
	stopDirective := &AudioPlayerStopDirective{
		Type: "AudioPlayer.Stop",
	}
	r.AddDirective(stopDirective)
	return stopDirective
}

// AddAudioPlayerClearQueueDirective creates a new clear queue directive for AudioPlayer interface.
func (r *Response) AddAudioPlayerClearQueueDirective(clearBehavior string) *AudioPlayerClearQueueDirective {
	// Must be one of the two values
	if clearBehavior != "CLEAR_ENQUEUED" && clearBehavior != "CLEAR_ALL" {
		log.Println("Invalid/ Unknown clearBehavior for ClearQueue directive!")
	}
	clearQueueDirective := &AudioPlayerClearQueueDirective{
		Type:          "AudioPlayer.ClearQueue",
		ClearBehavior: clearBehavior,
	}
	r.AddDirective(clearQueueDirective)
	return clearQueueDirective
}
