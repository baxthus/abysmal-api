package iplogger

type Response struct {
	Success bool `json:"success"`
}

type Whois struct {
	IP          string          `json:"ip"`
	Success     bool            `json:"success"`
	Type        string          `json:"type"`
	Continent   string          `json:"continent"`
	Country     string          `json:"country"`
	CountryCode string          `json:"country_code"`
	Region      string          `json:"region"`
	City        string          `json:"city"`
	Latitude    int             `json:"latitude"`
	Longitude   int             `json:"longitude"`
	Postal      string          `json:"postal"`
	Connection  WhoisConnection `json:"connection"`
	Timezone    WhoisTimezone   `json:"timezone"`
}

type WhoisConnection struct {
	ASN    int    `json:"asn"`
	Org    string `json:"org"`
	ISP    string `json:"isp"`
	Domain string `json:"domain"`
}

type WhoisTimezone struct {
	ID     string `json:"id"`
	Offset int    `json:"offset"`
	UTC    string `json:"utc"`
}

type Embed struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Thumbnail   EmbedThumbnail `json:"thumbnail"`
	Color       int32          `json:"color"`
	Fields      []EmbedField   `json:"fields"`
}

type EmbedThumbnail struct {
	URL string `json:"url"`
}

type EmbedField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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
	Components []ComponentsRow `json:"components"`
}
