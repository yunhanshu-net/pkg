package main

import (
	"fmt"
	"strings"
	"github.com/yunhanshu-net/pkg/x/tagx"
)

func main() {
	fmt.Println("=== function-go 类型推断器（修正版）示例 ===\n")
	
	// 运行类型推断示例
	tagx.ExampleTypeInference()
	
	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")
	
	// 运行类型映射表示例
	tagx.ExampleTypeMapping()
	
	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")
	
	// 运行类型常量示例
	tagx.ExampleConstants()
	
	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")
	
	// 展示不支持的类型
	fmt.Println("=== 不支持的类型说明 ===")
	fmt.Println("以下类型不支持，因为无法在程序启动时确定具体类型：")
	fmt.Println("- map[string]interface{}")
	fmt.Println("- interface{}")
	fmt.Println("- 抽象数组类型（如 []interface{}）")
	fmt.Println("- time.Time（统一使用 int64 时间戳存储）")
	fmt.Println("- 其他动态类型")
	fmt.Println("这些类型会被自动忽略，不会影响其他字段的渲染。")
	fmt.Println()
	fmt.Println("只支持定义好的具体类型，确保类型在程序启动时就能确定，避免运行时变化导致渲染问题。")
	fmt.Println("时间类型统一使用 int64 时间戳存储，前端通过 datetime 组件展示。")
	fmt.Println()
	fmt.Println("=== 自动类型推断说明 ===")
	fmt.Println("data.type 字段现在会自动推断，用户无需手动指定。")
	fmt.Println("用户只需要关注 Go 类型和组件类型，系统会自动处理类型映射。")
	fmt.Println("示例：string 类型 + input 组件 = 自动推断为 string 类型")
	fmt.Println("示例：int 类型 + number 组件 = 自动推断为 number 类型")
} 