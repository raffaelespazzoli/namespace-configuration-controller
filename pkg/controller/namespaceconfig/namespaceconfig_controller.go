package namespaceconfig

import (
	"container/list"
	"context"

	log "github.com/sirupsen/logrus"

	namespaceconfigv1alpha1 "github.com/raffaelespazzoli/namespace-configuration-controller/pkg/apis/namespaceconfig/v1alpha1"
	chandler "github.com/raffaelespazzoli/namespace-configuration-controller/pkg/controller/handler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new NamespaceConfig Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNamespaceConfig{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("namespaceconfig-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource NamespaceConfig
	err = c.Watch(&source.Kind{Type: &namespaceconfigv1alpha1.NamespaceConfig{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner NamespaceConfig
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &namespaceconfigv1alpha1.NamespaceConfig{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileNamespaceConfig{}

// ReconcileNamespaceConfig reconciles a NamespaceConfig object
type ReconcileNamespaceConfig struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a NamespaceConfig object and makes changes based on the state read
// and what is in the NamespaceConfig.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNamespaceConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Infof("Reconciling NamespaceConfig %s/%s\n", request.Namespace, request.Name)
	// Fetch the NamespaceConfig instance
	instance := &namespaceconfigv1alpha1.NamespaceConfig{}
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

	//if namespaceconfig still exists

	nl := corev1.NamespaceList{}
	log.Debugf("retrieving list of namespaces selected by this namespaceconfig: %s", request.NamespacedName.String())
	selector, err := metav1.LabelSelectorAsSelector(&instance.Spec.LabelSelector)
	if err != nil {
		return reconcile.Result{}, err
	}
	opts := &client.ListOptions{LabelSelector: selector}
	err = r.client.List(context.TODO(), opts, &nl)
	if err != nil {
		return reconcile.Result{}, err
	}
	log.Debugf("found the following namespaces: %s", nl.Items)
	errs := list.New()

	for _, namespace := range nl.Items {
		applyerrs := chandler.ApplyConfigToNamespace(*instance, namespace, r.scheme, r.client)
		errs.PushFrontList(&applyerrs)
	}

	return reconcile.Result{}, nil
}
