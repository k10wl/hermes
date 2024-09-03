package ai_clients

import (
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
	return io.ReadAll(res.Body)
}
