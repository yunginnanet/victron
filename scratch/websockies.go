package cereal

// build:linux

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/websocket"
	"github.com/lxc/lxd/shared"
	// "github.com/tarm/serial"
)

func concept(w http.ResponseWriter, r *http.Request) {
	// log.Info().Msgf("Got a request:\n %v", r)
	up := websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		// Subprotocols:      nil,
		// Error:             nil,
		// CheckOrigin:       nil,
		EnableCompression: true,
	}
	wsc, err := up.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade connection")
		return
	}
	plumbing := &shared.WebsocketIO{Conn: wsc}
	defer plumbing.Close()
	_, err = plumbing.Write([]byte("Hello, world!"))
}

func Listen(addrs ...string) error {
	http.HandleFunc("/", concept)
	addr := "localhost:8080"
	if len(addrs) > 0 {
		addr = addrs[0]
	}
	return http.ListenAndServe(addr, nil)
}
