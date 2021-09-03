package main

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/theaj/cloudflare-ip-updater/monitor"
)

func init() {
	if dev, exists := os.LookupEnv("DEV_MODE"); exists {
		if devMode, err := strconv.ParseBool(dev); err != nil {
			log.Err(err).Msgf("Could not parse DEV_MODE environment variable")
		} else if devMode {
			log.Logger = log.Output(zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: time.RFC3339,
			})
		}
	}
}

func main() {
	monitor.Start()
}
