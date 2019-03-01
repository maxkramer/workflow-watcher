package queue

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/project-interstellar/workflow-watcher/pkg"
	"github.com/sirupsen/logrus"
)

type PubSub struct {
	Log            *logrus.Logger
	Ctx            context.Context
	MessageFactory pkg.MessageFactory
	ProjectId      string
	TopicName      string
}

func (pubSub PubSub) Publish(workflow *v1alpha1.Workflow) *error {
	message := pubSub.MessageFactory.NewMessage(workflow)
	serializedMessage, _ := json.Marshal(message)

	client, err := pubsub.NewClient(pubSub.Ctx, pubSub.ProjectId)
	if err != nil {
		return &err
	}

	topic := client.Topic(pubSub.TopicName)

	result := topic.Publish(pubSub.Ctx, &pubsub.Message{Data: serializedMessage})
	_, err = result.Get(pubSub.Ctx)

	return &err
}
