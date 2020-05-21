package event

type state int

const (
	interference state = iota + 1
	power
	battery
	unknown
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
