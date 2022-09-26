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
	"github.com/gardener/gardener/extensions/pkg/controller"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"

	"github.com/Gerrit91/gardener-extension-registry-cache/pkg/apis/service"
	"github.com/Gerrit91/gardener-extension-registry-cache/pkg/apis/service/validation"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// . "github.com/onsi/gomega/gstruct"
	gomegatypes "github.com/onsi/gomega/types"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
)

var _ = Describe("Validation", func() {
	var (
		cluster = &controller.Cluster{
			Shoot: &gardencorev1beta1.Shoot{
				Spec: gardencorev1beta1.ShootSpec{
					Resources: []gardencorev1beta1.NamedResourceReference{
						{
							Name: "testref",
							ResourceRef: autoscalingv1.CrossVersionObjectReference{
								Kind:       "Secret",
								Name:       "referenced-secret",
								APIVersion: "v1",
							},
						},
					},
				},
			},
		}
	)
	DescribeTable("#ValidateRegistryConfig",
		func(config service.RegistryConfig, match gomegatypes.GomegaMatcher) {
			err := validation.ValidateRegistryConfig(&config, cluster)
			Expect(err).To(match)
		},
	)
})
