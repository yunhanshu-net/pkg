package files

import "context"

// Reader 文件读取接口
type Reader interface {
	// AddFileFromURL 从URL添加文件
	AddFileFromURL(url, name string, opts ...FileOption) error

	// DownloadFile 下载文件到指定路径
	DownloadFile(index int, localPath string) (string, error)

	// DownloadAll 下载所有文件到指定目录
	DownloadAll(dir string) ([]string, error)

	// FilterByType 按MIME类型过滤文件
	FilterByType(contentType string) []*File

	// FilterByExtension 按扩展名过滤文件
	FilterByExtension(ext string) []*File

	// GetFile 获取指定索引的文件
	GetFile(index int) (*File, error)

	// GetFiles 获取所有文件
	GetFiles() []*File

	// GetTotalSize 获取所有文件的总大小
	GetTotalSize() int64

	// SetSummary 设置文件集合摘要
	SetSummary(summary string)

	// GetTempDir 获取临时目录
	GetTempDir() (string, error)

	// Cleanup 清理临时文件
	Cleanup()
}

// Writer 文件写入接口
type Writer interface {
	// AddFile 添加本地文件
	AddFile(localPath string, opts ...FileOption) error

	// AddFileWithData 从数据创建文件
	AddFileWithData(name string, data []byte, opts ...FileOption) error

	// CreateTempFile 创建临时文件
	CreateTempFile(name string) (string, error)

	// SetSummary 设置文件集合摘要（支持链式调用）
	SetSummary(summary string) Writer

	// SetLifecycle 设置文件生命周期（支持链式调用）
	SetLifecycle(lifecycle Lifecycle) Writer

	// SetMetadata 设置元数据（支持链式调用）
	SetMetadata(key string, value interface{}) Writer

	// AddMetadata 批量添加元数据（支持链式调用）
	AddMetadata(metadata map[string]interface{}) Writer

	// GetFiles 获取所有文件
	GetFiles() []*File

	// GetTempDir 获取临时目录
	GetTempDir() (string, error)

	// Cleanup 清理临时文件
	Cleanup()
}

// 构造函数类型
type ReaderConstructor func(ctx context.Context) Reader
type WriterConstructor func(ctx context.Context) Writer
