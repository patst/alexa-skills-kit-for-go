package alexa

// OutgoingResponse is the complete object returned for a alexa POST request.
type OutgoingResponse struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          *Response              `json:"response,omitempty"`
}

type Response struct {
	OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
	Card             *Card         `json:"card,omitempty"`
	Reprompt         *Reprompt     `json:"reprompt,omitempty"`
	ShouldEndSession bool          `json:"shouldEndSession,omitempty"`
	Directives       []interface{} `json:"directives,omitempty"`
}

type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	Ssml string `json:"ssml,omitempty"`
}

type Reprompt struct {
	OutputSpeech OutputSpeech `json:"outputSpeech"`
}

type Card struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Text    string `json:"text,omitempty"`
	Image   struct {
		SmallImageURL string `json:"smallImageUrl,omitempty"`
		LargeImageURL string `json:"largeImageUrl,omitempty"`
	} `json:"image,omitempty"`
}

// NewOutgoingResponse creates a response skeletion for alexa responses
func NewOutgoingResponse(sessionAttributes map[string]interface{}) *OutgoingResponse {
	if sessionAttributes == nil {
		sessionAttributes = make(map[string]interface{})
	}
	return &OutgoingResponse{
		Version:           "1.0",
		Response:          &Response{},
		SessionAttributes: sessionAttributes,
	}
}

func (response *Response) SetOutputSpeech(text string) *Response {
	response.OutputSpeech = &OutputSpeech{
		Type: "PlainText",
		Text: text,
	}
	return response
}
func (response *Response) SetReprompt(text string) *Response {
	response.Reprompt = &Reprompt{
		OutputSpeech: OutputSpeech{
			Type: "PlainText",
			Text: text,
		},
	}
	return response
}
func (response *Response) SimpleCard(title string, content string) *Response {
	response.Card = &Card{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}
	return response
}
