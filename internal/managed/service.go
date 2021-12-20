package managed

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	godocApi "github.com/rpflynn22/godoc-operator/internal/api/v1alpha1"
)

func UpdateService(repo *godocApi.Repo, service *v1.Service) {
	if service.ObjectMeta.Labels == nil {
		service.ObjectMeta.Labels = make(map[string]string)
	}
	for k, v := range ResourceLabels(repo, serviceComponent) {
		service.ObjectMeta.Labels[k] = v
	}
	service.Spec = v1.ServiceSpec{
		Selector: ResourceLabels(repo, podComponent),
		Ports: []v1.ServicePort{
			{
				Protocol:   v1.ProtocolTCP,
				Port:       6060,
				TargetPort: intstr.FromInt(6060),
			},
		},
	}
}
