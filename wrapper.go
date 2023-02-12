package streamwebsocket

import (
	"encoding/json"
	"sync"
)

type Conn interface {
	WriteMessage(int, []byte) error
	Close() error
}

type WebsocketWrapper struct {
	Conn Conn
	mu   *sync.Mutex
}

func NewWrapper(c Conn) *WebsocketWrapper {
	return &WebsocketWrapper{
		Conn: c,
		mu:   &sync.Mutex{},
	}
}

func (w *WebsocketWrapper) WriteSafe(mt int, data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Conn.WriteMessage(mt, data)
}

func (w *WebsocketWrapper) WriteSafeJSON(data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return w.WriteSafe(TextMessage, j)
}

func (w *WebsocketWrapper) CloseSafe(data []byte) error {
	w.WriteSafe(CloseMessage, data)
	return w.Conn.Close()
}
