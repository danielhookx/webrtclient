package client

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/oofpgDLD/webrtclient/example/client/signal/testsignal"
	"net/url"
	"sync/atomic"

	"log"
)

type Client struct {
	conn *websocket.Conn
	queue chan []byte

	isClosed int32
}

func NewClient(name string, u url.URL) (*Client, error){
	c := &Client{
		queue:    make(chan []byte, 1024),
		isClosed: 0,
	}
	log.Println("connect to ", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	regMsg := testsignal.RegisterMsg{Name:name}
	data, err := json.Marshal(regMsg)
	if err != nil {
		return nil, err
	}
	msg := testsignal.Msg{
		Type: 1,
		Data: string(data),
	}
	msgData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	err = c.conn.WriteMessage(websocket.TextMessage, msgData)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			c.queue <- message
		}
	}()
	return c, nil
}

func (c *Client) Close() error{
	if atomic.CompareAndSwapInt32(&c.isClosed, 0, 1) {
		return c.conn.Close()
	}
	return fmt.Errorf("")
}

func (c *Client) Read() []byte{
	return <-c.queue
}