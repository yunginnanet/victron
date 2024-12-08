package victron

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/l0nax/go-spew/spew" // use normal package because debugging callers probably don't want color
	"github.com/rosenstand/go-vedirect/vedirect"
	"github.com/tarm/serial"
)

func NewStream(serialDevice string, port *serial.Port) *Stream {
	return &Stream{
		Stream: &vedirect.Stream{Device: serialDevice, Port: port},
	}
}

func (s *Stream) WithDebugPrinter(f func(string)) *Stream {
	s.mu.Lock()
	s.debugPrinter = f
	s.mu.Unlock()
	return s
}

type Stream struct {
	*vedirect.Stream
	dev          *Device
	mu           sync.RWMutex
	debugPrinter func(string)
}

func (s *Stream) DebugPrintln(msg string) {
	s.mu.RLock()
	dp := s.debugPrinter
	s.mu.RUnlock()
	if dp != nil {
		dp(msg)
	}
}

// AssociateDevice associates a device with a stream for control flow purposes.
func (s *Stream) AssociateDevice(d *Device) {
	s.mu.Lock()
	s.dev = d
	s.mu.Unlock()
}

func (s *Stream) Device() (*Device, bool) {
	s.mu.Lock()
	d := s.dev
	s.mu.Unlock()
	if d == nil {
		return nil, false
	}
	return d, true
}

func (s *Stream) Port() string {
	s.mu.RLock()
	p := s.Stream.Device
	s.mu.RUnlock()
	return p
}

var ErrDeviceNotAssociated = errors.New("device not associated with stream")

func (s *Stream) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.dev == nil {
		return ErrDeviceNotAssociated
	}
	return s.dev.Close()
}

func (s *Stream) ReadBlocks(ctx context.Context) (*Blocks, error) {
	combined := newBlocks()

	for {
		select {
		case <-ctx.Done():
			return combined, ctx.Err()
		default:
		}

		err := combined.readBlock(s.Stream)
		if errors.Is(err, ErrBadChecksumModulus) {
			s.DebugPrintln("bad checksum modulus, skipping block: ")
			s.DebugPrintln(spew.Sdump(combined.lastBlock))
			continue
		}
		if err != nil {
			s.DebugPrintln("failed to read block: " + err.Error())
			return combined, err
		}

		if combined.Validate() {
			break
		}

		s.DebugPrintln("failed to validate blocks: ")
		s.DebugPrintln(spew.Sdump(combined))
		s.DebugPrintln("dropping invalid blocks...")
		if n := combined.DropInvalid(); n > 0 {
			s.DebugPrintln(fmt.Sprintf("dropped %d invalid blocks", n))
		} else if n == 0 {
			s.DebugPrintln("no invalid blocks to drop...")
			continue
		}
		if combined.Validate() {
			break
		}
	}

	return combined, nil
}
