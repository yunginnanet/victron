package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/tarm/serial"

	"git.tcp.direct/kayos/cereal/pkg/victron"
	cereal "git.tcp.direct/kayos/cereal/scratch"
)

var (
	log    zerolog.Logger
	outLog zerolog.Logger
)

func init() {
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	outLog = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
}

func main() {
	entries, err := os.ReadDir("/dev/serial/by-id/")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list serial devices")
	}
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "Victron") {
			outLog.Info().Str("device", entry.Name()).Msg("found Victron device")
		}

		target := filepath.Join("/dev/serial/by-id/" + entry.Name())
		link := ""

		if link, err = os.Readlink(target); err == nil {
			if link, err = filepath.Abs(filepath.Join("/dev/serial/by-id/", link)); err == nil {
				target = link
			}
		}

		var sport io.ReadWriteCloser
		if sport, err = cereal.ConnectSerialBaud(target, 19200); err != nil {
			log.Error().Str("caller", entry.Name()).Err(err).Msg("failed to connect to serial port")
			continue
		}

		slog := log.With().Str("caller", entry.Name()).Logger()

		slog.Info().Msg("connected to serial port")

		stream := victron.NewStream(entry.Name(), sport.(*serial.Port))
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		var block *victron.Blocks
		if block, err = stream.ReadBlocks(ctx); err != nil {
			slog.Error().Err(err).Msg("failed to read blocks")
			cancel()
			continue
		}
		cancel()
		for key, value := range block.Fields {
			_, _ = os.Stdout.WriteString("\t" + key + ": " + value + "\n")
		}
	}
}
