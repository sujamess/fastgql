package client

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

const (
	connectionInitMsg = "connection_init" // Client -> Server
	startMsg          = "start"           // Client -> Server
	connectionAckMsg  = "connection_ack"  // Server -> Client
	connectionKaMsg   = "ka"              // Server -> Client
	dataMsg           = "data"            // Server -> Client
	errorMsg          = "error"           // Server -> Client
)

type operationMessage struct {
	Payload json.RawMessage `json:"payload,omitempty"`
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
}

type Subscription struct {
	Close func() error
	Next  func(response interface{}) error
}

func errorSubscription(err error) *Subscription {
	return &Subscription{
		Close: func() error { return nil },
		Next: func(response interface{}) error {
			return err
		},
	}
}

func (p *Client) Websocket(query string, options ...Option) *Subscription {
	return p.WebsocketWithPayload(query, nil, options...)
}

// Grab a single response from a websocket based query
func (p *Client) WebsocketOnce(query string, resp interface{}, options ...Option) error {
	sock := p.Websocket(query)
	defer sock.Close()
	return sock.Next(&resp)
}

func (p *Client) WebsocketWithPayload(query string, initPayload map[string]interface{}, options ...Option) *Subscription {
	r, err := p.newRequest(query, options...)
	if err != nil {
		return errorSubscription(fmt.Errorf("request: %s", err.Error()))
	}

	requestBody := r.Body()
	if requestBody == nil {
		return errorSubscription(fmt.Errorf("parse body: %s", err.Error()))
	}

	ln := startServerOnPort(1234, p.h)
	defer ln.Close()

	url := ln.Addr().String()
	url = strings.Replace(url, "http://", "ws://", -1)
	if !strings.HasPrefix(url, "ws://") {
		url = "ws://" + url
	}

	headers := make(http.Header)
	r.Header.VisitAll(func(key, value []byte) {
		headers.Add(string(key), string(value))
	})

	c, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		return errorSubscription(fmt.Errorf("dial: %s", err.Error()))
	}

	initMessage := operationMessage{Type: connectionInitMsg}
	if initPayload != nil {
		initMessage.Payload, err = json.Marshal(initPayload)
		if err != nil {
			return errorSubscription(fmt.Errorf("parse payload: %s", err.Error()))
		}
	}

	if err = c.WriteJSON(initMessage); err != nil {
		return errorSubscription(fmt.Errorf("init: %s", err.Error()))
	}

	var ack operationMessage
	if err = c.ReadJSON(&ack); err != nil {
		return errorSubscription(fmt.Errorf("ack: %s", err.Error()))
	}

	if ack.Type != connectionAckMsg {
		return errorSubscription(fmt.Errorf("expected ack message, got %#v", ack))
	}

	var ka operationMessage
	if err = c.ReadJSON(&ka); err != nil {
		return errorSubscription(fmt.Errorf("ack: %s", err.Error()))
	}

	if ka.Type != connectionKaMsg {
		return errorSubscription(fmt.Errorf("expected ack message, got %#v", ack))
	}

	if err = c.WriteJSON(operationMessage{Type: startMsg, ID: "1", Payload: requestBody}); err != nil {
		return errorSubscription(fmt.Errorf("start: %s", err.Error()))
	}

	return &Subscription{
		Close: func() error {
			ln.Close()
			return c.Close()
		},
		Next: func(response interface{}) error {
			var op operationMessage
			err := c.ReadJSON(&op)
			if err != nil {
				return err
			}
			if op.Type != dataMsg {
				if op.Type == errorMsg {
					return fmt.Errorf(string(op.Payload))
				} else {
					return fmt.Errorf("expected data message, got %#v", op)
				}
			}

			var respDataRaw Response
			err = json.Unmarshal(op.Payload, &respDataRaw)
			if err != nil {
				return fmt.Errorf("decode: %s", err.Error())
			}

			// we want to unpack even if there is an error, so we can see partial responses
			unpackErr := unpack(respDataRaw.Data, response)

			if respDataRaw.Errors != nil {
				return RawJsonError{respDataRaw.Errors}
			}
			return unpackErr
		},
	}
}

func startServerOnPort(port int, h fasthttp.RequestHandler) net.Listener {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(fmt.Errorf("cannot start tcp server on port %d: %s", port, err))
	}
	go fasthttp.Serve(ln, h)
	return ln
}
