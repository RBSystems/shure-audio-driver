package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/byuoitav/shure-audio-driver/db"
	"github.com/byuoitav/shure-audio-driver/event"
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

	pub := &event.Publisher{
		RoomID:     roomID,
		HubAddress: os.Getenv("HUB_ADDRESS"),
		RespCh:     make(chan string, 100),
	}

	// start waiting to publish events
	go pub.HandleEvents()

	if len(address) > 0 {
		fmt.Println(address)

		err = readEvents(address, pub)
		if err != nil {
			fmt.Printf("failed to read events: %s\n", err.Error())
			return
		}
	}
	fmt.Printf("there are no receivers in this room\n")
}

func readEvents(address string, pub *event.Publisher) error {
	// make a connection with address
	control := &shure.AudioControl{
		Address: address,
		Port:    "2202",
	}

	conn, err := control.GetConnection()
	if err != nil {
		return err
	}

	for {
		data, err := conn.ReadEvent()
		if err == io.EOF {
			fmt.Println("got an eof")
			conn.Conn.Close()
			conn, err = control.GetConnection()
		}
		if err != nil {
			return err
		}

		//send event to be published
		pub.RespCh <- data
	}
}
