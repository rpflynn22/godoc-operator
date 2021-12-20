package managed

import (
	appsApi "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	godocApi "github.com/rpflynn22/godoc-operator/internal/api/v1alpha1"
)

func UpdateDeployment(godoc *godocApi.Repo, deployment *appsApi.Deployment) {
	if deployment.ObjectMeta.Labels == nil {
		deployment.ObjectMeta.Labels = make(map[string]string)
	}
	for k, v := range ResourceLabels(godoc, deploymentComponent) {
		deployment.ObjectMeta.Labels[k] = v
	}
	deployment.Spec = appsApi.DeploymentSpec{
		Replicas: intP(2),
		Selector: &metav1.LabelSelector{
			MatchLabels: ResourceLabels(godoc, podComponent),
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: ResourceName(godoc.Name), // name todo
				Labels:       ResourceLabels(godoc, podComponent),
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:            ResourceName(godoc.Name),
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
	}
}

func intP(i int32) *int32 {
	return &i
}
