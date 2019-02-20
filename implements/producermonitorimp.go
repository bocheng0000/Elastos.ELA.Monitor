package implements

import (
	"github.com/elastos/Elastos.ELA.Monitor/servers"
	"github.com/elastos/Elastos.ELA.Monitor/utility/log"
)

var ProducerMonitorImp *producerMonitorImplements

type producerMonitorImplements struct {
}

func (pmi *producerMonitorImplements) Test(rpc *servers.Rpc) {
	 listProducers, err := rpc.GetListProducers(0, 100)
	 if err != nil{
	 	log.Errorf("call failed! %+v", err)
	 }

	log.Infof("producers is %+v", listProducers.Producers[0].NickName)
}

func (pmi *producerMonitorImplements) ReadBlock(rpc *servers.Rpc) {
	//block, err := rpc.GetBlockByHeight(100)
	//if err != nil{
	//	log.Errorf("call failed! %+v", err)
	//}

	height, err := rpc.GetChainHeight()
	if err != nil{
		log.Errorf("call failed! %+v", err)
	}

	log.Infof("block height is %+v", height)
}