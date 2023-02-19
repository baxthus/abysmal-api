package github

type Response struct {
	Success bool   `json:"success"`
	Hash    string `json:"hash,omitempty"`
}

type GithubResponse struct {
	Ref    string `json:"ref"`
	NodeID string `json:"node_id"`
	URL    string `json:"url"`
	Object struct {
		SHA  string `json:"sha"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"object"`
}
