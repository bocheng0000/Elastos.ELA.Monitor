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
}

type Nodes struct {
	MainChain *MainChain	`json:"MainChain"`
}

type MainChain struct {
	Host string					`json:"Host"`
	RpcPort int16				`json:"RpcPort"`
	RestfulPort int16			`json:"RestfulPort"`
	JarServer *JarServer		`json:"JarServer"`
	LogPath string				`json:"LogPath"`
	LogFreshInterval int64		`json:"LogFreshInterval"`
}

type JarServer struct {
	Url
	Binary string `json:"Binary"`
}