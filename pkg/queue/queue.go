package queue

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"k8s.io/apimachinery/pkg/watch"
)

type Queue interface {
	Publish(*v1alpha1.Workflow, watch.EventType) error
}
