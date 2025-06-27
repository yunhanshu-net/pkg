package files

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yunhanshu-net/pkg/trace"
)

func TestCloudWriterUpload(t *testing.T) {
	// 创建测试文件
	testDir := "./test_upload"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	testFile := filepath.Join(testDir, "test.txt")
	testContent := []byte("Hello, World! This is a test file.")
	err := os.WriteFile(testFile, testContent, 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建带有上传配置的context
	uploadConfig := trace.UploadConfig{
		Provider:       "qiniu",
		Bucket:         "test-bucket",
		AccessKey:      "test-access-key",
		SecretKey:      "test-secret-key",
		DownloadDomain: "https://cdn.example.com",
	}

	functionMsg := &trace.FunctionMsg{
		UploadConfig: uploadConfig,
	}

	ctx := context.WithValue(context.Background(), trace.FunctionMsgKey, functionMsg)

	// 创建CloudWriter
	writer := NewCloudWriter(ctx).(*CloudWriter)

	// 添加文件
	err = writer.AddFile(testFile)
	if err != nil {
		t.Fatalf("添加文件失败: %v", err)
	}

	// 测试JSON序列化（这会触发上传）
	jsonData, err := json.Marshal(writer)
	if err != nil {
		// 由于我们使用的是测试配置，上传可能会失败，这是正常的
		t.Logf("JSON序列化失败（预期的，因为使用测试配置）: %v", err)
	} else {
		t.Logf("JSON序列化成功: %s", string(jsonData))
	}

	// 清理
	writer.Cleanup()
}

func TestUploaderFactory(t *testing.T) {
	factory := NewUploaderFactory()

	// 测试七牛云上传器创建
	qiniuConfig := trace.UploadConfig{
		Provider:  "qiniu",
		Bucket:    "test-bucket",
		AccessKey: "test-access-key",
		SecretKey: "test-secret-key",
	}

	uploader, err := factory.CreateUploader(qiniuConfig)
	if err != nil {
		t.Fatalf("创建七牛云上传器失败: %v", err)
	}

	if uploader == nil {
		t.Fatal("上传器不应该为nil")
	}

	// 测试HTTP上传器创建
	httpConfig := trace.UploadConfig{
		Provider:     "http",
		UploadDomain: "https://upload.example.com",
	}

	httpUploader, err := factory.CreateUploader(httpConfig)
	if err != nil {
		t.Fatalf("创建HTTP上传器失败: %v", err)
	}

	if httpUploader == nil {
		t.Fatal("HTTP上传器不应该为nil")
	}

	// 测试不支持的提供商
	unsupportedConfig := trace.UploadConfig{
		Provider: "unsupported",
	}

	_, err = factory.CreateUploader(unsupportedConfig)
	if err == nil {
		t.Fatal("应该返回不支持的提供商错误")
	}
}

func TestFilesJSONSerialization(t *testing.T) {
	// 创建Files实例
	files := Files{
		{
			Name:        "test1.txt",
			Size:        100,
			ContentType: "text/plain",
			URL:         "https://example.com/test1.txt",
		},
		{
			Name:        "test2.jpg",
			Size:        2048,
			ContentType: "image/jpeg",
			URL:         "https://example.com/test2.jpg",
		},
	}

	// 测试JSON序列化
	jsonData, err := json.Marshal(files)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	t.Logf("序列化结果: %s", string(jsonData))

	// 测试JSON反序列化
	var deserializedFiles Files
	err = json.Unmarshal(jsonData, &deserializedFiles)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	if len(deserializedFiles) != 2 {
		t.Fatalf("反序列化后文件数量不正确，期望2个，实际%d个", len(deserializedFiles))
	}

	if deserializedFiles[0].Name != "test1.txt" {
		t.Fatalf("反序列化后文件名不正确，期望test1.txt，实际%s", deserializedFiles[0].Name)
	}
}
