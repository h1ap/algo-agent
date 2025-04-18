// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: oss/v1/oss.proto

package v1

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
	OSSService_TempUpload_FullMethodName       = "/api.oss.v1.OSSService/TempUpload"
	OSSService_UploadFile_FullMethodName       = "/api.oss.v1.OSSService/UploadFile"
	OSSService_DownloadFile_FullMethodName     = "/api.oss.v1.OSSService/DownloadFile"
	OSSService_DeleteFile_FullMethodName       = "/api.oss.v1.OSSService/DeleteFile"
	OSSService_BatchDeleteFiles_FullMethodName = "/api.oss.v1.OSSService/BatchDeleteFiles"
	OSSService_ListFiles_FullMethodName        = "/api.oss.v1.OSSService/ListFiles"
	OSSService_MergeFiles_FullMethodName       = "/api.oss.v1.OSSService/MergeFiles"
	OSSService_WriteString_FullMethodName      = "/api.oss.v1.OSSService/WriteString"
	OSSService_ReadString_FullMethodName       = "/api.oss.v1.OSSService/ReadString"
)

// OSSServiceClient is the client API for OSSService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// OSS服务接口定义
type OSSServiceClient interface {
	// 上传临时文件
	TempUpload(ctx context.Context, in *TempUploadRequest, opts ...grpc.CallOption) (*TempUploadReply, error)
	// 上传文件
	UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileReply, error)
	// 下载文件
	DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (*DownloadFileReply, error)
	// 删除文件
	DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileReply, error)
	// 批量删除文件
	BatchDeleteFiles(ctx context.Context, in *BatchDeleteFilesRequest, opts ...grpc.CallOption) (*BatchDeleteFilesReply, error)
	// 列出文件
	ListFiles(ctx context.Context, in *ListFilesRequest, opts ...grpc.CallOption) (*ListFilesReply, error)
	// 合并文件
	MergeFiles(ctx context.Context, in *MergeFilesRequest, opts ...grpc.CallOption) (*MergeFilesReply, error)
	// 写入字符串
	WriteString(ctx context.Context, in *WriteStringRequest, opts ...grpc.CallOption) (*WriteStringReply, error)
	// 读取字符串
	ReadString(ctx context.Context, in *ReadStringRequest, opts ...grpc.CallOption) (*ReadStringReply, error)
}

type oSSServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOSSServiceClient(cc grpc.ClientConnInterface) OSSServiceClient {
	return &oSSServiceClient{cc}
}

func (c *oSSServiceClient) TempUpload(ctx context.Context, in *TempUploadRequest, opts ...grpc.CallOption) (*TempUploadReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TempUploadReply)
	err := c.cc.Invoke(ctx, OSSService_TempUpload_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oSSServiceClient) UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadFileReply)
	err := c.cc.Invoke(ctx, OSSService_UploadFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oSSServiceClient) DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (*DownloadFileReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DownloadFileReply)
	err := c.cc.Invoke(ctx, OSSService_DownloadFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oSSServiceClient) DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteFileReply)
	err := c.cc.Invoke(ctx, OSSService_DeleteFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oSSServiceClient) BatchDeleteFiles(ctx context.Context, in *BatchDeleteFilesRequest, opts ...grpc.CallOption) (*BatchDeleteFilesReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchDeleteFilesReply)
	err := c.cc.Invoke(ctx, OSSService_BatchDeleteFiles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oSSServiceClient) ListFiles(ctx context.Context, in *ListFilesRequest, opts ...grpc.CallOption) (*ListFilesReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListFilesReply)
	err := c.cc.Invoke(ctx, OSSService_ListFiles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oSSServiceClient) MergeFiles(ctx context.Context, in *MergeFilesRequest, opts ...grpc.CallOption) (*MergeFilesReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MergeFilesReply)
	err := c.cc.Invoke(ctx, OSSService_MergeFiles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oSSServiceClient) WriteString(ctx context.Context, in *WriteStringRequest, opts ...grpc.CallOption) (*WriteStringReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WriteStringReply)
	err := c.cc.Invoke(ctx, OSSService_WriteString_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oSSServiceClient) ReadString(ctx context.Context, in *ReadStringRequest, opts ...grpc.CallOption) (*ReadStringReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadStringReply)
	err := c.cc.Invoke(ctx, OSSService_ReadString_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OSSServiceServer is the server API for OSSService service.
// All implementations must embed UnimplementedOSSServiceServer
// for forward compatibility.
//
// OSS服务接口定义
type OSSServiceServer interface {
	// 上传临时文件
	TempUpload(context.Context, *TempUploadRequest) (*TempUploadReply, error)
	// 上传文件
	UploadFile(context.Context, *UploadFileRequest) (*UploadFileReply, error)
	// 下载文件
	DownloadFile(context.Context, *DownloadFileRequest) (*DownloadFileReply, error)
	// 删除文件
	DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileReply, error)
	// 批量删除文件
	BatchDeleteFiles(context.Context, *BatchDeleteFilesRequest) (*BatchDeleteFilesReply, error)
	// 列出文件
	ListFiles(context.Context, *ListFilesRequest) (*ListFilesReply, error)
	// 合并文件
	MergeFiles(context.Context, *MergeFilesRequest) (*MergeFilesReply, error)
	// 写入字符串
	WriteString(context.Context, *WriteStringRequest) (*WriteStringReply, error)
	// 读取字符串
	ReadString(context.Context, *ReadStringRequest) (*ReadStringReply, error)
	mustEmbedUnimplementedOSSServiceServer()
}

// UnimplementedOSSServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedOSSServiceServer struct{}

func (UnimplementedOSSServiceServer) TempUpload(context.Context, *TempUploadRequest) (*TempUploadReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TempUpload not implemented")
}
func (UnimplementedOSSServiceServer) UploadFile(context.Context, *UploadFileRequest) (*UploadFileReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedOSSServiceServer) DownloadFile(context.Context, *DownloadFileRequest) (*DownloadFileReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedOSSServiceServer) DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFile not implemented")
}
func (UnimplementedOSSServiceServer) BatchDeleteFiles(context.Context, *BatchDeleteFilesRequest) (*BatchDeleteFilesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchDeleteFiles not implemented")
}
func (UnimplementedOSSServiceServer) ListFiles(context.Context, *ListFilesRequest) (*ListFilesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFiles not implemented")
}
func (UnimplementedOSSServiceServer) MergeFiles(context.Context, *MergeFilesRequest) (*MergeFilesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MergeFiles not implemented")
}
func (UnimplementedOSSServiceServer) WriteString(context.Context, *WriteStringRequest) (*WriteStringReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WriteString not implemented")
}
func (UnimplementedOSSServiceServer) ReadString(context.Context, *ReadStringRequest) (*ReadStringReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadString not implemented")
}
func (UnimplementedOSSServiceServer) mustEmbedUnimplementedOSSServiceServer() {}
func (UnimplementedOSSServiceServer) testEmbeddedByValue()                    {}

// UnsafeOSSServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OSSServiceServer will
// result in compilation errors.
type UnsafeOSSServiceServer interface {
	mustEmbedUnimplementedOSSServiceServer()
}

func RegisterOSSServiceServer(s grpc.ServiceRegistrar, srv OSSServiceServer) {
	// If the following call pancis, it indicates UnimplementedOSSServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&OSSService_ServiceDesc, srv)
}

func _OSSService_TempUpload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TempUploadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OSSServiceServer).TempUpload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OSSService_TempUpload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OSSServiceServer).TempUpload(ctx, req.(*TempUploadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OSSService_UploadFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OSSServiceServer).UploadFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OSSService_UploadFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OSSServiceServer).UploadFile(ctx, req.(*UploadFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OSSService_DownloadFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OSSServiceServer).DownloadFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OSSService_DownloadFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OSSServiceServer).DownloadFile(ctx, req.(*DownloadFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OSSService_DeleteFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OSSServiceServer).DeleteFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OSSService_DeleteFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OSSServiceServer).DeleteFile(ctx, req.(*DeleteFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OSSService_BatchDeleteFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchDeleteFilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OSSServiceServer).BatchDeleteFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OSSService_BatchDeleteFiles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OSSServiceServer).BatchDeleteFiles(ctx, req.(*BatchDeleteFilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OSSService_ListFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListFilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OSSServiceServer).ListFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OSSService_ListFiles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OSSServiceServer).ListFiles(ctx, req.(*ListFilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OSSService_MergeFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MergeFilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OSSServiceServer).MergeFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OSSService_MergeFiles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OSSServiceServer).MergeFiles(ctx, req.(*MergeFilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OSSService_WriteString_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteStringRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OSSServiceServer).WriteString(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OSSService_WriteString_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OSSServiceServer).WriteString(ctx, req.(*WriteStringRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OSSService_ReadString_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadStringRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OSSServiceServer).ReadString(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OSSService_ReadString_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OSSServiceServer).ReadString(ctx, req.(*ReadStringRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OSSService_ServiceDesc is the grpc.ServiceDesc for OSSService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OSSService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.oss.v1.OSSService",
	HandlerType: (*OSSServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TempUpload",
			Handler:    _OSSService_TempUpload_Handler,
		},
		{
			MethodName: "UploadFile",
			Handler:    _OSSService_UploadFile_Handler,
		},
		{
			MethodName: "DownloadFile",
			Handler:    _OSSService_DownloadFile_Handler,
		},
		{
			MethodName: "DeleteFile",
			Handler:    _OSSService_DeleteFile_Handler,
		},
		{
			MethodName: "BatchDeleteFiles",
			Handler:    _OSSService_BatchDeleteFiles_Handler,
		},
		{
			MethodName: "ListFiles",
			Handler:    _OSSService_ListFiles_Handler,
		},
		{
			MethodName: "MergeFiles",
			Handler:    _OSSService_MergeFiles_Handler,
		},
		{
			MethodName: "WriteString",
			Handler:    _OSSService_WriteString_Handler,
		},
		{
			MethodName: "ReadString",
			Handler:    _OSSService_ReadString_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "oss/v1/oss.proto",
}
