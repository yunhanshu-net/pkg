package workflow

import (
	"context"
	"fmt"
	"time"
)

// 步骤执行上下文
type StepContext struct {
	StepName  string                 `json:"step_name"`  // 步骤名称
	Function  string                 `json:"function"`   // 函数名
	InputArgs []*ArgumentInfo        `json:"input_args"` // 输入参数
	Returns   []*ArgumentInfo        `json:"returns"`    // 输出参数定义
	Metadata  map[string]interface{} `json:"metadata"`   // 元数据
	Variables map[string]interface{} `json:"variables"`  // 当前变量状态
	InputVars map[string]interface{} `json:"input_vars"` // 输入变量
	Context   context.Context        `json:"-"`          // 上下文
	StartTime time.Time              `json:"start_time"` // 开始时间
}

// 步骤执行结果
type StepCallbackResult struct {
	Success    bool          `json:"success"`     // 是否成功
	Error      error         `json:"error"`       // 错误信息
	OutputArgs []interface{} `json:"output_args"` // 输出参数值
	Duration   int64         `json:"duration_ms"` // 执行时长(毫秒)
}

// 步骤执行回调函数类型
type StepCallback func(ctx *StepContext) *StepCallbackResult

// 回调执行器
type CallbackExecutor struct {
	callbacks map[string]StepCallback
}

// 创建回调执行器
func NewCallbackExecutor() *CallbackExecutor {
	return &CallbackExecutor{
		callbacks: make(map[string]StepCallback),
	}
}

// 注册步骤回调
func (e *CallbackExecutor) RegisterStep(stepName string, callback StepCallback) {
	e.callbacks[stepName] = callback
}

// 执行步骤
func (e *CallbackExecutor) ExecuteStep(ctx *StepContext) *StepCallbackResult {
	callback, exists := e.callbacks[ctx.StepName]
	if !exists {
		return &StepCallbackResult{
			Success: false,
			Error:   fmt.Errorf("步骤 %s 的回调函数未注册", ctx.StepName),
		}
	}

	// 执行回调
	result := callback(ctx)
	if result == nil {
		result = &StepCallbackResult{
			Success: false,
			Error:   fmt.Errorf("步骤 %s 的回调函数返回nil", ctx.StepName),
		}
	}

	return result
}

// 默认回调 - 用于测试
func DefaultStepCallback(ctx *StepContext) *StepCallbackResult {
	// 模拟执行时间
	time.Sleep(100 * time.Millisecond)

	// 模拟执行结果
	outputArgs := make([]interface{}, len(ctx.Returns))
	for i, ret := range ctx.Returns {
		switch ret.Type {
		case "string":
			outputArgs[i] = fmt.Sprintf("模拟%s结果", ret.Value)
		case "int":
			outputArgs[i] = 12345
		case "bool":
			outputArgs[i] = true
		case "err":
			outputArgs[i] = nil // 模拟成功，无错误
		default:
			outputArgs[i] = "模拟结果"
		}
	}

	return &StepCallbackResult{
		Success:    true,
		Error:      nil,
		OutputArgs: outputArgs,
		Duration:   time.Since(ctx.StartTime).Milliseconds(),
	}
}
