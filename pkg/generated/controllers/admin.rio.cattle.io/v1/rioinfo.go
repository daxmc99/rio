/*
Copyright 2020 Rancher Labs.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/rancher/rio/pkg/apis/admin.rio.cattle.io/v1"
	clientset "github.com/rancher/rio/pkg/generated/clientset/versioned/typed/admin.rio.cattle.io/v1"
	informers "github.com/rancher/rio/pkg/generated/informers/externalversions/admin.rio.cattle.io/v1"
	listers "github.com/rancher/rio/pkg/generated/listers/admin.rio.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type RioInfoHandler func(string, *v1.RioInfo) (*v1.RioInfo, error)

type RioInfoController interface {
	generic.ControllerMeta
	RioInfoClient

	OnChange(ctx context.Context, name string, sync RioInfoHandler)
	OnRemove(ctx context.Context, name string, sync RioInfoHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() RioInfoCache
}

type RioInfoClient interface {
	Create(*v1.RioInfo) (*v1.RioInfo, error)
	Update(*v1.RioInfo) (*v1.RioInfo, error)
	UpdateStatus(*v1.RioInfo) (*v1.RioInfo, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v1.RioInfo, error)
	List(opts metav1.ListOptions) (*v1.RioInfoList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.RioInfo, err error)
}

type RioInfoCache interface {
	Get(name string) (*v1.RioInfo, error)
	List(selector labels.Selector) ([]*v1.RioInfo, error)

	AddIndexer(indexName string, indexer RioInfoIndexer)
	GetByIndex(indexName, key string) ([]*v1.RioInfo, error)
}

type RioInfoIndexer func(obj *v1.RioInfo) ([]string, error)

type rioInfoController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.RioInfosGetter
	informer          informers.RioInfoInformer
	gvk               schema.GroupVersionKind
}

func NewRioInfoController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.RioInfosGetter, informer informers.RioInfoInformer) RioInfoController {
	return &rioInfoController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromRioInfoHandlerToHandler(sync RioInfoHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.RioInfo
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.RioInfo))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *rioInfoController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.RioInfo))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateRioInfoDeepCopyOnChange(client RioInfoClient, obj *v1.RioInfo, handler func(obj *v1.RioInfo) (*v1.RioInfo, error)) (*v1.RioInfo, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *rioInfoController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *rioInfoController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *rioInfoController) OnChange(ctx context.Context, name string, sync RioInfoHandler) {
	c.AddGenericHandler(ctx, name, FromRioInfoHandlerToHandler(sync))
}

func (c *rioInfoController) OnRemove(ctx context.Context, name string, sync RioInfoHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromRioInfoHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *rioInfoController) Enqueue(name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), "", name)
}

func (c *rioInfoController) EnqueueAfter(name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), "", name, duration)
}

func (c *rioInfoController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *rioInfoController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *rioInfoController) Cache() RioInfoCache {
	return &rioInfoCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *rioInfoController) Create(obj *v1.RioInfo) (*v1.RioInfo, error) {
	return c.clientGetter.RioInfos().Create(obj)
}

func (c *rioInfoController) Update(obj *v1.RioInfo) (*v1.RioInfo, error) {
	return c.clientGetter.RioInfos().Update(obj)
}

func (c *rioInfoController) UpdateStatus(obj *v1.RioInfo) (*v1.RioInfo, error) {
	return c.clientGetter.RioInfos().UpdateStatus(obj)
}

func (c *rioInfoController) Delete(name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.RioInfos().Delete(name, options)
}

func (c *rioInfoController) Get(name string, options metav1.GetOptions) (*v1.RioInfo, error) {
	return c.clientGetter.RioInfos().Get(name, options)
}

func (c *rioInfoController) List(opts metav1.ListOptions) (*v1.RioInfoList, error) {
	return c.clientGetter.RioInfos().List(opts)
}

func (c *rioInfoController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.RioInfos().Watch(opts)
}

func (c *rioInfoController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.RioInfo, err error) {
	return c.clientGetter.RioInfos().Patch(name, pt, data, subresources...)
}

type rioInfoCache struct {
	lister  listers.RioInfoLister
	indexer cache.Indexer
}

func (c *rioInfoCache) Get(name string) (*v1.RioInfo, error) {
	return c.lister.Get(name)
}

func (c *rioInfoCache) List(selector labels.Selector) ([]*v1.RioInfo, error) {
	return c.lister.List(selector)
}

func (c *rioInfoCache) AddIndexer(indexName string, indexer RioInfoIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.RioInfo))
		},
	}))
}

func (c *rioInfoCache) GetByIndex(indexName, key string) (result []*v1.RioInfo, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1.RioInfo))
	}
	return result, nil
}

type RioInfoStatusHandler func(obj *v1.RioInfo, status v1.RioInfoStatus) (v1.RioInfoStatus, error)

type RioInfoGeneratingHandler func(obj *v1.RioInfo, status v1.RioInfoStatus) ([]runtime.Object, v1.RioInfoStatus, error)

func RegisterRioInfoStatusHandler(ctx context.Context, controller RioInfoController, condition condition.Cond, name string, handler RioInfoStatusHandler) {
	statusHandler := &rioInfoStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromRioInfoHandlerToHandler(statusHandler.sync))
}

func RegisterRioInfoGeneratingHandler(ctx context.Context, controller RioInfoController, apply apply.Apply,
	condition condition.Cond, name string, handler RioInfoGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &rioInfoGeneratingHandler{
		RioInfoGeneratingHandler: handler,
		apply:                    apply,
		name:                     name,
		gvk:                      controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	RegisterRioInfoStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type rioInfoStatusHandler struct {
	client    RioInfoClient
	condition condition.Cond
	handler   RioInfoStatusHandler
}

func (a *rioInfoStatusHandler) sync(key string, obj *v1.RioInfo) (*v1.RioInfo, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	obj.Status = newStatus
	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(obj, "", nil)
		} else {
			a.condition.SetError(obj, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, obj.Status) {
		var newErr error
		obj, newErr = a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
	}
	return obj, err
}

type rioInfoGeneratingHandler struct {
	RioInfoGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *rioInfoGeneratingHandler) Handle(obj *v1.RioInfo, status v1.RioInfoStatus) (v1.RioInfoStatus, error) {
	objs, newStatus, err := a.RioInfoGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	apply := a.apply

	if !a.opts.DynamicLookup {
		apply = apply.WithStrictCaching()
	}

	if !a.opts.AllowCrossNamespace && !a.opts.AllowClusterScoped {
		apply = apply.WithSetOwnerReference(true, false).
			WithDefaultNamespace(obj.GetNamespace()).
			WithListerNamespace(obj.GetNamespace())
	}

	if !a.opts.AllowClusterScoped {
		apply = apply.WithRestrictClusterScoped()
	}

	return newStatus, apply.
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
