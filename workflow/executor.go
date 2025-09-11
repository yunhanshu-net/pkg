package workflow

import (
	"fmt"
	"strings"
	"time"
)

// 执行引擎
type WorkflowExecutor struct {
	parser   *SimpleParser
	callback *TypeCallbackExecutor
}

// 执行结果
type ExecutionResult struct {
	Success   bool                   `json:"success"`     // 执行是否成功
	StartTime time.Time              `json:"start_time"`  // 开始时间
	EndTime   time.Time              `json:"end_time"`    // 结束时间
	Duration  int64                  `json:"duration_ms"` // 执行时长(毫秒)
	Steps     []StepExecutionResult  `json:"steps"`       // 步骤执行结果
	Error     string                 `json:"error"`       // 错误信息
	InputVars map[string]interface{} `json:"input_vars"`  // 输入变量
	Variables map[string]interface{} `json:"variables"`   // 执行过程中的变量
}

// 步骤执行结果
type StepExecutionResult struct {
	StepName   string                 `json:"step_name"`   // 步骤名称
	Function   string                 `json:"function"`    // 函数名
	StartTime  time.Time              `json:"start_time"`  // 开始时间
	EndTime    time.Time              `json:"end_time"`    // 结束时间
	Duration   int64                  `json:"duration_ms"` // 执行时长(毫秒)
	Success    bool                   `json:"success"`     // 是否成功
	Error      string                 `json:"error"`       // 错误信息
	InputArgs  []interface{}          `json:"input_args"`  // 输入参数
	OutputArgs []interface{}          `json:"output_args"` // 输出参数
	Metadata   map[string]interface{} `json:"metadata"`    // 元数据
}

// 创建执行引擎
func NewWorkflowExecutor() *WorkflowExecutor {
	return &WorkflowExecutor{
		parser:   NewSimpleParser(),
		callback: NewTypeCallbackExecutor(),
	}
}

// 创建带类型回调执行器的执行引擎
func NewWorkflowExecutorWithCallback(callback *TypeCallbackExecutor) *WorkflowExecutor {
	return &WorkflowExecutor{
		parser:   NewSimpleParser(),
		callback: callback,
	}
}

// 注册语句类型回调
func (e *WorkflowExecutor) RegisterType(statementType string, callback StatementCallback) {
	e.callback.RegisterType(statementType, callback)
}

// 执行工作流
func (e *WorkflowExecutor) ExecuteWorkflow(code string) *ExecutionResult {
	startTime := time.Now()

	// 解析工作流
	parseResult := e.parser.ParseWorkflow(code)
	if !parseResult.Success {
		return &ExecutionResult{
			Success:   false,
			StartTime: startTime,
			EndTime:   time.Now(),
			Duration:  time.Since(startTime).Milliseconds(),
			Error:     parseResult.Error,
			Variables: make(map[string]interface{}),
		}
	}

	// 执行工作流并返回结果
	return e.executeWorkflowWithResult(parseResult, startTime)
}

// 执行工作流并返回更新后的解析结果
func (e *WorkflowExecutor) ExecuteWorkflowWithResult(code string) (*ExecutionResult, *SimpleParseResult) {
	startTime := time.Now()

	// 解析工作流
	parseResult := e.parser.ParseWorkflow(code)
	if !parseResult.Success {
		return &ExecutionResult{
			Success:   false,
			StartTime: startTime,
			EndTime:   time.Now(),
			Duration:  time.Since(startTime).Milliseconds(),
			Error:     parseResult.Error,
			Variables: make(map[string]interface{}),
		}, parseResult
	}

	// 执行工作流并返回结果
	execResult := e.executeWorkflowWithResult(parseResult, startTime)
	return execResult, parseResult
}

// 执行工作流的核心逻辑
func (e *WorkflowExecutor) executeWorkflowWithResult(parseResult *SimpleParseResult, startTime time.Time) *ExecutionResult {

	// 初始化执行结果
	result := &ExecutionResult{
		Success:   true,
		StartTime: startTime,
		InputVars: make(map[string]interface{}),
		Variables: make(map[string]interface{}),
		Steps:     make([]StepExecutionResult, 0),
	}

	// 复制输入变量到专门的输入变量区域
	for key, value := range parseResult.InputVars {
		result.InputVars[key] = value
	}

	// 执行主函数
	err := e.executeMainFunction(parseResult.MainFunc, result, parseResult)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
	}

	result.EndTime = time.Now()
	result.Duration = time.Since(startTime).Milliseconds()

	return result
}

// 执行主函数
func (e *WorkflowExecutor) executeMainFunction(mainFunc *SimpleMainFunc, result *ExecutionResult, parseResult *SimpleParseResult) error {
	fmt.Println("🚀 开始执行工作流...")

	// 执行所有语句
	for _, stmt := range mainFunc.Statements {
		err := e.executeStatement(stmt, result, parseResult)
		if err != nil {
			return err
		}
	}

	fmt.Println("✅ 工作流执行完成！")
	return nil
}

// 执行语句
func (e *WorkflowExecutor) executeStatement(stmt *SimpleStatement, result *ExecutionResult, parseResult *SimpleParseResult) error {
	// 设置状态为正在执行
	stmt.SetStatus(StatusRunning)

	switch stmt.Type {
	case "print":
		return e.executePrintStatement(stmt, result, parseResult)
	case "function-call":
		return e.executeFunctionCall(stmt, result, parseResult)
	case "if":
		return e.executeIfStatement(stmt, result, parseResult)
	case "var":
		return e.executeVarStatement(stmt, result, parseResult)
	case "return":
		return e.executeReturnStatement(stmt, result, parseResult)
	default:
		// 添加全局日志
		parseResult.AddGlobalLog("warn", fmt.Sprintf("跳过未知语句类型: %s", stmt.Type), "system")
		fmt.Printf("⚠️ 跳过未知语句类型: %s\n", stmt.Type)
		stmt.SetStatus(StatusSkipped)
		return nil
	}
}

// 执行打印语句
func (e *WorkflowExecutor) executePrintStatement(stmt *SimpleStatement, result *ExecutionResult, parseResult *SimpleParseResult) error {
	// 检查是否是步骤级别的日志记录
	if strings.Contains(stmt.Content, ".Printf") || strings.Contains(stmt.Content, ".Println") {
		// 解析步骤名称和日志内容
		stepName, logMessage := e.parseStepLog(stmt.Content)

		// 找到对应的步骤并添加日志
		for i := range parseResult.Steps {
			if parseResult.Steps[i].Name == stepName {
				parseResult.Steps[i].AddLog("info", logMessage, stepName+".Printf")
				break
			}
		}
		fmt.Printf("   【%s】%s\n", stepName, logMessage)
	} else if strings.HasPrefix(stmt.Content, "fmt.Print") {
		// 全局日志
		parseResult.AddGlobalLog("info", stmt.Content, "fmt.Print")
		fmt.Printf("   【sys】%s\n", stmt.Content)
	} else {
		// 其他打印语句作为全局日志
		parseResult.AddGlobalLog("info", stmt.Content, "unknown")
		fmt.Printf("   【print】%s\n", stmt.Content)
	}

	stmt.SetStatus(StatusCompleted)
	return nil
}

// 执行函数调用
func (e *WorkflowExecutor) executeFunctionCall(stmt *SimpleStatement, result *ExecutionResult, parseResult *SimpleParseResult) error {
	stepStartTime := time.Now()

	// 记录开始执行的日志到对应步骤
	for i := range parseResult.Steps {
		if parseResult.Steps[i].Name == stmt.Function {
			parseResult.Steps[i].AddLog("info", fmt.Sprintf("开始执行步骤: %s", stmt.Function), stmt.Function)
			break
		}
	}
	fmt.Printf("🔧 [步骤] %s - 函数: %s\n", stmt.Function, stmt.Function)

	// 打印输入参数
	if len(stmt.Args) > 0 {
		fmt.Printf("   📥 输入参数: ")
		for i, arg := range stmt.Args {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%s", arg.Value)
		}
		fmt.Println()
	}

	// 模拟执行步骤
	time.Sleep(100 * time.Millisecond) // 模拟执行时间

	// 模拟执行结果
	stepResult := StepExecutionResult{
		StepName:   stmt.Function,
		Function:   stmt.Function,
		StartTime:  stepStartTime,
		EndTime:    time.Now(),
		Duration:   time.Since(stepStartTime).Milliseconds(),
		Success:    true,
		InputArgs:  make([]interface{}, 0),
		OutputArgs: make([]interface{}, 0),
		Metadata:   stmt.Metadata,
	}

	// 处理输入参数
	for _, arg := range stmt.Args {
		stepResult.InputArgs = append(stepResult.InputArgs, arg.Value)
	}

	// 处理输出参数
	for _, ret := range stmt.Returns {
		// 模拟输出值
		var outputValue interface{}
		switch ret.Type {
		case "string":
			outputValue = fmt.Sprintf("模拟%s结果", ret.Value)
		case "int":
			outputValue = 12345
		case "bool":
			outputValue = true
		case "err":
			outputValue = nil // 模拟成功，无错误
		default:
			outputValue = "模拟结果"
		}

		stepResult.OutputArgs = append(stepResult.OutputArgs, outputValue)

		// 更新变量映射
		result.Variables[ret.Value] = outputValue
	}

	// 打印输出参数
	if len(stepResult.OutputArgs) > 0 {
		fmt.Printf("   📤 输出参数: ")
		for i, output := range stepResult.OutputArgs {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%v", output)
		}
		fmt.Println()
	}

	// 打印元数据
	if len(stmt.Metadata) > 0 {
		fmt.Printf("   ⚙️ 元数据: ")
		for key, value := range stmt.Metadata {
			fmt.Printf("%s=%v ", key, value)
		}
		fmt.Println()
	}

	result.Steps = append(result.Steps, stepResult)

	// 记录执行完成的日志到对应步骤
	for i := range parseResult.Steps {
		if parseResult.Steps[i].Name == stmt.Function {
			parseResult.Steps[i].AddLog("info", fmt.Sprintf("步骤执行完成 (耗时: %dms)", stepResult.Duration), stmt.Function)
			break
		}
	}
	fmt.Printf("   ✅ 步骤执行完成 (耗时: %dms)\n", stepResult.Duration)

	// 设置状态为完成
	stmt.SetStatus(StatusCompleted)

	return nil
}

// 执行if语句
func (e *WorkflowExecutor) executeIfStatement(stmt *SimpleStatement, result *ExecutionResult, parseResult *SimpleParseResult) error {
	// 简单的条件判断逻辑
	// 这里可以根据实际需求实现更复杂的条件判断
	shouldExecute := true

	// 检查是否有错误变量
	if stmt.Condition != "" {
		// 简单的错误检查逻辑
		shouldExecute = e.evaluateCondition(stmt.Condition, result)
		parseResult.AddGlobalLog("info", fmt.Sprintf("条件判断: %s = %v", stmt.Condition, shouldExecute), "system")
	}

	if shouldExecute {
		parseResult.AddGlobalLog("info", "条件为真，执行子语句", "system")
		// 执行子语句
		for _, childStmt := range stmt.Children {
			err := e.executeStatement(childStmt, result, parseResult)
			if err != nil {
				parseResult.AddGlobalLog("error", fmt.Sprintf("子语句执行失败: %v", err), "system")
				stmt.SetStatus(StatusFailed)
				return err
			}
		}
	} else {
		parseResult.AddGlobalLog("info", "条件为假，跳过子语句", "system")
		stmt.SetStatus(StatusSkipped)
	}

	stmt.SetStatus(StatusCompleted)
	return nil
}

// 执行变量赋值语句
func (e *WorkflowExecutor) executeVarStatement(stmt *SimpleStatement, result *ExecutionResult, parseResult *SimpleParseResult) error {
	// 简单的变量赋值逻辑
	// 这里可以根据实际需求实现更复杂的变量处理
	result.Variables["变量值"] = "模拟变量值"

	parseResult.AddGlobalLog("info", "变量赋值完成", "system")
	stmt.SetStatus(StatusCompleted)
	return nil
}

// 执行return语句
func (e *WorkflowExecutor) executeReturnStatement(stmt *SimpleStatement, result *ExecutionResult, parseResult *SimpleParseResult) error {
	parseResult.AddGlobalLog("info", "执行返回语句", "system")
	stmt.SetStatus(StatusCompleted)
	return nil
}

// 获取变量值，优先从过程变量中查找，找不到再从输入变量中查找
func (e *WorkflowExecutor) getVariableValue(varName string, result *ExecutionResult) interface{} {
	// 先从过程变量中查找
	if value, exists := result.Variables[varName]; exists {
		return value
	}
	// 再从输入变量中查找
	if value, exists := result.InputVars[varName]; exists {
		return value
	}
	return nil
}

// 评估条件表达式
func (e *WorkflowExecutor) evaluateCondition(condition string, result *ExecutionResult) bool {
	// 简单的条件评估逻辑
	if condition == "" {
		return true
	}

	// 检查错误条件
	if condition == "step1Err != nil" || condition == "step2Err != nil" || condition == "step3Err != nil" {
		// 检查对应的错误变量是否存在且不为nil
		if condition == "step1Err != nil" {
			if err := e.getVariableValue("step1Err", result); err != nil {
				return err != nil // 如果错误不为nil，条件为真
			}
		} else if condition == "step2Err != nil" {
			if err := e.getVariableValue("step2Err", result); err != nil {
				return err != nil // 如果错误不为nil，条件为真
			}
		} else if condition == "step3Err != nil" {
			if err := e.getVariableValue("step3Err", result); err != nil {
				return err != nil // 如果错误不为nil，条件为真
			}
		}
		return false // 默认条件为假
	}

	// 检查其他条件，比如验证结果
	if condition == "验证结果" {
		if value := e.getVariableValue("验证结果", result); value != nil {
			if boolResult, ok := value.(bool); ok {
				return boolResult // 直接返回布尔值
			}
		}
		return false // 默认条件为假
	}

	// 检查其他布尔条件
	if condition == "工号 != \"\"" {
		if value := e.getVariableValue("工号", result); value != nil {
			if strValue, ok := value.(string); ok {
				return strValue != "" // 检查字符串是否不为空
			}
		}
		return false
	}

	return false
}

// 解析步骤日志
func (e *WorkflowExecutor) parseStepLog(content string) (stepName, logMessage string) {
	// 解析格式：step3.Printf("✅ 面试安排成功，时间: %s\n", 面试时间)
	// 或者：step3.Println("✅ 面试安排成功")

	// 查找第一个点号
	dotIndex := strings.Index(content, ".")
	if dotIndex == -1 {
		return "unknown", content
	}

	stepName = strings.TrimSpace(content[:dotIndex])

	// 查找括号
	parenStart := strings.Index(content, "(")
	parenEnd := strings.LastIndex(content, ")")
	if parenStart == -1 || parenEnd == -1 || parenStart >= parenEnd {
		return stepName, content
	}

	// 提取日志内容（去掉引号）
	logContent := strings.TrimSpace(content[parenStart+1 : parenEnd])
	if strings.HasPrefix(logContent, "\"") && strings.HasSuffix(logContent, "\"") {
		logMessage = logContent[1 : len(logContent)-1]
	} else {
		logMessage = logContent
	}

	return stepName, logMessage
}

// 打印执行结果
func (r *ExecutionResult) Print() {
	fmt.Println("\n📊 执行结果:")
	fmt.Printf("   成功: %v\n", r.Success)
	fmt.Printf("   开始时间: %s\n", r.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("   结束时间: %s\n", r.EndTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("   总耗时: %dms\n", r.Duration)

	if r.Error != "" {
		fmt.Printf("   错误: %s\n", r.Error)
	}

	fmt.Printf("   执行步骤数: %d\n", len(r.Steps))

	for i, step := range r.Steps {
		fmt.Printf("   %d. %s - %s (耗时: %dms)\n", i+1, step.StepName, step.Function, step.Duration)
	}

	fmt.Println("\n📈 执行统计:")
	fmt.Printf("   总步骤数: %d\n", len(r.Steps))
	successCount := 0
	for _, step := range r.Steps {
		if step.Success {
			successCount++
		}
	}
	fmt.Printf("   成功步骤: %d\n", successCount)
	fmt.Printf("   失败步骤: %d\n", len(r.Steps)-successCount)

	// 显示输入变量
	if len(r.InputVars) > 0 {
		fmt.Println("\n📥 输入变量:")
		for key, value := range r.InputVars {
			fmt.Printf("   %s: %v\n", key, value)
		}
	}

	// 显示过程变量
	if len(r.Variables) > 0 {
		fmt.Println("\n📋 过程变量:")
		for key, value := range r.Variables {
			fmt.Printf("   %s: %v\n", key, value)
		}
	}
}
