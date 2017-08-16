package k8s

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/bakins/alecton"
	"github.com/bakins/alecton/api"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
)

type Config struct {
	Kubeconfig     string
	Context        string
	SchemaCacheDir string
	Timeout        int
}

type Kubectl struct {
	Path     string
	Clusters map[string]Config
}

func New(raw map[string]interface{}) (*Kubectl, error) {
	var k Kubectl
	if err := alecton.ProviderConfigDecode(raw, &k); err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	if k.Path == "" {
		k.Path = "kubectl"
	}

	path, err := exec.LookPath(k.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find command \"%s\"in PATH ", k.Path)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find absolute path of %s", k.Path)
	}
	k.Path = absPath
	return nil, nil
}

func (k *Kubectl) buildKubectlCommand(cluster string, in ...string) (*exec.Cmd, error) {
	c, ok := k.Clusters[cluster]
	if !ok {
		return nil, errors.Errorf("unable to find config for cluster: \"%s\"", cluster)
	}
	args := []string{}

	if c.Context != "" {
		args = append(args, "--context")
		args = append(args, c.Context)
	}

	if c.Kubeconfig != "" {
		args = append(args, "--kubeconfig")
		args = append(args, c.Kubeconfig)
	}

	if c.SchemaCacheDir != "" {
		args = append(args, "--schema-cache-dir")
		args = append(args, c.SchemaCacheDir)
	}

	if c.Timeout <= 0 {
		c.Timeout = 5
	}

	args = append(args, "--timeout=")
	args = append(args, strconv.Itoa(c.Timeout))

	args = append(in)

	cmd := exec.Command(k.Path, args...)
	return cmd, nil
}

func (k *Kubectl) EnsureNamespace(ctx context.Context, cluster string, namespace string) error {
	cmd, err := k.buildKubectlCommand(cluster, "get", "namespace", namespace)
	if err != nil {
		return errors.Wrapf(err, "unable to build 'kubectl get namespace' command for %s/%s", cluster, namespace)
	}

	var buff bytes.Buffer
	cmd.Stderr = &buff
	cmd.Stdout = ioutil.Discard
	err := cmd.Run()
	switch {
	case err == nil:
		// namespace exists
		return nil
	case bytes.Contains(buff.Bytes(), []byte("NotFound")):
	// need to create, handled below
	default:
		return errors.Wrapf(err, "failed to get namespace %s/%s: %s", cluster, namespace, buff.String())
	}

	cmd, err = k.buildKubectlCommand(cluster, "create", "namespace", namespace)
	if err != nil {
		return errors.Wrapf(err, "unable to build 'kubectl create namespace' command for %s/%s", cluster, namespace)
	}

	buff.Reset()
	cmd.Stderr = buff
	cmd.Stdout = ioutil.Discard
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to create namespace %s/%s: %s", cluster, namespace, buff.String())
	}

	return nil
}

func (k *Kubectl) Deploy(ctx context.Context, cluster string, namespace string, resources []*api.Resource) error {
	// TODO: sort manifests like kind_sorter.go in tiller.
	// could sort in caller of this

	var buf bytes.Buffer
	buf.WriteString("---\n")
	for _, r := range resources {
		buf.WriteString(r.Data)
		buf.WriteString("\n---\n")
	}

	cmd, err := k.buildKubectlCommand(cluster, "apply", "--namespace", namespace, "--filename", "-")
	if err != nil {
		return errors.Wrapf(err, "unable to build 'kubectl apply' command for %s/%s", cluster, namespace)
	}

	var buff bytes.Buffer
	cmd.Stderr = &buff
	cmd.Stdout = ioutil.Discard

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to apply for namespace %s/%s: %s", cluster, namespace, buff.String())
	}

	return nil
}

func provider(c map[string]interface{}) (alecton.DeployProvider, error) {
	return New(c)
}

func init() {
	alecton.RegisterDeployProvider("k8s", provider)
}
