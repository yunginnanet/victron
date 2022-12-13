package cereal

import (
	"io"

	"github.com/tarm/serial"
	"golang.org/x/sys/unix"
)

// getBaudRate enumerates a serial devices baud rate in linux using ioctl TCGETS.
func getBaudRate(name string) (int, error) {
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

func connectSerial(name string) (io.ReadWriteCloser, error) {
	res, err := getBaudRate(name)
	if err != nil {
		return nil, err
	}
	port, err := serial.OpenPort(&serial.Config{Name: name, Baud: res})
	if err != nil {
		return nil, err
	}
	return port, nil
}
