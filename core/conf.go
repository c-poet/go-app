package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Conf 配置信息
type Conf struct {
	confFile string
	data     map[string]interface{}
	mu       sync.RWMutex
}

// NewConf 创建配置管理器
func NewConf(fileName string) (*Conf, error) {
	// 配置文件直接存放在工作目录中。
	confDir, err := GetWdHome()
	if err != nil {
		return nil, err
	}
	// 处理文件后缀：如果没有后缀则添加 ConfExt
	if filepath.Ext(fileName) == "" {
		fileName = fileName + ConfExt
	}
	confFile := filepath.Join(confDir, fileName)
	return &Conf{
		confFile: confFile,
	}, nil
}

// GetConfFile 获取配置文件路径
func (c *Conf) GetConfFile() string {
	return c.confFile
}

// Get 获取配置项（支持以.分隔的key）
func (c *Conf) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.data == nil {
		return nil
	}

	// 按.分割key
	keys := strings.Split(key, ".")

	// 从顶层map开始，一层一层查找
	var current interface{} = c.data
	for _, k := range keys {
		if currentMap, ok := current.(map[string]interface{}); ok {
			if val, exists := currentMap[k]; exists {
				current = val
			} else {
				return nil
			}
		} else {
			// 当前层不是map，无法继续查找
			return nil
		}
	}

	return current
}

// Set 设置配置项（支持以.分隔的key）
func (c *Conf) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.data == nil {
		c.data = make(map[string]interface{})
	}

	// 按.分割key
	keys := strings.Split(key, ".")

	// 从顶层map开始，一层一层查找或创建
	current := c.data
	for i := 0; i < len(keys)-1; i++ {
		k := keys[i]
		if nextMap, ok := current[k].(map[string]interface{}); ok {
			current = nextMap
		} else {
			// 创建新的map
			newMap := make(map[string]interface{})
			current[k] = newMap
			current = newMap
		}
	}

	// 设置最后一层的值
	current[keys[len(keys)-1]] = value
}

// Exists 检查文件是否存在
func (c *Conf) Exists(name string) bool {
	wdHome, err := GetWdHome()
	if err != nil {
		return false
	}
	path := filepath.Join(wdHome, name)
	_, err = os.Stat(path)
	return err == nil
}

// Read 从confFile中读取配置并反序列化到data
func (c *Conf) Read() error {
	data, err := os.ReadFile(c.confFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，保持data不变
			return nil
		}
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	if len(data) == 0 {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := json.Unmarshal(data, &c.data); err != nil {
		return fmt.Errorf("反序列化配置失败: %v", err)
	}

	return nil
}

// Write 把data序列化并写入confFile
func (c *Conf) Write() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.data == nil {
		return nil
	}

	data, err := json.Marshal(c.data)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(c.confFile, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}
