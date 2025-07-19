package files

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/yunhanshu-net/pkg/trace"
	"github.com/yunhanshu-net/pkg/typex"
)

// FileLifecycle 文件生命周期（可以为nil表示无限制）
type FileLifecycle struct {
	MaxDownloads *int       `json:"max_downloads,omitempty"` // 最大下载次数，nil表示无限制
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`    // 过期时间，nil表示永不过期
	Downloads    int        `json:"downloads,omitempty"`     // 当前下载次数（服务端维护）
}

// IsExpired 判断是否已过期或超出下载限制
func (lc *FileLifecycle) IsExpired() bool {
	if lc == nil {
		return false // nil表示无限制
	}

	// 检查时间过期
	if lc.ExpiresAt != nil && time.Now().After(*lc.ExpiresAt) {
		return true
	}

	// 检查下载次数超限
	if lc.MaxDownloads != nil && lc.Downloads >= *lc.MaxDownloads {
		return true
	}

	return false
}

// Files 文件集合对象，主要用于请求参数和响应参数
type Files struct {
	// 文件列表
	Files []*File `json:"files"`

	// 用户选择的参数/选项
	Options map[string]interface{} `json:"options"`

	// 用户备注/说明
	Note string `json:"note"`

	// 处理配置
	Config map[string]interface{} `json:"config"`

	// 创建时间
	CreatedAt typex.Time `json:"created_at"`

	// 默认生命周期（可以为nil表示无限制）
	DefaultLifecycle *FileLifecycle `json:"default_lifecycle,omitempty"`

	// 内部状态
	autoUpload bool
	ctx        context.Context
	tempDir    string
}

// FileInfo 简化的文件信息
type FileInfo struct {
	Name          string `json:"name"`
	Size          int64  `json:"size"`
	SizeFormatted string `json:"size_formatted"`
}

// NewFiles 创建新的文件集合（支持多种输入类型）
func NewFiles(input interface{}) *Files {
	fc := &Files{
		Files:      make([]*File, 0),
		Options:    make(map[string]interface{}),
		Config:     make(map[string]interface{}),
		CreatedAt:  typex.Time(time.Now()),
		autoUpload: true, // 默认自动上传
	}

	switch v := input.(type) {
	case []string:
		fc.addFromPaths(v)
	case []*File:
		fc.Files = v
	case string:
		fc.addFromPaths([]string{v})
	}

	return fc
}

// NewSingleFiles 创建新的文件(只有一个)
func NewSingleFiles(file *File) *Files {
	fc := &Files{
		Files:      []*File{file},
		Options:    make(map[string]interface{}),
		Config:     make(map[string]interface{}),
		CreatedAt:  typex.Time(time.Now()),
		autoUpload: true, // 默认自动上传
	}

	return fc
}

// NewFilesFromWriter 从Writer创建文件集合
func NewFilesFromWriter(writer Writer) *Files {
	return NewFiles(writer.GetFiles())
}

// addFromPaths 从路径数组添加文件
func (fc *Files) addFromPaths(paths []string) {
	for _, path := range paths {
		if err := fc.AddFileFromPath(path); err != nil {
			// 静默处理错误，避免阻塞整个流程
			continue
		}
	}
}

// SetContext 设置上下文
func (fc *Files) SetContext(ctx context.Context) *Files {
	fc.ctx = ctx
	return fc
}

// ===== 生命周期管理方法 =====

// SetMaxDownloads 设置下载次数限制
func (fc *Files) SetMaxDownloads(maxDownloads int) *Files {
	if fc.DefaultLifecycle == nil {
		fc.DefaultLifecycle = &FileLifecycle{}
	}
	fc.DefaultLifecycle.MaxDownloads = &maxDownloads
	return fc
}

// SetExpiresAt 设置过期时间
func (fc *Files) SetExpiresAt(expiresAt time.Time) *Files {
	if fc.DefaultLifecycle == nil {
		fc.DefaultLifecycle = &FileLifecycle{}
	}
	fc.DefaultLifecycle.ExpiresAt = &expiresAt
	return fc
}

// SetUnlimited 设置为无限制（清除生命周期限制）
func (fc *Files) SetUnlimited() *Files {
	fc.DefaultLifecycle = nil
	return fc
}

// SetTemporary 便捷方法：设置临时文件（下载一次后删除）
func (fc *Files) SetTemporary() *Files {
	return fc.SetMaxDownloads(1)
}

// SetExpiring7Days 便捷方法：设置7天后过期
func (fc *Files) SetExpiring7Days() *Files {
	return fc.SetExpiresAt(time.Now().Add(7 * 24 * time.Hour))
}

// ===== 文件添加方法 =====

// AddFileFromData 从数据创建文件并立即上传
func (fc *Files) AddFileFromData(filename string, data []byte, options ...FileOption) error {
	if data == nil {
		return fmt.Errorf("data不能为nil")
	}

	file := &File{
		Name:      filename,
		Size:      int64(len(data)),
		CreatedAt: typex.Time(time.Now()),
		Status:    "pending",
		Metadata:  make(map[string]string),
		Lifecycle: fc.DefaultLifecycle, // 使用默认生命周期
	}

	// 自动检测文件类型
	file.ContentType = fc.detectContentTypeFromData(data)

	// 应用选项
	for _, opt := range options {
		opt(file)
	}

	// 立即上传（如果启用自动上传）
	if fc.autoUpload {
		url, err := fc.uploadImmediately(filename, data)
		if err != nil {
			return fmt.Errorf("上传文件失败: %w", err)
		}
		file.URL = url
		file.Status = "uploaded"
		file.UpdatedAt = typex.Time(time.Now())
		file.uploaded = true
	}

	fc.Files = append(fc.Files, file)
	return nil
}

// AddFileFromPath 从本地路径添加文件并立即上传
func (fc *Files) AddFileFromPath(localPath string, options ...FileOption) error {
	// 检查文件是否存在
	info, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("文件不存在: %s", localPath)
	}

	// 读取文件数据
	data, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	filename := filepath.Base(localPath)

	// 创建文件对象
	file := &File{
		Name:      filename,
		Size:      info.Size(),
		CreatedAt: typex.Time(time.Now()),
		Status:    "pending",
		Metadata:  make(map[string]string),
		Lifecycle: fc.DefaultLifecycle,
	}

	// 自动检测文件类型
	file.ContentType = fc.detectContentTypeFromData(data)

	// 应用选项
	for _, opt := range options {
		opt(file)
	}

	// 立即上传
	if fc.autoUpload {
		url, err := fc.uploadImmediately(filename, data)
		if err != nil {
			return fmt.Errorf("上传文件失败: %w", err)
		}
		file.URL = url
		file.Status = "uploaded"
		file.UpdatedAt = typex.Time(time.Now())
		file.uploaded = true
	}

	fc.Files = append(fc.Files, file)
	return nil
}

// AddExistingFile 添加已存在的File对象（不上传）
func (fc *Files) AddExistingFile(file *File) *Files {
	if file != nil {
		fc.Files = append(fc.Files, file)
	}
	return fc
}

// AddFileFromURL 从URL添加文件引用（不下载）
func (fc *Files) AddFileFromURL(url, filename string, options ...FileOption) error {
	file := &File{
		Name:      filename,
		URL:       url,
		CreatedAt: typex.Time(time.Now()),
		Status:    "referenced", // 表示这是一个URL引用
		Metadata:  make(map[string]string),
		Lifecycle: fc.DefaultLifecycle,
	}

	// 应用选项
	for _, opt := range options {
		opt(file)
	}

	fc.Files = append(fc.Files, file)
	return nil
}

// ===== 元数据和配置方法 =====

// SetNote 设置备注（链式调用）
func (fc *Files) SetNote(note string) *Files {
	fc.Note = note
	return fc
}

// SetMetadata 设置元数据（链式调用）
func (fc *Files) SetMetadata(key string, value interface{}) *Files {
	if fc.Options == nil {
		fc.Options = make(map[string]interface{})
	}
	fc.Options[key] = value
	return fc
}

// AddMetadata 批量添加元数据（链式调用）
func (fc *Files) AddMetadata(metadata map[string]interface{}) *Files {
	if fc.Options == nil {
		fc.Options = make(map[string]interface{})
	}
	for k, v := range metadata {
		fc.Options[k] = v
	}
	return fc
}

// SetConfig 设置配置（链式调用）
func (fc *Files) SetConfig(key string, value interface{}) *Files {
	if fc.Config == nil {
		fc.Config = make(map[string]interface{})
	}
	fc.Config[key] = value
	return fc
}

// ===== 查询和统计方法 =====

// GetFiles 获取文件列表
func (fc *Files) GetFiles() []*File {
	return fc.Files
}

// GetNote 获取摘要信息
func (fc *Files) GetNote() string {
	return fc.Note
}

// GetMetadata 获取元数据
func (fc *Files) GetMetadata(key string) interface{} {
	if fc.Options == nil {
		return nil
	}
	return fc.Options[key]
}

// GetConfig 获取配置
func (fc *Files) GetConfig(key string) interface{} {
	if fc.Config == nil {
		return nil
	}
	return fc.Config[key]
}

// GetTotalSize 获取所有文件的总大小
func (fc *Files) GetTotalSize() int64 {
	var total int64
	for _, file := range fc.Files {
		total += file.Size
	}
	return total
}

// GetFileCount 获取文件数量
func (fc *Files) GetFileCount() int {
	return len(fc.Files)
}

// FilterByType 按文件类型过滤
func (fc *Files) FilterByType(contentType string) *Files {
	var filteredFiles []*File
	for _, file := range fc.Files {
		if file.ContentType == contentType {
			filteredFiles = append(filteredFiles, file)
		}
	}

	newCollection := &Files{
		Files:            filteredFiles,
		CreatedAt:        fc.CreatedAt,
		Options:          make(map[string]interface{}),
		Config:           make(map[string]interface{}),
		DefaultLifecycle: fc.DefaultLifecycle,
		autoUpload:       fc.autoUpload,
		ctx:              fc.ctx,
	}

	// 复制元数据
	for k, v := range fc.Options {
		newCollection.Options[k] = v
	}

	return newCollection
}

// FilterByExtension 按文件扩展名过滤
func (fc *Files) FilterByExtension(ext string) *Files {
	var filteredFiles []*File
	for _, file := range fc.Files {
		if getFileExtension(file.Name) == ext {
			filteredFiles = append(filteredFiles, file)
		}
	}

	newCollection := &Files{
		Files:            filteredFiles,
		CreatedAt:        fc.CreatedAt,
		Options:          make(map[string]interface{}),
		Config:           make(map[string]interface{}),
		DefaultLifecycle: fc.DefaultLifecycle,
		autoUpload:       fc.autoUpload,
		ctx:              fc.ctx,
	}

	// 复制元数据
	for k, v := range fc.Options {
		newCollection.Options[k] = v
	}

	return newCollection
}

// ===== 内部辅助方法 =====

// getTempDir 获取项目临时目录
func (fc *Files) getTempDir() (string, error) {
	if fc.tempDir != "" {
		return fc.tempDir, nil
	}

	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取工作目录失败: %w", err)
	}

	// 使用traceID作为临时目录名（与现有实现保持一致）
	traceID := "default"
	if fc.ctx != nil {
		if ctx := fc.ctx.Value("trace_id"); ctx != nil {
			if id, ok := ctx.(string); ok {
				traceID = id
			}
		}
	}

	// 创建临时目录: ./temp/traceID（与现有实现一致）
	tempBase := filepath.Join(workDir, "temp")
	if err := os.MkdirAll(tempBase, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	fc.tempDir = filepath.Join(tempBase, traceID)
	if err := os.MkdirAll(fc.tempDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	return fc.tempDir, nil
}

// uploadImmediately 立即上传文件
func (fc *Files) uploadImmediately(filename string, data []byte) (string, error) {
	// 使用项目临时目录而不是系统临时目录
	tempDir, err := fc.getTempDir()
	if err != nil {
		return "", err
	}

	tempPath := filepath.Join(tempDir, filename)
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return "", err
	}
	defer os.Remove(tempPath) // 上传完立即删除临时文件

	// 执行上传
	uploader := fc.getUploader()
	if uploader == nil {
		return "", fmt.Errorf("无法获取上传器")
	}

	return uploader.UploadFile(tempPath, filename)
}

// getUploader 获取上传器
func (fc *Files) getUploader() Uploader {
	if fc.ctx == nil {
		return nil
	}

	// 从context中获取上传配置
	functionMsg := fc.getFunctionMsg()
	if functionMsg == nil {
		return nil
	}

	uploadConfig := functionMsg.UploadConfig

	// 使用工厂创建上传器
	factory := NewUploaderFactory()
	uploader, err := factory.CreateUploaderWithMsg(uploadConfig, functionMsg)
	if err != nil {
		return nil
	}

	return uploader
}

// getFunctionMsg 从context中获取FunctionMsg
func (fc *Files) getFunctionMsg() *trace.FunctionMsg {
	if fc.ctx == nil {
		return nil
	}

	value := fc.ctx.Value(trace.FunctionMsgKey)
	if value == nil {
		return nil
	}

	if msg, ok := value.(*trace.FunctionMsg); ok {
		return msg
	}

	return nil
}

// detectContentTypeFromData 从数据内容检测MIME类型
func (fc *Files) detectContentTypeFromData(data []byte) string {
	mtype := mimetype.Detect(data)
	return mtype.String()
}

// CleanupLocalFiles 清理本地缓存文件
func (fc *Files) CleanupLocalFiles() error {
	var errors []error

	for _, file := range fc.Files {
		if file.LocalPath != "" && file.localCached {
			if err := os.Remove(file.LocalPath); err != nil && !os.IsNotExist(err) {
				errors = append(errors, fmt.Errorf("删除文件 %s 失败: %w", file.LocalPath, err))
			} else {
				file.LocalPath = ""
				file.localCached = false
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("清理过程中发生 %d 个错误", len(errors))
	}

	return nil
}

// ===== JSON序列化支持 =====

// MarshalJSON 实现json.Marshaler接口
func (fc *Files) MarshalJSON() ([]byte, error) {
	// 创建一个临时结构体用于序列化
	type filesAlias Files
	return json.Marshal((*filesAlias)(fc))
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (fc *Files) UnmarshalJSON(data []byte) error {
	// 创建一个临时结构体用于反序列化
	type filesAlias Files
	alias := (*filesAlias)(fc)

	// 先初始化 map，避免 nil pointer
	if alias.Options == nil {
		alias.Options = make(map[string]interface{})
	}
	if alias.Config == nil {
		alias.Config = make(map[string]interface{})
	}

	return json.Unmarshal(data, alias)
}

// ===== 数据库序列化支持 =====

// Scan 实现 sql.Scanner 接口，用于从数据库读取
func (fc *Files) Scan(value interface{}) error {
	if value == nil {
		*fc = Files{}
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("cannot scan %T into Files", value)
	}

	return json.Unmarshal(data, fc)
}

// Value 实现 driver.Valuer 接口，用于存储到数据库
func (fc Files) Value() (driver.Value, error) {
	if len(fc.Files) == 0 {
		return nil, nil
	}
	return json.Marshal(fc)
}

// ===== 兼容性方法 =====

// ToReader 转换为Reader接口（兼容旧的Files切片类型的方法）
func (fc *Files) ToReader(ctx context.Context) Reader {
	reader := NewURLReader(ctx)
	for _, file := range fc.Files {
		reader.AddFileFromURL(file.URL, file.Name)
	}
	return reader
}

// ToWriter 转换为Writer接口（兼容旧的Files切片类型的方法）
func (fc *Files) ToWriter(ctx context.Context) Writer {
	writer := NewCloudWriter(ctx)
	for _, file := range fc.Files {
		if file.LocalPath != "" {
			writer.AddFile(file.LocalPath)
		}
	}
	return writer
}

// getFileExtension 获取文件扩展名
func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}

// GORM 使用示例:
//
// type Task struct {
//     ID        uint           `gorm:"primarykey"`
//     Name      string         `gorm:"column:name;comment:任务名称"`
//     Files     *Files         `gorm:"type:json;column:files;comment:任务文件"`
//     CreatedAt time.Time
//     UpdatedAt time.Time
// }
//
// 使用方法:
// func CreateTask(db *gorm.DB, name string, files *Files) error {
//     task := &Task{
//         Name:  name,
//         Files: files,
//     }
//     return db.Create(task).Error
// }
//
// func GetTask(db *gorm.DB, id uint) (*Task, error) {
//     var task Task
//     err := db.First(&task, id).Error
//     return &task, err
// }
