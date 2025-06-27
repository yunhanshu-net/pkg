package files

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yunhanshu-net/pkg/typex"
)

// Files 文件集合对象，主要用于请求参数
// 用户上传文件时可以携带额外的元数据、参数等信息
type Files struct {
	// 文件列表
	Files []*File `json:"files"`

	// 用户选择的参数/选项
	Options map[string]interface{} `json:"options"`

	// 用户备注/说明
	Note string `json:"note"`

	// 处理配置
	Config map[string]interface{} `json:"config"`

	// 创建时间
	CreatedAt typex.Time `json:"created_at"`
}

// FileInfo 简化的文件信息
type FileInfo struct {
	Name          string `json:"name"`
	Size          int64  `json:"size"`
	SizeFormatted string `json:"size_formatted"`
}

// NewFiles 创建新的文件集合
func NewFiles(files []*File) *Files {
	collection := &Files{
		Files:     files,
		CreatedAt: typex.Time(time.Now()),
		Options:   make(map[string]interface{}),
		Config:    make(map[string]interface{}),
	}

	return collection
}

// NewFilesFromWriter 从Writer创建文件集合
func NewFilesFromWriter(writer Writer) *Files {
	return NewFiles(writer.GetFiles())
}

// SetSummaryChain 设置摘要信息（链式调用版本）
func (fc *Files) SetSummaryChain(summary string) *Files {
	fc.Note = summary
	return fc
}

// SetMetadata 设置元数据
func (fc *Files) SetMetadata(key string, value interface{}) *Files {
	if fc.Options == nil {
		fc.Options = make(map[string]interface{})
	}
	fc.Options[key] = value
	return fc
}

// AddMetadata 批量添加元数据
func (fc *Files) AddMetadata(metadata map[string]interface{}) *Files {
	if fc.Options == nil {
		fc.Options = make(map[string]interface{})
	}
	for k, v := range metadata {
		fc.Options[k] = v
	}
	return fc
}

// SetOption 设置选项
func (fc *Files) SetOption(key string, value interface{}) *Files {
	if fc.Options == nil {
		fc.Options = make(map[string]interface{})
	}
	fc.Options[key] = value
	return fc
}

// SetConfig 设置配置
func (fc *Files) SetConfig(key string, value interface{}) *Files {
	if fc.Config == nil {
		fc.Config = make(map[string]interface{})
	}
	fc.Config[key] = value
	return fc
}

// SetNote 设置备注
func (fc *Files) SetNote(note string) *Files {
	fc.Note = note
	return fc
}

// formatFileSize 格式化文件大小
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// MarshalJSON 实现json.Marshaler接口
func (fc *Files) MarshalJSON() ([]byte, error) {
	// 创建一个临时结构体用于序列化
	type filesAlias Files
	return json.Marshal((*filesAlias)(fc))
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (fc *Files) UnmarshalJSON(data []byte) error {
	// 创建一个临时结构体用于反序列化
	type filesAlias Files
	alias := (*filesAlias)(fc)

	// 先初始化 map，避免 nil pointer
	if alias.Options == nil {
		alias.Options = make(map[string]interface{})
	}
	if alias.Config == nil {
		alias.Config = make(map[string]interface{})
	}

	return json.Unmarshal(data, alias)
}

// GetFiles 获取文件列表
func (fc *Files) GetFiles() []*File {
	return fc.Files
}

// GetSummary 获取摘要信息
func (fc *Files) GetSummary() string {
	return fc.Note
}

// GetMetadata 获取元数据
func (fc *Files) GetMetadata() map[string]interface{} {
	return fc.Options
}

// getFileExtension 获取文件扩展名
func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}

// ToReader 转换为Reader接口（兼容旧的Files切片类型的方法）
func (fc *Files) ToReader(ctx context.Context) Reader {
	reader := NewURLReader(ctx)
	for _, file := range fc.Files {
		reader.AddFileFromURL(file.URL, file.Name)
	}
	return reader
}

// ToWriter 转换为Writer接口（兼容旧的Files切片类型的方法）
func (fc *Files) ToWriter(ctx context.Context) Writer {
	writer := NewCloudWriter(ctx)
	for _, file := range fc.Files {
		writer.AddFile(file.LocalPath)
	}
	return writer
}

// GetTotalSize 获取所有文件的总大小（兼容旧的Files切片类型的方法）
func (fc *Files) GetTotalSize() int64 {
	var total int64
	for _, file := range fc.Files {
		total += file.Size
	}
	return total
}

// GORM 使用示例:
//
// type Task struct {
//     ID        uint           `gorm:"primarykey"`
//     Name      string         `gorm:"column:name;comment:任务名称"`
//     Files     *Files         `gorm:"type:json;column:files;comment:任务文件"`
//     CreatedAt time.Time
//     UpdatedAt time.Time
// }
//
// 使用方法:
// func CreateTask(db *gorm.DB, name string, files *Files) error {
//     task := &Task{
//         Name:  name,
//         Files: files,
//     }
//     return db.Create(task).Error
// }
//
// func GetTask(db *gorm.DB, id uint) (*Task, error) {
//     var task Task
//     err := db.First(&task, id).Error
//     return &task, err
// }
