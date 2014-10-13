package zippo

type Payload struct {
	URL         string `json:"url"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type Response struct {
	Zipname     string    `json:"zipname"`
	Payloads    []Payload `json:"urls"`
	Length      int       `json:"length"`
	ContentType string    `json:"content_type"`
}
