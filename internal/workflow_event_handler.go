package internal

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/project-interstellar/workflow-watcher/pkg"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/watch"
)

type WorkflowEventHandler struct {
	Log    *logrus.Logger
	Queue  pkg.Queue
	Statsd *statsd.Client
}

func (handler WorkflowEventHandler) OnAdd(obj interface{}) {
	go handler.handleWorkflowChange(watch.Added, obj.(*v1alpha1.Workflow))
}

func (handler WorkflowEventHandler) OnUpdate(oldObj, newObj interface{}) {
	go handler.handleWorkflowChange(watch.Modified, newObj.(*v1alpha1.Workflow))
}

func (handler WorkflowEventHandler) OnDelete(obj interface{}) {
	go handler.handleWorkflowChange(watch.Deleted, obj.(*v1alpha1.Workflow))
}

func (handler WorkflowEventHandler) handleWorkflowChange(eventType watch.EventType, workflow *v1alpha1.Workflow) {
	handler.Log.Debugf("Received event with type %s", eventType)
	handler.Log.Debugf("Object's name %s resourceVersion %s", workflow.GetObjectMeta().GetName(),
		workflow.GetObjectMeta().GetResourceVersion())

	handler.incrementCounter(fmt.Sprintf("event-handler.%s", eventType))
	handler.pushWorkflowToQueue(workflow)
}

func (handler WorkflowEventHandler) incrementCounter(name string) {
	err := handler.Statsd.Count(name, 1, nil, 1)
	if err != nil {
		handler.Log.Errorf("Error incrementing %s counter", name)
	}
}

func (handler WorkflowEventHandler) pushWorkflowToQueue(workflow *v1alpha1.Workflow) {
	err := handler.Queue.Publish(workflow)
	if *err != nil {
		handler.Log.Error("Failed to add event to queue ", *err)
		handler.incrementCounter("queue.failure")
	} else {
		handler.incrementCounter("queue.success")
	}
}
