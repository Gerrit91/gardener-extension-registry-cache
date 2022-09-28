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
	"net/url"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/gerrit91/gardener-extension-registry-cache/pkg/apis/registry"
)

// ValidateRegistryConfig validates the passed configuration instance.
func ValidateRegistryConfig(config *registry.RegistryConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	for i, cache := range config.Caches {
		allErrs = append(allErrs, validateRegistry(cache, fldPath.Child("caches").Index(i))...)
	}

	return allErrs
}

func validateRegistry(registry registry.RegistryCache, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validateUpstream(fldPath.Child("upstreamURL"), registry.Upstream)...)

	return allErrs
}

func validateUpstream(fldPath *field.Path, upstream string) field.ErrorList {
	var allErrors field.ErrorList

	const form = "; desired format: host[:port]"
	if len(upstream) == 0 {
		allErrors = append(allErrors, field.Required(fldPath, "upstream must be provided"+form))
		return allErrors
	}

	if u, err := url.Parse(upstream); err != nil {
		allErrors = append(allErrors, field.Invalid(fldPath, upstream, "url must be a valid URL: "+err.Error()+form))
	} else {
		if u.Scheme != "" {
			allErrors = append(allErrors, field.Invalid(fldPath, u.Scheme, "scheme is not permitted in the registry URL"+form))
		}
		if len(u.Path) != 0 {
			allErrors = append(allErrors, field.Invalid(fldPath, u.Path, "path is not permitted in the registry URL"+form))
		}
		if len(u.Host) == 0 {
			allErrors = append(allErrors, field.Invalid(fldPath, u.Host, "host must be provided in the registry URL"+form))
		}
		if u.User != nil {
			allErrors = append(allErrors, field.Invalid(fldPath, u.User.String(), "user information is not permitted in the registry URL"))
		}
		if len(u.Fragment) != 0 {
			allErrors = append(allErrors, field.Invalid(fldPath, u.Fragment, "fragments are not permitted in the registry URL"))
		}
		if len(u.RawQuery) != 0 {
			allErrors = append(allErrors, field.Invalid(fldPath, u.RawQuery, "query parameters are not permitted in the registry URL"))
		}
	}

	return allErrors
}
