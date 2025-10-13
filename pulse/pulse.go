package pulse

import (
	"fmt"

	"github.com/jfreymuth/pulse"
	"github.com/jfreymuth/pulse/proto"
)

func NewPulse() (*Pulse, error) {
	c, err := pulse.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initiate a connection to audio server: %v", err)
	}

	return &Pulse{
		client: c,
	}, nil
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
