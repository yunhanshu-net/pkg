# 文件类型设计分析

## 📋 技术架构现状

### 🎯 链路限制
```
用户表单 -> function-server -> function-runtime -> function-go
```

**核心问题**：多层链路传递二进制文件会导致：
- 性能问题：文件在每层都需要序列化/反序列化
- 内存问题：大文件会占用大量内存
- 复杂度问题：错误处理、超时控制等复杂

**解决方案**：URL传递 + 框架封装
- 前端上传文件到OSS，获得URL
- 链路中只传递URL字符串
- function-go内部自动下载/上传文件
- 对用户完全透明

## 🔧 类型设计方案

### 方案一：统一Files类型（推荐）

```go
// 统一使用Files类型，支持单文件和多文件
type ProcessImageRequest struct {
    InputImages  typex.Files `runner:"code:input_images;name:输入图片;widget:upload;accept:image/*" json:"input_images"`
    OutputFormat string      `runner:"code:output_format;name:输出格式;widget:select;options:jpg,png,webp" json:"output_format"`
}

type ProcessImageResponse struct {
    OutputImages typex.Files `json:"output_images"` // 处理后的图片
    ProcessInfo  string      `json:"process_info"`  // 处理信息
}
```

**优势**：
- ✅ 请求响应类型一致
- ✅ 单文件多文件统一处理
- ✅ 丰富的文件元信息
- ✅ 支持集合操作（过滤、分组等）
- ✅ 扩展性强

### 方案二：File/Files分离

```go
// 根据业务场景选择File或Files
type ProcessImageRequest struct {
    InputImage   typex.File `runner:"code:input_image;name:输入图片;widget:upload;accept:image/*;max_count:1" json:"input_image"`
    OutputFormat string     `runner:"code:output_format;name:输出格式;widget:select;options:jpg,png,webp" json:"output_format"`
}

type ProcessImageResponse struct {
    OutputImage typex.File `json:"output_image"` // 处理后的图片
    ProcessInfo string     `json:"process_info"` // 处理信息
}
```

**优势**：
- ✅ 语义更明确
- ✅ 类型安全
- ❌ 需要维护两套类型
- ❌ 单文件转多文件需要重构

## 💡 推荐设计：统一Files类型

### 核心理由

1. **一致性**：请求和响应使用相同类型，减少认知负担
2. **扩展性**：业务从单文件扩展到多文件无需重构
3. **丰富性**：Files类型包含完整的文件元信息
4. **操作性**：支持丰富的集合操作

### 使用示例

```go
// 图片处理函数
func ProcessImages(ctx context.Context, req *ProcessImageRequest) (*ProcessImageResponse, error) {
    // 1. 获取输入文件（框架自动下载）
    inputFiles := req.InputImages
    
    // 2. 处理每个文件
    var outputFiles typex.Files
    for _, inputFile := range inputFiles {
        // 下载文件内容
        data, err := inputFile.Download()
        if err != nil {
            return nil, err
        }
        
        // 处理图片（转换格式）
        processedData, err := convertImage(data, req.OutputFormat)
        if err != nil {
            return nil, err
        }
        
        // 创建输出文件
        outputFile := typex.NewFile("", generateOutputName(inputFile.Name, req.OutputFormat))
        
        // 上传处理后的文件（框架自动上传）
        err = outputFile.Upload(processedData, getUploader())
        if err != nil {
            return nil, err
        }
        
        outputFiles.Add(*outputFile)
    }
    
    return &ProcessImageResponse{
        OutputImages: outputFiles,
        ProcessInfo:  fmt.Sprintf("处理了%d个文件", len(inputFiles)),
    }, nil
}

// 单文件场景的便利方法
func ProcessSingleImage(ctx context.Context, req *ProcessImageRequest) (*ProcessImageResponse, error) {
    if len(req.InputImages) == 0 {
        return nil, fmt.Errorf("没有输入文件")
    }
    
    // 只处理第一个文件
    inputFile := req.InputImages.First()
    // ... 处理逻辑
}
```

## 🎨 前端集成设计

### Upload组件配置

```go
// 单文件上传
Avatar typex.Files `runner:"code:avatar;name:头像;widget:upload;accept:image/*;max_count:1;preview:true"`

// 多文件上传
Gallery typex.Files `runner:"code:gallery;name:相册;widget:upload;accept:image/*;max_count:10;preview:true"`

// 文档上传
Documents typex.Files `runner:"code:documents;name:文档;widget:upload;accept:.pdf,.doc,.docx;max_count:5"`
```

### 前端渲染逻辑

```javascript
// 前端根据max_count判断单文件还是多文件
if (field.max_count === 1) {
    // 渲染单文件上传组件
    return <SingleFileUpload {...props} />
} else {
    // 渲染多文件上传组件
    return <MultiFileUpload {...props} />
}

// 数据格式统一
const fileData = {
    url: "https://oss.example.com/file.jpg",
    name: "image.jpg",
    size: 1024000,
    content_type: "image/jpeg",
    upload_time: "2025-01-13T10:30:00Z"
}

// 单文件：[fileData]
// 多文件：[fileData1, fileData2, ...]
```

## 🔄 数据流设计

### 请求流程
```
1. 用户选择文件 -> 前端上传到OSS -> 获得URL
2. 前端构造Files数组 -> 发送请求
3. function-go接收Files -> 框架自动下载文件
4. 用户代码处理文件 -> 生成新文件
5. 框架自动上传新文件 -> 返回Files数组
```

### 数据格式
```json
// 请求参数
{
    "input_images": [
        {
            "url": "https://oss.example.com/input.png",
            "name": "input.png",
            "size": 1024000,
            "content_type": "image/png"
        }
    ],
    "output_format": "jpg"
}

// 响应参数
{
    "output_images": [
        {
            "url": "https://oss.example.com/output.jpg",
            "name": "output.jpg", 
            "size": 856000,
            "content_type": "image/jpeg",
            "upload_time": "2025-01-13T10:35:00Z"
        }
    ],
    "process_info": "处理了1个文件"
}
```

## 🛠️ 框架封装设计

### 自动下载机制
```go
// 框架在调用用户函数前自动执行
func (f *File) ensureDownloaded() error {
    if f.localPath == "" {
        data, err := f.Download()
        if err != nil {
            return err
        }
        f.localPath = saveToTemp(data)
    }
    return nil
}
```

### 自动上传机制
```go
// 框架在用户函数返回后自动执行
func (f *File) ensureUploaded() error {
    if f.URL == "" && f.localPath != "" {
        data, err := os.ReadFile(f.localPath)
        if err != nil {
            return err
        }
        return f.Upload(data, getDefaultUploader())
    }
    return nil
}
```

## 📊 性能考虑

### 内存优化
- 大文件使用流式处理
- 及时清理临时文件
- 支持分片上传/下载

### 并发处理
- 多文件并发下载/上传
- 限制并发数避免资源耗尽
- 支持超时控制

### 缓存策略
- 相同URL的文件缓存
- 基于文件哈希的去重
- 临时文件自动清理

## 🎯 最终建议

**推荐使用统一的typex.Files类型**，理由：

1. **架构一致性**：请求响应类型统一，降低复杂度
2. **业务扩展性**：从单文件到多文件无缝扩展
3. **开发体验**：丰富的API和集合操作
4. **维护成本**：只需维护一套类型系统

**实现优先级**：
1. 实现typex.File和typex.Files基础类型
2. 实现Upload组件的Files支持
3. 实现框架的自动下载/上传机制
4. 优化性能和错误处理
5. 完善文档和示例

这样的设计既解决了技术架构的限制，又提供了良好的开发体验和扩展性。 