package workflow

import (
	"context"
	"fmt"
	"time"
)

// 语句执行上下文
type StatementContext struct {
	Type       string                 `json:"type"`        // 语句类型：function-call, if, var, print, return
	Content    string                 `json:"content"`     // 语句内容
	Function   string                 `json:"function"`    // 函数名（仅function-call）
	Args       []*ArgumentInfo        `json:"args"`        // 输入参数（仅function-call）
	Returns    []*ArgumentInfo        `json:"returns"`     // 输出参数（仅function-call）
	Condition  string                 `json:"condition"`   // 条件（仅if）
	Children   []*SimpleStatement     `json:"children"`    // 子语句（仅if）
	Metadata   map[string]interface{} `json:"metadata"`    // 元数据
	Variables  map[string]interface{} `json:"variables"`   // 当前变量状态
	InputVars  map[string]interface{} `json:"input_vars"`  // 输入变量
	Context    context.Context        `json:"-"`           // 上下文
	StartTime  time.Time              `json:"start_time"`  // 开始时间
	LineNumber int                    `json:"line_number"` // 行号
}

// 语句执行结果
type StatementCallbackResult struct {
	Success    bool          `json:"success"`     // 是否成功
	Error      error         `json:"error"`       // 错误信息
	OutputArgs []interface{} `json:"output_args"` // 输出参数值（仅function-call）
	ShouldSkip bool          `json:"should_skip"` // 是否跳过执行（仅if）
	Duration   int64         `json:"duration_ms"` // 执行时长(毫秒)
}

// 语句类型回调函数
type StatementCallback func(ctx *StatementContext) *StatementCallbackResult

// 类型回调执行器
type TypeCallbackExecutor struct {
	callbacks map[string]StatementCallback
}

// 创建类型回调执行器
func NewTypeCallbackExecutor() *TypeCallbackExecutor {
	return &TypeCallbackExecutor{
		callbacks: make(map[string]StatementCallback),
	}
}

// 注册语句类型回调
func (e *TypeCallbackExecutor) RegisterType(statementType string, callback StatementCallback) {
	e.callbacks[statementType] = callback
}

// 执行语句回调
func (e *TypeCallbackExecutor) ExecuteStatement(ctx *StatementContext) *StatementCallbackResult {
	callback, exists := e.callbacks[ctx.Type]
	if !exists {
		// 如果没有注册回调，使用默认处理
		return e.getDefaultResult(ctx)
	}

	// 执行回调
	result := callback(ctx)
	if result == nil {
		result = &StatementCallbackResult{
			Success: false,
			Error:   fmt.Errorf("语句类型 %s 的回调函数返回nil", ctx.Type),
		}
	}

	return result
}

// 获取默认结果
func (e *TypeCallbackExecutor) getDefaultResult(ctx *StatementContext) *StatementCallbackResult {
	switch ctx.Type {
	case "function-call":
		// 默认函数调用：模拟执行
		time.Sleep(100 * time.Millisecond)
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
				outputArgs[i] = nil
			default:
				outputArgs[i] = "模拟结果"
			}
		}
		return &StatementCallbackResult{
			Success:    true,
			Error:      nil,
			OutputArgs: outputArgs,
			Duration:   time.Since(ctx.StartTime).Milliseconds(),
		}
	case "if":
		// 默认条件判断：总是执行
		return &StatementCallbackResult{
			Success:    true,
			Error:      nil,
			ShouldSkip: false,
			Duration:   time.Since(ctx.StartTime).Milliseconds(),
		}
	case "var":
		// 默认变量赋值：直接成功
		return &StatementCallbackResult{
			Success:  true,
			Error:    nil,
			Duration: time.Since(ctx.StartTime).Milliseconds(),
		}
	case "print":
		// 默认打印：直接成功
		return &StatementCallbackResult{
			Success:  true,
			Error:    nil,
			Duration: time.Since(ctx.StartTime).Milliseconds(),
		}
	case "return":
		// 默认返回：直接成功
		return &StatementCallbackResult{
			Success:  true,
			Error:    nil,
			Duration: time.Since(ctx.StartTime).Milliseconds(),
		}
	default:
		return &StatementCallbackResult{
			Success:  false,
			Error:    fmt.Errorf("未知的语句类型: %s", ctx.Type),
			Duration: time.Since(ctx.StartTime).Milliseconds(),
		}
	}
}
