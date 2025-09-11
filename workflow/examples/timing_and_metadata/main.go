package main

import (
	"context"
	"fmt"
	"time"

	"github.com/yunhanshu-net/pkg/workflow"
)

func main() {
	// 工作流代码 - 展示执行耗时记录和元数据功能
	code := `
var input = map[string]interface{}{
	"用户名": "张三",
	"手机号": "13800138000",
	"邮箱": "zhangsan@example.com",
}

step1 = beiluo.test1.user.create_user(username: string "用户名", phone: string "手机号") -> (userId: string "用户ID", err: string "错误信息");
step2 = beiluo.test1.user.send_email(email: string "邮箱", userId: string "用户ID") -> (success: bool "是否成功", err: string "错误信息");
step3 = beiluo.test1.user.activate_user(userId: string "用户ID") -> (success: bool "是否成功", err: string "错误信息");

func main() {
	//desc: 开始执行用户注册流程
	sys.Print("开始执行用户注册流程...")
	
	//desc: 创建用户账户 - 使用自定义元数据
	userId, err := step1(input["用户名"], input["手机号"]){timeout: 10000, retry_count: 3, debug: true, priority: 1}
	
	//desc: 检查用户创建是否成功
	if err != nil {
		//desc: 创建失败，记录错误
		sys.Print("用户创建失败: {{err}}")
		return
	}
	
	//desc: 创建成功，记录用户ID
	sys.Print("用户创建成功，用户ID: {{userId}}")
	
	//desc: 发送欢迎邮件 - 使用不同的元数据配置
	success, err := step2(input["邮箱"], userId){timeout: 5000, retry_count: 1, async: true, log_level: "debug"}
	
	//desc: 检查邮件发送是否成功
	if err != nil {
		//desc: 邮件发送失败，记录错误
		sys.Print("邮件发送失败: {{err}}")
		return
	}
	
	//desc: 邮件发送成功
	sys.Print("邮件发送成功: {{success}}")
	
	//desc: 激活用户账户 - 使用默认元数据（无超时限制）
	success, err = step3(userId)
	
	//desc: 检查用户激活是否成功
	if err != nil {
		//desc: 激活失败，记录错误
		sys.Print("用户激活失败: {{err}}")
		return
	}
	
	//desc: 激活成功，流程完成
	sys.Print("用户激活成功: {{success}}")
	sys.Print("🎉 用户注册流程执行完成！")
}
`

	// 创建解析器
	parser := workflow.NewSimpleParser()
	result := parser.ParseWorkflow(code)

	// 检查解析结果
	if !result.Success {
		sys.Printf("❌ 解析失败: %s\n", result.Error)
		return
	}

	// 设置FlowID
	result.FlowID = "user-registration-" + fmt.Sprintf("%d", time.Now().Unix())

	// 创建执行器
	executor := workflow.NewExecutor()

	// 设置回调函数
	executor.OnFunctionCall = func(ctx context.Context, step workflow.SimpleStep, in *workflow.ExecutorIn) (*workflow.ExecutorOut, error) {
		sys.Printf("【print】执行步骤: %s - %s\n", step.Name, in.StepDesc)
		sys.Printf("【print】输入参数: %+v\n", in.RealInput)
		sys.Printf("【print】元数据配置: %+v\n", in.Metadata)

		// 模拟不同的执行时间
		var sleepTime time.Duration
		switch step.Name {
		case "step1":
			sleepTime = 200 * time.Millisecond // 用户创建需要200ms
		case "step2":
			sleepTime = 150 * time.Millisecond // 邮件发送需要150ms
		case "step3":
			sleepTime = 100 * time.Millisecond // 用户激活需要100ms
		default:
			sleepTime = 50 * time.Millisecond
		}

		time.Sleep(sleepTime)

		// 模拟业务逻辑
		return &workflow.ExecutorOut{
			Success: true,
			WantOutput: map[string]interface{}{
				"userId":  "USER_" + fmt.Sprintf("%d", time.Now().Unix()),
				"success": true,
				"err":     nil,
			},
			Error: "",
			Logs:  []string{fmt.Sprintf("步骤 %s 执行成功", step.Name)},
		}, nil
	}

	executor.OnWorkFlowUpdate = func(ctx context.Context, current *workflow.SimpleParseResult) error {
		sys.Printf("【print】工作流状态更新: FlowID=%s\n", current.FlowID)
		return nil
	}

	executor.OnWorkFlowExit = func(ctx context.Context, current *workflow.SimpleParseResult) error {
		sys.Println("【print】工作流正常结束")
		return nil
	}

	// 执行工作流
	ctx := context.Background()
	startTime := time.Now()

	err := executor.Start(ctx, result)
	if err != nil {
		sys.Printf("❌ 执行失败: %v\n", err)
		return
	}

	totalDuration := time.Since(startTime)

	// 输出执行统计信息
	sys.Println("\n=== 执行统计信息 ===")
	sys.Printf("总执行时间: %v\n", totalDuration)
	sys.Printf("工作流ID: %s\n", result.FlowID)

	// 输出每个语句的执行时间
	sys.Println("\n=== 语句执行时间 ===")
	for i, stmt := range result.MainFunc.Statements {
		if stmt.StartTime != nil && stmt.EndTime != nil {
			sys.Printf("语句 %d: %s\n", i+1, stmt.Content)
			sys.Printf("  开始时间: %v\n", stmt.StartTime.Format("15:04:05.000"))
			sys.Printf("  结束时间: %v\n", stmt.EndTime.Format("15:04:05.000"))
			sys.Printf("  执行耗时: %v\n", stmt.Duration)
			sys.Printf("  状态: %s\n", stmt.Status)

			// 如果是function-call，显示元数据信息
			if stmt.Type == "function-call" {
				sys.Printf("  元数据: %+v\n", stmt.GetMergedMetadata())
			}
			sys.Println()
		}
	}

	// 输出变量信息
	sys.Println("=== 变量信息 ===")
	for name, varInfo := range result.Variables {
		sys.Printf("变量 %s: 类型=%s, 值=%v, 来源=%s\n",
			name, varInfo.Type, varInfo.Value, varInfo.Source)
	}

	sys.Println("\n✅ 用户注册流程执行完成！")
}
