package internal

import (
	"errors"
	"testing"

	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/project-interstellar/workflow-watcher/mocks"
	"github.com/project-interstellar/workflow-watcher/pkg/metrics"
)

func TestNewWorkflowEventHandler(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	queue := mocks.NewMockQueue(mockController)
	agent := mocks.NewMockAgent(mockController)
	channel := make(chan string)

	handler := NewWorkflowEventHandler(queue, agent, channel)
	assert.Equal(t, handler.queue, queue)
	assert.Equal(t, handler.metricsAgent, agent)
	assert.Equal(t, handler.resourceVersionChannel, channel)
}

func TestWorkflowEventHandler_OnAdd_ShouldIncrementCounter(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	handler := verifyCounterIncremented(mockController, metrics.EventHandlerAdded)
	workflow := v1alpha1.Workflow{}
	handler.OnAdd(&workflow)
}

func TestWorkflowEventHandler_OnAdd_ShouldPushEventToQueue(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	workflow := v1alpha1.Workflow{}
	handler := verifyPushedToQueue(mockController, watch.Added, &workflow)
	handler.OnAdd(&workflow)
}

func TestWorkflowEventHandler_OnAdd_ShouldPushResourceVersion(t *testing.T) {
	handler, channel := verifyPushedResourceVersion()

	workflow := v1alpha1.Workflow{ObjectMeta: v1.ObjectMeta{
		ResourceVersion: "1234",
	}}

	handler.OnAdd(&workflow)
	resourceVersion := <-channel
	assert.Equal(t, resourceVersion, "1234")
}

func TestWorkflowEventHandler_OnUpdate_ShouldIncrementCounter(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	handler := verifyCounterIncremented(mockController, metrics.EventHandlerModified)
	workflow := v1alpha1.Workflow{}
	handler.OnUpdate(&workflow, &workflow)
}

func TestWorkflowEventHandler_OnUpdate_ShouldPushEventToQueue(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	workflow := v1alpha1.Workflow{}
	handler := verifyPushedToQueue(mockController, watch.Modified, &workflow)
	handler.OnUpdate(&workflow, &workflow)
}

func TestWorkflowEventHandler_OnUpdate_ShouldPushResourceVersion(t *testing.T) {
	handler, channel := verifyPushedResourceVersion()

	workflow := v1alpha1.Workflow{ObjectMeta: v1.ObjectMeta{
		ResourceVersion: "1234",
	}}

	handler.OnUpdate(&workflow, &workflow)
	resourceVersion := <-channel
	assert.Equal(t, resourceVersion, "1234")
}

func TestWorkflowEventHandler_OnDelete_ShouldIncrementCounter(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	handler := verifyCounterIncremented(mockController, metrics.EventHandlerDeleted)
	workflow := v1alpha1.Workflow{}
	handler.OnDelete(&workflow)
}

func TestWorkflowEventHandler_OnDelete_ShouldPushEventToQueue(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	workflow := v1alpha1.Workflow{}
	handler := verifyPushedToQueue(mockController, watch.Deleted, &workflow)
	handler.OnDelete(&workflow)
}

func TestWorkflowEventHandler_OnDelete_ShouldPushResourceVersion(t *testing.T) {
	handler, channel := verifyPushedResourceVersion()

	workflow := v1alpha1.Workflow{ObjectMeta: v1.ObjectMeta{
		ResourceVersion: "1234",
	}}

	handler.OnDelete(&workflow)
	resourceVersion := <-channel
	assert.Equal(t, resourceVersion, "1234")
}

func TestWorkflowEventHandler_IncrementsFailureOnQueueFailure(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	agent := mocks.NewMockAgent(mockController)
	queue := mocks.NewMockQueue(mockController)
	handler := NewWorkflowEventHandler(queue, agent, nil)

	workflow := v1alpha1.Workflow{}
	queue.EXPECT().Publish(gomock.Eq(&workflow), gomock.Eq(watch.Added)).Return(errors.New("fail"))

	first := agent.EXPECT().IncrementCounter(gomock.Eq(metrics.EventHandlerAdded), gomock.Eq(float64(1))).Times(1)
	second := agent.EXPECT().IncrementCounter(gomock.Eq(metrics.QueueFailure), gomock.Eq(float64(1))).Times(1)

	gomock.InOrder(first, second)
	handler.OnAdd(&workflow)
}

func TestWorkflowEventHandler_IncrementsFailureOnQueueSuccess(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	agent := mocks.NewMockAgent(mockController)
	queue := mocks.NewMockQueue(mockController)
	handler := NewWorkflowEventHandler(queue, agent, nil)

	workflow := v1alpha1.Workflow{}
	queue.EXPECT().Publish(gomock.Eq(&workflow), gomock.Eq(watch.Added)).Return(nil)

	first := agent.EXPECT().IncrementCounter(gomock.Eq(metrics.EventHandlerAdded), gomock.Eq(float64(1))).Times(1)
	second := agent.EXPECT().IncrementCounter(gomock.Eq(metrics.QueueSuccess), gomock.Eq(float64(1))).Times(1)

	gomock.InOrder(first, second)
	handler.OnAdd(&workflow)
}

func verifyPushedToQueue(mockController *gomock.Controller, eventType watch.EventType,
	workflow *v1alpha1.Workflow) *WorkflowEventHandler {
	queue := mocks.NewMockQueue(mockController)
	handler := NewWorkflowEventHandler(queue, nil, nil)

	queue.EXPECT().Publish(gomock.Eq(workflow), gomock.Eq(eventType))

	return handler
}

func verifyCounterIncremented(mockController *gomock.Controller, metric metrics.Metric) *WorkflowEventHandler {
	agent := mocks.NewMockAgent(mockController)
	handler := NewWorkflowEventHandler(nil, agent, nil)

	agent.EXPECT().IncrementCounter(gomock.Eq(metric), gomock.Eq(float64(1))).Times(1)

	return handler
}

func verifyPushedResourceVersion() (*WorkflowEventHandler, chan string) {
	channel := make(chan string, 1)
	handler := NewWorkflowEventHandler(nil, nil, channel)

	return handler, channel
}
