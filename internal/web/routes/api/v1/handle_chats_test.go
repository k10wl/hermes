package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"testing"

	"github.com/k10wl/hermes/internal/models"
	"github.com/k10wl/hermes/internal/test_helpers"
	"github.com/k10wl/hermes/internal/test_helpers/db_helpers"
)

func TestHandleChats(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	seeder := db_helpers.NewSeeder(db, context.Background())
	seeder.SeedChatsN(10)

	srv := httptest.NewServer(handleChats(coreInstance))
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s/api/v1/chats", srv.URL))

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

	generated := db_helpers.GenerateChatsSliceN(10)
	slices.Reverse(generated)
	for _, val := range generated {
		val.TimestampsToNilForTest__()
	}

	resData := []*models.Chat{}
	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err)
	}
	for _, val := range resData {
		val.TimestampsToNilForTest__()
	}

	expected := test_helpers.UnpointerSlice(generated)
	actual := test_helpers.UnpointerSlice(resData)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"return chats has bad output\nexpected: %+v\nactual:   %+v\n\n",
			expected,
			actual,
		)
	}
}

func TestHandleChatsWithLimitAndOffset(t *testing.T) {
	coreInstance, db := test_helpers.CreateCore()
	seeder := db_helpers.NewSeeder(db, context.Background())
	seeder.SeedChatsN(10)

	srv := httptest.NewServer(handleChats(coreInstance))
	defer srv.Close()

	res, err := http.Get(
		fmt.Sprintf("%s/api/v1/chats?limit=5&start-after-id=5", srv.URL),
	)

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

	generated := db_helpers.GenerateChatsSliceN(10)[5:]
	slices.Reverse(generated)
	for _, val := range generated {
		val.TimestampsToNilForTest__()
	}

	resData := []*models.Chat{}
	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err)
	}
	for _, val := range resData {
		val.TimestampsToNilForTest__()
	}

	expected := test_helpers.UnpointerSlice(generated)
	actual := test_helpers.UnpointerSlice(resData)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"bad output with chat limit and start after id\nexpected: %+v\nactual:   %+v\n\n",
			expected,
			actual,
		)
	}
}
