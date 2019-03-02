package metrics

import (
	"testing"

	"github.com/argoproj/argo/errors"
	"github.com/stretchr/testify/assert"
)

type MockStatsdClient struct {
	Error    error
	incrName string
	incrTags []string
	incrRate float64
}

func (client *MockStatsdClient) Incr(name string, tags []string, rate float64) error {
	client.incrName = name
	client.incrTags = tags
	client.incrRate = rate

	return client.Error
}

func TestNewStatsdAgent(t *testing.T) {
	client := &MockStatsdClient{}
	agent := NewStatsdAgent(client)
	assert.Equal(t, agent.client, client)
}

func TestStatsdAgent_IncrementCounter(t *testing.T) {
	client := &MockStatsdClient{}
	agent := NewStatsdAgent(client)

	err := agent.IncrementCounter(EventHandlerAdded, 1)

	assert.Nil(t, err)
	assert.Equal(t, client.incrName, string(EventHandlerAdded))
	assert.Nil(t, client.incrTags)
	assert.Equal(t, client.incrRate, float64(1))
}

func TestStatsdAgent_IncrementCounter_ReturnsErrors(t *testing.T) {
	client := &MockStatsdClient{Error: errors.New("error", "Something went wrong")}
	agent := NewStatsdAgent(client)

	err := agent.IncrementCounter(EventHandlerAdded, 1)
	assert.Equal(t, client.Error, err)
}
