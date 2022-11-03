package fake

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type KeptnWorkloadInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            Status `json:"status,omitempty"`
}

type Status struct {
	PreDeploymentEvaluationStatus string `json:"preDeploymentEvaluationStatus,omitempty"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeptnWorkloadInstance) DeepCopyInto(out *KeptnWorkloadInstance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Status) DeepCopyInto(out *Status) {
	*out = *in
	if in.PreDeploymentEvaluationStatus != "" {
		out.PreDeploymentEvaluationStatus = in.PreDeploymentEvaluationStatus
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeptnWorkloadInstance.
func (in *KeptnWorkloadInstance) DeepCopy() *KeptnWorkloadInstance {
	if in == nil {
		return nil
	}
	out := new(KeptnWorkloadInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KeptnWorkloadInstance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
