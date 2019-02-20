package logparse

import (
	"github.com/elastos/Elastos.ELA.Monitor/config"
	"github.com/elastos/Elastos.ELA.Monitor/models"
	"github.com/elastos/Elastos.ELA.Monitor/utility/constants"
	"github.com/elastos/Elastos.ELA.Monitor/utility/error"
	"github.com/elastos/Elastos.ELA.Monitor/utility/file"
	"github.com/elastos/Elastos.ELA.Monitor/utility/log"
	"path"
	"strconv"
	"strings"
	"time"
)

//var LogParser *LogParse

type LogParse struct {
	FileName 	string
	Path		string
	CurrentLine int64
}

func NewLogParse() *LogParse {
	return &LogParse{CurrentLine: 1}
}

func (logParse *LogParse) ParseOldLogs(logData *LogData) (currentLogFile string) {
	logPath := &config.ConfigManager.MonitorConfig.Nodes.MainChain.LogPath
	logFiles, err := file.GetDirectoryFileNames(*logPath)
	if err != nil {
		log.Errorf("get log filenames error: %+v", err)
		return ""
		//return nil, err
	}

	currentLogFile = logFiles.Back().Value.(string)
	logFiles.Remove(logFiles.Back())

	for element := logFiles.Front(); element != nil; element = element.Next() {
		logParse.RestLinePosition(1)
		logPath := path.Join(element.Value.(string), *logPath)
		err := logParse.ReadLogFile(logData, logPath)
		errorhelper.Warn(err , "read log file failed!")
	}

	//logParser := NewLogParse(currentLogFile, *logPath)
	//for {
	//	time.Sleep(time.Duration(config.ConfigManager.MonitorConfig.Nodes.MainChain.LogFreshInterval) * time.Second)
	//	logParser.ReadLogFile()
	//}
	return currentLogFile
}

func (logParse *LogParse) ReadLogFile(logData *LogData, logPath string) error {
	//logPath := path.Join(logParse.Path, logParse.FileName)
	log.Infof("read log from %s", logPath)
	lines, err := file.ReadLastLinesFromFile(logPath, logParse.CurrentLine)
	if err != nil {
		log.Error("read log file in lines failed!")
		return err
	}

	for line := lines.Front(); line != nil; line = line.Next() {
		logParse.ReadLine(logData, line.Value.(string))
	}

	logParse.CurrentLine = logParse.CurrentLine + int64(lines.Len())
	return nil
}

func (logParse *LogParse) RestLinePosition(position int64)  {
	logParse.CurrentLine = position
}

func (logParse *LogParse) ReadLine(logData *LogData, line string) {
	switch {
	case strings.Contains(line, constants_logmark.NodeVersion):
		logData.Version = logParse.readVersion(line)

	case strings.Contains(line, constants_logmark.InternalNbrAddress):
		logData.Network.PushBack(logParse.readInternalNbr(line))

	case strings.Contains(line, constants_logmark.OnVoteArrived):
		logData.VoteArrived.PushBack(logParse.readOnVoteArrived(line))

	case strings.Contains(line, constants_logmark.OnProposalArrived):
		logData.ProposalArrived.PushBack(logParse.readOnProposalArrived(line))

	case strings.Contains(line, constants_logmark.OnProposalFinished):
		logData.ProposalFinished.PushBack(logParse.readOnProposalFinished(line))

	case strings.Contains(line, constants_logmark.OnViewStarted):
		logData.ViewStarted.PushBack(logParse.readOnViewStarted(line))

	case strings.Contains(line, constants_logmark.OnConsensusStarted):
		logData.ConsensusStarted.PushBack(logParse.readOnConsensusStarted(line))

	case strings.Contains(line, constants_logmark.OnConsensusFinished):
		logData.ConsensusFinished.PushBack(logParse.readOnConsensusFinished(line))
	}
}

func (logParse *LogParse) readVersion(line string) string {
	version, _ := logParse.readContent(line, constants_logmark.NodeVersion)
	return version
}

func (logParse *LogParse) readInternalNbr(line string) *models.Network {
	logTime, content := logParse.readContent(line, constants_logmark.InternalNbrAddress)
	content = strings.TrimLeft(content, "[")
	content = strings.TrimRight(content, "]")
	nbrHosts := strings.Split(content, " ")
	network := models.Network{ *logParse.parseLogTime(logTime),&nbrHosts}
	return &network
}

func (logParse *LogParse) readOnViewStarted(line string) *models.ViewStart {
	valueMap := logParse.readProperties(line, constants_logmark.OnViewStarted)
	startTime, _ := time.Parse(constants_logmark.LogTimeParseLayout1, (*valueMap)["StartTime"])
	offset, _ := strconv.ParseInt((*valueMap)["Offset"], 10, 16)
	height, _ := strconv.ParseUint((*valueMap)["Height"], 10, 32)
	viewStart := models.ViewStart{
		*logParse.parseLogTime((*valueMap)["LogTime"]),
		(*valueMap)["OnDutyArbitrator"],
		startTime,
		int16(offset),
		uint32(height)}
	return &viewStart
}

func (logParse *LogParse) readOnVoteStarted(line string) *models.ViewStart {
	valueMap := logParse.readProperties(line, constants_logmark.OnViewStarted)
	startTime, _ := time.Parse(constants_logmark.LogTimeParseLayout1, (*valueMap)["StartTime"])
	offset, _ := strconv.ParseInt((*valueMap)["Offset"], 10, 16)
	height, _ := strconv.ParseUint((*valueMap)["Height"], 10, 32)
	viewStart := models.ViewStart{
		*logParse.parseLogTime((*valueMap)["LogTime"]),
		(*valueMap)["OnDutyArbitrator"],
		startTime, int16(offset),
		uint32(height)}
	return &viewStart
}

func (logParse *LogParse) readOnProposalArrived(line string) *models.ProposalMessage {
	return logParse.readOnProposal(line, constants_logmark.OnProposalArrived)
}

func (logParse *LogParse) readOnProposalFinished(line string) *models.ProposalMessage {
	return logParse.readOnProposal(line, constants_logmark.OnProposalFinished)
}

func (logParse *LogParse) readOnProposal(line, logMark string) *models.ProposalMessage {
	valueMap := logParse.readProperties(line, logMark)
	receivedTime, _ := time.Parse(constants_logmark.LogTimeParseLayout1, (*valueMap)["ReceivedTime"])
	result, _ := strconv.ParseBool((*valueMap)["Result"])
	message := models.ProposalMessage {
		*logParse.parseLogTime((*valueMap)["LogTime"]),
		(*valueMap)["Proposal"],
		(*valueMap)["BlockHash"],
		receivedTime,
		result}
	return &message
}

func (logParse *LogParse) readOnVoteArrived(line string) *models.VoteArrivedMessage {
	valueMap := logParse.readProperties(line, constants_logmark.OnVoteArrived)
	receivedTime, _ := time.Parse(constants_logmark.LogTimeParseLayout1, (*valueMap)["ReceivedTime"])
	result, _ := strconv.ParseBool((*valueMap)["Result"])
	message := models.VoteArrivedMessage {
		*logParse.parseLogTime((*valueMap)["LogTime"]),
		(*valueMap)["Signer"],
		receivedTime,
		result}
	return &message
}

func (logParse *LogParse) readOnConsensusStarted(line string) *models.ConsensusMessage {
	return logParse.readOnConsensus(line, constants_logmark.OnConsensusStarted, "StartTime")
}

func (logParse *LogParse) readOnConsensusFinished(line string) *models.ConsensusMessage {
	return logParse.readOnConsensus(line, constants_logmark.OnConsensusFinished, "EndTime")
}

func (logParse *LogParse) readOnConsensus(line, logMark, timeProperty string) *models.ConsensusMessage {
	valueMap := logParse.readProperties(line, logMark)
	receivedTime, _ := time.Parse(constants_logmark.LogTimeParseLayout1, (*valueMap)[timeProperty])
	height, _ := strconv.ParseUint((*valueMap)["Height"], 10, 32)
	message := models.ConsensusMessage {
		*logParse.parseLogTime((*valueMap)["LogTime"]),
		receivedTime,
		uint32(height)}
	return &message
}

func (logParse *LogParse) readProperties(line, logMark string) *map[string]string {
	logTime, content := logParse.readContent(line, logMark)
	properties := strings.Split(content, ", ")
	values := make(map[string]string, len(properties))
	values["LogTime"] = logTime
	for _, property := range properties {
		value := strings.Split(property, ":")
		values[value[0]] = strings.TrimLeft(value[1], "")
	}

	return &values
}

func (logParse *LogParse) readContent(line, logMark string) (logTime, content string) {
	position := strings.Index(line, logMark)
	position = position + len(logMark)
	content = strings.TrimLeft(line[position:], "")
	return line[:26], content
}

func (logParse *LogParse) parseLogTime(timeStr string) *time.Time {
	logTime, err := time.Parse(constants_logmark.LogTimeParseLayout2, timeStr)
	if err != nil {
		log.Errorf("invalidate log time parse: %s", timeStr)
		panic(err)
	}

	return &logTime
}