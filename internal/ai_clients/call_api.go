package ai_clients

import (
	"fmt"
	"io"
	"net/http"
)

func callApi(
	url string,
	body io.Reader,
	fillHeaders func(*http.Request) error,
) ([]byte, error) {
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	err = fillHeaders(request)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || 399 < res.StatusCode {
		return nil, fmt.Errorf("API error - %s\n", data)
	}
	return data, nil
}
