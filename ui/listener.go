package ui

import (
	"time"
)

func (ui *UI) startPulseMonitoring() {
	sinks, err := ui.pulse.ListSinks()
	if err == nil {
		ui.state.sinkCount = len(sinks)
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

	currentSinkCount := len(sinks)
	var currentDefaultSink string
	for _, sink := range sinks {
		if ui.pulse.IsDefaultSink(sink) {
			currentDefaultSink = sink.ID
			break
		}
	}

	// Check if either sink count or default sink changed
	if currentSinkCount != ui.state.sinkCount || currentDefaultSink != ui.state.lastKnownDefaultSink {
		ui.state.sinkCount = currentSinkCount
		ui.state.lastKnownDefaultSink = currentDefaultSink
		ui.refreshMenu()
	}
}

func (ui *UI) refreshMenu() {
	ui.clearMenuItems()
	
	// Update sink count when refreshing
	sinks, err := ui.pulse.ListSinks()
	if err == nil {
		ui.state.sinkCount = len(sinks)
	}
	
	ui.renderSinks()
	ui.renderOptions()
}

func (ui *UI) triggerManualRefresh() {
	ui.pulse.UpdateConfig()
	ui.refreshMenu()
}

func (ui *UI) resetMonitoringState() {
	ui.state.lastKnownDefaultSink = ""
	ui.state.sinkCount = 0
}
