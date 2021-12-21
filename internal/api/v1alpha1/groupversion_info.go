// Package v1alpha1 contains API Schema definitions for the godoc v1alpha1 API group
//+kubebuilder:object:generate=true
//+groupName=godoc.rpflynn22.io
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "godoc.rpflynn22.io", Version: "v1alpha1"}
)
