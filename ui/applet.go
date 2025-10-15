package ui

import (
	"log"

	"deedles.dev/tray"
	"github.com/ony-boom/swoosh/pulse"
)

func Init(p *pulse.Pulse) {
	done := make(chan struct{})

	item, err := tray.New(
		tray.ItemID("ony.world.swoosh"),
		tray.ItemTitle("Switch audio sinks/sources"),
		tray.ItemIsMenu(true),
		tray.ItemIconName("audio-card"),
		tray.ItemHandler(tray.ActivateHandler(func(x, y int) error {
			log.Println("ActivateHandler")
			return nil
		})),
	)

	// Set the global done channel for quit functionality
	setGlobalDoneChannel(done)

	// Start pulse audio monitoring to detect changes
	startPulseMonitoring(item, p)

	if err != nil {
		log.Fatal(err)
	}

	defer item.Close()

	renderSinks(item, p)
	renderOptions(item, p)

	<-done
}
