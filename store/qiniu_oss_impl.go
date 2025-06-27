package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"encoding/base64"
	"encoding/hex"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/yunhanshu-net/pkg/x/httpx"
)

func NewDefaultQiNiu() *QiNiu {

	return &QiNiu{
		Domain:           "http://cdn.geeleo.com",
		DownloadRootPath: os.TempDir(),
		Bucket:           "geeleo",
		AccessKey:        "ehF_E4x_EyO_wSN_nwqExyhXPe5hGl5Xjo89_cZ6",
		SecretKey:        "FjfIpqUevEcVx9bQxdgiuX9Di-CUOrKFkR88CZAj",
	}
}

type QiNiu struct {
	DownloadRootPath string `json:"download_root_path"` //下载文件的默认本地根路径
	Domain           string `json:"domain"`             //云存储的域名
	Bucket           string `json:"bucket"`             //空间
	AccessKey        string `json:"access_key"`
	SecretKey        string `json:"secret_key"`
	ImgPath          string `json:"img_path"` //存储路径
}

func (q *QiNiu) GetFullPath(savePath string) string {
	return q.Domain + savePath
}

func qiniuConfig() *storage.Config {
	cfg := storage.Config{UseHTTPS: false, UseCdnDomains: false}
	return &cfg
}

//type File struct {
//	Path     string `json:"path"`
//	Filename string `json:"filename"`
//}

// FileSave 上传本地文件到七牛云
func (q *QiNiu) FileSave(localFilePath string, ossPath string) (*FileSaveInfo, error) {
	ossPath = strings.TrimPrefix(ossPath, "\\")
	ossPath = strings.TrimPrefix(ossPath, "/")
	file, err := os.Open(localFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	name := file.Name()
	fileType := filepath.Ext(name)
	if len(fileType) > 0 {
		fileType = fileType[1:] // 去掉点号
	}
	putPolicy := storage.PutPolicy{Scope: q.Bucket}
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := qiniuConfig()
	formUploader := storage.NewFormUploader(cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{Params: map[string]string{"x:name": "github logo"}}

	fileKey := fmt.Sprintf("%s", ossPath) // 文件名格式 自己可以改 建议保证唯一性
	putErr := formUploader.PutFile(context.Background(), &ret, upToken, fileKey, localFilePath, &putExtra)
	if putErr != nil {
		return nil, putErr
	}

	// 七牛云返回的Hash字段是七牛云自己的ETag算法，不是标准MD5
	// ETag算法说明：
	// - 小文件(≤4M): hash = UrlsafeBase64([0x16, sha1(FileContent)])
	// - 大文件(>4M): hash = UrlsafeBase64([0x96, sha1([sha1(Block1), sha1(Block2), ...])])
	// 这个哈希值可以用于文件完整性校验和去重，但不是标准的MD5
	qiniuETag := ""
	if ret.Hash != "" {
		// 七牛云返回的是URL安全的Base64编码，需要先转换为标准Base64
		base64Hash := ret.Hash
		// 将URL安全的Base64转换为标准Base64
		base64Hash = strings.ReplaceAll(base64Hash, "-", "+")
		base64Hash = strings.ReplaceAll(base64Hash, "_", "/")

		// 解码Base64得到原始字节
		hashBytes, decodeErr := base64.StdEncoding.DecodeString(base64Hash)
		if decodeErr == nil {
			// 转换为十六进制字符串（这是七牛云ETag的十六进制表示，不是MD5）
			qiniuETag = hex.EncodeToString(hashBytes)
		}
	}

	return &FileSaveInfo{
		SavePath:     "/" + fileKey,
		SaveFullPath: q.Domain + "/" + fileKey,
		FileName:     name,
		FileType:     fileType,
		MD5:          qiniuETag, // 注意：这里存储的是七牛云ETag，不是标准MD5
	}, nil
}

// GetFile 下载文件到本地
func (q *QiNiu) GetFile(savePath string) (*GetFileInfo, error) {
	if savePath == "" {
		return nil, fmt.Errorf("请正确填写文件地址")
	}
	if savePath[0] != '/' {
		savePath = "/" + savePath
	}

	tempPath := filepath.Join(q.DownloadRootPath, savePath)
	url := q.GetFullPath(savePath)
	err := httpx.DownloadFile(url, tempPath)
	if err != nil {
		return nil, err
	}
	p := strings.ReplaceAll(savePath, "\\", "/")
	fileName := filepath.Base(p)  // 使用标准库获取文件名
	fileSuffix := filepath.Ext(p) // 使用标准库获取文件扩展名
	if fileSuffix != "" && fileSuffix[0] == '.' {
		fileSuffix = fileSuffix[1:] // 移除开头的点
	}

	return &GetFileInfo{
		FileSaveInfo: FileSaveInfo{
			SavePath:     savePath,
			SaveFullPath: q.Domain + savePath,
			FileName:     fileName,
		},
		FileLocalPath: tempPath,
		FileType:      fileSuffix,
	}, nil
}

// GetUploadToken 获取上传token
func (q *QiNiu) GetUploadToken() string {
	putPolicy := storage.PutPolicy{
		Scope: q.Bucket,
		//InsertOnly:1,
	}
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}

// DeleteFile 删除文件
func (q *QiNiu) DeleteFile(savePath string) error {
	savePath = strings.TrimPrefix(savePath, "/")
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	cfg := qiniuConfig()
	bucketManager := storage.NewBucketManager(mac, cfg)
	if err := bucketManager.Delete(q.Bucket, savePath); err != nil {
		return errors.New("function bucketManager.Delete() failed, err:" + err.Error())
	}
	return nil
}
