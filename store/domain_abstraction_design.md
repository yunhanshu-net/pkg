# Store域名抽象设计方案

## 📋 需求分析

### 🎯 核心需求
1. **域名隐藏**：将CDN域名从代码中抽离，支持动态配置
2. **向后兼容**：不破坏现有的store接口和实现
3. **typex.File集成**：与新的File类型系统无缝集成
4. **多环境支持**：开发、测试、生产环境使用不同域名

### 🔍 现有问题
- 域名硬编码：`Domain: "http://cdn.geeleo.com"`
- 路径混乱：SavePath和SaveFullPath概念重叠
- 缺乏统一的文件标识符

## 🔧 技术方案

### 方案一：文件ID抽象（推荐）

#### 核心思路
使用**文件ID**替代完整URL，通过ID映射到实际存储路径：

```go
// 文件ID格式：store://bucket/path/file.ext
// 示例：store://geeleo/images/avatar/user123.jpg

type FileID string

func (id FileID) ToURL(resolver URLResolver) string {
    return resolver.Resolve(id)
}

func (id FileID) ToPath() string {
    // store://geeleo/images/avatar/user123.jpg -> /images/avatar/user123.jpg
    return extractPath(id)
}
```

#### 接口设计
```go
// URLResolver 域名解析器接口
type URLResolver interface {
    Resolve(fileID FileID) string
    ParseURL(url string) (FileID, error)
}

// 更新后的FileStore接口
type FileStore interface {
    FileSave(localFilePath string, savePath string) (*FileSaveResult, error)
    GetFile(fileID FileID) (*GetFileResult, error)
    DeleteFile(fileID FileID) error
    GetResolver() URLResolver
}

type FileSaveResult struct {
    FileID   FileID `json:"file_id"`   // store://geeleo/path/file.ext
    FileName string `json:"file_name"` // 原始文件名
    FileType string `json:"file_type"` // 文件类型
    FileSize int64  `json:"file_size"` // 文件大小
}

type GetFileResult struct {
    FileSaveResult
    LocalPath string `json:"local_path"` // 下载后的本地路径
}
```

### 方案二：配置化域名管理

#### 配置结构
```go
type StoreConfig struct {
    Provider string            `json:"provider"` // qiniu, aliyun, aws
    Domains  map[string]string `json:"domains"`  // 环境域名映射
    Default  string            `json:"default"`  // 默认环境
}

// 配置示例
{
    "provider": "qiniu",
    "domains": {
        "dev": "http://dev-cdn.geeleo.com",
        "test": "http://test-cdn.geeleo.com", 
        "prod": "http://cdn.geeleo.com"
    },
    "default": "prod"
}
```

#### 域名解析器实现
```go
type ConfigURLResolver struct {
    config *StoreConfig
    env    string
}

func (r *ConfigURLResolver) Resolve(fileID FileID) string {
    domain := r.getDomain()
    path := fileID.ToPath()
    return domain + path
}

func (r *ConfigURLResolver) getDomain() string {
    if domain, ok := r.config.Domains[r.env]; ok {
        return domain
    }
    return r.config.Domains[r.config.Default]
}
```

## 🔗 与typex.File集成

### 集成设计
```go
// 扩展typex.File支持FileID
type File struct {
    ID          FileID    `json:"id"`           // 文件标识符
    URL         string    `json:"url"`          // 完整URL（动态生成）
    Name        string    `json:"name"`         // 文件名
    Size        int64     `json:"size"`         // 文件大小
    ContentType string    `json:"content_type"` // MIME类型
    UploadTime  time.Time `json:"upload_time"`  // 上传时间
}

// 实现FileUploader接口
type StoreFileUploader struct {
    store    FileStore
    resolver URLResolver
}

func (u *StoreFileUploader) Upload(data []byte, filename, contentType string) (string, error) {
    // 1. 保存到临时文件
    tempFile := saveTempFile(data, filename)
    defer os.Remove(tempFile)
    
    // 2. 上传到存储
    result, err := u.store.FileSave(tempFile, generateSavePath(filename))
    if err != nil {
        return "", err
    }
    
    // 3. 返回FileID（而非完整URL）
    return string(result.FileID), nil
}

// File类型的URL动态生成
func (f *File) GetURL(resolver URLResolver) string {
    if f.ID != "" {
        return resolver.Resolve(f.ID)
    }
    return f.URL // 兼容旧格式
}
```

### JSON序列化优化
```go
func (f *File) MarshalJSON() ([]byte, error) {
    type fileAlias File
    
    // 如果有resolver，动态生成URL
    if f.ID != "" && globalResolver != nil {
        f.URL = globalResolver.Resolve(f.ID)
    }
    
    return json.Marshal((*fileAlias)(f))
}

func (f *File) UnmarshalJSON(b []byte) error {
    // 支持多种格式
    if isSimpleURL(b) {
        // 简单URL格式："http://cdn.geeleo.com/path/file.ext"
        var url string
        json.Unmarshal(b, &url)
        f.URL = url
        f.ID = parseURLToFileID(url) // 尝试解析为FileID
    } else {
        // 完整对象格式
        type fileAlias File
        json.Unmarshal(b, (*fileAlias)(f))
    }
    
    return nil
}
```

## 🛠️ 实现策略

### 阶段一：向后兼容的域名抽象
1. **保持现有接口不变**
2. **添加URLResolver支持**
3. **配置化域名管理**
4. **渐进式迁移**

```go
// 兼容性包装
type LegacyQiNiu struct {
    *QiNiu
    resolver URLResolver
}

func (q *LegacyQiNiu) FileSave(localFilePath string, savePath string) (*FileSaveInfo, error) {
    // 调用新接口
    result, err := q.QiNiu.FileSaveV2(localFilePath, savePath)
    if err != nil {
        return nil, err
    }
    
    // 转换为旧格式
    return &FileSaveInfo{
        SavePath:     result.FileID.ToPath(),
        SaveFullPath: q.resolver.Resolve(result.FileID),
        FileName:     result.FileName,
        FileType:     result.FileType,
    }, nil
}
```

### 阶段二：typex.File集成
1. **实现FileUploader接口**
2. **File类型支持FileID**
3. **自动URL解析**
4. **统一文件操作API**

### 阶段三：完整迁移
1. **废弃旧接口**
2. **统一使用FileID**
3. **优化性能和缓存**

## 📊 技术对比

| 方案 | 优势 | 劣势 | 适用场景 |
|------|------|------|----------|
| FileID抽象 | 完全隐藏域名，支持多存储 | 需要重构现有代码 | 新项目，长期规划 |
| 配置化域名 | 改动最小，快速实现 | 仍然暴露URL结构 | 现有项目，快速迁移 |
| 混合方案 | 兼容性好，渐进迁移 | 复杂度较高 | 大型项目，平滑过渡 |

## 🎯 推荐方案

**推荐使用混合方案**，分阶段实施：

1. **短期**：配置化域名管理，快速解决域名硬编码问题
2. **中期**：FileID抽象，与typex.File集成
3. **长期**：完全迁移到FileID体系，废弃URL传递

### 核心优势
- ✅ 向后兼容，不破坏现有代码
- ✅ 渐进式迁移，风险可控
- ✅ 最终实现完全的域名抽象
- ✅ 与typex.File无缝集成

### 实现优先级
1. **P0**：配置化域名管理（1天）
2. **P1**：URLResolver接口设计（1天）
3. **P2**：typex.File集成（2天）
4. **P3**：FileID抽象实现（2天）
5. **P4**：完整迁移和优化（3天）

这个方案既解决了域名硬编码问题，又为未来的扩展和优化奠定了基础。 