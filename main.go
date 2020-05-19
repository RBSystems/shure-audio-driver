package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/byuoitav/shure-audio-driver/db"
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

	// read events on that receiver
	err = readEvents(address)
}

func readEvents(address string) error {
	// make a connection with address
	// wait for events to come in from that device
	// push events to event hub
	return nil
}
