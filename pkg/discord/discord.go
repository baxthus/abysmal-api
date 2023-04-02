package discord

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   struct {
		URL string `json:"url"`
	} `json:"thumbnail"`
	Color  int32 `json:"color"`
	Fields []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"fields"`
}

type ComponentsRow struct {
	Type       int8        `json:"type"`
	Components []Component `json:"components"`
}

type Component struct {
	Style    int8   `json:"style"`
	Label    string `json:"label"`
	URL      string `json:"url"`
	Type     int8   `json:"type"`
	Disabled bool   `json:"disabled"`
}

type Webhook struct {
	Username   string          `json:"username"`
	AvatarURL  string          `json:"avatar_url"`
	Embeds     []Embed         `json:"embeds"`
	Components []ComponentsRow `json:"components,omitempty"`
}
