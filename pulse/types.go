package pulse

import (
	"github.com/jfreymuth/pulse"
	"github.com/ony-boom/swoosh/config"
)

type Pulse struct {
	client *pulse.Client
	config config.Config
}

type SimpleSink struct {
	ID   string
	Name string
}
