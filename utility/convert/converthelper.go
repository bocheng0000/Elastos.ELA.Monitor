package convert

import (
	"fmt"
	"github.com/elastos/Elastos.ELA.Monitor/utility/error"
	"strconv"
)

func StringToInt64(input string, defaultValue int64) int64 {
	value, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		errorhelper.Warn(err, fmt.Sprintf("convert %s to int64 failed!", input))
		return defaultValue
	}

	return value
}