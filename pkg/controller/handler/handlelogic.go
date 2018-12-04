package handler

import (
	"container/list"
	"context"

	namespaceconfigv1alpha1 "github.com/raffaelespazzoli/namespace-configuration-controller/pkg/apis/namespaceconfig/v1alpha1"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	settingsv1alpha1 "k8s.io/api/settings/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func ApplyConfigToNamespace(namespaceconfig namespaceconfigv1alpha1.NamespaceConfig, namespace corev1.Namespace, scheme *runtime.Scheme, client client.Client) list.List {
	log.Debugf("starting to apply namespace config %s to namespace %s", namespaceconfig.Name, namespace.Name)
	errs := list.New()
	nperrs := reconcileNetworkPolicy(namespaceconfig, namespace, scheme, client)
	errs.PushFrontList(&nperrs)
	nperrs = reconcileConfigmaps(namespaceconfig, namespace, scheme, client)
	errs.PushFrontList(&nperrs)
	nperrs = reconcileLimitRanges(namespaceconfig, namespace, scheme, client)
	errs.PushFrontList(&nperrs)
	nperrs = reconcileClusterRoleBindings(namespaceconfig, namespace, scheme, client)
	errs.PushFrontList(&nperrs)
	nperrs = reconcilePodPreset(namespaceconfig, namespace, scheme, client)
	errs.PushFrontList(&nperrs)
	nperrs = reconcileQuotas(namespaceconfig, namespace, scheme, client)
	errs.PushFrontList(&nperrs)
	nperrs = reconcileRoleBindings(namespaceconfig, namespace, scheme, client)
	errs.PushFrontList(&nperrs)
	nperrs = reconcileServiceAccount(namespaceconfig, namespace, scheme, client)
	errs.PushFrontList(&nperrs)
	return *errs
}

func reconcileNetworkPolicy(namespaceconfig namespaceconfigv1alpha1.NamespaceConfig, namespace corev1.Namespace, scheme *runtime.Scheme, client client.Client) list.List {
	log.Debugf("reconciling configured network policies in namespace %s", namespace.Name)
	errs := list.New()
	log.Debugf("there are %d networkpolicy to apply", len(namespaceconfig.Spec.NetworkPolicies))
	for _, np := range namespaceconfig.Spec.NetworkPolicies {
		log.Debugf("reconciling network policy %s in namespace %s", np.Name, namespace.Name)
		cnp := np.DeepCopy()
		err := controllerutil.SetControllerReference(&namespaceconfig, cnp, scheme)
		if err != nil {
			log.Printf("Error Setting The Object Owner: %s", err)
			errs.PushFront(err)
			continue
		}
		cnp.Namespace = namespace.Name
		found := &networkingv1.NetworkPolicy{}
		err = client.Get(context.TODO(), types.NamespacedName{Name: np.Name, Namespace: namespace.Name}, found)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Debugf("network policy %s not found in namespace %s , creating it", np.Name, namespace.Name)
				err = client.Create(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error creating the NetworkPolicy Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			} else {
				log.Debugf("network policy %s found in namespace %s , updating it", np.Name, namespace.Name)
				err = client.Update(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error updating the NetworkPolicy Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			}
		}
	}
	return *errs
}

func reconcileConfigmaps(namespaceconfig namespaceconfigv1alpha1.NamespaceConfig, namespace corev1.Namespace, scheme *runtime.Scheme, client client.Client) list.List {
	log.Debugf("reconciling configured configmaps in namespace %s", namespace.Name)
	errs := list.New()
	log.Debugf("there are %d confimaps to apply", len(namespaceconfig.Spec.Configmaps))
	for _, np := range namespaceconfig.Spec.Configmaps {
		log.Debugf("reconciling configmaps %s in namespace %s", np.Name, namespace.Name)
		cnp := np.DeepCopy()
		err := controllerutil.SetControllerReference(&namespaceconfig, cnp, scheme)
		if err != nil {
			log.Printf("Error Setting The Object Owner: %s", err)
			errs.PushFront(err)
			continue
		}
		cnp.Namespace = namespace.Name
		found := &corev1.ConfigMap{}
		err = client.Get(context.TODO(), types.NamespacedName{Name: np.Name, Namespace: namespace.Name}, found)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Debugf("configmap %s not found in namespace %s , creating it", np.Name, namespace.Name)
				err = client.Create(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error creating the Configmap Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			} else {
				log.Debugf("configmap %s found in namespace %s , updating it", np.Name, namespace.Name)
				err = client.Update(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error updating the Configmap Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			}
		}
	}
	return *errs
}

func reconcilePodPreset(namespaceconfig namespaceconfigv1alpha1.NamespaceConfig, namespace corev1.Namespace, scheme *runtime.Scheme, client client.Client) list.List {
	log.Debugf("reconciling configured podpresets in namespace %s", namespace.Name)
	errs := list.New()
	log.Debugf("there are %d podpresets to apply", len(namespaceconfig.Spec.PodPresets))
	for _, np := range namespaceconfig.Spec.PodPresets {
		log.Debugf("reconciling podpresets %s in namespace %s", np.Name, namespace.Name)
		cnp := np.DeepCopy()
		err := controllerutil.SetControllerReference(&namespaceconfig, cnp, scheme)
		if err != nil {
			log.Printf("Error Setting The Object Owner: %s", err)
			errs.PushFront(err)
			continue
		}
		cnp.Namespace = namespace.Name
		found := &settingsv1alpha1.PodPreset{}
		err = client.Get(context.TODO(), types.NamespacedName{Name: np.Name, Namespace: namespace.Name}, found)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Debugf("podpreset %s not found in namespace %s , creating it", np.Name, namespace.Name)
				err = client.Create(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error creating the podpreset Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			} else {
				log.Debugf("podpreset %s found in namespace %s , updating it", np.Name, namespace.Name)
				err = client.Update(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error updating the podpreset Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			}
		}
	}
	return *errs
}

func reconcileQuotas(namespaceconfig namespaceconfigv1alpha1.NamespaceConfig, namespace corev1.Namespace, scheme *runtime.Scheme, client client.Client) list.List {
	log.Debugf("reconciling configured quotas in namespace %s", namespace.Name)
	errs := list.New()
	log.Debugf("there are %d quotas to apply", len(namespaceconfig.Spec.Quotas))
	for _, np := range namespaceconfig.Spec.Quotas {
		log.Debugf("reconciling quotas %s in namespace %s", np.Name, namespace.Name)
		cnp := np.DeepCopy()
		err := controllerutil.SetControllerReference(&namespaceconfig, cnp, scheme)
		if err != nil {
			log.Printf("Error Setting The Object Owner: %s", err)
			errs.PushFront(err)
			continue
		}
		cnp.Namespace = namespace.Name
		found := &corev1.ResourceQuota{}
		err = client.Get(context.TODO(), types.NamespacedName{Name: np.Name, Namespace: namespace.Name}, found)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Debugf("quotas %s not found in namespace %s , creating it", np.Name, namespace.Name)
				err = client.Create(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error creating the quotas Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			} else {
				log.Debugf("quotas %s found in namespace %s , updating it", np.Name, namespace.Name)
				err = client.Update(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error updating the quotas Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			}
		}
	}
	return *errs
}

func reconcileLimitRanges(namespaceconfig namespaceconfigv1alpha1.NamespaceConfig, namespace corev1.Namespace, scheme *runtime.Scheme, client client.Client) list.List {
	log.Debugf("reconciling configured limitranges in namespace %s", namespace.Name)
	errs := list.New()
	log.Debugf("there are %d limitranges to apply", len(namespaceconfig.Spec.LimitRanges))
	for _, np := range namespaceconfig.Spec.LimitRanges {
		log.Debugf("reconciling limitranges %s in namespace %s", np.Name, namespace.Name)
		cnp := np.DeepCopy()
		err := controllerutil.SetControllerReference(&namespaceconfig, cnp, scheme)
		if err != nil {
			log.Printf("Error Setting The Object Owner: %s", err)
			errs.PushFront(err)
			continue
		}
		cnp.Namespace = namespace.Name
		found := &corev1.LimitRange{}
		err = client.Get(context.TODO(), types.NamespacedName{Name: np.Name, Namespace: namespace.Name}, found)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Debugf("limitranges %s not found in namespace %s , creating it", np.Name, namespace.Name)
				err = client.Create(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error creating the limitranges Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			} else {
				log.Debugf("limitranges %s found in namespace %s , updating it", np.Name, namespace.Name)
				err = client.Update(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error updating the limitranges Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			}
		}
	}
	return *errs
}

func reconcileRoleBindings(namespaceconfig namespaceconfigv1alpha1.NamespaceConfig, namespace corev1.Namespace, scheme *runtime.Scheme, client client.Client) list.List {
	log.Debugf("reconciling configured rolebindings in namespace %s", namespace.Name)
	errs := list.New()
	log.Debugf("there are %d rolebindings to apply", len(namespaceconfig.Spec.RoleBindings))
	for _, np := range namespaceconfig.Spec.RoleBindings {
		log.Debugf("reconciling rolebindings %s in namespace %s", np.Name, namespace.Name)
		cnp := np.DeepCopy()
		err := controllerutil.SetControllerReference(&namespaceconfig, cnp, scheme)
		if err != nil {
			log.Printf("Error Setting The Object Owner: %s", err)
			errs.PushFront(err)
			continue
		}
		cnp.Namespace = namespace.Name
		found := &rbacv1.RoleBinding{}
		err = client.Get(context.TODO(), types.NamespacedName{Name: np.Name, Namespace: namespace.Name}, found)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Debugf("rolebindings %s not found in namespace %s , creating it", np.Name, namespace.Name)
				err = client.Create(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error creating the rolebindings Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			} else {
				log.Debugf("rolebindings %s found in namespace %s , updating it", np.Name, namespace.Name)
				err = client.Update(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error updating the rolebindings Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			}
		}
	}
	return *errs
}

func reconcileClusterRoleBindings(namespaceconfig namespaceconfigv1alpha1.NamespaceConfig, namespace corev1.Namespace, scheme *runtime.Scheme, client client.Client) list.List {
	log.Debugf("reconciling configured clusterrolebindings in namespace %s", namespace.Name)
	errs := list.New()
	log.Debugf("there are %d clusterrolebindings to apply", len(namespaceconfig.Spec.ClusterRoleBindings))
	for _, np := range namespaceconfig.Spec.ClusterRoleBindings {
		log.Debugf("reconciling clusterrolebindings %s in namespace %s", np.Name, namespace.Name)
		cnp := np.DeepCopy()
		err := controllerutil.SetControllerReference(&namespaceconfig, cnp, scheme)
		if err != nil {
			log.Printf("Error Setting The Object Owner: %s", err)
			errs.PushFront(err)
			continue
		}
		cnp.Namespace = namespace.Name
		found := &rbacv1.ClusterRoleBinding{}
		err = client.Get(context.TODO(), types.NamespacedName{Name: np.Name, Namespace: namespace.Name}, found)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Debugf("clusterrolebindings %s not found in namespace %s , creating it", np.Name, namespace.Name)
				err = client.Create(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error creating the clusterrolebindings Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			} else {
				log.Debugf("clusterrolebindings %s found in namespace %s , updating it", np.Name, namespace.Name)
				err = client.Update(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error updating the clusterrolebindings Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			}
		}
	}
	return *errs
}

func reconcileServiceAccount(namespaceconfig namespaceconfigv1alpha1.NamespaceConfig, namespace corev1.Namespace, scheme *runtime.Scheme, client client.Client) list.List {
	log.Debugf("reconciling configured serviceaccounts in namespace %s", namespace.Name)
	errs := list.New()
	log.Debugf("there are %d serviceaccounts to apply", len(namespaceconfig.Spec.ServiceAccounts))
	for _, np := range namespaceconfig.Spec.ServiceAccounts {
		log.Debugf("reconciling serviceaccounts %s in namespace %s", np.Name, namespace.Name)
		cnp := np.DeepCopy()
		err := controllerutil.SetControllerReference(&namespaceconfig, cnp, scheme)
		if err != nil {
			log.Printf("Error Setting The Object Owner: %s", err)
			errs.PushFront(err)
			continue
		}
		cnp.Namespace = namespace.Name
		found := &corev1.ServiceAccount{}
		err = client.Get(context.TODO(), types.NamespacedName{Name: np.Name, Namespace: namespace.Name}, found)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Debugf("serviceaccounts %s not found in namespace %s , creating it", np.Name, namespace.Name)
				err = client.Create(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error creating the serviceaccounts Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			} else {
				log.Debugf("serviceaccounts %s found in namespace %s , updating it", np.Name, namespace.Name)
				err = client.Update(context.TODO(), cnp)
				if err != nil {
					log.Printf("Error updating the serviceaccounts Object: %s", err)
					errs.PushFront(err)
					continue
				}
				continue
			}
		}
	}
	return *errs
}
