package monitors

import (
	"github.com/elastos/Elastos.ELA.Monitor/implements"
	"github.com/elastos/Elastos.ELA.Monitor/nodes"
)

var ProducerMonitor *producerMonitor

type producerMonitor struct {
	Name string
}

func (producerMonitor *producerMonitor) Start(node *nodes.Ela) {
	//implements.ProducerMonitorImp.Test(node.Rpc)
	implements.ProducerMonitorImp.ReadBlock(node.Rpc)
}