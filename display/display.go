package display

import (
	"bytes"
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
	Networks    			*[]string
	View					*View
	Evil					*Evil
	NextView				*NextView
}

func NewDisplay(content *Content, networks *[]string, view *View, evil *Evil, nextView *NextView) *Display {
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
	var hostsBuf bytes.Buffer
	for index := 0; index < len(*display.Networks); index ++ {
		hostsBuf.WriteString(fmt.Sprintf("%s ", (*display.Networks)[index]))
	}
	log.Info(hostsBuf.String())

	log.Info("View:")
	if display.View == nil {
		log.Warn("No available view data")
	} else {
		log.Infof("Change times: %d", display.View.ChangeTimes)
		log.Infof("On duty producer: %s \t %s", display.View.OnDutyProducer.NickName, display.View.OnDutyProducer.OwnerPublicKey)
		log.Infof("Proposal: Total: %d \t Approval: %d \t Reject: %d", display.View.Proposal.Total, len(display.View.Proposal.Approval), len(display.View.Proposal.Reject))
		log.Info("Producers:")
		for index := 0; index < len(*display.View.Producers); index ++ {
			log.Infof("NickName: %s \t NodePublicKey: %s \t Vote: %d \t IsActive: %v",
				(*display.View.Producers)[index].NickName,
				(*display.View.Producers)[index].NodePublicKey,
				(*display.View.Producers)[index].Vote,
				(*display.View.Producers)[index].IsActive)
		}
	}


	log.Info("Evil:")
	if display.Evil == nil {
		log.Warn("No available evil data")
	} else {
		for index := 0; index < len(*display.Evil.Producers); index ++ {
			log.Infof("NickName: %s \t NodePublicKey: %s \t Vote: %d \t IsActive: %v \t Evidence: %s",
				(*display.Evil.Producers)[index].NickName,
				(*display.Evil.Producers)[index].NodePublicKey,
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
	errorhelper.Warn(err, "init display call failed!")

	display.CurrentConsensusTime = display.initCurrentConsensusTime(logData)
	display.Content = display.initContent(logData, ela)
	display.Networks = display.initNetworks(logData)
	display.View = display.initView(logData, ela, listProducers)
	display.Evil = display.initEvil(logData, ela, listProducers)
	display.NextView = display.initNextView(logData, ela, listProducers)
}

func (display *Display) initCurrentConsensusTime(logData *logparse.LogData) time.Time {
	if logData.ConsensusStarted.Len() == 0 {
		return time.Now().AddDate(0,0,-1)
	}

	return (*logData.ConsensusStarted.Back().Value.(*models.ConsensusMessage)).LogTime
}

func (display *Display) initContent(logData *logparse.LogData, ela *nodes.Ela) *Content {
	height, _ := ela.Rpc.GetChainHeight()
	virtualMemory, _ := mem.VirtualMemory()
	duration, _ := time.ParseDuration("1s")
	cpuPercent, _ := cpu.Percent(duration, false)
	return &Content{logData.Version, height, cpuPercent[0],virtualMemory.UsedPercent,ela.Host, ela.Rpc.Port, ela.Restful.Port}
}

func (display *Display) initNetworks(logData *logparse.LogData) *[]string {
	if logData.Network.Len() == 0 {
		return &[]string{"Network unavailable!"}
	}
	return (*logData.Network.Back().Value.(*models.Network)).NbrHosts
}

func (display *Display) initView(logData *logparse.LogData, ela *nodes.Ela, listProducers *models.ListProducersResponse) *View {
	if logData.ViewStarted.Len() == 0 {
		return nil
	}

	currentView := *logData.ViewStarted.Back().Value.(*models.ViewStart)
	var onDutyProducer *ProducerMonitor
	var producerMonitors []ProducerMonitor
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

		if listProducers.Producers[index].OwnerPublicKey == currentView.OnDutyArbitrator {
			onDutyProducer = &producerMonitors[index]
		}
	}

	if onDutyProducer == nil {
		onDutyProducer = &ProducerMonitor{
			"CRC Producer",
			currentView.OnDutyArbitrator,
			currentView.OnDutyArbitrator,
			0,
			0,
			true,
			"",
		}
	}

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

	return &View{currentView.Offset,onDutyProducer,&proposal,&producerMonitors}
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