// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: users.proto

package users

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// HealthClient is the client API for Health service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HealthClient interface {
	Check(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
	Watch(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (Health_WatchClient, error)
}

type healthClient struct {
	cc grpc.ClientConnInterface
}

func NewHealthClient(cc grpc.ClientConnInterface) HealthClient {
	return &healthClient{cc}
}

func (c *healthClient) Check(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, "/users.Health/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *healthClient) Watch(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (Health_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &Health_ServiceDesc.Streams[0], "/users.Health/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &healthWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Health_WatchClient interface {
	Recv() (*HealthCheckResponse, error)
	grpc.ClientStream
}

type healthWatchClient struct {
	grpc.ClientStream
}

func (x *healthWatchClient) Recv() (*HealthCheckResponse, error) {
	m := new(HealthCheckResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// HealthServer is the server API for Health service.
// All implementations should embed UnimplementedHealthServer
// for forward compatibility
type HealthServer interface {
	Check(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
	Watch(*emptypb.Empty, Health_WatchServer) error
}

// UnimplementedHealthServer should be embedded to have forward compatible implementations.
type UnimplementedHealthServer struct {
}

func (UnimplementedHealthServer) Check(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedHealthServer) Watch(*emptypb.Empty, Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}

// UnsafeHealthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HealthServer will
// result in compilation errors.
type UnsafeHealthServer interface {
	mustEmbedUnimplementedHealthServer()
}

func RegisterHealthServer(s grpc.ServiceRegistrar, srv HealthServer) {
	s.RegisterService(&Health_ServiceDesc, srv)
}

func _Health_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.Health/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthServer).Check(ctx, req.(*HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Health_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(HealthServer).Watch(m, &healthWatchServer{stream})
}

type Health_WatchServer interface {
	Send(*HealthCheckResponse) error
	grpc.ServerStream
}

type healthWatchServer struct {
	grpc.ServerStream
}

func (x *healthWatchServer) Send(m *HealthCheckResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Health_ServiceDesc is the grpc.ServiceDesc for Health service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Health_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "users.Health",
	HandlerType: (*HealthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Check",
			Handler:    _Health_Check_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _Health_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "users.proto",
}

// UserGroupServiceClient is the client API for UserGroupService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserGroupServiceClient interface {
	CreateGroup(ctx context.Context, in *CreateGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpdateGroup(ctx context.Context, in *UpdateGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteGroup(ctx context.Context, in *DeleteGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeactivateGroup(ctx context.Context, in *DeactivateGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ActivateGroup(ctx context.Context, in *ActivateGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	IsGroupAvailable(ctx context.Context, in *IsGroupAvailableRequest, opts ...grpc.CallOption) (*IsGroupAvailableResponse, error)
	GetGroup(ctx context.Context, in *GetGroupRequest, opts ...grpc.CallOption) (*GetGroupResponse, error)
}

type userGroupServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserGroupServiceClient(cc grpc.ClientConnInterface) UserGroupServiceClient {
	return &userGroupServiceClient{cc}
}

func (c *userGroupServiceClient) CreateGroup(ctx context.Context, in *CreateGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/users.UserGroupService/CreateGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userGroupServiceClient) UpdateGroup(ctx context.Context, in *UpdateGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/users.UserGroupService/UpdateGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userGroupServiceClient) DeleteGroup(ctx context.Context, in *DeleteGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/users.UserGroupService/DeleteGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userGroupServiceClient) DeactivateGroup(ctx context.Context, in *DeactivateGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/users.UserGroupService/DeactivateGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userGroupServiceClient) ActivateGroup(ctx context.Context, in *ActivateGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/users.UserGroupService/ActivateGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userGroupServiceClient) IsGroupAvailable(ctx context.Context, in *IsGroupAvailableRequest, opts ...grpc.CallOption) (*IsGroupAvailableResponse, error) {
	out := new(IsGroupAvailableResponse)
	err := c.cc.Invoke(ctx, "/users.UserGroupService/IsGroupAvailable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userGroupServiceClient) GetGroup(ctx context.Context, in *GetGroupRequest, opts ...grpc.CallOption) (*GetGroupResponse, error) {
	out := new(GetGroupResponse)
	err := c.cc.Invoke(ctx, "/users.UserGroupService/GetGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserGroupServiceServer is the server API for UserGroupService service.
// All implementations should embed UnimplementedUserGroupServiceServer
// for forward compatibility
type UserGroupServiceServer interface {
	CreateGroup(context.Context, *CreateGroupRequest) (*emptypb.Empty, error)
	UpdateGroup(context.Context, *UpdateGroupRequest) (*emptypb.Empty, error)
	DeleteGroup(context.Context, *DeleteGroupRequest) (*emptypb.Empty, error)
	DeactivateGroup(context.Context, *DeactivateGroupRequest) (*emptypb.Empty, error)
	ActivateGroup(context.Context, *ActivateGroupRequest) (*emptypb.Empty, error)
	IsGroupAvailable(context.Context, *IsGroupAvailableRequest) (*IsGroupAvailableResponse, error)
	GetGroup(context.Context, *GetGroupRequest) (*GetGroupResponse, error)
}

// UnimplementedUserGroupServiceServer should be embedded to have forward compatible implementations.
type UnimplementedUserGroupServiceServer struct {
}

func (UnimplementedUserGroupServiceServer) CreateGroup(context.Context, *CreateGroupRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGroup not implemented")
}
func (UnimplementedUserGroupServiceServer) UpdateGroup(context.Context, *UpdateGroupRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateGroup not implemented")
}
func (UnimplementedUserGroupServiceServer) DeleteGroup(context.Context, *DeleteGroupRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGroup not implemented")
}
func (UnimplementedUserGroupServiceServer) DeactivateGroup(context.Context, *DeactivateGroupRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeactivateGroup not implemented")
}
func (UnimplementedUserGroupServiceServer) ActivateGroup(context.Context, *ActivateGroupRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivateGroup not implemented")
}
func (UnimplementedUserGroupServiceServer) IsGroupAvailable(context.Context, *IsGroupAvailableRequest) (*IsGroupAvailableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsGroupAvailable not implemented")
}
func (UnimplementedUserGroupServiceServer) GetGroup(context.Context, *GetGroupRequest) (*GetGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroup not implemented")
}

// UnsafeUserGroupServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserGroupServiceServer will
// result in compilation errors.
type UnsafeUserGroupServiceServer interface {
	mustEmbedUnimplementedUserGroupServiceServer()
}

func RegisterUserGroupServiceServer(s grpc.ServiceRegistrar, srv UserGroupServiceServer) {
	s.RegisterService(&UserGroupService_ServiceDesc, srv)
}

func _UserGroupService_CreateGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserGroupServiceServer).CreateGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UserGroupService/CreateGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserGroupServiceServer).CreateGroup(ctx, req.(*CreateGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserGroupService_UpdateGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserGroupServiceServer).UpdateGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UserGroupService/UpdateGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserGroupServiceServer).UpdateGroup(ctx, req.(*UpdateGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserGroupService_DeleteGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserGroupServiceServer).DeleteGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UserGroupService/DeleteGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserGroupServiceServer).DeleteGroup(ctx, req.(*DeleteGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserGroupService_DeactivateGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeactivateGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserGroupServiceServer).DeactivateGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UserGroupService/DeactivateGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserGroupServiceServer).DeactivateGroup(ctx, req.(*DeactivateGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserGroupService_ActivateGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivateGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserGroupServiceServer).ActivateGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UserGroupService/ActivateGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserGroupServiceServer).ActivateGroup(ctx, req.(*ActivateGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserGroupService_IsGroupAvailable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsGroupAvailableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserGroupServiceServer).IsGroupAvailable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UserGroupService/IsGroupAvailable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserGroupServiceServer).IsGroupAvailable(ctx, req.(*IsGroupAvailableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserGroupService_GetGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserGroupServiceServer).GetGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UserGroupService/GetGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserGroupServiceServer).GetGroup(ctx, req.(*GetGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserGroupService_ServiceDesc is the grpc.ServiceDesc for UserGroupService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserGroupService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "users.UserGroupService",
	HandlerType: (*UserGroupServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateGroup",
			Handler:    _UserGroupService_CreateGroup_Handler,
		},
		{
			MethodName: "UpdateGroup",
			Handler:    _UserGroupService_UpdateGroup_Handler,
		},
		{
			MethodName: "DeleteGroup",
			Handler:    _UserGroupService_DeleteGroup_Handler,
		},
		{
			MethodName: "DeactivateGroup",
			Handler:    _UserGroupService_DeactivateGroup_Handler,
		},
		{
			MethodName: "ActivateGroup",
			Handler:    _UserGroupService_ActivateGroup_Handler,
		},
		{
			MethodName: "IsGroupAvailable",
			Handler:    _UserGroupService_IsGroupAvailable_Handler,
		},
		{
			MethodName: "GetGroup",
			Handler:    _UserGroupService_GetGroup_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "users.proto",
}
