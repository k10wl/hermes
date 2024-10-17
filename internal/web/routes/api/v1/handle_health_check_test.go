package v1

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHealthCheck(t *testing.T) {
	srv := httptest.NewServer(handleCheckHeath())
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s", srv.URL))

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("status not OK")
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Fatalf("body is not OK, actual %q\n", string(body))
	}
}
