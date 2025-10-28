package ui

import (
	"os"
	"time"

	"deedles.dev/tray"
	"github.com/ony-boom/swoosh/logger"
)

const MAX_LENGTH = 32 // Maximum length for a menu item label

func (ui *UI) clearMenuItems() {
	for _, item := range ui.menuItems {
		if item != nil {
			item.Remove()
		}
	}
	ui.menuItems = nil
}

func (ui *UI) addMenuItem(item *tray.MenuItem) {
	ui.menuItems = append(ui.menuItems, item)
}

func (ui *UI) renderSinks() {
	sinks, err := ui.pulse.ListSinks()
	sinksHeader, _ := ui.item.Menu().AddChild(
		tray.MenuItemLabel("Sinks"),
		tray.MenuItemEnabled(false),
	)
	ui.addMenuItem(sinksHeader)

	if err != nil {
		errorItem, _ := ui.item.Menu().AddChild(
			tray.MenuItemLabel("Error loading sinks"),
		)
		ui.addMenuItem(errorItem)
		return
	}

	var sinkItems []*tray.MenuItem
	for i, sink := range sinks {
		defaultState := tray.Off
		if ui.pulse.IsDefaultSink(sink) {
			defaultState = tray.On
		}

		label := " " + sink.Name // Add a space for better readability
		if len(label) > MAX_LENGTH {
			label = label[:MAX_LENGTH-3] + "..."
		}

		sinkItem, _ := ui.item.Menu().AddChild(
			tray.MenuItemToggleType(tray.Radio),
			tray.MenuItemLabel(label),
			tray.MenuItemToggleState(defaultState),
		)
		ui.addMenuItem(sinkItem)

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
			if err := ui.pulse.SetDefaultSink(currentSink.ID); err != nil {
				logger.Log.Error("Failed to set default sink: %v", err)
				os.Exit(1)
			}
			// Update the monitoring state to trigger refresh on next check
			go func() {
				// Small delay to allow PulseAudio to process the change
				time.Sleep(100 * time.Millisecond)
				// Force an update by clearing the last known sink
				// This will cause the monitoring to detect a change
				ui.resetMonitoringState()
			}()
			return nil
		}))
	}
}

func (ui *UI) renderOptions() {
	separator, _ := ui.item.Menu().AddChild(tray.MenuItemType(tray.Separator))
	ui.addMenuItem(separator)

	refreshItem, _ := ui.item.Menu().AddChild(
		tray.MenuItemLabel("Refresh"),
		tray.MenuItemHandler(tray.ClickedHandler(func(data any, timestamp uint32) error {
			ui.triggerManualRefresh()
			return nil
		})),
	)

	ui.addMenuItem(refreshItem)

	quitItem, _ := ui.item.Menu().AddChild(
		tray.MenuItemLabel("Quit"),
		tray.MenuItemHandler(tray.ClickedHandler(func(data any, timestamp uint32) error {
			ui.stopPulseMonitoring()
			if ui.state.done != nil {
				close(ui.state.done)
			}
			return nil
		})),
	)
	ui.addMenuItem(quitItem)
}
