package util

import (
	"github.com/jhonnli/container-orchestration-api/model/k8s"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CastToK8sListOptions(options k8s.ListOptions) v1.ListOptions {
	return v1.ListOptions{
		LabelSelector:        options.LabelSelector,
		FieldSelector:        options.FieldSelector,
		IncludeUninitialized: options.IncludeUninitialized,
		Watch:                options.Watch,
		ResourceVersion:      options.ResourceVersion,
		TimeoutSeconds:       options.TimeoutSeconds,
		Limit:                options.Limit,
		Continue:             options.Continue,
	}
}
