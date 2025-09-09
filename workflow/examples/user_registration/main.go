package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yunhanshu-net/pkg/workflow"
)

func main() {
	fmt.Println("🚀 用户注册工作流演示")
	fmt.Println("========================")

	// 1. 定义工作流代码
	workflowCode := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
    "部门": "技术部",
    "职位": "高级工程师"
}

// 步骤1：创建用户账号
step1 = beiluo.test1.user.create_user(
    username: string "用户名",
    phone: int "手机号",
    email: string "邮箱"
) -> (
    userId: string "用户ID",
    username: string "用户名",
    err: error "是否失败"
);

// 步骤2：分配部门
step2 = beiluo.test1.user.assign_department(
    userId: string "用户ID",
    department: string "部门",
    position: string "职位"
) -> (
    success: bool "是否成功",
    message: string "消息",
    err: error "是否失败"
);

// 步骤3：发送欢迎邮件
step3 = beiluo.test1.notify.send_welcome_email(
    userId: string "用户ID",
    username: string "用户名",
    email: string "邮箱",
    department: string "部门"
) -> (
    success: bool "是否成功",
    err: error "是否失败"
);

// 步骤4：创建用户档案
step4 = beiluo.test1.user.create_profile(
    userId: string "用户ID",
    username: string "用户名",
    department: string "部门",
    position: string "职位"
) -> (
    profileId: string "档案ID",
    err: error "是否失败"
);

func main() {
    // 步骤1：创建用户账号
    用户ID, 用户名, step1Err := step1(input["用户名"], input["手机号"], input["邮箱"]){retry:3, timeout:5000}
    if step1Err != nil {
        step1.Printf("❌ 用户创建失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，用户ID: %s", 用户ID)
    
    // 步骤2：分配部门
    分配成功, 消息, step2Err := step2(用户ID, input["部门"], input["职位"]){retry:2, timeout:3000}
    if step2Err != nil {
        step2.Printf("❌ 部门分配失败: %v", step2Err)
        return
    }
    step2.Printf("✅ 部门分配成功: %s", 消息)
    
    // 步骤3：发送欢迎邮件
    邮件成功, step3Err := step3(用户ID, 用户名, input["邮箱"], input["部门"]){retry:1, timeout:2000}
    if step3Err != nil {
        step3.Printf("❌ 欢迎邮件发送失败: %v", step3Err)
        return
    }
    step3.Printf("✅ 欢迎邮件发送成功")
    
    // 步骤4：创建用户档案
    档案ID, step4Err := step4(用户ID, 用户名, input["部门"], input["职位"]){retry:2, timeout:4000}
    if step4Err != nil {
        step4.Printf("❌ 用户档案创建失败: %v", step4Err)
        return
    }
    step4.Printf("✅ 用户档案创建成功，档案ID: %s", 档案ID)
    
    fmt.Printf("🎉 用户注册流程完成！用户: %s, ID: %s, 档案: %s\n", 用户名, 用户ID, 档案ID)
}`

	// 2. 解析工作流
	parser := workflow.NewSimpleParser()
	parseResult := parser.ParseWorkflow(workflowCode)
	if !parseResult.Success {
		log.Fatalf("❌ 工作流解析失败: %s", parseResult.Error)
	}

	// 3. 设置FlowID
	parseResult.FlowID = "user-registration-" + fmt.Sprintf("%d", time.Now().Unix())

	// 4. 创建执行器
	executor := workflow.NewExecutor()

	// 5. 设置回调函数
	executor.OnFunctionCall = func(ctx context.Context, step workflow.SimpleStep, in *workflow.ExecutorIn) (*workflow.ExecutorOut, error) {
		fmt.Printf("\n📋 执行步骤: %s - %s\n", step.Name, in.StepDesc)
		fmt.Printf("📥 输入参数: %+v\n", in.RealInput)

		// 模拟不同的业务逻辑
		switch step.Name {
		case "step1":
			// 模拟用户创建
			time.Sleep(100 * time.Millisecond) // 模拟网络延迟
			return &workflow.ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"userId":   "ks_beiluo",
					"username": in.RealInput["username"],
					"err":      nil,
				},
				Error: "",
				Logs:  []string{"用户账号创建成功"},
			}, nil

		case "step2":
			// 模拟部门分配
			time.Sleep(80 * time.Millisecond)
			return &workflow.ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"success": true,
					"message": fmt.Sprintf("已分配到 %s 部门", in.RealInput["department"]),
					"err":     nil,
				},
				Error: "",
				Logs:  []string{"部门分配成功"},
			}, nil

		case "step3":
			// 模拟邮件发送
			time.Sleep(120 * time.Millisecond)
			return &workflow.ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"success": true,
					"err":     nil,
				},
				Error: "",
				Logs:  []string{"欢迎邮件发送成功"},
			}, nil

		case "step4":
			// 模拟档案创建
			time.Sleep(90 * time.Millisecond)
			return &workflow.ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"profileId": "PROFILE_" + fmt.Sprintf("%d", time.Now().Unix()),
					"err":       nil,
				},
				Error: "",
				Logs:  []string{"用户档案创建成功"},
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
		log.Fatalf("❌ 工作流执行失败: %v", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("\n⏱️  总执行时间: %v\n", duration)

	// 7. 显示最终结果
	fmt.Println("\n📊 最终变量状态:")
	for name, varInfo := range parseResult.Variables {
		fmt.Printf("  %s: %v (%s)\n", name, varInfo.Value, varInfo.Type)
	}

	// 8. 显示步骤日志
	fmt.Println("\n📝 步骤执行日志:")
	for _, step := range parseResult.Steps {
		if len(step.Logs) > 0 {
			fmt.Printf("  %s:\n", step.Name)
			for _, log := range step.Logs {
				fmt.Printf("    %s\n", log.Message)
			}
		}
	}
}
