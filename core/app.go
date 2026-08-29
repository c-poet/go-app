package core

import (
	"log/slog"
	"strings"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

// buildTime 编译时间，格式：年.两位月份两位日.两位时两位分 (例如: 26.0515.1019)
// 该变量在编译时通过 -ldflags 注入
var buildTime = "unknown"

// Application 应用，单例，用于传递全局上下文
type Application struct {
	name       string
	confName   string
	logName    string
	logEnabled bool
	cfg        *Conf
	props      map[string]interface{}
}

var (
	instance *Application
	once     sync.Once
)

// App 获取应用单例实例
func App() *Application {
	return instance
}

// AppOptions 应用初始化选项
type AppOptions struct {
	Name       string // 应用名称，默认为"app"
	ConfName   string // 配置名称，默认为应用名称
	LogName    string // 日志名称，默认为应用名称
	LogEnabled bool   // 是否启用日志，默认为 true
}

// AppOption 选项函数类型
type AppOption func(*AppOptions)

// WithName 设置应用名称的选项函数
func WithName(name string) AppOption {
	return func(opts *AppOptions) {
		opts.Name = name
	}
}

// WithConfName 设置配置名称的选项函数
func WithConfName(confName string) AppOption {
	return func(opts *AppOptions) {
		opts.ConfName = confName
	}
}

// WithLogName 设置日志名称的选项函数
func WithLogName(logName string) AppOption {
	return func(opts *AppOptions) {
		opts.LogName = logName
	}
}

// WithLogEnabled 设置是否启用日志的选项函数
func WithLogEnabled(enabled bool) AppOption {
	return func(opts *AppOptions) {
		opts.LogEnabled = enabled
	}
}

// InitApp 初始化应用单例
func InitApp(opts ...AppOption) (*Application, error) {
	appOpts := newAppOptions(opts...)

	var initErr error
	once.Do(func() {
		cfg, err := NewConf(appOpts.ConfName)
		if err != nil {
			initErr = err
			return
		}
		// 初始化应用时读取配置
		err = cfg.Read()
		if err != nil {
			initErr = err
			return
		}
		instance = &Application{
			name:       appOpts.Name,
			confName:   appOpts.ConfName,
			logName:    appOpts.LogName,
			logEnabled: appOpts.LogEnabled,
			cfg:        cfg,
		}
	})
	if initErr != nil {
		return nil, initErr
	}
	return instance, nil
}

func newAppOptions(opts ...AppOption) *AppOptions {
	// 应用默认选项
	appOpts := &AppOptions{
		Name:       "app",
		LogEnabled: true,
	}

	// 应用传入的选项
	for _, opt := range opts {
		opt(appOpts)
	}
	if appOpts.LogName == "" {
		appOpts.LogName = appOpts.Name
	}
	if appOpts.ConfName == "" {
		appOpts.ConfName = appOpts.Name
	}
	return appOpts
}

// Name 获取应用名称
func (app *Application) Name() string {
	return app.name
}

// ConfName 获取配置名称
func (app *Application) ConfName() string {
	return app.confName
}

// LogName 获取日志名称
func (app *Application) LogName() string {
	return app.logName
}

// LogEnabled 获取日志是否启用
func (app *Application) LogEnabled() bool {
	return app.logEnabled
}

// Cfg 获取应用配置
func (app *Application) Cfg() *Conf {
	return app.cfg
}

// Version 返回应用版本号（编译时间）
func (app *Application) Version() string {
	return buildTime
}

// GetProps 获取属性（支持以.分隔的key）
func (app *Application) GetProps(key string) interface{} {
	if app.props == nil {
		return nil
	}

	// 按.分割key
	keys := strings.Split(key, ".")

	// 从顶层map开始，一层一层查找
	var current interface{} = app.props
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

// SetProps 设置属性（支持以.分隔的key）
func (app *Application) SetProps(key string, value interface{}) *Application {
	if app.props == nil {
		app.props = make(map[string]interface{})
	}

	// 按.分割key
	keys := strings.Split(key, ".")

	// 从顶层map开始，一层一层查找或创建
	current := app.props
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

	return app
}

// GetConf 获取配置（支持以.分隔的key），先从props查找，如果不存在则从cfg查找
func (app *Application) GetConf(key string) interface{} {
	// 先从props查找
	value := app.GetProps(key)
	if value != nil {
		return value
	}
	// 如果props中不存在，则从cfg查找
	return app.cfg.Get(key)
}

func (app *Application) ReloadLogConf() *Application {
	if !app.LogEnabled() {
		return app
	}
	roller := &lumberjack.Logger{
		Filename:   GetWdLogFilePath(app.LogName() + LogExt),
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     28,
		Compress:   true,
	}
	handler := slog.NewTextHandler(roller, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))
	return app
}

// Destroy 关闭应用时调用，同步写入配置文件
func (app *Application) Destroy() {
	if err := app.cfg.Write(); err != nil {
		slog.Error("保存配置文件失败", "error", err)
	}
}
