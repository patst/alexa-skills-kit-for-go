package alexa

// OutgoingResponse is the complete object returned for a alexa POST request.
type OutgoingResponse struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          *Response              `json:"response,omitempty"`
}

// Response payload for alexa requests
type Response struct {
	OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
	Card             *Card         `json:"card,omitempty"`
	Reprompt         *Reprompt     `json:"reprompt,omitempty"`
	ShouldEndSession bool          `json:"shouldEndSession,omitempty"`
	Directives       []interface{} `json:"directives,omitempty"`
}

//OutputSpeech containing the speech to render to the user.
type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	Ssml string `json:"ssml,omitempty"`
}

// Reprompt containing the outputSpeech to use if a re-prompt is necessary.
type Reprompt struct {
	OutputSpeech OutputSpeech `json:"outputSpeech"`
}

// Card containing a card to render to the Amazon Alexa App
type Card struct {
	// A string describing the type of card to render. Values: 'Simple', 'Standard', 'LinkAccount'
	Type string `json:"type,omitempty"`
	// A string containing the title of the card. (not applicable for cards of type LinkAccount).
	Title string `json:"title,omitempty"`
	// A string containing the contents of a Simple card (not applicable for cards of type Standard or LinkAccount).
	Content string `json:"content,omitempty"`
	// A string containing the contents of a Simple card (not applicable for cards of type Standard or LinkAccount).
	Text string `json:"text,omitempty"`
	// An image object that specifies the URLs for the image to display on a Standard card. Only applicable for Standard cards.
	Image struct {
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

// SetOutputSpeech creates a PlainText output speech object for the response. Any present output speech is overwritten.
func (response *Response) SetOutputSpeech(text string) *Response {
	response.OutputSpeech = &OutputSpeech{
		Type: "PlainText",
		Text: text,
	}
	return response
}

// SetReprompt creates a PlainText reprompt output speech object for the response. Any present reprompt is overwritten.
func (response *Response) SetReprompt(text string) *Response {
	response.Reprompt = &Reprompt{
		OutputSpeech: OutputSpeech{
			Type: "PlainText",
			Text: text,
		},
	}
	return response
}

// SimpleCard creates a simple card for the response. Any present card is overwritten.
func (response *Response) SimpleCard(title string, content string) *Response {
	response.Card = &Card{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}
	return response
}
