package types

type ContactResponse struct {
	Success bool `json:"success"`
}

type ContactEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int32  `json:"color"`
	Fields      []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"fields"`
}

type ContactBody struct {
	Username  string         `json:"username"`
	AvatarURL string         `json:"avatar_url"`
	Embeds    []ContactEmbed `json:"embeds"`
}
