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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RegistryResourceName is the name for registry resources in the shoot.
const RegistryResourceName = "extension-registry-cache"

// RegistryEnsurerResourceName is the name for registry cri ensurer resources in the shoot.
const RegistryEnsurerResourceName = "extension-registry-cache-cri-ensurer"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RegistryConfig configuration resource
type RegistryConfig struct {
	metav1.TypeMeta `json:",inline"`

	// Caches is a slice of registry cache to deploy
	Caches []RegistryCache `json:"caches"`
}

// RegistryCache defines a registry cache to deploy
type RegistryCache struct {
	// Upstream is the remote registry host (and optionally port) to cache
	Upstream string `json:"upstream"`
	// Size is the size of the registry cache, defaults to 10Gi.
	// +optional
	Size *resource.Quantity `json:"size,omitempty"`
	// GarbageCollectionEnabled enables/disables cache garbage collection, defaults to true.
	// +optional
	GarbageCollectionEnabled *bool `json:"garbageCollectionEnabled,omitempty"`
	// StorageClassName specifies the storageclass to use for the cache's volume
	// +optional
	StorageClassName *string `json:"storageClassName,omitempty"`
}
