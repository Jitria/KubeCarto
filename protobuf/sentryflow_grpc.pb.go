// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.26.0
// source: sentryflow.proto

package protobuf

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SentryFlow_GetAPILog_FullMethodName        = "/protobuf.SentryFlow/GetAPILog"
	SentryFlow_GetEnvoyMetrics_FullMethodName  = "/protobuf.SentryFlow/GetEnvoyMetrics"
	SentryFlow_GiveAPILog_FullMethodName       = "/protobuf.SentryFlow/GiveAPILog"
	SentryFlow_GiveEnvoyMetrics_FullMethodName = "/protobuf.SentryFlow/GiveEnvoyMetrics"
	SentryFlow_GiveDeployInfo_FullMethodName   = "/protobuf.SentryFlow/GiveDeployInfo"
	SentryFlow_GivePodInfo_FullMethodName      = "/protobuf.SentryFlow/GivePodInfo"
	SentryFlow_GiveSvcInfo_FullMethodName      = "/protobuf.SentryFlow/GiveSvcInfo"
)

// SentryFlowClient is the client API for SentryFlow service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SentryFlowClient interface {
	// enovy -> agent
	GetAPILog(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (grpc.ServerStreamingClient[APILog], error)
	GetEnvoyMetrics(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (grpc.ServerStreamingClient[EnvoyMetrics], error)
	// agent -> operator
	GiveAPILog(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[APILog, Response], error)
	GiveEnvoyMetrics(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[EnvoyMetrics, Response], error)
	GiveDeployInfo(ctx context.Context, in *Deploy, opts ...grpc.CallOption) (*Response, error)
	GivePodInfo(ctx context.Context, in *Pod, opts ...grpc.CallOption) (*Response, error)
	GiveSvcInfo(ctx context.Context, in *Service, opts ...grpc.CallOption) (*Response, error)
}

type sentryFlowClient struct {
	cc grpc.ClientConnInterface
}

func NewSentryFlowClient(cc grpc.ClientConnInterface) SentryFlowClient {
	return &sentryFlowClient{cc}
}

func (c *sentryFlowClient) GetAPILog(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (grpc.ServerStreamingClient[APILog], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &SentryFlow_ServiceDesc.Streams[0], SentryFlow_GetAPILog_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ClientInfo, APILog]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SentryFlow_GetAPILogClient = grpc.ServerStreamingClient[APILog]

func (c *sentryFlowClient) GetEnvoyMetrics(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (grpc.ServerStreamingClient[EnvoyMetrics], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &SentryFlow_ServiceDesc.Streams[1], SentryFlow_GetEnvoyMetrics_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ClientInfo, EnvoyMetrics]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SentryFlow_GetEnvoyMetricsClient = grpc.ServerStreamingClient[EnvoyMetrics]

func (c *sentryFlowClient) GiveAPILog(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[APILog, Response], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &SentryFlow_ServiceDesc.Streams[2], SentryFlow_GiveAPILog_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[APILog, Response]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SentryFlow_GiveAPILogClient = grpc.ClientStreamingClient[APILog, Response]

func (c *sentryFlowClient) GiveEnvoyMetrics(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[EnvoyMetrics, Response], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &SentryFlow_ServiceDesc.Streams[3], SentryFlow_GiveEnvoyMetrics_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[EnvoyMetrics, Response]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SentryFlow_GiveEnvoyMetricsClient = grpc.ClientStreamingClient[EnvoyMetrics, Response]

func (c *sentryFlowClient) GiveDeployInfo(ctx context.Context, in *Deploy, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, SentryFlow_GiveDeployInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sentryFlowClient) GivePodInfo(ctx context.Context, in *Pod, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, SentryFlow_GivePodInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sentryFlowClient) GiveSvcInfo(ctx context.Context, in *Service, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, SentryFlow_GiveSvcInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SentryFlowServer is the server API for SentryFlow service.
// All implementations should embed UnimplementedSentryFlowServer
// for forward compatibility.
type SentryFlowServer interface {
	// enovy -> agent
	GetAPILog(*ClientInfo, grpc.ServerStreamingServer[APILog]) error
	GetEnvoyMetrics(*ClientInfo, grpc.ServerStreamingServer[EnvoyMetrics]) error
	// agent -> operator
	GiveAPILog(grpc.ClientStreamingServer[APILog, Response]) error
	GiveEnvoyMetrics(grpc.ClientStreamingServer[EnvoyMetrics, Response]) error
	GiveDeployInfo(context.Context, *Deploy) (*Response, error)
	GivePodInfo(context.Context, *Pod) (*Response, error)
	GiveSvcInfo(context.Context, *Service) (*Response, error)
}

// UnimplementedSentryFlowServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSentryFlowServer struct{}

func (UnimplementedSentryFlowServer) GetAPILog(*ClientInfo, grpc.ServerStreamingServer[APILog]) error {
	return status.Errorf(codes.Unimplemented, "method GetAPILog not implemented")
}
func (UnimplementedSentryFlowServer) GetEnvoyMetrics(*ClientInfo, grpc.ServerStreamingServer[EnvoyMetrics]) error {
	return status.Errorf(codes.Unimplemented, "method GetEnvoyMetrics not implemented")
}
func (UnimplementedSentryFlowServer) GiveAPILog(grpc.ClientStreamingServer[APILog, Response]) error {
	return status.Errorf(codes.Unimplemented, "method GiveAPILog not implemented")
}
func (UnimplementedSentryFlowServer) GiveEnvoyMetrics(grpc.ClientStreamingServer[EnvoyMetrics, Response]) error {
	return status.Errorf(codes.Unimplemented, "method GiveEnvoyMetrics not implemented")
}
func (UnimplementedSentryFlowServer) GiveDeployInfo(context.Context, *Deploy) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GiveDeployInfo not implemented")
}
func (UnimplementedSentryFlowServer) GivePodInfo(context.Context, *Pod) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GivePodInfo not implemented")
}
func (UnimplementedSentryFlowServer) GiveSvcInfo(context.Context, *Service) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GiveSvcInfo not implemented")
}
func (UnimplementedSentryFlowServer) testEmbeddedByValue() {}

// UnsafeSentryFlowServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SentryFlowServer will
// result in compilation errors.
type UnsafeSentryFlowServer interface {
	mustEmbedUnimplementedSentryFlowServer()
}

func RegisterSentryFlowServer(s grpc.ServiceRegistrar, srv SentryFlowServer) {
	// If the following call pancis, it indicates UnimplementedSentryFlowServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SentryFlow_ServiceDesc, srv)
}

func _SentryFlow_GetAPILog_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ClientInfo)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SentryFlowServer).GetAPILog(m, &grpc.GenericServerStream[ClientInfo, APILog]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SentryFlow_GetAPILogServer = grpc.ServerStreamingServer[APILog]

func _SentryFlow_GetEnvoyMetrics_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ClientInfo)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SentryFlowServer).GetEnvoyMetrics(m, &grpc.GenericServerStream[ClientInfo, EnvoyMetrics]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SentryFlow_GetEnvoyMetricsServer = grpc.ServerStreamingServer[EnvoyMetrics]

func _SentryFlow_GiveAPILog_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SentryFlowServer).GiveAPILog(&grpc.GenericServerStream[APILog, Response]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SentryFlow_GiveAPILogServer = grpc.ClientStreamingServer[APILog, Response]

func _SentryFlow_GiveEnvoyMetrics_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SentryFlowServer).GiveEnvoyMetrics(&grpc.GenericServerStream[EnvoyMetrics, Response]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SentryFlow_GiveEnvoyMetricsServer = grpc.ClientStreamingServer[EnvoyMetrics, Response]

func _SentryFlow_GiveDeployInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Deploy)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentryFlowServer).GiveDeployInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SentryFlow_GiveDeployInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentryFlowServer).GiveDeployInfo(ctx, req.(*Deploy))
	}
	return interceptor(ctx, in, info, handler)
}

func _SentryFlow_GivePodInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Pod)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentryFlowServer).GivePodInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SentryFlow_GivePodInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentryFlowServer).GivePodInfo(ctx, req.(*Pod))
	}
	return interceptor(ctx, in, info, handler)
}

func _SentryFlow_GiveSvcInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Service)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentryFlowServer).GiveSvcInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SentryFlow_GiveSvcInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentryFlowServer).GiveSvcInfo(ctx, req.(*Service))
	}
	return interceptor(ctx, in, info, handler)
}

// SentryFlow_ServiceDesc is the grpc.ServiceDesc for SentryFlow service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SentryFlow_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.SentryFlow",
	HandlerType: (*SentryFlowServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GiveDeployInfo",
			Handler:    _SentryFlow_GiveDeployInfo_Handler,
		},
		{
			MethodName: "GivePodInfo",
			Handler:    _SentryFlow_GivePodInfo_Handler,
		},
		{
			MethodName: "GiveSvcInfo",
			Handler:    _SentryFlow_GiveSvcInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAPILog",
			Handler:       _SentryFlow_GetAPILog_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetEnvoyMetrics",
			Handler:       _SentryFlow_GetEnvoyMetrics_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GiveAPILog",
			Handler:       _SentryFlow_GiveAPILog_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GiveEnvoyMetrics",
			Handler:       _SentryFlow_GiveEnvoyMetrics_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "sentryflow.proto",
}
