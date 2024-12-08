package main

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	safeSpew "github.com/davecgh/go-spew/spew"
	"github.com/l0nax/go-spew/spew"
	"github.com/rs/zerolog"
	"github.com/tarm/serial"

	"git.tcp.direct/kayos/common/pool"

	"git.tcp.direct/kayos/cereal/pkg/victron"
	cereal "git.tcp.direct/kayos/cereal/scratch"

	"git.tcp.direct/tcp.direct/database/bitcask"
)

var (
	log  zerolog.Logger
	bufs = pool.NewBufferFactory()
	db   = bitcask.OpenDB(".victron_log")
)

func flags() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	for i, arg := range os.Args {
		switch arg {
		case "-v":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
		case "-vv":
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
		default:
			//
		}
	}
}

func hijackConsole() {
	// fix annoying use of fmt.Println in external libraries

	outR, outW, outE := os.Pipe()
	errR, errW, errE := os.Pipe()
	if outE != nil {
		panic(outE.Error())
	}
	if errE != nil {
		panic(errE.Error())
	}
	stdOut := bufio.NewScanner(outR)
	stdErr := bufio.NewScanner(errR)
	os.Stdout = outW
	os.Stderr = errW

	go func() {
		time.Sleep(10 * time.Millisecond)
		for {
			for !stdOut.Scan() && stdOut.Text() == "" {
				time.Sleep(5 * time.Millisecond)
			}
			log.Trace().Str("caller", "stdout").Msg(stdOut.Text())
		}
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		for {
			for !stdErr.Scan() && stdErr.Text() == "" {
				time.Sleep(5 * time.Millisecond)
			}
			log.Error().Str("caller", "stderr").Msg(stdErr.Text())
		}
	}()
}

func watchForOSSignals() chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	return sig
}

func mustCreateLogFile() io.Writer {
	targetFile := "./victron.log"
	if os.Getenv("VICTRON_LOG") != "" {
		targetFile = os.Getenv("VICTRON_LOG")
	}
	f, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	return f
}

func init() {
	realStdout := os.Stdout

	cw := zerolog.ConsoleWriter{Out: realStdout}
	f := mustCreateLogFile()
	mw := zerolog.MultiLevelWriter(cw, f)

	log = zerolog.New(mw).With().Timestamp().Logger()

	hijackConsole()
	flags()
	log.Trace().Msg("stdout and stderr redirected to logger")
}

func debugPrinter(prefix string, s string) {
	log.Debug().Str("caller", prefix).Msg(s)
}

func setup() []*victron.Stream {
	if len(os.Args) < 2 {
		log.Fatal().Msg("Must specify a serial port")
	}
	var devices []*victron.Stream
	for _, path := range os.Args[1:] {
		sport, err := cereal.ConnectSerialBaud(path, 19200)
		if err != nil {
			log.Fatal().Str("caller", path).Err(err).Msg("Failed to connect to serial port")
		}

		log.Info().Str("caller", path).Msgf("Connected to serial port: %v", sport)
		devices = append(devices, victron.NewStream(path, sport.(*serial.Port)).
			WithDebugPrinter(func(s string) { debugPrinter(path, s) }),
		)
	}
	return devices
}

var (
	SerialToDevice = make(map[string]*victron.Device)
	PathToDevice   = make(map[string]*victron.Device)
	DevicesMutex   = sync.RWMutex{}
)

var (
	debug *spew.ConfigState
)

func init() {
	debug = spew.NewDefaultConfig()
	debug.DisablePointerAddresses = true
	debug.DisableCapacities = true
	debug.DisableMethods = true
	debug.DisablePointerMethods = true
	debug.HighlightHex = true
	debug.HighlightValues = true
	debug.Indent = "\t"
	debug.MaxDepth = 2
}

func debugSdump(dev *victron.Device) (spewed string) {
	defer func() {
		if r := recover(); r != nil {
			println("[WARN] colored spew panic, falling back to upstream spew...")
			spewed = safeSpew.Sdump(dev.Status())
		}
	}()
	spewed = debug.Sdump(dev.Status())
	return
}

func stream(ctx context.Context, serialStream *victron.Stream) {
	slog := log.With().Str("caller", serialStream.Port()).Logger()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		block, err := serialStream.ReadBlocks(ctx)
		if err != nil {
			slog.Error().Err(err).Msg("Failed to read block")
			continue
		}
		if !block.Validate() {
			slog.Warn().Msg("Invalid block, sleeping...")
			time.Sleep(500 * time.Millisecond)
			continue
		}

		var dev *victron.Device
		var ok bool

		findDevice := func() (dev *victron.Device, ok bool) {
			DevicesMutex.RLock()
			defer DevicesMutex.RUnlock()
			dev, ok = PathToDevice[serialStream.Port()]
			if ok {
				return
			}

			blockSerial, serialFieldPresent := block.Fields[victron.PrefixSerial]

			if !serialFieldPresent || len(blockSerial) == 0 {
				slog.Trace().Msg("No serial field present in block")
				if zerolog.GlobalLevel() == zerolog.TraceLevel {
					debug.Dump(block)
				}
				return
			}

			dev, ok = SerialToDevice[blockSerial]
			if ok {
				DevicesMutex.RUnlock()
				DevicesMutex.Lock()
				PathToDevice[serialStream.Port()] = dev
				DevicesMutex.Unlock()
				DevicesMutex.RLock()
				return
			}
			DevicesMutex.RUnlock()
			DevicesMutex.Lock()
			if dev, err = victron.NewDevice(blockSerial); err != nil {
				slog.Panic().Err(err).Msg("Failed to create device timeseries")
			}
			dev.Serial = blockSerial
			SerialToDevice[blockSerial] = dev
			PathToDevice[serialStream.Port()] = dev
			DevicesMutex.Unlock()
			DevicesMutex.RLock()
			ok = true

			return
		}

		if dev, ok = serialStream.Device(); !ok {
			if _, wErr := os.Stderr.WriteString(victron.ErrDeviceNotAssociated.Error() + "\n"); wErr != nil {
				panic(wErr)
			}
			dev, ok = findDevice()
			if !ok {
				continue
			}
			serialStream.AssociateDevice(dev)
		}

		slog = log.With().Str("caller", dev.Serial).Logger()

		if err = dev.Update(block); err != nil {
			slog.Panic().Err(err).Msg("Failed to update device")
		}

		slog.Debug().Msg(debugSdump(dev))

	}
}

func main() {
	devices := setup()
	ctx, cancel := context.WithCancel(context.Background())
	for _, vdev := range devices {
		go stream(ctx, vdev)
	}
	caught := <-watchForOSSignals()
	cancel()
	log.Warn().Msgf("%s: caught signal, shutting down...", caught.String())
	for _, device := range devices {
		if err := device.Close(); err != nil {
			log.Error().Str("caller", device.Port()).Err(err).Msg("Failed to close device")
		}
	}
}
