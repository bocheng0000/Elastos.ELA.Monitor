package main

import (
	"github.com/elastos/Elastos.ELA.Monitor/config"
	"github.com/elastos/Elastos.ELA.Monitor/display"
	"github.com/elastos/Elastos.ELA.Monitor/logparse"
	"github.com/elastos/Elastos.ELA.Monitor/nodes"
	"github.com/elastos/Elastos.ELA.Monitor/utility/log"
)

func init() {
	log.Init(
		config.ConfigManager.MonitorConfig.Log.Path,
		config.ConfigManager.MonitorConfig.Log.Level,
		config.ConfigManager.MonitorConfig.Log.MaxPerLogSizeMb,
		config.ConfigManager.MonitorConfig.Log.MaxLogsSizeMb,
	)
}

func main() {
	log.Info("Welcome To Elastos ELA Monitor")

	//err := email.SendMonitorEMail("Elastos ELA Monitor Notify Email", "Hello Wen.Zhang")
	//errorhelper.Warn(err, "send email test failed!")

	mainChain := config.ConfigManager.MonitorConfig.Nodes.MainChain

	logData := logparse.NewLogData()
	logParse := logparse.NewLogParse()
	elaNode := nodes.NewEla(mainChain)
	displayer := display.NewDisplay(nil, nil, nil, nil, nil)
	displayer.Start(logData, logParse, elaNode)

	//monitors.ProducerMonitor.Start(elaNode)
}
