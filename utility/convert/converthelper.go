package convert

import (
	"fmt"
	"github.com/elastos/Elastos.ELA.Monitor/utility/error"
	"strconv"
	"time"
)

func StringToInt64(input string, defaultValue int64) int64 {
	value, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		errorhelper.Warn(err, fmt.Sprintf("convert %s to int64 failed!", input))
		return defaultValue
	}

	return value
}

func StringToTime(layout, input string) time.Time {
	outputTime, err := time.Parse(layout, input[:len(layout)])
	if err != nil {
		errorhelper.Warn(err, fmt.Sprintf("convert %s to Time failed!", input))
		panic(err)
	}

	return outputTime
}