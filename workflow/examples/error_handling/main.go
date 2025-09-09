package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yunhanshu-net/pkg/workflow"
)

func main() {
	fmt.Println("🚨 错误处理工作流演示")
	fmt.Println("========================")

	// 1. 定义工作流代码 - 包含可能失败的步骤
	workflowCode := `var input = map[string]interface{}{
    "用户名": "李四",
    "手机号": 13900139000,
    "邮箱": "lisi@example.com"
}

// 步骤1：验证用户信息
step1 = beiluo.test1.user.validate_user(
    username: string "用户名",
    phone: int "手机号",
    email: string "邮箱"
) -> (
    valid: bool "是否有效",
    message: string "验证消息",
    err: error "是否失败"
);

// 步骤2：创建用户账号
step2 = beiluo.test1.user.create_user(
    username: string "用户名",
    phone: int "手机号",
    email: string "邮箱"
) -> (
    userId: string "用户ID",
    err: error "是否失败"
);

// 步骤3：发送验证邮件
step3 = beiluo.test1.notify.send_verification_email(
    userId: string "用户ID",
    email: string "邮箱"
) -> (
    success: bool "是否成功",
    err: error "是否失败"
);

func main() {
    // 步骤1：验证用户信息
    验证通过, 验证消息, step1Err := step1(input["用户名"], input["手机号"], input["邮箱"]){retry:2, timeout:3000}
    if step1Err != nil {
        step1.Printf("❌ 用户信息验证失败: %v", step1Err)
        return
    }
    if !验证通过 {
        step1.Printf("❌ 用户信息无效: %s", 验证消息)
        return
    }
    step1.Printf("✅ 用户信息验证通过: %s", 验证消息)
    
    // 步骤2：创建用户账号
    用户ID, step2Err := step2(input["用户名"], input["手机号"], input["邮箱"]){retry:3, timeout:5000}
    if step2Err != nil {
        step2.Printf("❌ 用户账号创建失败: %v", step2Err)
        return
    }
    step2.Printf("✅ 用户账号创建成功，用户ID: %s", 用户ID)
    
    // 步骤3：发送验证邮件
    邮件成功, step3Err := step3(用户ID, input["邮箱"]){retry:1, timeout:2000}
    if step3Err != nil {
        step3.Printf("❌ 验证邮件发送失败: %v", step3Err)
        return
    }
    if !邮件成功 {
        step3.Printf("❌ 验证邮件发送失败")
        return
    }
    step3.Printf("✅ 验证邮件发送成功")
    
    fmt.Printf("🎉 用户注册流程完成！用户ID: %s\n", 用户ID)
}`

	// 2. 解析工作流
	parser := workflow.NewSimpleParser()
	parseResult := parser.ParseWorkflow(workflowCode)
	if !parseResult.Success {
		log.Fatalf("❌ 工作流解析失败: %s", parseResult.Error)
	}

	// 3. 设置FlowID
	parseResult.FlowID = "error-handling-demo-" + fmt.Sprintf("%d", time.Now().Unix())

	// 4. 创建执行器
	executor := workflow.NewExecutor()

	// 5. 设置回调函数 - 模拟不同的失败场景
	executor.OnFunctionCall = func(ctx context.Context, step workflow.SimpleStep, in *workflow.ExecutorIn) (*workflow.ExecutorOut, error) {
		fmt.Printf("\n📋 执行步骤: %s - %s\n", step.Name, in.StepDesc)
		fmt.Printf("📥 输入参数: %+v\n", in.RealInput)

		switch step.Name {
		case "step1":
			// 模拟验证逻辑 - 可能失败
			time.Sleep(100 * time.Millisecond)
			username := in.RealInput["username"].(string)

			// 模拟验证失败场景
			if username == "李四" {
				return &workflow.ExecutorOut{
					Success: true,
					WantOutput: map[string]interface{}{
						"valid":   false,
						"message": "用户名已存在，请选择其他用户名",
						"err":     nil,
					},
					Error: "",
					Logs:  []string{"用户名验证失败"},
				}, nil
			}

			return &workflow.ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"valid":   true,
					"message": "用户信息验证通过",
					"err":     nil,
				},
				Error: "",
				Logs:  []string{"用户信息验证成功"},
			}, nil

		case "step2":
			// 模拟用户创建 - 总是成功
			time.Sleep(80 * time.Millisecond)
			return &workflow.ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"userId": "USER_" + fmt.Sprintf("%d", time.Now().Unix()),
					"err":    nil,
				},
				Error: "",
				Logs:  []string{"用户账号创建成功"},
			}, nil

		case "step3":
			// 模拟邮件发送 - 可能失败
			time.Sleep(120 * time.Millisecond)

			// 模拟邮件发送失败
			return &workflow.ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"success": false,
					"err":     nil,
				},
				Error: "",
				Logs:  []string{"邮件服务暂时不可用"},
			}, nil

		default:
			return &workflow.ExecutorOut{Success: false, Error: "未知步骤"}, nil
		}
	}

	executor.OnWorkFlowUpdate = func(ctx context.Context, current *workflow.SimpleParseResult) error {
		fmt.Printf("🔄 工作流状态更新: FlowID=%s, 变量数量=%d\n", current.FlowID, len(current.Variables))
		return nil
	}

	executor.OnWorkFlowExit = func(ctx context.Context, current *workflow.SimpleParseResult) error {
		fmt.Println("\n✅ 工作流正常结束")
		return nil
	}

	executor.OnWorkFlowReturn = func(ctx context.Context, current *workflow.SimpleParseResult) error {
		fmt.Println("\n❌ 工作流因错误中断")
		return nil
	}

	// 6. 执行工作流
	ctx := context.Background()
	startTime := time.Now()

	fmt.Println("\n🚀 开始执行工作流...")
	if err := executor.Start(ctx, parseResult); err != nil {
		fmt.Printf("❌ 工作流执行失败: %v\n", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("\n⏱️  总执行时间: %v\n", duration)

	// 7. 显示最终结果
	fmt.Println("\n📊 最终变量状态:")
	for name, varInfo := range parseResult.Variables {
		fmt.Printf("  %s: %v (%s)\n", name, varInfo.Value, varInfo.Type)
	}
}
