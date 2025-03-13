package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/web/routes/api/v1/messages"
)

func handleChats(c *core.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		limit, err := strconv.Atoi(params.Get("limit"))
		if err != nil {
			limit = -1
		}
		startBeforeID, err := strconv.Atoi(params.Get("start-before-id"))
		if err != nil {
			startBeforeID = -1
		}

		query := core.NewGetChatsQuery(c, int64(limit), int64(startBeforeID))
		err = query.Execute(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s\n", err)
			return
		}
		bytes, err := json.Marshal(query.Result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s\n", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(bytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s\n", err)
			return
		}
	}
}

func handleCheckHeath() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	}
}

func handleRelay(relay chan []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		messages.BroadcastData(relay, data)
		w.WriteHeader(http.StatusOK)
	}
}
