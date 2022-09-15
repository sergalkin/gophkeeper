// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: api/proto/secret.proto

package proto

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

// SecretClient is the client API for Secret service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SecretClient interface {
	CreateSecret(ctx context.Context, in *CreateSecretRequest, opts ...grpc.CallOption) (*CreateSecretResponse, error)
	GetSecret(ctx context.Context, in *GetSecretRequest, opts ...grpc.CallOption) (*GetSecretResponse, error)
}

type secretClient struct {
	cc grpc.ClientConnInterface
}

func NewSecretClient(cc grpc.ClientConnInterface) SecretClient {
	return &secretClient{cc}
}

func (c *secretClient) CreateSecret(ctx context.Context, in *CreateSecretRequest, opts ...grpc.CallOption) (*CreateSecretResponse, error) {
	out := new(CreateSecretResponse)
	err := c.cc.Invoke(ctx, "/proto.Secret/CreateSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretClient) GetSecret(ctx context.Context, in *GetSecretRequest, opts ...grpc.CallOption) (*GetSecretResponse, error) {
	out := new(GetSecretResponse)
	err := c.cc.Invoke(ctx, "/proto.Secret/GetSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SecretServer is the server API for Secret service.
// All implementations must embed UnimplementedSecretServer
// for forward compatibility
type SecretServer interface {
	CreateSecret(context.Context, *CreateSecretRequest) (*CreateSecretResponse, error)
	GetSecret(context.Context, *GetSecretRequest) (*GetSecretResponse, error)
	mustEmbedUnimplementedSecretServer()
}

// UnimplementedSecretServer must be embedded to have forward compatible implementations.
type UnimplementedSecretServer struct {
}

func (UnimplementedSecretServer) CreateSecret(context.Context, *CreateSecretRequest) (*CreateSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSecret not implemented")
}
func (UnimplementedSecretServer) GetSecret(context.Context, *GetSecretRequest) (*GetSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSecret not implemented")
}
func (UnimplementedSecretServer) mustEmbedUnimplementedSecretServer() {}

// UnsafeSecretServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SecretServer will
// result in compilation errors.
type UnsafeSecretServer interface {
	mustEmbedUnimplementedSecretServer()
}

func RegisterSecretServer(s grpc.ServiceRegistrar, srv SecretServer) {
	s.RegisterService(&Secret_ServiceDesc, srv)
}

func _Secret_CreateSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretServer).CreateSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Secret/CreateSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretServer).CreateSecret(ctx, req.(*CreateSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Secret_GetSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretServer).GetSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Secret/GetSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretServer).GetSecret(ctx, req.(*GetSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Secret_ServiceDesc is the grpc.ServiceDesc for Secret service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Secret_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Secret",
	HandlerType: (*SecretServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateSecret",
			Handler:    _Secret_CreateSecret_Handler,
		},
		{
			MethodName: "GetSecret",
			Handler:    _Secret_GetSecret_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/secret.proto",
}
