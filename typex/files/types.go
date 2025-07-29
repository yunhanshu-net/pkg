package files

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	// 生命周期管理（可以为nil表示无限制）
	Lifecycle *FileLifecycle `json:"lifecycle,omitempty"`

	// 内部状态字段
	uploaded    bool `json:"-"` // 是否已上传
	localCached bool `json:"-"` // 是否为本地缓存文件
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

// GetLocalPath 获取本地路径（按需下载）
func (f *File) GetLocalPath(ctx context.Context) (string, error) {
	// 如果已有本地路径且文件存在
	if f.LocalPath != "" {
		if _, err := os.Stat(f.LocalPath); err == nil {
			return f.LocalPath, nil
		}
	}

	// 如果有URL，下载到项目临时目录
	if f.URL != "" {
		localPath, err := f.downloadToProjectTemp(ctx)
		if err != nil {
			return "", err
		}
		f.LocalPath = localPath
		f.localCached = true
		return localPath, nil
	}

	return "", fmt.Errorf("文件既没有本地路径也没有URL")
}

// downloadToProjectTemp 下载文件到项目临时目录
func (f *File) downloadToProjectTemp(ctx context.Context) (string, error) {
	// 获取项目临时目录
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取工作目录失败: %w", err)
	}

	// 使用traceID作为临时目录名
	traceID := "default"
	if ctx != nil {
		if ctxTraceID := ctx.Value("trace_id"); ctxTraceID != nil {
			if id, ok := ctxTraceID.(string); ok {
				traceID = id
			}
		}
	}

	// 创建临时目录: ./temp/traceID
	tempBase := filepath.Join(workDir, "temp")
	if err := os.MkdirAll(tempBase, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	tempDir := filepath.Join(tempBase, traceID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 下载文件到临时目录
	localPath := filepath.Join(tempDir, f.Name)

	// 执行下载逻辑
	req, err := http.NewRequestWithContext(ctx, "GET", f.URL, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	// 创建本地文件
	file, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return localPath, nil
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

// WithLifecycle 设置文件生命周期
func WithLifecycle(lifecycle *FileLifecycle) FileOption {
	return func(f *File) {
		f.Lifecycle = lifecycle
	}
}

// WithURL 设置文件URL
func WithURL(url string) FileOption {
	return func(f *File) {
		f.URL = url
	}
}
