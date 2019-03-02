package internal

import (
	"fmt"
	"strings"

	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/project-interstellar/workflow-watcher/pkg/metrics"
	"github.com/project-interstellar/workflow-watcher/pkg/queue"
)

type WorkflowEventHandler struct {
	queue                  queue.Queue
	metricsAgent           metrics.Agent
	resourceVersionChannel chan string
}

func NewWorkflowEventHandler(queue queue.Queue, metricsAgent metrics.Agent,
	resourceVersionChannel chan string) *WorkflowEventHandler {
	return &WorkflowEventHandler{
		queue:                  queue,
		metricsAgent:           metricsAgent,
		resourceVersionChannel: resourceVersionChannel,
	}
}

func (handler *WorkflowEventHandler) OnAdd(obj interface{}) {
	handler.handleWorkflowChange(watch.Added, obj.(*v1alpha1.Workflow))
}

func (handler *WorkflowEventHandler) OnUpdate(oldObj, newObj interface{}) {
	handler.handleWorkflowChange(watch.Modified, newObj.(*v1alpha1.Workflow))
}

func (handler *WorkflowEventHandler) OnDelete(obj interface{}) {
	handler.handleWorkflowChange(watch.Deleted, obj.(*v1alpha1.Workflow))
}

func (handler *WorkflowEventHandler) handleWorkflowChange(eventType watch.EventType, workflow *v1alpha1.Workflow) {
	log.Debugf("Received event with type %s", eventType)
	log.Debugf("Object's name %s resourceVersion %s", workflow.GetObjectMeta().GetName(),
		workflow.GetObjectMeta().GetResourceVersion())

	handler.incrementCounter(metrics.Metric(fmt.Sprintf("event-handler.%s", strings.ToLower(string(eventType)))))
	handler.pushWorkflowToQueue(workflow, eventType)
	handler.pushResourceVersion(workflow)
}

func (handler *WorkflowEventHandler) incrementCounter(name metrics.Metric) {
	if handler.metricsAgent != nil {
		err := handler.metricsAgent.IncrementCounter(name, 1)
		if err != nil {
			log.Errorf("Error incrementing %s counter", name)
		}
	}
}

func (handler *WorkflowEventHandler) pushResourceVersion(workflow *v1alpha1.Workflow) {
	if handler.resourceVersionChannel != nil {
		handler.resourceVersionChannel <- workflow.ResourceVersion
	}
}

func (handler *WorkflowEventHandler) pushWorkflowToQueue(workflow *v1alpha1.Workflow, eventType watch.EventType) {
	if handler.queue != nil {
		err := handler.queue.Publish(workflow, eventType)

		if err != nil {
			log.Error("Failed to add event to queue ", err)
		}

		if handler.metricsAgent != nil {
			if err != nil {
				handler.incrementCounter(metrics.QueueFailure)
			} else {
				handler.incrementCounter(metrics.QueueSuccess)
			}
		}
	}
}
