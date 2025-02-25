//go:build !ignore_autogenerated

/*
Copyright 2023.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIServer) DeepCopyInto(out *APIServer) {
	*out = *in
	in.OidcConfig.DeepCopyInto(&out.OidcConfig)
	if in.AdditionalOidcConfig != nil {
		in, out := &in.AdditionalOidcConfig, &out.AdditionalOidcConfig
		*out = new([]v1beta1.OIDCConfig)
		if **in != nil {
			in, out := *in, *out
			*out = make([]v1beta1.OIDCConfig, len(*in))
			for i := range *in {
				(*in)[i].DeepCopyInto(&(*out)[i])
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIServer.
func (in *APIServer) DeepCopy() *APIServer {
	if in == nil {
		return nil
	}
	out := new(APIServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomConfig) DeepCopyInto(out *CustomConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomConfig.
func (in *CustomConfig) DeepCopy() *CustomConfig {
	if in == nil {
		return nil
	}
	out := new(CustomConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomConfigList) DeepCopyInto(out *CustomConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CustomConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomConfigList.
func (in *CustomConfigList) DeepCopy() *CustomConfigList {
	if in == nil {
		return nil
	}
	out := new(CustomConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CustomConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomConfigSpec) DeepCopyInto(out *CustomConfigSpec) {
	*out = *in
	in.RegistryCache.DeepCopyInto(&out.RegistryCache)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomConfigSpec.
func (in *CustomConfigSpec) DeepCopy() *CustomConfigSpec {
	if in == nil {
		return nil
	}
	out := new(CustomConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Egress) DeepCopyInto(out *Egress) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Egress.
func (in *Egress) DeepCopy() *Egress {
	if in == nil {
		return nil
	}
	out := new(Egress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Filter) DeepCopyInto(out *Filter) {
	*out = *in
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(Ingress)
		**out = **in
	}
	out.Egress = in.Egress
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Filter.
func (in *Filter) DeepCopy() *Filter {
	if in == nil {
		return nil
	}
	out := new(Filter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GardenerCluster) DeepCopyInto(out *GardenerCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GardenerCluster.
func (in *GardenerCluster) DeepCopy() *GardenerCluster {
	if in == nil {
		return nil
	}
	out := new(GardenerCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GardenerCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GardenerClusterList) DeepCopyInto(out *GardenerClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GardenerCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GardenerClusterList.
func (in *GardenerClusterList) DeepCopy() *GardenerClusterList {
	if in == nil {
		return nil
	}
	out := new(GardenerClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GardenerClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GardenerClusterSpec) DeepCopyInto(out *GardenerClusterSpec) {
	*out = *in
	out.Kubeconfig = in.Kubeconfig
	out.Shoot = in.Shoot
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GardenerClusterSpec.
func (in *GardenerClusterSpec) DeepCopy() *GardenerClusterSpec {
	if in == nil {
		return nil
	}
	out := new(GardenerClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GardenerClusterStatus) DeepCopyInto(out *GardenerClusterStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GardenerClusterStatus.
func (in *GardenerClusterStatus) DeepCopy() *GardenerClusterStatus {
	if in == nil {
		return nil
	}
	out := new(GardenerClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageRegistryCache) DeepCopyInto(out *ImageRegistryCache) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageRegistryCache.
func (in *ImageRegistryCache) DeepCopy() *ImageRegistryCache {
	if in == nil {
		return nil
	}
	out := new(ImageRegistryCache)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ingress) DeepCopyInto(out *Ingress) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ingress.
func (in *Ingress) DeepCopy() *Ingress {
	if in == nil {
		return nil
	}
	out := new(Ingress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Kubeconfig) DeepCopyInto(out *Kubeconfig) {
	*out = *in
	out.Secret = in.Secret
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Kubeconfig.
func (in *Kubeconfig) DeepCopy() *Kubeconfig {
	if in == nil {
		return nil
	}
	out := new(Kubeconfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Kubernetes) DeepCopyInto(out *Kubernetes) {
	*out = *in
	if in.Version != nil {
		in, out := &in.Version, &out.Version
		*out = new(string)
		**out = **in
	}
	in.KubeAPIServer.DeepCopyInto(&out.KubeAPIServer)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Kubernetes.
func (in *Kubernetes) DeepCopy() *Kubernetes {
	if in == nil {
		return nil
	}
	out := new(Kubernetes)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Networking) DeepCopyInto(out *Networking) {
	*out = *in
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Networking.
func (in *Networking) DeepCopy() *Networking {
	if in == nil {
		return nil
	}
	out := new(Networking)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkingSecurity) DeepCopyInto(out *NetworkingSecurity) {
	*out = *in
	in.Filter.DeepCopyInto(&out.Filter)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkingSecurity.
func (in *NetworkingSecurity) DeepCopy() *NetworkingSecurity {
	if in == nil {
		return nil
	}
	out := new(NetworkingSecurity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Provider) DeepCopyInto(out *Provider) {
	*out = *in
	if in.Workers != nil {
		in, out := &in.Workers, &out.Workers
		*out = make([]v1beta1.Worker, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AdditionalWorkers != nil {
		in, out := &in.AdditionalWorkers, &out.AdditionalWorkers
		*out = new([]v1beta1.Worker)
		if **in != nil {
			in, out := *in, *out
			*out = make([]v1beta1.Worker, len(*in))
			for i := range *in {
				(*in)[i].DeepCopyInto(&(*out)[i])
			}
		}
	}
	if in.ControlPlaneConfig != nil {
		in, out := &in.ControlPlaneConfig, &out.ControlPlaneConfig
		*out = new(runtime.RawExtension)
		(*in).DeepCopyInto(*out)
	}
	if in.InfrastructureConfig != nil {
		in, out := &in.InfrastructureConfig, &out.InfrastructureConfig
		*out = new(runtime.RawExtension)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Provider.
func (in *Provider) DeepCopy() *Provider {
	if in == nil {
		return nil
	}
	out := new(Provider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Runtime) DeepCopyInto(out *Runtime) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Runtime.
func (in *Runtime) DeepCopy() *Runtime {
	if in == nil {
		return nil
	}
	out := new(Runtime)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Runtime) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuntimeList) DeepCopyInto(out *RuntimeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Runtime, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuntimeList.
func (in *RuntimeList) DeepCopy() *RuntimeList {
	if in == nil {
		return nil
	}
	out := new(RuntimeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RuntimeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuntimeShoot) DeepCopyInto(out *RuntimeShoot) {
	*out = *in
	if in.LicenceType != nil {
		in, out := &in.LicenceType, &out.LicenceType
		*out = new(string)
		**out = **in
	}
	if in.EnforceSeedLocation != nil {
		in, out := &in.EnforceSeedLocation, &out.EnforceSeedLocation
		*out = new(bool)
		**out = **in
	}
	in.Kubernetes.DeepCopyInto(&out.Kubernetes)
	in.Provider.DeepCopyInto(&out.Provider)
	in.Networking.DeepCopyInto(&out.Networking)
	if in.ControlPlane != nil {
		in, out := &in.ControlPlane, &out.ControlPlane
		*out = new(v1beta1.ControlPlane)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuntimeShoot.
func (in *RuntimeShoot) DeepCopy() *RuntimeShoot {
	if in == nil {
		return nil
	}
	out := new(RuntimeShoot)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuntimeSpec) DeepCopyInto(out *RuntimeSpec) {
	*out = *in
	in.Shoot.DeepCopyInto(&out.Shoot)
	in.Security.DeepCopyInto(&out.Security)
	out.Caching = in.Caching
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuntimeSpec.
func (in *RuntimeSpec) DeepCopy() *RuntimeSpec {
	if in == nil {
		return nil
	}
	out := new(RuntimeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuntimeStatus) DeepCopyInto(out *RuntimeStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuntimeStatus.
func (in *RuntimeStatus) DeepCopy() *RuntimeStatus {
	if in == nil {
		return nil
	}
	out := new(RuntimeStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Secret) DeepCopyInto(out *Secret) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Secret.
func (in *Secret) DeepCopy() *Secret {
	if in == nil {
		return nil
	}
	out := new(Secret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Security) DeepCopyInto(out *Security) {
	*out = *in
	if in.Administrators != nil {
		in, out := &in.Administrators, &out.Administrators
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.Networking.DeepCopyInto(&out.Networking)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Security.
func (in *Security) DeepCopy() *Security {
	if in == nil {
		return nil
	}
	out := new(Security)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Shoot) DeepCopyInto(out *Shoot) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Shoot.
func (in *Shoot) DeepCopy() *Shoot {
	if in == nil {
		return nil
	}
	out := new(Shoot)
	in.DeepCopyInto(out)
	return out
}
