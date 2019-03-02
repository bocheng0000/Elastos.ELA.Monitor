package servers

import (
	"encoding/json"
	"fmt"
	"github.com/elastos/Elastos.ELA.Monitor/models"
	"github.com/elastos/Elastos.ELA.Monitor/utility/error"
	"github.com/elastos/Elastos.ELA.Monitor/utility/http"
	"github.com/elastos/Elastos.ELA.Monitor/utility/log"
	"github.com/goinggo/mapstructure"
)

type Rpc struct {
	Host string
	Port uint16
	Url string
}

func NewRpc(host string, port uint16) *Rpc {
	url := fmt.Sprintf("http://%s:%d", host, port)
	return &Rpc{host, port, url}
}

func (rpc *Rpc) GetDposPeersInfos() (*[]models.DposPeersInfo, error) {
	response, err :=rpc.callAndReadRpc("getdpospeersinfo", nil)
	if err != nil {
		return nil, err
	}

	dposPeersInfos := []models.DposPeersInfo{}
	err = mapstructure.Decode(response.Result, dposPeersInfos)
	errorhelper.Warn(err, "decode block failed!")

	return &dposPeersInfos, err
}

func (rpc *Rpc) GetChainHeight() (uint32, error) {
	count, err := rpc.GetBlockCount()
	if err != nil {
		return 0, err
	}

	return count - 1, err
}

func (rpc *Rpc) GetBlockCount() (uint32, error) {
	response, err :=rpc.callAndReadRpc("getblockcount", nil)
	if err != nil {
		return 0, err
	}

	return uint32(response.Result.(float64)), err
}

func (rpc *Rpc) GetBlockByHeight(height uint32) (*models.Block, error) {
	data := models.Height{height}
	response, err :=rpc.callAndReadRpc("getblockbyheight", data)
	if err != nil {
		return nil, err
	}

	block := models.Block{}
	err = mapstructure.Decode(response.Result, &block)
	errorhelper.Warn(err, "decode block failed!")

	return &block, err
}

func (rpc *Rpc) GetListProducers(start, limit uint16) (*models.ListProducersResponse, error) {
	data := models.ListProducers{start, limit}
	response, err :=rpc.callAndReadRpc("listproducers", data)
	if err != nil {
		return nil, err
	}

	listProducersResponse := models.ListProducersResponse{}
	err = mapstructure.Decode(response.Result, &listProducersResponse)
	errorhelper.Warn(err, "decode list producer response failed!")

	return &listProducersResponse, err
}

func (rpc *Rpc) rpcPost(url, method string, params interface{}) ([]byte, error) {
	httpData := &models.HttpData{method, params}
	data, err := json.Marshal(httpData)
	if err != nil {
		log.Warnf("json marshal failed: %+v", data)
		log.Warnf("Error: %+v", err)
		return nil, err
	}

	sendData := string(data)
	log.Debug(fmt.Sprintf("call %s with %s", method, sendData))
	return http.Post(url, sendData)
}

func (rpc *Rpc) callAndReadRpc(method string, params interface{}) (*models.RpcResponse, error) {
	response, err := rpc.rpcPost(rpc.Url, method, params)
	if err != nil {
		return nil, err
	}

	log.Debugf("response is %+v", string(response))
	rpcResponse := &models.RpcResponse{}
	err = json.Unmarshal(response, rpcResponse)
	if err != nil {
		log.Errorf("Unmarshal json file err %v", err)
		log.Errorf("respone is %v", string(response))
	}
	return rpcResponse, err
}