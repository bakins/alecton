package helm

import (
	"github.com/bakins/alecton"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
)

type clusterConfig struct {
	Address string
	// TODO: tls?
	client services.ReleaseServiceClient
}

type helmConfig struct {
	Clusters map[string]clusterConfig
}

// Helm is a simple client for tiller
type Helm struct {
	helmConfig
}

func newGrpcClient(address string) (services.ReleaseServiceClient, error) {
	// TODO: timeout
	conn, err := grpc.DialContext(context.Background(), address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return services.NewReleaseServiceClient(conn), nil
}

// New creates a deploy provider using help
func New(c *helmConfig) (*Helm, error) {
	for k, v := range c.Clusters {
		if v.Address == "" {
			return nil, errors.Errorf("address is required for cluster %s", k)
			client, err := newGrpcClient(v.Address)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create client for cluster %s", k)
			}
			v.client = client
		}
	}
	return &Helm{*c}, nil
}

func (h *Helm) getClient(cluster string) (services.ReleaseServiceClient, error) {
	c, ok := h.Clusters[cluster]
	if !ok {
		return nil, errors.Errorf("unknown cluster %s", cluster)
	}
	return c, nil
}

func (h *Helm) InstallRelease(ctx context.Context, cluster string, req *services.InstallReleaseRequest) (*release.Release, error) {
	client, err := h.getClient(cluster)
	if err != nil {
		return nil, err
	}
	res, err := client.InstallRelease(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.Release, nil
}

func (h *Helm) UpdateRelease(ctx context.Context, cluster string, req *services.UpdateReleaseRequest) (*release.Release, error) {
	client, err := h.getClient(cluster)
	if err != nil {
		return nil, err
	}

	res, err := client.UpdateRelease(req)
	if err != nil {
		return nil, err
	}

	return res.Release, nil
}

func (h *Helm) RollbackRelease(ctx context.Context, cluster string, req *services.RollbackReleaseRequest) (*release.Release, error) {
	client, err := h.getClient(cluster)
	if err != nil {
		return nil, err
	}

	res, err := client.RollbackRelease(req)

	if err != nil {
		return nil, err
	}

	return res.Release, nil
}

func (h *Helm) ReleaseHistory(ctx context.Context, req *services.GetHistoryRequest) ([]*release.Release, error) {
	client, err := h.getClient(cluster)
	if err != nil {
		return nil, err
	}

	res, err := client.GetHistory(req)
	if err != nil {
		return nil, err
	}
	1
	releases := []*release.Release{}

	for _, r := range res.Releases {
		releases = append(releases, r)
	}
	return releases, nil
}

func provider(config map[string]interface{}) (alecton.DeployProvider, error) {
	var c helmConfig
	if err := alecton.ProviderConfigDecode(config, &c); err != nil {
		return errors.Wrap(err, "failed to deconde helm config")
	}
	return New(&config)
}

func init() {
	alecton.RegisterDeployProvider("helm", provider)
}
