package util

import "testing"

func TestParseBatteryVoltage(t *testing.T) {
	voltage, err := ParseBatteryVoltage("\t25450")
	if err != nil {
		t.Errorf("ParseBatteryVoltage returned an error: %v", err)
	}
	if voltage != 25.45 {
		t.Errorf("ParseBatteryVoltage returned %f, expected 25.45", voltage)
	}

	// Test invalid input
	_, err = ParseBatteryVoltage("Invalid input")
	if err == nil {
		t.Error("ParseBatteryVoltage did not return an error for invalid input")
	}
}
