package cloudclub

import (
	"github.com/cloud-club/cloudclub-operator/internal/driver"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Manager struct {
	ApplicationClient *driver.ApplicationClient
}

func NewManager(kube client.Client, schema *runtime.Scheme) (*Manager, error) {
	applicationClient, err := driver.NewApplicationClient(kube, schema)
	if err != nil {
		return nil, err
	}
	return &Manager{
		ApplicationClient: applicationClient,
	}, nil
}
