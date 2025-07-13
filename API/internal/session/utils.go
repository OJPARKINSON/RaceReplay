package session

import "math"

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

// Add options for how many decimal places
func GetFloatValue(val interface{}) float64 {
	if val == nil {
		return 0.0
	}

	if f, ok := val.(float64); ok {
		return math.Round(f*10) / 10
	}

	if i, ok := val.(int64); ok {
		return float64(i)
	}
	return 0.0
}
