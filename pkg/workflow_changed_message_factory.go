package pkg

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/watch"
)

type WorkflowChangedMessageFactory struct{}

func (factory WorkflowChangedMessageFactory) NewMessage(workflow *v1alpha1.Workflow) interface{} {
	return WorkflowChangedMessage{
		Id:         uuid.MustParse(workflow.GetObjectMeta().GetName()),
		EventType:  watch.Added,
		Status:     workflow.Status.Phase,
		Nodes:      factory.parseNodes(workflow.GetObjectMeta().GetName(), workflow.Status.Nodes),
		StartedAt:  workflow.Status.StartedAt,
		FinishedAt: workflow.Status.FinishedAt,
	}
}

func (WorkflowChangedMessageFactory) parseArtifact(artifactType ArtifactType, artifacts []v1alpha1.Artifact) *ArtifactLocation {
	var artifactLocation *v1alpha1.S3Artifact
	for _, artifact := range artifacts {
		if artifact.Name == "main-logs" && artifactType == Logs {
			artifactLocation = artifact.S3
			break
		} else if artifact.Name == "input" && artifactType == Input {
			artifactLocation = artifact.S3
			break
		} else if artifact.Name == "output" && artifactType == Output {
			artifactLocation = artifact.S3
			break
		}
	}

	if artifactLocation == nil {
		return nil
	}

	return &ArtifactLocation{Bucket: artifactLocation.Bucket, Key: artifactLocation.Key}
}

func (factory WorkflowChangedMessageFactory) parseNodes(workflowName string, statuses map[string]v1alpha1.NodeStatus) map[string]WorkflowNode {
	phases := make(map[string]WorkflowNode)

	for nodeName, status := range statuses {
		if nodeName != "DAG" && nodeName != "dag" && nodeName != workflowName {
			node := WorkflowNode{
				Id:            uuid.MustParse(status.DisplayName),
				Status:        status.Phase,
				StatusMessage: status.Message,
				StartedAt:     status.StartedAt,
				FinishedAt:    status.FinishedAt,
			}

			if status.Outputs.HasOutputs() {
				node.Logs = factory.parseArtifact(Logs, status.Outputs.Artifacts)
				node.Input = factory.parseArtifact(Input, status.Outputs.Artifacts)
				node.Output = factory.parseArtifact(Output, status.Outputs.Artifacts)
			}

			phases[nodeName] = node
		}
	}

	return phases
}
