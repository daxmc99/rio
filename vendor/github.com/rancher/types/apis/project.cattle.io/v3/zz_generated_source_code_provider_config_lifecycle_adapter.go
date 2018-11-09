package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type SourceCodeProviderConfigLifecycle interface {
	Create(obj *SourceCodeProviderConfig) (runtime.Object, error)
	Remove(obj *SourceCodeProviderConfig) (runtime.Object, error)
	Updated(obj *SourceCodeProviderConfig) (runtime.Object, error)
}

type sourceCodeProviderConfigLifecycleAdapter struct {
	lifecycle SourceCodeProviderConfigLifecycle
}

func (w *sourceCodeProviderConfigLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*SourceCodeProviderConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *sourceCodeProviderConfigLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*SourceCodeProviderConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *sourceCodeProviderConfigLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*SourceCodeProviderConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewSourceCodeProviderConfigLifecycleAdapter(name string, clusterScoped bool, client SourceCodeProviderConfigInterface, l SourceCodeProviderConfigLifecycle) SourceCodeProviderConfigHandlerFunc {
	adapter := &sourceCodeProviderConfigLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *SourceCodeProviderConfig) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
