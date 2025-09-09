package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/yunhanshu-net/pkg/workflow"
)

func main() {
	fmt.Println("🚀 并发工作流执行演示")
	fmt.Println("========================")

	// 1. 定义工作流代码
	workflowCode := `var input = map[string]interface{}{
    "用户名": "王五",
    "手机号": 13700137000,
    "邮箱": "wangwu@example.com"
}

// 步骤1：创建用户
step1 = beiluo.test1.user.create_user(
    username: string "用户名",
    phone: int "手机号",
    email: string "邮箱"
) -> (
    userId: string "用户ID",
    err: error "是否失败"
);

// 步骤2：发送欢迎邮件
step2 = beiluo.test1.notify.send_welcome_email(
    userId: string "用户ID",
    email: string "邮箱"
) -> (
    success: bool "是否成功",
    err: error "是否失败"
);

// 步骤3：创建用户档案
step3 = beiluo.test1.user.create_profile(
    userId: string "用户ID",
    username: string "用户名"
) -> (
    profileId: string "档案ID",
    err: error "是否失败"
);

func main() {
    用户ID, step1Err := step1(input["用户名"], input["手机号"], input["邮箱"]){retry:3, timeout:5000}
    if step1Err != nil {
        step1.Printf("❌ 用户创建失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，用户ID: %s", 用户ID)
    
    邮件成功, step2Err := step2(用户ID, input["邮箱"]){retry:2, timeout:3000}
    if step2Err != nil {
        step2.Printf("❌ 邮件发送失败: %v", step2Err)
        return
    }
    step2.Printf("✅ 邮件发送成功")
    
    档案ID, step3Err := step3(用户ID, input["用户名"]){retry:2, timeout:4000}
    if step3Err != nil {
        step3.Printf("❌ 档案创建失败: %v", step3Err)
        return
    }
    step3.Printf("✅ 档案创建成功，档案ID: %s", 档案ID)
    
    fmt.Printf("🎉 用户注册完成！用户ID: %s, 档案ID: %s\n", 用户ID, 档案ID)
}`

	// 2. 解析工作流
	parser := workflow.NewSimpleParser()
	parseResult := parser.ParseWorkflow(workflowCode)
	if !parseResult.Success {
		log.Fatalf("❌ 工作流解析失败: %s", parseResult.Error)
	}

	// 3. 创建执行器
	executor := workflow.NewExecutor()

	// 4. 设置回调函数
	executor.OnFunctionCall = func(ctx context.Context, step workflow.SimpleStep, in *workflow.ExecutorIn) (*workflow.ExecutorOut, error) {
		fmt.Printf("[%s] 📋 执行步骤: %s - %s\n", time.Now().Format("15:04:05"), step.Name, in.StepDesc)
		fmt.Printf("[%s] 📥 输入参数: %+v\n", time.Now().Format("15:04:05"), in.RealInput)

		// 模拟不同的执行时间
		var sleepTime time.Duration
		switch step.Name {
		case "step1":
			sleepTime = 200 * time.Millisecond
		case "step2":
			sleepTime = 150 * time.Millisecond
		case "step3":
			sleepTime = 180 * time.Millisecond
		default:
			sleepTime = 100 * time.Millisecond
		}

		time.Sleep(sleepTime)

		switch step.Name {
		case "step1":
			return &workflow.ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"userId": "USER_" + fmt.Sprintf("%d", time.Now().UnixNano()),
					"err":    nil,
				},
				Error: "",
				Logs:  []string{"用户创建成功"},
			}, nil

		case "step2":
			return &workflow.ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"success": true,
					"err":     nil,
				},
				Error: "",
				Logs:  []string{"邮件发送成功"},
			}, nil

		case "step3":
			return &workflow.ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"profileId": "PROFILE_" + fmt.Sprintf("%d", time.Now().UnixNano()),
					"err":       nil,
				},
				Error: "",
				Logs:  []string{"档案创建成功"},
			}, nil

		default:
			return &workflow.ExecutorOut{Success: false, Error: "未知步骤"}, nil
		}
	}

	executor.OnWorkFlowUpdate = func(ctx context.Context, current *workflow.SimpleParseResult) error {
		fmt.Printf("[%s] 🔄 工作流状态更新: FlowID=%s\n", time.Now().Format("15:04:05"), current.FlowID)
		return nil
	}

	executor.OnWorkFlowExit = func(ctx context.Context, current *workflow.SimpleParseResult) error {
		fmt.Printf("[%s] ✅ 工作流正常结束\n", time.Now().Format("15:04:05"))
		return nil
	}

	executor.OnWorkFlowReturn = func(ctx context.Context, current *workflow.SimpleParseResult) error {
		fmt.Printf("[%s] ❌ 工作流因错误中断\n", time.Now().Format("15:04:05"))
		return nil
	}

	// 5. 并发执行多个工作流实例
	const numWorkflows = 5
	var wg sync.WaitGroup
	results := make(chan string, numWorkflows)

	fmt.Printf("\n🚀 开始并发执行 %d 个工作流实例...\n", numWorkflows)

	startTime := time.Now()

	for i := 0; i < numWorkflows; i++ {
		wg.Add(1)
		go func(instanceID int) {
			defer wg.Done()

			// 为每个实例创建独立的工作流
			instanceResult := *parseResult
			instanceResult.FlowID = fmt.Sprintf("concurrent-demo-%d-%d", instanceID, time.Now().Unix())

			ctx := context.Background()
			if err := executor.Start(ctx, &instanceResult); err != nil {
				results <- fmt.Sprintf("实例 %d 执行失败: %v", instanceID, err)
				return
			}

			results <- fmt.Sprintf("实例 %d 执行成功: FlowID=%s", instanceID, instanceResult.FlowID)
		}(i)
	}

	// 等待所有工作流完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	successCount := 0
	for result := range results {
		fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), result)
		if result != "" {
			successCount++
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\n⏱️  总执行时间: %v\n", duration)
	fmt.Printf("📊 成功执行: %d/%d 个工作流实例\n", successCount, numWorkflows)
	fmt.Printf("🚀 平均每个工作流执行时间: %v\n", duration/time.Duration(numWorkflows))
}
