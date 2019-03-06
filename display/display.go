package display

import (
	"fmt"
	"github.com/elastos/Elastos.ELA.Monitor/utility/utility"
	"github.com/shirou/gopsutil/cpu"
	"path"
	"time"

	"github.com/elastos/Elastos.ELA.Monitor/config"
	"github.com/elastos/Elastos.ELA.Monitor/logparse"
	"github.com/elastos/Elastos.ELA.Monitor/models"
	"github.com/elastos/Elastos.ELA.Monitor/nodes"
	"github.com/elastos/Elastos.ELA.Monitor/utility/error"
	"github.com/elastos/Elastos.ELA.Monitor/utility/log"
	"github.com/shirou/gopsutil/mem"
)

type Display struct {
	CurrentConsensusTime	time.Time
	Content					*Content
	Networks    			*Network
	View					*View
	Evil					*Evil
	NextView				*NextView
}

func NewDisplay(content *Content, networks *Network, view *View, evil *Evil, nextView *NextView) *Display {
	return &Display{
		Content: content,
		Networks: networks,
		View: view,
		Evil: evil,
		NextView: nextView,
	}
}

func (display *Display)Start(logData *logparse.LogData, logParse *logparse.LogParse, ela *nodes.Ela) {
	currentNodeLogFile := logParse.ParseOldLogs(logData, "node")
	currentDposLogFile := logParse.ParseOldLogs(logData, "dpos")
	nodeLogPath := path.Join(config.ConfigManager.MonitorConfig.Nodes.MainChain.LogPath, "node", currentNodeLogFile)
	dposLogPath := path.Join(config.ConfigManager.MonitorConfig.Nodes.MainChain.LogPath, "dpos", currentDposLogFile)

	for {
		time.Sleep(time.Duration(config.ConfigManager.MonitorConfig.Nodes.MainChain.LogFreshInterval) * time.Second)

		err := logParse.ReadLogFile(logData, nodeLogPath)
		errorhelper.Warn(err, "fresh node log data failed!")

		err = logParse.ReadLogFile(logData, dposLogPath)
		errorhelper.Warn(err, "fresh dpos log data failed!")

		display.initDisplay(logData, ela)
		display.Show()
	}
}

func (display *Display)Show() {
	log.Infof("------------------------ %s ------------------------", config.ConfigManager.MonitorConfig.AppName)
	log.Infof("Version: %s \t Height: %d \t CPUUsed: %.2f%% \t MemoryUsed: %.2f%%",
		display.Content.Version, display.Content.Height, display.Content.CPUUsed, display.Content.MemoryUsed)
	log.Infof("Host: %s \t RpcPort: %d \t RestfulPort: %d", display.Content.Host, display.Content.RpcPort, display.Content.RestfulPort)

	log.Info("Network:")
	log.Infof("Total: %d \t TwoWayConnection: %d \t OutboundOnly: %d \t InboundOnly: %d \t NoneConnection: %d",
		display.Networks.Total,
		display.Networks.TwoWayConnection,
		display.Networks.InboundOnly,
		display.Networks.OutboundOnly,
		display.Networks.NoneConnection)

	log.Info("View:")
	if display.View == nil {
		log.Warn("No available view data")
	} else {
		log.Infof("Change times: %d", display.View.ChangeTimes)
		log.Infof("On duty producer: %s \t %s", display.View.OnDutyProducer.NickName, display.View.OnDutyProducer.OwnerPublicKey)
		log.Infof("Proposal: Total: %d \t Approval: %d \t Reject: %d", display.View.Proposal.Total, len(display.View.Proposal.Approval), len(display.View.Proposal.Reject))
		log.Info("Producers:")
		for index := 0; index < len(*display.View.Producers); index ++ {
			log.Infof("NickName: %s \t OwnerPublicKey: %s \t Vote: %d \t IsActive: %v",
				(*display.View.Producers)[index].NickName,
				(*display.View.Producers)[index].OwnerPublicKey,
				(*display.View.Producers)[index].Vote,
				(*display.View.Producers)[index].IsActive)
		}
	}

	log.Info("Evil:")
	if display.Evil == nil {
		log.Warn("No available evil data")
	} else {
		for index := 0; index < len(*display.Evil.Producers); index ++ {
			log.Infof("NickName: %s \t OwnerPublicKey: %s \t Vote: %d \t IsActive: %v \t Evidence: %s",
				(*display.Evil.Producers)[index].NickName,
				(*display.Evil.Producers)[index].OwnerPublicKey,
				(*display.Evil.Producers)[index].Vote,
				(*display.Evil.Producers)[index].IsActive,
				(*display.Evil.Producers)[index].Evidence)
		}
	}

	log.Info("NextView:")
	if display.NextView == nil {
		log.Warn("No available next view data")
	} else {
		for index := 0; index < len(*display.NextView.Producers); index ++ {
			log.Infof("NickName: %s \t NodePublicKey: %s \t Vote: %d \t",
				(*display.NextView.Producers)[index].NickName,
				(*display.NextView.Producers)[index].NodePublicKey,
				(*display.NextView.Producers)[index].Vote)
		}
	}
}

func (display *Display) initDisplay(logData *logparse.LogData, ela *nodes.Ela) {
	listProducers, err := ela.Rpc.GetListProducers(0, 100)
	errorhelper.Warn(err, "get list producers failed!")

	display.CurrentConsensusTime = display.initCurrentConsensusTime(logData)
	display.Content = display.initContent(logData, ela)
	display.Networks = display.initNetworks(logData, ela)
	display.View = display.initView(logData, ela, listProducers)
	display.Evil = display.initEvil(logData, ela, listProducers)
	display.NextView = display.initNextView(logData, ela, listProducers)
}

func (display *Display) initCurrentConsensusTime(logData *logparse.LogData) time.Time {
	if logData.ChangeView.Len() == 0 {
		return time.Now().AddDate(0,0,-1)
	}

	return *logData.ChangeView.Back().Value.(*time.Time)
}

func (display *Display) initContent(logData *logparse.LogData, ela *nodes.Ela) *Content {
	height, _ := ela.Rpc.GetChainHeight()
	virtualMemory, _ := mem.VirtualMemory()
	duration, _ := time.ParseDuration("1s")
	cpuPercent, _ := cpu.Percent(duration, false)
	return &Content{logData.Version, height, cpuPercent[0],virtualMemory.UsedPercent,ela.Host, ela.Rpc.Port, ela.Restful.Port}
}

func (display *Display) initNetworks(logData *logparse.LogData, ela *nodes.Ela) *Network {
	return &Network{0,0,0,0,0}
	//dposPeersInfos, err := ela.Rpc.GetDposPeersInfos()
	//if err != nil {
	//	return &netWork
	//}
	//
	//netWork.Total = uint8(len(*dposPeersInfos))
	//for _, dposPeersInfo := range *dposPeersInfos {
	//	switch dposPeersInfo.ConnectState {
	//	case "2WayConnection":
	//		netWork.TwoWayConnection ++
	//	case "OutboundOnly":
	//		netWork.OutboundOnly ++
	//	case "InboundOnly":
	//		netWork.InboundOnly ++
	//	case "NoneConnection":
	//		netWork.NoneConnection ++
	//	default:
	//	}
	//}
	//return &netWork
}

func (display *Display) initView(logData *logparse.LogData, ela *nodes.Ela, listProducers *models.ListProducersResponse) *View {
	if logData.ViewStarted.Len() == 0 {
		return nil
	}

	currentView := *logData.ViewStarted.Back().Value.(*models.ViewStart)
	producerMonitors := display.initProducerMonitors(listProducers)
	onDutyProducer := display.initOnDutyProducer(currentView.OnDutyArbitrator, producerMonitors)
	proposal := display.initProposalInfo(logData)

	return &View{currentView.Offset,onDutyProducer,proposal,producerMonitors}
}

func (display *Display) initProducerMonitors(listProducers *models.ListProducersResponse) *[]ProducerMonitor {
	producerMonitors := make([]ProducerMonitor, 0)
	crcProducers := &config.ConfigManager.MonitorConfig.Nodes.MainChain.CRCProducers

	for index, crcProducer := range *crcProducers {
		producerMonitors = append(producerMonitors, ProducerMonitor {
			fmt.Sprintf("CRCProducer%d", index),
			crcProducer,
			crcProducer,
			0,
			0,
			true,
			"",
		})
	}

	for index := 0; index < len(listProducers.Producers); index ++ {
		vote, _ := utility.ElaStringToSelaInt64(listProducers.Producers[index].Votes, 64)
		producerMonitors = append(producerMonitors, ProducerMonitor {
			listProducers.Producers[index].NickName,
			listProducers.Producers[index].OwnerPublicKey,
			listProducers.Producers[index].NodePublicKey,
			0,
			vote,
			listProducers.Producers[index].Active,
			"",
		})
	}

	return &producerMonitors
}

func (display *Display) initOnDutyProducer(onDutyArbitrator string, producerMonitors *[]ProducerMonitor) *ProducerMonitor {
	for _, producerMonitor := range *producerMonitors {
		if producerMonitor.OwnerPublicKey == onDutyArbitrator {
			return &producerMonitor
		}
	}

	return &ProducerMonitor{
		"CRC Producer",
		onDutyArbitrator,
		onDutyArbitrator,
		0,
		0,
		true,
		"",
	}
}

func (display *Display) initProposalInfo(logData *logparse.LogData) *Proposal {
	var approvals, rejects []*string
	for voteArrived := logData.VoteArrived.Front(); voteArrived != nil; voteArrived = voteArrived.Next() {
		vote := *voteArrived.Value.(*models.VoteArrivedMessage)
		if display.CurrentConsensusTime.After(vote.LogTime) {
			continue
		}

		if vote.Result {
			approvals = append(approvals, &vote.Signer)
		} else {
			rejects = append(rejects, &vote.Signer)
		}
	}

	proposal := Proposal {
		uint8(logData.VoteArrived.Len()),
		approvals,
		rejects,
	}

	return &proposal
}

func (display *Display) initEvil(logData *logparse.LogData, ela *nodes.Ela, listProducers *models.ListProducersResponse) *Evil {
	var producerMonitors []ProducerMonitor
	vote, _ := utility.ElaStringToSelaInt64(listProducers.Producers[0].Votes, 64)
	producerMonitors = append(producerMonitors, ProducerMonitor {
		listProducers.Producers[0].NickName,
		listProducers.Producers[0].OwnerPublicKey,
		listProducers.Producers[0].NodePublicKey,
		0,
		vote,
		listProducers.Producers[0].Active,
		"",
	})

	return &Evil{&producerMonitors}
}

func (display *Display) initNextView(logData *logparse.LogData, ela *nodes.Ela, listProducers *models.ListProducersResponse) *NextView {
	var producerMonitors []ProducerMonitor
	vote, _ := utility.ElaStringToSelaInt64(listProducers.Producers[0].Votes, 64)
	producerMonitors = append(producerMonitors, ProducerMonitor {
		listProducers.Producers[0].NickName,
		listProducers.Producers[0].OwnerPublicKey,
		listProducers.Producers[0].NodePublicKey,
		0,
		vote,
		listProducers.Producers[0].Active,
		"",
	})

	return &NextView{&producerMonitors}
}