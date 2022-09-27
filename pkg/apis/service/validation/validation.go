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
	"strings"

	"github.com/gerrit91/gardener-extension-registry-cache/pkg/apis/service"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateRegistryConfig validates the passed configuration instance.
func ValidateRegistryConfig(config *service.RegistryConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	for i, mirror := range config.Mirrors {
		allErrs = append(allErrs, validateRegistry(mirror, fldPath.Child("registries").Index(i))...)
	}

	return allErrs
}

func validateRegistry(registry service.RegistryMirror, fldPath *field.Path) field.ErrorList {
	var (
		allErrs = field.ErrorList{}
	)

	allErrs = append(allErrs, validateRegistryURL(fldPath.Child("remoteURL"), registry.RemoteURL)...)

	if errs := validation.IsValidPortNum(registry.Port); errs != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("port"), registry.Port, "port is not valid: "+strings.Join(errs, ", ")))
	}

	return allErrs
}

func validateRegistryURL(fldPath *field.Path, URL string) field.ErrorList {
	var allErrors field.ErrorList
	const form = "; desired format: https://host[:port]"
	if u, err := url.Parse(URL); err != nil {
		allErrors = append(allErrors, field.Required(fldPath, "url must be a valid URL: "+err.Error()+form))
	} else {
		if u.Scheme != "https" {
			allErrors = append(allErrors, field.Invalid(fldPath, u.Scheme, "'https' is the only allowed URL scheme"+form))
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
