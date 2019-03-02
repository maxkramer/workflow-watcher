package pkg

import (
	"strings"

	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/project-interstellar/workflow-watcher/pkg/types"
)

type WorkflowChangedMessageFactory struct{}

func (factory WorkflowChangedMessageFactory) NewMessage(workflow *v1alpha1.Workflow, eventType watch.EventType) interface{} {
	name := workflow.GetObjectMeta().GetName()
	return types.WorkflowChangedMessage{
		Id:         uuid.MustParse(name),
		EventType:  eventType,
		Status:     workflow.Status.Phase,
		Nodes:      factory.parseNodes(name, workflow.Status.Nodes),
		StartedAt:  workflow.Status.StartedAt,
		FinishedAt: workflow.Status.FinishedAt,
	}
}

func (factory *WorkflowChangedMessageFactory) parseNodes(workflowName string, statuses map[string]v1alpha1.NodeStatus) map[string]types.WorkflowNode {
	phases := make(map[string]types.WorkflowNode)

	for nodeName, status := range statuses {
		if !strings.EqualFold(nodeName, "dag") && nodeName != workflowName {
			node := types.WorkflowNode{
				Id:            uuid.MustParse(status.DisplayName),
				Status:        status.Phase,
				StatusMessage: status.Message,
				StartedAt:     status.StartedAt,
				FinishedAt:    status.FinishedAt,
			}

			if status.Outputs != nil && status.Outputs.HasOutputs() {
				factory.parseArtifacts(&node, status.Outputs.Artifacts)
			}

			phases[nodeName] = node
		}
	}

	return phases
}

func (*WorkflowChangedMessageFactory) parseArtifacts(node *types.WorkflowNode, artifacts []v1alpha1.Artifact) {
	for _, artifact := range artifacts {
		if strings.EqualFold(artifact.Name, "main-logs") {
			node.Logs = &types.ArtifactLocation{
				Bucket: artifact.S3.Bucket, Key: artifact.S3.Key,
			}
		} else if strings.EqualFold(artifact.Name, "input") {
			node.Input = &types.ArtifactLocation{
				Bucket: artifact.S3.Bucket, Key: artifact.S3.Key,
			}
		} else if strings.EqualFold(artifact.Name, "output") {
			node.Output = &types.ArtifactLocation{
				Bucket: artifact.S3.Bucket, Key: artifact.S3.Key,
			}
		}
	}
}
