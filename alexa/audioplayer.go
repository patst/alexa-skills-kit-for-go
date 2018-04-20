package alexa

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
