package oauthsdk

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type mockResponse struct {
	Url        string
	StatusCode int
	Body       string
}

type requestClientMock struct {
	Responses []mockResponse
	Error     error
}

func (r requestClientMock) Do(req *http.Request) (*http.Response, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	var mockRes *mockResponse
	for _, res := range r.Responses {
		if strings.Contains(req.URL.String(), res.Url) {
			mockRes = &res
			break
		}
	}

	httpRes := &http.Response{}
	if mockRes == nil {
		return httpRes, fmt.Errorf(fmt.Sprintf("Not found mock response for url %v", req.RequestURI))
	}

	httpRes.StatusCode = mockRes.StatusCode
	httpRes.Body = io.NopCloser(strings.NewReader(mockRes.Body))

	return httpRes, nil
}
