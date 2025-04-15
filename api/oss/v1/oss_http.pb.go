// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             v5.29.3
// source: oss/v1/oss.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationOSSServiceBatchDeleteFiles = "/api.oss.v1.OSSService/BatchDeleteFiles"
const OperationOSSServiceDeleteFile = "/api.oss.v1.OSSService/DeleteFile"
const OperationOSSServiceDownloadFile = "/api.oss.v1.OSSService/DownloadFile"
const OperationOSSServiceListFiles = "/api.oss.v1.OSSService/ListFiles"
const OperationOSSServiceMergeFiles = "/api.oss.v1.OSSService/MergeFiles"
const OperationOSSServiceReadString = "/api.oss.v1.OSSService/ReadString"
const OperationOSSServiceTempUpload = "/api.oss.v1.OSSService/TempUpload"
const OperationOSSServiceUploadFile = "/api.oss.v1.OSSService/UploadFile"
const OperationOSSServiceWriteString = "/api.oss.v1.OSSService/WriteString"

type OSSServiceHTTPServer interface {
	// BatchDeleteFiles 批量删除文件
	BatchDeleteFiles(context.Context, *BatchDeleteFilesRequest) (*BatchDeleteFilesReply, error)
	// DeleteFile 删除文件
	DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileReply, error)
	// DownloadFile 下载文件
	DownloadFile(context.Context, *DownloadFileRequest) (*DownloadFileReply, error)
	// ListFiles 列出文件
	ListFiles(context.Context, *ListFilesRequest) (*ListFilesReply, error)
	// MergeFiles 合并文件
	MergeFiles(context.Context, *MergeFilesRequest) (*MergeFilesReply, error)
	// ReadString 读取字符串
	ReadString(context.Context, *ReadStringRequest) (*ReadStringReply, error)
	// TempUpload 上传临时文件
	TempUpload(context.Context, *TempUploadRequest) (*TempUploadReply, error)
	// UploadFile 上传文件
	UploadFile(context.Context, *UploadFileRequest) (*UploadFileReply, error)
	// WriteString 写入字符串
	WriteString(context.Context, *WriteStringRequest) (*WriteStringReply, error)
}

func RegisterOSSServiceHTTPServer(s *http.Server, srv OSSServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/oss/temp/upload", _OSSService_TempUpload0_HTTP_Handler(srv))
	r.POST("/v1/oss/upload", _OSSService_UploadFile0_HTTP_Handler(srv))
	r.GET("/v1/oss/download/{bucket_name}/{file_path}", _OSSService_DownloadFile0_HTTP_Handler(srv))
	r.DELETE("/v1/oss/file/{bucket_name}/{file_path}", _OSSService_DeleteFile0_HTTP_Handler(srv))
	r.POST("/v1/oss/batch/delete", _OSSService_BatchDeleteFiles0_HTTP_Handler(srv))
	r.GET("/v1/oss/list/{bucket_name}/{prefix}", _OSSService_ListFiles0_HTTP_Handler(srv))
	r.POST("/v1/oss/merge", _OSSService_MergeFiles0_HTTP_Handler(srv))
	r.POST("/v1/oss/write", _OSSService_WriteString0_HTTP_Handler(srv))
	r.GET("/v1/oss/read/{bucket_name}/{file_path}", _OSSService_ReadString0_HTTP_Handler(srv))
}

func _OSSService_TempUpload0_HTTP_Handler(srv OSSServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in TempUploadRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationOSSServiceTempUpload)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.TempUpload(ctx, req.(*TempUploadRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*TempUploadReply)
		return ctx.Result(200, reply)
	}
}

func _OSSService_UploadFile0_HTTP_Handler(srv OSSServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UploadFileRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationOSSServiceUploadFile)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UploadFile(ctx, req.(*UploadFileRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UploadFileReply)
		return ctx.Result(200, reply)
	}
}

func _OSSService_DownloadFile0_HTTP_Handler(srv OSSServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DownloadFileRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationOSSServiceDownloadFile)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DownloadFile(ctx, req.(*DownloadFileRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DownloadFileReply)
		return ctx.Result(200, reply)
	}
}

func _OSSService_DeleteFile0_HTTP_Handler(srv OSSServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteFileRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationOSSServiceDeleteFile)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteFile(ctx, req.(*DeleteFileRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteFileReply)
		return ctx.Result(200, reply)
	}
}

func _OSSService_BatchDeleteFiles0_HTTP_Handler(srv OSSServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in BatchDeleteFilesRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationOSSServiceBatchDeleteFiles)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.BatchDeleteFiles(ctx, req.(*BatchDeleteFilesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*BatchDeleteFilesReply)
		return ctx.Result(200, reply)
	}
}

func _OSSService_ListFiles0_HTTP_Handler(srv OSSServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListFilesRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationOSSServiceListFiles)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListFiles(ctx, req.(*ListFilesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListFilesReply)
		return ctx.Result(200, reply)
	}
}

func _OSSService_MergeFiles0_HTTP_Handler(srv OSSServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in MergeFilesRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationOSSServiceMergeFiles)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.MergeFiles(ctx, req.(*MergeFilesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*MergeFilesReply)
		return ctx.Result(200, reply)
	}
}

func _OSSService_WriteString0_HTTP_Handler(srv OSSServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in WriteStringRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationOSSServiceWriteString)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.WriteString(ctx, req.(*WriteStringRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*WriteStringReply)
		return ctx.Result(200, reply)
	}
}

func _OSSService_ReadString0_HTTP_Handler(srv OSSServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ReadStringRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationOSSServiceReadString)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ReadString(ctx, req.(*ReadStringRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ReadStringReply)
		return ctx.Result(200, reply)
	}
}

type OSSServiceHTTPClient interface {
	BatchDeleteFiles(ctx context.Context, req *BatchDeleteFilesRequest, opts ...http.CallOption) (rsp *BatchDeleteFilesReply, err error)
	DeleteFile(ctx context.Context, req *DeleteFileRequest, opts ...http.CallOption) (rsp *DeleteFileReply, err error)
	DownloadFile(ctx context.Context, req *DownloadFileRequest, opts ...http.CallOption) (rsp *DownloadFileReply, err error)
	ListFiles(ctx context.Context, req *ListFilesRequest, opts ...http.CallOption) (rsp *ListFilesReply, err error)
	MergeFiles(ctx context.Context, req *MergeFilesRequest, opts ...http.CallOption) (rsp *MergeFilesReply, err error)
	ReadString(ctx context.Context, req *ReadStringRequest, opts ...http.CallOption) (rsp *ReadStringReply, err error)
	TempUpload(ctx context.Context, req *TempUploadRequest, opts ...http.CallOption) (rsp *TempUploadReply, err error)
	UploadFile(ctx context.Context, req *UploadFileRequest, opts ...http.CallOption) (rsp *UploadFileReply, err error)
	WriteString(ctx context.Context, req *WriteStringRequest, opts ...http.CallOption) (rsp *WriteStringReply, err error)
}

type OSSServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewOSSServiceHTTPClient(client *http.Client) OSSServiceHTTPClient {
	return &OSSServiceHTTPClientImpl{client}
}

func (c *OSSServiceHTTPClientImpl) BatchDeleteFiles(ctx context.Context, in *BatchDeleteFilesRequest, opts ...http.CallOption) (*BatchDeleteFilesReply, error) {
	var out BatchDeleteFilesReply
	pattern := "/v1/oss/batch/delete"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationOSSServiceBatchDeleteFiles))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *OSSServiceHTTPClientImpl) DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...http.CallOption) (*DeleteFileReply, error) {
	var out DeleteFileReply
	pattern := "/v1/oss/file/{bucket_name}/{file_path}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationOSSServiceDeleteFile))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *OSSServiceHTTPClientImpl) DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...http.CallOption) (*DownloadFileReply, error) {
	var out DownloadFileReply
	pattern := "/v1/oss/download/{bucket_name}/{file_path}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationOSSServiceDownloadFile))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *OSSServiceHTTPClientImpl) ListFiles(ctx context.Context, in *ListFilesRequest, opts ...http.CallOption) (*ListFilesReply, error) {
	var out ListFilesReply
	pattern := "/v1/oss/list/{bucket_name}/{prefix}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationOSSServiceListFiles))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *OSSServiceHTTPClientImpl) MergeFiles(ctx context.Context, in *MergeFilesRequest, opts ...http.CallOption) (*MergeFilesReply, error) {
	var out MergeFilesReply
	pattern := "/v1/oss/merge"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationOSSServiceMergeFiles))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *OSSServiceHTTPClientImpl) ReadString(ctx context.Context, in *ReadStringRequest, opts ...http.CallOption) (*ReadStringReply, error) {
	var out ReadStringReply
	pattern := "/v1/oss/read/{bucket_name}/{file_path}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationOSSServiceReadString))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *OSSServiceHTTPClientImpl) TempUpload(ctx context.Context, in *TempUploadRequest, opts ...http.CallOption) (*TempUploadReply, error) {
	var out TempUploadReply
	pattern := "/v1/oss/temp/upload"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationOSSServiceTempUpload))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *OSSServiceHTTPClientImpl) UploadFile(ctx context.Context, in *UploadFileRequest, opts ...http.CallOption) (*UploadFileReply, error) {
	var out UploadFileReply
	pattern := "/v1/oss/upload"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationOSSServiceUploadFile))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *OSSServiceHTTPClientImpl) WriteString(ctx context.Context, in *WriteStringRequest, opts ...http.CallOption) (*WriteStringReply, error) {
	var out WriteStringReply
	pattern := "/v1/oss/write"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationOSSServiceWriteString))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
