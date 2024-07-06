package entity

type ShortenRequestBody struct {
	Url string `json:"url"`
}

type ShortenResponseBody struct {
	Result string `json:"result"`
}
