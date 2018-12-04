package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	settingsv1alpha1 "k8s.io/api/settings/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NamespaceConfigSpec defines the desired state of NamespaceConfig
type NamespaceConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	LabelSelector   metav1.LabelSelector         `json:"selector,omitempty"`
	NetworkPolicies []networkingv1.NetworkPolicy `json:"networkpolicies,omitempty"`
	//Secrets             []corev1.Secret              `json:"secrets,omitempty"`
	Configmaps          []corev1.ConfigMap           `json:"configmaps,omitempty"`
	PodPresets          []settingsv1alpha1.PodPreset `json:"podpresets,omitempty"`
	Quotas              []corev1.ResourceQuota       `json:"quotas,omitempty"`
	LimitRanges         []corev1.LimitRange          `json:"limitranges,omitempty"`
	RoleBindings        []rbacv1.RoleBinding         `json:"rolebingings,omitempty"`
	ClusterRoleBindings []rbacv1.ClusterRoleBinding  `json:"clusterrolebindings,omitempty"`
	ServiceAccounts     []corev1.ServiceAccount      `json:"serviceaccounts,omitempty"`
}

// NamespaceConfigStatus defines the observed state of NamespaceConfig
type NamespaceConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceConfig is the Schema for the namespaceconfigs API
// +k8s:openapi-gen=true
type NamespaceConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamespaceConfigSpec   `json:"spec,omitempty"`
	Status NamespaceConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceConfigList contains a list of NamespaceConfig
type NamespaceConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NamespaceConfig{}, &NamespaceConfigList{})
}
