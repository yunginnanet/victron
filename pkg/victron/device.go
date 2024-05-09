package victron

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/tarm/serial"

	"github.com/rosenstand/go-vedirect/vedirect"

	"git.tcp.direct/kayos/cereal/pkg/timeseries"
	"git.tcp.direct/kayos/cereal/pkg/util"
)

func NewStream(serialDevice string, port *serial.Port) *Stream {
	return &Stream{Stream: &vedirect.Stream{Device: serialDevice, Port: port}}
}

type Stream struct {
	*vedirect.Stream
}

type Block struct {
	*vedirect.Block
	Fields map[string]string
}

func (s *Stream) ReadBlock() (Block, error) {
	block, n := s.Stream.ReadBlock()
	if n == 0 || !block.Validate() {
		return Block{}, fmt.Errorf("%w: no block available", io.EOF)
	}
	reflected := reflect.ValueOf(&block).Elem().FieldByName("fields")
	if !reflected.IsValid() {
		return Block{}, errors.New("reflection failure: no fields in block")
	}
	if reflected.Kind() != reflect.Map {
		return Block{}, errors.New("reflection failure: fields is not a map")
	}
	fields := make(map[string]string)
	for _, key := range reflected.MapKeys() {
		fields[key.String()] = reflected.MapIndex(key).String()
	}

	return Block{Block: &block, Fields: fields}, nil
}

func NewDevice() (*Device, error) {
	ts, err := timeseries.New(".data/victron")
	if err != nil {
		return nil, err
	}
	return &Device{RWMutex: &sync.RWMutex{}, statusLog: ts}, nil
}

type Device struct {
	Product  string
	Serial   string
	Firmware string

	*sync.RWMutex
	lastStatus Status
	statusLog  *timeseries.TimeSeries
}

func (d *Device) GetAllHistory() []Status {
	d.RLock()
	defer d.RUnlock()
	var history []Status
	for ts, status := range d.statusLog.Items() {
		dat, err := sonic.Marshal(status)
		if err != nil {
			_, _ = os.Stderr.WriteString(err.Error())
			continue
		}
		entry := Status{}
		if err = sonic.Unmarshal(dat, &entry); err != nil {
			_, _ = os.Stderr.WriteString(err.Error())
			continue
		}
		entry.ts = ts
		history = append(history, entry)
	}
	return history
}

func (d *Device) GetLastStatus() Status {
	d.RLock()
	defer d.RUnlock()
	return d.lastStatus
}

func (d *Device) Update(block Block) error {
	d.Lock()
	status := Status{Data: make(map[string]any), ts: time.Now()}

	if len(block.Fields) == 0 {
		reflected := reflect.ValueOf(&block).FieldByName("fields")
		if !reflected.IsValid() {
			return fmt.Errorf("no fields in block")
		}
		if reflected.Kind() != reflect.Map {
			return fmt.Errorf("fields is not a map")
		}
		fields := reflected.Interface().(map[string]string)
		block.Fields = fields
	}

	for key, value := range block.Fields {
		switch key {
		case PrefixBatteryVoltage:
			var err error
			if status.BatteryVoltage, err = util.ParseBatteryVoltage(value); err != nil {
				return fmt.Errorf("failed to parse battery voltage: %w", err)
			}
			status.Data[key] = status.BatteryVoltage
		case PrefixPVWattage:
			var err error
			if status.PVWattage, err = strconv.Atoi(value); err != nil {
				return fmt.Errorf("failed to parse PV wattage: %w", err)
			}
			status.Data[key] = status.PVWattage
		default:
			status.Data[key] = value
		}
	}
	d.lastStatus = status
	err := d.statusLog.IngestVictron(status)
	d.Unlock()
	return err
}

type Status struct {
	ts             time.Time      // key in timeseries, so not needed in json
	BatteryVoltage float32        `json:"batt_voltage"`
	Data           map[string]any `json:"fields"`
	PVWattage      int            `json:"pv_wattage"`
	// BatteryCurrent
	// SolarVoltage
	// SolarAmperage
	// ErrorState
	// Relay State
}

func (s Status) Timestamp() time.Time {
	return s.ts
}

func (s Status) Fields() map[string]any {
	return s.Data
}

const (
	PrefixProductID      = "PID"
	PrefixFirmware       = "FWE"
	PrefixSerial         = "SER#"
	PrefixBatteryVoltage = "V"
	PrefixPVWattage      = "PPV"
)
