package models

type Producer struct {
	OwnerPublicKey	string	`json:"ownerpublickey"`
	NodePublicKey 	string	`json:"nodepublickey"`
	NickName 		string	`json:"nickname"`
	Url 			string	`json:"url"`
	Location 		int64	`json:"location"`
	Active 			bool 	`json:"active"`
	Votes 			string	`json:"votes"`
	Ip 				string	`json:"ip"`
	Index 			int16	`json:"index"`
}
