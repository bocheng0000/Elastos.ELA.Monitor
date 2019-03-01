package servers

import (
	"fmt"
)

type Restful struct {
	Host string
	Port uint16
	Url string
}

func NewRestful(host string, port uint16) *Restful {
	url := fmt.Sprintf("http://%s:%d", host, port)
	return &Restful{host, port, url}
}