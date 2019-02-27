package pkg

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/sirupsen/logrus"
)

type PubSub struct {
	Log *logrus.Logger
}

func (pubsub PubSub) Append(workflow *v1alpha1.Workflow) *error {
	pubsub.Log.Debugf("PubSub name %s resourceVersion %s", workflow.GetObjectMeta().GetName(),
		workflow.GetObjectMeta().GetResourceVersion())
	return nil
}
