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
	"path/filepath"
	"time"

	"github.com/Gerrit91/gardener-extension-registry-cache/charts"
	"github.com/Gerrit91/gardener-extension-registry-cache/pkg/apis/config"
	"github.com/Gerrit91/gardener-extension-registry-cache/pkg/apis/service"
	"github.com/Gerrit91/gardener-extension-registry-cache/pkg/apis/service/v1alpha1"
	"github.com/Gerrit91/gardener-extension-registry-cache/pkg/apis/service/validation"

	"github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/extension"
	"github.com/gardener/gardener/extensions/pkg/util"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/chartrenderer"
	"github.com/gardener/gardener/pkg/utils/managedresources"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ActuatorName is the name of the registry service actuator.
const ActuatorName = "registry-cache-actuator"

// NewActuator returns an actuator responsible for Extension resources.
func NewActuator(config config.Configuration) extension.Actuator {
	return &actuator{
		logger:        log.Log.WithName(ActuatorName),
		serviceConfig: config,
	}
}

type actuator struct {
	client  client.Client
	config  *rest.Config
	decoder runtime.Decoder

	serviceConfig config.Configuration

	logger logr.Logger
}

// Reconcile the Extension resource.
func (a *actuator) Reconcile(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()

	cluster, err := controller.GetCluster(ctx, a.client, namespace)
	if err != nil {
		return err
	}

	RegistryConfig := &service.RegistryConfig{}
	if ex.Spec.ProviderConfig != nil {
		if _, _, err := a.decoder.Decode(ex.Spec.ProviderConfig.Raw, nil, RegistryConfig); err != nil {
			return fmt.Errorf("failed to decode provider config: %w", err)
		}
		if errs := validation.ValidateRegistryConfig(RegistryConfig, cluster); len(errs) > 0 {
			return errs.ToAggregate()
		}
	}

	if err := a.createResources(ctx, RegistryConfig, cluster, namespace); err != nil {
		return err
	}

	return a.updateStatus(ctx, ex, RegistryConfig)
}

// Delete the Extension resource.
func (a *actuator) Delete(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()
	a.logger.Info("Component is being deleted", "component", "registry-cache", "namespace", namespace)

	return a.deleteResources(ctx, namespace)
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

func (a *actuator) createResources(ctx context.Context, registryConfig *service.RegistryConfig, cluster *controller.Cluster, namespace string) error {
	var mirrors []map[string]interface{}
	for _, m := range registryConfig.Mirrors {
		mirrors = append(mirrors, map[string]interface{}{
			"remoteURL": m.RemoteURL,
			"port":      m.Port,
		})
	}

	values := map[string]interface{}{
		"mirrors": mirrors,
	}

	renderer, err := util.NewChartRendererForShoot(cluster.Shoot.Spec.Kubernetes.Version)
	if err != nil {
		return fmt.Errorf("could not create chart renderer: %w", err)
	}

	return a.createManagedResource(ctx, namespace, v1alpha1.RegistryResourceName, "", renderer, v1alpha1.RegistryChartName, metav1.NamespaceSystem, values, nil)
}

func (a *actuator) deleteResources(ctx context.Context, namespace string) error {
	a.logger.Info("Deleting managed resource for seed", "namespace", namespace)

	if err := managedresources.Delete(ctx, a.client, namespace, v1alpha1.RegistryResourceName, false); err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	return managedresources.WaitUntilDeleted(timeoutCtx, a.client, namespace, v1alpha1.RegistryResourceName)
}

func (a *actuator) createManagedResource(ctx context.Context, namespace, name, class string, renderer chartrenderer.Interface, chartName, chartNamespace string, chartValues map[string]interface{}, injectedLabels map[string]string) error {
	chartPath := filepath.Join(charts.ChartsPath, chartName)
	chart, err := renderer.RenderEmbeddedFS(charts.Internal, chartPath, chartName, chartNamespace, chartValues)
	if err != nil {
		return err
	}

	data := map[string][]byte{chartName: chart.Manifest()}
	keepObjects := false
	forceOverwriteAnnotations := false
	return managedresources.Create(ctx, a.client, namespace, name, false, class, data, &keepObjects, injectedLabels, &forceOverwriteAnnotations)
}

func (a *actuator) updateStatus(ctx context.Context, ex *extensionsv1alpha1.Extension, RegistryConfig *service.RegistryConfig) error {
	patch := client.MergeFrom(ex.DeepCopy())
	// ex.Status.Resources = resources
	return a.client.Status().Patch(ctx, ex, patch)
}
