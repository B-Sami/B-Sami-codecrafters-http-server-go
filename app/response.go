package main

import (
	"fmt"
)

type StatusLine struct {
	HTTPVersion  string
	StatusCode   int
	StatusPhrase string
}

type HTTPResponse struct {
	StatusLine   StatusLine
	Headers      map[string]string
	ResponseBody string
}

func NewServerResponse(statusCode int, statusPhrase string, headers map[string]string, body string) *HTTPResponse {
	return &HTTPResponse{
		StatusLine: StatusLine{
			HTTPVersion:  "HTTP/1.1",
			StatusCode:   statusCode,
			StatusPhrase: statusPhrase,
		},
		Headers:      headers,
		ResponseBody: body,
	}
}

func (r *HTTPResponse) FormatStatusLine() string {
	return fmt.Sprintf("%s %d %s\r\n", r.StatusLine.HTTPVersion, r.StatusLine.StatusCode, r.StatusLine.StatusPhrase)
}

func (r *HTTPResponse) FormatHeaders() string {
	var headers string
	for key, value := range r.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	if headers == "" {
		return "\r\n"
	}
	return headers + "\r\n"
}

func (r *HTTPResponse) FormatBody() string {
	return r.ResponseBody
}

func (r *HTTPResponse) FullResponse() string {
	return r.FormatStatusLine() + r.FormatHeaders() + r.FormatBody()
}
