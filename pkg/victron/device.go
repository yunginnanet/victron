package victron

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/bytedance/sonic"

	"git.tcp.direct/kayos/cereal/pkg/timeseries"
	"git.tcp.direct/kayos/cereal/pkg/util"
)

func NewDevice(serial string) (*Device, error) {
	ts, err := timeseries.New(filepath.Join(".data", "victron", serial))
	if err != nil {
		return nil, err
	}
	return &Device{statusLog: ts}, nil
}

func (d *Device) AssociateStream(s *Stream) {
	d.Lock()
	d.c = s
	s.AssociateDevice(d)
	d.Unlock()
}

type Device struct {
	Product  string
	Serial   string
	Firmware string

	c *Stream
	sync.RWMutex
	lastStatus   Status
	statusLog    *timeseries.TimeSeries
	debugPrinter func(string)
}

func (d *Device) Close() error {
	errs := make([]error, 0, 2)
	if d.c != nil {
		errs = append(errs, d.c.dev.Close())
	}
	if d.statusLog != nil {
		errs = append(errs, d.statusLog.Close())
	}
	return errors.Join(errs...)
}

func (d *Device) History() ([]Status, error) {
	d.RLock()
	defer d.RUnlock()
	var history []Status
	var errs []error
	for ts, status := range d.statusLog.Items() {
		// Marshal and unmarshal to ensure that the data fits in the struct.
		dat, err := sonic.Marshal(status)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		entry := Status{}
		if err = sonic.Unmarshal(dat, &entry); err != nil {
			errs = append(errs, err)
			continue
		}
		entry.ts = ts
		history = append(history, entry)
	}
	return history, errors.Join(errs...)
}

// Status returns the last known status of the device.
func (d *Device) Status() Status {
	d.RLock()
	defer d.RUnlock()
	return d.lastStatus
}

func (d *Device) Update(block *Blocks) error {
	if !block.Validate() {
		return errors.New("corrupt/incomplete blocks in update")

	}
	d.Lock()
	defer d.Unlock()
	status := Status{Data: make(map[string]any), ts: time.Now()}

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

	return err
}
