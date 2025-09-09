package workflow

import (
	"fmt"
	"strings"
)

func main() {
	// 演示工作流代码
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
    "部门": "技术部",
}

step1 = beiluo.test1.devops.devops_script_create(string 用户名, int 手机号, string 邮箱) -> (string 工号, string 用户名, err 是否失败);
step2 = beiluo.test1.crm.crm_interview_schedule(string 用户名, string 部门) -> (string 面试时间, string 面试官名称, err 是否失败);
step3 = beiluo.test1.notification.send_email(string 邮箱, string 内容) -> (err 是否失败);

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
    
    // 发送通知邮件
    通知内容 := "你收到了:{{用户名}},时间：{{面试时间}}的面试安排，请关注"
    step3Err := step3(input["邮箱"], 通知内容){retry:1, timeout:2000, priority:"low"}
    if step3Err != nil {
        fmt.Printf("⚠️ 发送邮件失败: %v\n", step3Err)
    } else {
        fmt.Printf("✅ 邮件发送成功\n")
    }
    
    fmt.Printf("🎉 流程完成！工号: %s，面试时间: %s\n", 工号, 面试时间)
}`

	// 创建执行引擎
	executor := NewWorkflowExecutor()

	// 执行工作流
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("AI工作流编排语言 - 执行引擎演示")
	fmt.Println(strings.Repeat("=", 60))

	result := executor.ExecuteWorkflow(code)

	// 打印执行结果
	result.Print()

	// 打印详细的执行统计
	fmt.Println("\n📈 执行统计:")
	fmt.Printf("   总步骤数: %d\n", len(result.Steps))
	fmt.Printf("   成功步骤: %d\n", countSuccessfulSteps(result.Steps))
	fmt.Printf("   失败步骤: %d\n", len(result.Steps)-countSuccessfulSteps(result.Steps))

	// 打印变量状态
	fmt.Println("\n📋 变量状态:")
	for key, value := range result.Variables {
		fmt.Printf("   %s: %v\n", key, value)
	}
}

// 统计成功步骤数
func countSuccessfulSteps(steps []StepExecutionResult) int {
	count := 0
	for _, step := range steps {
		if step.Success {
			count++
		}
	}
	return count
}
