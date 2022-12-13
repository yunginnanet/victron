package cereal

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// connectToWebsocketServer is a test helper function that connects to a websocket server.
func connectToWebsocketServer(addr string, t *testing.T) *websocket.Conn {
	t.Helper()
	conn, resp, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		t.Fatalf("Failed to connect to websocket server: %v", err)
	}
	if resp.StatusCode != 101 {
		t.Fatalf("Expected status code 101, got %v", resp.StatusCode)
	}
	conn.SetReadLimit(1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	return conn
}

func TestListen(t *testing.T) {
	go func() {
		if err := Listen(); err != nil {
			t.Errorf("Failed to listen: %v", err)
		}
	}()
	time.Sleep(150 * time.Millisecond)
	conn := connectToWebsocketServer("ws://localhost:8080", t)
	if conn == nil {
		t.Fatal("Failed to connect to websocket server")
	}
	defer conn.Close()

	mt, m, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}
	if mt != websocket.BinaryMessage {
		t.Fatalf("Expected text message, got %v", mt)
	}
	if string(m) != "Hello, world!" {
		t.Fatalf("Expected 'Hello, world!', got %v", string(m))
	}

}
