package files

import (
	"context"
	"strings"
	"time"

	"github.com/yunhanshu-net/pkg/typex"
)

// File 文件结构体
type File struct {
	Name        string            `json:"name"`         // 文件名
	Size        int64             `json:"size"`         // 文件大小
	ContentType string            `json:"content_type"` // 文件类型
	URL         string            `json:"url"`          // 上传/下载URL
	PreviewURL  string            `json:"preview_url"`  // 预览URL（响应时使用）
	LocalPath   string            `json:"-"`            // 本地路径（请求时使用，不序列化到JSON）
	CreatedAt   typex.Time        `json:"created_at"`   // 创建时间
	UpdatedAt   typex.Time        `json:"updated_at"`   // 更新时间（响应时使用）
	Status      string            `json:"status"`       // 文件状态（响应时使用）
	Metadata    map[string]string `json:"metadata"`     // 元数据
	Downloaded  bool              `json:"downloaded,omitempty"`
	DownloadAt  typex.Time        `json:"download_at,omitempty"`
	Hash        string            `json:"hash,omitempty"`
	Description string            `json:"description,omitempty"`
	AutoDelete  bool              `json:"auto_delete,omitempty"`
	Compressed  bool              `json:"compressed,omitempty"`
}

// 注意：Files 类型已移动到 files.go 中，现在是一个增强的文件集合对象

// FileList 文件列表类型，用于API请求参数
type FileList []*File

// ToReader 转换为Reader接口
func (f FileList) ToReader(ctx context.Context) Reader {
	reader := NewURLReader(ctx)
	for _, file := range f {
		reader.AddFileFromURL(file.URL, file.Name)
	}
	return reader
}

// ToWriter 转换为Writer接口
func (f FileList) ToWriter(ctx context.Context) Writer {
	writer := NewCloudWriter(ctx)
	for _, file := range f {
		writer.AddFile(file.LocalPath)
	}
	return writer
}

// GetTotalSize 获取所有文件的总大小
func (f FileList) GetTotalSize() int64 {
	var total int64
	for _, file := range f {
		total += file.Size
	}
	return total
}

// FilterByType 按文件类型过滤
func (f FileList) FilterByType(contentType string) FileList {
	var result FileList
	for _, file := range f {
		if file.ContentType == contentType {
			result = append(result, file)
		}
	}
	return result
}

// FilterByExtension 按文件扩展名过滤
func (f FileList) FilterByExtension(ext string) FileList {
	var result FileList
	for _, file := range f {
		if strings.HasSuffix(strings.ToLower(file.Name), strings.ToLower(ext)) {
			result = append(result, file)
		}
	}
	return result
}

// ToUploadFile 转换为上传文件
func (f *File) ToUploadFile() *File {
	return &File{
		Name:        f.Name,
		Size:        f.Size,
		ContentType: f.ContentType,
		URL:         f.URL,       // 上传URL
		LocalPath:   f.LocalPath, // 本地路径
		CreatedAt:   typex.Time(time.Now()),
		Metadata:    f.Metadata,
	}
}

// ToDisplayFile 转换为展示文件
func (f *File) ToDisplayFile() *File {
	return &File{
		Name:        f.Name,
		Size:        f.Size,
		ContentType: f.ContentType,
		URL:         f.URL,        // 下载URL
		PreviewURL:  f.PreviewURL, // 预览URL
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
		Status:      f.Status,
		Metadata:    f.Metadata,
	}
}

// 文件类型常量
const (
	TypeImage    = "image"
	TypeDocument = "document"
	TypeVideo    = "video"
	TypeAudio    = "audio"
	TypeArchive  = "archive"
	TypeOther    = "other"
)

// 文件状态常量
const (
	StatusPending   = "pending"   // 待处理
	StatusUploading = "uploading" // 上传中
	StatusSuccess   = "success"   // 成功
	StatusFailed    = "failed"    // 失败
	StatusDeleted   = "deleted"   // 已删除
)

// FileOption 文件选项函数类型
type FileOption func(*File)

// Lifecycle 文件生命周期类型
type Lifecycle string

// 常量定义
const (
	TraceID = "trace_id"

	// 生命周期常量
	LifecycleTemporary Lifecycle = "temporary"  // 处理完立即删除
	LifecycleShortTerm Lifecycle = "short_term" // 下载后删除
	LifecycleLongTerm  Lifecycle = "long_term"  // 永久保存
	LifecycleCache     Lifecycle = "cache"      // 定期清理

	// 文件类型常量
	FileTypeImage    = "image"
	FileTypeDocument = "document"
	FileTypeVideo    = "video"
	FileTypeAudio    = "audio"
	FileTypeArchive  = "archive"
	FileTypeOther    = "other"
)

// 文件类型映射
var FileTypeMap = map[string][]string{
	FileTypeImage: {
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
	},
	FileTypeDocument: {
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"text/plain",
		"text/csv",
	},
	FileTypeVideo: {
		"video/mp4",
		"video/webm",
		"video/quicktime",
	},
	FileTypeAudio: {
		"audio/mpeg",
		"audio/wav",
		"audio/ogg",
	},
	FileTypeArchive: {
		"application/zip",
		"application/x-rar-compressed",
		"application/x-7z-compressed",
		"application/x-tar",
		"application/gzip",
	},
}

// 选项函数

// WithSize 设置文件大小
func WithSize(size int64) FileOption {
	return func(f *File) {
		f.Size = size
	}
}

// WithContentType 设置文件MIME类型
func WithContentType(contentType string) FileOption {
	return func(f *File) {
		f.ContentType = contentType
	}
}

// WithHash 设置文件哈希
func WithHash(hash string) FileOption {
	return func(f *File) {
		f.Hash = hash
	}
}

// WithDescription 设置文件描述
func WithDescription(desc string) FileOption {
	return func(f *File) {
		f.Description = desc
	}
}

// WithMetadata 设置文件元数据
func WithMetadata(key, value string) FileOption {
	return func(f *File) {
		if f.Metadata == nil {
			f.Metadata = make(map[string]string)
		}
		f.Metadata[key] = value
	}
}

// WithAutoDelete 设置自动删除
func WithAutoDelete(autoDelete bool) FileOption {
	return func(f *File) {
		f.AutoDelete = autoDelete
	}
}

// WithCompressed 设置压缩标志
func WithCompressed(compressed bool) FileOption {
	return func(f *File) {
		f.Compressed = compressed
	}
}
