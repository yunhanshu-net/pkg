package trace

const (
	FunctionMsgKey = "FunctionMsg"
)

type FunctionMsg struct {
	UploadConfig UploadConfig `json:"upload_config"` //上传配置
	User         string       `json:"user"`          //app所属租户
	Version      string       `json:"version"`       //当前版本
	Runner       string       `json:"runner"`        //app
	TraceID      string       `json:"trace_id"`      //trace_id
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
