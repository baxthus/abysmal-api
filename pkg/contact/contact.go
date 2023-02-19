package contact

type Response struct {
	Success bool `json:"success"`
}

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int32  `json:"color"`
	Fields      []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"fields"`
}

type Body struct {
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Embeds    []Embed `json:"embeds"`
}
