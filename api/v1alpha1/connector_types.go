package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// SecretKeySelector identifies a key inside a Kubernetes Secret.
type SecretKeySelector struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// ConnectorSpec defines the desired state of a Connector.
type ConnectorSpec struct {
	Type string `json:"type"`

	Config map[string]string `json:"config,omitempty"`

	SecretRef *SecretKeySelector `json:"secretRef,omitempty"`
}

// ConnectorStatus defines the observed state of a Connector.
type ConnectorStatus struct {
	Phase           string       `json:"phase,omitempty"`
	ResourcesSynced int32        `json:"resourcesSynced,omitempty"`
	LastSyncTime    *metav1.Time `json:"lastSyncTime,omitempty"`
	Message         string       `json:"message,omitempty"`
}

// Connector represents a KubeBridge Connector resource.
type Connector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConnectorSpec   `json:"spec,omitempty"`
	Status ConnectorStatus `json:"status,omitempty"`
}

// ConnectorList contains a list of Connector resources.
type ConnectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Connector `json:"items"`
}

func (in *Connector) DeepCopyInto(out *Connector) {
	*out = *in

	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = *in.ObjectMeta.DeepCopy()

	out.Spec.Type = in.Spec.Type

	if in.Spec.Config != nil {
		out.Spec.Config = make(map[string]string, len(in.Spec.Config))

		for key, value := range in.Spec.Config {
			out.Spec.Config[key] = value
		}
	}

	if in.Spec.SecretRef != nil {
		out.Spec.SecretRef = &SecretKeySelector{
			Name: in.Spec.SecretRef.Name,
			Key:  in.Spec.SecretRef.Key,
		}
	}

	out.Status = in.Status

	if in.Status.LastSyncTime != nil {
		out.Status.LastSyncTime = in.Status.LastSyncTime.DeepCopy()
	}
}

func (in *Connector) DeepCopy() *Connector {
	if in == nil {
		return nil
	}

	out := new(Connector)
	in.DeepCopyInto(out)

	return out
}

func (in *Connector) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}

func (in *ConnectorList) DeepCopyInto(out *ConnectorList) {
	*out = *in

	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]Connector, len(in.Items))

		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}

func (in *ConnectorList) DeepCopy() *ConnectorList {
	if in == nil {
		return nil
	}

	out := new(ConnectorList)
	in.DeepCopyInto(out)

	return out
}

func (in *ConnectorList) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}
