package pkg

import (
	"context"
	"encoding/json"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/sirupsen/logrus"
)

type PubSub struct {
	Log *logrus.Logger
	ctx context.Context
	messageFactory WorkflowChangedMessageFactory
}


func (pubsub PubSub) Append(workflow *v1alpha1.Workflow) *error {
	message := pubsub.messageFactory.NewMessage(workflow)

	res, _ := json.Marshal(message)
	pubsub.Log.Debug(string(res))

	// Sets your Google Cloud Platform project ID.
	//projectID := "interstellar-staging-env"
	//
	//// Creates a client.
	//client, err := gpubsub.NewClient(pubsub.ctx, projectID)
	//if err != nil {
	//	return &err
	//}
	//
	//topic := client.Topic("")
	//
	//gpubsub.Message{Data:}
	//topic.Publish(pubsub.ctx)
	return nil
}
