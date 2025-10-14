package ui

import (
	"time"

	"deedles.dev/tray"
	"github.com/ony-boom/swoosh/pulse"
)

var (
	lastKnownDefaultSink string
	refreshTicker        *time.Ticker
	stopRefresh          chan bool
	globalDoneChannel    chan struct{}
)

func startPulseMonitoring(item *tray.Item, p *pulse.Pulse) {
	sinks, err := p.ListSinks()
	if err == nil {
		for _, sink := range sinks {
			if p.IsDefaultSink(sink) {
				lastKnownDefaultSink = sink.ID
				break
			}
		}
	}

	refreshTicker = time.NewTicker(2 * time.Second)
	stopRefresh = make(chan bool, 1)

	go func() {
		for {
			select {
			case <-refreshTicker.C:
				checkForPulseChanges(item, p)
			case <-stopRefresh:
				return
			}
		}
	}()
}

func stopPulseMonitoring() {
	if refreshTicker != nil {
		refreshTicker.Stop()
	}
	if stopRefresh != nil {
		select {
		case stopRefresh <- true:
		default:
		}
	}
}

func checkForPulseChanges(item *tray.Item, p *pulse.Pulse) {
	sinks, err := p.ListSinks()
	if err != nil {
		return
	}

	var currentDefaultSink string
	for _, sink := range sinks {
		if p.IsDefaultSink(sink) {
			currentDefaultSink = sink.ID
			break
		}
	}

	if currentDefaultSink != lastKnownDefaultSink {
		lastKnownDefaultSink = currentDefaultSink
		refreshMenu(item, p)
	}
}

func refreshMenu(item *tray.Item, p *pulse.Pulse) {
	clearMenuItems()

	renderSinks(item, p)

	separator, _ := item.Menu().AddChild(tray.MenuItemType(tray.Separator))
	addMenuItem(separator)

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

// func triggerManualRefresh(item *tray.Item, p *pulse.Pulse) {
// 	refreshMenu(item, p)
// }

func resetMonitoringState() {
	lastKnownDefaultSink = ""
}

func setGlobalDoneChannel(done chan struct{}) {
	globalDoneChannel = done
}
