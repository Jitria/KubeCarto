// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.26.0
// source: agent.proto

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
	Agent_GetAPILog_FullMethodName        = "/protobuf.Agent/GetAPILog"
	Agent_GetEnvoyMetrics_FullMethodName  = "/protobuf.Agent/GetEnvoyMetrics"
	Agent_GiveAPILog_FullMethodName       = "/protobuf.Agent/GiveAPILog"
	Agent_GiveEnvoyMetrics_FullMethodName = "/protobuf.Agent/GiveEnvoyMetrics"
)

// AgentClient is the client API for Agent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentClient interface {
	GetAPILog(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (grpc.ServerStreamingClient[APILog], error)
	GetEnvoyMetrics(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (grpc.ServerStreamingClient[EnvoyMetrics], error)
	GiveAPILog(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[APILog, Response], error)
	GiveEnvoyMetrics(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[EnvoyMetrics, Response], error)
}

type agentClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentClient(cc grpc.ClientConnInterface) AgentClient {
	return &agentClient{cc}
}

func (c *agentClient) GetAPILog(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (grpc.ServerStreamingClient[APILog], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Agent_ServiceDesc.Streams[0], Agent_GetAPILog_FullMethodName, cOpts...)
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
type Agent_GetAPILogClient = grpc.ServerStreamingClient[APILog]

func (c *agentClient) GetEnvoyMetrics(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (grpc.ServerStreamingClient[EnvoyMetrics], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Agent_ServiceDesc.Streams[1], Agent_GetEnvoyMetrics_FullMethodName, cOpts...)
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
type Agent_GetEnvoyMetricsClient = grpc.ServerStreamingClient[EnvoyMetrics]

func (c *agentClient) GiveAPILog(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[APILog, Response], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Agent_ServiceDesc.Streams[2], Agent_GiveAPILog_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[APILog, Response]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Agent_GiveAPILogClient = grpc.ClientStreamingClient[APILog, Response]

func (c *agentClient) GiveEnvoyMetrics(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[EnvoyMetrics, Response], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Agent_ServiceDesc.Streams[3], Agent_GiveEnvoyMetrics_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[EnvoyMetrics, Response]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Agent_GiveEnvoyMetricsClient = grpc.ClientStreamingClient[EnvoyMetrics, Response]

// AgentServer is the server API for Agent service.
// All implementations should embed UnimplementedAgentServer
// for forward compatibility.
type AgentServer interface {
	GetAPILog(*ClientInfo, grpc.ServerStreamingServer[APILog]) error
	GetEnvoyMetrics(*ClientInfo, grpc.ServerStreamingServer[EnvoyMetrics]) error
	GiveAPILog(grpc.ClientStreamingServer[APILog, Response]) error
	GiveEnvoyMetrics(grpc.ClientStreamingServer[EnvoyMetrics, Response]) error
}

// UnimplementedAgentServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAgentServer struct{}

func (UnimplementedAgentServer) GetAPILog(*ClientInfo, grpc.ServerStreamingServer[APILog]) error {
	return status.Errorf(codes.Unimplemented, "method GetAPILog not implemented")
}
func (UnimplementedAgentServer) GetEnvoyMetrics(*ClientInfo, grpc.ServerStreamingServer[EnvoyMetrics]) error {
	return status.Errorf(codes.Unimplemented, "method GetEnvoyMetrics not implemented")
}
func (UnimplementedAgentServer) GiveAPILog(grpc.ClientStreamingServer[APILog, Response]) error {
	return status.Errorf(codes.Unimplemented, "method GiveAPILog not implemented")
}
func (UnimplementedAgentServer) GiveEnvoyMetrics(grpc.ClientStreamingServer[EnvoyMetrics, Response]) error {
	return status.Errorf(codes.Unimplemented, "method GiveEnvoyMetrics not implemented")
}
func (UnimplementedAgentServer) testEmbeddedByValue() {}

// UnsafeAgentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentServer will
// result in compilation errors.
type UnsafeAgentServer interface {
	mustEmbedUnimplementedAgentServer()
}

func RegisterAgentServer(s grpc.ServiceRegistrar, srv AgentServer) {
	// If the following call pancis, it indicates UnimplementedAgentServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Agent_ServiceDesc, srv)
}

func _Agent_GetAPILog_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ClientInfo)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AgentServer).GetAPILog(m, &grpc.GenericServerStream[ClientInfo, APILog]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Agent_GetAPILogServer = grpc.ServerStreamingServer[APILog]

func _Agent_GetEnvoyMetrics_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ClientInfo)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AgentServer).GetEnvoyMetrics(m, &grpc.GenericServerStream[ClientInfo, EnvoyMetrics]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Agent_GetEnvoyMetricsServer = grpc.ServerStreamingServer[EnvoyMetrics]

func _Agent_GiveAPILog_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AgentServer).GiveAPILog(&grpc.GenericServerStream[APILog, Response]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Agent_GiveAPILogServer = grpc.ClientStreamingServer[APILog, Response]

func _Agent_GiveEnvoyMetrics_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AgentServer).GiveEnvoyMetrics(&grpc.GenericServerStream[EnvoyMetrics, Response]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Agent_GiveEnvoyMetricsServer = grpc.ClientStreamingServer[EnvoyMetrics, Response]

// Agent_ServiceDesc is the grpc.ServiceDesc for Agent service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Agent_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.Agent",
	HandlerType: (*AgentServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAPILog",
			Handler:       _Agent_GetAPILog_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetEnvoyMetrics",
			Handler:       _Agent_GetEnvoyMetrics_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GiveAPILog",
			Handler:       _Agent_GiveAPILog_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GiveEnvoyMetrics",
			Handler:       _Agent_GiveEnvoyMetrics_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "agent.proto",
}
