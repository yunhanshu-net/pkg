package workflow

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type ExecutorIn struct {
	StepName  string                 `json:"step_name"`  // 当前步骤名
	StepDesc  string                 `json:"step_desc"`  // 步骤描述
	RealInput map[string]interface{} `json:"real_input"` // 实际输入参数
	Metadata  map[string]interface{} `json:"metadata"`   // 步骤元数据
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
	cancelCtx    context.CancelFunc
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
		// 检查取消信号
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 设置语句状态为运行中
		stmt.Status = "running"

		// 触发状态更新回调
		if e.OnWorkFlowUpdate != nil {
			if err := e.OnWorkFlowUpdate(ctx, workflow); err != nil {
				return err
			}
		}

		// 执行语句
		if err := e.executeStatement(ctx, stmt, workflow); err != nil {
			// 执行失败，设置状态为失败
			stmt.Status = "failed"

			// 触发状态更新回调
			if e.OnWorkFlowUpdate != nil {
				if err := e.OnWorkFlowUpdate(ctx, workflow); err != nil {
					return err
				}
			}

			return err
		}
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
	case "print":
		return e.executePrintStatement(ctx, stmt, workflow)
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

	// 2. 构建输入参数
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

	// 3. 调用业务回调
	executorIn := &ExecutorIn{
		StepName:  step.Name,
		StepDesc:  stmt.Desc,
		RealInput: realInput,
		Metadata:  stmt.Metadata,
	}

	executorOut, err := e.OnFunctionCall(ctx, *step, executorIn)
	if err != nil {
		return err
	}

	// 4. 处理输出参数
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
	} else {
		// 执行失败，更新状态为失败
		stmt.Status = "failed"

		// 触发状态更新回调
		if e.OnWorkFlowUpdate != nil {
			if err := e.OnWorkFlowUpdate(ctx, workflow); err != nil {
				return err
			}
		}

		return fmt.Errorf("步骤执行失败: %s", executorOut.Error)
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
			if err := e.executeStatement(ctx, child, workflow); err != nil {
				return err
			}
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

// executePrintStatement 执行print语句
func (e *Executor) executePrintStatement(ctx context.Context, stmt *SimpleStatement, workflow *SimpleParseResult) error {
	// 1. 解析打印内容
	content := stmt.Content
	if content == "" {
		return fmt.Errorf("print语句缺少内容")
	}

	// 2. 处理变量替换
	processedContent, err := e.processTemplate(content, workflow.Variables)
	if err != nil {
		return fmt.Errorf("模板处理失败: %v", err)
	}

	// 3. 记录到步骤日志
	stepName := e.extractStepNameFromPrint(content)
	if stepName != "" {
		// 找到对应的步骤并添加日志
		for _, step := range workflow.Steps {
			if step.Name == stepName {
				log := &StepLog{
					Timestamp: time.Now(),
					Level:     "info",
					Message:   processedContent,
					Source:    stepName + ".Printf",
				}
				step.Logs = append(step.Logs, log)
				break
			}
		}
	} else {
		// 全局日志
		log := &StepLog{
			Timestamp: time.Now(),
			Level:     "info",
			Message:   processedContent,
			Source:    "sys.Print",
		}
		workflow.GlobalLogs = append(workflow.GlobalLogs, log)
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

// extractStepNameFromPrint 从print语句中提取步骤名
func (e *Executor) extractStepNameFromPrint(content string) string {
	// 匹配 step1.Printf("...") 格式
	re := regexp.MustCompile(`(\w+)\.Printf\(`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

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
