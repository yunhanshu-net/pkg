package typex

import (
	"crypto/md5"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// File 文件类型，支持URL传递和文件操作
type File struct {
	URL         string    `json:"url"`          // 文件URL地址
	Name        string    `json:"name"`         // 文件名
	Size        int64     `json:"size"`         // 文件大小（字节）
	ContentType string    `json:"content_type"` // MIME类型
	UploadTime  time.Time `json:"upload_time"`  // 上传时间
	Hash        string    `json:"hash"`         // 文件哈希值（MD5）
}

// Files 文件集合类型，支持多文件操作
type Files []File

// NewFile 创建新的文件对象
func NewFile(url, name string) *File {
	f := &File{
		URL:        url,
		Name:       name,
		UploadTime: time.Now(),
	}

	// 自动推断ContentType
	if ext := filepath.Ext(name); ext != "" {
		f.ContentType = mime.TypeByExtension(ext)
	}

	return f
}

// NewFileFromURL 从URL创建文件对象
func NewFileFromURL(url string) *File {
	name := filepath.Base(url)
	// 移除URL参数
	if idx := strings.Index(name, "?"); idx != -1 {
		name = name[:idx]
	}

	return NewFile(url, name)
}

// NewFileWithHash 创建带MD5哈希的文件对象
func NewFileWithHash(url, name, hash string) *File {
	f := NewFile(url, name)
	f.Hash = hash
	return f
}

// ===== JSON序列化支持 =====

func (f *File) UnmarshalJSON(b []byte) error {
	// 支持简单字符串URL格式
	if b[0] == '"' {
		var url string
		if err := json.Unmarshal(b, &url); err != nil {
			return err
		}
		*f = *NewFileFromURL(url)
		return nil
	}

	// 支持完整对象格式
	type fileAlias File
	return json.Unmarshal(b, (*fileAlias)(f))
}

func (f File) MarshalJSON() ([]byte, error) {
	// 如果只有URL，返回简单字符串格式
	if f.Name == "" && f.Size == 0 && f.ContentType == "" {
		return json.Marshal(f.URL)
	}

	// 返回完整对象格式
	type fileAlias File
	return json.Marshal(fileAlias(f))
}

// ===== 数据库支持 =====

func (f *File) Scan(value interface{}) error {
	if value == nil {
		*f = File{}
		return nil
	}

	switch v := value.(type) {
	case string:
		*f = *NewFileFromURL(v)
	case []byte:
		*f = *NewFileFromURL(string(v))
	default:
		return fmt.Errorf("cannot scan %T into File", value)
	}

	return nil
}

func (f File) Value() (driver.Value, error) {
	if f.URL == "" {
		return nil, nil
	}
	return f.URL, nil
}

// ===== 文件操作方法 =====

// Download 下载文件内容
func (f *File) Download() ([]byte, error) {
	if f.URL == "" {
		return nil, fmt.Errorf("file URL is empty")
	}

	resp, err := http.Get(f.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	// 更新文件信息
	f.Size = int64(len(data))
	if f.ContentType == "" {
		f.ContentType = resp.Header.Get("Content-Type")
	}

	return data, nil
}

// DownloadToReader 下载文件并返回Reader
func (f *File) DownloadToReader() (io.ReadCloser, error) {
	if f.URL == "" {
		return nil, fmt.Errorf("file URL is empty")
	}

	resp, err := http.Get(f.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to download file: status %d", resp.StatusCode)
	}

	// 更新文件信息
	if f.ContentType == "" {
		f.ContentType = resp.Header.Get("Content-Type")
	}

	return resp.Body, nil
}

// Upload 上传文件内容（需要实现具体的上传逻辑）
func (f *File) Upload(data []byte, uploader FileUploader) error {
	if uploader == nil {
		return fmt.Errorf("uploader is nil")
	}

	url, err := uploader.Upload(data, f.Name, f.ContentType)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	f.URL = url
	f.Size = int64(len(data))
	f.UploadTime = time.Now()

	return nil
}

// ===== File哈希相关方法 =====

// GetHash 获取文件哈希值
func (f *File) GetHash() string {
	return f.Hash
}

// SetHash 设置文件哈希值
func (f *File) SetHash(hash string) {
	f.Hash = hash
}

// HasValidHash 检查是否有有效的哈希值
func (f *File) HasValidHash() bool {
	return len(strings.TrimSpace(f.Hash)) > 0
}

// GetMD5 获取文件MD5值（兼容性方法，实际可能是其他哈希算法）
func (f *File) GetMD5() string {
	return f.Hash
}

// SetMD5 设置文件MD5值（兼容性方法，实际可能是其他哈希算法）
func (f *File) SetMD5(hash string) {
	f.Hash = hash
}

// HasValidMD5 检查是否有有效的MD5值（兼容性方法，实际检查任何哈希值）
func (f *File) HasValidMD5() bool {
	return f.HasValidHash()
}

// CalculateMD5FromData 从数据计算MD5并设置到Hash字段
func (f *File) CalculateMD5FromData(data []byte) {
	hash := md5.Sum(data)
	f.Hash = hex.EncodeToString(hash[:])
}

// VerifyMD5 验证文件的MD5值（需要下载文件）
// 注意：如果Hash字段存储的不是MD5（如七牛云ETag），此方法会返回false
func (f *File) VerifyMD5() (bool, error) {
	if !f.HasValidHash() {
		return false, fmt.Errorf("文件没有哈希值")
	}

	// 下载文件内容
	resp, err := http.Get(f.URL)
	if err != nil {
		return false, fmt.Errorf("下载文件失败: %v", err)
	}
	defer resp.Body.Close()

	// 计算MD5
	hash := md5.New()
	_, err = io.Copy(hash, resp.Body)
	if err != nil {
		return false, fmt.Errorf("读取文件内容失败: %v", err)
	}

	calculatedMD5 := hex.EncodeToString(hash.Sum(nil))
	return strings.EqualFold(f.Hash, calculatedMD5), nil
}

// IsQiniuETag 判断哈希值是否可能是七牛云ETag格式
func (f *File) IsQiniuETag() bool {
	if !f.HasValidHash() {
		return false
	}

	// 七牛云ETag解码后的长度：
	// - 小文件：1字节前缀 + 20字节SHA1 = 21字节 = 42个十六进制字符
	// - 大文件：1字节前缀 + 20字节SHA1 = 21字节 = 42个十六进制字符
	return len(f.Hash) == 42
}

// IsStandardMD5 判断哈希值是否可能是标准MD5格式
func (f *File) IsStandardMD5() bool {
	if !f.HasValidHash() {
		return false
	}

	// 标准MD5是32个十六进制字符
	return len(f.Hash) == 32
}

// IsImage 判断是否为图片文件
func (f *File) IsImage() bool {
	return strings.HasPrefix(f.ContentType, "image/")
}

// IsVideo 判断是否为视频文件
func (f *File) IsVideo() bool {
	return strings.HasPrefix(f.ContentType, "video/")
}

// IsAudio 判断是否为音频文件
func (f *File) IsAudio() bool {
	return strings.HasPrefix(f.ContentType, "audio/")
}

// GetExtension 获取文件扩展名
func (f *File) GetExtension() string {
	return filepath.Ext(f.Name)
}

// String 返回文件的字符串表示
func (f File) String() string {
	if f.Name != "" {
		return fmt.Sprintf("%s (%s)", f.Name, f.URL)
	}
	return f.URL
}

// ===== Files集合操作 =====

// Add 添加文件到集合
func (fs *Files) Add(file File) {
	*fs = append(*fs, file)
}

// AddURL 通过URL添加文件
func (fs *Files) AddURL(url string) {
	*fs = append(*fs, *NewFileFromURL(url))
}

// Filter 过滤文件
func (fs Files) Filter(predicate func(File) bool) Files {
	var result Files
	for _, f := range fs {
		if predicate(f) {
			result = append(result, f)
		}
	}
	return result
}

// Images 获取所有图片文件
func (fs Files) Images() Files {
	return fs.Filter(func(f File) bool { return f.IsImage() })
}

// Videos 获取所有视频文件
func (fs Files) Videos() Files {
	return fs.Filter(func(f File) bool { return f.IsVideo() })
}

// First 获取第一个文件
func (fs Files) First() *File {
	if len(fs) == 0 {
		return nil
	}
	return &fs[0]
}

// URLs 获取所有文件的URL
func (fs Files) URLs() []string {
	urls := make([]string, len(fs))
	for i, f := range fs {
		urls[i] = f.URL
	}
	return urls
}

// TotalSize 获取总文件大小
func (fs Files) TotalSize() int64 {
	var total int64
	for _, f := range fs {
		total += f.Size
	}
	return total
}

// ===== Files集合哈希相关方法 =====

// WithValidHash 获取所有有有效哈希值的文件
func (fs Files) WithValidHash() Files {
	return fs.Filter(func(f File) bool { return f.HasValidHash() })
}

// WithoutHash 获取所有没有哈希值的文件
func (fs Files) WithoutHash() Files {
	return fs.Filter(func(f File) bool { return !f.HasValidHash() })
}

// FindByHash 根据哈希值查找文件
func (fs Files) FindByHash(hash string) *File {
	hash = strings.ToLower(hash)
	for i, f := range fs {
		if strings.ToLower(f.Hash) == hash {
			return &fs[i]
		}
	}
	return nil
}

// FindDuplicatesByHash 查找具有相同哈希值的重复文件
func (fs Files) FindDuplicatesByHash() map[string]Files {
	hashMap := make(map[string]Files)

	for _, f := range fs {
		if f.HasValidHash() {
			hash := strings.ToLower(f.Hash)
			hashMap[hash] = append(hashMap[hash], f)
		}
	}

	// 只返回有重复的哈希值
	duplicates := make(map[string]Files)
	for hash, files := range hashMap {
		if len(files) > 1 {
			duplicates[hash] = files
		}
	}

	return duplicates
}

// GetHashList 获取所有文件的哈希值列表
func (fs Files) GetHashList() []string {
	var hashes []string
	for _, f := range fs {
		if f.HasValidHash() {
			hashes = append(hashes, f.Hash)
		}
	}
	return hashes
}

// GroupByHashType 按哈希类型分组文件
func (fs Files) GroupByHashType() map[string]Files {
	groups := make(map[string]Files)

	for _, f := range fs {
		if !f.HasValidHash() {
			groups["no_hash"] = append(groups["no_hash"], f)
		} else if f.IsStandardMD5() {
			groups["md5"] = append(groups["md5"], f)
		} else if f.IsQiniuETag() {
			groups["qiniu_etag"] = append(groups["qiniu_etag"], f)
		} else {
			groups["other"] = append(groups["other"], f)
		}
	}

	return groups
}

// ===== 上传器接口 =====

// FileUploader 文件上传器接口
type FileUploader interface {
	Upload(data []byte, filename, contentType string) (url string, err error)
}

// ===== 辅助函数 =====

// ParseFileURL 解析文件URL，提取文件信息
func ParseFileURL(url string) *File {
	return NewFileFromURL(url)
}

// ParseFileURLs 解析多个文件URL
func ParseFileURLs(urls []string) Files {
	files := make(Files, len(urls))
	for i, url := range urls {
		files[i] = *NewFileFromURL(url)
	}
	return files
}

// ===== 兼容性方法（MD5相关） =====

// WithValidMD5 获取所有有有效MD5的文件（兼容性方法）
func (fs Files) WithValidMD5() Files {
	return fs.WithValidHash()
}

// WithoutMD5 获取所有没有MD5的文件（兼容性方法）
func (fs Files) WithoutMD5() Files {
	return fs.WithoutHash()
}

// FindByMD5 根据MD5查找文件（兼容性方法）
func (fs Files) FindByMD5(hash string) *File {
	return fs.FindByHash(hash)
}

// FindDuplicatesByMD5 查找具有相同MD5的重复文件（兼容性方法）
func (fs Files) FindDuplicatesByMD5() map[string]Files {
	return fs.FindDuplicatesByHash()
}
