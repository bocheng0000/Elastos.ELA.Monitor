package display

type ProducerMonitor struct {
	NickName       	string
	OwnerPublicKey 	string
	NodePublicKey 	string
	//Host		   	string
	Deposit     	int16
	Vote        	int64
	IsActive    	bool
	Evidence    	string
}
