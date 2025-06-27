# Files 包

## 概述

`files` 包提供了文件处理相关的高级类型，包括文件读取（`Reader`）和文件写入（`Writer`）功能。该包设计用于处理文件的上传、下载、存储和管理。

## 重构说明

为了避免 `typex` 包变得过于臃肿，我们将文件处理相关的复杂类型移动到了独立的 `files` 子包中。这样的设计有以下优势：

1. **模块化**: 文件处理功能独立成包，职责更清晰
2. **可维护性**: 减少单个包的复杂度，便于维护
3. **扩展性**: 未来可以更容易地添加新的文件处理功能

## 主要类型

### File
基础文件类型，包含以下属性：
- `URL`: 文件URL
- `Name`: 文件名
- `Size`: 文件大小
- `ContentType`: 文件MIME类型
- `Hash`: 文件哈希值
- `Description`: 文件描述
- `Metadata`: 文件元数据
- `LocalPath`: 本地路径
- `Downloaded`: 是否已下载
- `DownloadAt`: 下载时间
- `CreatedAt`: 创建时间
- `AutoDelete`: 是否自动删除
- `Compressed`: 是否压缩

### Reader
文件读取接口，提供以下功能：
- 从URL添加文件
- 下载文件到本地
- 按类型或扩展名过滤文件
- 获取文件信息
- 管理临时目录

### Writer
文件写入接口，提供以下功能：
- 添加本地文件
- 从数据创建文件
- 创建临时文件
- 设置文件集合摘要
- 设置文件生命周期
- 管理临时目录

## 生命周期

文件支持以下生命周期：
- `LifecycleTemporary`: 临时文件
- `LifecycleShortTerm`: 短期存储
- `LifecycleLongTerm`: 长期存储
- `LifecycleCache`: 缓存文件

## 使用示例

```go
import (
    "context"
    "github.com/your-org/your-repo/pkg/typex/files"
)

// 创建Reader
ctx := context.Background()
reader := files.NewReader(ctx)
defer reader.Cleanup()

// 从URL添加文件
err := reader.AddFileFromURL(
    "https://example.com/file.txt",
    "file.txt",
    files.WithContentType("text/plain"),
    files.WithDescription("示例文件"),
)

// 下载文件
localPath, err := reader.DownloadFile(0, "/path/to/save")

// 创建Writer
writer := files.NewWriter(ctx)
defer writer.Cleanup()

// 添加本地文件
err = writer.AddFile(
    localPath,
    files.WithContentType("text/plain"),
    files.WithDescription("处理后的文件"),
)

// 从数据创建文件
err = writer.AddFileWithData(
    "data.txt",
    []byte("Hello, World!"),
    files.WithContentType("text/plain"),
)

// 设置生命周期
writer.SetLifecycle(files.LifecycleShortTerm)
```

## 目录结构

```
pkg/typex/files/
├── README.md           # 本文档
├── types.go           # 基础类型定义
├── interfaces.go      # 接口定义
├── url_reader.go      # URL读取实现
├── cloud_writer.go    # 云存储写入实现
├── constructors.go    # 构造函数
├── interface_test.go  # 接口测试
├── url_reader_test.go # URL读取测试
└── cloud_writer_test.go # 云存储写入测试
```

## 向后兼容

为了保持向后兼容，在 `typex` 包中提供了迁移说明文件 `files_compat.go`，指导用户如何迁移到新的包结构。

## 临时目录管理

- 使用 `GetTempDir()` 获取临时目录
- 使用 `Cleanup()` 清理临时文件
- 支持自动创建和删除临时目录
- 临时目录格式：`./temp/{trace_id}/`

## 文件选项

支持以下文件选项：
- `WithSize`: 设置文件大小
- `WithContentType`: 设置MIME类型
- `WithHash`: 设置文件哈希
- `WithDescription`: 设置文件描述
- `WithMetadata`: 设置元数据
- `WithAutoDelete`: 设置自动删除
- `WithCompressed`: 设置压缩状态 