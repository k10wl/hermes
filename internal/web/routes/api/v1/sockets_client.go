package v1

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/k10wl/hermes/internal/core"
	"github.com/k10wl/hermes/internal/web/routes/api/v1/messages"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump(coreInstanse *core.Core) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(
		func(string) error {
			c.conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		},
	)
	stdout := coreInstanse.GetConfig().Stdoout
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
		if err != nil {
			fmt.Fprintf(stdout, "%s\n", err)
			err := messages.Broadcast(
				c.send,
				messages.NewErrorMessage(
					fmt.Sprintf("unhandled message: %s", message),
				),
			)
			if err != nil {
				fmt.Fprintf(stdout, "failed to broadcast %s\n", err)
			}
			continue
		}
		fmt.Fprintf(stdout, "received %+v\n", clientMessage)

		serverMessage, receiver, err := clientMessage.Process(coreInstanse)
		if err != nil {
			fmt.Fprintf(stdout, "%s\n", err)
			err := messages.Broadcast(
				c.send,
				messages.NewErrorMessage(
					fmt.Sprintf("failed to create server message: %s", message),
				),
			)
			if err != nil {
				fmt.Fprintf(stdout, "failed to broadcast %s\n", err)
			}
			continue
		}

		var to chan []byte
		switch receiver {
		case messages.None:
			continue
		case messages.Sender:
			to = c.send
		case messages.All:
			to = c.hub.broadcast
		}
		if err := messages.Broadcast(to, serverMessage); err != nil {
			fmt.Fprintf(stdout, "failed to broadcast %s\n", err)
		}
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

func handleServeWebSockets(core *core.Core, hub *Hub) http.HandlerFunc {
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
			// err := client.sendMessage(newMessage("reload", nil))
			if err != nil {
				fmt.Fprintf(config.Stderr, "failed to send refresh message: %s\n", err)
			}
		}

		go client.writePump()
		go client.readPump(core)
	}
}
