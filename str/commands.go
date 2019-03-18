package str

type CommandAction int

const (
	RS CommandAction = iota
	MD
	STOP
	GetMS
	GetMF
	Reconnect
	Connect
)
