package stub

import (
	"reflect"

	ambassadorshimv1alpha1 "admiralty.io/ambassador-shim-operator-sdk/pkg/apis/ambassadorshim/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"github.com/operator-framework/operator-sdk/pkg/sdk/handler"
	"github.com/operator-framework/operator-sdk/pkg/sdk/query"
	"github.com/operator-framework/operator-sdk/pkg/sdk/types"
	yaml "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewHandler() handler.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx types.Context, event types.Event) error {
	switch o := event.Object.(type) {
	case *ambassadorshimv1alpha1.Mapping:
		m := o
		if event.Deleted {
			// The Mapping was deleted:
			// garbage collection will take care of the dummy Service.
			return nil
		}

		// generate the desired Service from the Mapping
		ds, err := dummyService(m)
		if err != nil {
			return err
		}
		addOwnerRefToObject(ds, asOwner(m))

		os := ds.DeepCopy()
		if err := query.Get(os); err != nil { // Get() uses TypeMeta and ObjectMeta from os, then overwrites it.
			if errors.IsNotFound(err) {
				// if the Service doesn't exist, create it
				// (update Mapping status for current observed state)
				m.Status = ambassadorshimv1alpha1.MappingStatus{
					Configured: false,
					UpToDate:   false,
				}
				if err := action.Update(m); err != nil {
					return err
				}

				err := action.Create(ds)
				return err
			}
			return err
		}

		// if the Service exist and its annotation matches the MappingSpec
		// do nothing but update the Mapping status
		if reflect.DeepEqual(ds.Annotations, os.Annotations) {
			m.Status = ambassadorshimv1alpha1.MappingStatus{
				Configured: true,
				UpToDate:   true,
			}
			err := action.Update(m)
			return err
		}

		// if the Service exists but its annotation doesn't match
		// update it accordingly
		m.Status = ambassadorshimv1alpha1.MappingStatus{
			Configured: true,
			UpToDate:   false,
		}
		if err := action.Update(m); err != nil {
			return err
		}

		os.Annotations = ds.Annotations
		err = action.Update(os)
		return err
	}
	return nil
}

// LegacyMapping is a representation of a Service annotation,
// which can be marshalled to YAML.
type LegacyMapping struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Name       string `yaml:"name"`
	Prefix     string `yaml:"prefix"`
	Service    string `yaml:"service"`
}

// dummyService creates a Service object from a Mapping object.
// The Service's annotations match the Mapping's spec.
func dummyService(m *ambassadorshimv1alpha1.Mapping) (*corev1.Service, error) {
	// Let's build the annotation as a struct,
	// before marshalling it to YAML.
	lm := LegacyMapping{
		APIVersion: "ambassador/v0",
		Kind:       "Mapping",
		Name:       m.Name,
		Prefix:     m.Spec.Prefix,
		Service:    m.Spec.Service,
	}

	y, err := yaml.Marshal(&lm)
	if err != nil {
		return nil, err
	}

	s := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-ambassadorshim",
			Namespace: m.Namespace,
			Annotations: map[string]string{
				"getambassador.io/config": string(y),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				corev1.ServicePort{Port: 80},
			}, // dummy port (required in ServiceSpec)
		},
	}

	return s, nil
}

// addOwnerRefToObject appends the desired OwnerReference to the object
func addOwnerRefToObject(obj metav1.Object, ownerRef metav1.OwnerReference) {
	obj.SetOwnerReferences(append(obj.GetOwnerReferences(), ownerRef))
}

// asOwner returns an OwnerReference set as the memcached CR
func asOwner(m *ambassadorshimv1alpha1.Mapping) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: m.APIVersion,
		Kind:       m.Kind,
		Name:       m.Name,
		UID:        m.UID,
		Controller: &trueVar,
	}
}
