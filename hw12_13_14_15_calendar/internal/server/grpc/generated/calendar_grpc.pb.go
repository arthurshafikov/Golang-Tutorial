// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package generated

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

// CalendarServiceClient is the client API for CalendarService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CalendarServiceClient interface {
	Create(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*ServerResponse, error)
	Update(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*ServerResponse, error)
	Delete(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*ServerResponse, error)
	ListEventsOnADay(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ServerResponse, error)
	ListEventsOnAWeek(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ServerResponse, error)
	ListEventsOnAMonth(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ServerResponse, error)
}

type calendarServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCalendarServiceClient(cc grpc.ClientConnInterface) CalendarServiceClient {
	return &calendarServiceClient{cc}
}

func (c *calendarServiceClient) Create(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*ServerResponse, error) {
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, "/api.CalendarService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) Update(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*ServerResponse, error) {
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, "/api.CalendarService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) Delete(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*ServerResponse, error) {
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, "/api.CalendarService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) ListEventsOnADay(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ServerResponse, error) {
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, "/api.CalendarService/ListEventsOnADay", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) ListEventsOnAWeek(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ServerResponse, error) {
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, "/api.CalendarService/ListEventsOnAWeek", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) ListEventsOnAMonth(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ServerResponse, error) {
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, "/api.CalendarService/ListEventsOnAMonth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarServiceServer is the server API for CalendarService service.
// All implementations must embed UnimplementedCalendarServiceServer
// for forward compatibility
type CalendarServiceServer interface {
	Create(context.Context, *EventRequest) (*ServerResponse, error)
	Update(context.Context, *EventRequest) (*ServerResponse, error)
	Delete(context.Context, *EventRequest) (*ServerResponse, error)
	ListEventsOnADay(context.Context, *ListEventsRequest) (*ServerResponse, error)
	ListEventsOnAWeek(context.Context, *ListEventsRequest) (*ServerResponse, error)
	ListEventsOnAMonth(context.Context, *ListEventsRequest) (*ServerResponse, error)
	mustEmbedUnimplementedCalendarServiceServer()
}

// UnimplementedCalendarServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCalendarServiceServer struct {
}

func (UnimplementedCalendarServiceServer) Create(context.Context, *EventRequest) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedCalendarServiceServer) Update(context.Context, *EventRequest) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedCalendarServiceServer) Delete(context.Context, *EventRequest) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedCalendarServiceServer) ListEventsOnADay(context.Context, *ListEventsRequest) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventsOnADay not implemented")
}
func (UnimplementedCalendarServiceServer) ListEventsOnAWeek(context.Context, *ListEventsRequest) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventsOnAWeek not implemented")
}
func (UnimplementedCalendarServiceServer) ListEventsOnAMonth(context.Context, *ListEventsRequest) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventsOnAMonth not implemented")
}
func (UnimplementedCalendarServiceServer) mustEmbedUnimplementedCalendarServiceServer() {}

// UnsafeCalendarServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalendarServiceServer will
// result in compilation errors.
type UnsafeCalendarServiceServer interface {
	mustEmbedUnimplementedCalendarServiceServer()
}

func RegisterCalendarServiceServer(s grpc.ServiceRegistrar, srv CalendarServiceServer) {
	s.RegisterService(&CalendarService_ServiceDesc, srv)
}

func _CalendarService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CalendarService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).Create(ctx, req.(*EventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CalendarService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).Update(ctx, req.(*EventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CalendarService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).Delete(ctx, req.(*EventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_ListEventsOnADay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).ListEventsOnADay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CalendarService/ListEventsOnADay",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).ListEventsOnADay(ctx, req.(*ListEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_ListEventsOnAWeek_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).ListEventsOnAWeek(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CalendarService/ListEventsOnAWeek",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).ListEventsOnAWeek(ctx, req.(*ListEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_ListEventsOnAMonth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).ListEventsOnAMonth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.CalendarService/ListEventsOnAMonth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).ListEventsOnAMonth(ctx, req.(*ListEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CalendarService_ServiceDesc is the grpc.ServiceDesc for CalendarService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CalendarService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.CalendarService",
	HandlerType: (*CalendarServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _CalendarService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _CalendarService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _CalendarService_Delete_Handler,
		},
		{
			MethodName: "ListEventsOnADay",
			Handler:    _CalendarService_ListEventsOnADay_Handler,
		},
		{
			MethodName: "ListEventsOnAWeek",
			Handler:    _CalendarService_ListEventsOnAWeek_Handler,
		},
		{
			MethodName: "ListEventsOnAMonth",
			Handler:    _CalendarService_ListEventsOnAMonth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "calendar.proto",
}
