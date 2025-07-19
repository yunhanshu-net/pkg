# Files API 重构说明

## 概述

本次重构统一了文件处理接口，将原本分离的 `files.Files`（输入）和 `files.Writer`（输出）合并为统一的 `files.Files` 类型，简化了API设计，提高了易用性。

## 核心设计理念

### 1. **概念简化**
- **统一类型**: 输入和输出都使用 `files.Files` 类型
- **自动化处理**: 程序内部智能判断和处理，无需暴露复杂配置
- **无声操作**: 下载、处理、上传都在后台自动完成

### 2. **立即上传机制**
- **添加时上传**: 调用 `AddFileFromData()` 或 `AddFileFromPath()` 时立即上传
- **避免延迟**: 不依赖JSON序列化时上传，避免时机问题
- **自动清理**: 上传完成后自动清理本地临时文件

### 3. **生命周期管理**
- **简化设计**: 可以为nil表示无限制
- **灵活配置**: 支持下载次数限制和过期时间
- **便捷方法**: 提供常用生命周期的快捷设置

## 主要功能

### 1. **文件添加方法**

```go
// 从数据创建文件并立即上传
func (fc *Files) AddFileFromData(filename string, data []byte, options ...FileOption) error

// 从本地路径添加文件并立即上传  
func (fc *Files) AddFileFromPath(localPath string, options ...FileOption) error

// 添加已存在的File对象（不上传）
func (fc *Files) AddExistingFile(file *File) *Files

// 从URL添加文件引用（不下载）
func (fc *Files) AddFileFromURL(url, filename string, options ...FileOption) error
```

### 2. **生命周期管理**

```go
// 设置下载次数限制
func (fc *Files) SetMaxDownloads(maxDownloads int) *Files

// 设置过期时间
func (fc *Files) SetExpiresAt(expiresAt time.Time) *Files

// 便捷方法
func (fc *Files) SetTemporary() *Files        // 下载一次后删除
func (fc *Files) SetExpiring7Days() *Files    // 7天后过期
func (fc *Files) SetUnlimited() *Files        // 无限制
```

### 3. **按需下载**

```go
// File 级别的按需下载
func (f *File) GetLocalPath(ctx context.Context) (string, error)
```

### 4. **元数据和配置**

```go
// 链式调用设置
func (fc *Files) SetNote(note string) *Files
func (fc *Files) SetMetadata(key string, value interface{}) *Files
func (fc *Files) SetConfig(key string, value interface{}) *Files
```

### 5. **查询和统计**

```go
func (fc *Files) GetFileCount() int
func (fc *Files) GetTotalSize() int64
func (fc *Files) FilterByType(contentType string) *Files
func (fc *Files) FilterByExtension(ext string) *Files
```

### 6. **清理管理**

```go
// 清理本地缓存文件
func (fc *Files) CleanupLocalFiles() error
```

### 7. **Context便捷方法**

在 `runner.Context` 上提供了便捷的创建方法，自动设置context：

```go
// 基础方法：创建文件集合并自动设置context
func (c *Context) NewFiles(input interface{}) *files.Files

// 便捷方法：创建临时文件集合（下载一次后删除）
func (c *Context) NewTemporaryFiles() *files.Files

// 便捷方法：创建有效期文件集合（7天后过期）
func (c *Context) NewExpiringFiles() *files.Files

// 便捷方法：创建永久文件集合（无限制）
func (c *Context) NewPermanentFiles() *files.Files
```

## Context便捷方法使用示例

```go
func ProcessFiles(ctx *runner.Context, req *ProcessFilesReq, resp response.Response) error {
    // 根据业务需求选择合适的便捷方法
    var outputFiles *files.Files
    
    switch req.OutputType {
    case "temporary":
        outputFiles = ctx.NewTemporaryFiles() // 临时文件，下载一次后删除
    case "expiring":
        outputFiles = ctx.NewExpiringFiles()  // 7天后过期
    case "permanent":
        outputFiles = ctx.NewPermanentFiles() // 永久保存
    default:
        outputFiles = ctx.NewFiles([]string{}) // 默认无限制
    }
    
    // 设置文件集合的元数据
    outputFiles.SetNote("处理结果文件").
        SetMetadata("process_type", req.ProcessType)
    
    // 处理文件...
    for _, file := range req.InputFiles.Files {
        localPath, err := file.GetLocalPath(ctx)
        if err != nil {
            continue
        }
        
        processedData := processFile(localPath)
        outputFiles.AddFileFromData("processed_"+file.Name, processedData)
    }
    
    return resp.Form(&ProcessFilesResp{
        OutputFiles: outputFiles,
    }).Build()
}
```

## 使用场景

### 场景1：仅存储
```go
// 用户上传文件 -> 直接存储（文件已经在OSS上）
input: files.Files{Files: []*File{{URL: "https://oss.com/file1.pdf"}}}
output: files.Files{Files: []*File{{URL: "https://oss.com/file1.pdf"}}} // 原样返回
```

### 场景2：处理后返回
```go
// 用户上传文件 -> 下载到本地 -> 处理 -> 生成新文件 -> 上传 -> 返回
input: files.Files{Files: []*File{{URL: "https://oss.com/input.pdf"}}}

// 函数内部处理
for _, file := range req.Files.Files {
    localPath, err := file.GetLocalPath(ctx)  // 按需下载
    processedData := processFile(localPath)   // 处理文件
    
    // 添加处理后的文件（立即上传）
    outputFiles.AddFileFromData("processed.txt", processedData)
}
```

### 场景3：混合处理
```go
// 部分文件处理，部分直接存储
for _, file := range req.Files.Files {
    if shouldProcess(file) {
        // 处理文件
        localPath, err := file.GetLocalPath(ctx)
        processedData := processFile(localPath)
        outputFiles.AddFileFromData("processed_"+file.Name, processedData)
    } else {
        // 直接引用
        outputFiles.AddFileFromURL(file.URL, file.Name)
    }
}
```

## 标准示例

```go
type ProcessFilesReq struct {
    Files *files.Files `json:"files" runner:"code:files;name:输入文件" widget:"type:file_upload;multiple:true" data:"type:files" validate:"required"`
    ProcessType string `json:"process_type" form:"process_type" runner:"code:process_type;name:处理类型" widget:"type:select;options:仅存储,文本处理,图片处理" data:"type:string;default_value:仅存储" validate:"required"`
}

type ProcessFilesResp struct {
    ProcessedFiles *files.Files `json:"processed_files" runner:"code:processed_files;name:处理后的文件" widget:"type:file_display;display_mode:card" data:"type:files"`
    Summary string `json:"summary" runner:"code:summary;name:处理摘要" widget:"type:input;mode:text_area;readonly:true" data:"type:string"`
}

func ProcessFiles(ctx *runner.Context, req *ProcessFilesReq, resp response.Response) error {
    // 创建输出文件集合，设置为临时文件（使用Context便捷方法）
    outputFiles := ctx.NewTemporaryFiles().
        SetNote("文件处理结果")
    
    // 遍历输入文件，选择性处理
    for _, file := range req.Files.Files {
        switch req.ProcessType {
        case "仅存储":
            // 直接引用，不下载处理
            outputFiles.AddFileFromURL(file.URL, file.Name)
            
        case "文本处理":
            // 按需下载处理
            localPath, err := file.GetLocalPath(ctx)
            if err != nil {
                continue
            }
            
            processedData := processTextFile(localPath)
            outputFiles.AddFileFromData("processed_"+file.Name, processedData)
            
        case "图片处理":
            // 类似处理...
        }
    }
    
    // 清理本地缓存文件
    defer func() {
        req.Files.CleanupLocalFiles()
        outputFiles.CleanupLocalFiles()
    }()
    
    return resp.Form(&ProcessFilesResp{
        ProcessedFiles: outputFiles,
        Summary: "处理完成",
    }).Build()
}
```

## 核心优势

### 1. **简化概念**
- 统一使用 `files.Files` 类型，大模型不会混淆
- 自动化处理，减少配置复杂度

### 2. **灵活处理**
- 支持按需下载，不强制下载所有文件
- 支持并发处理（业务层自己实现）
- 支持部分处理、部分存储

### 3. **生命周期管理**
- 简化的生命周期设计
- 支持临时文件、有效期文件、永久文件
- 便捷的设置方法

### 4. **立即上传**
- 添加文件时立即上传，避免序列化时机问题
- 上传完成后立即清理本地文件
- 支持自动清理机制

### 5. **兼容性**
- 提供向后兼容的方法
- 支持数据库序列化（Scanner/Valuer接口）
- 支持JSON序列化

## 数据库支持

`Files` 类型实现了 `sql.Scanner` 和 `driver.Valuer` 接口，可以直接存储到数据库：

```go
type MyModel struct {
    ID    int          `gorm:"primaryKey"`
    Files *files.Files `gorm:"type:json;comment:文件列表"`
}
```

## 注意事项

1. **临时目录**: 使用项目工作目录下的 `./temp/traceID` 而不是系统临时目录
2. **上传时机**: 文件添加时立即上传，不依赖JSON序列化
3. **清理机制**: 需要手动调用 `CleanupLocalFiles()` 清理本地缓存
4. **并发控制**: 业务层自行控制并发，避免带宽问题
5. **错误处理**: 使用 `github.com/pkg/errors` 包装错误，在API层统一处理

## 迁移指南

### 从旧API迁移

```go
// 旧API
writer := files.NewWriter(ctx)
writer.AddFileWithData("file.txt", data)
return resp.Form(&Resp{Files: writer}).Build()

// 新API（方式1：直接创建）
outputFiles := files.NewFiles([]string{}).SetContext(ctx)
outputFiles.AddFileFromData("file.txt", data)
return resp.Form(&Resp{Files: outputFiles}).Build()

// 新API（方式2：使用Context便捷方法，推荐）
outputFiles := ctx.NewTemporaryFiles() // 或 NewPermanentFiles()、NewExpiringFiles()
outputFiles.AddFileFromData("file.txt", data)
return resp.Form(&Resp{Files: outputFiles}).Build()
```

### 请求响应结构体

```go
// 统一使用 *files.Files 类型
type Req struct {
    Files *files.Files `json:"files" data:"type:files"`
}

type Resp struct {
    Files *files.Files `json:"files" data:"type:files"`
}
```

## JSON示例

### 完整的Files JSON结构

```json
{
  "files": [
    {
      "name": "document.pdf",
      "size": 2048576,
      "content_type": "application/pdf",
      "url": "https://oss.example.com/files/document_20250115.pdf",
      "preview_url": "https://oss.example.com/preview/document_20250115.jpg",
      "created_at": "2025-01-15T09:30:00Z",
      "updated_at": "2025-01-15T09:30:15Z",
      "status": "uploaded",
      "metadata": {
        "original_name": "用户文档.pdf",
        "process_type": "text_extraction",
        "upload_user": "user123"
      },
      "hash": "sha256:abc123def456...",
      "description": "用户上传的PDF文档",
      "lifecycle": {
        "max_downloads": 5,
        "expires_at": "2025-01-22T09:30:00Z",
        "downloads": 0
      }
    },
    {
      "name": "image.jpg",
      "size": 1024000,
      "content_type": "image/jpeg",
      "url": "https://oss.example.com/files/image_20250115.jpg",
      "preview_url": "https://oss.example.com/preview/image_20250115_thumb.jpg",
      "created_at": "2025-01-15T09:31:00Z",
      "updated_at": "2025-01-15T09:31:10Z",
      "status": "uploaded",
      "metadata": {
        "original_name": "照片.jpg",
        "process_type": "image_analysis",
        "width": "1920",
        "height": "1080"
      },
      "hash": "sha256:def456ghi789...",
      "description": "用户上传的图片文件",
      "lifecycle": null
    }
  ],
  "options": {
    "process_mode": "auto_extract",
    "quality": 80,
    "compress": true,
    "backup_original": true
  },
  "note": "文件处理完成 - 共处理2个文件",
  "config": {
    "upload_provider": "qiniu",
    "storage_region": "cn-east-1",
    "auto_cleanup": true
  },
  "created_at": "2025-01-15T09:30:00Z",
  "default_lifecycle": {
    "max_downloads": 10,
    "expires_at": "2025-01-22T09:30:00Z",
    "downloads": 0
  }
}
```

### 简化的Files JSON结构（仅存储）

```json
{
  "files": [
    {
      "name": "simple.txt",
      "size": 1024,
      "content_type": "text/plain",
      "url": "https://oss.example.com/files/simple.txt",
      "created_at": "2025-01-15T09:30:00Z",
      "status": "referenced",
      "metadata": {}
    }
  ],
  "options": {},
  "note": "",
  "config": {},
  "created_at": "2025-01-15T09:30:00Z",
  "default_lifecycle": null
}
```

### 临时文件的JSON结构

```json
{
  "files": [
    {
      "name": "temp_report.txt",
      "size": 2048,
      "content_type": "text/plain",
      "url": "https://oss.example.com/temp/temp_report_20250115.txt",
      "created_at": "2025-01-15T09:30:00Z",
      "updated_at": "2025-01-15T09:30:05Z",
      "status": "uploaded",
      "metadata": {
        "type": "report",
        "auto_generated": "true"
      },
      "lifecycle": {
        "max_downloads": 1,
        "expires_at": null,
        "downloads": 0
      }
    }
  ],
  "options": {
    "temp_file": true
  },
  "note": "临时处理结果 - 下载一次后删除",
  "config": {},
  "created_at": "2025-01-15T09:30:00Z",
  "default_lifecycle": {
    "max_downloads": 1,
    "expires_at": null,
    "downloads": 0
  }
}
```

### 有效期文件的JSON结构

```json
{
  "files": [
    {
      "name": "expiring_data.csv",
      "size": 4096,
      "content_type": "text/csv",
      "url": "https://oss.example.com/expiring/data_20250115.csv",
      "created_at": "2025-01-15T09:30:00Z",
      "status": "uploaded",
      "metadata": {
        "export_type": "user_data",
        "records_count": "150"
      },
      "lifecycle": {
        "max_downloads": null,
        "expires_at": "2025-01-22T09:30:00Z",
        "downloads": 0
      }
    }
  ],
  "options": {},
  "note": "数据导出结果 - 7天后过期",
  "config": {},
  "created_at": "2025-01-15T09:30:00Z",
  "default_lifecycle": {
    "max_downloads": null,
    "expires_at": "2025-01-22T09:30:00Z",
    "downloads": 0
  }
}
```

### 字段说明

#### Files 对象字段
- `files`: 文件列表数组
- `options`: 用户选择的参数/选项
- `note`: 用户备注/说明
- `config`: 处理配置
- `created_at`: 创建时间
- `default_lifecycle`: 默认生命周期（可为null表示无限制）

#### File 对象字段
- `name`: 文件名
- `size`: 文件大小（字节）
- `content_type`: 文件MIME类型
- `url`: 文件访问URL
- `preview_url`: 预览URL（可选）
- `created_at`: 创建时间
- `updated_at`: 更新时间（可选）
- `status`: 文件状态（pending/uploading/uploaded/failed/referenced）
- `metadata`: 元数据键值对
- `hash`: 文件哈希值（可选）
- `description`: 文件描述（可选）
- `lifecycle`: 文件生命周期（可为null表示无限制）

#### FileLifecycle 对象字段
- `max_downloads`: 最大下载次数（null表示无限制）
- `expires_at`: 过期时间（null表示永不过期）
- `downloads`: 当前下载次数（服务端维护）

## 数据库支持 