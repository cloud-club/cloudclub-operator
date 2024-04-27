package driver

import (
	"context"
	appv1alpha1 "github.com/cloud-club/cloudclub-operator/api/v1alpha1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ApplicationClient struct {
	Kubernetes client.Client
}

func NewApplicationClient(kube client.Client) (*ApplicationClient, error) {
	return &ApplicationClient{
		Kubernetes: kube,
	}, nil
}

func (a *ApplicationClient) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	app := &appv1alpha1.Application{}
	err := a.Kubernetes.Get(ctx, req.NamespacedName, app)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		return ctrl.Result{}, err
	}

	if err := a.UpsertDeployment(ctx, req, app); err != nil {
		return ctrl.Result{}, err
	}

	if err := a.UpsertService(ctx, req, app); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (a *ApplicationClient) UpsertDeployment(ctx context.Context, req ctrl.Request, app *appv1alpha1.Application) error {
	deployment := &v1.Deployment{}
	err := a.Kubernetes.Get(ctx, req.NamespacedName, deployment)
	if err != nil {
		if errors.IsNotFound(err) {
			newDeployment := &v1.Deployment{}
			return a.Kubernetes.Create(ctx, newDeployment)
		}
		return err
	}
	if app.Spec.Replicas != deployment.Spec.Replicas {
		deployment.Spec.Replicas = app.Spec.Replicas
		return a.Kubernetes.Update(ctx, deployment)
	}
	return nil
}

func (a *ApplicationClient) UpsertService(ctx context.Context, req ctrl.Request, app *appv1alpha1.Application) error {
	service := &corev1.Service{}
	err := a.Kubernetes.Get(ctx, req.NamespacedName, service)
	if err != nil {
		if errors.IsNotFound(err) {
			newService := a.createNewService(app)
			return a.Kubernetes.Create(ctx, newService)
		}
	}
	appServicePort := intstr.IntOrString{IntVal: app.Spec.ServicePort}
	if service.Spec.Ports[0].TargetPort != appServicePort {
		service.Spec.Ports[0].TargetPort = appServicePort
		return a.Kubernetes.Update(ctx, service)
	}
	return nil
}

func (a *ApplicationClient) createNewService(app *appv1alpha1.Application) *corev1.Service {
	newService := &corev1.Service{
		Spec: corev1.ServiceSpec{
			Type: "ClusterIP",
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: app.Spec.ServicePort,
					TargetPort: intstr.IntOrString{
						IntVal: app.Spec.ServicePort,
					},
				},
			},
			Selector: map[string]string{},
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
	}
	return newService
}