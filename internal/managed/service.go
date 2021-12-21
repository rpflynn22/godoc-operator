package managed

import (
	"fmt"

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
	if service.ObjectMeta.Annotations == nil {
		service.ObjectMeta.Annotations = make(map[string]string)
	}
	for k, v := range svcAnnotations(repo) {
		service.ObjectMeta.Annotations[k] = v
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

func svcAnnotations(repo *godocApi.Repo) map[string]string {
	return map[string]string{
		"alb.ingress.kubernetes.io/backend-protocol": "HTTP",
		"alb.ingress.kubernetes.io/listen-ports":     "[{\"HTTP\":80}]",
		"alb.ingress.kubernetes.io/scheme":           "internal",
		"alb.ingress.kubernetes.io/security-groups":  repo.Spec.ALBConfig.AlbSg,
		"alb.ingress.kubernetes.io/target-type":      "ip",
		"external-dns.alpha.kubernetes.io/hostname":  dnsHostname(repo),
		"kubernetes.io/ingress.class":                "alb",
	}
}

func dnsHostname(repo *godocApi.Repo) string {
	return fmt.Sprintf("%s.%s", repo.GetName(), repo.Spec.ALBConfig.DNSParent)
}
