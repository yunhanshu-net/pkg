package dev

// Upload组件设计文档
// 文件上传相关的组件，支持多种文件类型和上传方式

// ===== 重要说明 =====
/*
以下参数是在外层通用结构中定义的，组件内部不需要重复定义：
- Type: string (固定类型，对应Go类型系统，如typex.File、typex.Files)
- Widget: string (组件类型，固定为"upload")
- Code: string (字段标识)
- Name: string (显示名称)
- Required: bool (是否必填)

FileWidget结构体只需要定义组件特有的配置参数。

文件类型设计：
- 单文件：typex.File (包含URL、文件名、大小、类型等信息)
- 多文件：typex.Files (File数组，支持集合操作)
- 传输方式：通过URL传递，避免二进制文件在链路中传输
- 框架封装：自动下载、处理、上传，对用户透明
*/

// ===== Upload组件定义 =====

// FileWidget 文件上传组件定义
type FileWidget struct {
	Accept     string `json:"accept"`      // 接受的文件类型：.jpg,.png,.pdf 等
	MaxSize    string `json:"max_size"`    // 最大文件大小：1MB、2GB等
	MaxCount   int    `json:"max_count"`   // 最大文件数量，默认1
	Preview    bool   `json:"preview"`     // 是否显示预览，默认false
	DragDrop   bool   `json:"drag_drop"`   // 是否支持拖拽上传，默认false
	UploadText string `json:"upload_text"` // 上传按钮文本
}

// 使用示例：
// Avatar string `runner:"code:avatar;name:头像;widget:upload;accept:.jpg,.png,.gif;max_size:2MB;upload_text:选择头像" validate:"required"`
// Documents []string `runner:"code:documents;name:文档;widget:upload;accept:.pdf,.doc,.docx;max_count:5;max_size:10MB"`

// ===== Upload组件设计 =====

/*
Upload组件用于处理文件上传，支持多种文件类型和上传场景。

设计目标：
1. 支持单文件和多文件上传
2. 支持文件类型限制和大小限制
3. 提供上传进度显示
4. 支持拖拽上传
5. 与现有验证规则无缝集成

使用场景：
- 头像上传：用户头像、商品图片等
- 文档上传：合同、证书、报告等
- 批量上传：相册、文档批量处理等
- 特定格式：Excel导入、PDF文档等
*/

// ===== 组件示例 =====

// UploadExample 文件上传组件示例
type UploadExample struct {
	// 单文件上传 - 头像
	Avatar string `runner:"code:avatar;name:头像;widget:upload;accept:image/*;max_size:2MB;max_count:1" validate:"required" json:"avatar"`
	// 注释：单个图片文件上传，限制2MB大小，必填验证

	// 多文件上传 - 相册
	Photos string `runner:"code:photos;name:相册;widget:upload;accept:image/jpeg,image/png;max_size:5MB;max_count:10" json:"photos"`
	// 注释：多个图片上传，支持JPEG和PNG格式，最多10个文件

	// 文档上传 - PDF
	Contract string `runner:"code:contract;name:合同文件;widget:upload;accept:.pdf;max_size:10MB;max_count:1" validate:"required" json:"contract"`
	// 注释：PDF文件上传，限制10MB大小，必填验证

	// Excel文件上传
	DataFile string `runner:"code:data_file;name:数据文件;widget:upload;accept:.xlsx,.xls;max_size:20MB;max_count:1" json:"data_file"`
	// 注释：Excel文件上传，支持新旧格式，限制20MB

	// 压缩包上传
	Archive string `runner:"code:archive;name:压缩包;widget:upload;accept:.zip,.rar,.7z;max_size:50MB;max_count:1" json:"archive"`
	// 注释：压缩包上传，支持多种压缩格式，限制50MB

	// 视频文件上传
	Video string `runner:"code:video;name:视频文件;widget:upload;accept:video/*;max_size:100MB;max_count:1" json:"video"`
	// 注释：视频文件上传，支持所有视频格式，限制100MB

	// 音频文件上传
	Audio string `runner:"code:audio;name:音频文件;widget:upload;accept:audio/*;max_size:20MB;max_count:5" json:"audio"`
	// 注释：音频文件上传，支持所有音频格式，最多5个文件

	// 任意文件上传
	Attachment string `runner:"code:attachment;name:附件;widget:upload;max_size:30MB;max_count:20" json:"attachment"`
	// 注释：不限制文件类型，限制大小和数量

	// 带预览的图片上传
	ProductImage string `runner:"code:product_image;name:商品图片;widget:upload;accept:image/*;max_size:3MB;max_count:6;preview:true" json:"product_image"`
	// 注释：商品图片上传，支持预览功能，最多6张图片

	// 拖拽上传区域
	BulkFiles string `runner:"code:bulk_files;name:批量文件;widget:upload;accept:*;max_size:10MB;max_count:50;drag_drop:true" json:"bulk_files"`
	// 注释：支持拖拽的批量文件上传，最多50个文件
}

// ===== 标签配置详解 =====

/*
Upload组件支持的标签：

核心标签：
- code: 字段代码（必需）
- name: 显示名称（必需）
- widget: upload（必需）

文件限制：
- accept: 允许的文件类型，支持MIME类型和扩展名
- max_size: 单个文件最大大小，如1MB、2GB等
- max_count: 最大文件数量，默认为1

功能配置：
- preview: 是否显示预览，默认false
- drag_drop: 是否支持拖拽上传，默认false
- multiple: 是否支持多文件选择，根据max_count自动判断

上传配置：
- upload_url: 自定义上传接口地址
- chunk_size: 分片上传大小，默认2MB
- auto_upload: 是否自动上传，默认true

显示控制：
- show: 显示场景控制
- hidden: 隐藏场景控制

常用accept值：
- image/*: 所有图片类型
- image/jpeg,image/png: 指定图片格式
- .pdf: PDF文件
- .xlsx,.xls: Excel文件
- video/*: 所有视频类型
- audio/*: 所有音频类型
- *: 所有文件类型
*/

// ===== 实现要点 =====

/*
前端实现要点：

1. 文件选择：
   - 支持点击选择和拖拽上传
   - 根据accept限制文件类型
   - 多文件选择支持

2. 文件验证：
   - 文件类型验证
   - 文件大小验证
   - 文件数量验证
   - 提供友好的错误提示

3. 上传处理：
   - 支持分片上传大文件
   - 显示上传进度
   - 支持暂停和恢复上传
   - 错误重试机制

4. 预览功能：
   - 图片预览
   - 文件信息显示
   - 删除已上传文件

5. 用户体验：
   - 拖拽区域高亮
   - 上传状态指示
   - 批量操作支持

后端实现要点：

1. 文件接收：
   - 处理multipart/form-data
   - 支持分片上传
   - 文件临时存储

2. 文件验证：
   - 服务端文件类型验证
   - 文件大小限制
   - 安全性检查

3. 文件存储：
   - 本地存储或云存储
   - 文件路径管理
   - 重复文件处理

4. 文件管理：
   - 文件元信息存储
   - 文件访问权限
   - 文件清理机制
*/

// ===== 验证规则集成 =====

/*
支持的验证规则：

1. 基础验证：
   - required: 必填验证
   - omitempty: 可选字段

2. 文件验证：
   - 自动根据accept进行文件类型验证
   - 根据max_size进行大小验证
   - 根据max_count进行数量验证

3. 自定义验证：
   - 可以添加自定义的文件验证规则
   - 支持业务逻辑验证

示例验证规则：
```go
// 必填的文件上传
Avatar string `validate:"required"`

// 文件类型和大小自动验证
Document string `runner:"accept:.pdf;max_size:10MB"`

// 多文件数量验证
Photos string `runner:"max_count:5"`
```
*/

// ===== 使用最佳实践 =====

/*
最佳实践建议：

1. 文件类型限制：
   - 根据业务需求严格限制文件类型
   - 使用具体的MIME类型而非通配符
   - 考虑安全性，避免可执行文件

2. 文件大小控制：
   - 设置合理的文件大小限制
   - 考虑网络传输和存储成本
   - 提供清晰的大小限制说明

3. 用户体验：
   - 提供拖拽上传功能
   - 显示上传进度和状态
   - 支持预览和删除操作

4. 性能优化：
   - 大文件使用分片上传
   - 图片自动压缩和缩略图
   - 异步处理和后台任务

5. 安全考虑：
   - 服务端验证文件类型
   - 文件内容安全扫描
   - 访问权限控制

示例：
```go
// 推荐的配置
type GoodExample struct {
    // 头像上传
    Avatar string `runner:"code:avatar;name:头像;widget:upload;accept:image/jpeg,image/png;max_size:2MB;max_count:1;preview:true"`

    // 文档上传
    Document string `runner:"code:document;name:文档;widget:upload;accept:.pdf,.doc,.docx;max_size:10MB;max_count:1"`

    // 批量图片上传
    Gallery string `runner:"code:gallery;name:相册;widget:upload;accept:image/*;max_size:5MB;max_count:10;drag_drop:true;preview:true"`
}

// 不推荐的配置
type BadExample struct {
    // 没有文件类型限制（安全风险）
    File1 string `runner:"code:file1;name:文件;widget:upload"` // 应该限制文件类型

    // 文件大小过大
    File2 string `runner:"code:file2;name:文件;widget:upload;max_size:1GB"` // 考虑网络传输

    // 文件数量过多
    File3 string `runner:"code:file3;name:文件;widget:upload;max_count:100"` // 考虑性能影响
}
```
*/

// ===== 实现优先级 =====

/*
实现步骤：

第一步：基础上传功能
- 单文件上传
- 文件类型和大小验证
- 基础UI界面

第二步：多文件上传
- 多文件选择
- 批量上传处理
- 上传进度显示

第三步：拖拽上传
- 拖拽区域实现
- 拖拽事件处理
- 视觉反馈

第四步：预览功能
- 图片预览
- 文件信息显示
- 删除功能

第五步：高级功能
- 分片上传
- 断点续传
- 图片压缩

预计总工期：4-5天
- 第一步：1天
- 第二步：1天
- 第三步：1天
- 第四步：1天
- 第五步：1-2天
*/
