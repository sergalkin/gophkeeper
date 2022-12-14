// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: api/proto/secret_type.proto

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

// SecretTypeClient is the client API for SecretType service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SecretTypeClient interface {
	GetSecretTypesList(ctx context.Context, in *SecretTypesListRequest, opts ...grpc.CallOption) (*SecretTypesListResponse, error)
}

type secretTypeClient struct {
	cc grpc.ClientConnInterface
}

func NewSecretTypeClient(cc grpc.ClientConnInterface) SecretTypeClient {
	return &secretTypeClient{cc}
}

func (c *secretTypeClient) GetSecretTypesList(ctx context.Context, in *SecretTypesListRequest, opts ...grpc.CallOption) (*SecretTypesListResponse, error) {
	out := new(SecretTypesListResponse)
	err := c.cc.Invoke(ctx, "/proto.SecretType/GetSecretTypesList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SecretTypeServer is the server API for SecretType service.
// All implementations must embed UnimplementedSecretTypeServer
// for forward compatibility
type SecretTypeServer interface {
	GetSecretTypesList(context.Context, *SecretTypesListRequest) (*SecretTypesListResponse, error)
	mustEmbedUnimplementedSecretTypeServer()
}

// UnimplementedSecretTypeServer must be embedded to have forward compatible implementations.
type UnimplementedSecretTypeServer struct {
}

func (UnimplementedSecretTypeServer) GetSecretTypesList(context.Context, *SecretTypesListRequest) (*SecretTypesListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSecretTypesList not implemented")
}
func (UnimplementedSecretTypeServer) mustEmbedUnimplementedSecretTypeServer() {}

// UnsafeSecretTypeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SecretTypeServer will
// result in compilation errors.
type UnsafeSecretTypeServer interface {
	mustEmbedUnimplementedSecretTypeServer()
}

func RegisterSecretTypeServer(s grpc.ServiceRegistrar, srv SecretTypeServer) {
	s.RegisterService(&SecretType_ServiceDesc, srv)
}

func _SecretType_GetSecretTypesList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecretTypesListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretTypeServer).GetSecretTypesList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.SecretType/GetSecretTypesList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretTypeServer).GetSecretTypesList(ctx, req.(*SecretTypesListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SecretType_ServiceDesc is the grpc.ServiceDesc for SecretType service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SecretType_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.SecretType",
	HandlerType: (*SecretTypeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSecretTypesList",
			Handler:    _SecretType_GetSecretTypesList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/secret_type.proto",
}
