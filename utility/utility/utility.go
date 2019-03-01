package utility

import (
	"fmt"
	"github.com/elastos/Elastos.ELA.Monitor/utility/constants"
	"github.com/elastos/Elastos.ELA.Monitor/utility/error"
	"strconv"
)

func ElaStringToSelaInt64(input string, bitSize int) (int64, error) {
	value, err := strconv.ParseFloat(input, bitSize)
	value = value * constants.ElaToSelaRate
	if err != nil {
		errorhelper.Warn(err, fmt.Sprintf("convert ela string %s to sela int64 failed!", input))
		return 0, err
	}

	return int64(value), err
}
