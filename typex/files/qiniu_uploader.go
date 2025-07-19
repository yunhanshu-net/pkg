package files

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/yunhanshu-net/pkg/trace"
)

// QiniuUploader 七牛云上传器
type QiniuUploader struct {
	config      trace.UploadConfig
	functionMsg *trace.FunctionMsg // 添加FunctionMsg用于获取user和runner信息
}

// NewQiniuUploader 创建七牛云上传器
func NewQiniuUploader(config trace.UploadConfig) *QiniuUploader {
	return &QiniuUploader{
		config: config,
	}
}

// NewQiniuUploaderWithMsg 创建带有FunctionMsg的七牛云上传器
func NewQiniuUploaderWithMsg(config trace.UploadConfig, functionMsg *trace.FunctionMsg) *QiniuUploader {
	return &QiniuUploader{
		config:      config,
		functionMsg: functionMsg,
	}
}

// UploadFile 上传文件到七牛云
func (q *QiniuUploader) UploadFile(localPath, filename string) (string, error) {
	// 打开本地文件
	file, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 获取文件信息
	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 生成文件Key（存储路径）
	fileKey := q.generateFileKeyWithPath(localPath, filename)

	// 获取上传Token
	upToken := q.getUploadToken()

	// 配置七牛云
	cfg := &storage.Config{
		UseHTTPS:      true,
		UseCdnDomains: false,
	}

	// 创建表单上传器
	formUploader := storage.NewFormUploader(cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": filename,
		},
	}

	// 执行上传
	err = formUploader.Put(context.Background(), &ret, upToken, fileKey, file, stat.Size(), &putExtra)
	if err != nil {
		return "", fmt.Errorf("七牛云上传失败: %w", err)
	}

	// 构建下载URL
	downloadURL := q.buildDownloadURL(ret.Key)
	return downloadURL, nil
}

// getUploadToken 获取上传Token
func (q *QiniuUploader) getUploadToken() string {
	// 如果配置中已有Token，直接使用
	if q.config.UploadToken != "" {
		return q.config.UploadToken
	}

	// 否则使用AccessKey和SecretKey生成Token
	putPolicy := storage.PutPolicy{
		Scope: q.config.Bucket,
	}
	mac := qbox.NewMac(q.config.AccessKey, q.config.SecretKey)
	return putPolicy.UploadToken(mac)
}

// generateFileKeyWithPath 生成文件存储Key（带本地路径，现在使用UUID）
func (q *QiniuUploader) generateFileKeyWithPath(localPath, filename string) string {
	// 如果有FunctionMsg信息，使用规范的上传路径
	if q.functionMsg != nil {
		// 使用FunctionMsg.GetUploadPath()获取规范路径
		// 格式：租户/应用/函数/方法/output/日期
		basePath := q.functionMsg.GetUploadPath()

		// 生成UUID作为唯一值（性能更好，避免MD5计算）
		uniqueID := uuid.New().String()

		// 组合完整路径：租户/应用/函数/方法/output/日期/UUID/原文件名
		fullPath := fmt.Sprintf("%s/%s/%s", basePath, uniqueID, filename)

		return fullPath
	}

	// 兜底：如果没有FunctionMsg信息，使用简单的时间戳路径
	timestamp := time.Now().Format("2006/01/02/15-04-05")
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	uniqueId := time.Now().UnixNano()
	return fmt.Sprintf("uploads/%s/%s_%d%s", timestamp, name, uniqueId, ext)
}

// generateFileKey 生成文件存储Key（兼容旧接口）
func (q *QiniuUploader) generateFileKey(filename string) string {
	// 如果有FunctionMsg信息，使用规范的上传路径
	if q.functionMsg != nil {
		// 使用FunctionMsg.GetUploadPath()获取规范路径
		// 格式：租户/应用/函数/方法/output/日期
		basePath := q.functionMsg.GetUploadPath()

		// 生成UUID作为唯一值（性能更好）
		uniqueID := uuid.New().String()

		// 组合完整路径：租户/应用/函数/方法/output/日期/UUID/原文件名
		fullPath := fmt.Sprintf("%s/%s/%s", basePath, uniqueID, filename)

		return fullPath
	}

	// 兜底：如果没有FunctionMsg信息，使用简单的时间戳路径
	timestamp := time.Now().Format("2006/01/02/15-04-05")
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	uniqueId := time.Now().UnixNano()
	return fmt.Sprintf("uploads/%s/%s_%d%s", timestamp, name, uniqueId, ext)
}

// buildDownloadURL 构建下载URL
func (q *QiniuUploader) buildDownloadURL(key string) string {
	if q.config.DownloadDomain == "" {
		return key // 如果没有下载域名，返回Key
	}

	// 确保下载域名以http://或https://开头
	domain := q.config.DownloadDomain
	if domain[len(domain)-1] == '/' {
		domain = domain[:len(domain)-1]
	}

	return fmt.Sprintf("%s/%s", domain, key)
}

// GenerateFileKey 公开的路径生成方法（用于测试）
func (q *QiniuUploader) GenerateFileKey(filename string) string {
	return q.generateFileKey(filename)
}
