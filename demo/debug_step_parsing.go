package main

import (
	"fmt"
	"strings"

	"github.com/yunhanshu-net/pkg/workflow"
)

func main() {
	fmt.Println("=== 调试步骤解析 ===")

	// 测试步骤定义行
	stepLine := `step1 = beiluo.test1.devops.devops_script_create(
    username: string "用户名",
    phone: int "手机号", 
    email: string "邮箱"
) -> (
    workId: string "工号",
    username: string "用户名", 
    err: error "是否失败"
);`

	fmt.Printf("步骤定义:\n%s\n", stepLine)
	fmt.Printf("包含 '=' : %v\n", strings.Contains(stepLine, "="))
	fmt.Printf("包含 '->' : %v\n", strings.Contains(stepLine, "->"))

	// 按行分割测试
	lines := strings.Split(stepLine, "\n")
	fmt.Printf("\n按行分割结果 (%d 行):\n", len(lines))
	for i, line := range lines {
		fmt.Printf("  行 %d: %s\n", i+1, line)
		fmt.Printf("    包含 '=' : %v\n", strings.Contains(line, "="))
		fmt.Printf("    包含 '->' : %v\n", strings.Contains(line, "->"))
	}

	// 测试解析
	parser := workflow.NewSimpleParser()
	step, err := parser.ParseStep(stepLine)
	if err != nil {
		fmt.Printf("解析失败: %v\n", err)
	} else {
		fmt.Printf("解析成功: %+v\n", step)
	}
}
