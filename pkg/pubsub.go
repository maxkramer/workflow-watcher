package pkg

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type PubSub struct {
	Log *logrus.Logger
	Ctx context.Context
	MessageFactory MessageFactory
	ProjectId string
}


func (pubSub PubSub) Append(workflow *v1alpha1.Workflow) *error {
	message := pubSub.MessageFactory.NewMessage(workflow)
	serializedMessage, _ := json.Marshal(message)

	client, err := pubsub.NewClient(pubSub.Ctx, pubSub.ProjectId, option.WithCredentialsFile("/Users/kramer_max/.config/gcloud/interstellar-admin.json"))
	if err != nil {
		return &err
	}

	topic := client.Topic("some-topic")

	_, err = topic.Publish(pubSub.Ctx, &pubsub.Message{Data: serializedMessage})
	return nil
}
