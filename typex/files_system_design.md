# 文件系统设计文档

## 概述

本文件系统是一个基于context的多租户文件处理框架，专为function-go SDK设计。系统使用**trace_id**进行目录隔离，确保链路追踪的一致性，同时提供强大的文件上传、下载和处理能力。

## 核心特性

### 1. 链路追踪一致性
- **trace_id隔离**：使用请求的trace_id作为临时目录名，确保文件操作与请求链路完美关联
- **调试便利**：可以通过trace_id直接定位到对应的临时文件目录
- **日志关联**：文件操作日志与请求日志使用相同的trace_id，便于问题排查

### 2. 多租户安全
- **用户隔离**：云存储路径包含用户ID，确保多租户数据隔离
- **权限控制**：通过context传递用户信息，支持权限验证扩展
- **安全路径**：使用相对路径和UUID，避免路径遍历攻击

### 3. 自动化管理
- **自动清理**：请求结束后自动清理临时文件，防止磁盘空间泄漏
- **自动上传**：FilesWriter实现json.Marshaler接口，JSON序列化时自动上传文件
- **生命周期管理**：支持短期、临时、长期等不同的文件生命周期策略

## 架构设计

### 目录结构

```
./temp/
└── {trace_id}/                 # 基于trace_id的隔离目录
    ├── downloads/              # FilesReader下载的文件
    │   ├── input1.txt
    │   └── input2.pdf
    ├── uploads/                # FilesWriter准备上传的文件
    │   ├── result1.txt
    │   └── summary.json
    └── processing/             # 中间处理文件
        └── temp_data.csv
```

### 云存储路径

```
/{userID}/{runnerID}/{uniqueID}/{filename}
```

- `userID`: 用户标识，确保多租户隔离
- `runnerID`: 运行器标识，支持不同应用隔离
- `uniqueID`: 唯一标识符，防止文件名冲突
- `filename`: 原始文件名，保持用户友好性

## 核心组件

### 1. FilesWriter - 文件输出器

用于创建和上传文件到云存储。

```go
// 创建文件输出器（使用context中的trace_id）
writer := typex.NewFilesWriter(ctx)
defer writer.Cleanup() // 确保清理

// 添加文件数据
err := writer.AddFileWithData("result.txt", data, 
    typex.WithDescription("处理结果"))

// 设置生命周期
writer.SetLifecycle(typex.LifecycleShortTerm)

// JSON序列化时自动上传
jsonData, _ := json.Marshal(writer)
```

**特性**：
- 实现`json.Marshaler`接口，序列化时自动上传
- 支持多种文件生命周期策略
- 自动生成文件哈希和元数据
- 支持子目录组织

### 2. FilesReader - 文件输入器

用于下载和处理用户上传的文件。

```go
// 创建文件输入器（使用context中的trace_id）
reader := typex.NewFilesReader(ctx)
defer reader.Cleanup() // 确保清理

// 添加文件URL
reader.AddFileFromURL(url, name, 
    typex.WithSize(size),
    typex.WithContentType(contentType))

// 下载文件到本地
localPath, err := reader.DownloadFile(0)

// 批量下载
paths, err := reader.DownloadAll()
```

**特性**：
- 支持HTTP/HTTPS文件下载
- 自动文件类型检测
- 支持文件过滤和搜索
- 提供丰富的文件元数据

### 3. Context工具函数

从context中提取必要信息的工具函数。

```go
// 提取trace_id（优先级：context.Value > gin.Context.GetString > UUID fallback）
traceID := getTraceIDFromContext(ctx)

// 提取用户信息
userID := getUserFromContext(ctx)      // 默认: "anonymous"
runnerID := getRunnerFromContext(ctx)  // 默认: "default"
```

## 使用示例

### 基础文件处理

```go
func ProcessFiles(ctx *runner.Context, req *Request, resp response.Response) error {
    // 创建文件读取器（使用trace_id隔离）
    reader := typex.NewFilesReader(ctx)
    defer reader.Cleanup()
    
    // 添加输入文件
    for _, file := range req.Files {
        reader.AddFileFromURL(file.URL, file.Name,
            typex.WithSize(file.Size),
            typex.WithContentType(file.ContentType))
    }
    
    // 创建文件输出器（使用相同的trace_id）
    writer := typex.NewFilesWriter(ctx)
    defer writer.Cleanup()
    
    // 处理每个文件
    files := reader.GetFiles()
    for i, inputFile := range files {
        // 下载到本地（./temp/{trace_id}/downloads/）
        localPath, err := reader.DownloadFile(i)
        if err != nil {
            return err
        }
        
        // 处理文件内容
        content, err := os.ReadFile(localPath)
        if err != nil {
            return err
        }
        
        processedContent := processContent(content)
        
        // 添加处理结果（将上传到云存储）
        err = writer.AddFileWithData(
            fmt.Sprintf("processed_%s", inputFile.Name),
            processedContent,
            typex.WithDescription("处理后的文件"))
        if err != nil {
            return err
        }
    }
    
    // 返回响应（FilesWriter会自动上传）
    return resp.Form(&Response{
        Message: "处理完成",
        Files:   writer, // 自动序列化并上传
    }).Build()
}
```

### 高级文件组织

```go
// 使用子目录组织文件
tempDir, err := writer.GetTempDir("images", "processed")
filePath, err := writer.CreateTempFile("result.jpg", "images", "thumbnails")

// 文件过滤
textFiles := reader.FilterByExtension(".txt", ".md")
imageFiles := reader.FilterByType("image/")

// 批量操作
totalSize := reader.GetTotalSize()
allPaths, err := reader.DownloadAll("downloads", "batch1")
```

## 配置选项

### FilesWriter选项

```go
// 文件选项
typex.WithDescription("文件描述")
typex.WithMetadata("key", "value")
typex.WithAutoDelete(true)
typex.WithCompression(true)

// 生命周期选项
writer.SetLifecycle(typex.LifecycleShortTerm)   // 短期保存
writer.SetLifecycle(typex.LifecycleTemporary)   // 临时文件
writer.SetLifecycle(typex.LifecycleLongTerm)    // 长期保存
```

### FilesReader选项

```go
// 文件信息选项
typex.WithSize(1024)
typex.WithContentType("text/plain")
typex.WithHash("sha256hash")
typex.WithFileDescription("输入文件")
typex.WithFileMetadata("source", "upload")
```

## 最佳实践

### 1. 资源管理

```go
// ✅ 正确：使用defer确保清理
writer := typex.NewFilesWriter(ctx)
defer writer.Cleanup()

// ❌ 错误：忘记清理资源
writer := typex.NewFilesWriter(ctx)
// 没有调用Cleanup()
```

### 2. 错误处理

```go
// ✅ 正确：检查所有错误
localPath, err := reader.DownloadFile(0)
if err != nil {
    return fmt.Errorf("下载文件失败: %v", err)
}

// ❌ 错误：忽略错误
localPath, _ := reader.DownloadFile(0)
```

### 3. 目录组织

```go
// ✅ 正确：使用有意义的子目录
tempDir, err := writer.GetTempDir("images", "processed")
filePath, err := writer.CreateTempFile("thumb.jpg", "images", "thumbnails")

// ❌ 错误：所有文件放在根目录
tempDir, err := writer.GetTempDir()
```

### 4. 生命周期管理

```go
// ✅ 正确：根据用途设置合适的生命周期
writer.SetLifecycle(typex.LifecycleShortTerm)   // 临时处理结果
writer.SetLifecycle(typex.LifecycleLongTerm)    // 重要报告文件

// ❌ 错误：不设置生命周期（使用默认值）
```

## 链路追踪优势

### 1. 调试便利性

```bash
# 通过trace_id直接定位文件
ls ./temp/trace-abc123/
# downloads/  uploads/  processing/

# 日志中的trace_id与文件目录一致
[INFO] trace-abc123: 开始处理文件
[INFO] trace-abc123: 文件下载到 ./temp/trace-abc123/downloads/input.txt
[INFO] trace-abc123: 文件上传完成
```

### 2. 问题排查

```go
// 日志中包含trace_id和文件路径
logger.InfoContextf(ctx, "文件处理完成，临时目录: %s", tempDir)
// 输出: [INFO] trace-abc123: 文件处理完成，临时目录: ./temp/trace-abc123
```

### 3. 监控集成

```go
// 可以基于trace_id进行文件操作监控
metrics.RecordFileOperation(ctx.Value(constants.TraceID), "upload", fileSize)
```

## 安全考虑

### 1. 路径安全
- 使用相对路径避免路径遍历
- trace_id作为目录名，确保隔离
- 自动清理防止文件泄漏

### 2. 多租户隔离
- 云存储路径包含用户ID
- Context传递用户信息
- 权限验证扩展点

### 3. 资源限制
- 自动清理临时文件
- 文件大小限制（可配置）
- 并发下载控制

## 扩展性

### 1. 存储后端
- 当前支持七牛云
- 接口设计支持多云存储
- 可扩展本地存储、AWS S3等

### 2. 文件处理
- 支持自定义文件处理器
- 可扩展文件格式支持
- 支持异步处理

### 3. 监控集成
- 文件操作指标收集
- 性能监控
- 错误追踪

## 性能优化

### 1. 并发处理
```go
// 并发下载多个文件
var wg sync.WaitGroup
for i := range files {
    wg.Add(1)
    go func(index int) {
        defer wg.Done()
        reader.DownloadFile(index)
    }(i)
}
wg.Wait()
```

### 2. 内存优化
- 流式文件处理
- 分块上传大文件
- 及时释放资源

### 3. 网络优化
- 连接池复用
- 断点续传支持
- 压缩传输

## 总结

本文件系统通过使用trace_id进行目录隔离，实现了：

1. **链路追踪一致性** - 文件操作与请求链路完美关联
2. **多租户安全** - 用户数据完全隔离
3. **自动化管理** - 资源自动清理，操作自动化
4. **开发友好** - 简单易用的API，丰富的功能
5. **生产就绪** - 完善的错误处理、监控和扩展性

这种设计特别适合微服务架构中的文件处理需求，既保证了安全性和可靠性，又提供了出色的开发体验和运维便利性。 