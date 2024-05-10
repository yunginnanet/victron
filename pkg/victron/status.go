package victron

import "time"

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
