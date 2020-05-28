package main

import (
	"io"

	"github.com/byuoitav/shure-audio-driver/log"
	"github.com/byuoitav/shure-audio-driver/publish"
	"github.com/byuoitav/shure-audio-library"
	"go.uber.org/zap"
)

func monitorEvents(roomID, address string, pub *publish.EventPublisher) {
	err := pub.StartMessenger()
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
