// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.10.1
// source: tgo.proto

package v1

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

// TgoClient is the client API for Tgo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TgoClient interface {
	Tg(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type tgoClient struct {
	cc grpc.ClientConnInterface
}

func NewTgoClient(cc grpc.ClientConnInterface) TgoClient {
	return &tgoClient{cc}
}

func (c *tgoClient) Tg(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/pto.Tgo/Tg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TgoServer is the server API for Tgo service.
// All implementations must embed UnimplementedTgoServer
// for forward compatibility
type TgoServer interface {
	Tg(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedTgoServer()
}

// UnimplementedTgoServer must be embedded to have forward compatible implementations.
type UnimplementedTgoServer struct {
}

func (UnimplementedTgoServer) Tg(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Tg not implemented")
}
func (UnimplementedTgoServer) mustEmbedUnimplementedTgoServer() {}

// UnsafeTgoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TgoServer will
// result in compilation errors.
type UnsafeTgoServer interface {
	mustEmbedUnimplementedTgoServer()
}

func RegisterTgoServer(s grpc.ServiceRegistrar, srv TgoServer) {
	s.RegisterService(&Tgo_ServiceDesc, srv)
}

func _Tgo_Tg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TgoServer).Tg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pto.Tgo/Tg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TgoServer).Tg(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Tgo_ServiceDesc is the grpc.ServiceDesc for Tgo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Tgo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pto.Tgo",
	HandlerType: (*TgoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Tg",
			Handler:    _Tgo_Tg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tgo.proto",
}