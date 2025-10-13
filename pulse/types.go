package pulse

import (
	"github.com/jfreymuth/pulse"
)



type Pulse struct {
	client *pulse.Client
}

type SimpleSink struct {
	ID   string
	Name string
}
