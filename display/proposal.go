package display

type Proposal struct {
	Total    uint8
	Approval []*string
	Reject   []*string
}