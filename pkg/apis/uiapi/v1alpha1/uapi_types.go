package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Ui struct {
	Size            int    `json:"size"`
	Name            string `json:"name"`
	ServiceNodePort int32  `json:"serviceNodePort"`
	ApiUrl          string `json:"apiUrl"`
	Image           string `json:"image"`
}

type Api struct {
	Size            int    `json:"size"`
	Name            string `json:"name"`
	ServiceNodePort int32  `json:"serviceNodePort"`
	ConfSecretName  string `json:"confSecretName"`
	Image           string `json:"image"`
}

type Db struct {
	Image string `json:"image"`
	Host  string `json:"host"`
	Port  int32  `json:"port"`
	Name  string `json:"name"`
}

// UapiSpec defines the desired state of Uapi
// +k8s:openapi-gen=true
type UapiSpec struct {
	Namespace string `json:"namespace"`
	Ui        Ui     `json:"ui"`
	Api       Api    `json:"api"`
	Db        Db     `json:"db"`
}

// UapiStatus defines the observed state of Uapi
// +k8s:openapi-gen=true
type UapiStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	OpStatus []string `json:"opStatus"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Uapi is the Schema for the uapis API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Uapi struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UapiSpec   `json:"spec,omitempty"`
	Status UapiStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UapiList contains a list of Uapi
type UapiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Uapi `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Uapi{}, &UapiList{})
}
