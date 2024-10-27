package clients

import (
	"bytes"
	"net/http"
)

type RestClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

func NewRestClient(baseURL, apiKey string) *RestClient {
	return &RestClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		apiKey:     apiKey,
	}
}

// DoRequest sends a request to the specified endpoint with the given method and body.
func (rc *RestClient) DoRequest(method, endpoint string, body *bytes.Buffer) (*http.Response, error) {
	url := rc.baseURL + endpoint
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+rc.apiKey)

	return rc.httpClient.Do(req)
}
