package http

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/elastos/Elastos.ELA.Monitor/utility/log"
)

func Post(url, data string) ([]byte, error){
	response, err := http.Post(url, "application/json", strings.NewReader(data))
	if err != nil {
		log.Errorf("http post failed: %+v", err)
		log.Errorf("http url is %s", url)
		log.Errorf("data is %+v", data)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("read response failed! %+v", err)
		log.Errorf("response is %+v", response)
	}

	return body, err
}
