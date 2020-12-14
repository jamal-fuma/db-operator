// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AmazonInstance) DeepCopyInto(out *AmazonInstance) {
	*out = *in
	out.Generic = in.Generic
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AmazonInstance.
func (in *AmazonInstance) DeepCopy() *AmazonInstance {
	if in == nil {
		return nil
	}
	out := new(AmazonInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackendServer) DeepCopyInto(out *BackendServer) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackendServer.
func (in *BackendServer) DeepCopy() *BackendServer {
	if in == nil {
		return nil
	}
	out := new(BackendServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Database) DeepCopyInto(out *Database) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Database.
func (in *Database) DeepCopy() *Database {
	if in == nil {
		return nil
	}
	out := new(Database)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Database) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseBackup) DeepCopyInto(out *DatabaseBackup) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseBackup.
func (in *DatabaseBackup) DeepCopy() *DatabaseBackup {
	if in == nil {
		return nil
	}
	out := new(DatabaseBackup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseList) DeepCopyInto(out *DatabaseList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Database, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseList.
func (in *DatabaseList) DeepCopy() *DatabaseList {
	if in == nil {
		return nil
	}
	out := new(DatabaseList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DatabaseList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseProxyStatus) DeepCopyInto(out *DatabaseProxyStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseProxyStatus.
func (in *DatabaseProxyStatus) DeepCopy() *DatabaseProxyStatus {
	if in == nil {
		return nil
	}
	out := new(DatabaseProxyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseSpec) DeepCopyInto(out *DatabaseSpec) {
	*out = *in
	out.Backup = in.Backup
	if in.Extensions != nil {
		in, out := &in.Extensions, &out.Extensions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseSpec.
func (in *DatabaseSpec) DeepCopy() *DatabaseSpec {
	if in == nil {
		return nil
	}
	out := new(DatabaseSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseStatus) DeepCopyInto(out *DatabaseStatus) {
	*out = *in
	if in.InstanceRef != nil {
		in, out := &in.InstanceRef, &out.InstanceRef
		*out = new(DbInstance)
		(*in).DeepCopyInto(*out)
	}
	out.ProxyStatus = in.ProxyStatus
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseStatus.
func (in *DatabaseStatus) DeepCopy() *DatabaseStatus {
	if in == nil {
		return nil
	}
	out := new(DatabaseStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstance) DeepCopyInto(out *DbInstance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstance.
func (in *DbInstance) DeepCopy() *DbInstance {
	if in == nil {
		return nil
	}
	out := new(DbInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DbInstance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceBackup) DeepCopyInto(out *DbInstanceBackup) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceBackup.
func (in *DbInstanceBackup) DeepCopy() *DbInstanceBackup {
	if in == nil {
		return nil
	}
	out := new(DbInstanceBackup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceList) DeepCopyInto(out *DbInstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DbInstance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceList.
func (in *DbInstanceList) DeepCopy() *DbInstanceList {
	if in == nil {
		return nil
	}
	out := new(DbInstanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DbInstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceMonitoring) DeepCopyInto(out *DbInstanceMonitoring) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceMonitoring.
func (in *DbInstanceMonitoring) DeepCopy() *DbInstanceMonitoring {
	if in == nil {
		return nil
	}
	out := new(DbInstanceMonitoring)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceSSLConnection) DeepCopyInto(out *DbInstanceSSLConnection) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceSSLConnection.
func (in *DbInstanceSSLConnection) DeepCopy() *DbInstanceSSLConnection {
	if in == nil {
		return nil
	}
	out := new(DbInstanceSSLConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceSource) DeepCopyInto(out *DbInstanceSource) {
	*out = *in
	if in.Google != nil {
		in, out := &in.Google, &out.Google
		*out = new(GoogleInstance)
		**out = **in
	}
	if in.Generic != nil {
		in, out := &in.Generic, &out.Generic
		*out = new(GenericInstance)
		**out = **in
	}
	if in.Percona != nil {
		in, out := &in.Percona, &out.Percona
		*out = new(PerconaCluster)
		(*in).DeepCopyInto(*out)
	}
	if in.Amazon != nil {
		in, out := &in.Amazon, &out.Amazon
		*out = new(AmazonInstance)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceSource.
func (in *DbInstanceSource) DeepCopy() *DbInstanceSource {
	if in == nil {
		return nil
	}
	out := new(DbInstanceSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceSpec) DeepCopyInto(out *DbInstanceSpec) {
	*out = *in
	out.AdminUserSecret = in.AdminUserSecret
	out.Backup = in.Backup
	out.Monitoring = in.Monitoring
	out.SSLConnection = in.SSLConnection
	in.DbInstanceSource.DeepCopyInto(&out.DbInstanceSource)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceSpec.
func (in *DbInstanceSpec) DeepCopy() *DbInstanceSpec {
	if in == nil {
		return nil
	}
	out := new(DbInstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceStatus) DeepCopyInto(out *DbInstanceStatus) {
	*out = *in
	if in.Info != nil {
		in, out := &in.Info, &out.Info
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Checksums != nil {
		in, out := &in.Checksums, &out.Checksums
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceStatus.
func (in *DbInstanceStatus) DeepCopy() *DbInstanceStatus {
	if in == nil {
		return nil
	}
	out := new(DbInstanceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GenericInstance) DeepCopyInto(out *GenericInstance) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GenericInstance.
func (in *GenericInstance) DeepCopy() *GenericInstance {
	if in == nil {
		return nil
	}
	out := new(GenericInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleInstance) DeepCopyInto(out *GoogleInstance) {
	*out = *in
	out.ConfigmapName = in.ConfigmapName
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleInstance.
func (in *GoogleInstance) DeepCopy() *GoogleInstance {
	if in == nil {
		return nil
	}
	out := new(GoogleInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PerconaCluster) DeepCopyInto(out *PerconaCluster) {
	*out = *in
	if in.ServerList != nil {
		in, out := &in.ServerList, &out.ServerList
		*out = make([]BackendServer, len(*in))
		copy(*out, *in)
	}
	out.MonitorUserSecret = in.MonitorUserSecret
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PerconaCluster.
func (in *PerconaCluster) DeepCopy() *PerconaCluster {
	if in == nil {
		return nil
	}
	out := new(PerconaCluster)
	in.DeepCopyInto(out)
	return out
}
