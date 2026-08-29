# go-app

Go 应用基础组件库，提供应用单例、JSON 配置、工作目录文件操作，以及基于 `slog` 的滚动日志配置。

## 安装

```bash
go get github.com/c-poet/go-app
```

## 快速开始

```go
package main

import (
	"log/slog"

	"github.com/c-poet/go-app/core"
)

func main() {
	app, err := core.InitApp(core.WithName("example"))
	if err != nil {
		panic(err)
	}
	defer app.Destroy()

	app.Cfg().Set("server.port", 8080)
	app.ReloadLogConf()
	slog.Info("application started", "version", app.Version())
}
```

`InitApp` 会读取 `<应用名>.json` 配置文件；`Destroy` 会将配置写回。默认工作目录是可执行文件所在目录，也可通过 `ctj.cmd.home` 环境变量指定。日志会写入工作目录的 `logs/` 子目录。

## 开发

```bash
gofmt -w core/*.go
go test ./...
```

## 包结构

- `core`：应用初始化、配置管理、工作目录帮助函数和日志配置。
