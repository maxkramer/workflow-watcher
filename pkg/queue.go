package pkg

import "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"

type Queue interface {
	Publish(workflow *v1alpha1.Workflow) *error
}
