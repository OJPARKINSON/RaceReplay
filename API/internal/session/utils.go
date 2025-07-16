package session

import (
	"fmt"
	"math"
	"time"
)

func GetIntValue(val interface{}) int {
	if val == nil {
		return 0
	}

	if f, ok := val.(float64); ok {
		return int(math.Round(f))
	}

	if i, ok := val.(int64); ok {
		return int(i)
	}
	return 0
}

func GetFloatValue(val interface{}, decimalPlaces int) float64 {
	if val == nil {
		return 0.0
	}

	var result float64

	if f, ok := val.(float64); ok {
		result = f
	} else if i, ok := val.(int64); ok {
		result = float64(i)
	} else {
		return 0.0
	}

	multiplier := math.Pow(10, float64(decimalPlaces))
	return math.Round(result*multiplier) / multiplier
}

func GetPressureInBar(val interface{}, decimalPlaces int) float64 {
	kpa := GetFloatValue(val, decimalPlaces+2)
	bar := kpa / 100.0

	// Round to specified decimal places
	multiplier := math.Pow(10, float64(decimalPlaces))
	return math.Round(bar*multiplier) / multiplier
}

func GetTimeFormatted(val interface{}) string {
	seconds := GetFloatValue(val, 0)

	if seconds < 0 {
		seconds = 0
	}

	minutes := int(seconds) / 60
	remainingSeconds := int(seconds) % 60

	return fmt.Sprintf("%02d:%02d", minutes, remainingSeconds)
}

func GetTimeFormattedWithMillis(val interface{}) string {
	totalSeconds := GetFloatValue(val, 3)

	if totalSeconds < 0 {
		totalSeconds = 0
	}

	minutes := int(totalSeconds) / 60
	seconds := math.Mod(totalSeconds, 60)

	return fmt.Sprintf("%02d:%02f", minutes, seconds)
}

func GetDurationFromSeconds(val interface{}) time.Duration {
	seconds := GetFloatValue(val, 3)
	return time.Duration(seconds * float64(time.Second))
}

func GetIntSpeedValue(val interface{}) int {
	if val == nil {
		return 0
	}

	if f, ok := val.(float64); ok {
		return int(math.Round(f) * 2.237)
	}

	if i, ok := val.(int64); ok {
		return int(i)
	}
	return 0
}

type stringFormat int

const (
	CurrentLapTime = iota
	Delta
	LastLapTime
	BatteryVolts
	FuelUsed
	FuelLevel
	FuelLastLap
	TyrePressure
	TyreTemp
)

var FormatString = map[stringFormat]string{
	CurrentLapTime: "%02d:%02d:%06.3f", // 00:00:00.000
	Delta:          "%02d:%06.3f",      // 00:00.000
	LastLapTime:    "%02d:%02d:%06.3f", // 00:00:00.000
	BatteryVolts:   "%04.1f",           // 00.0
	FuelUsed:       "%.1f",             // 0.0
	FuelLevel:      "%04.1f",           // 00.0
	FuelLastLap:    "%.3f",             // 0.000
	TyrePressure:   "%.2f",             // 0.00
	TyreTemp:       "%02.0f",           // 00
}
