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

package validation_test

import (
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega/gstruct"

	api "github.com/gerrit91/gardener-extension-registry-cache/pkg/apis/registry"
	. "github.com/gerrit91/gardener-extension-registry-cache/pkg/apis/registry/validation"
)

var _ = Describe("Validation", func() {
	var (
		fldPath = field.NewPath("providerConfig")

		registryConfig *api.RegistryConfig
	)

	BeforeEach(func() {
		registryConfig = &api.RegistryConfig{
			Mirrors: []api.RegistryMirror{{
				UpstreamURL: "https://registry-1.docker.io",
				Port:        5000,
			}},
		}
	})

	Describe("#ValidateRegistryConfig", func() {
		It("should allow valid configuration", func() {
			Expect(ValidateRegistryConfig(registryConfig, fldPath)).To(BeEmpty())
		})

		It("should require upstream URL", func() {
			registryConfig.Mirrors[0].UpstreamURL = ""

			path := fldPath.Child("mirrors").Index(0).Child("upstreamURL").String()
			Expect(ValidateRegistryConfig(registryConfig, fldPath)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeRequired),
					"Field":  Equal(path),
					"Detail": ContainSubstring("upstream URL must be provided"),
				})),
			))
		})

		It("should deny invalid upstream URL", func() {
			registryConfig.Mirrors[0].UpstreamURL = "https://registry-1:docker:io/"

			path := fldPath.Child("mirrors").Index(0).Child("upstreamURL").String()
			Expect(ValidateRegistryConfig(registryConfig, fldPath)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal(path),
					"Detail": ContainSubstring("must be a valid URL"),
				})),
			))
		})

		It("should deny upstream URL without host", func() {
			registryConfig.Mirrors[0].UpstreamURL = "https://"

			path := fldPath.Child("mirrors").Index(0).Child("upstreamURL").String()
			Expect(ValidateRegistryConfig(registryConfig, fldPath)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal(path),
					"Detail": ContainSubstring("host must be provided"),
				})),
			))
		})

		It("should deny upstream URL with non-permitted parts", func() {
			registryConfig.Mirrors[0].UpstreamURL = "http://user:pass@registry-1.docker.io/path/to/foo?query#fragment"

			path := fldPath.Child("mirrors").Index(0).Child("upstreamURL").String()
			Expect(ValidateRegistryConfig(registryConfig, fldPath)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal(path),
					"Detail": ContainSubstring("https' is the only allowed URL scheme"),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal(path),
					"Detail": ContainSubstring("path is not permitted"),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal(path),
					"Detail": ContainSubstring("user information is not permitted"),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal(path),
					"Detail": ContainSubstring("fragments are not permitted"),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal(path),
					"Detail": ContainSubstring("query parameters are not permitted"),
				})),
			))
		})

		It("should deny invalid port numbers", func() {
			registryConfig.Mirrors = append(registryConfig.Mirrors, *registryConfig.Mirrors[0].DeepCopy())
			registryConfig.Mirrors = append(registryConfig.Mirrors, *registryConfig.Mirrors[0].DeepCopy())
			registryConfig.Mirrors = append(registryConfig.Mirrors, *registryConfig.Mirrors[0].DeepCopy())
			registryConfig.Mirrors[0].Port = 0
			registryConfig.Mirrors[1].Port = -1
			registryConfig.Mirrors[2].Port = 65536

			Expect(ValidateRegistryConfig(registryConfig, fldPath)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal(fldPath.Child("mirrors").Index(0).Child("port").String()),
					"Detail": ContainSubstring("port is invalid"),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal(fldPath.Child("mirrors").Index(1).Child("port").String()),
					"Detail": ContainSubstring("port is invalid"),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal(fldPath.Child("mirrors").Index(2).Child("port").String()),
					"Detail": ContainSubstring("port is invalid"),
				})),
			))
		})
	})
})
