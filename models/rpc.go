package models

type ListProducers struct {
	Start uint16 `json:"start"`
	Limit uint16 `json:"limit"`
}

type ListProducersResponse struct {
	Producers	[]Producer	`json:"producers"`
	TotalVotes	string 		`json:"totalvotes"`
	TotalCounts int64 		`json:"totalcounts"`
}

type Height struct {
	Height uint32	`json:"height"`
}

//type DposPeersInfoParameter struct {
//	Count uint8	`json:"count"`
//}

type RpcResponse struct {
	Error 	string 		`json:"error"`
	Id 		int16 		`json:"id"`
	JsonRpc string 		`json:"jsonrpc"`
	Result 	interface{} `json:"result"`
}

type DposPeersInfo struct {
	OwnerPublicKey string `json:"ownerpublickey"`
	NodePublicKey  string `json:"nodepublickey"`
	Ip             string `json:"ip"`
	ConnectState   string `json:"connectstate"`
}
