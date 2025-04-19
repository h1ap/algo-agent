package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ZipFileExtractor ZIP文件处理工具类
type ZipFileExtractor struct{}

// Unzip 解压ZIP文件到指定目录
//
// 参数：
//   - zipFilePath: ZIP文件路径
//   - outputDirectory: 解压后的目录名称
//
// 返回值：
//   - error: 如果发生错误则返回错误信息
func (ze *ZipFileExtractor) Unzip(zipFilePath, outputDirectory string) error {
	// 创建输出目录
	if err := os.MkdirAll(outputDirectory, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 打开ZIP文件
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("打开ZIP文件失败: %v", err)
	}
	defer reader.Close()

	// 遍历处理ZIP文件中的每个条目
	for _, file := range reader.File {
		err := ze.extractFile(file, outputDirectory)
		if err != nil {
			return fmt.Errorf("解压文件 %s 失败: %v", file.Name, err)
		}
	}

	return nil
}

// UnzipWithoutTopDirectory 解压ZIP文件到指定目录，并剥离第一层目录
//
// 参数：
//   - zipFilePath: ZIP文件路径
//   - outputDirectory: 解压后的目录名称
//
// 返回值：
//   - error: 如果发生错误则返回错误信息
func (ze *ZipFileExtractor) UnzipWithoutTopDirectory(zipFilePath, outputDirectory string) error {
	// 创建输出目录
	if err := os.MkdirAll(outputDirectory, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 打开ZIP文件
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("打开ZIP文件失败: %v", err)
	}
	defer reader.Close()

	// 遍历处理ZIP文件中的每个条目
	for _, file := range reader.File {
		// 去掉第一层目录
		strippedName := ze.stripTopDirectory(file.Name)
		if strippedName == "" {
			continue // 跳过顶层目录
		}

		// 创建新的zip.File，但使用剥离后的路径
		err := ze.extractFileWithPath(file, outputDirectory, strippedName)
		if err != nil {
			return fmt.Errorf("解压文件 %s 失败: %v", file.Name, err)
		}
	}

	return nil
}

// UnzipCurrentFolder 将当前目录下的所有zip文件解压到当前目录
//
// 参数：
//   - datasetAnnotationPath: 数据集注释路径
//
// 返回值：
//   - error: 如果发生错误则返回错误信息
func (ze *ZipFileExtractor) UnzipCurrentFolder(datasetAnnotationPath string) error {
	files, err := os.ReadDir(datasetAnnotationPath)
	if err != nil {
		return fmt.Errorf("读取目录失败: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".zip") {
			zipPath := filepath.Join(datasetAnnotationPath, file.Name())
			err := ze.Unzip(zipPath, datasetAnnotationPath)
			if err != nil {
				return fmt.Errorf("解压文件 %s 失败: %v", zipPath, err)
			}
		}
	}

	return nil
}

// extractFile 解压单个文件
func (ze *ZipFileExtractor) extractFile(file *zip.File, destDir string) error {
	// 构建完整的目标路径
	destPath := filepath.Join(destDir, file.Name)

	// 如果是目录，创建它
	if file.FileInfo().IsDir() {
		return os.MkdirAll(destPath, 0755)
	}

	// 确保父目录存在
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// 创建目标文件
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 打开ZIP中的文件
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 复制内容
	_, err = io.Copy(destFile, srcFile)
	return err
}

// extractFileWithPath 使用指定路径解压单个文件
func (ze *ZipFileExtractor) extractFileWithPath(file *zip.File, destDir, targetPath string) error {
	// 构建完整的目标路径
	destPath := filepath.Join(destDir, targetPath)

	// 如果是目录，创建它
	if file.FileInfo().IsDir() {
		return os.MkdirAll(destPath, 0755)
	}

	// 确保父目录存在
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// 创建目标文件
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 打开ZIP中的文件
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 复制内容
	_, err = io.Copy(destFile, srcFile)
	return err
}

// stripTopDirectory 去掉路径的第一个目录层级
func (ze *ZipFileExtractor) stripTopDirectory(entryName string) string {
	parts := strings.Split(entryName, "/")
	if len(parts) <= 1 {
		return "" // 没有子目录，返回空字符串
	}
	return filepath.Join(parts[1:]...)
}
