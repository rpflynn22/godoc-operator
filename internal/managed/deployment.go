package managed

import (
	"fmt"
	"strings"

	appsApi "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	godocApi "github.com/rpflynn22/godoc-operator/internal/api/v1alpha1"
)

func DeploymentName(name string) string {
	return fmt.Sprintf("godoc-server-%s", name)
}

func Label(repo string) string {
	return strings.Join(strings.Split(repo, "/")[1:], "")
}

func Deployment(godoc *godocApi.Repo) *appsApi.Deployment {
	return &appsApi.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DeploymentName(godoc.Name),
			Namespace: godoc.Namespace,
			Labels:    map[string]string{},
		},
		Spec: appsApi.DeploymentSpec{
			Replicas: intP(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  "godoc-server",
					"repo": Label(godoc.Spec.Repo),
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: DeploymentName(godoc.Name), // name todo
					Labels: map[string]string{
						"app":  "godoc-server",
						"repo": Label(godoc.Spec.Repo),
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            DeploymentName(godoc.Name),
							Image:           "rpflynn22/godoc-server:latest",
							ImagePullPolicy: v1.PullNever,
							Ports: []v1.ContainerPort{
								{
									ContainerPort: 6060,
								},
							},
							Env: []v1.EnvVar{
								{
									Name:  "GO_REPO",
									Value: godoc.Spec.Repo,
								},
								{
									Name: "GH_PAT",
									ValueFrom: &v1.EnvVarSource{
										SecretKeyRef: &v1.SecretKeySelector{
											LocalObjectReference: v1.LocalObjectReference{
												Name: godoc.Spec.GithubPATSecret.Name,
											},
											Key: godoc.Spec.GithubPATSecret.Key,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func intP(i int32) *int32 {
	return &i
}
