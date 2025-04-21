package biz

import (
	"context"
	"io"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// File 文件模型
type File struct {
	Name        string
	Size        int64
	ContentType string
	Content     []byte
}

// FileProgress 文件上传进度
type FileProgress struct {
	UUID            string
	CurrentFile     string
	ChunksPath      string
	SHA256Sign      string
	Progress        float64
	DateTime        time.Time
	ChunkNumber     int
	NextChunkNumber int
	TotalChunks     int
	IsEnd           bool
	ProgressType    int
}

// OSSService OSS存储接口
type OSSService interface {
	// TempUpload 上传临时文件
	TempUpload(ctx context.Context, bucketName string, file *File) (string, error)

	// UploadFile 上传文件
	UploadFile(ctx context.Context, bucketName string, file *File, filePath string) error

	// UploadFileWithProgress 上传文件并监听进度
	UploadFileWithProgress(ctx context.Context, bucketName string, file *File, filePath string, progressCallback func(float64)) error

	// UploadStream 上传流数据
	UploadStream(ctx context.Context, bucketName string, stream io.Reader, size int64, filePath string) error

	// MergeFiles 合并文件
	MergeFiles(ctx context.Context, bucketName string, sourceFiles []string, targetFile string) error

	// DownloadFile 下载文件
	DownloadFile(ctx context.Context, bucketName string, key string) ([]byte, error)

	// GetFileStream 获取文件流
	GetFileStream(ctx context.Context, bucketName string, key string) (io.ReadCloser, error)

	// DeleteFile 删除文件
	DeleteFile(ctx context.Context, bucketName string, key string) error

	// BatchDeleteFiles 批量删除文件
	BatchDeleteFiles(ctx context.Context, bucketName string, keys []string) error

	// ListFiles 列出指定路径下的文件
	ListFiles(ctx context.Context, bucketName string, prefix string) ([]string, error)

	// WriteString 写入字符串内容
	WriteString(ctx context.Context, bucketName string, content string, filePath string) error

	// ReadString 读取字符串内容
	ReadString(ctx context.Context, bucketName string, filePath string) (string, error)

	// DownloadSingleFile 下载单个文件到指定路径
	DownloadSingleFile(ctx context.Context, bucketName, objectName, destinationPath, fileName string) error

	// Close 关闭连接
	Close()
}

// OSSUsecase OSS用例
type OSSUsecase struct {
	store OSSService
	log   *log.Helper
}

// TempUpload 上传临时文件
func (uc *OSSUsecase) TempUpload(ctx context.Context, bucketName string, file *File) (string, error) {
	uc.log.WithContext(ctx).Infof("TempUpload: fileName=%s, fileSize=%d", file.Name, file.Size)
	return uc.store.TempUpload(ctx, bucketName, file)
}

// UploadFile 上传文件
func (uc *OSSUsecase) UploadFile(ctx context.Context, bucketName string, file *File, filePath string) error {
	uc.log.WithContext(ctx).Infof("UploadFile: fileName=%s, filePath=%s, fileSize=%d", file.Name, filePath, file.Size)
	return uc.store.UploadFile(ctx, bucketName, file, filePath)
}

// UploadFileWithProgress 上传文件并监听进度
func (uc *OSSUsecase) UploadFileWithProgress(ctx context.Context, bucketName string, file *File, filePath string, progressCallback func(float64)) error {
	uc.log.WithContext(ctx).Infof("UploadFileWithProgress: fileName=%s, filePath=%s, fileSize=%d", file.Name, filePath, file.Size)
	return uc.store.UploadFileWithProgress(ctx, bucketName, file, filePath, progressCallback)
}

// UploadStream 上传流数据
func (uc *OSSUsecase) UploadStream(ctx context.Context, bucketName string, stream io.Reader, size int64, filePath string) error {
	uc.log.WithContext(ctx).Infof("UploadStream: filePath=%s, size=%d", filePath, size)
	return uc.store.UploadStream(ctx, bucketName, stream, size, filePath)
}

// MergeFiles 合并文件
func (uc *OSSUsecase) MergeFiles(ctx context.Context, bucketName string, sourceFiles []string, targetFile string) error {
	uc.log.WithContext(ctx).Infof("MergeFiles: targetFile=%s, sourceCount=%d", targetFile, len(sourceFiles))
	return uc.store.MergeFiles(ctx, bucketName, sourceFiles, targetFile)
}

// DownloadFile 下载文件
func (uc *OSSUsecase) DownloadFile(ctx context.Context, bucketName string, key string) ([]byte, error) {
	uc.log.WithContext(ctx).Infof("DownloadFile: key=%s", key)
	return uc.store.DownloadFile(ctx, bucketName, key)
}

// GetFileStream 获取文件流
func (uc *OSSUsecase) GetFileStream(ctx context.Context, bucketName string, key string) (io.ReadCloser, error) {
	uc.log.WithContext(ctx).Infof("GetFileStream: key=%s", key)
	return uc.store.GetFileStream(ctx, bucketName, key)
}

// DeleteFile 删除文件
func (uc *OSSUsecase) DeleteFile(ctx context.Context, bucketName string, key string) error {
	uc.log.WithContext(ctx).Infof("DeleteFile: key=%s", key)
	return uc.store.DeleteFile(ctx, bucketName, key)
}

// BatchDeleteFiles 批量删除文件
func (uc *OSSUsecase) BatchDeleteFiles(ctx context.Context, bucketName string, keys []string) error {
	uc.log.WithContext(ctx).Infof("BatchDeleteFiles: keyCount=%d", len(keys))
	return uc.store.BatchDeleteFiles(ctx, bucketName, keys)
}

// ListFiles 列出指定路径下的文件
func (uc *OSSUsecase) ListFiles(ctx context.Context, bucketName string, prefix string) ([]string, error) {
	uc.log.WithContext(ctx).Infof("ListFiles: prefix=%s", prefix)
	return uc.store.ListFiles(ctx, bucketName, prefix)
}

// WriteString 写入字符串内容
func (uc *OSSUsecase) WriteString(ctx context.Context, bucketName string, content string, filePath string) error {
	uc.log.WithContext(ctx).Infof("WriteString: filePath=%s, contentLength=%d", filePath, len(content))
	return uc.store.WriteString(ctx, bucketName, content, filePath)
}

// ReadString 读取字符串内容
func (uc *OSSUsecase) ReadString(ctx context.Context, bucketName string, filePath string) (string, error) {
	uc.log.WithContext(ctx).Infof("ReadString: filePath=%s", filePath)
	return uc.store.ReadString(ctx, bucketName, filePath)
}
