package managed

import (
	"fmt"

	netApi "k8s.io/api/networking/v1"

	godocApi "github.com/rpflynn22/godoc-operator/internal/api/v1alpha1"
)

func UpdateIngress(repo *godocApi.Repo, ingress *netApi.Ingress) {
	if ingress.ObjectMeta.Labels == nil {
		ingress.ObjectMeta.Labels = make(map[string]string)
	}
	if ingress.ObjectMeta.Annotations == nil {
		ingress.ObjectMeta.Annotations = make(map[string]string)
	}

	for k, v := range ResourceLabels(repo, ingressComponent) {
		ingress.ObjectMeta.Labels[k] = v
	}
	for k, v := range ingAnnotations(repo) {
		ingress.ObjectMeta.Annotations[k] = v
	}

	ingress.Spec = netApi.IngressSpec{
		Rules: []netApi.IngressRule{
			ingRule(repo),
		},
	}
}

func ingAnnotations(repo *godocApi.Repo) map[string]string {
	return map[string]string{
		"alb.ingress.kubernetes.io/backend-protocol": "HTTP",
		"alb.ingress.kubernetes.io/listen-ports":     "'[{\"HTTP\":80}]'",
		"alb.ingress.kubernetes.io/scheme":           "internal",
		"alb.ingress.kubernetes.io/security-groups":  repo.Spec.ALBSecurityGroup,
		"alb.ingress.kubernetes.io/target-type":      "ip",
		"external-dns.alpha.kubernetes.io/hostname":  ingHostname(repo),
		"kubernetes.io/ingress.class":                "alb",
	}
}

func ingHostname(repo *godocApi.Repo) string {
	return fmt.Sprintf("%s.%s", repo.GetName(), repo.Spec.DNSParent)
}

func ingRule(repo *godocApi.Repo) netApi.IngressRule {
	return netApi.IngressRule{
		Host: ingHostname(repo),
		IngressRuleValue: netApi.IngressRuleValue{
			HTTP: &netApi.HTTPIngressRuleValue{
				Paths: []netApi.HTTPIngressPath{
					{
						Path:     ingPath(repo),
						PathType: ingPathTypeFix(),
						Backend: netApi.IngressBackend{
							Service: &netApi.IngressServiceBackend{
								Name: ResourceName(repo.GetName()),
								Port: netApi.ServiceBackendPort{
									Number: 6060,
								},
							},
						},
					},
				},
			},
		},
	}
}

func ingPath(repo *godocApi.Repo) string {
	return fmt.Sprintf("/pkg/%s", repo.Spec.Repo)
}

func ingPathTypeFix() *netApi.PathType {
	result := netApi.PathTypePrefix
	return &result
}
