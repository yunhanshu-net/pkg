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

type UploadConfig struct {
	UploadDomain   string `json:"upload_domain"`   //上传地址
	DownloadDomain string `json:"download_domain"` //上传后下载的域名
	UploadToken    string `json:"upload_token"`    //上传Token（七牛云等需要）
	Bucket         string `json:"bucket"`          //存储桶名称
	AccessKey      string `json:"access_key"`      //访问密钥
	SecretKey      string `json:"secret_key"`      //私钥
	Provider       string `json:"provider"`        //存储提供商（qiniu/aliyun/aws等）
}
