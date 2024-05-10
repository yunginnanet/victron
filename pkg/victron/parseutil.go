package victron

/*
	Reference: https://www.victronenergy.com/upload/documents/VE.Direct-Protocol-3.33.pdf
	Pages: 6-8
*/

const (
	PrefixProductID          = "PID"
	PrefixProductDescription = "BMV" // Note: deprecated in newer firmware/models

	PrefixFirmwareVersion = "FWE"
	PrefixSerial          = "SER#"

	PrefixPVWattage      = "PPV"
	PrefixBatteryVoltage = "V"

	// PrefixStateOfCharge is the current state of charge.
	//  - When the BMV is not synchronised, these statistics have no meaning, so “---” will be sent instead of a value.
	//  - When device is configured as DC monitor, “---” will be sent instead of a value.
	PrefixStateOfCharge = "SOC"

	// PrefixTimeToGo is the time remaining until the battery is empty.
	//   - When the BMV is not synchronised, these statistics have no meaning, so “---” will be sent instead of a value.
	//   - When device is configured as DC monitor, “---” will be sent instead of a value.
	//   - When the battery is not discharging the time- to-go is infinite. This is represented as -1.
	PrefixTimeToGo = "TTG"

	PrefixAlarm       = "ALRM"
	PrefixAlarmReason = "AR"
	PrefixRelayState  = "Relay"
	PrefixOffReason   = "OR"

	// PrefixHistoryDeepestDischargeAh is the all time deepest discharge known.
	//   - When device is configured as DC monitor, “---” will be sent instead of a value.
	PrefixHistoryDeepestDischargeAh = "H1"

	// PrefixHistoryLastDischargeAh is the last discharge depth.
	//   - When device is configured as DC monitor, “---” will be sent instead of a value.
	PrefixHistoryLastDischargeAh = "H2"

	// PrefixHistoryAverageDischargeAh is the average discharge depth.
	//   - When device is configured as DC monitor, “---” will be sent instead of a value.
	PrefixHistoryAverageDischargeAh = "H3"

	// PrefixHistoryChargeCycles is the number of charge cycles.
	//   - When device is configured as DC monitor, “---” will be sent instead of a value.
	PrefixHistoryChargeCycles = "H4"

	// PrefixHistoryFullDischarges is the number of times the battery has been fully discharged.
	//   - When device is configured as DC monitor, “---” will be sent instead of a value.
	PrefixHistoryFullDischarges = "H5"

	// PrefixHistoryTotalAhDrawn is the total Ah drawn from the battery.
	//   - When device is configured as DC monitor, “---” will be sent instead of a value.
	PrefixHistoryTotalAhDrawn   = "H6"
	PrefixHistoryMinimumVoltage = "H7"
	PrefixHistoryMaximumVoltage = "H8"

	// PrefixHistorySecondsSinceFullCharge is the number of seconds since the last full charge.
	//   - When device is configured as DC monitor, “---” will be sent instead of a value.
	PrefixHistorySecondsSinceFullCharge = "H9"

	// PrefixHistoryAutomaticSyncCount is the number of times the BMV has been automatically synchronised.
	//   - When device is configured as DC monitor, “---” will be sent instead of a value.
	PrefixHistoryAutomaticSyncCount = "H10"

	PrefixHistoryLowVoltageAlarmCount     = "H11"
	PrefixHistoryHighVoltageAlarmCount    = "H12"
	PrefixHistoryLowAuxVoltageAlarmCount  = "H13"
	PrefixHistoryHighAuxVoltageAlarmCount = "H14"
	PrefixHistoryMinimumAuxVoltage        = "H15"
	PrefixHistoryMaximumAuxVoltage        = "H16"

	// PrefixHistoryTotalOutputKwh is the total Ah that has passed through the BMV, shunt, monitor, or battery outwards.
	// A.K.A: Discharged
	PrefixHistoryTotalOutputKwh = "H17"

	// PrefixHistoryTotalInputKwh is the total Ah that has passed through the BMV, shunt, monitor, or battery inwards.
	// A.K.A: Charged
	PrefixHistoryTotalInputKwh = "H18"

	PrefixHistoryUserResettableYieldKwh = "H19"
	PrefixHistoryTodaysYieldKwh         = "H20"
	PrefixHistoryMaxPowerTodayW         = "H21"
	PrefixHistoryYesterdaysYieldKwh     = "H22"
	PrefixHistoryMaxPowerYesterdayW     = "H23"

	PrefixErrorCode = "ERR"

	PrefixPanelVoltageMv      = "VPV"
	PrefixPanelPowerW         = "PPV"
	PrefixBatteryCurrentMa    = "I"
	PrefixLoadCurrentMa       = "IL"
	PrefixDeviceTemperatureC  = "T"
	PrefixInstantaneousPowerW = "P"
	PrefixConsumedAh          = "CE"

	PrefixAlarmCondition = "Alarm"

	PrefixMaximumPowerTodayW      = "H21"
	PrefixMaximumPowerYesterdayW  = "H23"
	PrefixErrorReason             = "ERR"
	PrefixWarningReason           = "WARN"
	PrefixDaySequenceNumber       = "HSDS"
	PrefixDeviceMode              = "MODE"
	PrefixTrackerOperationMode    = "MPPT"
	PrefixDCMonitorMode           = "MON"
	PrefixACOutputVoltageC        = "AC_OUT_V"
	PrefixACOutputCurrentA        = "AC_OUT_I"
	PrefixACOutputApparentPowerVA = "AC_OUT_S"
)
