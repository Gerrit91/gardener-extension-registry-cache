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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RegistryResourceName is the name for registry resources in the shoot.
const RegistryResourceName = "extension-registry-cache"

// RegistryChartNameSeed is the name of the chart for registry in the seed.
const RegistryChartName = "registry-cache"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RegistryConfig configuration resource
type RegistryConfig struct {
	metav1.TypeMeta `json:",inline"`

	// Mirrors is a slice of registry mirrors to deploy
	Mirrors []RegistryMirror `json:"registries"`
}

// RegistryMirror defines a registry mirror to deploy
type RegistryMirror struct {
	// RemoteURL is the remote URL of registry to mirror
	RemoteURL string `json:"remoteURL"`
	// Port is the port on which the registry mirror is going to serve
	Port int `json:"port"`
	// CacheSize is the size of the registry cache
	// +optional
	CacheSize *string `json:"cacheSize,omitempty"`
	// CacheGarbageCollectionEnabled enables/disables cache garbage collection
	// +optional
	CacheGarbageCollectionEnabled *bool `json:"cacheGarbageCollectionEnabled,omitempty"`
}
