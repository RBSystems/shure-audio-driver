package publish

import (
	"fmt"
	"regexp"
	"time"

	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/shure-audio-driver/event"
	"github.com/byuoitav/shure-audio-driver/log"
	"go.uber.org/zap"
)

//EventPublisher receives device responses through a channel, parses them, then sends them on to the event hub
type EventPublisher struct {
	RoomID     string
	HubAddress string
	RoomSys    string
	RespCh     chan string
	msg        *messenger.Messenger
}

//StartMessenger will build the event messenger
func (p *EventPublisher) StartMessenger() error {
	// start the messenger
	var err error
	p.msg, err = messenger.BuildMessenger(p.HubAddress, base.Messenger, 1000)
	if err != nil {
		return err
	}
	return nil
}

//PublishEvents blocks on the response channel until a new response arrives to be published
func (p *EventPublisher) PublishEvents() {
	for {
		go p.handle(<-p.RespCh)
	}
}

func (p *EventPublisher) handle(resp string) {
	e, err := p.parseResponse(resp)
	if err != nil {
		// event is invalid, skip it
		return
	}

	e.FillEventInfo(resp)

	if e.E.Value != event.IGNORE {
		log.L.Info("publishing event", zap.String("key", e.E.Key), zap.String("value", e.E.Value))
		err = p.publish(e)
		if err != nil {
			log.L.Error("failed to publish event", zap.Error(err))
		}
	}
}

func (p *EventPublisher) parseResponse(resp string) (*event.ShureEvent, error) {
	re := regexp.MustCompile(`REP (\d)`)
	channel := re.FindStringSubmatch(resp)
	if len(channel) == 0 {
		//no data
		return nil, fmt.Errorf("ignore event")
	}
	deviceID := fmt.Sprintf("%s-MIC%s", p.RoomID, channel[1])

	hubEvent := &events.Event{
		TargetDevice: events.GenerateBasicDeviceInfo(deviceID),
	}

	e := &event.ShureEvent{
		E: hubEvent,
	}
	e.SetEventType(resp)
	if e.Type == event.Unknown {
		return nil, fmt.Errorf("unknown event type")
	}

	return e, nil
}

func (p *EventPublisher) publish(event *event.ShureEvent) error {
	event.E.GeneratingSystem = event.E.TargetDevice.DeviceID
	event.E.Timestamp = time.Now()
	event.E.AffectedRoom = events.GenerateBasicRoomInfo(p.RoomID)

	if len(p.RoomSys) > 0 {
		event.E.AddToTags(events.RoomSystem)
	}

	p.msg.SendEvent(*event.E)
	return nil
}
