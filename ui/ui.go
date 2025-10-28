package ui

import (
	"time"

	"deedles.dev/tray"
	"github.com/ony-boom/swoosh/pulse"
)

type state struct {
	refreshTicker        *time.Ticker
	done                 chan struct{}
	lastKnownDefaultSink string
	stopRefresh          chan bool
	sinkCount            int
}

type UI struct {
	item      *tray.Item
	state     *state
	pulse     *pulse.Pulse
	menuItems []*tray.MenuItem
}

func Init(p *pulse.Pulse) {
	ui := UI{
		pulse: p,
		item:  newTray(),
		state: &state{
			done:      make(chan struct{}),
			sinkCount: 0,
		},
	}

	defer ui.item.Close()

	ui.renderSinks()
	ui.renderOptions()
	ui.startPulseMonitoring()

	<-ui.state.done
}
