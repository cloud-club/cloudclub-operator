package driver

import (
	"context"

	appv1alpha1 "github.com/cloud-club/cloudclub-operator/api/v1alpha1"
	"github.com/cloud-club/cloudclub-operator/internal/log"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ApplicationClient struct {
	Kubernetes client.Client
	Schema     *runtime.Scheme
}

func NewApplicationClient(kube client.Client, schema *runtime.Scheme) (*ApplicationClient, error) {
	return &ApplicationClient{
		Kubernetes: kube,
		Schema:     schema,
	}, nil
}

func (a *ApplicationClient) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Info("start application reconcile")
	app := &appv1alpha1.Application{}
	err := a.Kubernetes.Get(ctx, req.NamespacedName, app)
	if err != nil {
		log.Errorf(err)
		if errors.IsNotFound(err) {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		return ctrl.Result{}, err
	}

	if err := a.UpsertDeployment(ctx, req, app); err != nil {
		log.Errorf(err)
		return ctrl.Result{}, err
	}

	if err := a.UpsertService(ctx, req, app); err != nil {
		log.Errorf(err)
		return ctrl.Result{}, err
	}
	log.Info("finish application reconcile")
	return ctrl.Result{}, nil
}

func (a *ApplicationClient) UpsertDeployment(ctx context.Context, req ctrl.Request, app *appv1alpha1.Application) error {
	deployment := &v1.Deployment{}
	err := a.Kubernetes.Get(ctx, req.NamespacedName, deployment)
	if err != nil {
		if errors.IsNotFound(err) {
			if deployment.DeletionTimestamp == nil {
				newDeployment := &v1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      req.Name,
						Namespace: req.Namespace,
						// OwnerReferences: []metav1.OwnerReference{
						//  *metav1.NewControllerRef(deployment, appv1alpha1.GroupVersion.WithKind("Application")),
						// },
					},
					Spec: v1.DeploymentSpec{
						Replicas: app.Spec.App.Replicas,
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": req.Name,
							},
						},
						Template: corev1.PodTemplateSpec{
							ObjectMeta: metav1.ObjectMeta{
								Labels: map[string]string{
									"app": req.Name,
								},
							},
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Name:  req.Name,
										Image: app.Spec.App.Image,
									},
								},
							},
						},
					},
				}
				// affinity 설정이 존재하면 새로운 deployment에 추가
				if app.Spec.Scheduler.Affinity != nil {
					newDeployment.Spec.Template.Spec.Affinity = app.Spec.Scheduler.Affinity
				} else {
					// Node Affinity
					newDeployment.Spec.Template.Spec.Affinity = &corev1.Affinity{
						NodeAffinity: &corev1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
								NodeSelectorTerms: []corev1.NodeSelectorTerm{
									{
										MatchExpressions: []corev1.NodeSelectorRequirement{
											{
												Key:      "key",
												Operator: corev1.NodeSelectorOpIn,
												Values:   []string{"value"},
											},
										},
									},
								},
							},
						},
					}
				}
				ctrl.SetControllerReference(app, newDeployment, a.Schema)
				return a.Kubernetes.Create(ctx, newDeployment)
			}
		}
		return err
	}

	if app.Spec.App.Replicas != deployment.Spec.Replicas {
		deployment.Spec.Replicas = app.Spec.App.Replicas
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
	appServicePort := intstr.IntOrString{IntVal: app.Spec.App.ContainerPort}
	if service.Spec.Ports[0].TargetPort != appServicePort {
		service.Spec.Ports[0].Port = app.Spec.App.ContainerPort
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
					Port: app.Spec.App.ContainerPort,
					TargetPort: intstr.IntOrString{
						IntVal: app.Spec.App.ContainerPort,
					},
				},
			},
			Selector: map[string]string{
				"app": app.Name,
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
	}
	return newService
}
