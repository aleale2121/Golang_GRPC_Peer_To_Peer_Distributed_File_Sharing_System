// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package connect

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// LiveConnectionClient is the client API for LiveConnection service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LiveConnectionClient interface {
	Subscribe(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (LiveConnection_SubscribeClient, error)
	Unsubscribe(ctx context.Context, in *UnSubscribeRequest, opts ...grpc.CallOption) (*Response, error)
}

type liveConnectionClient struct {
	cc grpc.ClientConnInterface
}

func NewLiveConnectionClient(cc grpc.ClientConnInterface) LiveConnectionClient {
	return &liveConnectionClient{cc}
}

func (c *liveConnectionClient) Subscribe(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (LiveConnection_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_LiveConnection_serviceDesc.Streams[0], "/protos.LiveConnection/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &liveConnectionSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type LiveConnection_SubscribeClient interface {
	Recv() (*Response, error)
	grpc.ClientStream
}

type liveConnectionSubscribeClient struct {
	grpc.ClientStream
}

func (x *liveConnectionSubscribeClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *liveConnectionClient) Unsubscribe(ctx context.Context, in *UnSubscribeRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/protos.LiveConnection/Unsubscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LiveConnectionServer is the server API for LiveConnection service.
// All implementations must embed UnimplementedLiveConnectionServer
// for forward compatibility
type LiveConnectionServer interface {
	Subscribe(*ConnectRequest, LiveConnection_SubscribeServer) error
	Unsubscribe(context.Context, *UnSubscribeRequest) (*Response, error)
	mustEmbedUnimplementedLiveConnectionServer()
}

// UnimplementedLiveConnectionServer must be embedded to have forward compatible implementations.
type UnimplementedLiveConnectionServer struct {
}

func (*UnimplementedLiveConnectionServer) Subscribe(*ConnectRequest, LiveConnection_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (*UnimplementedLiveConnectionServer) Unsubscribe(context.Context, *UnSubscribeRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unsubscribe not implemented")
}
func (*UnimplementedLiveConnectionServer) mustEmbedUnimplementedLiveConnectionServer() {}

func RegisterLiveConnectionServer(s *grpc.Server, srv LiveConnectionServer) {
	s.RegisterService(&_LiveConnection_serviceDesc, srv)
}

func _LiveConnection_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConnectRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LiveConnectionServer).Subscribe(m, &liveConnectionSubscribeServer{stream})
}

type LiveConnection_SubscribeServer interface {
	Send(*Response) error
	grpc.ServerStream
}

type liveConnectionSubscribeServer struct {
	grpc.ServerStream
}

func (x *liveConnectionSubscribeServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func _LiveConnection_Unsubscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnSubscribeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveConnectionServer).Unsubscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.LiveConnection/Unsubscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveConnectionServer).Unsubscribe(ctx, req.(*UnSubscribeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _LiveConnection_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.LiveConnection",
	HandlerType: (*LiveConnectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Unsubscribe",
			Handler:    _LiveConnection_Unsubscribe_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _LiveConnection_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "connection.proto",
}
