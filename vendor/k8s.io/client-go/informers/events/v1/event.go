/*
Copyright The Kubernetes Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	eventsv1 "k8s.io/api/events/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
	kubernetes "k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/events/v1"
	cache "k8s.io/client-go/tools/cache"
)

// EventInformer provides access to a shared informer and lister for
// Events.
type EventInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.EventLister
}

type eventInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewEventInformer constructs a new informer for Event type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewEventInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredEventInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredEventInformer constructs a new informer for Event type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredEventInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return NewFilteredEventInformerWithOptions(client, namespace, tweakListOptions, cache.WithResyncPeriod(resyncPeriod), cache.WithIndexers(indexers))
}

func NewFilteredEventInformerWithOptions(client kubernetes.Interface, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc, opts ...cache.SharedInformerOption) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformerWithOptions(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.EventsV1().Events(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.EventsV1().Events(namespace).Watch(context.TODO(), options)
			},
		},
		&eventsv1.Event{},
		opts...,
	)
}

func (f *eventInformer) defaultInformer(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	indexers := cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}
	for k, v := range f.factory.ExtraNamespaceScopedIndexers() {
		indexers[k] = v
	}

	return NewFilteredEventInformerWithOptions(client, f.namespace,
		f.tweakListOptions,
		cache.WithResyncPeriod(resyncPeriod),
		cache.WithIndexers(indexers),
		cache.WithKeyFunction(f.factory.KeyFunction()),
	)
}

func (f *eventInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&eventsv1.Event{}, f.defaultInformer)
}

func (f *eventInformer) Lister() v1.EventLister {
	return v1.NewEventLister(f.Informer().GetIndexer())
}
