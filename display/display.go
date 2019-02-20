package display

import (
	"bytes"
	"fmt"
	"github.com/elastos/Elastos.ELA.Monitor/config"
	"github.com/elastos/Elastos.ELA.Monitor/logparse"
	"github.com/elastos/Elastos.ELA.Monitor/models"
	"github.com/elastos/Elastos.ELA.Monitor/nodes"
	"github.com/elastos/Elastos.ELA.Monitor/utility/convert"
	"github.com/elastos/Elastos.ELA.Monitor/utility/error"
	"github.com/elastos/Elastos.ELA.Monitor/utility/log"
	"path"
	"time"
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
	currentLogFile := logParse.ParseOldLogs(logData)
	logPath := path.Join(config.ConfigManager.MonitorConfig.Nodes.MainChain.LogPath, currentLogFile)

	for {
		time.Sleep(time.Duration(config.ConfigManager.MonitorConfig.Nodes.MainChain.LogFreshInterval) * time.Second)

		err := logParse.ReadLogFile(logData, logPath)
		errorhelper.Warn(err, "fresh log data failed!")

		display.initDisplay(logData, ela)
		display.Show()
	}
}

func (display *Display)Show() {
	log.Infof("-------------- %s --------------", config.ConfigManager.MonitorConfig.AppName)
	log.Infof("Version: %s \t Height: %d", display.Content.Version, display.Content.Height)
	log.Infof("Host: %s \t RpcPort: %d \t RestfulPort: %d", display.Content.Version, display.Content.RpcPort, display.Content.RestfulPort)

	log.Info("Network:")
	var hostsBuf bytes.Buffer
	for index := 0; index < len(*display.Networks); index ++ {
		hostsBuf.WriteString(fmt.Sprintf("%s ", (*display.Networks)[index]))
	}
	log.Info(hostsBuf.String())

	log.Info("View:")
	log.Infof("Change times: %d", display.View.ChangeTimes)
	log.Infof("On duty producer: %s \t %s", display.View.OnDutyProducer.NickName, display.View.OnDutyProducer.OwnerPublicKey)
	log.Infof("Proposal: Total: %s \t Approval: %s \t Decline: %s", len(display.View.Proposal.Approval), len(display.View.Proposal.Decline))
	log.Info("Producers:")
	for index := 0; index < len(*display.View.Producers); index ++ {
		log.Infof("NickName: %s \t NodePublicKey: %s \t Vote: %s \t IsActive: %v",
			(*display.View.Producers)[index].NickName,
			(*display.View.Producers)[index].NodePublicKey,
			(*display.View.Producers)[index].Vote,
			(*display.View.Producers)[index].IsActive)
	}

	log.Info("Evil:")
	for index := 0; index < len(*display.Evil.Producers); index ++ {
		log.Infof("NickName: %s \t NodePublicKey: %s \t Vote: %s \t IsActive: %v \t Evidence: %s",
			(*display.Evil.Producers)[index].NickName,
			(*display.Evil.Producers)[index].NodePublicKey,
			(*display.Evil.Producers)[index].Vote,
			(*display.Evil.Producers)[index].IsActive,
			(*display.Evil.Producers)[index].Evidence)
	}

	log.Info("NextView:")
	for index := 0; index < len(*display.NextView.Producers); index ++ {
		log.Infof("NickName: %s \t NodePublicKey: %s \t Vote: %s \t",
			(*display.NextView.Producers)[index].NickName,
			(*display.NextView.Producers)[index].NodePublicKey,
			(*display.NextView.Producers)[index].Vote)
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
	return logData.ConsensusStarted.Back().Value.(models.ConsensusMessage).LogTime
}

func (display *Display) initContent(logData *logparse.LogData, ela *nodes.Ela) *Content {
	height, _ := ela.Rpc.GetChainHeight()
	return &Content{logData.Version, height, ela.Host, ela.Rpc.Port, ela.Restful.Port}
}

func (display *Display) initNetworks(logData *logparse.LogData) *[]string {
	return logData.Network.Back().Value.(models.Network).NbrHosts
}

func (display *Display) initView(logData *logparse.LogData, ela *nodes.Ela, listProducers *models.ListProducersResponse) *View {
	currentView := logData.ViewStarted.Back().Value.(models.ViewStart)
	var onDutyProducer *ProducerMonitor
	var producerMonitors []ProducerMonitor
	for index := 0; index < len(listProducers.Producers); index ++ {
		producerMonitors[index] = ProducerMonitor {
			listProducers.Producers[index].NickName,
			listProducers.Producers[index].OwnerPublicKey,
			listProducers.Producers[index].NodePublicKey,
			0,
			convert.StringToInt64(listProducers.Producers[index].Votes, 0),
			listProducers.Producers[index].Active,
			"",
		}

		if listProducers.Producers[index].OwnerPublicKey == currentView.OnDutyArbitrator {
			onDutyProducer = &producerMonitors[index]
		}
	}

	var approvals, declines []*string
	for voteArrived := logData.VoteArrived.Front(); voteArrived != nil; voteArrived = voteArrived.Next() {
		vote := voteArrived.Value.(models.VoteArrivedMessage)
		if display.CurrentConsensusTime.After(vote.LogTime) {
			continue
		}

		if vote.Result {
			approvals = append(approvals, &vote.Signer)
		} else {
			declines = append(declines, &vote.Signer)
		}
	}

	proposal := Proposal {
		uint8(logData.VoteArrived.Len()),
		approvals,
		declines,
	}

	return &View{currentView.Offset,onDutyProducer,&proposal,&producerMonitors}
}

func (display *Display) initEvil(logData *logparse.LogData, ela *nodes.Ela, listProducers *models.ListProducersResponse) *Evil {
	var producerMonitors []ProducerMonitor
	producerMonitors[0] = ProducerMonitor {
		listProducers.Producers[0].NickName,
		listProducers.Producers[0].OwnerPublicKey,
		listProducers.Producers[0].NodePublicKey,
		0,
		convert.StringToInt64(listProducers.Producers[0].Votes, 0),
		listProducers.Producers[0].Active,
		"",
	}

	return &Evil{&producerMonitors}
}

func (display *Display) initNextView(logData *logparse.LogData, ela *nodes.Ela, listProducers *models.ListProducersResponse) *NextView {
	var producerMonitors []ProducerMonitor
	producerMonitors[0] = ProducerMonitor {
		listProducers.Producers[0].NickName,
		listProducers.Producers[0].OwnerPublicKey,
		listProducers.Producers[0].NodePublicKey,
		0,
		convert.StringToInt64(listProducers.Producers[0].Votes, 0),
		listProducers.Producers[0].Active,
		"",
	}
	return &NextView{&producerMonitors}
}