package files

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/yunhanshu-net/pkg/trace"
	"github.com/yunhanshu-net/pkg/typex"
)

// CloudWriter 云存储文件写入器实现
type CloudWriter struct {
	ctx       context.Context
	files     []*File
	tempDir   string
	summary   string
	lifecycle Lifecycle
	metadata  map[string]interface{} // 添加元数据字段
}

// NewCloudWriter 创建新的云存储文件写入器
func NewCloudWriter(ctx context.Context) Writer {
	return &CloudWriter{
		ctx:       ctx,
		files:     make([]*File, 0),
		lifecycle: LifecycleTemporary,
	}
}

// AddFile 添加本地文件到输出列表
func (w *CloudWriter) AddFile(localPath string, options ...FileOption) error {
	// 检查文件是否存在
	info, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("file not found: %s", localPath)
	}

	// 创建文件对象
	file := &File{
		Name:      filepath.Base(localPath),
		Size:      info.Size(),
		LocalPath: localPath,
		CreatedAt: typex.Time(time.Now()),
		Metadata:  make(map[string]string),
	}

	// 自动检测文件类型（内部兜底）
	if detectedType, err := w.detectContentType(localPath); err == nil {
		file.ContentType = detectedType
	}

	// 应用选项（用户手动设置的会覆盖自动检测的）
	for _, opt := range options {
		opt(file)
	}

	// 如果没有设置URL，生成一个占位符（实际上传时会替换）
	if file.URL == "" {
		file.URL = fmt.Sprintf("pending://upload/%s", file.Name)
	}

	w.files = append(w.files, file)
	return nil
}

// AddFileWithData 从数据创建文件并添加到输出列表
func (w *CloudWriter) AddFileWithData(filename string, data []byte, options ...FileOption) error {
	// 创建临时文件
	tempPath, err := w.CreateTempFile(filename)
	if err != nil {
		return err
	}

	// 写入数据
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return err
	}

	// 创建文件对象，先自动检测类型
	file := &File{
		Name:      filename,
		Size:      int64(len(data)),
		LocalPath: tempPath,
		CreatedAt: typex.Time(time.Now()),
		Metadata:  make(map[string]string),
	}

	// 自动检测文件类型（基于数据内容）
	if detectedType := w.detectContentTypeFromData(data); detectedType != "" {
		file.ContentType = detectedType
	}

	// 应用选项（用户手动设置的会覆盖自动检测的）
	for _, opt := range options {
		opt(file)
	}

	// 如果没有设置URL，生成一个占位符（实际上传时会替换）
	if file.URL == "" {
		file.URL = fmt.Sprintf("pending://upload/%s", file.Name)
	}

	w.files = append(w.files, file)
	return nil
}

// CreateTempFile 在临时目录中创建文件路径
func (w *CloudWriter) CreateTempFile(filename string) (string, error) {
	tempDir, err := w.GetTempDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(tempDir, filename), nil
}

// SetSummary 设置文件集合摘要（支持链式调用）
func (w *CloudWriter) SetSummary(summary string) Writer {
	w.summary = summary
	return w
}

// SetLifecycle 设置文件生命周期（支持链式调用）
func (w *CloudWriter) SetLifecycle(lifecycle Lifecycle) Writer {
	w.lifecycle = lifecycle
	return w
}

// SetMetadata 设置元数据（支持链式调用）
func (w *CloudWriter) SetMetadata(key string, value interface{}) Writer {
	if w.metadata == nil {
		w.metadata = make(map[string]interface{})
	}
	w.metadata[key] = value
	return w
}

// AddMetadata 批量添加元数据（支持链式调用）
func (w *CloudWriter) AddMetadata(metadata map[string]interface{}) Writer {
	if w.metadata == nil {
		w.metadata = make(map[string]interface{})
	}
	for k, v := range metadata {
		w.metadata[k] = v
	}
	return w
}

// GetTempDir 获取临时目录
func (w *CloudWriter) GetTempDir() (string, error) {
	if w.tempDir != "" {
		return w.tempDir, nil
	}

	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取工作目录失败: %w", err)
	}

	// 使用traceID作为临时目录名
	traceID := "default"
	if ctx := w.ctx.Value("trace_id"); ctx != nil {
		if id, ok := ctx.(string); ok {
			traceID = id
		}
	}

	// 创建临时目录: ./temp/traceID
	tempBase := filepath.Join(workDir, "temp")
	if err := os.MkdirAll(tempBase, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	w.tempDir = filepath.Join(tempBase, traceID)
	if err := os.MkdirAll(w.tempDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	return w.tempDir, nil
}

// Cleanup 手动清理临时文件和目录
func (w *CloudWriter) Cleanup() {
	if w.tempDir != "" {
		os.RemoveAll(w.tempDir)
		w.tempDir = ""
	}
}

// GetFiles 获取所有文件信息（辅助方法）
func (w *CloudWriter) GetFiles() []*File {
	return w.files
}

// GetSummary 获取文件集合摘要（辅助方法）
func (w *CloudWriter) GetSummary() string {
	return w.summary
}

// GetLifecycle 获取文件生命周期（辅助方法）
func (w *CloudWriter) GetLifecycle() Lifecycle {
	return w.lifecycle
}

// MarshalJSON 实现json.Marshaler接口，在序列化时进行文件上传
func (w *CloudWriter) MarshalJSON() ([]byte, error) {
	// 在这里进行文件上传操作
	err := w.uploadFiles()
	if err != nil {
		return nil, fmt.Errorf("文件上传失败: %w", err)
	}

	// 创建用于序列化的结构体
	// LocalPath字段已设置为json:"-"，会自动被忽略
	result := struct {
		Files     []*File                `json:"files"`
		Summary   string                 `json:"summary,omitempty"`
		Lifecycle Lifecycle              `json:"lifecycle,omitempty"`
		Metadata  map[string]interface{} `json:"metadata,omitempty"`
	}{
		Files:     w.files,
		Summary:   w.summary,
		Lifecycle: w.lifecycle,
		Metadata:  w.metadata,
	}

	// 返回完整对象的JSON
	return json.Marshal(result)
}

// uploadFiles 执行文件上传操作
func (w *CloudWriter) uploadFiles() error {
	for _, file := range w.files {
		if file.LocalPath == "" {
			continue
		}

		// 这里应该调用实际的云存储上传服务
		// 暂时生成一个模拟的URL
		uploadedURL, err := w.uploadToCloud(file.LocalPath, file.Name)
		if err != nil {
			return fmt.Errorf("上传文件 %s 失败: %w", file.Name, err)
		}

		// 更新文件URL
		file.URL = uploadedURL
		file.Status = "uploaded"
		file.UpdatedAt = typex.Time(time.Now())
	}

	return nil
}

// uploadToCloud 使用真实的云存储上传配置进行文件上传
func (w *CloudWriter) uploadToCloud(localPath, filename string) (string, error) {
	// 从context中获取上传配置
	functionMsg := w.getFunctionMsg()
	if functionMsg == nil {
		return "", fmt.Errorf("未找到上传配置信息")
	}

	uploadConfig := functionMsg.UploadConfig

	// 检查文件是否存在
	if _, err := os.Stat(localPath); err != nil {
		return "", fmt.Errorf("文件不存在: %s", localPath)
	}

	// 使用工厂创建上传器
	factory := NewUploaderFactory()
	uploader, err := factory.CreateUploaderWithMsg(uploadConfig, functionMsg)
	if err != nil {
		return "", fmt.Errorf("创建上传器失败: %w", err)
	}

	// 执行上传
	downloadURL, err := uploader.UploadFile(localPath, filename)
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}

	return downloadURL, nil
}

// getFunctionMsg 从context中获取FunctionMsg
func (w *CloudWriter) getFunctionMsg() *trace.FunctionMsg {
	if w.ctx == nil {
		return nil
	}

	value := w.ctx.Value(trace.FunctionMsgKey)
	if value == nil {
		return nil
	}

	if msg, ok := value.(*trace.FunctionMsg); ok {
		return msg
	}

	return nil
}

// uploadFile 执行实际的文件上传操作
func (w *CloudWriter) uploadFile(localPath, filename string, config trace.UploadConfig) (string, error) {
	// 打开文件
	file, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 创建multipart表单
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文件字段
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("创建表单文件字段失败: %w", err)
	}

	// 复制文件内容
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("复制文件内容失败: %w", err)
	}

	// 关闭writer
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("关闭multipart writer失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", config.UploadDomain, &buf)
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 执行请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("执行HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("上传失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var uploadResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			URL      string `json:"url"`
			Filename string `json:"filename"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return "", fmt.Errorf("解析上传响应失败: %w", err)
	}

	if !uploadResp.Success {
		return "", fmt.Errorf("上传失败: %s", uploadResp.Message)
	}

	// 构建完整的下载URL
	downloadURL := uploadResp.Data.URL
	if config.DownloadDomain != "" && downloadURL != "" {
		// 如果返回的是相对路径，则拼接下载域名
		if downloadURL[0] == '/' {
			downloadURL = config.DownloadDomain + downloadURL
		}
	}

	return downloadURL, nil
}

// detectContentType 从文件路径检测MIME类型
func (w *CloudWriter) detectContentType(filePath string) (string, error) {
	mtype, err := mimetype.DetectFile(filePath)
	if err != nil {
		return "", err
	}

	// 为文本文件添加UTF-8编码声明（如果还没有的话）
	mimeType := mtype.String()
	if (mtype.Is("text/plain") || mtype.Is("text/html") || mtype.Is("text/xml") ||
		mtype.Is("application/json") || mtype.Is("text/css") || mtype.Is("text/javascript")) &&
		!strings.Contains(mimeType, "charset") {
		mimeType += "; charset=utf-8"
	}

	return mimeType, nil
}

// detectContentTypeFromData 从数据内容检测MIME类型
func (w *CloudWriter) detectContentTypeFromData(data []byte) string {
	mtype := mimetype.Detect(data)

	// 为文本文件添加UTF-8编码声明（如果还没有的话）
	mimeType := mtype.String()
	if (mtype.Is("text/plain") || mtype.Is("text/html") || mtype.Is("text/xml") ||
		mtype.Is("application/json") || mtype.Is("text/css") || mtype.Is("text/javascript")) &&
		!strings.Contains(mimeType, "charset") {
		mimeType += "; charset=utf-8"
	}

	return mimeType
}
