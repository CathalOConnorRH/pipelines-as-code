package adapter

import (
	"context"
	"github.com/kcp-dev/logicalcluster/v2"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const tlsMountPath = "/etc/pipelines-as-code/tls"

// isTLSEnabled validates if tls secret exist and if the required fields are defined
func (l listener) isTLSEnabled() (bool, string, string) {
	tlsSecret := os.Getenv("TLS_SECRET_NAME")
	tlsKey := os.Getenv("TLS_KEY")
	tlsCert := os.Getenv("TLS_CERT")

	tls, err := l.run.Clients.Kube.Cluster(logicalcluster.Name{}).CoreV1().Secrets(os.Getenv("SYSTEM_NAMESPACE")).
		Get(context.Background(), tlsSecret, v1.GetOptions{})
	if err != nil {
		return false, "", ""
	}
	_, ok := tls.Data[tlsKey]
	if !ok {
		return false, "", ""
	}
	_, ok = tls.Data[tlsCert]
	if !ok {
		return false, "", ""
	}

	return true,
		filepath.Join(tlsMountPath, tlsCert),
		filepath.Join(tlsMountPath, tlsKey)
}
