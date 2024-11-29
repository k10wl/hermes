package test_helpers

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
)

func CreateWebsocketConnection(url string) (*websocket.Conn, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, cancel, err
	}
	return conn, cancel, err
}
