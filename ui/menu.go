package ui

import (
	"log"

	"deedles.dev/tray"
)

func newTray() *tray.Item {
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
	if err != nil {
		log.Fatal(err)
	}

	return item
}
