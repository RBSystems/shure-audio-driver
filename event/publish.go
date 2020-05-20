package event

import (
	"fmt"
	"regexp"
)

//Publisher receives device responses through a channel, parses them, then sends them on to the event hub
type Publisher struct {
	RoomID     string
	HubAddress string
	RespCh     chan string
}

//HandleEvents blocks on the response channel until a new response arrives to be published
func (p *Publisher) HandleEvents() {
	for {
		go p.handle(<-p.RespCh)
	}
}

func (p *Publisher) handle(resp string) {
	err := p.parseResponse(resp)
	if err != nil {

	}
	// identify type
	// publish event
}

func (p *Publisher) parseResponse(resp string) error {

	re := regexp.MustCompile(`REP (\d)`)
	channel := re.FindStringSubmatch(resp)
	if len(channel) == 0 {
		//no data
		return fmt.Errorf("ignore event")
	}

	deviceName := fmt.Sprintf("%s-MIC%s", p.RoomID, channel[1])
	fmt.Println(deviceName)

	return nil
}
