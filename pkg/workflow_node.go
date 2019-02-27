package pkg

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/google/uuid"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WorkflowNode struct {
	Id            uuid.UUID          `json:"id"`
	Status        v1alpha1.NodePhase `json:"status"`
	StatusMessage string             `json:"statusMessage"`
	StartedAt     meta.Time          `json:"startedAt"`
	FinishedAt    meta.Time          `json:"finishedAt"`
	Logs          *ArtifactLocation  `json:"logs"`
	Input         *ArtifactLocation  `json:"input"`
	Output        *ArtifactLocation  `json:"output"`
}
