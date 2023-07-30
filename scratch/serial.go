package cereal

import (
	"io"
	"os"

	"github.com/tarm/serial"
	"golang.org/x/sys/unix"
)

// GetBaudRate enumerates a serial devices baud rate in linux using ioctl TCGETS.
func GetBaudRate(name string) (int, error) {
	// modes are a reproduction of what i found when running `strace stty -F /dev/tty0`
	fd, err := unix.Open(name, unix.O_RDONLY|unix.O_NONBLOCK|unix.O_LARGEFILE, 0)
	if err != nil {
		return 0, err
	}
	term, err := unix.IoctlGetTermios(fd, unix.TCGETS2)
	_ = unix.Close(fd)
	if err != nil {
		return 0, err
	}
	return int(term.Ispeed), nil
}

// newSerial creates a new virtual serial device.
func newSerial(name string, termies unix.Termios) (io.ReadWriteCloser, error) {
	fd, err := unix.Open(name, unix.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK|unix.O_LARGEFILE, 0666)
	if err != nil {
		return nil, err
	}

	term, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil, err
	}

	term.Iflag = termies.Iflag
	term.Oflag = termies.Oflag
	term.Cflag = termies.Cflag
	term.Lflag = termies.Lflag
	term.Line = termies.Line
	term.Ispeed = termies.Ispeed
	term.Ospeed = termies.Ospeed
	term.Cc = termies.Cc

	if err := unix.IoctlSetTermios(fd, unix.TCSETS, term); err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(fd), name), nil
}

func ConnectSerial(name string) (io.ReadWriteCloser, error) {
	res, err := GetBaudRate(name)
	if err != nil {
		return nil, err
	}
	port, err := serial.OpenPort(&serial.Config{Name: name, Baud: res})
	if err != nil {
		return nil, err
	}
	return port, nil
}

func ConnectSerialBaud(name string, baud int) (io.ReadWriteCloser, error) {
	port, err := serial.OpenPort(&serial.Config{Name: name, Baud: baud})
	if err != nil {
		return nil, err
	}
	return port, nil
}
