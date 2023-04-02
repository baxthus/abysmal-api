package iplogger

type Response struct {
	Success bool `json:"success"`
}

type Whois struct {
	IP          string `json:"ip"`
	Success     bool   `json:"success"`
	Type        string `json:"type"`
	Continent   string `json:"continent"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Region      string `json:"region"`
	City        string `json:"city"`
	Latitude    int    `json:"latitude"`
	Longitude   int    `json:"longitude"`
	Postal      string `json:"postal"`
	Connection  struct {
		ASN    int    `json:"asn"`
		Org    string `json:"org"`
		ISP    string `json:"isp"`
		Domain string `json:"domain"`
	} `json:"connection"`
	Timezone struct {
		ID     string `json:"id"`
		Offset int    `json:"offset"`
		UTC    string `json:"utc"`
	} `json:"timezone"`
}
