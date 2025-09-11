package workflow

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type ExecutorIn struct {
	StepName   string                 `json:"step_name"`   // 当前步骤名
	StepDesc   string                 `json:"step_desc"`   // 步骤描述
	RealInput  map[string]interface{} `json:"real_input"`  // 实际输入参数
	WantParams []ParameterInfo        `json:"want_params"` // 预期返回参数信息
	Options    *ExecutorOptions       `json:"options"`     // 执行选项
}

// 执行器选项
type ExecutorOptions struct {
	Timeout    *time.Duration `json:"timeout"`     // 超时时间，nil表示无超时限制
	RetryCount int            `json:"retry_count"` // 重试次数
	Async      bool           `json:"async"`       // 是否异步执行
	Priority   int            `json:"priority"`    // 优先级
	Debug      bool           `json:"debug"`       // 是否调试模式
	LogLevel   string         `json:"log_level"`   // 日志级别
	AIModel    string         `json:"ai_model"`    // AI模型
}

type ExecutorOut struct {
	Success    bool                   `json:"success"`     // 执行是否成功
	WantOutput map[string]interface{} `json:"want_output"` // 输出参数
	Error      string                 `json:"error"`       // 错误信息
	Logs       []string               `json:"logs"`        // 执行日志
}

// OnFunctionCall 执行function-call 回调
type OnFunctionCall func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error)

// OnWorkFlowUpdate 每次执行后我们需要回调，业务侧需要保存整个工作流的状态
type OnWorkFlowUpdate func(ctx context.Context, current *SimpleParseResult) error

// OnWorkFlowExit 正常结束
type OnWorkFlowExit func(ctx context.Context, current *SimpleParseResult) error

// OnWorkFlowReturn 走到某个节点中断了
type OnWorkFlowReturn func(ctx context.Context, current *SimpleParseResult) error

type Executor struct {
	OnFunctionCall   OnFunctionCall //函数执行回调
	OnWorkFlowUpdate OnWorkFlowUpdate
	OnWorkFlowExit   OnWorkFlowExit
	OnWorkFlowReturn OnWorkFlowReturn

	// 流程管理
	FlowMap      map[string]*SimpleParseResult
	RunningFlows map[string]context.CancelFunc // 正在运行的流程
}

type ExecutorResp struct {
}

// NewExecutor 创建新的执行器
func NewExecutor() *Executor {
	return &Executor{
		FlowMap:      make(map[string]*SimpleParseResult),
		RunningFlows: make(map[string]context.CancelFunc),
	}
}

// Start 启动工作流执行
func (e *Executor) Start(ctx context.Context, workflow *SimpleParseResult) error {
	// 1. 检查流程是否已经在运行
	if _, exists := e.RunningFlows[workflow.FlowID]; exists {
		return fmt.Errorf("流程 %s 已在运行中", workflow.FlowID)
	}

	// 2. 创建子上下文用于取消
	flowCtx, cancel := context.WithCancel(ctx)
	e.RunningFlows[workflow.FlowID] = cancel
	defer func() {
		delete(e.RunningFlows, workflow.FlowID)
		cancel()
	}()

	// 3. 保存流程到映射表
	e.FlowMap[workflow.FlowID] = workflow

	// 4. 执行主函数语句
	return e.executeMainFunction(flowCtx, workflow)
}

// Stop 停止指定流程
func (e *Executor) Stop(ctx context.Context, flowId string) error {
	if cancel, exists := e.RunningFlows[flowId]; exists {
		cancel()
		delete(e.RunningFlows, flowId)
		return nil
	}
	return fmt.Errorf("流程 %s 未在运行中", flowId)
}

// executeMainFunction 执行主函数
func (e *Executor) executeMainFunction(ctx context.Context, workflow *SimpleParseResult) error {
	if workflow.MainFunc == nil {
		return fmt.Errorf("工作流没有主函数")
	}

	// 遍历执行每个语句
	for _, stmt := range workflow.MainFunc.Statements {
		// 检查取消信号（服务侧取消，如服务器关机等）
		select {
		case <-ctx.Done():
			// 标记当前语句为取消状态
			stmt.Status = "cancelled"
			stmt.EndExecution()
			// 触发兜底回调，更新节点状态
			if e.OnWorkFlowUpdate != nil {
				_ = e.OnWorkFlowUpdate(ctx, workflow)
			}
			return fmt.Errorf("工作流执行被取消: %v", ctx.Err())
		default:
		}

		// 开始执行计时
		stmt.StartExecution()

		// 触发状态更新回调
		if e.OnWorkFlowUpdate != nil {
			if err := e.OnWorkFlowUpdate(ctx, workflow); err != nil {
				stmt.EndExecution()
				return err
			}
		}

		// 执行语句
		if err := e.executeStatement(ctx, stmt, workflow); err != nil {
			// 检查是否是取消错误
			if strings.Contains(err.Error(), "被取消") || strings.Contains(err.Error(), "被主动取消") {
				// 取消错误，状态已经在 executeStatement 中设置
				return err
			}

			// 执行失败，设置状态为失败
			stmt.EndExecution()
			stmt.Status = "failed"

			// 触发状态更新回调
			if e.OnWorkFlowUpdate != nil {
				if err := e.OnWorkFlowUpdate(ctx, workflow); err != nil {
					return err
				}
			}

			return err
		}

		// 结束执行计时
		stmt.EndExecution()
	}

	// 正常结束
	if e.OnWorkFlowExit != nil {
		return e.OnWorkFlowExit(ctx, workflow)
	}

	return nil
}

// executeStatement 执行单个语句
func (e *Executor) executeStatement(ctx context.Context, stmt *SimpleStatement, workflow *SimpleParseResult) error {
	switch stmt.Type {
	case "function-call":
		return e.executeFunctionCall(ctx, stmt, workflow)
	case "if":
		return e.executeIfStatement(ctx, stmt, workflow)
	// 注意：已移除打印语句支持，执行引擎会自动处理日志记录
	case "var":
		return e.executeVarStatement(ctx, stmt, workflow)
	case "return":
		return e.executeReturnStatement(ctx, stmt, workflow)
	default:
		// 其他类型语句暂时跳过
		return nil
	}
}

// executeFunctionCall 执行函数调用
func (e *Executor) executeFunctionCall(ctx context.Context, stmt *SimpleStatement, workflow *SimpleParseResult) error {
	// 1. 找到对应的步骤定义
	var step *SimpleStep
	for _, s := range workflow.Steps {
		if s.Name == stmt.Function {
			step = s
			break
		}
	}

	if step == nil {
		return fmt.Errorf("未找到步骤定义: %s", stmt.Function)
	}

	// 2. 获取元数据配置
	timeout := stmt.GetTimeout()
	retryCount := stmt.GetRetryCount()
	isDebug := stmt.IsDebug()

	// 3. 构建输入参数
	realInput := make(map[string]interface{})
	for i, arg := range stmt.Args {
		// 获取对应的输入参数定义
		if i < len(step.InputParams) {
			paramDef := step.InputParams[i]

			// 处理不同类型的参数
			if arg.IsInput {
				// 输入参数：input["用户名"] -> 从InputVars中获取
				if strings.HasPrefix(arg.Value, "input[") && strings.HasSuffix(arg.Value, "]") {
					// 提取键名
					key := strings.TrimSpace(arg.Value[6 : len(arg.Value)-1])
					key = strings.Trim(key, "\"")
					if value, exists := workflow.InputVars[key]; exists {
						realInput[paramDef.Name] = value
					} else {
						realInput[paramDef.Name] = arg.Value
					}
				} else {
					realInput[paramDef.Name] = arg.Value
				}
			} else if varInfo, exists := workflow.Variables[arg.Value]; exists {
				// 从变量映射中获取实际值
				realInput[paramDef.Name] = varInfo.Value
			} else {
				// 如果是字面量，直接使用
				realInput[paramDef.Name] = arg.Value
			}
		}
	}

	// 4. 超时控制现在由业务回调自己处理，通过元数据传递超时信息

	// 5. 执行重试逻辑
	var lastErr error
	for attempt := 0; attempt <= retryCount; attempt++ {
		// 检查原始上下文取消信号（服务侧取消，如服务器关机等）
		select {
		case <-ctx.Done():
			// 触发兜底回调，更新节点状态为取消
			stmt.Status = "cancelled"
			stmt.EndExecution()
			if e.OnWorkFlowUpdate != nil {
				// 异步执行，不阻塞主流程
				go func() {
					_ = e.OnWorkFlowUpdate(ctx, workflow)
				}()
			}
			return fmt.Errorf("步骤执行被取消: %v", ctx.Err())
		default:
		}

		// 超时控制现在由业务回调自己处理

		// 构建执行选项
		options := &ExecutorOptions{
			Timeout:    timeout,
			RetryCount: retryCount,
			Async:      stmt.IsAsync(),
			Priority:   stmt.GetPriority(),
			Debug:      isDebug,
			LogLevel:   stmt.GetLogLevel(),
			AIModel:    stmt.GetAIModel(),
		}

		// 调用业务回调
		executorIn := &ExecutorIn{
			StepName:   step.Name,
			StepDesc:   stmt.Desc,
			RealInput:  realInput,
			WantParams: step.OutputParams, // 从步骤定义中获取预期返回参数
			Options:    options,
		}

		// 如果是调试模式，记录调试信息
		if isDebug {
			timeoutStr := "无限制"
			if timeout != nil {
				timeoutStr = timeout.String()
			}
			fmt.Printf("【print】调试模式 - 执行步骤: %s, 尝试次数: %d/%d, 超时: %s\n",
				step.Name, attempt+1, retryCount+1, timeoutStr)
		}

		executorOut, err := e.OnFunctionCall(ctx, *step, executorIn)
		if err != nil {
			// 检查是否是上下文相关错误（取消或超时）
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || ctx.Err() != nil {
				// 上下文被取消或超时
				stmt.Status = "cancelled"
				stmt.EndExecution()
				if e.OnWorkFlowUpdate != nil {
					_ = e.OnWorkFlowUpdate(ctx, workflow)
				}
				return fmt.Errorf("步骤执行被取消: %v", err)
			}

			// 业务逻辑错误，检查err_continue元数据
			errContinue := false
			if step.Metadata != nil {
				if val, exists := step.Metadata["err_continue"]; exists {
					if boolVal, ok := val.(bool); ok {
						errContinue = boolVal
					}
				}
			}

			lastErr = err

			if errContinue {
				// err_continue: true - 记录错误但继续执行
				stmt.Status = "failed_continue"

				// 记录错误日志
				log := &StepLog{
					Timestamp: time.Now(),
					Level:     "error",
					Message:   fmt.Sprintf("步骤执行失败但继续执行: %v", err),
					Source:    step.Name + ".Error",
				}
				step.Logs = append(step.Logs, log)

				// 触发状态更新回调
				if e.OnWorkFlowUpdate != nil {
					if err := e.OnWorkFlowUpdate(ctx, workflow); err != nil {
						return err
					}
				}

				// 继续执行，不返回错误
				return nil
			} else {
				// err_continue: false 或不设置 - 执行失败时终止工作流
				if attempt < retryCount {
					// 还有重试机会，等待一段时间后重试
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}
				return err
			}
		}

		// 业务回调执行后，再次检查原始上下文是否被取消
		select {
		case <-ctx.Done():
			stmt.Status = "cancelled"
			stmt.EndExecution()
			if e.OnWorkFlowUpdate != nil {
				_ = e.OnWorkFlowUpdate(ctx, workflow)
			}
			return fmt.Errorf("步骤执行被取消: %v", ctx.Err())
		default:
		}

		// 6. 处理输出参数
		if executorOut.Success {
			// 将输出参数存储到变量映射中
			// 需要将形参名映射到实例名
			for i, returnVar := range stmt.Returns {
				// 获取对应的输出参数定义
				if i < len(step.OutputParams) {
					paramDef := step.OutputParams[i]
					// 从WantOutput中获取形参名对应的值
					if value, exists := executorOut.WantOutput[paramDef.Name]; exists {
						// 使用实例名作为变量名，而不是形参名
						workflow.Variables[returnVar.Value] = VariableInfo{
							Name:    returnVar.Value, // 实例名：工号、用户名、step1Err
							Type:    paramDef.Type,   // 从步骤定义获取类型
							Value:   value,           // 实际值
							Source:  stmt.Function,
							LineNum: stmt.LineNumber,
							IsInput: false,
						}
					}
				}
			}

			// 执行成功，更新状态为完成
			stmt.Status = "completed"

			// 触发状态更新回调
			if e.OnWorkFlowUpdate != nil {
				if err := e.OnWorkFlowUpdate(ctx, workflow); err != nil {
					return err
				}
			}

			return nil
		} else {
			// 执行失败，检查err_continue元数据
			errContinue := false
			if step.Metadata != nil {
				if val, exists := step.Metadata["err_continue"]; exists {
					if boolVal, ok := val.(bool); ok {
						errContinue = boolVal
					}
				}
			}

			lastErr = fmt.Errorf("步骤执行失败: %s", executorOut.Error)

			if errContinue {
				// err_continue: true - 记录错误但继续执行
				stmt.Status = "failed_continue"

				// 记录错误日志
				log := &StepLog{
					Timestamp: time.Now(),
					Level:     "error",
					Message:   fmt.Sprintf("步骤执行失败但继续执行: %s", executorOut.Error),
					Source:    step.Name + ".Error",
				}
				step.Logs = append(step.Logs, log)

				// 触发状态更新回调
				if e.OnWorkFlowUpdate != nil {
					if err := e.OnWorkFlowUpdate(ctx, workflow); err != nil {
						return err
					}
				}

				// 继续执行，不返回错误
				return nil
			} else {
				// err_continue: false 或不设置 - 执行失败时终止工作流
				if attempt < retryCount {
					// 还有重试机会，等待一段时间后重试
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}
			}
		}
	}

	// 所有重试都失败了
	stmt.Status = "failed"

	// 触发状态更新回调
	if e.OnWorkFlowUpdate != nil {
		if err := e.OnWorkFlowUpdate(ctx, workflow); err != nil {
			return err
		}
	}

	return lastErr
}

// executeIfStatement 执行if语句
func (e *Executor) executeIfStatement(ctx context.Context, stmt *SimpleStatement, workflow *SimpleParseResult) error {
	// 1. 解析条件表达式
	condition := stmt.Condition
	if condition == "" {
		return fmt.Errorf("if语句缺少条件表达式")
	}

	// 2. 评估条件
	result, err := e.evaluateCondition(condition, workflow.Variables)
	if err != nil {
		return fmt.Errorf("条件评估失败: %v", err)
	}

	// 3. 根据条件结果执行子语句
	if result {
		// 条件为真，执行子语句
		for _, child := range stmt.Children {
			// 为子语句开始计时
			child.StartExecution()
			if err := e.executeStatement(ctx, child, workflow); err != nil {
				child.EndExecution()
				return err
			}
			child.EndExecution()
		}
	}

	// 4. 更新语句状态为完成
	stmt.Status = "completed"

	// 5. 触发状态更新回调
	if e.OnWorkFlowUpdate != nil {
		if err := e.OnWorkFlowUpdate(ctx, workflow); err != nil {
			return err
		}
	}

	return nil
}

// 注意：已移除 executePrintStatement 函数，执行引擎会自动处理日志记录

// executeVarStatement 执行var语句
func (e *Executor) executeVarStatement(ctx context.Context, stmt *SimpleStatement, workflow *SimpleParseResult) error {
	// 1. 解析变量赋值语句
	content := stmt.Content
	if content == "" {
		return fmt.Errorf("var语句缺少内容")
	}

	// 2. 解析变量名和值
	varName, varValue, err := e.parseVarAssignment(content)
	if err != nil {
		return fmt.Errorf("变量赋值解析失败: %v", err)
	}

	// 3. 处理模板变量替换
	processedValue, err := e.processTemplate(varValue, workflow.Variables)
	if err != nil {
		return fmt.Errorf("模板处理失败: %v", err)
	}

	// 4. 创建变量信息
	varInfo := VariableInfo{
		Name:    varName,
		Type:    "string", // 默认为string类型，后续可以改进类型推断
		Value:   processedValue,
		Source:  "assignment",
		LineNum: stmt.LineNumber,
		IsInput: false,
	}

	// 5. 存储到变量映射
	workflow.Variables[varName] = varInfo

	// 6. 更新语句状态为完成
	stmt.Status = "completed"

	// 7. 触发状态更新回调
	if e.OnWorkFlowUpdate != nil {
		if err := e.OnWorkFlowUpdate(ctx, workflow); err != nil {
			return err
		}
	}

	return nil
}

// executeReturnStatement 执行return语句
func (e *Executor) executeReturnStatement(ctx context.Context, stmt *SimpleStatement, workflow *SimpleParseResult) error {
	// 1. 更新语句状态为完成
	stmt.Status = "completed"

	// 2. 触发状态更新回调
	if e.OnWorkFlowUpdate != nil {
		if err := e.OnWorkFlowUpdate(ctx, workflow); err != nil {
			return err
		}
	}

	// 3. 触发return回调
	if e.OnWorkFlowReturn != nil {
		return e.OnWorkFlowReturn(ctx, workflow)
	}
	return nil
}

// evaluateCondition 评估条件表达式
func (e *Executor) evaluateCondition(condition string, variables map[string]VariableInfo) (bool, error) {
	// 简单的条件评估，支持常见的比较操作
	// 例如: "step1Err != nil", "验证通过 == true", "count > 0"

	// 处理 != nil 条件
	if strings.Contains(condition, "!= nil") {
		varName := strings.TrimSpace(strings.Split(condition, "!=")[0])
		if varInfo, exists := variables[varName]; exists {
			return varInfo.Value == nil, nil
		}
		return false, nil
	}

	// 处理 == true 条件
	if strings.Contains(condition, "== true") {
		varName := strings.TrimSpace(strings.Split(condition, "==")[0])
		if varInfo, exists := variables[varName]; exists {
			return varInfo.Value == true, nil
		}
		return false, nil
	}

	// 处理 == false 条件
	if strings.Contains(condition, "== false") {
		varName := strings.TrimSpace(strings.Split(condition, "==")[0])
		if varInfo, exists := variables[varName]; exists {
			return varInfo.Value == false, nil
		}
		return true, nil
	}

	// 处理 != true 条件
	if strings.Contains(condition, "!= true") {
		varName := strings.TrimSpace(strings.Split(condition, "!=")[0])
		if varInfo, exists := variables[varName]; exists {
			return varInfo.Value != true, nil
		}
		return true, nil
	}

	// 默认返回false
	return false, nil
}

// processTemplate 处理模板变量替换
func (e *Executor) processTemplate(content string, variables map[string]VariableInfo) (string, error) {
	// 处理 {{变量名}} 格式的模板
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	result := re.ReplaceAllStringFunc(content, func(match string) string {
		// 提取变量名
		varName := strings.TrimSpace(match[2 : len(match)-2])

		// 查找变量
		if varInfo, exists := variables[varName]; exists {
			return fmt.Sprintf("%v", varInfo.Value)
		}

		// 变量不存在，保持原样
		return match
	})

	return result, nil
}

// 注意：已移除 extractStepNameFromPrint 函数，不再需要打印语句支持

// parseVarAssignment 解析变量赋值语句
func (e *Executor) parseVarAssignment(content string) (string, string, error) {
	// 匹配 变量名 := "值" 格式，支持中文变量名
	re := regexp.MustCompile(`([\w\p{Han}]+)\s*:=\s*(.+)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 3 {
		return "", "", fmt.Errorf("无法解析变量赋值语句: %s", content)
	}

	varName := strings.TrimSpace(matches[1])
	varValue := strings.TrimSpace(matches[2])

	// 去掉引号
	if strings.HasPrefix(varValue, `"`) && strings.HasSuffix(varValue, `"`) {
		varValue = varValue[1 : len(varValue)-1]
	}

	return varName, varValue, nil
}

// Get 获取指定工作流的全部信息
func (e *Executor) Get(flowID string) (*SimpleParseResult, error) {
	// 检查工作流是否存在
	workflow, exists := e.FlowMap[flowID]
	if !exists {
		return nil, fmt.Errorf("工作流不存在: %s", flowID)
	}

	// 直接返回工作流信息
	return workflow, nil
}
