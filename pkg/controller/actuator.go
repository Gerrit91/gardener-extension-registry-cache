// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	corev1 "k8s.io/api/core/v1"

	"github.com/gerrit91/gardener-extension-registry-cache/pkg/apis/config"
	"github.com/gerrit91/gardener-extension-registry-cache/pkg/apis/registry"
	"github.com/gerrit91/gardener-extension-registry-cache/pkg/apis/registry/v1alpha1"
	"github.com/gerrit91/gardener-extension-registry-cache/pkg/imagevector"

	extensionsconfig "github.com/gardener/gardener/extensions/pkg/apis/config"
	"github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/extension"
	"github.com/gardener/gardener/extensions/pkg/util"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"

	"github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/gardener/gardener/pkg/utils/managedresources"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewActuator returns an actuator responsible for Extension resources.
func NewActuator(config config.Configuration) extension.Actuator {
	return &actuator{
		config: config,
	}
}

type actuator struct {
	client  client.Client
	decoder runtime.Decoder
	config  config.Configuration
}

// InjectClient injects the controller runtime client into the reconciler.
func (a *actuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}

// InjectScheme injects the given scheme into the reconciler.
func (a *actuator) InjectScheme(scheme *runtime.Scheme) error {
	a.decoder = serializer.NewCodecFactory(scheme, serializer.EnableStrict).UniversalDecoder()
	return nil
}

// Reconcile the Extension resource.
func (a *actuator) Reconcile(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()

	cluster, err := controller.GetCluster(ctx, a.client, namespace)
	if err != nil {
		return err
	}

	registryConfig := &registry.RegistryConfig{}
	if ex.Spec.ProviderConfig != nil {
		if _, _, err := a.decoder.Decode(ex.Spec.ProviderConfig.Raw, nil, registryConfig); err != nil {
			return fmt.Errorf("failed to decode provider config: %w", err)
		}
	}

	if err := a.createResources(ctx, log, registryConfig, cluster, namespace); err != nil {
		return err
	}

	return nil
}

// Delete the Extension resource.
func (a *actuator) Delete(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.deleteResources(ctx, log, ex.GetNamespace())
}

// Restore the Extension resource.
func (a *actuator) Restore(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Reconcile(ctx, log, ex)
}

// Migrate the Extension resource.
func (a *actuator) Migrate(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return nil
}

func (a *actuator) createResources(ctx context.Context, log logr.Logger, registryConfig *registry.RegistryConfig, cluster *controller.Cluster, namespace string) error {
	registryImage, err := imagevector.ImageVector().FindImage("registry")
	if err != nil {
		return fmt.Errorf("failed to find registry image: %w", err)
	}
	ensurerImage, err := imagevector.ImageVector().FindImage("cri-config-ensurer")
	if err != nil {
		return fmt.Errorf("failed to find ensurer image: %w", err)
	}

	objects := []client.Object{
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: registryCacheNamespaceName,
			},
		},
	}

	for _, cache := range registryConfig.Caches {
		c := registryCache{
			Namespace:                registryCacheNamespaceName,
			Upstream:                 cache.Upstream,
			VolumeSize:               *cache.Size,
			GarbageCollectionEnabled: *cache.GarbageCollectionEnabled,
			RegistryImage:            registryImage,
		}

		os, err := c.Ensure()
		if err != nil {
			return err
		}

		objects = append(objects, os...)
	}

	resources, err := managedresources.NewRegistry(kubernetes.SeedScheme, kubernetes.SeedCodec, kubernetes.SeedSerializer).AddAllAndSerialize(objects...)
	if err != nil {
		return err
	}

	// create ManagedResource for the registryCache
	err = a.createManagedResources(ctx, v1alpha1.RegistryResourceName, namespace, "", resources, nil)
	if err != nil {
		return err
	}

	_, shootClient, err := util.NewClientForShoot(ctx, a.client, cluster.ObjectMeta.Name, client.Options{}, extensionsconfig.RESTOptions{})
	if err != nil {
		return fmt.Errorf("shoot client cannot be crated: %w", err)
	}

	var criMirrors map[string]string
	selector := labels.NewSelector()
	r, err := labels.NewRequirement(registryCacheServiceUpstreamLabel, selection.Exists, nil)
	if err != nil {
		return err
	}
	selector = selector.Add(*r)

	err = retry.Do(func() error {
		services := &corev1.ServiceList{}
		if err := shootClient.List(ctx, services, &client.ListOptions{
			Namespace:     registryCacheNamespaceName,
			LabelSelector: selector,
		}); err != nil {
			log.Error(err, "could not read extension from shoot namespace")
			return err
		}

		if len(services.Items) != len(registryConfig.Caches) {
			log.Info("not all services for all configured caches exist")
			return err
		}

		criMirrors = map[string]string{}

		for i := range services.Items {
			svc := services.Items[i]
			criMirrors[svc.Labels[registryCacheServiceUpstreamLabel]] = fmt.Sprintf("http://%s:%d", svc.Spec.ClusterIP, svc.Spec.Ports[0].Port)
		}

		return nil
	}, retry.Context(ctx), retry.LastErrorOnly(true))
	if err != nil {
		return err
	}

	e := criEnsurer{
		Name:            criEnsurerName,
		Namespace:       registryCacheNamespaceName,
		CRIEnsurerImage: ensurerImage,
		RegistryMirrors: criMirrors,
	}

	os := e.Ensure()
	objects = []client.Object{}
	objects = append(objects, os...)

	resources, err = managedresources.NewRegistry(kubernetes.SeedScheme, kubernetes.SeedCodec, kubernetes.SeedSerializer).AddAllAndSerialize(objects...)
	if err != nil {
		return err
	}

	err = a.createManagedResources(ctx, v1alpha1.RegistryEnsurerResourceName, namespace, "", resources, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *actuator) deleteResources(ctx context.Context, log logr.Logger, namespace string) error {
	log.Info("deleting managed resource for registry cache")

	if err := managedresources.Delete(ctx, a.client, namespace, v1alpha1.RegistryResourceName, false); err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	return managedresources.WaitUntilDeleted(timeoutCtx, a.client, namespace, v1alpha1.RegistryResourceName)
}

func (a *actuator) createManagedResources(ctx context.Context, name, namespace, class string, resources map[string][]byte, injectedLabels map[string]string) error {
	keepObjects := false
	forceOverwriteAnnotations := false
	secretsWithPrefix := false

	return managedresources.Create(ctx, a.client, namespace, name, secretsWithPrefix, class, resources, &keepObjects, injectedLabels, &forceOverwriteAnnotations)
}

func (a *actuator) updateStatus(ctx context.Context, ex *extensionsv1alpha1.Extension, _ *registry.RegistryConfig) error {
	patch := client.MergeFrom(ex.DeepCopy())
	// ex.Status.Resources = resources
	return a.client.Status().Patch(ctx, ex, patch)
}
