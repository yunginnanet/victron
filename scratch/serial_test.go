package cereal

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func Test_connectSerial(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("Must be root to run this test")
	}
	port, err := connectSerial("/dev/tty0")
	if err != nil {
		t.Fatal(err)
	}
	defer port.Close()
	t.Logf("Connected to serial port: %v", port)
}

func Test_getBaudRate(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("Must be root to run this test")
	}
	stty := exec.Command("stty", "-F", "/dev/tty0")
	out, err := stty.Output()
	if err != nil {
		t.Fatal(err)
	}
	split := strings.Split(string(out), " ")
	baud := split[1]
	t.Logf("Baud rate from `stty`: %v", baud)

	res, err := getBaudRate("/dev/tty0")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Baud rate from `getBaudRate`: %v", res)

	if baud != strconv.Itoa(res) {
		t.Fatalf("Baud rate mismatch: %v != %v", baud, res)
	}
}
