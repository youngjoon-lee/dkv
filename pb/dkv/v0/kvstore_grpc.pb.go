// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: dkv/v0/kvstore.proto

package v0

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// KVStoreClient is the client API for KVStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KVStoreClient interface {
	Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutReply, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
}

type kVStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewKVStoreClient(cc grpc.ClientConnInterface) KVStoreClient {
	return &kVStoreClient{cc}
}

func (c *kVStoreClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutReply, error) {
	out := new(PutReply)
	err := c.cc.Invoke(ctx, "/dkv.v0.KVStore/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVStoreClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := c.cc.Invoke(ctx, "/dkv.v0.KVStore/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KVStoreServer is the server API for KVStore service.
// All implementations must embed UnimplementedKVStoreServer
// for forward compatibility
type KVStoreServer interface {
	Put(context.Context, *PutRequest) (*PutReply, error)
	Get(context.Context, *GetRequest) (*GetReply, error)
	mustEmbedUnimplementedKVStoreServer()
}

// UnimplementedKVStoreServer must be embedded to have forward compatible implementations.
type UnimplementedKVStoreServer struct {
}

func (UnimplementedKVStoreServer) Put(context.Context, *PutRequest) (*PutReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedKVStoreServer) Get(context.Context, *GetRequest) (*GetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedKVStoreServer) mustEmbedUnimplementedKVStoreServer() {}

// UnsafeKVStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KVStoreServer will
// result in compilation errors.
type UnsafeKVStoreServer interface {
	mustEmbedUnimplementedKVStoreServer()
}

func RegisterKVStoreServer(s grpc.ServiceRegistrar, srv KVStoreServer) {
	s.RegisterService(&KVStore_ServiceDesc, srv)
}

func _KVStore_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVStoreServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dkv.v0.KVStore/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVStoreServer).Put(ctx, req.(*PutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVStore_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVStoreServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dkv.v0.KVStore/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVStoreServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KVStore_ServiceDesc is the grpc.ServiceDesc for KVStore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KVStore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dkv.v0.KVStore",
	HandlerType: (*KVStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _KVStore_Put_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _KVStore_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dkv/v0/kvstore.proto",
}
