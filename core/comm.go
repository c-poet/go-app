package core

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetWdHome 获取工作目录。优先使用启动参数设置的目录，未设置时使用可执行文件所在目录。
func GetWdHome() (string, error) {
	confDir := os.Getenv("app.home")
	if confDir == "" {
		executable, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("获取可执行文件路径失败: %v", err)
		}
		confDir = filepath.Dir(executable)
	}
	// 确保工作目录存在
	if err := os.MkdirAll(confDir, 0755); err != nil {
		return "", fmt.Errorf("创建 %s 目录失败: %v", confDir, err)
	}
	return confDir, nil
}

// GetWdSpecifyDir 工作目录下的获取特定目录
func GetWdSpecifyDir(name string) (string, error) {
	wdHome, err := GetWdHome()
	if err != nil {
		return "", err
	}
	specDir := filepath.Join(wdHome, name)
	if err := os.MkdirAll(specDir, 0755); err != nil {
		return "", fmt.Errorf("创建 %s 目录失败: %v", specDir, err)
	}
	return specDir, nil
}

// GetWdDataFilePath 获取指定数据文件的路径
func GetWdDataFilePath(filename string) string {
	dataDir, err := GetWdSpecifyDir(WorkDataDir)
	if err != nil {
		return filename
	}
	return filepath.Join(dataDir, filename)
}

// GetWdLogFilePath 获取指定数据文件的路径
func GetWdLogFilePath(filename string) string {
	dataDir, err := GetWdSpecifyDir(WorkLogDir)
	if err != nil {
		return filename
	}
	return filepath.Join(dataDir, filename)
}

// GetWdTempFilePath 获取指定临时文件的路径
func GetWdTempFilePath(filename string) string {
	dataDir, err := GetWdSpecifyDir(WorkTempDir)
	if err != nil {
		return filename
	}
	return filepath.Join(dataDir, filename)
}

// ReadWdFile 读取工作目录中的文件内容
func ReadWdFile(name string) ([]byte, error) {
	wdHome, err := GetWdHome()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(wdHome, name)
	return os.ReadFile(path)
}

// WriteWdFile 写入文件内容到工作目录
func WriteWdFile(name string, data []byte) error {
	wdHome, err := GetWdHome()
	if err != nil {
		return err
	}
	path := filepath.Join(wdHome, name)
	return os.WriteFile(path, data, 0644)
}
