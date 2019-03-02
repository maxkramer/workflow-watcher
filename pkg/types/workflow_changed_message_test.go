package types

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/google/uuid"
	"github.com/project-interstellar/workflow-watcher/internal/test_utils"
	"github.com/stretchr/testify/assert"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"testing"
	"time"
)

func TestShouldRenderCorrectJson_WorkflowChangedMessage(t *testing.T) {
	message := WorkflowChangedMessage{
		Id:         uuid.Must(uuid.NewRandom()),
		EventType:  watch.Added,
		Status:     v1alpha1.NodeRunning,
		Nodes:      make(map[string]WorkflowNode),
		StartedAt:  meta.Now(),
		FinishedAt: meta.Now(),
	}

	res := test_utils.MarshalToMap(t, message)

	assert.Equal(t, message.Id.String(), res["id"])
	assert.Equal(t, string(message.EventType), res["eventType"])
	assert.Equal(t, string(message.Status), res["status"])
	assert.Equal(t, message.StartedAt.UTC().Format(time.RFC3339), res["startedAt"])
	assert.Equal(t, message.FinishedAt.UTC().Format(time.RFC3339), res["finishedAt"])
	assert.Empty(t, res["nodes"])
}
