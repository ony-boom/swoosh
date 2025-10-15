package ui

import (
	"log"
	"time"

	"deedles.dev/tray"
	"github.com/ony-boom/swoosh/pulse"
)

// Store menu items globally so they can be cleared
var menuItems []*tray.MenuItem

const MAX_LENGTH = 32 // Maximum length for a menu item label

func clearMenuItems() {
	for _, item := range menuItems {
		if item != nil {
			item.Remove()
		}
	}
	menuItems = nil
}

func addMenuItem(item *tray.MenuItem) {
	menuItems = append(menuItems, item)
}

func renderSinks(item *tray.Item, p *pulse.Pulse) {
	sinks, err := p.ListSinks()
	sinksHeader, _ := item.Menu().AddChild(
		tray.MenuItemLabel("Sinks"),
		tray.MenuItemEnabled(false),
	)
	addMenuItem(sinksHeader)
	if err != nil {
		errorItem, _ := item.Menu().AddChild(
			tray.MenuItemLabel("Error loading sinks"),
		)
		addMenuItem(errorItem)
		return
	}

	var sinkItems []*tray.MenuItem
	for i, sink := range sinks {
		defaultState := tray.Off
		if p.IsDefaultSink(sink) {
			defaultState = tray.On
		}

		label := " " + sink.Name // Add a space for better readability
		if len(label) > MAX_LENGTH {
			label = label[:MAX_LENGTH-3] + "..."
		}

		sinkItem, _ := item.Menu().AddChild(
			tray.MenuItemToggleType(tray.Radio),
			tray.MenuItemLabel(label),
			tray.MenuItemToggleState(defaultState),
		)
		addMenuItem(sinkItem)

		sinkItem.RequestActivation(uint32(i))

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
				log.Fatalf("Failed to set default sink: %v", err)
			}
			// Update the monitoring state to trigger refresh on next check
			go func() {
				// Small delay to allow PulseAudio to process the change
				time.Sleep(100 * time.Millisecond)
				// Force an update by clearing the last known sink
				// This will cause the monitoring to detect a change
				resetMonitoringState()
			}()
			return nil
		}))
	}
}

func renderOptions(item *tray.Item, p *pulse.Pulse) {
	separator, _ := item.Menu().AddChild(tray.MenuItemType(tray.Separator))
	addMenuItem(separator)

	refreshItem, _ := item.Menu().AddChild(
		tray.MenuItemLabel("Refresh"),
		tray.MenuItemHandler(tray.ClickedHandler(func(data any, timestamp uint32) error {
			triggerManualRefresh(item, p)
			return nil
		})),
	)
	addMenuItem(refreshItem)

	quitItem, _ := item.Menu().AddChild(
		tray.MenuItemLabel("Quit"),
		tray.MenuItemHandler(tray.ClickedHandler(func(data any, timestamp uint32) error {
			stopPulseMonitoring()
			if globalDoneChannel != nil {
				close(globalDoneChannel)
			}
			return nil
		})),
	)
	addMenuItem(quitItem)
}
