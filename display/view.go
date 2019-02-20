package display

type View struct {
	ChangeTimes		int16
	OnDutyProducer	*ProducerMonitor
	Proposal		*Proposal
	Producers       *[]ProducerMonitor
}