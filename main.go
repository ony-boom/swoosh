package main

import (
	"log"

	"github.com/ony-boom/swoosh/pulse"
	"github.com/ony-boom/swoosh/ui"
)

func main() {
	p, err := pulse.NewPulse()
	if err != nil {
		log.Fatal(err)
	}

	ui.Init(p)
}
