package models

type Config struct {
	MonitorConfig Configuration `json:"Configuration"`
}

type LogConfig struct {
	Path string 			`json:"Path"`
	Level uint8    			`json:"Level"`
	MaxPerLogSizeMb int64	`json:"MaxPerLogSizeMb"`
	MaxLogsSizeMb int64		`json:"MaxLogsSizeMb"`
}

type Configuration struct {
	Version  int	`json:"Version"`
	AppName  string	`json:"AppName"`
	Log *LogConfig	`json:"Log"`
	Nodes *Nodes	`json:"Nodes"`
	EMail *EMail	`json:"EMail"`
}

type Nodes struct {
	MainChain *MainChain	`json:"MainChain"`
}

type EMail struct {
	Host 		string		`json:"Host"`
	UserName 	string		`json:"UserName"`
	PassWord 	string		`json:"PassWord"`
	NotifyUser 	[]string	`json:"NotifyUser"`
}

type MainChain struct {
	Host string					`json:"Host"`
	RpcPort uint16				`json:"RpcPort"`
	RestfulPort uint16			`json:"RestfulPort"`
	JarServer *JarServer		`json:"JarServer"`
	LogPath string				`json:"LogPath"`
	LogFreshInterval int64		`json:"LogFreshInterval"`
}

type JarServer struct {
	Url
	Binary string `json:"Binary"`
}