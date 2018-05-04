package alexa

// ResponseEnvelope is the envelope for the object returned for a alexa POST request.
type ResponseEnvelope struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          *Response              `json:"response,omitempty"`
}

// Response payload for alexa requests
type Response struct {
	OutputSpeech *OutputSpeech `json:"outputSpeech,omitempty"`
	Card         *Card         `json:"card,omitempty"`
	Reprompt     *Reprompt     `json:"reprompt,omitempty"`
	// Use a pointer to be able to specify true,false and do not set it
	ShouldEndSession *bool         `json:"shouldEndSession,omitempty"`
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
	OutputSpeech *OutputSpeech `json:"outputSpeech"`
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
	// A list of scope strings that maps to Alexa permissions.
	// Include only those Alexa permissions that are both needed by your skill and that are declared in your skill metadata on the Amazon Developer Portal.
	Permissions []string `json:"permissions"`
}

// NewResponseEnvelope creates a response skeletion for alexa responses
func newResponseEnvelope(sessionAttributes map[string]interface{}) *ResponseEnvelope {
	if sessionAttributes == nil {
		sessionAttributes = make(map[string]interface{})
	}
	return &ResponseEnvelope{
		Version:           "1.0",
		Response:          &Response{},
		SessionAttributes: sessionAttributes,
	}
}

// SetOutputSpeech creates a SSML output speech object for the response. Any present output speech is overwritten.
func (response *Response) SetOutputSpeech(text string) *Response {
	response.OutputSpeech = &OutputSpeech{
		Type: "SSML",
		Ssml: "<speak> " + text + " </speak>",
	}
	return response
}

// SetReprompt creates a PlainText reprompt output speech object for the response. Any present reprompt is overwritten.
func (response *Response) SetReprompt(text string) *Response {
	response.Reprompt = &Reprompt{
		OutputSpeech: &OutputSpeech{
			Type: "SSML",
			Ssml: "<speak> " + text + " </speak>",
		},
	}
	return response
}

// SetSimpleCard creates a simple card for the response. Any present card is overwritten.
func (response *Response) SetSimpleCard(title string, content string) *Response {
	response.Card = &Card{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}
	return response
}

// SetAskForPermissionsConsentCard creates a card to ask for permissions to read user or list data. Any present card is overwritten.
// Permission examples are 'read::alexa:device:all:address' or 'read::alexa:device:all:address:country_and_postal_code'
func (response *Response) SetAskForPermissionsConsentCard(title, content string, permissions []string) *Response {
	response.Card = &Card{
		Type:        "AskForPermissionsConsent",
		Title:       title,
		Content:     content,
		Permissions: permissions,
	}
	return response
}

// AddDirective adds a directive to the slice of existing directives for a response.
func (response *Response) AddDirective(directive interface{}) {
	response.Directives = append(response.Directives, directive)
}
