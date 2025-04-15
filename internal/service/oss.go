package service

import (
	pb "algo-agent/api/oss/v1"
	"algo-agent/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// OSSServer 实现OSS服务API接口
type OSSServer struct {
	pb.UnimplementedOSSServiceServer
	uc  *biz.OSSUsecase
	log *log.Helper
}

// TempUpload 上传临时文件
func (s *OSSServer) TempUpload(ctx context.Context, req *pb.TempUploadRequest) (*pb.TempUploadReply, error) {
	s.log.WithContext(ctx).Infof("TempUpload: bucketName=%s, fileName=%s", req.BucketName, req.File.Name)

	file := &biz.File{
		Name:        req.File.Name,
		Size:        req.File.Size,
		ContentType: req.File.ContentType,
		Content:     req.File.Content,
	}

	path, err := s.uc.TempUpload(ctx, req.BucketName, file)
	if err != nil {
		s.log.WithContext(ctx).Errorf("TempUpload failed: %v", err)
		return nil, err
	}

	return &pb.TempUploadReply{FilePath: path}, nil
}

// UploadFile 上传文件
func (s *OSSServer) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileReply, error) {
	s.log.WithContext(ctx).Infof("UploadFile: bucketName=%s, fileName=%s, filePath=%s",
		req.BucketName, req.File.Name, req.FilePath)

	file := &biz.File{
		Name:        req.File.Name,
		Size:        req.File.Size,
		ContentType: req.File.ContentType,
		Content:     req.File.Content,
	}

	err := s.uc.UploadFile(ctx, req.BucketName, file, req.FilePath)
	if err != nil {
		s.log.WithContext(ctx).Errorf("UploadFile failed: %v", err)
		return nil, err
	}

	return &pb.UploadFileReply{Success: true}, nil
}

// DownloadFile 下载文件
func (s *OSSServer) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileReply, error) {
	s.log.WithContext(ctx).Infof("DownloadFile: bucketName=%s, filePath=%s", req.BucketName, req.FilePath)

	content, err := s.uc.DownloadFile(ctx, req.BucketName, req.FilePath)
	if err != nil {
		s.log.WithContext(ctx).Errorf("DownloadFile failed: %v", err)
		return nil, err
	}

	return &pb.DownloadFileReply{Content: content}, nil
}

// DeleteFile 删除文件
func (s *OSSServer) DeleteFile(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileReply, error) {
	s.log.WithContext(ctx).Infof("DeleteFile: bucketName=%s, filePath=%s", req.BucketName, req.FilePath)

	err := s.uc.DeleteFile(ctx, req.BucketName, req.FilePath)
	if err != nil {
		s.log.WithContext(ctx).Errorf("DeleteFile failed: %v", err)
		return nil, err
	}

	return &pb.DeleteFileReply{Success: true}, nil
}

// BatchDeleteFiles 批量删除文件
func (s *OSSServer) BatchDeleteFiles(ctx context.Context, req *pb.BatchDeleteFilesRequest) (*pb.BatchDeleteFilesReply, error) {
	s.log.WithContext(ctx).Infof("BatchDeleteFiles: bucketName=%s, fileCount=%d", req.BucketName, len(req.FilePaths))

	err := s.uc.BatchDeleteFiles(ctx, req.BucketName, req.FilePaths)
	if err != nil {
		s.log.WithContext(ctx).Errorf("BatchDeleteFiles failed: %v", err)
		return nil, err
	}

	return &pb.BatchDeleteFilesReply{Success: true}, nil
}

// ListFiles 列出文件
func (s *OSSServer) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesReply, error) {
	s.log.WithContext(ctx).Infof("ListFiles: bucketName=%s, prefix=%s", req.BucketName, req.Prefix)

	files, err := s.uc.ListFiles(ctx, req.BucketName, req.Prefix)
	if err != nil {
		s.log.WithContext(ctx).Errorf("ListFiles failed: %v", err)
		return nil, err
	}

	return &pb.ListFilesReply{FilePaths: files}, nil
}

// MergeFiles 合并文件
func (s *OSSServer) MergeFiles(ctx context.Context, req *pb.MergeFilesRequest) (*pb.MergeFilesReply, error) {
	s.log.WithContext(ctx).Infof("MergeFiles: bucketName=%s, sourceCount=%d, targetFile=%s",
		req.BucketName, len(req.SourceFiles), req.TargetFile)

	err := s.uc.MergeFiles(ctx, req.BucketName, req.SourceFiles, req.TargetFile)
	if err != nil {
		s.log.WithContext(ctx).Errorf("MergeFiles failed: %v", err)
		return nil, err
	}

	return &pb.MergeFilesReply{Success: true}, nil
}

// WriteString 写入字符串
func (s *OSSServer) WriteString(ctx context.Context, req *pb.WriteStringRequest) (*pb.WriteStringReply, error) {
	s.log.WithContext(ctx).Infof("WriteString: bucketName=%s, filePath=%s, contentLength=%d",
		req.BucketName, req.FilePath, len(req.Content))

	err := s.uc.WriteString(ctx, req.BucketName, req.Content, req.FilePath)
	if err != nil {
		s.log.WithContext(ctx).Errorf("WriteString failed: %v", err)
		return nil, err
	}

	return &pb.WriteStringReply{Success: true}, nil
}

// ReadString 读取字符串
func (s *OSSServer) ReadString(ctx context.Context, req *pb.ReadStringRequest) (*pb.ReadStringReply, error) {
	s.log.WithContext(ctx).Infof("ReadString: bucketName=%s, filePath=%s", req.BucketName, req.FilePath)

	content, err := s.uc.ReadString(ctx, req.BucketName, req.FilePath)
	if err != nil {
		s.log.WithContext(ctx).Errorf("ReadString failed: %v", err)
		return nil, err
	}

	return &pb.ReadStringReply{Content: content}, nil
}
