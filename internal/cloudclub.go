package cloudclub

import (
	"github.com/cloud-club/cloudclub-operator/internal/driver"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Manager struct {
	ApplicationClient *driver.ApplicationClient
}

func NewManager(kube client.Client) (*Manager, error) {
	applicationClient, err := driver.NewApplicationClient(kube)
	if err != nil {
		return nil, err
	}
	return &Manager{
		ApplicationClient: applicationClient,
	}, nil
}
