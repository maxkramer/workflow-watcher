package pkg

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"k8s.io/apimachinery/pkg/watch"
)

type MessageFactory interface {
	NewMessage(*v1alpha1.Workflow, watch.EventType) interface{}
}