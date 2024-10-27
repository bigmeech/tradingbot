package clients

import (
	"github.com/gorilla/websocket"
	"net/url"
)

type WebSocketClient struct {
	conn *websocket.Conn
}

func NewWebSocketClient(endpoint string) (*WebSocketClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return &WebSocketClient{conn: conn}, nil
}

func (ws *WebSocketClient) ReadMessage() ([]byte, error) {
	_, message, err := ws.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (ws *WebSocketClient) WriteMessage(messageType int, data []byte) error {
	return ws.conn.WriteMessage(messageType, data)
}

func (ws *WebSocketClient) Close() error {
	return ws.conn.Close()
}
