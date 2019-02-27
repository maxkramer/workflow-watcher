package internal

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/project-interstellar/workflow-watcher/pkg"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/watch"
)

type WorkflowEventHandler struct {
	Log   *logrus.Logger
	Queue pkg.Queue
}

func (handler WorkflowEventHandler) OnAdd(obj interface{}) {

	handler.handleWorkflowChange(watch.Added, obj.(*v1alpha1.Workflow))
}

func (handler WorkflowEventHandler) OnUpdate(oldObj, newObj interface{}) {
	handler.handleWorkflowChange(watch.Modified, newObj.(*v1alpha1.Workflow))
}

func (handler WorkflowEventHandler) OnDelete(obj interface{}) {
	handler.handleWorkflowChange(watch.Deleted, obj.(*v1alpha1.Workflow))
}

func (handler WorkflowEventHandler) handleWorkflowChange(eventType watch.EventType, workflow *v1alpha1.Workflow) {
	handler.Log.Debugf("Received event with type %s", eventType)
	handler.Log.Debugf("Object's name %s resourceVersion %s", workflow.GetObjectMeta().GetName(),
		workflow.GetObjectMeta().GetResourceVersion())

	go handler.Queue.Append(workflow)
}
