package zippo

type Response struct {
	Zipname     string    `json:"zipname"`
	Payloads    []Payload `json:"payloads"`
	Length      int       `json:"length"`
	ContentType string    `json:"content_type"`
}
