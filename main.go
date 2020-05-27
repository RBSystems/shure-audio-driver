package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/byuoitav/shure-audio-driver/db"
	"github.com/byuoitav/shure-audio-driver/log"
	"github.com/byuoitav/shure-audio-driver/publish"
	"github.com/byuoitav/shure-audio-library"
	"go.uber.org/zap"
)

func main() {
	log.StartLogger()

	// get receiver address from db
	db := &db.Database{
		Address:  os.Getenv("DB_ADDRESS"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	systemID := os.Getenv("SYSTEM_ID")
	s := strings.Split(systemID, "-")
	roomID := fmt.Sprintf("%s-%s", s[0], s[1])

	address, err := db.GetReceiverAddress(roomID)
	if err != nil {
		log.L.Fatal("failed to get receiver address", zap.Error(err))
	}

	pub := &publish.EventPublisher{
		RoomID:     roomID,
		HubAddress: os.Getenv("HUB_ADDRESS"),
		RoomSys:    os.Getenv("ROOM_SYSTEM"),
		RespCh:     make(chan string, 100),
	}

	err = pub.StartMessenger()
	if err != nil {
		log.L.Fatal("failure when building event hub messenger", zap.Error(err))
	}

	// start waiting to publish events
	go pub.PublishEvents()

	if len(address) > 0 {
		err = readEvents(address, pub)
		if err != nil {
			log.L.Fatal("Failure when connecting and reading events", zap.Error(err))
		}
	}
	log.L.Error("There are no receivers in this room. Stopping service...")
}

func readEvents(address string, pub *publish.EventPublisher) error {
	log.L.Info("connecting to receiver", zap.String("address", address))
	control := &shure.AudioControl{
		Address: address,
	}

	conn, err := control.GetConnection()
	if err != nil {
		return err
	}

	log.L.Info("connected to receiver", zap.String("address", address))
	for {
		data, err := conn.ReadEvent()
		if err == io.EOF {
			conn.Conn.Close()
			conn, err = control.GetConnection()
		} else if err != nil {
			return err
		}

		//send event to be published
		pub.RespCh <- data
	}
}
