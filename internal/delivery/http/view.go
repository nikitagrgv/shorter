package http

type ResultPageData struct {
	ShortURL string
}

type ErrorPageData struct {
	ErrorCode        int
	ErrorDescription string
}
