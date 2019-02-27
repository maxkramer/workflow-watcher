package pkg

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/google/uuid"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type WorkflowChangedMessage struct {
	Id         uuid.UUID               `json:"id"`
	EventType  watch.EventType         `json:"eventType"`
	Status     v1alpha1.NodePhase      `json:"status"`
	Nodes      map[string]WorkflowNode `json:"nodes"`
	StartedAt  meta.Time               `json:"startedAt"`
	FinishedAt meta.Time               `json:"finishedAt"`
}
