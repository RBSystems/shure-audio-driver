package event

import (
	"fmt"
	"regexp"
	"time"

	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/v2/events"
)

//Publisher receives device responses through a channel, parses them, then sends them on to the event hub
type Publisher struct {
	RoomID     string
	HubAddress string
	RoomSys    string
	RespCh     chan string
	msg        *messenger.Messenger
}

//PublishEvents blocks on the response channel until a new response arrives to be published
func (p *Publisher) PublishEvents() {
	// start the messenger
	var err error
	p.msg, err = messenger.BuildMessenger(p.HubAddress, base.Messenger, 1000)
	if err != nil {

	}

	for {
		go p.handle(<-p.RespCh)
	}
}

func (p *Publisher) handle(resp string) {
	event, err := p.parseResponse(resp)
	if err != nil {
		// event is invalid, skip it
		return
	}

	//fill in event
	event.FillEventInfo(resp)

	// publish event
	if event.E.Value != flag {
		err = p.publish(event)
		if err != nil {

		}
	}
}

func (p *Publisher) parseResponse(resp string) (*shureEvent, error) {
	re := regexp.MustCompile(`REP (\d)`)
	channel := re.FindStringSubmatch(resp)
	if len(channel) == 0 {
		//no data
		return nil, fmt.Errorf("ignore event")
	}
	deviceID := fmt.Sprintf("%s-MIC%s", p.RoomID, channel[1])

	event := &events.Event{
		TargetDevice: events.GenerateBasicDeviceInfo(deviceID),
	}

	e := &shureEvent{
		DeviceName: fmt.Sprintf("%s-MIC%s", p.RoomID, channel[1]),
		E:          event,
	}
	e.SetEventType(resp)
	if e.Type == unknown {
		return nil, fmt.Errorf("unknown event type")
	}

	return e, nil
}

func (p *Publisher) publish(event *shureEvent) error {
	event.E.GeneratingSystem = event.E.TargetDevice.DeviceID
	event.E.Timestamp = time.Now()
	event.E.AffectedRoom = events.GenerateBasicRoomInfo(p.RoomID)

	if len(p.RoomSys) > 0 {
		event.E.AddToTags(events.RoomSystem)
	}

	// add error to tags? if necessary

	p.msg.SendEvent(*event.E)
	return nil
}
