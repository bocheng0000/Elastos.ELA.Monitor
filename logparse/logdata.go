package logparse

import "container/list"

type LogData struct {
	Version 			string
	Network				*list.List
	ViewStarted 		*list.List
	ProposalArrived 	*list.List
	ProposalFinished 	*list.List
	ConsensusStarted 	*list.List
	ConsensusFinished 	*list.List
	ChangeView 			*list.List
	VoteArrived 		*list.List
}

func NewLogData() *LogData {
	return &LogData {
		Version: "UnKnown",
		Network: list.New(),
		ViewStarted: list.New(),
		ProposalArrived: list.New(),
		ProposalFinished: list.New(),
		ConsensusStarted: list.New(),
		ConsensusFinished: list.New(),
		ChangeView: list.New(),
		VoteArrived: list.New(),
	}
}