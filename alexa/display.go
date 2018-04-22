package alexa

// DisplayRenderTemplateDirective directive to render display text, images or items on an device with screen.
type DisplayRenderTemplateDirective struct {
	Type string `json:"type,omitempty"`
	// Template is the body template to render
	Template DisplayTemplate `json:"template"`
}

// DisplayTemplate displays text and images. Types may either be BodyTemplate* or ListTemplate*.
// For a body template these images cannot be made selectable.
// List template displays a scrollable list of items, each with associated text and optional images.
// These images can be made selectable, as described in this reference.
type DisplayTemplate struct {
	Type  string `json:"type"`
	Token string `json:"token"`
	// BackButton state (e.g. 'VISIBLE' or 'HIDDEN')
	BackButton      string             `json:"backButton,omitempty"`
	BackgroundImage DisplayImageObject `json:"backgroundImage,omitempty"`
	Title           string             `json:"title,omitempty"`
	TextContent     struct {
		PrimaryText   DisplayTextContent `json:"primaryText,omitempty"`
		SecondaryText DisplayTextContent `json:"secondaryText,omitempty"`
		TertiaryText  DisplayTextContent `json:"tertiaryText,omitempty"`
	} `json:"textContent,omitempty"`
	// ListItems contains the text and images of the list items.
	ListItems []struct{} `json:"listItems,omitempty"`
}

// DisplayTextContent contains text and a text type for displaying text with the Display interface.
type DisplayTextContent struct {
	//Type must be PlainText or RichtText
	Type string `json:"type"`
	Text string `json:"text"`
}

// DisplayImageObject references and describes the image. Multiple sources for the image can be provided.
type DisplayImageObject struct {
	ContentDescription string `json:"contentDescription"`
	Sources            []struct {
		URL          string `json:"url"`
		Size         string `json:"size,omitempty"`
		WidthPixels  int    `json:"widthPixels,omitempty"`
		HeightPixels int    `json:"heightPixels,omitempty"`
	} `json:"sources"`
}

// AddDisplayRenderTemplateDirective creates a new directive to render a body or list template for the Alexa Display Interface.
func (r *Response) AddDisplayRenderTemplateDirective(templateType string) *DisplayRenderTemplateDirective {
	d := &DisplayRenderTemplateDirective{
		Type: "Display.RenderTemplate",
		Template: DisplayTemplate{
			Type: templateType,
		},
	}
	r.AddDirective(d)
	return d
}
