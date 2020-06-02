package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/byuoitav/shure-audio-driver/db"
	"github.com/byuoitav/shure-audio-driver/handlers"
	"github.com/byuoitav/shure-audio-driver/log"
	"github.com/byuoitav/shure-audio-driver/publish"
	"github.com/labstack/echo"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	log.StartLogger()

	var port int
	var logLevel string
	pflag.IntVarP(&port, "port", "p", 8013, "port to run the server on")
	pflag.StringVarP(&logLevel, "log-level", "l", "Info", "level of logging wanted. Debug, Info, Warn, Error, Panic")
	pflag.Parse()

	// set the initial log level
	if err := log.SetLogLevel(logLevel); err != nil {
		log.L.Fatal("unable to set log level", zap.Error(err), zap.String("got", logLevel))
	}

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

	go monitorEvents(roomID, address, pub)

	// build echo server
	e := echo.New()

	e.GET("/:channel/battery/:format", handlers.GetBattery)
	e.GET("/:channel/power", handlers.GetPower)

	e.POST("/command", handlers.SendCommand)

	addr := fmt.Sprintf(":%d", port)
	err = e.Start(addr)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.L.Fatal("failed to start server", zap.Error(err))
	}
}
