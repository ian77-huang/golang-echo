package cast

import (
	"log"
	"strconv"
)

func Int(val string, def int) int {
	if val == "" {
		return def
	}

	castValue, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("WARN: Failed to cast string %q to int, falling back to default: %d", val, def)
		return def
	}

	return castValue
}
