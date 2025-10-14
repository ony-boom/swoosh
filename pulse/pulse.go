package pulse

import (
	"fmt"
	"log"
	"time"

	"github.com/jfreymuth/pulse"
	"github.com/jfreymuth/pulse/proto"
)

func NewPulse() (*Pulse, error) {
	return NewPulseWithRetry(30*time.Second, 500*time.Millisecond)
}

func NewPulseWithRetry(maxWaitTime, initialDelay time.Duration) (*Pulse, error) {
	var c *pulse.Client
	var err error

	delay := initialDelay
	maxDelay := 5 * time.Second
	start := time.Now()

	log.Printf("Attempting to connect to PulseAudio...")

	for time.Since(start) < maxWaitTime {
		c, err = pulse.NewClient()
		if err == nil {
			return &Pulse{
				client: c,
			}, nil
		}
		time.Sleep(delay)
		delay = min(delay*2, maxDelay)
	}

	return nil, fmt.Errorf("failed to connect to PulseAudio after %v: %v", maxWaitTime, err)
}

func IsTheAudioServerAvailable() bool {
	c, err := pulse.NewClient()
	if err != nil {
		return false
	}
	c.Close()
	return true
}

func simpleSinkFromPulse(pulseSink *pulse.Sink) *SimpleSink {
	return &SimpleSink{
		ID:   pulseSink.ID(),
		Name: pulseSink.Name(),
	}
}

func (p *Pulse) DefaultSink() (*SimpleSink, error) {
	pulseDefault, err := p.client.DefaultSink()
	if err != nil {
		return nil, fmt.Errorf("failed to get the default sink: %v", err)
	}

	return simpleSinkFromPulse(pulseDefault), nil
}

func (p *Pulse) ListSinks() ([]*SimpleSink, error) {
	var sinks []*SimpleSink

	pulseSinks, err := p.client.ListSinks()
	if err != nil {
		return nil, fmt.Errorf("failed to list all sinks: %v", err)
	}

	for _, sink := range pulseSinks {
		sinks = append(sinks, simpleSinkFromPulse(sink))
	}

	return sinks, err
}

func (p *Pulse) SetDefaultSink(sinkId string) error {
	err := p.client.RawRequest(&proto.SetDefaultSink{
		SinkName: sinkId,
	}, nil)
	if err != nil {
		return fmt.Errorf("could'nt update the default sink")
	}

	return nil
}

func (p *Pulse) IsDefaultSink(s *SimpleSink) bool {
	defaultSink, _ := p.DefaultSink()

	return defaultSink.ID == s.ID
}

func (p *Pulse) Close() {
	p.client.Close()
}
