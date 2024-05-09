package main

import (
	"os"
	"sync"
	"time"

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

func setup() []*victron.Stream {
	if len(os.Args) < 2 {
		log.Fatal().Msg("Must specify a serial port")
	}
	var devices []*victron.Stream
	for _, path := range os.Args[1:] {
		portal, err := cereal.ConnectSerialBaud(path, 19200)
		if err != nil {
			log.Fatal().Str("caller", path).Err(err).Msg("Failed to connect to serial port")
		}

		log.Info().Str("caller", path).Msgf("Connected to serial port: %v", portal)
		devices = append(devices, victron.NewStream(path, portal.(*serial.Port)))
	}
	return devices
}

var (
	Devices      = make(map[string]*victron.Device)
	DevicesMutex = sync.RWMutex{}
)

func main() {
	devices := setup()
	quit := make(chan struct{})
	for _, vdev := range devices {
		go func(device *victron.Stream) {
			for {
				block, _ := device.ReadBlock()
				if !block.Validate() {
					log.Warn().Msg("Invalid block, sleeping...")
					time.Sleep(500 * time.Millisecond)
					continue
				}
				blockDeviceSerial, ok := block.Fields[victron.PrefixSerial]
				if !ok {
					// FIXME, there is a second block with the serial number that shpuld be handled
					// log.Warn().Msg("No serial number in block, skipping...")
					// fmt.Printf("%+v\n", block)
					continue
				}
				DevicesMutex.RLock()
				blockDevice, ok := Devices[blockDeviceSerial]
				DevicesMutex.RUnlock()
				if !ok {
					var err error
					blockDevice, err = victron.NewDevice()
					if err != nil {
						panic(err.Error())
					}
					DevicesMutex.Lock()
					Devices[blockDeviceSerial] = blockDevice
					DevicesMutex.Unlock()
				}
				if err := blockDevice.Update(block); err != nil {
					log.Panic().Err(err).Msg("Failed to update device")
				}
				// log.Info().Msgf("Updated device: %v", blockDevice)
			}
		}(vdev)
	}
	<-quit

}
