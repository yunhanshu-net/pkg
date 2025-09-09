package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/yunhanshu-net/pkg/workflow"
)

func main() {
	// 创建类型回调执行器
	callbackExecutor := workflow.NewTypeCallbackExecutor()

	// 注册 function-call 类型回调 - 处理HTTP接口调用
	callbackExecutor.RegisterType("function-call", func(ctx *workflow.StatementContext) *workflow.StatementCallbackResult {
		fmt.Printf("🌐 [HTTP调用] %s - 参数: %v\n", ctx.Function, ctx.Args)

		// 1. 查询数据库获取步骤配置
		stepConfig, err := queryStepConfig(ctx.Function)
		if err != nil {
			return &workflow.StatementCallbackResult{
				Success: false,
				Error:   fmt.Errorf("查询步骤配置失败: %v", err),
			}
		}

		// 2. 准备HTTP请求参数
		requestData := map[string]interface{}{
			"function": ctx.Function,
			"args":     ctx.Args,
			"metadata": ctx.Metadata,
		}

		// 3. 调用HTTP接口
		response, err := callHTTPAPI(stepConfig.URL, requestData)
		if err != nil {
			return &workflow.StatementCallbackResult{
				Success: false,
				Error:   fmt.Errorf("HTTP调用失败: %v", err),
			}
		}

		// 4. 解析响应并返回结果
		outputArgs := make([]interface{}, len(ctx.Returns))
		for i := range ctx.Returns {
			if i < len(response.OutputArgs) {
				outputArgs[i] = response.OutputArgs[i]
			} else {
				outputArgs[i] = nil
			}
		}

		return &workflow.StatementCallbackResult{
			Success:    response.Success,
			Error:      response.Error,
			OutputArgs: outputArgs,
			Duration:   time.Since(ctx.StartTime).Milliseconds(),
		}
	})

	// 注册 if 类型回调 - 处理条件判断
	callbackExecutor.RegisterType("if", func(ctx *workflow.StatementContext) *workflow.StatementCallbackResult {
		fmt.Printf("🔀 [条件判断] %s\n", ctx.Condition)

		// 根据条件判断是否跳过执行
		shouldSkip := evaluateCondition(ctx.Condition, ctx.Variables, ctx.InputVars)

		return &workflow.StatementCallbackResult{
			Success:    true,
			Error:      nil,
			ShouldSkip: shouldSkip,
			Duration:   time.Since(ctx.StartTime).Milliseconds(),
		}
	})

	// 注册 var 类型回调 - 处理变量赋值
	callbackExecutor.RegisterType("var", func(ctx *workflow.StatementContext) *workflow.StatementCallbackResult {
		fmt.Printf("📝 [变量赋值] %s\n", ctx.Content)

		// 处理变量赋值逻辑
		// 这里可以添加变量验证、类型转换等逻辑

		return &workflow.StatementCallbackResult{
			Success:  true,
			Error:    nil,
			Duration: time.Since(ctx.StartTime).Milliseconds(),
		}
	})

	// 创建执行引擎
	executor := workflow.NewWorkflowExecutorWithCallback(callbackExecutor)

	// 工作流代码
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

	// 执行工作流
	fmt.Println("============================================================")
	fmt.Println("AI工作流编排语言 - 生产环境演示")
	fmt.Println("============================================================")

	result := executor.ExecuteWorkflow(code)

	// 打印执行结果
	result.Print()
}

// 步骤配置
type StepConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Timeout int               `json:"timeout"`
}

// 查询步骤配置
func queryStepConfig(functionName string) (*StepConfig, error) {
	// 模拟数据库查询
	configs := map[string]*StepConfig{
		"step1": {
			URL:     "http://api.example.com/devops/create-user",
			Method:  "POST",
			Headers: map[string]string{"Content-Type": "application/json"},
			Timeout: 5000,
		},
		"step2": {
			URL:     "http://api.example.com/crm/schedule-interview",
			Method:  "POST",
			Headers: map[string]string{"Content-Type": "application/json"},
			Timeout: 3000,
		},
		"step3": {
			URL:     "http://api.example.com/notification/send-email",
			Method:  "POST",
			Headers: map[string]string{"Content-Type": "application/json"},
			Timeout: 2000,
		},
	}

	config, exists := configs[functionName]
	if !exists {
		return nil, fmt.Errorf("步骤 %s 的配置未找到", functionName)
	}

	return config, nil
}

// HTTP API响应
type APIResponse struct {
	Success    bool          `json:"success"`
	Error      error         `json:"error"`
	OutputArgs []interface{} `json:"output_args"`
}

// 调用HTTP接口
func callHTTPAPI(url string, data map[string]interface{}) (*APIResponse, error) {
	// 模拟HTTP调用
	jsonData, _ := json.Marshal(data)

	// 这里应该是真实的HTTP调用
	// resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	fmt.Printf("   📡 调用接口: %s\n", url)
	fmt.Printf("   📤 请求数据: %s\n", string(jsonData))

	// 模拟响应
	time.Sleep(100 * time.Millisecond)

	return &APIResponse{
		Success: true,
		Error:   nil,
		OutputArgs: []interface{}{
			"真实工号结果",
			"真实用户名结果",
			nil,
		},
	}, nil
}

// 评估条件
func evaluateCondition(condition string, variables map[string]interface{}, inputVars map[string]interface{}) bool {
	// 简单的条件评估逻辑
	if condition == "" {
		return false
	}

	// 检查错误条件
	if strings.Contains(condition, "Err != nil") {
		// 提取错误变量名
		parts := strings.Split(condition, "Err != nil")
		if len(parts) > 0 {
			errVarName := strings.TrimSpace(parts[0])
			if errVarName == "" {
				errVarName = "step1Err" // 默认
			}

			// 从变量中查找错误值
			if err, exists := variables[errVarName]; exists {
				return err != nil
			}
		}
	}

	return false
}
