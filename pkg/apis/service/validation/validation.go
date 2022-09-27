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

package validation

import (
	"github.com/gerrit91/gardener-extension-registry-cache/pkg/apis/service"

	"github.com/gardener/gardener/extensions/pkg/controller"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateRegistryConfig validates the passed configuration instance.
func ValidateRegistryConfig(config *service.RegistryConfig, cluster *controller.Cluster) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateRegistries(cluster, nil, field.NewPath("registries"))...)

	return allErrs
}

func validateRegistries(cluster *controller.Cluster, registries any, fldPath *field.Path) field.ErrorList {
	var (
		allErrs = field.ErrorList{}
	)

	return allErrs
}