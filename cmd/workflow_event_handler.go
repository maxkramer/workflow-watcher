package main

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/watch"
)

type WorkflowEventHandler struct{}

func (WorkflowEventHandler) OnAdd(obj interface{}) {
	handleWorkflowChange(watch.Added, obj.(*v1alpha1.Workflow))
}

func (WorkflowEventHandler) OnUpdate(oldObj, newObj interface{}) {
	handleWorkflowChange(watch.Modified, newObj.(*v1alpha1.Workflow))
}

func (WorkflowEventHandler) OnDelete(obj interface{}) {
	handleWorkflowChange(watch.Deleted, obj.(*v1alpha1.Workflow))
}

func handleWorkflowChange(eventType watch.EventType, workflow *v1alpha1.Workflow) {
	log.Debugf("Received event with type %s", eventType)
	log.Debugf("Object's name %s resourceVersion %s", workflow.GetObjectMeta().GetName(),
		workflow.GetObjectMeta().GetResourceVersion())
}
