package models

type Block struct {
	Auxpow            string `json:"auxpow"`
	Bits              int    `json:"bits"`
	Chainwork         string `json:"chainwork"`
	Confirmations     int    `json:"confirmations"`
	Difficulty        string `json:"difficulty"`
	Hash              string `json:"hash"`
	Height            uint32    `json:"height"`
	Mediantime        int    `json:"mediantime"`
	Merkleroot        string `json:"merkleroot"`
	Minerinfo         string `json:"minerinfo"`
	Nextblockhash     string `json:"nextblockhash"`
	Nonce             int    `json:"nonce"`
	Previousblockhash string `json:"previousblockhash"`
	Size              int    `json:"size"`
	Strippedsize      int    `json:"strippedsize"`
	Time              int    `json:"time"`
	Tx                []struct {
		Attributes []struct {
			Data  string `json:"data"`
			Usage int    `json:"usage"`
		} `json:"attributes"`
		Blockhash     string `json:"blockhash"`
		Blocktime     int    `json:"blocktime"`
		Confirmations int    `json:"confirmations"`
		Hash          string `json:"hash"`
		Locktime      int    `json:"locktime"`
		Payload       struct {
			CoinbaseData string `json:"CoinbaseData"`
		} `json:"payload"`
		Payloadversion int           `json:"payloadversion"`
		Programs       []interface{} `json:"programs"`
		Size           int           `json:"size"`
		Time           int           `json:"time"`
		Txid           string        `json:"txid"`
		Type           int           `json:"type"`
		Version        int           `json:"version"`
		Vin            []struct {
			Sequence int    `json:"sequence"`
			Txid     string `json:"txid"`
			Vout     int    `json:"vout"`
		} `json:"vin"`
		Vout []struct {
			Address    string      `json:"address"`
			Assetid    string      `json:"assetid"`
			N          int         `json:"n"`
			Outputlock int         `json:"outputlock"`
			Payload    interface{} `json:"payload"`
			Type       int         `json:"type"`
			Value      string      `json:"value"`
		} `json:"vout"`
		Vsize int `json:"vsize"`
	} `json:"tx"`
	Version    int    `json:"version"`
	Versionhex string `json:"versionhex"`
	Weight     int    `json:"weight"`
}

