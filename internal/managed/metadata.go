package managed

import (
	"fmt"

	godocApi "github.com/rpflynn22/godoc-operator/internal/api/v1alpha1"
)

const (
	deploymentComponent = "deployment"
	podComponent        = "pod"
	serviceComponent    = "service"
)

func ResourceName(repoResourceName string) string {
	return fmt.Sprintf("godoc-server-%s", repoResourceName)
}

func ResourceLabels(repo *godocApi.Repo, component string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "godoc",
		"app.kubernetes.io/instance":   fmt.Sprintf("godoc-%s", repo.GetName()),
		"app.kubernetes.io/component":  component,
		"app.kubernetes.io/managed-by": "godoc-operator",
	}
}
