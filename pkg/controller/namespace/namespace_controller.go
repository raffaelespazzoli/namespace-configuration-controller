package namespace

import (
	"container/list"
	"context"

	log "github.com/sirupsen/logrus"
	namespaceconfigv1alpha1 "github.com/raffaelespazzoli/namespace-configuration-controller/pkg/apis/namespaceconfig/v1alpha1"
	chandler "github.com/raffaelespazzoli/namespace-configuration-controller/pkg/controller/handler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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

// Add creates a new Namespace Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNamespace{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("namespace-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Namespace
	err = c.Watch(&source.Kind{Type: &corev1.Namespace{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Namespace
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &corev1.Namespace{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileNamespace{}

// ReconcileNamespace reconciles a Namespace object
type ReconcileNamespace struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Namespace object and makes changes based on the state read
// and what is in the Namespace.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNamespace) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Infof("Reconciling Namespace %s/%s\n", request.Namespace, request.Name)

	// Fetch the Namespace instance
	instance := &corev1.Namespace{}
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

	//find all namepaceconfig which apply to this namespace
	namespace := *instance
	ncl, err := r.findApplicableNameSpaceConfigs(*instance)

	errs := list.New()
	//apply config to namespace
	for _, nc := range ncl.Items {
		applyerrs := chandler.ApplyConfigToNamespace(nc, *instance, r.scheme, r.client)
		errs.PushFrontList(&applyerrs)
	}
	if errs.Len() > 0 {
		log.Infof("reconciliation on namespace %s executed with the following errors (if any) %+v\n", namespace.Name, errs)
	} else {
		log.Infof("reconciliation on namespace %s executed without errors\n", namespace.Name)
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileNamespace) findApplicableNameSpaceConfigs(namespace corev1.Namespace) (namespaceconfigv1alpha1.NamespaceConfigList, error) {
	log.Debugf("finding namespaceconfigs applicable to namespace %s", namespace.Name)
	//find all the namespacepolicies
	nclist := list.New()
	ncl := namespaceconfigv1alpha1.NamespaceConfigList{}
	opts := &client.ListOptions{}
	err := r.client.List(context.TODO(), opts, &ncl)
	if err != nil {
		return namespaceconfigv1alpha1.NamespaceConfigList{}, err
	}
	//for each netwrokpolicies get the selected namespaces
	for _, nc := range ncl.Items {
		selector, err := metav1.LabelSelectorAsSelector(&nc.Spec.LabelSelector)
		if err != nil {
			return namespaceconfigv1alpha1.NamespaceConfigList{}, err
		}
		if selector.Matches(labels.Set(namespace.Labels)) {
			nclist.PushFront(nc)
		}
	}
	ret := namespaceconfigv1alpha1.NamespaceConfigList{Items: make([]namespaceconfigv1alpha1.NamespaceConfig, nclist.Len())}
	e := nclist.Front()
	for i := range ret.Items {
		ret.Items[i] = e.Value.(namespaceconfigv1alpha1.NamespaceConfig)
		e = e.Next()
	}
	log.Debugf("found namespaceconfigs %+v applicable to namespace %s", ret.Items, namespace.Name)
	return ret, nil
}
