package event

type state int

const (
	Interference state = iota + 1
	Power
	Battery
	Unknown
)

var states = [...]string{
	"RF_INT_DET",
	"TX_TYPE",
	"BATT",
	"UNKNOWN",
}

func (s state) String() string {
	return states[s-1]
}
