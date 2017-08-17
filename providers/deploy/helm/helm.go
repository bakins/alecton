package helm

import "k8s.io/helm/pkg/helm"

func New() *helm.Client {
	return helm.NewClient()
}
