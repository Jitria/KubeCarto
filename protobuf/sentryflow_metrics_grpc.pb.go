// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.26.0
// source: sentryflow_metrics.proto

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
	APIClassifier_ClassifyAPIs_FullMethodName = "/protobuf.APIClassifier/ClassifyAPIs"
)

// APIClassifierClient is the client API for APIClassifier service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type APIClassifierClient interface {
	ClassifyAPIs(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[APIClassifierRequest, APIClassifierResponse], error)
}

type aPIClassifierClient struct {
	cc grpc.ClientConnInterface
}

func NewAPIClassifierClient(cc grpc.ClientConnInterface) APIClassifierClient {
	return &aPIClassifierClient{cc}
}

func (c *aPIClassifierClient) ClassifyAPIs(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[APIClassifierRequest, APIClassifierResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &APIClassifier_ServiceDesc.Streams[0], APIClassifier_ClassifyAPIs_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[APIClassifierRequest, APIClassifierResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type APIClassifier_ClassifyAPIsClient = grpc.BidiStreamingClient[APIClassifierRequest, APIClassifierResponse]

// APIClassifierServer is the server API for APIClassifier service.
// All implementations should embed UnimplementedAPIClassifierServer
// for forward compatibility.
type APIClassifierServer interface {
	ClassifyAPIs(grpc.BidiStreamingServer[APIClassifierRequest, APIClassifierResponse]) error
}

// UnimplementedAPIClassifierServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAPIClassifierServer struct{}

func (UnimplementedAPIClassifierServer) ClassifyAPIs(grpc.BidiStreamingServer[APIClassifierRequest, APIClassifierResponse]) error {
	return status.Errorf(codes.Unimplemented, "method ClassifyAPIs not implemented")
}
func (UnimplementedAPIClassifierServer) testEmbeddedByValue() {}

// UnsafeAPIClassifierServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to APIClassifierServer will
// result in compilation errors.
type UnsafeAPIClassifierServer interface {
	mustEmbedUnimplementedAPIClassifierServer()
}

func RegisterAPIClassifierServer(s grpc.ServiceRegistrar, srv APIClassifierServer) {
	// If the following call pancis, it indicates UnimplementedAPIClassifierServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&APIClassifier_ServiceDesc, srv)
}

func _APIClassifier_ClassifyAPIs_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(APIClassifierServer).ClassifyAPIs(&grpc.GenericServerStream[APIClassifierRequest, APIClassifierResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type APIClassifier_ClassifyAPIsServer = grpc.BidiStreamingServer[APIClassifierRequest, APIClassifierResponse]

// APIClassifier_ServiceDesc is the grpc.ServiceDesc for APIClassifier service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var APIClassifier_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.APIClassifier",
	HandlerType: (*APIClassifierServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ClassifyAPIs",
			Handler:       _APIClassifier_ClassifyAPIs_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "sentryflow_metrics.proto",
}
