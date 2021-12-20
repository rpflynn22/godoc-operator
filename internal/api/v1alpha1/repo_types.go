package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RepoSpec defines the desired state of Repo
type RepoSpec struct {
	Repo             string          `json:"repo"`
	GithubPATSecret  GithubPATSecret `json:"githubPATSecret"`
	DNSParent        string          `json:"dnsParent"`
	ALBSecurityGroup string          `json:"albSecurityGroup"`
}

type GithubPATSecret struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

//+kubebuilder:subresource:status
//+kubebuilder:resource:path=repos

// RepoStatus defines the observed state of Repo
type RepoStatus struct {
	State string `json:"state"`
}

//+kubebuilder:object:root=true

// Repo is the Schema for the repos API
type Repo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RepoSpec   `json:"spec,omitempty"`
	Status RepoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RepoList contains a list of Repo
type RepoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Repo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Repo{}, &RepoList{})
}
