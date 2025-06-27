package files

import "context"

// 便捷构造函数

// NewReader 创建默认的Reader实现（URLReader）
func NewReader(ctx context.Context) Reader {
	return NewURLReader(ctx)
}

// NewWriter 创建默认的Writer实现（CloudWriter）
func NewWriter(ctx context.Context) Writer {
	return NewCloudWriter(ctx)
}

// 具体类型构造函数映射
var (
	// ReaderConstructors 可用的Reader构造函数
	ReaderConstructors = map[string]ReaderConstructor{
		"url":   NewURLReader,
		"http":  NewURLReader,
		"https": NewURLReader,
	}

	// WriterConstructors 可用的Writer构造函数
	WriterConstructors = map[string]WriterConstructor{
		"cloud": NewCloudWriter,
		"qiniu": NewCloudWriter,
		"oss":   NewCloudWriter,
	}
)

// GetReader 根据类型获取Reader
func GetReader(readerType string, ctx context.Context) Reader {
	if constructor, exists := ReaderConstructors[readerType]; exists {
		return constructor(ctx)
	}
	// 默认返回URLReader
	return NewURLReader(ctx)
}

// GetWriter 根据类型获取Writer
func GetWriter(writerType string, ctx context.Context) Writer {
	if constructor, exists := WriterConstructors[writerType]; exists {
		return constructor(ctx)
	}
	// 默认返回CloudWriter
	return NewCloudWriter(ctx)
}
