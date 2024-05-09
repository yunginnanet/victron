package util

import (
	"strconv"
	"strings"
)

func ParseBatteryVoltage(in string) (float32, error) {
	// Example line from the Device that translates to 25.45v: "V	25450"
	in = strings.TrimSpace(in)
	voltage, err := strconv.ParseFloat(in, 32)
	if err != nil {
		return 0, err
	}
	return float32(voltage) / 1000, nil
}
