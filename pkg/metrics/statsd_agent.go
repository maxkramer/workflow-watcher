package metrics

import (
	"strings"

	"github.com/DataDog/datadog-go/statsd"
)

type StatsdAgent struct {
	client *statsd.Client
}

func NewStatsdAgent(client *statsd.Client) *StatsdAgent {
	return &StatsdAgent{client: client}
}

func (agent *StatsdAgent) IncrementCounter(metric Metric, amount float64) error {
	return agent.client.Incr(strings.ToLower(string(metric)), nil, amount)
}
