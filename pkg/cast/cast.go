package cast

import (
	"fmt"
	"strconv"
)

func StringToInt(val string, def int) (int, error) {
	if val == "" {
		return def, nil
	}

	castValue, err := strconv.Atoi(val)
	if err != nil {
		return def, fmt.Errorf("WARN: Failed to cast string %q to int, falling back to default: %d", val, def)
	}

	return castValue, nil
}

func IntToString(val int) string {
	return strconv.Itoa(val)
}

func MBToBytes(val int) int {
	return 1024 * 1024 * val
}
