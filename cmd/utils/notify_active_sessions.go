package utils

import (
	"bytes"
	"net/http"

	"github.com/k10wl/hermes/internal/core"
)

func NotifyActiveSessions(c *core.Core, id string, data []byte) {
	config := c.GetConfig()
	db := c.GetDB()
	activeSession, err := db.GetActiveSessionByDatabaseDNS(config.DatabaseDSN)
	if err != nil {
		return
	}
	body := new(bytes.Buffer)
	body.Write(data)
	req, err := http.NewRequest(
		"POST",
		activeSession.Address+"/api/v1/relay",
		body,
	)
	req.Header.Add("ID", id)
	if err != nil {
		return
	}
	if _, err = http.DefaultClient.Do(req); err != nil {
		db.RemoveActiveSession(activeSession) // unreachable, remove from active sessions
	}
}
