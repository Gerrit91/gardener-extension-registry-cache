//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright (c) SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	service "github.com/gerrit91/gardener-extension-registry-cache/pkg/apis/service"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*RegistryConfig)(nil), (*service.RegistryConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_RegistryConfig_To_service_RegistryConfig(a.(*RegistryConfig), b.(*service.RegistryConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*service.RegistryConfig)(nil), (*RegistryConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_service_RegistryConfig_To_v1alpha1_RegistryConfig(a.(*service.RegistryConfig), b.(*RegistryConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*RegistryMirror)(nil), (*service.RegistryMirror)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_RegistryMirror_To_service_RegistryMirror(a.(*RegistryMirror), b.(*service.RegistryMirror), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*service.RegistryMirror)(nil), (*RegistryMirror)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_service_RegistryMirror_To_v1alpha1_RegistryMirror(a.(*service.RegistryMirror), b.(*RegistryMirror), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_RegistryConfig_To_service_RegistryConfig(in *RegistryConfig, out *service.RegistryConfig, s conversion.Scope) error {
	out.Mirrors = *(*[]service.RegistryMirror)(unsafe.Pointer(&in.Mirrors))
	return nil
}

// Convert_v1alpha1_RegistryConfig_To_service_RegistryConfig is an autogenerated conversion function.
func Convert_v1alpha1_RegistryConfig_To_service_RegistryConfig(in *RegistryConfig, out *service.RegistryConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_RegistryConfig_To_service_RegistryConfig(in, out, s)
}

func autoConvert_service_RegistryConfig_To_v1alpha1_RegistryConfig(in *service.RegistryConfig, out *RegistryConfig, s conversion.Scope) error {
	out.Mirrors = *(*[]RegistryMirror)(unsafe.Pointer(&in.Mirrors))
	return nil
}

// Convert_service_RegistryConfig_To_v1alpha1_RegistryConfig is an autogenerated conversion function.
func Convert_service_RegistryConfig_To_v1alpha1_RegistryConfig(in *service.RegistryConfig, out *RegistryConfig, s conversion.Scope) error {
	return autoConvert_service_RegistryConfig_To_v1alpha1_RegistryConfig(in, out, s)
}

func autoConvert_v1alpha1_RegistryMirror_To_service_RegistryMirror(in *RegistryMirror, out *service.RegistryMirror, s conversion.Scope) error {
	out.RemoteURL = in.RemoteURL
	out.Port = in.Port
	return nil
}

// Convert_v1alpha1_RegistryMirror_To_service_RegistryMirror is an autogenerated conversion function.
func Convert_v1alpha1_RegistryMirror_To_service_RegistryMirror(in *RegistryMirror, out *service.RegistryMirror, s conversion.Scope) error {
	return autoConvert_v1alpha1_RegistryMirror_To_service_RegistryMirror(in, out, s)
}

func autoConvert_service_RegistryMirror_To_v1alpha1_RegistryMirror(in *service.RegistryMirror, out *RegistryMirror, s conversion.Scope) error {
	out.RemoteURL = in.RemoteURL
	out.Port = in.Port
	return nil
}

// Convert_service_RegistryMirror_To_v1alpha1_RegistryMirror is an autogenerated conversion function.
func Convert_service_RegistryMirror_To_v1alpha1_RegistryMirror(in *service.RegistryMirror, out *RegistryMirror, s conversion.Scope) error {
	return autoConvert_service_RegistryMirror_To_v1alpha1_RegistryMirror(in, out, s)
}
