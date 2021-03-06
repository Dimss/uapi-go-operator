package uapi

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	uiapiv1alpha1 "github.com/uapi-go-operator/pkg/apis/uiapi/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strconv"
)

var log = logf.Log.WithName("controller_uapi")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Uapi Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileUapi{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("uapi-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Uapi
	err = c.Watch(&source.Kind{Type: &uiapiv1alpha1.Uapi{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Uapi
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &uiapiv1alpha1.Uapi{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileUapi implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileUapi{}

// ReconcileUapi reconciles a Uapi object
type ReconcileUapi struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Uapi object and makes changes based on the state read
// and what is in the Uapi.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileUapi) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Uapi")

	// Fetch the Uapi instance
	instance := &uiapiv1alpha1.Uapi{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	mongoDeployment := newMongoDbDeployment(instance)
	mongoService := newMongoService(instance)
	apiSecret := newApiSecret(instance)
	apiDeployment := newApiDeployment(instance)
	apiService := newApiService(instance)
	uiService := newUiService(instance)
	uiDeployment := newUiDeployment(instance)

	// Set Uapi instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, mongoDeployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Set Uapi instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, mongoService, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Set Uapi instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, apiSecret, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Set Uapi instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, apiDeployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Set Uapi instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, apiService, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Set Uapi instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, uiService, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Set Uapi instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, uiDeployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: mongoDeployment.Name, Namespace: mongoDeployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new MongoDB", "Deployment.Namespace", mongoDeployment.Namespace, "Deployment.Name", mongoDeployment.Name)
		err = r.client.Create(context.TODO(), mongoDeployment)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), mongoService)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), apiSecret)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), apiDeployment)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), apiService)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), uiService)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), uiDeployment)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

func newPodForDB(cr *uiapiv1alpha1.Uapi) corev1.PodTemplateSpec {
	labels := map[string]string{
		"app": cr.Spec.Db.Host,
	}

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.Db.Host + "-pod",
			Namespace: cr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            cr.Spec.Db.Name,
					Image:           cr.Spec.Db.Image,
					ImagePullPolicy: corev1.PullAlways,
					Ports:           []corev1.ContainerPort{{ContainerPort: cr.Spec.Db.Port}},
					Env: []corev1.EnvVar{
						{Name: "MONGODB_USER", Value: cr.Spec.Db.Host},
						{Name: "MONGODB_PASSWORD", Value: cr.Spec.Db.Host},
						{Name: "MONGODB_ADMIN_PASSWORD", Value: cr.Spec.Db.Host},
						{Name: "MONGODB_DATABASE", Value: cr.Spec.Db.Name},
					},
				},
			},
		},
	}
}

func newPodForApi(cr *uiapiv1alpha1.Uapi) corev1.PodTemplateSpec {
	labels := map[string]string{
		"app": cr.Spec.Api.Name,
	}

	envs := []corev1.EnvVar{
		{
			Name: "PROFILE",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: cr.Spec.Api.ConfSecretName},
					Key:                  "profile",
				},
			},
		},
		{
			Name: "DB_HOST",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: cr.Spec.Api.ConfSecretName},
					Key:                  "db_host",
				},
			},
		},
		{
			Name: "DB_PORT",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: cr.Spec.Api.ConfSecretName},
					Key:                  "db_port",
				},
			},
		},
		{
			Name: "DB_NAME",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: cr.Spec.Api.ConfSecretName},
					Key:                  "db_name",
				},
			},
		},
		{
			Name: "DB_USER",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: cr.Spec.Api.ConfSecretName},
					Key:                  "db_host",
				},
			},
		},
		{
			Name: "DB_PASS",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: cr.Spec.Api.ConfSecretName},
					Key:                  "db_host",
				},
			},
		},
	}
	readinessProbe := &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/ready",
				Port: intstr.IntOrString{IntVal: int32(8080)},
			},
		},
		InitialDelaySeconds: 3,
		PeriodSeconds:       3,
	}

	livenessProbe := &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/healthy",
				Port: intstr.IntOrString{IntVal: int32(8080)},
			},
		},
		InitialDelaySeconds: 3,
		PeriodSeconds:       3,
	}

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.Db.Name + "-pod",
			Namespace: cr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            cr.Spec.Api.Name,
					Image:           cr.Spec.Api.Image,
					ImagePullPolicy: corev1.PullAlways,
					Ports:           []corev1.ContainerPort{{ContainerPort: int32(8080)}},
					Env:             envs,
					ReadinessProbe:  readinessProbe,
					LivenessProbe:   livenessProbe,
				},
			},
		},
	}
}

func newMongoService(cr *uiapiv1alpha1.Uapi) *corev1.Service {

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.Db.Host,
			Namespace: cr.Spec.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": cr.Spec.Db.Host},
			Ports:    []corev1.ServicePort{{Name: "mongo", Port: cr.Spec.Db.Port}},
		},
	}
	return &service
}

func newApiService(cr *uiapiv1alpha1.Uapi) *corev1.Service {

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.Api.Name,
			Namespace: cr.Spec.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": cr.Spec.Api.Name},
			Ports:    []corev1.ServicePort{{Name: "http", Port: 8080, NodePort: cr.Spec.Api.ServiceNodePort}},
			Type:     corev1.ServiceTypeNodePort,
		},
	}
	return &service
}

func newMongoDbDeployment(cr *uiapiv1alpha1.Uapi) (deployment *appsv1.Deployment) {
	var size int32
	size = 1
	labels := map[string]string{
		"app": cr.Spec.Db.Host,
	}
	metadata := metav1.ObjectMeta{
		Namespace: cr.Spec.Namespace,
		Name:      cr.Spec.Db.Host,
		Labels:    labels,
	}
	spec := appsv1.DeploymentSpec{
		Replicas: &size,
		Selector: &metav1.LabelSelector{MatchLabels: labels},
		Template: newPodForDB(cr),
	}
	deployment = &appsv1.Deployment{
		ObjectMeta: metadata,
		Spec:       spec,
	}
	return
}

func newApiSecret(cr *uiapiv1alpha1.Uapi) (secret *corev1.Secret) {

	secret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: cr.Spec.Api.ConfSecretName, Namespace: cr.Spec.Namespace},
		Type:       corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"profile": "prod",
			"db_host": cr.Spec.Db.Host,
			"db_port": strconv.Itoa(int(cr.Spec.Db.Port)),
			"db_name": cr.Spec.Db.Name,
		},
	}
	return
}

func newPodForUi(cr *uiapiv1alpha1.Uapi) corev1.PodTemplateSpec {
	labels := map[string]string{
		"app": cr.Spec.Ui.Name,
	}

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.Ui.Name + "-pod",
			Namespace: cr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            cr.Spec.Ui.Name,
					Image:           cr.Spec.Ui.Image,
					ImagePullPolicy: corev1.PullAlways,
					Ports:           []corev1.ContainerPort{{ContainerPort: int32(3000)}},
					Env: []corev1.EnvVar{
						{Name: "API_URL", Value: cr.Spec.Ui.ApiUrl},
					},
				},
			},
		},
	}
}

func newApiDeployment(cr *uiapiv1alpha1.Uapi) (deployment *appsv1.Deployment) {
	var size int32
	size = 1
	labels := map[string]string{
		"app": cr.Spec.Api.Name,
	}
	metadata := metav1.ObjectMeta{
		Namespace: cr.Spec.Namespace,
		Name:      cr.Spec.Api.Name,
		Labels:    labels,
	}
	spec := appsv1.DeploymentSpec{
		Replicas: &size,
		Selector: &metav1.LabelSelector{MatchLabels: labels},
		Template: newPodForApi(cr),
	}
	deployment = &appsv1.Deployment{
		ObjectMeta: metadata,
		Spec:       spec,
	}
	return
}

func newUiService(cr *uiapiv1alpha1.Uapi) *corev1.Service {

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.Ui.Name,
			Namespace: cr.Spec.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": cr.Spec.Ui.Name},
			Ports:    []corev1.ServicePort{{Name: "http", Port: 3000, NodePort: cr.Spec.Ui.ServiceNodePort}},
			Type:     corev1.ServiceTypeNodePort,
		},
	}
	return &service
}

func newUiDeployment(cr *uiapiv1alpha1.Uapi) (deployment *appsv1.Deployment) {
	var size int32
	size = 1
	labels := map[string]string{
		"app": cr.Spec.Ui.Name,
	}
	metadata := metav1.ObjectMeta{
		Namespace: cr.Spec.Namespace,
		Name:      cr.Spec.Ui.Name,
		Labels:    labels,
	}
	spec := appsv1.DeploymentSpec{
		Replicas: &size,
		Selector: &metav1.LabelSelector{MatchLabels: labels},
		Template: newPodForUi(cr),
	}
	deployment = &appsv1.Deployment{
		ObjectMeta: metadata,
		Spec:       spec,
	}
	return
}
