package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ony-boom/swoosh/logger"
	"github.com/ony-boom/swoosh/pulse"
	"github.com/ony-boom/swoosh/ui"
)

func main() {
	if delayStr := os.Getenv("SWOOSH_STARTUP_DELAY"); delayStr != "" {
		if delay, err := strconv.Atoi(delayStr); err == nil && delay > 0 {
			log.Printf("Waiting %d seconds before starting (SWOOSH_STARTUP_DELAY)...", delay)
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}

	p, err := pulse.NewPulse()
	if err != nil {
		logger.Log.Error("")
		logger.Log.Info("Tip: You can set SWOOSH_STARTUP_DELAY=5 to wait 5 seconds before starting.")
		logger.Log.Error("Error: %v", err)
		// TODO: Consider adding a retry mechanism here
		// if not use os.Exit(1) to indicate failure
		return
	}

	ui.Init(p)
}
