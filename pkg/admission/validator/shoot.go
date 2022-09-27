// Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package validator

import (
	"context"
	"fmt"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/core"
	"github.com/gerrit91/gardener-extension-registry-cache/pkg/controller"
	"github.com/gerrit91/gardener-extension-registry-cache/pkg/apis/service/validation"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// shoot validates shoots
type shoot struct {
	client         client.Client
	decoder        runtime.Decoder
}

// NewShootValidator returns a new instance of a shoot validator.
func NewShootValidator() extensionswebhook.Validator {
	return &shoot{}
}

// InjectScheme injects the given scheme into the validator.
func (s *shoot) InjectScheme(scheme *runtime.Scheme) error {
	s.decoder = serializer.NewCodecFactory(scheme, serializer.EnableStrict).UniversalDecoder()
	return nil
}

// InjectClient injects the given client into the validator.
func (s *shoot) InjectClient(client client.Client) error {
	s.client = client
	return nil
}

// Validate validates the given shoot object
func (s *shoot) Validate(ctx context.Context, new, old client.Object) error {
	shoot, ok := new.(*core.Shoot)
	if !ok {
		return fmt.Errorf("wrong object type %T", new)
	}

	var ext *core.Extension
	var fldPath *field.Path
	for i, ex := range shoot.Spec.Extensions {
		if ex.Type == controller.Type {
			ext = ex.DeepCopy()
			fldPath = field.NewPath("spec", "extensions").Index(i)
			break
		} 
	}
	if ext == nil {
		return nil
	}

	providerConfigPath := fldPath.Child("providerConfig")

	registryConfig, err := decodeRegistryConfig(s.decoder, ext.ProviderConfig, providerConfigPath)
	if err != nil {
		return err
	}

	return validation.ValidateRegistryConfig(registryConfig, providerConfigPath).ToAggregate()

}

