package typex

// 文件处理相关的类型已经移动到 pkg/x/typex/files 包中
// 为了保持向后兼容，请使用以下导入：
//
//   import "github.com/yunhanshu-net/function-go/pkg/x/typex/files"
//
// 然后使用：
//   files.NewFilesReader(ctx)
//   files.NewFilesWriter(ctx)
//
// 这样的设计可以避免typex包变得过于臃肿，同时保持良好的模块化结构。
