package queue

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/project-interstellar/workflow-watcher/pkg"
	"k8s.io/apimachinery/pkg/watch"
)

type PubSub struct {
	context        context.Context
	messageFactory pkg.MessageFactory
	projectId string
	topicId string
	topic *pubsub.Topic
}

func NewPubSub(context context.Context, messageFactory pkg.MessageFactory, projectId string, topicId string) *PubSub {
	return &PubSub{
		context:        context,
		messageFactory: messageFactory,
		projectId: projectId,
		topicId: topicId,
	}
}

func (p PubSub) Publish(workflow *v1alpha1.Workflow, eventType watch.EventType) error {
	if p.topic == nil {
		connErr := p.connect()
		if connErr != nil {
			return connErr
		}
	}

	message := p.messageFactory.NewMessage(workflow, eventType)
	serializedMessage, _ := json.Marshal(message)

	result := p.topic.Publish(p.context, &pubsub.Message{Data: serializedMessage})
	_, err := result.Get(p.context)
	return err
}

func (p PubSub) connect() error {
	client, err := pubsub.NewClient(p.context, p.projectId)
	if err != nil {
		return err
	}

	p.topic = client.Topic(p.topicId)
	return nil
}
