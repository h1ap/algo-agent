package utils

import "os"

// EnsureDirectoryExists 检查指定目录是否存在，如果不存在则创建目录
//
// 参数：
//   - directoryPath: 要检查或创建的目录路径
//
// 返回值：
//   - bool: 如果目录已存在或创建成功，返回 true；如果创建失败，返回 false
//   - error: 如果发生错误，返回相应的错误信息
func EnsureDirectoryExists(directoryPath string) (bool, error) {
	// 获取目录信息
	info, err := os.Stat(directoryPath)

	// 如果目录已存在且是目录类型
	if err == nil && info.IsDir() {
		return true, nil
	}

	// 如果错误不是"不存在"，则返回错误
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	// 创建目录（包括所有必需的父目录）
	err = os.MkdirAll(directoryPath, 0755)
	if err != nil {
		return false, err
	}

	return true, nil
}

// RemoveDirectory 删除指定的目录及其中的所有内容
//
// 参数：
//   - directoryPath: 要删除的目录路径
//
// 返回值：
//   - error: 如果发生错误，返回相应的错误信息
func RemoveDirectory(directoryPath string) error {
	return os.RemoveAll(directoryPath)
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
