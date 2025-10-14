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

	// Initial render of sinks
	renderSinks(item, p)

	separator, _ := item.Menu().AddChild(tray.MenuItemType(tray.Separator))
	addMenuItem(separator)

	quitItem, _ := item.Menu().AddChild(
		tray.MenuItemLabel("Quit"),
		tray.MenuItemHandler(tray.ClickedHandler(func(data any, timestamp uint32) error {
			stopPulseMonitoring()
			close(done)
			return nil
		})),
	)
	addMenuItem(quitItem)

	<-done
}
