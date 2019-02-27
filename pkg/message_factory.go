package pkg

import "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"

type MessageFactory interface {
	NewMessage(workflow *v1alpha1.Workflow) interface{}
}