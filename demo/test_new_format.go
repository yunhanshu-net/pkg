package main

import (
	"fmt"

	"github.com/yunhanshu-net/pkg/workflow"
)

func main() {
	fmt.Println("=== 测试新的参数定义格式 ===")

	parser := workflow.NewSimpleParser()

	// 测试新的工作流代码格式
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
    "部门": "技术部",
}

step1 = beiluo.test1.devops.devops_script_create(
    username: string "用户名",
    phone: int "手机号", 
    email: string "邮箱"
) -> (
    workId: string "工号",
    username: string "用户名", 
    err: error "是否失败"
);

func main() {
    fmt.Println("🚀 开始用户注册和面试安排流程...")
    
    // 创建用户
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"], input["邮箱"]){retry:3, timeout:5000, priority:"high"}
    if step1Err != nil {
        fmt.Printf("❌ 创建用户失败: %v\n", step1Err)
        return
    }
    fmt.Printf("✅ 用户创建成功，工号: %s\n", 工号)
    
    fmt.Printf("🎉 流程完成！工号: %s\n", 工号)
}`

	result := parser.ParseWorkflow(code)

	if !result.Success {
		fmt.Printf("解析失败: %s\n", result.Error)
		return
	}

	fmt.Printf("解析成功！\n")
	fmt.Printf("步骤数量: %d\n", len(result.Steps))

	for i, step := range result.Steps {
		fmt.Printf("\n步骤 %d: %s\n", i+1, step.Name)
		fmt.Printf("  函数: %s\n", step.Function)
		fmt.Printf("  输入参数: %d 个\n", len(step.InputParams))
		for j, param := range step.InputParams {
			fmt.Printf("    %d. %s (%s) - %s\n", j+1, param.Name, param.Type, param.Desc)
		}
		fmt.Printf("  输出参数: %d 个\n", len(step.OutputParams))
		for j, param := range step.OutputParams {
			fmt.Printf("    %d. %s (%s) - %s\n", j+1, param.Name, param.Type, param.Desc)
		}
	}

	// 打印详细结果
	fmt.Println("\n=== 详细解析结果 ===")
	result.Print()
}
