package http

type CreateLinkRequest struct {
	LongURL string `json:"long_url" validate:"required,url"`
}

type CreateLinkResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}
