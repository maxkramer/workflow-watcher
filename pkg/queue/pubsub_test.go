package queue

import (
	"context"
	"testing"

	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/project-interstellar/workflow-watcher/pkg/types"
)

type MockMessageFactory struct {
	mock.Mock
}

func (mf MockMessageFactory) NewMessage(workflow *v1alpha1.Workflow, eventType watch.EventType) interface{} {
	return types.WorkflowChangedMessage{}
}

func TestNewPubSub(t *testing.T) {
	ctx := context.Background()
	messageFactory := MockMessageFactory{}
	projectId := "project-id"
	topicId := "topic"

	pubsub := NewPubSub(ctx, messageFactory, projectId, topicId)

	assert.Equal(t, ctx, pubsub.context)
	assert.Equal(t, messageFactory, pubsub.messageFactory)
	assert.Equal(t, projectId, pubsub.projectId)
	assert.Equal(t, topicId, pubsub.topicId)

	assert.Nil(t, pubsub.topic)
}

func TestPublishShouldConnectToTopic(t *testing.T) {
	ctx := context.Background()
	messageFactory := MockMessageFactory{}
	projectId := "project-id"
	topicId := "topic"

	pubsub := NewPubSub(ctx, messageFactory, projectId, topicId)
	assert.Nil(t, pubsub.Publish(&v1alpha1.Workflow{}, watch.Added))

}

func TestPublishShouldReturnError(t *testing.T) {

}
