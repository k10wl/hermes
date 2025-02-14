package v1

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/k10wl/hermes/internal/ai_clients"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/web/routes/api/v1/messages"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 0,
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c Client) Single() chan []byte {
	return c.send
}

func (c Client) All() chan []byte {
	return c.hub.broadcast
}

func (c *Client) readPump(coreInstanse *core.Core, completionFn ai_clients.CompletionFn) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(
		func(string) error {
			c.conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		},
	)
	stderr := coreInstanse.GetConfig().Stderr

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %v", err)
			}
			break
		}

		clientMessage, err := messages.ReadMessage(message)
		if err != nil || clientMessage == nil {
			fmt.Fprintf(stderr, "errored upon message read %s\n", err.Error())
			if err := messages.BroadcastServerEmittedMessage(
				c.send,
				messages.NewServerError(
					uuid.New().String(),
					fmt.Sprintf(
						"malformed message, error: %s\ndata: %q\n",
						err,
						message,
					),
				),
			); err != nil {
				fmt.Fprintf(stderr, "failed to broadcast %s\n", err.Error())
			}
			continue
		}

		go func() {
			if err := clientMessage.Process(
				c,
				coreInstanse,
				completionFn,
			); err != nil {
				fmt.Fprintf(stderr, "processing error %s\n", err.Error())
				if err := messages.BroadcastServerEmittedMessage(
					c.send,
					messages.NewServerError(
						clientMessage.GetID(),
						fmt.Sprintf("malformed message, error: %s\ndata: %q\n",
							err,
							message,
						),
					),
				); err != nil {
					fmt.Fprintf(
						stderr,
						"failed to broadcast %s\n",
						err.Error(),
					)
				}
			}
		}()
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type message interface {
	encode() ([]byte, error)
}

func (c *Client) sendMessage(m message) error {
	data, err := m.encode()
	if err != nil {
		return err
	}
	c.send <- data
	return nil
}

func handleServeWebSockets(
	core *core.Core,
	hub *Hub,
	completionFn ai_clients.CompletionFn,
) http.HandlerFunc {
	config := core.GetConfig()

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		client := &Client{
			hub:  hub,
			conn: conn,
			send: make(chan []byte, 256),
		}

		client.hub.register <- client
		if r.URL.Query().Get("reconnect") == "true" {
			err := messages.BroadcastServerEmittedMessage(client.send, messages.NewServerReload())
			if err != nil {
				fmt.Fprintf(config.Stderr, "failed to send refresh message: %s\n", err)
			}
		}

		go client.writePump()
		go client.readPump(core, completionFn)
	}
}
