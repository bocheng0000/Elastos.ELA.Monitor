package display

type Content struct {
	Version			string
	Height 			uint32
	CPUUsed			float64
	MemoryUsed		float64
	Host 			string
	RpcPort 		uint16
	RestfulPort 	uint16
}