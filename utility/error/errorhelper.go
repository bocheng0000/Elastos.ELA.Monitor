package errorhelper

import (
	"github.com/elastos/Elastos.ELA.Monitor/utility/log"
)

func Panic(err error, message string) {
	if err != nil {
		log.Error(message)
		panic(err)
	}
}

func Warn(err error, message string) {
	if err != nil {
		log.Warnf("%s:\n %v", message, err)
	}
}

func WarnThenReturn(err error, message string) error {
	Warn(err, message)
	return err
}
