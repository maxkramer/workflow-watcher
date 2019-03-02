package internal

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/project-interstellar/workflow-watcher/pkg/queue"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/watch"
	"strings"
)

type WorkflowEventHandler struct {
	queue                  queue.Queue
	statsd                 *statsd.Client
	resourceVersionChannel chan string
}

func NewWorkflowEventHandler(queue queue.Queue, statsd *statsd.Client, resourceVersionChannel chan string) *WorkflowEventHandler {
	return &WorkflowEventHandler{
		queue:                  queue,
		statsd:                 statsd,
		resourceVersionChannel: resourceVersionChannel,
	}
}

func (handler *WorkflowEventHandler) OnAdd(obj interface{}) {
	go handler.handleWorkflowChange(watch.Added, obj.(*v1alpha1.Workflow))
}

func (handler *WorkflowEventHandler) OnUpdate(oldObj, newObj interface{}) {
	go handler.handleWorkflowChange(watch.Modified, newObj.(*v1alpha1.Workflow))
}

func (handler *WorkflowEventHandler) OnDelete(obj interface{}) {
	go handler.handleWorkflowChange(watch.Deleted, obj.(*v1alpha1.Workflow))
}

func (handler *WorkflowEventHandler) handleWorkflowChange(eventType watch.EventType, workflow *v1alpha1.Workflow) {
	log.Debugf("Received event with type %s", eventType)
	log.Debugf("Object's name %s resourceVersion %s", workflow.GetObjectMeta().GetName(),
		workflow.GetObjectMeta().GetResourceVersion())

	handler.incrementCounter(fmt.Sprintf("event-handler.%s", strings.ToLower(string(eventType))))
	handler.pushWorkflowToQueue(workflow, eventType)
}

func (handler *WorkflowEventHandler) incrementCounter(name string) {
	err := handler.statsd.Count(name, 1, nil, 1)
	if err != nil {
		log.Errorf("Error incrementing %s counter", name)
	}
}

func (handler *WorkflowEventHandler) pushWorkflowToQueue(workflow *v1alpha1.Workflow, eventType watch.EventType) {
	err := handler.queue.Publish(workflow, eventType)
	if err != nil {
		log.Error("Failed to add event to queue ", err)
		handler.incrementCounter("queue.failure")
	} else {
		handler.incrementCounter("queue.success")
	}

	handler.resourceVersionChannel <- workflow.ResourceVersion
}
