package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/byuoitav/shure-audio-driver/db"
	"github.com/byuoitav/shure-audio-library"
)

type command struct {
	Message string `json:"message"`
}

func runBatteryCommand(channel, format string) (string, error) {
	ch, err := strconv.Atoi(channel)
	if err != nil {
		return "", fmt.Errorf("channel is invalid: %s", err.Error())
	}

	address, err := getAddress()
	if err != nil {
		return "", fmt.Errorf("failed to get receiver address: %s", err.Error())
	}

	control := &shure.AudioControl{
		Address: address,
	}

	conn, err := control.GetConnection()
	if err != nil {
		return "", fmt.Errorf("failed to open connection to receiver: %s", err.Error())
	}

	var resp string
	switch format {
	case "percentage":
		resp, err = conn.GetBatteryCharge(ch)
	case "time":
		resp, err = conn.GetBatteryRunTime(ch)
	case "bars":
		resp, err = conn.GetBatteryBars(ch)
	default:
		return "", fmt.Errorf("format is invalid")
	}

	if err != nil {
		return "", fmt.Errorf("failed to run battery command: %s", err.Error())
	}

	return resp, nil
}

func runPowerCommand(channel string) (string, error) {
	ch, err := strconv.Atoi(channel)
	if err != nil {
		return "", fmt.Errorf("channel is invalid: %s", err.Error())
	}

	address, err := getAddress()
	if err != nil {
		return "", fmt.Errorf("failed to get receiver address: %s", err.Error())
	}

	control := &shure.AudioControl{
		Address: address,
	}

	conn, err := control.GetConnection()
	if err != nil {
		return "", fmt.Errorf("failed to open connection to receiver: %s", err.Error())
	}

	resp, err := conn.GetPowerStatus(ch)
	if err != nil {
		return "", fmt.Errorf("failed to get power status: %s", err.Error())
	}

	return resp, nil
}

func sendRawCommand(cmd command) (string, error) {
	address, err := getAddress()
	if err != nil {
		return "", fmt.Errorf("failed to get address: %s", err.Error())
	}

	control := &shure.AudioControl{
		Address: address,
	}

	conn, err := control.GetConnection()
	if err != nil {
		return "", fmt.Errorf("failed to open connection to receiver: %s", err.Error())
	}

	resp, err := conn.SendCommand(cmd.Message)
	if err != nil {
		return "", fmt.Errorf("failed to send command: %s", err.Error())
	}

	return resp, nil
}

func getAddress() (string, error) {
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
		return "", fmt.Errorf("failed to get receiver address: %s", err.Error())
	}
	return address, nil
}
