package files

import (
	"context"
	"testing"
)

func TestNewReader(t *testing.T) {
	ctx := context.Background()

	// 测试默认Reader
	reader := NewReader(ctx)
	if reader == nil {
		t.Fatal("NewReader returned nil")
	}

	// 测试URLReader
	urlReader := NewURLReader(ctx)
	if urlReader == nil {
		t.Fatal("NewURLReader returned nil")
	}

	// 测试接口实现
	var _ Reader = reader
	var _ Reader = urlReader
}

func TestNewWriter(t *testing.T) {
	ctx := context.Background()

	// 测试默认Writer
	writer := NewWriter(ctx)
	if writer == nil {
		t.Fatal("NewWriter returned nil")
	}

	// 测试CloudWriter
	cloudWriter := NewCloudWriter(ctx)
	if cloudWriter == nil {
		t.Fatal("NewCloudWriter returned nil")
	}

	// 测试接口实现
	var _ Writer = writer
	var _ Writer = cloudWriter
}

func TestGetReader(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		readerType string
		expected   string
	}{
		{"url", "URLReader"},
		{"http", "URLReader"},
		{"https", "URLReader"},
		{"unknown", "URLReader"}, // 默认类型
	}

	for _, tt := range tests {
		t.Run(tt.readerType, func(t *testing.T) {
			reader := GetReader(tt.readerType, ctx)
			if reader == nil {
				t.Fatalf("GetReader(%s) returned nil", tt.readerType)
			}
		})
	}
}

func TestGetWriter(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		writerType string
		expected   string
	}{
		{"cloud", "CloudWriter"},
		{"qiniu", "CloudWriter"},
		{"oss", "CloudWriter"},
		{"unknown", "CloudWriter"}, // 默认类型
	}

	for _, tt := range tests {
		t.Run(tt.writerType, func(t *testing.T) {
			writer := GetWriter(tt.writerType, ctx)
			if writer == nil {
				t.Fatalf("GetWriter(%s) returned nil", tt.writerType)
			}
		})
	}
}

func TestFileOptions(t *testing.T) {
	file := &File{}

	// 测试选项函数
	WithSize(1024)(file)
	if file.Size != 1024 {
		t.Errorf("WithSize failed: expected 1024, got %d", file.Size)
	}

	WithContentType("text/plain")(file)
	if file.ContentType != "text/plain" {
		t.Errorf("WithContentType failed: expected 'text/plain', got '%s'", file.ContentType)
	}

	WithHash("abc123")(file)
	if file.Hash != "abc123" {
		t.Errorf("WithHash failed: expected 'abc123', got '%s'", file.Hash)
	}

	WithDescription("test file")(file)
	if file.Description != "test file" {
		t.Errorf("WithDescription failed: expected 'test file', got '%s'", file.Description)
	}

	WithMetadata("key", "value")(file)
	if file.Metadata["key"] != "value" {
		t.Errorf("WithMetadata failed: expected 'value', got '%s'", file.Metadata["key"])
	}

	WithAutoDelete(true)(file)
	if !file.AutoDelete {
		t.Error("WithAutoDelete failed: expected true")
	}

	WithCompressed(true)(file)
	if !file.Compressed {
		t.Error("WithCompressed failed: expected true")
	}
}

func TestLifecycleConstants(t *testing.T) {
	tests := []struct {
		lifecycle Lifecycle
		expected  string
	}{
		{LifecycleTemporary, "temporary"},
		{LifecycleShortTerm, "short_term"},
		{LifecycleLongTerm, "long_term"},
		{LifecycleCache, "cache"},
	}

	for _, tt := range tests {
		if string(tt.lifecycle) != tt.expected {
			t.Errorf("Lifecycle constant mismatch: expected '%s', got '%s'", tt.expected, string(tt.lifecycle))
		}
	}
}

// 集成测试示例
func TestReaderWriterIntegration(t *testing.T) {
	ctx := context.Background()

	// 创建Reader和Writer
	reader := NewReader(ctx)
	writer := NewWriter(ctx)

	// 确保清理
	defer func() {
		reader.Cleanup()
		writer.Cleanup()
	}()

	// 测试基本功能
	reader.SetSummary("Test files")
	writer.SetSummary("Output files")
	writer.SetLifecycle(LifecycleTemporary)

	// 验证接口方法存在
	files := reader.GetFiles()
	if files == nil {
		t.Error("GetFiles returned nil")
	}

	totalSize := reader.GetTotalSize()
	if totalSize < 0 {
		t.Error("GetTotalSize returned negative value")
	}
}

func TestURLReaderBasicFunctionality(t *testing.T) {
	ctx := context.Background()
	reader := NewURLReader(ctx)
	defer reader.Cleanup()

	// 测试添加文件
	err := reader.AddFileFromURL(
		"https://httpbin.org/json",
		"test.json",
		WithContentType("application/json"),
		WithDescription("测试JSON文件"),
	)
	if err != nil {
		t.Fatalf("AddFileFromURL failed: %v", err)
	}

	// 测试获取文件列表
	files := reader.GetFiles()
	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	if files[0].Name != "test.json" {
		t.Errorf("Expected file name 'test.json', got '%s'", files[0].Name)
	}

	// 测试获取单个文件
	file, err := reader.GetFile(0)
	if err != nil {
		t.Fatalf("GetFile failed: %v", err)
	}

	if file.Description != "测试JSON文件" {
		t.Errorf("Expected description '测试JSON文件', got '%s'", file.Description)
	}
}

func TestCloudWriterBasicFunctionality(t *testing.T) {
	ctx := context.Background()
	writer := NewCloudWriter(ctx)
	defer writer.Cleanup()

	// 测试创建临时文件
	_, err := writer.CreateTempFile("test.txt")
	if err != nil {
		t.Fatalf("CreateTempFile failed: %v", err)
	}

	// 测试从数据添加文件
	testData := []byte("Hello, World!")
	err = writer.AddFileWithData("hello.txt", testData,
		WithDescription("测试文件"),
		WithMetadata("type", "greeting"),
	)
	if err != nil {
		t.Fatalf("AddFileWithData failed: %v", err)
	}

	// 测试生命周期设置
	writer.SetLifecycle(LifecycleShortTerm)

	// 测试摘要设置
	writer.SetSummary("测试文件集合")
}
