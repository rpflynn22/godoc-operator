package managed

import (
	"fmt"

	godocApi "github.com/rpflynn22/godoc-operator/internal/api/v1alpha1"
	netApi "k8s.io/api/networking/v1"
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
	if repo.Spec.AWSAlbConfig != nil {
		out := map[string]string{
			"alb.ingress.kubernetes.io/backend-protocol": "HTTP",
			"alb.ingress.kubernetes.io/scheme":           "internal",
			"alb.ingress.kubernetes.io/security-groups":  repo.Spec.AWSAlbConfig.SecurityGroup,
			"alb.ingress.kubernetes.io/target-type":      "ip",
			"alb.ingress.kubernetes.io/healthcheck-path": "/pkg/",
			"external-dns.alpha.kubernetes.io/hostname":  ingHostname(repo),
			"kubernetes.io/ingress.class":                "alb",
		}
		if repo.Spec.AWSAlbConfig.CertificateArn != "" {
			out["alb.ingress.kubernetes.io/certificate-arn"] =
				repo.Spec.AWSAlbConfig.CertificateArn
			out["alb.ingress.kubernetes.io/listen-ports"] = "[{\"HTTPS\": 443}]"
		}
		return out
	}
	return make(map[string]string)
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
						Path:     "/*",
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

func ingPathTypeFix() *netApi.PathType {
	result := netApi.PathTypeImplementationSpecific
	return &result
}
