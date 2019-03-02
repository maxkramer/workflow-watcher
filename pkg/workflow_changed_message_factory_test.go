package pkg

import (
	"testing"

	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/project-interstellar/workflow-watcher/pkg/types"
)

func TestWorkflowChangedMessageFactory_NewMessage_CreatesMessage(t *testing.T) {
	id := createUUID()
	workflow := v1alpha1.Workflow{
		ObjectMeta: v1.ObjectMeta{Name: id},
		Status: v1alpha1.WorkflowStatus{
			Phase:      v1alpha1.NodePending,
			StartedAt:  v1.Now(),
			FinishedAt: v1.Now(),
		},
	}

	message := WorkflowChangedMessageFactory{}.NewMessage(&workflow, watch.Added).(types.WorkflowChangedMessage)
	assert.Equal(t, message.Id, uuid.MustParse(id))
	assert.Equal(t, message.EventType, watch.Added)
	assert.Equal(t, message.Status, workflow.Status.Phase)
	assert.Equal(t, message.StartedAt, workflow.Status.StartedAt)
	assert.Equal(t, message.FinishedAt, workflow.Status.FinishedAt)
	assert.Empty(t, message.Nodes)

}

func TestWorkflowChangedMessageFactory_NewMessage_UsesEventType(t *testing.T) {
	workflow := v1alpha1.Workflow{
		ObjectMeta: v1.ObjectMeta{Name: createUUID()},
	}

	addedMessage := WorkflowChangedMessageFactory{}.NewMessage(&workflow, watch.Added).(types.WorkflowChangedMessage)
	assert.Equal(t, addedMessage.EventType, watch.Added)

	modifiedMessage := WorkflowChangedMessageFactory{}.NewMessage(&workflow, watch.Modified).(types.WorkflowChangedMessage)
	assert.Equal(t, modifiedMessage.EventType, watch.Modified)

	deletedMessage := WorkflowChangedMessageFactory{}.NewMessage(&workflow, watch.Deleted).(types.WorkflowChangedMessage)
	assert.Equal(t, deletedMessage.EventType, watch.Deleted)
}

func TestWorkflowChangedMessageFactory_NewMessage_ParsesNodes(t *testing.T) {
	nodes := make(map[string]v1alpha1.NodeStatus)

	node1 := v1alpha1.NodeStatus{
		DisplayName: createUUID(),
		Phase:       v1alpha1.NodeSucceeded,
		Message:     "message",
		StartedAt:   v1.Now(),
		FinishedAt:  v1.Now(),
	}
	nodes[node1.DisplayName] = node1

	node2 := v1alpha1.NodeStatus{
		DisplayName: createUUID(),
		Phase:       v1alpha1.NodeFailed,
		Message:     "",
		StartedAt:   v1.Now(),
		FinishedAt:  v1.Now(),
	}
	nodes[node2.DisplayName] = node2

	node3 := v1alpha1.NodeStatus{
		DisplayName: createUUID(),
		Phase:       v1alpha1.NodePending,
		Message:     "message",
		StartedAt:   v1.Now(),
		FinishedAt:  v1.Now(),
	}
	nodes[node3.DisplayName] = node3

	workflow := v1alpha1.Workflow{
		ObjectMeta: v1.ObjectMeta{Name: createUUID()},
		Status: v1alpha1.WorkflowStatus{
			Nodes: nodes,
		},
	}

	message := WorkflowChangedMessageFactory{}.NewMessage(&workflow, watch.Added).(types.WorkflowChangedMessage)
	parsedNodes := message.Nodes
	verifyNodeParsedCorrectly(t, node1, parsedNodes[node1.DisplayName])
	verifyNodeParsedCorrectly(t, node2, parsedNodes[node2.DisplayName])
	verifyNodeParsedCorrectly(t, node3, parsedNodes[node3.DisplayName])
}

func verifyNodeParsedCorrectly(t *testing.T, status v1alpha1.NodeStatus, parsedNode types.WorkflowNode) {
	assert.Equal(t, parsedNode.Id.String(), status.DisplayName)
	assert.Equal(t, parsedNode.Status, status.Phase)
	assert.Equal(t, parsedNode.StatusMessage, status.Message)
	assert.Equal(t, parsedNode.StartedAt, status.StartedAt)
	assert.Equal(t, parsedNode.FinishedAt, status.FinishedAt)
}

func TestWorkflowChangedMessageFactory_NewMessage_IgnoresDAG(t *testing.T) {
	nodes := make(map[string]v1alpha1.NodeStatus)

	node1 := v1alpha1.NodeStatus{
		DisplayName: createUUID(),
		Phase:       v1alpha1.NodeSucceeded,
		Message:     "message",
		StartedAt:   v1.Now(),
		FinishedAt:  v1.Now(),
	}
	nodes["DAG"] = node1

	node2 := v1alpha1.NodeStatus{
		DisplayName: createUUID(),
		Phase:       v1alpha1.NodeFailed,
		Message:     "",
		StartedAt:   v1.Now(),
		FinishedAt:  v1.Now(),
	}
	nodes[node2.DisplayName] = node2

	workflow := v1alpha1.Workflow{
		ObjectMeta: v1.ObjectMeta{Name: createUUID()},
		Status: v1alpha1.WorkflowStatus{
			Nodes: nodes,
		},
	}

	message := WorkflowChangedMessageFactory{}.NewMessage(&workflow, watch.Added).(types.WorkflowChangedMessage)
	assert.Len(t, message.Nodes, 1)

	assert.NotNil(t, message.Nodes[node2.DisplayName])
}

func TestWorkflowChangedMessageFactory_NewMessage_IgnoresWorkflowName(t *testing.T) {
	workflowId := createUUID()

	nodes := make(map[string]v1alpha1.NodeStatus)

	node1 := v1alpha1.NodeStatus{
		DisplayName: createUUID(),
		Phase:       v1alpha1.NodeSucceeded,
		Message:     "message",
		StartedAt:   v1.Now(),
		FinishedAt:  v1.Now(),
	}
	nodes[workflowId] = node1

	node2 := v1alpha1.NodeStatus{
		DisplayName: createUUID(),
		Phase:       v1alpha1.NodeFailed,
		Message:     "",
		StartedAt:   v1.Now(),
		FinishedAt:  v1.Now(),
	}
	nodes[node2.DisplayName] = node2

	workflow := v1alpha1.Workflow{
		ObjectMeta: v1.ObjectMeta{Name: workflowId},
		Status: v1alpha1.WorkflowStatus{
			Nodes: nodes,
		},
	}

	message := WorkflowChangedMessageFactory{}.NewMessage(&workflow, watch.Added).(types.WorkflowChangedMessage)
	assert.Len(t, message.Nodes, 1)

	assert.NotNil(t, message.Nodes[node2.DisplayName])
}

func TestWorkflowChangedMessageFactory_NewMessage_ExtractsArtifacts(t *testing.T) {
	nodes := make(map[string]v1alpha1.NodeStatus)

	node1 := v1alpha1.NodeStatus{
		DisplayName: createUUID(),
	}
	nodes[node1.DisplayName] = node1

	node2 := v1alpha1.NodeStatus{
		DisplayName: createUUID(),
		Outputs: &v1alpha1.Outputs{
			Artifacts: []v1alpha1.Artifact{
				{
					Name: "main-logs",
					ArtifactLocation: v1alpha1.ArtifactLocation{
						S3: &v1alpha1.S3Artifact{
							Key: "key",
							S3Bucket: v1alpha1.S3Bucket{
								Bucket: "bucket",
							},
						},
					},
				},
			},
		},
	}
	nodes[node2.DisplayName] = node2

	node3 := v1alpha1.NodeStatus{
		DisplayName: createUUID(),
		Outputs: &v1alpha1.Outputs{
			Artifacts: []v1alpha1.Artifact{
				{
					Name: "input",
					ArtifactLocation: v1alpha1.ArtifactLocation{
						S3: &v1alpha1.S3Artifact{
							Key: "key",
							S3Bucket: v1alpha1.S3Bucket{
								Bucket: "bucket",
							},
						},
					},
				},
				{
					Name: "output",
					ArtifactLocation: v1alpha1.ArtifactLocation{
						S3: &v1alpha1.S3Artifact{
							Key: "key",
							S3Bucket: v1alpha1.S3Bucket{
								Bucket: "bucket",
							},
						},
					},
				},
			},
		},
	}
	nodes[node3.DisplayName] = node3

	workflow := v1alpha1.Workflow{
		ObjectMeta: v1.ObjectMeta{Name: createUUID()},
		Status: v1alpha1.WorkflowStatus{
			Nodes: nodes,
		},
	}

	message := WorkflowChangedMessageFactory{}.NewMessage(&workflow, watch.Added).(types.WorkflowChangedMessage)
	assert.Nil(t, message.Nodes[node1.DisplayName].Logs)
	assert.Nil(t, message.Nodes[node1.DisplayName].Input)
	assert.Nil(t, message.Nodes[node1.DisplayName].Output)

	assert.Equal(t, "bucket", message.Nodes[node2.DisplayName].Logs.Bucket)
	assert.Equal(t, "key", message.Nodes[node2.DisplayName].Logs.Key)
	assert.Nil(t, message.Nodes[node2.DisplayName].Input)
	assert.Nil(t, message.Nodes[node2.DisplayName].Output)

	assert.Nil(t, message.Nodes[node3.DisplayName].Logs)
	assert.Equal(t, "bucket", message.Nodes[node3.DisplayName].Input.Bucket)
	assert.Equal(t, "key", message.Nodes[node3.DisplayName].Input.Key)
	assert.Equal(t, "bucket", message.Nodes[node3.DisplayName].Output.Bucket)
	assert.Equal(t, "key", message.Nodes[node3.DisplayName].Output.Key)
}

func createUUID() string {
	return uuid.Must(uuid.NewRandom()).String()
}
