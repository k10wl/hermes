package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/k10wl/hermes/internal/core"
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
