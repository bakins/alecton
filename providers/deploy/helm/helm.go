package helm

import (
	"fmt"

	"github.com/bakins/alecton"
	"golang.org/x/net/context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
)

type Helm struct {
	client *helm.Client
}

func New(map[string]interface{}) *Helm {
	return &Helm{client: helm.NewClient()}
}

func (h *Helm) InstallRelease(ctx context.Context, req *services.InstallReleaseRequest) (*release.Release, error) {

	// helm client doesn't directly expose passing in the
	// grpc messages, so we have to use the exported functions.
	res, err := h.client.InstallReleaseFromChart(
		req.Chart,
		req.Namespace,
		helm.InstallTimeout(req.Timeout),
		helm.ReleaseName(req.Name),
		helm.ValueOverrides([]byte(req.Values.Raw)),
	)

	if err != nil {
		return nil, err
	}

	return res.Release, nil
}

func (h *Helm) UpdateRelease(ctx context.Context, req *services.UpdateReleaseRequest) (*release.Release, error) {
	res, err := h.client.UpdateReleaseFromChart(
		req.Name,
		req.Chart,
		helm.ResetValues(true),
		helm.UpdateValueOverrides([]byte(req.Values.Raw)),
		helm.UpgradeTimeout(req.Timeout),
	)

	if err != nil {
		return nil, err
	}

	return res.Release, nil
}

func (h *Helm) RollbackRelease(ctx context.Context, req *services.RollbackReleaseRequest) (*release.Release, error) {
	res, err := h.client.RollbackRelease(
		req.Name,
		helm.RollbackTimeout(req.Timeout),
		helm.RollbackVersion(req.Version),
	)

	if err != nil {
		return nil, err
	}

	return res.Release, nil
}

func (h *Helm) ReleaseHistory(ctx context.Context, req *services.GetHistoryRequest) ([]*release.Release, error) {
	res, err := h.client.ReleaseHistory(
		req.Name,
		helm.WithMaxHistory(req.Max),
	)

	if err != nil {
		return nil, err
	}

	releases := []*release.Release{}

	for _, r := range res.Releases {
		releases = append(releases, r)
	}
	return releases, nil
}

func provider(config map[string]interface{}) (alecton.DeployProvider, error) {
	return New(config), nil
}

func init() {
	alecton.RegisterDeployProvider("helm", provider)
}

// getKubeClient creates a Kubernetes config and client for a given kubeconfig context.
func getKubeClient(context string) (*rest.Config, kubernetes.Interface, error) {
	config, err := configForContext(context)
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get Kubernetes client: %s", err)
	}
	return config, client, nil
}

func configForContext(context string) (*rest.Config, error) {
	// TODO: pass in kubeconfig and actually use context
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", context, err)
	}
	return config, nil
}
