package main

import (
	"bufio"
	"bytes"
	"os"
	//	"strings"

	"github.com/rs/zerolog"

	"git.tcp.direct/kayos/common/pool"

	cereal "git.tcp.direct/kayos/cereal/scratch"

	"git.tcp.direct/tcp.direct/database/bitcask"
)

var (
	log  zerolog.Logger
	bufs = pool.NewBufferFactory()
	db   = bitcask.OpenDB(".victron_log")
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	for _, arg := range os.Args {
		switch arg {
		case "-v":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "-vv":
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		default:
			//
		}
	}
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
}

type MPPT struct {
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal().Msg("Must specify a serial port")
	}
	portal, err := cereal.ConnectSerialBaud(os.Args[1], 19200)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to serial port")
	}
	defer portal.Close()
	log.Info().Msgf("Connected to serial port: %v", portal)
	xerox := bufio.NewScanner(portal)

	logFile, err := os.Create("victron.log")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create log file")
	}
	defer logFile.Close()

	for xerox.Scan() {
		lineBytes := xerox.Bytes()
		_, _ = logFile.Write(lineBytes)
		_, _ = logFile.Write([]byte("\n"))
		line := string(bytes.TrimSpace(xerox.Bytes()))
		log.Trace().Msgf("Got a line: %s", string(line))
		switch {
		case len(bytes.TrimSpace(lineBytes)) == 0:
			//
		default:

		}
	}
}
