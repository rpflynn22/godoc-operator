package managed

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	godocApi "github.com/rpflynn22/godoc-operator/internal/api/v1alpha1"
)

func ServiceName(name string) string {
	return fmt.Sprintf("godoc-server-%s", name)
}

func Service(godoc *godocApi.Repo) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceName(godoc.Name),
			Namespace: godoc.Namespace,
			Labels:    map[string]string{},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"app":  "godoc-server",
				"repo": Label(godoc.Spec.Repo),
			},
			Ports: []v1.ServicePort{
				{
					Protocol:   v1.ProtocolTCP,
					Port:       6060,
					TargetPort: intstr.FromInt(6060),
				},
			},
		},
	}
}
