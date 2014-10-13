package zippo

type Payload struct {
	URL         string `json:"url"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

func (p Payload) String() string {
	return p.Filename + "::" + p.URL
}

type Response struct {
	Zipname     string    `json:"zipname"`
	Payloads    []Payload `json:"payloads"`
	Length      int       `json:"length"`
	ContentType string    `json:"content_type"`
}
