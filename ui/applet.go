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
	)
	if err != nil {
		log.Fatal(err)
	}

	defer item.Close()

	renderSinks(item, p)

	item.Menu().AddChild(tray.MenuItemType(tray.Separator))

	item.Menu().AddChild(
		tray.MenuItemLabel("Quit"),
		tray.MenuItemHandler(tray.ClickedHandler(func(data any, timestamp uint32) error {
			close(done)
			return nil
		})),
	)

	<-done
}

func renderSinks(item *tray.Item, p *pulse.Pulse) {
	sinks, err := p.ListSinks()

	item.Menu().AddChild(
		tray.MenuItemLabel("Sinks"),
		tray.MenuItemEnabled(false),
	)

	if err != nil {
		item.Menu().AddChild(
			tray.MenuItemLabel("Error loading sinks"),
		)
		return
	}

	var sinkItems []*tray.MenuItem

	for _, sink := range sinks {
		defaultState := tray.Off
		if p.IsDefaultSink(sink) {
			defaultState = tray.On
		}

		sinkItem, _ := item.Menu().AddChild(
			tray.MenuItemToggleType(tray.Radio),
			tray.MenuItemLabel(sink.Name),
			tray.MenuItemToggleState(defaultState),
		)

		sinkItems = append(sinkItems, sinkItem)

		currentSink := sink
		currentItem := sinkItem

		sinkItem.SetProps(tray.MenuItemHandler(func(eventID tray.MenuEventID, data any, timestamp uint32) error {
			// Turn this one ON
			currentItem.SetProps(tray.MenuItemToggleState(tray.On))

			// Turn all others OFF
			for _, other := range sinkItems {
				if other != currentItem {
					other.SetProps(tray.MenuItemToggleState(tray.Off))
				}
			}

			if err := p.SetDefaultSink(currentSink.ID); err != nil {
				log.Println("Failed to set default sink:", err)
			}

			return nil
		}))
	}
}
