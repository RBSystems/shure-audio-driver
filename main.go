package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/byuoitav/shure-audio-driver/db"
	"github.com/byuoitav/shure-audio-library"
)

func main() {
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

	}

	if len(address) > 0 {
		// read events on that receiver
		err = readEvents(address)
		if err != nil {
			fmt.Printf("failed to read events: %s\n", err.Error())
			return
		}
	}
	// report error
	fmt.Printf("there are no receivers in this room\n")
}

func readEvents(address string) error {
	// make a connection with address
	control := &shure.AudioControl{
		Address: address,
	}

	conn, err := control.GetConnection()
	if err != nil {
		return err
	}

	for {
		// wait for events to come in from that device
		data, err := conn.ReadEvent()
		if err != nil {

		}
		fmt.Printf(data)

		// process data

		// push events to event hub

	}
	return nil
}
