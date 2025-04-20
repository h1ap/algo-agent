package data

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"algo-agent/internal/biz"
	"algo-agent/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// OSSRepo 实现OSS客户端，满足biz.OSSStore接口
type OSSRepo struct {
	client *minio.Client
	conf   *conf.Data_Oss
	log    *log.Helper
}

// TempUpload 上传临时文件
func (r *OSSRepo) TempUpload(ctx context.Context, bucketName string, file *biz.File) (string, error) {
	r.log.WithContext(ctx).Infof("TempUpload: fileName=%s, fileSize=%d", file.Name, file.Size)

	tempPath := generateTempPath(file.Name)
	reader := bytes.NewReader(file.Content)

	_, err := r.client.PutObject(ctx, bucketName, tempPath, reader, file.Size,
		minio.PutObjectOptions{ContentType: file.ContentType})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to upload temp file: %v", err)
		return "", err
	}

	return tempPath, nil
}

// UploadFile 上传文件
func (r *OSSRepo) UploadFile(ctx context.Context, bucketName string, file *biz.File, filePath string) error {
	r.log.WithContext(ctx).Infof("UploadFile: fileName=%s, filePath=%s, fileSize=%d", file.Name, filePath, file.Size)

	reader := bytes.NewReader(file.Content)
	opts := minio.PutObjectOptions{ContentType: file.ContentType}

	// 文件大于50MB时使用分片上传
	if file.Size > 50*1024*1024 {
		opts.PartSize = 10 * 1024 * 1024 // 10MB分片
		r.log.WithContext(ctx).Infof("文件尺寸大于50MB，使用分片上传: %s", filePath)
	}

	_, err := r.client.PutObject(ctx, bucketName, filePath, reader, file.Size, opts)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to upload file: %v", err)
		return err
	}

	return nil
}

// UploadFileWithProgress 上传文件并监听进度
func (r *OSSRepo) UploadFileWithProgress(ctx context.Context, bucketName string, file *biz.File, filePath string, progressCallback func(float64)) error {
	r.log.WithContext(ctx).Infof("UploadFileWithProgress: fileName=%s, filePath=%s, fileSize=%d", file.Name, filePath, file.Size)

	reader := bytes.NewReader(file.Content)
	opts := minio.PutObjectOptions{ContentType: file.ContentType}

	// 设置进度回调
	if progressCallback != nil {
		opts.Progress = &progressReader{
			reader:   reader,
			size:     file.Size,
			callback: progressCallback,
		}
	}

	// 文件大于50MB时使用分片上传
	if file.Size > 50*1024*1024 {
		opts.PartSize = 10 * 1024 * 1024 // 10MB分片
		r.log.WithContext(ctx).Infof("文件尺寸大于50MB，使用分片上传: %s", filePath)
	}

	_, err := r.client.PutObject(ctx, bucketName, filePath, reader, file.Size, opts)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to upload file with progress: %v", err)
		return err
	}

	return nil
}

// UploadStream 上传流数据
func (r *OSSRepo) UploadStream(ctx context.Context, bucketName string, stream io.Reader, size int64, filePath string) error {
	r.log.WithContext(ctx).Infof("UploadStream: filePath=%s, size=%d", filePath, size)

	opts := minio.PutObjectOptions{}

	// 文件大于50MB时使用分片上传
	if size > 50*1024*1024 {
		opts.PartSize = 10 * 1024 * 1024 // 10MB分片
		r.log.WithContext(ctx).Infof("文件尺寸大于50MB，使用分片上传: %s", filePath)
	}

	_, err := r.client.PutObject(ctx, bucketName, filePath, stream, size, opts)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to upload stream: %v", err)
		return err
	}

	return nil
}

// MergeFiles 合并文件
func (r *OSSRepo) MergeFiles(ctx context.Context, bucketName string, sourceFiles []string, targetFile string) error {
	r.log.WithContext(ctx).Infof("MergeFiles: targetFile=%s, sourceCount=%d", targetFile, len(sourceFiles))

	if len(sourceFiles) == 0 {
		return errors.New("source files list is empty")
	}

	var sources []minio.CopySrcOptions
	for _, src := range sourceFiles {
		sources = append(sources, minio.CopySrcOptions{
			Bucket: bucketName,
			Object: src,
		})
	}

	// 使用ComposeObject合并文件
	dst := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: targetFile,
	}

	_, err := r.client.ComposeObject(ctx, dst, sources...)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to merge files: %v", err)
		return err
	}

	// 合并完成后删除源文件
	for _, src := range sourceFiles {
		err := r.client.RemoveObject(ctx, bucketName, src, minio.RemoveObjectOptions{})
		if err != nil {
			r.log.WithContext(ctx).Warnf("failed to delete source file after merge: %s, %v", src, err)
			// 继续处理，不中断流程
		}
	}

	return nil
}

// DownloadFile 下载文件
func (r *OSSRepo) DownloadFile(ctx context.Context, bucketName string, key string) ([]byte, error) {
	r.log.WithContext(ctx).Infof("DownloadFile: key=%s", key)

	obj, err := r.client.GetObject(ctx, bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to get object: %v", err)
		return nil, err
	}
	defer obj.Close()

	return ioutil.ReadAll(obj)
}

// GetFileStream 获取文件流
func (r *OSSRepo) GetFileStream(ctx context.Context, bucketName string, key string) (io.ReadCloser, error) {
	r.log.WithContext(ctx).Infof("GetFileStream: key=%s", key)

	obj, err := r.client.GetObject(ctx, bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to get object stream: %v", err)
		return nil, err
	}

	return obj, nil
}

// DeleteFile 删除文件
func (r *OSSRepo) DeleteFile(ctx context.Context, bucketName string, key string) error {
	r.log.WithContext(ctx).Infof("DeleteFile: key=%s", key)

	err := r.client.RemoveObject(ctx, bucketName, key, minio.RemoveObjectOptions{})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to delete file: %v", err)
		return err
	}

	return nil
}

// BatchDeleteFiles 批量删除文件
func (r *OSSRepo) BatchDeleteFiles(ctx context.Context, bucketName string, keys []string) error {
	r.log.WithContext(ctx).Infof("BatchDeleteFiles: keyCount=%d", len(keys))

	if len(keys) == 0 {
		return nil
	}

	objectsCh := make(chan minio.ObjectInfo)

	// 生成待删除对象
	go func() {
		defer close(objectsCh)
		for _, key := range keys {
			objectsCh <- minio.ObjectInfo{
				Key: key,
			}
		}
	}()

	// 执行删除
	for err := range r.client.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{}) {
		if err.Err != nil {
			r.log.WithContext(ctx).Errorf("failed to delete object %s: %v", err.ObjectName, err.Err)
			// 继续处理其他文件，不中断流程
		}
	}

	return nil
}

// ListFiles 列出指定路径下的文件
func (r *OSSRepo) ListFiles(ctx context.Context, bucketName string, prefix string) ([]string, error) {
	r.log.WithContext(ctx).Infof("ListFiles: prefix=%s", prefix)

	var files []string

	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}

	for obj := range r.client.ListObjects(ctx, bucketName, opts) {
		if obj.Err != nil {
			r.log.WithContext(ctx).Errorf("failed to list objects: %v", obj.Err)
			return nil, obj.Err
		}
		files = append(files, obj.Key)
	}

	return files, nil
}

// WriteString 写入字符串内容
func (r *OSSRepo) WriteString(ctx context.Context, bucketName string, content string, filePath string) error {
	r.log.WithContext(ctx).Infof("WriteString: filePath=%s, contentLength=%d", filePath, len(content))

	reader := bytes.NewReader([]byte(content))

	_, err := r.client.PutObject(ctx, bucketName, filePath, reader, int64(len(content)),
		minio.PutObjectOptions{ContentType: "text/plain; charset=utf-8"})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to write string: %v", err)
		return err
	}

	return nil
}

// ReadString 读取字符串内容
func (r *OSSRepo) ReadString(ctx context.Context, bucketName string, filePath string) (string, error) {
	r.log.WithContext(ctx).Infof("ReadString: filePath=%s", filePath)

	obj, err := r.client.GetObject(ctx, bucketName, filePath, minio.GetObjectOptions{})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to get object for reading: %v", err)
		return "", err
	}
	defer obj.Close()

	data, err := ioutil.ReadAll(obj)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to read object content: %v", err)
		return "", err
	}

	return string(data), nil
}

// DownloadSingleFile 下载单个文件到指定路径
func (r *OSSRepo) DownloadSingleFile(ctx context.Context, bucketName, objectName, destinationPath, fileName string) error {
	r.log.WithContext(ctx).Infof("DownloadSingleFile: objectName=%s, destinationPath=%s", objectName, destinationPath)

	// 如果文件名为空，沿用对象名中的文件名
	if fileName == "" {
		fileName = filepath.Base(objectName) // 提取对象中的文件名
	}

	// 拼接完整的目标文件路径
	destinationFile := filepath.Join(destinationPath, fileName)

	// 确保目标目录存在
	parentDir := filepath.Dir(destinationFile)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		r.log.WithContext(ctx).Errorf("failed to create destination directory: %v", err)
		return err
	}

	// 下载文件
	obj, err := r.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to get object for download: %v", err)
		return err
	}
	defer obj.Close()

	// 创建目标文件
	file, err := os.Create(destinationFile)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to create destination file: %v", err)
		return err
	}
	defer file.Close()

	// 复制内容
	buffer := make([]byte, 1024)
	_, err = io.CopyBuffer(file, obj, buffer)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to copy file content: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("successfully downloaded file: %s to %s", objectName, destinationFile)
	return nil
}

// Close 关闭连接
func (r *OSSRepo) Close() {
	// MinIO客户端不需要显式关闭
	r.log.Info("OSS connection closed")
}

// 生成临时文件路径
func generateTempPath(fileName string) string {
	now := time.Now()
	dir := fmt.Sprintf("temp/%d/%02d/%02d", now.Year(), now.Month(), now.Day())
	uniqueID := uuid.New().String()
	ext := filepath.Ext(fileName)

	return fmt.Sprintf("%s/%s%s", dir, uniqueID, ext)
}

// progressReader 用于跟踪进度的读取器
type progressReader struct {
	reader      io.Reader
	size        int64
	read        int64
	callback    func(float64)
	lastPercent float64
}

// Read 实现io.Reader接口
func (r *progressReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if n > 0 {
		r.read += int64(n)
		percent := float64(r.read) / float64(r.size)

		// 只在进度有明显变化时回调，避免过多调用
		if percent-r.lastPercent >= 0.01 {
			r.callback(percent)
			r.lastPercent = percent
		}
	}
	return
}

// 示例代码：
// modelPath := deployRequest.ModelPath
// fileName := filepath.Base(modelPath) // 提取文件路径中的文件名部分
