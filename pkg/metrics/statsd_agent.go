package metrics

import (
	"strings"
)

type StatsdAgent struct {
	client statsdClient
}

func NewStatsdAgent(client statsdClient) *StatsdAgent {
	return &StatsdAgent{client: client}
}

func (agent *StatsdAgent) IncrementCounter(metric Metric, amount float64) error {
	return agent.client.Incr(strings.ToLower(string(metric)), nil, amount)
}
