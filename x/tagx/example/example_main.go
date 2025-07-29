package main

import (
	"fmt"
	"strings"
	"github.com/yunhanshu-net/pkg/x/tagx"
)

func main() {
	fmt.Println("=== function-go 类型推断器示例 ===\n")
	
	// 运行类型推断示例
	tagx.ExampleTypeInference()
	
	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")
	
	// 运行类型映射表示例
	tagx.ExampleTypeMapping()
	
	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")
	
	// 运行类型常量示例
	tagx.ExampleConstants()
} 