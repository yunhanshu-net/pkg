package files

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/yunhanshu-net/pkg/typex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// URLReader 从URL读取文件的实现
type URLReader struct {
	ctx     context.Context
	files   []*File
	tempDir string
	summary string
}

// NewURLReader 创建新的URL文件读取器
func NewURLReader(ctx context.Context) Reader {
	return &URLReader{
		ctx:   ctx,
		files: make([]*File, 0),
	}
}

// AddFileFromURL 从URL添加文件
func (r *URLReader) AddFileFromURL(url, name string, opts ...FileOption) error {
	file := &File{
		URL:       url,
		Name:      name,
		CreatedAt: typex.Time(time.Now()),
		Metadata:  make(map[string]string),
	}

	// 应用选项
	for _, opt := range opts {
		opt(file)
	}

	// 如果没有设置大小和类型，尝试从HTTP头获取
	if file.Size == 0 || file.ContentType == "" {
		if err := r.fetchFileInfo(file); err != nil {
			// 如果获取失败，不阻塞添加过程，只记录警告
			fmt.Printf("Warning: failed to fetch file info for %s: %v\n", url, err)
		}
	}

	r.files = append(r.files, file)
	return nil
}

// fetchFileInfo 从HTTP头获取文件信息
func (r *URLReader) fetchFileInfo(file *File) error {
	req, err := http.NewRequestWithContext(r.ctx, "HEAD", file.URL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if file.Size == 0 && resp.ContentLength > 0 {
		file.Size = resp.ContentLength
	}

	if file.ContentType == "" {
		file.ContentType = resp.Header.Get("Content-Type")
	}

	return nil
}

// DownloadFile 下载文件到指定路径
func (r *URLReader) DownloadFile(index int, localPath string) (string, error) {
	if index < 0 || index >= len(r.files) {
		return "", fmt.Errorf("文件索引超出范围: %d", index)
	}

	file := r.files[index]
	if file.Downloaded {
		return file.LocalPath, nil
	}

	// 如果未指定本地路径，则使用临时目录
	if localPath == "" {
		tempDir, err := r.GetTempDir()
		if err != nil {
			return "", fmt.Errorf("获取临时目录失败: %v", err)
		}
		localPath = filepath.Join(tempDir, file.Name)
	}

	// 下载文件
	if err := r.downloadFileToPath(file, localPath); err != nil {
		return "", err
	}

	// 计算MD5
	hash, err := r.calculateMD5(localPath)
	if err != nil {
		return "", fmt.Errorf("计算MD5失败: %v", err)
	}

	// 更新文件信息
	file.LocalPath = localPath
	file.Downloaded = true
	file.DownloadAt = typex.Time(time.Now())
	file.Hash = hash

	return localPath, nil
}

// downloadFileToPath 下载文件到指定路径
func (r *URLReader) downloadFileToPath(file *File, localPath string) error {
	req, err := http.NewRequestWithContext(r.ctx, "GET", file.URL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	// 创建本地文件
	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// 写入文件
	_, err = io.Copy(f, resp.Body)
	return err
}

// calculateMD5 计算文件MD5哈希
func (r *URLReader) calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// DownloadAll 下载所有文件到指定目录
func (r *URLReader) DownloadAll(dir string) ([]string, error) {
	if dir == "" {
		tempDir, err := r.GetTempDir()
		if err != nil {
			return nil, fmt.Errorf("获取临时目录失败: %v", err)
		}
		dir = tempDir
	}

	var paths []string
	for i := range r.files {
		localPath := filepath.Join(dir, r.files[i].Name)
		path, err := r.DownloadFile(i, localPath)
		if err != nil {
			return nil, fmt.Errorf("下载文件失败: %v", err)
		}
		paths = append(paths, path)
	}

	return paths, nil
}

// GetFile 获取指定索引的文件信息
func (r *URLReader) GetFile(index int) (*File, error) {
	if index < 0 || index >= len(r.files) {
		return nil, fmt.Errorf("file index %d out of range", index)
	}
	return r.files[index], nil
}

// GetFiles 获取所有文件
func (r *URLReader) GetFiles() []*File {
	return r.files
}

// FilterByType 按MIME类型过滤文件
func (r *URLReader) FilterByType(contentType string) []*File {
	var result []*File
	for _, file := range r.files {
		if file.ContentType == contentType {
			result = append(result, file)
		}
	}
	return result
}

// FilterByExtension 按扩展名过滤文件
func (r *URLReader) FilterByExtension(ext string) []*File {
	var result []*File
	for _, file := range r.files {
		if strings.HasSuffix(strings.ToLower(file.Name), strings.ToLower(ext)) {
			result = append(result, file)
		}
	}
	return result
}

// GetTotalSize 获取所有文件的总大小
func (r *URLReader) GetTotalSize() int64 {
	var total int64
	for _, file := range r.files {
		total += file.Size
	}
	return total
}

// SetSummary 设置文件集合摘要
func (r *URLReader) SetSummary(summary string) {
	r.summary = summary
}

// GetTempDir 获取临时目录
func (r *URLReader) GetTempDir() (string, error) {
	if r.tempDir != "" {
		return r.tempDir, nil
	}

	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取工作目录失败: %w", err)
	}

	// 使用traceID作为临时目录名
	traceID := "default"
	if ctx := r.ctx.Value("trace_id"); ctx != nil {
		if id, ok := ctx.(string); ok {
			traceID = id
		}
	}

	// 创建临时目录: ./temp/traceID
	tempBase := filepath.Join(workDir, "temp")
	if err := os.MkdirAll(tempBase, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	r.tempDir = filepath.Join(tempBase, traceID)
	if err := os.MkdirAll(r.tempDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	return r.tempDir, nil
}

// Cleanup 清理临时文件
func (r *URLReader) Cleanup() {
	if r.tempDir != "" {
		os.RemoveAll(r.tempDir)
		r.tempDir = ""
	}
}

// GetFileType 根据MIME类型获取文件类型
func GetFileType(contentType string) string {
	for fileType, mimeTypes := range FileTypeMap {
		for _, mimeType := range mimeTypes {
			if contentType == mimeType {
				return fileType
			}
		}
	}
	return FileTypeOther
}

// IsImage 判断是否为图片文件
func IsImage(contentType string) bool {
	return GetFileType(contentType) == FileTypeImage
}

// IsDocument 判断是否为文档文件
func IsDocument(contentType string) bool {
	return GetFileType(contentType) == FileTypeDocument
}

// IsVideo 判断是否为视频文件
func IsVideo(contentType string) bool {
	return GetFileType(contentType) == FileTypeVideo
}

// IsAudio 判断是否为音频文件
func IsAudio(contentType string) bool {
	return GetFileType(contentType) == FileTypeAudio
}

// IsArchive 判断是否为压缩文件
func IsArchive(contentType string) bool {
	return GetFileType(contentType) == FileTypeArchive
}
