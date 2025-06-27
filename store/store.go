package store

import "os"

// FileStore 存储接口
type FileStore interface {
	// FileSave 保存文件
	FileSave(localFilePath string, savePath string) (*FileSaveInfo, error)

	// GetFile 下载文件
	GetFile(savePath string) (*GetFileInfo, error)

	// DeleteFile 删除文件
	DeleteFile(savePath string) error
}

type FileSaveInfo struct {
	SavePath     string `json:"save_path"`      //   /beiluo/d.zip
	SaveFullPath string `json:"save_full_path"` //http://cdn.geeleo.com/beiluo/d.zip
	//FileHomedir  string `json:"file_homedir"`   //文件解压后目录名字
	FileName string `json:"file_name"` //文件名称
	FileType string `json:"file_type"` //文件类型
	MD5      string `json:"md5"`       //文件哈希值，用于完整性校验和去重（注意：不同存储服务商可能使用不同的哈希算法，如七牛云使用ETag算法）
}

type GetFileInfo struct {
	FileSaveInfo
	FileLocalPath string `json:"file_local_path"` //下载后本地存储的文件地址
	FileSize      int64  `json:"file_size"`
	FileType      string `json:"file_type"`
}

func (i *GetFileInfo) RemoveLocalFile() error {
	return os.Remove(i.FileLocalPath)
}
