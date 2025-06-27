package files

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCloudWriter_AddFile(t *testing.T) {
	ctx := context.Background()
	writer := NewCloudWriter(ctx)
	defer writer.Cleanup()

	// 创建测试文件
	tempDir, err := os.MkdirTemp("", "cloud-writer-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("Hello, World!"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 测试添加文件
	err = writer.AddFile(testFile,
		WithContentType("text/plain"),
		WithDescription("测试文件"),
		WithMetadata("type", "greeting"),
	)
	if err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	// 验证文件是否被添加
	cw := writer.(*CloudWriter)
	files := cw.files
	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	if files[0].Name != "test.txt" {
		t.Errorf("Expected file name 'test.txt', got '%s'", files[0].Name)
	}

	if files[0].Description != "测试文件" {
		t.Errorf("Expected description '测试文件', got '%s'", files[0].Description)
	}

	if files[0].Metadata["type"] != "greeting" {
		t.Errorf("Expected metadata type 'greeting', got '%s'", files[0].Metadata["type"])
	}
}

func TestCloudWriter_AddFileWithData(t *testing.T) {
	ctx := context.Background()
	writer := NewCloudWriter(ctx)
	defer writer.Cleanup()

	// 测试从数据添加文件
	testData := []byte("Hello, World!")
	err := writer.AddFileWithData("test.txt", testData,
		WithContentType("text/plain"),
		WithDescription("测试文件"),
		WithMetadata("type", "greeting"),
	)
	if err != nil {
		t.Fatalf("AddFileWithData failed: %v", err)
	}

	// 验证文件是否被添加
	cw := writer.(*CloudWriter)
	files := cw.files
	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	if files[0].Name != "test.txt" {
		t.Errorf("Expected file name 'test.txt', got '%s'", files[0].Name)
	}

	// 验证文件内容
	tempDir, err := writer.GetTempDir()
	if err != nil {
		t.Fatalf("GetTempDir failed: %v", err)
	}

	filePath := filepath.Join(tempDir, "test.txt")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(content) != "Hello, World!" {
		t.Errorf("Expected content 'Hello, World!', got '%s'", string(content))
	}
}

func TestCloudWriter_CreateTempFile(t *testing.T) {
	ctx := context.Background()
	writer := NewCloudWriter(ctx)
	defer writer.Cleanup()

	// 测试创建临时文件
	filePath, err := writer.CreateTempFile("test.txt")
	if err != nil {
		t.Fatalf("CreateTempFile failed: %v", err)
	}

	// 验证文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Created file does not exist")
	}

	// 验证文件路径
	expectedDir, err := writer.GetTempDir()
	if err != nil {
		t.Fatalf("GetTempDir failed: %v", err)
	}

	expectedPath := filepath.Join(expectedDir, "test.txt")
	if filePath != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, filePath)
	}
}

func TestCloudWriter_SetSummary(t *testing.T) {
	ctx := context.Background()
	writer := NewCloudWriter(ctx)
	defer writer.Cleanup()

	// 测试设置摘要
	summary := "测试文件集合"
	writer.SetSummary(summary)

	// 验证摘要是否被设置
	cw := writer.(*CloudWriter)
	if cw.summary != summary {
		t.Errorf("Expected summary '%s', got '%s'", summary, cw.summary)
	}
}

func TestCloudWriter_SetLifecycle(t *testing.T) {
	ctx := context.Background()
	writer := NewCloudWriter(ctx)
	defer writer.Cleanup()

	// 测试设置生命周期
	lifecycle := LifecycleShortTerm
	writer.SetLifecycle(lifecycle)

	// 验证生命周期是否被设置
	cw := writer.(*CloudWriter)
	if cw.lifecycle != lifecycle {
		t.Errorf("Expected lifecycle '%s', got '%s'", lifecycle, cw.lifecycle)
	}
}

func TestCloudWriter_Cleanup(t *testing.T) {
	ctx := context.Background()
	writer := NewCloudWriter(ctx)

	// 获取临时目录
	tempDir, err := writer.GetTempDir()
	if err != nil {
		t.Fatalf("GetTempDir failed: %v", err)
	}

	// 创建一些测试文件
	testFiles := []string{"test1.txt", "test2.txt", "test3.txt"}
	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		err := os.WriteFile(filePath, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
	}

	// 执行清理
	writer.Cleanup()

	// 验证临时目录已被删除
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Error("Temp directory was not cleaned up")
	}
}
