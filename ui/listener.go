package ui

import (
	"time"
)

func (ui *UI) startPulseMonitoring() {
	sinks, err := ui.pulse.ListSinks()
	if err == nil {
		for _, sink := range sinks {
			if ui.pulse.IsDefaultSink(sink) {
				ui.state.lastKnownDefaultSink = sink.ID
				break
			}
		}
	}

	ui.state.refreshTicker = time.NewTicker(ui.pulse.GetPollInterval())
	ui.state.stopRefresh = make(chan bool, 1)

	go func() {
		for {
			select {
			case <-ui.state.refreshTicker.C:
				ui.checkForPulseChanges()
			case <-ui.state.stopRefresh:
				return
			}
		}
	}()
}

func (ui *UI) stopPulseMonitoring() {
	if ui.state.refreshTicker != nil {
		ui.state.refreshTicker.Stop()
	}
	if ui.state.stopRefresh != nil {
		select {
		case ui.state.stopRefresh <- true:
		default:
		}
	}
}

func (ui *UI) checkForPulseChanges() {
	sinks, err := ui.pulse.ListSinks()
	if err != nil {
		return
	}

	var currentDefaultSink string
	for _, sink := range sinks {
		if ui.pulse.IsDefaultSink(sink) {
			currentDefaultSink = sink.ID
			break
		}
	}

	if currentDefaultSink != ui.state.lastKnownDefaultSink {
		ui.state.lastKnownDefaultSink = currentDefaultSink
		ui.refreshMenu()
	}

	// If the number of sinks has changed, refresh the menu
	if len(sinks) != len(ui.menuItems)-3 {
		ui.refreshMenu()
	}
}

func (ui *UI) refreshMenu() {
	ui.clearMenuItems()
	ui.renderSinks()
	ui.renderOptions()
}

func (ui *UI) triggerManualRefresh() {
	ui.pulse.UpdateConfig()
	ui.refreshMenu()
}

func (ui *UI) resetMonitoringState() {
	ui.state.lastKnownDefaultSink = ""
}
