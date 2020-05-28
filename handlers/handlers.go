package handlers

import (
	"net/http"

	"github.com/byuoitav/shure-audio-driver/log"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func GetBattery(c echo.Context) error {
	channel := c.Param("channel")
	format := c.Param("format")
	log.L.Info("getting battery status", zap.String("channel", channel), zap.String("format", format))

	resp, err := runBatteryCommand(channel, format)
	if err != nil {
		log.L.Error("failed to get battery status", zap.String("channel", channel), zap.String("format", format), zap.Error(err))
		return c.String(http.StatusInternalServerError, "failed to get battery status")
	}

	log.L.Info("got batter status", zap.String("channel", channel), zap.String("format", format), zap.String("status", resp))
	return c.JSON(http.StatusOK, resp)
}

func GetPower(c echo.Context) error {
	channel := c.Param("channel")
	log.L.Info("getting power status", zap.String("channel", channel))

	resp, err := runPowerCommand(channel)
	if err != nil {
		log.L.Error("failed to get power status", zap.String("channel", channel), zap.Error(err))
		return c.String(http.StatusInternalServerError, "failed to get power status")
	}

	log.L.Info("got power status", zap.String("channel", channel), zap.String("status", resp))
	return c.JSON(http.StatusOK, resp)
}
