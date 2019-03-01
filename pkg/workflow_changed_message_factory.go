package pkg

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/google/uuid"
	"github.com/project-interstellar/workflow-watcher/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

type WorkflowChangedMessageFactory struct{}

func (factory WorkflowChangedMessageFactory) NewMessage(workflow *v1alpha1.Workflow) interface{} {
	return types.WorkflowChangedMessage{
		Id:         uuid.MustParse(workflow.GetObjectMeta().GetName()),
		EventType:  watch.Added,
		Status:     workflow.Status.Phase,
		Nodes:      factory.parseNodes(workflow.GetObjectMeta().GetName(), workflow.Status.Nodes),
		StartedAt:  workflow.Status.StartedAt,
		FinishedAt: workflow.Status.FinishedAt,
	}
}

func (WorkflowChangedMessageFactory) parseArtifact(artifactType types.ArtifactType, artifacts []v1alpha1.Artifact) *types.ArtifactLocation {
	var artifactLocation *v1alpha1.S3Artifact
	for _, artifact := range artifacts {
		if artifact.Name == "main-logs" && artifactType == types.Logs {
			artifactLocation = artifact.S3
			break
		} else if artifact.Name == "input" && artifactType == types.Input {
			artifactLocation = artifact.S3
			break
		} else if artifact.Name == "output" && artifactType == types.Output {
			artifactLocation = artifact.S3
			break
		}
	}

	if artifactLocation == nil {
		return nil
	}

	return &types.ArtifactLocation{Bucket: artifactLocation.Bucket, Key: artifactLocation.Key}
}

func (factory WorkflowChangedMessageFactory) parseNodes(workflowName string, statuses map[string]v1alpha1.NodeStatus) map[string]types.WorkflowNode {
	phases := make(map[string]types.WorkflowNode)

	for nodeName, status := range statuses {
		if nodeName != "DAG" && nodeName != "dag" && nodeName != workflowName {
			node := types.WorkflowNode{
				Id:            uuid.MustParse(status.DisplayName),
				Status:        status.Phase,
				StatusMessage: status.Message,
				StartedAt:     status.StartedAt,
				FinishedAt:    status.FinishedAt,
			}

			if status.Outputs != nil && status.Outputs.HasOutputs() {
				node.Logs = factory.parseArtifact(types.Logs, status.Outputs.Artifacts)
				node.Input = factory.parseArtifact(types.Input, status.Outputs.Artifacts)
				node.Output = factory.parseArtifact(types.Output, status.Outputs.Artifacts)
			}

			phases[nodeName] = node
		}
	}

	return phases
}
