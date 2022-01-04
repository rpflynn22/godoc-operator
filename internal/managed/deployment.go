package managed

import (
	appsApi "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	godocApi "github.com/rpflynn22/godoc-operator/internal/api/v1alpha1"
)

func UpdateDeployment(repo *godocApi.Repo, deployment *appsApi.Deployment) {
	if deployment.ObjectMeta.Labels == nil {
		deployment.ObjectMeta.Labels = make(map[string]string)
	}
	for k, v := range ResourceLabels(repo, deploymentComponent) {
		deployment.ObjectMeta.Labels[k] = v
	}
	deployment.Spec = appsApi.DeploymentSpec{
		Replicas: intP(2),
		Selector: &metav1.LabelSelector{
			MatchLabels: ResourceLabels(repo, podComponent),
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: ResourceName(repo.Name),
				Labels:       ResourceLabels(repo, podComponent),
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:            ResourceName(repo.Name),
						Image:           "rpflynn22/godoc-server:latest",
						ImagePullPolicy: v1.PullNever, // todo: change with real images
						Ports: []v1.ContainerPort{
							{
								ContainerPort: 6060,
							},
						},
						Env: []v1.EnvVar{
							{
								Name:  "GO_REPO",
								Value: repo.Spec.GoConfig.Repo,
							},
							{
								Name:  "GH_USER",
								Value: repo.Spec.GHCreds.Username,
							},
							{
								Name: "GH_PAT",
								ValueFrom: &v1.EnvVarSource{
									SecretKeyRef: &v1.SecretKeySelector{
										LocalObjectReference: v1.LocalObjectReference{
											Name: repo.Spec.GHCreds.PATSecret.Name,
										},
										Key: repo.Spec.GHCreds.PATSecret.Key,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if repo.Spec.GoConfig.GoPrivate != "" {
		deployment.Spec.Template.Spec.Containers[0].Env = append(
			deployment.Spec.Template.Spec.Containers[0].Env,
			v1.EnvVar{
				Name:  "GOPRIVATE_PATTERN",
				Value: repo.Spec.GoConfig.GoPrivate,
			},
		)
	}
	if repo.Spec.GoConfig.ModuleVersion != "" {
		deployment.Spec.Template.Spec.Containers[0].Env = append(
			deployment.Spec.Template.Spec.Containers[0].Env,
			v1.EnvVar{
				Name:  "MOD_VERSION",
				Value: repo.Spec.GoConfig.ModuleVersion,
			},
		)
	}
}

func intP(i int32) *int32 {
	return &i
}
