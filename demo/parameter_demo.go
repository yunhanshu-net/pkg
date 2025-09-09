package main

import (
	"fmt"

	"github.com/yunhanshu-net/pkg/workflow"
)

func main() {
	// 测试新的参数定义格式
	parser := workflow.NewSimpleParser()

	// 测试参数定义解析
	testParameterDefinitions(parser)

	// 测试完整的工作流解析
	testWorkflowParsing(parser)
}

func testParameterDefinitions(parser *workflow.SimpleParser) {
	fmt.Println("=== 测试参数定义解析 ===")

	// 测试单个参数定义
	testCases := []string{
		`username: string "用户名"`,
		`phone: int "手机号"`,
		`email: string "邮箱"`,
		`workId: string "工号"`,
		`err: error "是否失败"`,
	}

	for _, testCase := range testCases {
		fmt.Printf("测试: %s\n", testCase)

		// 这里需要调用解析函数，但我们需要先实现它
		// params := parser.parseParameterDefinitions(testCase)
		// fmt.Printf("结果: %+v\n", params)
		fmt.Println("  (解析函数待实现)")
	}
}

func testWorkflowParsing(parser *workflow.SimpleParser) {
	fmt.Println("\n=== 测试工作流解析 ===")

	// 新的工作流代码格式
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

step2 = beiluo.test1.crm.crm_interview_schedule(
    username: string "用户名",
    department: string "部门"
) -> (
    interviewTime: string "面试时间",
    interviewer: string "面试官名称", 
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
    
    // 安排面试
    面试时间, 面试官名称, step2Err := step2(用户名, input["部门"]){retry:2, timeout:3000, priority:"normal"}
    if step2Err != nil {
        fmt.Printf("❌ 安排面试失败: %v\n", step2Err)
        return
    }
    fmt.Printf("✅ 面试安排成功，时间: %s，面试官: %s\n", 面试时间, 面试官名称)
    
    fmt.Printf("🎉 流程完成！工号: %s，面试时间: %s\n", 工号, 面试时间)
}`

	result := parser.ParseWorkflow(code)

	if !result.Success {
		fmt.Printf("解析失败: %s\n", result.Error)
		return
	}

	fmt.Printf("解析成功！\n")
	fmt.Printf("步骤数量: %d\n", len(result.Steps))

	for i, step := range result.Steps {
		fmt.Printf("步骤 %d: %s\n", i+1, step.Name)
		fmt.Printf("  函数: %s\n", step.Function)
		fmt.Printf("  输入参数: %d 个\n", len(step.InputParams))
		fmt.Printf("  输出参数: %d 个\n", len(step.OutputParams))
	}
}
