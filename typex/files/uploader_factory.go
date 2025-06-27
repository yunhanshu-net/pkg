package files

import (
	"fmt"

	"github.com/yunhanshu-net/pkg/trace"
)

// Uploader 上传器接口
type Uploader interface {
	UploadFile(localPath, filename string) (string, error)
}

// UploaderFactory 上传器工厂
type UploaderFactory struct{}

// NewUploaderFactory 创建上传器工厂
func NewUploaderFactory() *UploaderFactory {
	return &UploaderFactory{}
}

// CreateUploader 根据配置创建对应的上传器
func (f *UploaderFactory) CreateUploader(config trace.UploadConfig) (Uploader, error) {
	return f.CreateUploaderWithMsg(config, nil)
}

// CreateUploaderWithMsg 根据配置和FunctionMsg创建对应的上传器
func (f *UploaderFactory) CreateUploaderWithMsg(config trace.UploadConfig, functionMsg *trace.FunctionMsg) (Uploader, error) {
	switch config.Provider {
	case "qiniu":
		if config.Bucket == "" {
			return nil, fmt.Errorf("七牛云上传需要配置Bucket")
		}
		if config.AccessKey == "" || config.SecretKey == "" {
			if config.UploadToken == "" {
				return nil, fmt.Errorf("七牛云上传需要配置AccessKey/SecretKey或UploadToken")
			}
		}
		// 如果有FunctionMsg，使用带消息的构造函数
		if functionMsg != nil {
			return NewQiniuUploaderWithMsg(config, functionMsg), nil
		}
		return NewQiniuUploader(config), nil

	case "aliyun":
		// TODO: 实现阿里云OSS上传器
		return nil, fmt.Errorf("阿里云OSS上传器暂未实现")

	case "aws":
		// TODO: 实现AWS S3上传器
		return nil, fmt.Errorf("AWS S3上传器暂未实现")

	case "http", "":
		// 默认HTTP multipart上传
		if config.UploadDomain == "" {
			return nil, fmt.Errorf("HTTP上传需要配置UploadDomain")
		}
		return NewHTTPUploader(config), nil

	default:
		return nil, fmt.Errorf("不支持的上传提供商: %s", config.Provider)
	}
}

// HTTPUploader HTTP上传器（原有的multipart实现）
type HTTPUploader struct {
	config trace.UploadConfig
}

// NewHTTPUploader 创建HTTP上传器
func NewHTTPUploader(config trace.UploadConfig) *HTTPUploader {
	return &HTTPUploader{
		config: config,
	}
}

// UploadFile 通过HTTP上传文件
func (h *HTTPUploader) UploadFile(localPath, filename string) (string, error) {
	// 这里可以复用CloudWriter中的uploadFile方法
	// 为了简化，我们创建一个临时的CloudWriter实例
	writer := &CloudWriter{}
	return writer.uploadFile(localPath, filename, h.config)
}
