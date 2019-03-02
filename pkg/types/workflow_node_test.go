package types

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/google/uuid"
	"github.com/project-interstellar/workflow-watcher/internal/test_utils"
	"github.com/stretchr/testify/assert"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func TestShouldRenderCorrectJson_WorkflowNode(t *testing.T) {
	message := WorkflowNode{
		Id:         uuid.Must(uuid.NewRandom()),
		Status:     v1alpha1.NodeRunning,
		StatusMessage: "message",
		StartedAt:  meta.Now(),
		FinishedAt: meta.Now(),
		Logs: &ArtifactLocation{ Bucket: "bucket", Key: "logs"},
		Input: &ArtifactLocation{ Bucket: "bucket", Key: "input"},
		Output: &ArtifactLocation{ Bucket: "bucket", Key: "output"},
	}

	res := test_utils.MarshalToMap(t, message)

	assert.Equal(t, message.Id.String(), res["id"])
	assert.Equal(t, string(message.Status), res["status"])
	assert.Equal(t, message.StatusMessage, res["statusMessage"])
	assert.Equal(t, message.StartedAt.UTC().Format(time.RFC3339), res["startedAt"])
	assert.Equal(t, message.FinishedAt.UTC().Format(time.RFC3339), res["finishedAt"])

	logsLocation := res["logs"].(map[string]interface{})
	inputLocation := res["input"].(map[string]interface{})
	outputLocation := res["output"].(map[string]interface{})

	assert.Equal(t, message.Logs.Key, logsLocation["key"])
	assert.Equal(t, message.Logs.Bucket, logsLocation["bucket"])

	assert.Equal(t, message.Input.Key, inputLocation["key"])
	assert.Equal(t, message.Input.Bucket, inputLocation["bucket"])

	assert.Equal(t, message.Output.Key, outputLocation["key"])
	assert.Equal(t, message.Output.Bucket, outputLocation["bucket"])

}
