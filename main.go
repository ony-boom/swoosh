package main

import (
	"log"
	"os"
	"strconv"
	"time"

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
		log.Printf("Error: %v", err)
		log.Println("Make sure your audio server is running and try again.")
		log.Println("Tip: You can set SWOOSH_STARTUP_DELAY=5 to wait 5 seconds before starting.")
		return
	}

	ui.Init(p)
}
