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

	"github.com/gerrit91/gardener-extension-registry-cache/pkg/apis/config"
	"github.com/gerrit91/gardener-extension-registry-cache/pkg/apis/registry"
	"github.com/gerrit91/gardener-extension-registry-cache/pkg/apis/registry/v1alpha1"
	"github.com/gerrit91/gardener-extension-registry-cache/pkg/imagevector"

	"github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/extension"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/utils/managedresources"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ActuatorName is the name of the registry service actuator.
const ActuatorName = "registry-cache-actuator"

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

	if err := a.createResources(ctx, registryConfig, cluster, namespace); err != nil {
		return err
	}

	return nil
}

// Delete the Extension resource.
func (a *actuator) Delete(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.deleteResources(ctx, ex.GetNamespace())
}

// Restore the Extension resource.
func (a *actuator) Restore(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Reconcile(ctx, log, ex)
}

// Migrate the Extension resource.
func (a *actuator) Migrate(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return nil
}

// InjectConfig injects the rest config to this actuator.
func (a *actuator) InjectConfig(config *rest.Config) error {
	a.config = config
	return nil
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

func (a *actuator) createResources(ctx context.Context, registryConfig *service.RegistryConfig, _ *controller.Cluster, namespace string) error {
	registryImage, err := imagevector.ImageVector().FindImage("registry")
	if err != nil {
		return fmt.Errorf("failed to find registry image: %w", err)
	}

	for _, m := range registryConfig.Mirrors {
		c := registryCache{
			Client:                        a.client,
			Ctx:                           ctx,
			RemoteURL:                     m.RemoteURL,
			CacheVolumeSize:               m.CacheSize,
			CacheGarbageCollectionEnabled: m.CacheGarbageCollectionEnabled,
			RegistryImage:                 registryImage,
		}

		resources, err := c.EnsureRegistryCache()
		if err != nil {
			return err
		}

		// create manageresource from the registryCache
		err = a.createManagedResources(ctx, c.Name, namespace, "", resources, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *actuator) deleteResources(ctx context.Context, namespace string) error {
	a.logger.Info("Deleting managed resource for registry cache", "namespace", namespace)

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
