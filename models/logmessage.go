package models

import "time"

type ProposalMessage struct {
	LogTime			time.Time
	Proposal		string
	BlockHash		string
	ReceivedTime	time.Time
	Result			bool
}

type ConsensusMessage struct {
	LogTime			time.Time
	ReceivedTime	time.Time
	Height			uint32
}

type VoteArrivedMessage struct {
	LogTime			time.Time
	Signer			string
	ReceivedTime	time.Time
	Result			bool
}

type ViewStart struct {
	LogTime				time.Time
	OnDutyArbitrator	string
	StartTime			time.Time
	Offset				int16
	Height				uint32
}

type Network struct {
	LogTime			time.Time
	NbrHosts		*[]string
}