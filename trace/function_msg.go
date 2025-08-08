package trace

import (
	"fmt"
	"strings"
	"time"
)

const (
	FunctionMsgKey = "FunctionMsg"
)

type FunctionMsg struct {
	UploadConfig UploadConfig `json:"upload_config"` //上传配置
	User         string       `json:"user"`          //app所属租户
	Version      string       `json:"version"`       //当前版本
	Runner       string       `json:"runner"`        //app
	TraceID      string       `json:"trace_id"`      //trace_id
	Method       string       `json:"method"`
	Router       string       `json:"router"`
}

func (f *FunctionMsg) GetUploadPath() string {
	//每个函数的文件上传目录都不一样，前端上传的属于input，后端处理后上传返回给前端的属于output
	//租户/应用/函数/方法/output/日期/文件key/文件名称 按照这种格式来，这样可以防止文件重复
	s := fmt.Sprintf("%s/%s/%s/%s/output/%s", f.User, f.Runner, f.Router, f.Method,
		time.Now().Format("20060102"))
	return strings.ReplaceAll(s, "//", "/")
}

// ToMap 返回FunctionMsg的字段作为map，方便日志系统使用
// includeFields: 指定要包含的字段，如果为空则返回所有字段
func (f *FunctionMsg) ToMap(includeFields ...string) map[string]interface{} {
	if f == nil {
		return nil
	}
	
	// 定义所有可用字段
	allFields := map[string]interface{}{
		"trace_id": f.TraceID,
		"method":   f.Method,
		"router":   f.Router,
		"user":     f.User,
		"version":  f.Version,
		"runner":   f.Runner,
		// UploadConfig相关字段（如果需要的话）
		// "upload_provider": f.UploadConfig.Provider,
		// "upload_bucket":   f.UploadConfig.Bucket,
	}
	
	// 如果没有指定字段，返回所有字段
	if len(includeFields) == 0 {
		return allFields
	}
	
	// 如果指定了字段，只返回指定的字段
	result := make(map[string]interface{})
	for _, fieldName := range includeFields {
		if value, exists := allFields[fieldName]; exists {
			result[fieldName] = value
		}
	}
	
	return result
}

// ToLogMap 专门为日志系统提供的方法，返回适合日志的字段
func (f *FunctionMsg) ToLogMap() map[string]interface{} {
	// 返回完整的日志字段，包含 router、method 等参数
	return f.ToMap("trace_id", "method", "router", "user", "version", "runner")
}

type UploadConfig struct {
	UploadDomain   string `json:"upload_domain"`   //上传地址
	DownloadDomain string `json:"download_domain"` //上传后下载的域名
	UploadToken    string `json:"upload_token"`    //上传Token（七牛云等需要）
	Bucket         string `json:"bucket"`          //存储桶名称
	AccessKey      string `json:"access_key"`      //访问密钥
	SecretKey      string `json:"secret_key"`      //私钥
	Provider       string `json:"provider"`        //存储提供商（qiniu/aliyun/aws等）
}
