package files

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestURLReader_AddFileFromURL(t *testing.T) {
	ctx := context.Background()
	reader := NewURLReader(ctx)
	defer reader.Cleanup()

	tests := []struct {
		name        string
		url         string
		fileName    string
		options     []FileOption
		expectError bool
	}{
		{
			name:     "有效的JSON文件",
			url:      "https://httpbin.org/json",
			fileName: "test.json",
			options: []FileOption{
				WithContentType("application/json"),
				WithDescription("测试JSON文件"),
			},
			expectError: false,
		},
		{
			name:        "无效的URL",
			url:         "https://invalid-url-that-does-not-exist.com",
			fileName:    "test.txt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := reader.AddFileFromURL(tt.url, tt.fileName, tt.options...)
			if (err != nil) != tt.expectError {
				t.Errorf("AddFileFromURL() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestURLReader_DownloadFile(t *testing.T) {
	ctx := context.Background()
	reader := NewURLReader(ctx)
	defer reader.Cleanup()

	// 添加测试文件
	err := reader.AddFileFromURL(
		"https://httpbin.org/json",
		"test.json",
		WithContentType("application/json"),
	)
	if err != nil {
		t.Fatalf("AddFileFromURL failed: %v", err)
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "url-reader-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 测试下载文件
	downloadPath := filepath.Join(tempDir, "downloaded.json")
	_, err = reader.DownloadFile(0, downloadPath)
	if err != nil {
		t.Fatalf("DownloadFile failed: %v", err)
	}

	// 验证文件是否存在
	if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
		t.Error("Downloaded file does not exist")
	}
}

func TestURLReader_DownloadAll(t *testing.T) {
	ctx := context.Background()
	reader := NewURLReader(ctx)
	defer reader.Cleanup()

	// 添加多个测试文件
	testFiles := []struct {
		url      string
		fileName string
	}{
		{"https://httpbin.org/json", "test1.json"},
		{"https://httpbin.org/bytes/100", "test2.bin"},
	}

	for _, tf := range testFiles {
		err := reader.AddFileFromURL(tf.url, tf.fileName)
		if err != nil {
			t.Fatalf("AddFileFromURL failed for %s: %v", tf.fileName, err)
		}
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "url-reader-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 测试下载所有文件
	_, err = reader.DownloadAll(tempDir)
	if err != nil {
		t.Fatalf("DownloadAll failed: %v", err)
	}

	// 验证所有文件是否都已下载
	for _, tf := range testFiles {
		filePath := filepath.Join(tempDir, tf.fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("File %s was not downloaded", tf.fileName)
		}
	}
}

func TestURLReader_FilterByType(t *testing.T) {
	ctx := context.Background()
	reader := NewURLReader(ctx)
	defer reader.Cleanup()

	// 添加不同类型的文件
	testFiles := []struct {
		url         string
		fileName    string
		contentType string
	}{
		{"https://httpbin.org/json", "test1.json", "application/json"},
		{"https://httpbin.org/bytes/100", "test2.bin", "application/octet-stream"},
		{"https://httpbin.org/stream/1", "test3.txt", "text/plain"},
	}

	for _, tf := range testFiles {
		err := reader.AddFileFromURL(tf.url, tf.fileName, WithContentType(tf.contentType))
		if err != nil {
			t.Fatalf("AddFileFromURL failed for %s: %v", tf.fileName, err)
		}
	}

	// 测试按类型过滤
	jsonFiles := reader.FilterByType("application/json")
	if len(jsonFiles) != 1 {
		t.Errorf("Expected 1 JSON file, got %d", len(jsonFiles))
	}

	textFiles := reader.FilterByType("text/plain")
	if len(textFiles) != 1 {
		t.Errorf("Expected 1 text file, got %d", len(textFiles))
	}
}

func TestURLReader_FilterByExtension(t *testing.T) {
	ctx := context.Background()
	reader := NewURLReader(ctx)
	defer reader.Cleanup()

	// 添加不同扩展名的文件
	testFiles := []struct {
		url      string
		fileName string
	}{
		{"https://httpbin.org/json", "test1.json"},
		{"https://httpbin.org/bytes/100", "test2.bin"},
		{"https://httpbin.org/stream/1", "test3.txt"},
	}

	for _, tf := range testFiles {
		err := reader.AddFileFromURL(tf.url, tf.fileName)
		if err != nil {
			t.Fatalf("AddFileFromURL failed for %s: %v", tf.fileName, err)
		}
	}

	// 测试按扩展名过滤
	jsonFiles := reader.FilterByExtension(".json")
	if len(jsonFiles) != 1 {
		t.Errorf("Expected 1 JSON file, got %d", len(jsonFiles))
	}

	txtFiles := reader.FilterByExtension(".txt")
	if len(txtFiles) != 1 {
		t.Errorf("Expected 1 TXT file, got %d", len(txtFiles))
	}
}

func TestURLReader_GetTotalSize(t *testing.T) {
	ctx := context.Background()
	reader := NewURLReader(ctx)
	defer reader.Cleanup()

	// 添加测试文件
	err := reader.AddFileFromURL(
		"https://httpbin.org/bytes/100",
		"test.bin",
		WithSize(100),
	)
	if err != nil {
		t.Fatalf("AddFileFromURL failed: %v", err)
	}

	// 测试获取总大小
	totalSize := reader.GetTotalSize()
	if totalSize != 100 {
		t.Errorf("Expected total size 100, got %d", totalSize)
	}
}

func TestURLReader_Cleanup(t *testing.T) {
	ctx := context.Background()
	reader := NewURLReader(ctx)

	// 获取临时目录
	tempDir, err := reader.GetTempDir()
	if err != nil {
		t.Fatalf("GetTempDir failed: %v", err)
	}
	if tempDir == "" {
		t.Fatal("GetTempDir returned empty string")
	}

	// 验证临时目录存在
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Fatal("Temp directory does not exist")
	}

	// 执行清理
	reader.Cleanup()

	// 验证临时目录已被删除
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Error("Temp directory was not cleaned up")
	}
}

func TestFileTypeFunctions(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    string
	}{
		{"JPEG图片", "image/jpeg", FileTypeImage},
		{"PNG图片", "image/png", FileTypeImage},
		{"PDF文档", "application/pdf", FileTypeDocument},
		{"Word文档", "application/msword", FileTypeDocument},
		{"MP4视频", "video/mp4", FileTypeVideo},
		{"MP3音频", "audio/mpeg", FileTypeAudio},
		{"ZIP压缩包", "application/zip", FileTypeArchive},
		{"未知类型", "application/unknown", FileTypeOther},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileType := GetFileType(tt.contentType)
			if fileType != tt.expected {
				t.Errorf("GetFileType(%s) = %s, want %s", tt.contentType, fileType, tt.expected)
			}
		})
	}
}

func TestFileTypePredicates(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		isImage     bool
		isDocument  bool
		isVideo     bool
		isAudio     bool
		isArchive   bool
	}{
		{
			name:        "JPEG图片",
			contentType: "image/jpeg",
			isImage:     true,
			isDocument:  false,
			isVideo:     false,
			isAudio:     false,
			isArchive:   false,
		},
		{
			name:        "PDF文档",
			contentType: "application/pdf",
			isImage:     false,
			isDocument:  true,
			isVideo:     false,
			isAudio:     false,
			isArchive:   false,
		},
		{
			name:        "MP4视频",
			contentType: "video/mp4",
			isImage:     false,
			isDocument:  false,
			isVideo:     true,
			isAudio:     false,
			isArchive:   false,
		},
		{
			name:        "MP3音频",
			contentType: "audio/mpeg",
			isImage:     false,
			isDocument:  false,
			isVideo:     false,
			isAudio:     true,
			isArchive:   false,
		},
		{
			name:        "ZIP压缩包",
			contentType: "application/zip",
			isImage:     false,
			isDocument:  false,
			isVideo:     false,
			isAudio:     false,
			isArchive:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsImage(tt.contentType) != tt.isImage {
				t.Errorf("IsImage(%s) = %v, want %v", tt.contentType, IsImage(tt.contentType), tt.isImage)
			}
			if IsDocument(tt.contentType) != tt.isDocument {
				t.Errorf("IsDocument(%s) = %v, want %v", tt.contentType, IsDocument(tt.contentType), tt.isDocument)
			}
			if IsVideo(tt.contentType) != tt.isVideo {
				t.Errorf("IsVideo(%s) = %v, want %v", tt.contentType, IsVideo(tt.contentType), tt.isVideo)
			}
			if IsAudio(tt.contentType) != tt.isAudio {
				t.Errorf("IsAudio(%s) = %v, want %v", tt.contentType, IsAudio(tt.contentType), tt.isAudio)
			}
			if IsArchive(tt.contentType) != tt.isArchive {
				t.Errorf("IsArchive(%s) = %v, want %v", tt.contentType, IsArchive(tt.contentType), tt.isArchive)
			}
		})
	}
}
