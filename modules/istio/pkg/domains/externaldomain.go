package domains

import (
	"fmt"

	riov1 "github.com/rancher/rio/pkg/apis/rio.cattle.io/v1"
	"github.com/rancher/rio/pkg/constants"
)

// IsPublic returns true if any port of the svc is exposed
func IsPublic(svc *riov1.Service) bool {
	public := false
	for _, port := range svc.Spec.Ports {
		if port.Expose != nil && *port.Expose {
			public = true
			break
		}
	}
	return public
}

func GetPublicGateway(systemNamespace string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", constants.RioGateway, systemNamespace)
}

func GetExternalDomain(name, namespace, clusterDomain string) string {
	return fmt.Sprintf("%s-%s.%s", name, namespace, clusterDomain)
}
