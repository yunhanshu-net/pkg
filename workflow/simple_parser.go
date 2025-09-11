package workflow

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 简单解析器
type SimpleParser struct{}

// 生成FlowID
func generateFlowID() string {
	return fmt.Sprintf("flow_%d_%d", time.Now().UnixNano(), rand.Intn(10000))
}

// 默认元数据配置
type DefaultMetadata struct {
	Timeout     *time.Duration `json:"timeout"`      // 超时时间，nil表示无超时
	RetryCount  int            `json:"retry_count"`  // 重试次数，默认0
	Async       bool           `json:"async"`        // 是否异步执行，默认false
	Priority    int            `json:"priority"`     // 优先级，默认0
	Debug       bool           `json:"debug"`        // 是否调试模式，默认false
	LogLevel    string         `json:"log_level"`    // 日志级别，默认info
	AIModel     string         `json:"ai_model"`     // AI模型，默认空
	ErrContinue bool           `json:"err_continue"` // 错误时是否继续执行，默认false（出错终止）
}

// 获取默认元数据
func GetDefaultMetadata() DefaultMetadata {
	return DefaultMetadata{
		Timeout:     nil, // 默认无超时限制
		RetryCount:  0,
		Async:       false,
		Priority:    0,
		Debug:       false,
		LogLevel:    "info",
		AIModel:     "",
		ErrContinue: false, // 默认出错终止
	}
}

// 解析结果
type SimpleParseResult struct {
	FlowID string `json:"flow_id"`

	Success    bool                    `json:"success"`     // 解析是否成功
	InputVars  map[string]interface{}  `json:"input_vars"`  // 输入变量
	Steps      []*SimpleStep           `json:"steps"`       // 工作流步骤
	MainFunc   *SimpleMainFunc         `json:"main_func"`   // 主函数
	Variables  map[string]VariableInfo `json:"variables"`   // 变量映射表
	GlobalLogs []*StepLog              `json:"global_logs"` // 全局日志
	Error      string                  `json:"error"`       // 错误信息
}

// 变量信息
type VariableInfo struct {
	Name    string      `json:"name"`     // 变量名
	Type    string      `json:"type"`     // 变量类型
	Value   interface{} `json:"value"`    // 变量值
	Source  string      `json:"source"`   // 来源函数名
	LineNum int         `json:"line_num"` // 定义行号
	IsInput bool        `json:"is_input"` // 是否来自input
}

// 参数信息结构体
type ArgumentInfo struct {
	Value      string `json:"value"`       // 参数值
	Type       string `json:"type"`        // 参数类型
	Desc       string `json:"desc"`        // 参数描述
	IsVariable bool   `json:"is_variable"` // 是否为变量引用
	IsLiteral  bool   `json:"is_literal"`  // 是否为字面量
	IsInput    bool   `json:"is_input"`    // 是否为输入参数
	Source     string `json:"source"`      // 来源（变量名或函数名）
	LineNum    int    `json:"line_num"`    // 定义行号
}

// 参数定义信息（用于步骤定义中的参数）
type ParameterInfo struct {
	Name string `json:"name"` // 英文参数名
	Type string `json:"type"` // 参数类型
	Desc string `json:"desc"` // 中文描述
}

// 主函数
type SimpleMainFunc struct {
	Statements []*SimpleStatement `json:"statements"` // 语句列表
}

// 语句状态
type StatementStatus string

const (
	StatusPending   StatementStatus = "pending"   // 待执行
	StatusRunning   StatementStatus = "running"   // 正在执行
	StatusCompleted StatementStatus = "completed" // 执行完成
	StatusFailed    StatementStatus = "failed"    // 执行失败
	StatusSkipped   StatementStatus = "skipped"   // 跳过执行
)

// 步骤日志
type StepLog struct {
	Timestamp time.Time `json:"timestamp"` // 日志时间
	Level     string    `json:"level"`     // 日志级别 (info, warn, error)
	Message   string    `json:"message"`   // 日志内容
	Source    string    `json:"source"`    // 日志来源 (step1.Printf, fmt.Print等)
}

// 语句
type SimpleStatement struct {
	Type       string                 `json:"type"`        // 语句类型
	Content    string                 `json:"content"`     // 语句内容
	LineNumber int                    `json:"line_number"` // 行号
	Children   []*SimpleStatement     `json:"children"`    // 嵌套语句，如if语句的body
	Condition  string                 `json:"condition"`   // 条件表达式，如if语句的条件
	Function   string                 `json:"function"`    // 函数名，如step1()
	Args       []*ArgumentInfo        `json:"args"`        // 函数输入参数信息
	Returns    []*ArgumentInfo        `json:"returns"`     // 函数输出参数信息
	Metadata   map[string]interface{} `json:"metadata"`    // 元数据配置，如 {retry:1, timeout:2000}
	Status     StatementStatus        `json:"status"`      // 执行状态
	RetryCount int                    `json:"retry_count"` // 重试次数
	Desc       string                 `json:"desc"`        // 步骤描述信息
	StartTime  *time.Time             `json:"start_time"`  // 开始执行时间
	EndTime    *time.Time             `json:"end_time"`    // 结束执行时间
	Duration   time.Duration          `json:"duration"`    // 执行耗时
}

// 简单步骤定义
type SimpleStep struct {
	Name         string                 `json:"name"`          // 步骤名称
	Function     string                 `json:"function"`      // 函数名
	InputParams  []ParameterInfo        `json:"input_params"`  // 输入参数定义
	OutputParams []ParameterInfo        `json:"output_params"` // 输出参数定义
	IsStatic     bool                   `json:"is_static"`     // 是否为静态工作流
	CaseID       string                 `json:"case_id"`       // 用例ID
	Logs         []*StepLog             `json:"logs"`          // 步骤日志
	Desc         string                 `json:"desc"`          // 步骤描述信息
	Metadata     map[string]interface{} `json:"metadata"`      // 元数据配置
}

// 添加步骤日志
func (s *SimpleStep) AddLog(level, message, source string) {
	log := &StepLog{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Source:    source,
	}
	s.Logs = append(s.Logs, log)
}

// 添加全局日志
func (r *SimpleParseResult) AddGlobalLog(level, message, source string) {
	log := &StepLog{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Source:    source,
	}
	r.GlobalLogs = append(r.GlobalLogs, log)
}

// 简单类型定义
type SimpleTypeDef struct {
	Type string `json:"type"` // 类型
	Name string `json:"name"` // 名称
}

// 创建简单解析器
func NewSimpleParser() *SimpleParser {
	return &SimpleParser{}
}

// 解析工作流
func (p *SimpleParser) ParseWorkflow(code string) *SimpleParseResult {
	result := &SimpleParseResult{
		FlowID:     generateFlowID(),
		Success:    true,
		InputVars:  make(map[string]interface{}),
		Steps:      []*SimpleStep{},
		MainFunc:   &SimpleMainFunc{Statements: []*SimpleStatement{}},
		Variables:  make(map[string]VariableInfo),
		GlobalLogs: make([]*StepLog, 0),
	}

	// 检查空代码
	code = strings.TrimSpace(code)
	if code == "" {
		result.Success = false
		result.Error = "代码为空"
		return result
	}

	lines := strings.Split(code, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// 解析输入变量
		if strings.HasPrefix(line, "var input") {
			inputVars, err := p.parseInputVars(code)
			if err != nil {
				result.Success = false
				result.Error = err.Error()
				return result
			}
			result.InputVars = inputVars
			continue
		}

		// 解析步骤定义 - 单行格式
		if strings.Contains(line, "=") && (strings.Contains(line, "->") || strings.Contains(line, "beiluo.")) {
			step, err := p.parseStep(line, i+1, lines)
			if err != nil {
				result.Success = false
				result.Error = err.Error()
				return result
			}
			result.Steps = append(result.Steps, step)
		}

		// 解析main函数
		if strings.HasPrefix(line, "func main()") {
			mainFunc, err := p.parseMainFunction(code, lines, result)
			if err != nil {
				result.Success = false
				result.Error = err.Error()
				return result
			}
			result.MainFunc = mainFunc
		}
	}

	// 检查是否有main函数
	if result.MainFunc == nil || len(result.MainFunc.Statements) == 0 {
		result.Success = false
		result.Error = "缺少main函数"
		return result
	}

	return result
}

// 解析输入变量
func (p *SimpleParser) parseInputVars(code string) (map[string]interface{}, error) {
	// 找到 var input = map[string]interface{}{ ... } 部分
	start := strings.Index(code, "var input")
	if start == -1 {
		return nil, fmt.Errorf("未找到输入变量定义")
	}

	// 找到第二个 { (map[string]interface{ 后面的 {)
	firstBrace := strings.Index(code[start:], "{")
	if firstBrace == -1 {
		return nil, fmt.Errorf("输入变量定义格式错误")
	}
	firstBrace += start

	// 找到第二个 {
	braceStart := strings.Index(code[firstBrace+1:], "{")
	if braceStart == -1 {
		return nil, fmt.Errorf("输入变量定义格式错误")
	}
	braceStart += firstBrace + 1

	// 找到匹配的 }
	braceCount := 0
	braceEnd := -1
	for i := braceStart; i < len(code); i++ {
		if code[i] == '{' {
			braceCount++
		} else if code[i] == '}' {
			braceCount--
			if braceCount == 0 {
				braceEnd = i
				break
			}
		}
	}

	if braceEnd == -1 {
		return nil, fmt.Errorf("输入变量定义括号不匹配")
	}

	// 提取map内容
	mapContent := code[braceStart+1 : braceEnd]
	return p.parseMapContent(mapContent)
}

// 解析map内容
func (p *SimpleParser) parseMapContent(content string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 按行解析
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "," || line == "}" {
			continue
		}

		// 移除末尾的逗号
		line = strings.TrimSuffix(line, ",")

		// 分割键值对
		colonIndex := strings.Index(line, ":")
		if colonIndex == -1 {
			continue
		}

		key := strings.Trim(strings.TrimSpace(line[:colonIndex]), "\"")
		value := strings.TrimSpace(line[colonIndex+1:])

		// 解析值
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			// 字符串
			result[key] = strings.Trim(value, "\"")
		} else if value == "true" {
			result[key] = true
		} else if value == "false" {
			result[key] = false
		} else if num, err := strconv.Atoi(value); err == nil {
			result[key] = num
		} else {
			result[key] = value
		}
	}

	return result, nil
}

// 解析步骤定义（公开方法）
func (p *SimpleParser) ParseStep(line string) (*SimpleStep, error) {
	return p.parseStep(line, 0, nil)
}

// 解析步骤定义
func (p *SimpleParser) parseStep(line string, lineNumber int, lines []string) (*SimpleStep, error) {
	step := &SimpleStep{}

	// 根据 = 分隔
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return step, fmt.Errorf("步骤定义格式错误: %s", line)
	}

	step.Name = strings.TrimSpace(parts[0])

	// 根据 -> 分隔输入和输出
	arrowParts := strings.SplitN(parts[1], "->", 2)
	if len(arrowParts) != 2 {
		return step, fmt.Errorf("步骤定义缺少 -> 分隔符: %s", line)
	}

	inputPart := strings.TrimSpace(arrowParts[0])
	outputPart := strings.TrimSpace(arrowParts[1])

	// 解析输入部分
	inputTypes, function, isStatic, caseID, err := p.parseInputPart(inputPart)
	if err != nil {
		return step, err
	}

	step.Function = function
	step.InputParams = inputTypes
	step.IsStatic = isStatic
	step.CaseID = caseID

	// 解析输出部分和元数据
	outputTypes, metadata, err := p.parseOutputPartWithMetadata(outputPart)
	if err != nil {
		return step, err
	}

	step.OutputParams = outputTypes
	step.Metadata = metadata
	step.Logs = make([]*StepLog, 0)

	// 提取描述信息
	if lines != nil && lineNumber > 0 {
		step.Desc = p.extractDescription(lines, lineNumber-1)
	}

	return step, nil
}

// 解析输入部分
func (p *SimpleParser) parseInputPart(inputPart string) ([]ParameterInfo, string, bool, string, error) {
	// 检查是否为静态工作流 [用例ID]
	// 静态工作流的特征是：函数名[用例ID]，且用例ID不包含复杂类型字符
	if strings.Contains(inputPart, "[") && strings.Contains(inputPart, "]") {
		// 使用更精确的正则表达式匹配静态工作流格式
		// 静态工作流格式：functionName[caseID]，其中caseID通常是简单的字符串或数字
		re := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_.]*)\[([^\]]+)\]$`)
		matches := re.FindStringSubmatch(inputPart)
		if len(matches) == 3 {
			function := strings.TrimSpace(matches[1])
			caseID := strings.TrimSpace(matches[2])
			return []ParameterInfo{}, function, true, caseID, nil
		}
		// 如果不匹配静态工作流格式，继续处理为动态工作流
	}

	// 动态工作流 function(param: type "desc", param: type "desc")
	if strings.Contains(inputPart, "(") && strings.Contains(inputPart, ")") {
		// 提取函数名
		parenIndex := strings.Index(inputPart, "(")
		function := strings.TrimSpace(inputPart[:parenIndex])

		// 提取参数部分
		paramPart := inputPart[parenIndex+1 : strings.LastIndex(inputPart, ")")]

		// 解析参数定义
		inputParams := p.parseParameterDefinitions(paramPart)

		return inputParams, function, false, "", nil
	}

	return nil, "", false, "", fmt.Errorf("无法解析输入部分: %s", inputPart)
}

// 解析输出部分和元数据
func (p *SimpleParser) parseOutputPartWithMetadata(outputPart string) ([]ParameterInfo, map[string]interface{}, error) {
	// 移除分号
	outputPart = strings.TrimSuffix(outputPart, ";")

	// 检查是否有元数据 {key: value, key: value}
	// 需要区分真正的元数据和复杂类型中的括号
	metadata := make(map[string]interface{})
	if strings.Contains(outputPart, "{") && strings.Contains(outputPart, "}") {
		// 检查是否是真正的元数据格式：{key: value} 在字符串末尾
		// 元数据应该在字符串末尾，且不包含在复杂类型中
		lastBraceIndex := strings.LastIndex(outputPart, "}")
		firstBraceIndex := strings.LastIndex(outputPart, "{")

		if lastBraceIndex > firstBraceIndex && firstBraceIndex > 0 {
			// 检查是否是真正的元数据：大括号前有空格，且大括号在字符串末尾
			// 元数据格式：参数定义 {元数据}
			if firstBraceIndex > 0 && outputPart[firstBraceIndex-1] == ' ' &&
				lastBraceIndex == len(outputPart)-1 {
				// 提取元数据部分
				metadataStr := outputPart[firstBraceIndex : lastBraceIndex+1]
				metadata = p.parseMetadata(metadataStr)

				// 移除元数据部分，保留输出参数部分
				outputPart = strings.TrimSpace(outputPart[:firstBraceIndex])
			}
		}
	}

	// 检查是否有括号
	if strings.HasPrefix(outputPart, "(") && strings.HasSuffix(outputPart, ")") {
		outputPart = outputPart[1 : len(outputPart)-1]
	}

	// 解析输出参数
	var outputParams []ParameterInfo
	// 检查是否为新格式（包含冒号和引号）
	if strings.Contains(outputPart, ":") && strings.Contains(outputPart, "\"") {
		outputParams = p.parseParameterDefinitions(outputPart)
	} else {
		// 旧格式：bool 验证结果, err 是否失败
		outputParams = p.parseLegacyOutputPart(outputPart)
	}

	return outputParams, metadata, nil
}

// 解析输出部分（保持向后兼容）
func (p *SimpleParser) parseOutputPart(outputPart string) ([]ParameterInfo, error) {
	params, _, err := p.parseOutputPartWithMetadata(outputPart)
	return params, err
}

// 解析旧格式输出部分
func (p *SimpleParser) parseLegacyOutputPart(outputPart string) []ParameterInfo {
	params := make([]ParameterInfo, 0)

	// 按逗号分割参数
	paramList := strings.Split(outputPart, ",")

	for _, param := range paramList {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}

		// 解析格式：bool 验证结果
		fields := strings.Fields(param)
		if len(fields) >= 2 {
			paramType := fields[0]
			paramName := strings.Join(fields[1:], " ")

			params = append(params, ParameterInfo{
				Name: paramName,
				Type: paramType,
				Desc: paramName, // 使用参数名作为描述
			})
		}
	}

	return params
}

// 解析类型列表
func (p *SimpleParser) parseTypeList(typeList string) ([]SimpleTypeDef, error) {
	var result []SimpleTypeDef

	// 按逗号分割
	parts := strings.Split(typeList, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 分割类型和名称
		fields := strings.Fields(part)
		if len(fields) < 2 {
			return nil, fmt.Errorf("类型定义格式错误: %s", part)
		}

		typeDef := SimpleTypeDef{
			Type: fields[0],
			Name: strings.Join(fields[1:], " "),
		}

		result = append(result, typeDef)
	}

	return result, nil
}

// 解析main函数
func (p *SimpleParser) parseMainFunction(code string, lines []string, result *SimpleParseResult) (*SimpleMainFunc, error) {
	mainFunc := &SimpleMainFunc{Statements: []*SimpleStatement{}}

	// 找到main函数开始位置
	mainStart := -1
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "func main()") {
			mainStart = i
			break
		}
	}

	if mainStart == -1 {
		return mainFunc, nil
	}

	// 找到main函数结束位置（匹配的右括号）
	braceCount := 0
	mainEnd := -1
	for i := mainStart; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// 计算大括号
		for _, char := range line {
			if char == '{' {
				braceCount++
			} else if char == '}' {
				braceCount--
				if braceCount == 0 {
					mainEnd = i
					break
				}
			}
		}

		if mainEnd != -1 {
			break
		}
	}

	if mainEnd == -1 {
		mainEnd = len(lines) - 1
	}

	// 解析main函数内的语句（支持嵌套）
	statements, _ := p.parseStatements(lines, mainStart+1, mainEnd, result)
	mainFunc.Statements = statements

	return mainFunc, nil
}

// 提取描述信息
func (p *SimpleParser) extractDescription(lines []string, currentIndex int) string {
	// 检查当前行之前是否有 //desc: 注释
	for j := currentIndex - 1; j >= 0; j-- {
		prevLine := strings.TrimSpace(lines[j])
		if prevLine == "" {
			continue
		}
		if strings.HasPrefix(prevLine, "//desc:") {
			return strings.TrimSpace(prevLine[7:]) // 去掉 "//desc:" 前缀
		}
		if !strings.HasPrefix(prevLine, "//") {
			break // 遇到非注释行，停止查找
		}
	}
	return ""
}

// 解析语句列表（支持嵌套）
func (p *SimpleParser) parseStatements(lines []string, start, end int, result *SimpleParseResult) ([]*SimpleStatement, int) {
	var statements []*SimpleStatement

	for i := start; i < end; i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || line == "{" || line == "}" {
			continue
		}

		// 跳过纯注释行
		if strings.HasPrefix(line, "//") {
			continue
		}

		// 检查是否是if语句
		if strings.HasPrefix(line, "if ") {
			ifStmt, nextIndex := p.parseIfStatement(lines, i, result)
			statements = append(statements, ifStmt)
			i = nextIndex - 1 // -1 因为循环会+1
			continue
		}

		// 解析普通语句
		statement := p.parseStatement(line, i+1, result)
		if statement != nil {
			// 提取描述信息
			statement.Desc = p.extractDescription(lines, i)
			statements = append(statements, statement)
		}
	}

	return statements, end
}

// 解析if语句
func (p *SimpleParser) parseIfStatement(lines []string, start int, result *SimpleParseResult) (*SimpleStatement, int) {
	// 简化为单分支if语句，多分支通过递归解析处理
	line := strings.TrimSpace(lines[start])

	// 提取条件
	condition := ""
	if strings.HasPrefix(line, "if ") {
		condition = strings.TrimSpace(line[3:])
		// 移除末尾的 {
		condition = strings.TrimSuffix(condition, "{")
		condition = strings.TrimSpace(condition)
	}

	// 找到if语句的结束位置
	braceCount := 0
	ifEnd := -1
	for i := start; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// 计算大括号
		for _, char := range line {
			if char == '{' {
				braceCount++
			} else if char == '}' {
				braceCount--
				if braceCount == 0 {
					ifEnd = i
					break
				}
			}
		}

		if ifEnd != -1 {
			break
		}
	}

	if ifEnd == -1 {
		ifEnd = len(lines) - 1
	}

	// 解析if语句体内的语句
	children, _ := p.parseStatements(lines, start+1, ifEnd, result)

	// 提取描述信息
	desc := p.extractDescription(lines, start)

	return &SimpleStatement{
		Type:       "if",
		Content:    line,
		LineNumber: start + 1,
		Condition:  condition,
		Children:   children,
		Status:     StatusPending,
		RetryCount: 0,
		Desc:       desc,
	}, ifEnd + 1
}

// 解析语句
func (p *SimpleParser) parseStatement(line string, lineNumber int, result *SimpleParseResult) *SimpleStatement {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	// 移除末尾的分号
	line = strings.TrimSuffix(line, ";")

	// 判断语句类型
	// 注意：已移除打印语句支持，执行引擎会自动处理日志记录

	// 注意：已移除所有打印语句支持，执行引擎会自动处理日志记录

	if strings.HasPrefix(line, "return") {
		return &SimpleStatement{
			Type:       "return",
			Content:    line,
			LineNumber: lineNumber,
			Status:     StatusPending,
			RetryCount: 0,
		}
	}

	if strings.Contains(line, " := ") && !strings.Contains(line, "(") {
		// 变量赋值 - 解析为var
		parts := strings.SplitN(line, " := ", 2)
		if len(parts) == 2 {
			varName := strings.TrimSpace(parts[0])
			_ = strings.TrimSpace(parts[1]) // 暂时不使用，后续可以用于类型推断

			// 建立变量映射
			result.Variables[varName] = VariableInfo{
				Name:    varName,
				Type:    "string", // 默认为string，可以根据值推断
				Source:  "assignment",
				LineNum: lineNumber,
				IsInput: false,
			}

			return &SimpleStatement{
				Type:       "var",
				Content:    line,
				LineNumber: lineNumber,
				Status:     StatusPending,
				RetryCount: 0,
			}
		}
	}

	if strings.Contains(line, " := ") && strings.Contains(line, "(") {
		// 函数调用赋值 - 解析为function-call
		stmt := &SimpleStatement{
			Type:       "function-call",
			Content:    line,
			LineNumber: lineNumber,
			Metadata:   make(map[string]interface{}),
			Status:     StatusPending,
			RetryCount: 0,
		}

		// 解析函数名和参数
		parts := strings.SplitN(line, " := ", 2)
		if len(parts) == 2 {
			funcCall := strings.TrimSpace(parts[1])
			if strings.Contains(funcCall, "(") && strings.Contains(funcCall, ")") {
				// 检查是否有元数据
				if strings.Contains(funcCall, "){") && strings.Contains(funcCall, "}") {
					// 分离函数调用和元数据
					braceIndex := strings.Index(funcCall, "){")
					funcPart := funcCall[:braceIndex+1] // 包含右括号
					metadataPart := funcCall[braceIndex+1:]

					// 解析元数据
					stmt.Metadata = p.parseMetadata(metadataPart)
					funcCall = funcPart
				}

				// 提取函数名
				funcStart := strings.Index(funcCall, "(")
				funcName := strings.TrimSpace(funcCall[:funcStart])
				stmt.Function = funcName

				// 提取参数
				paramStart := funcStart + 1
				paramEnd := strings.LastIndex(funcCall, ")")
				if paramEnd > paramStart {
					paramStr := strings.TrimSpace(funcCall[paramStart:paramEnd])
					if paramStr != "" {
						stmt.Args = p.parseArguments(paramStr, funcName, result)
					}
				}

				// 解析返回变量并建立映射
				stmt.Returns = p.parseReturnVariables(parts[0], funcName, lineNumber, result)
			}
		}

		return stmt
	}

	if strings.Contains(line, " = ") && strings.Contains(line, "(") {
		// 函数调用赋值 - 解析为function-call
		stmt := &SimpleStatement{
			Type:       "function-call",
			Content:    line,
			LineNumber: lineNumber,
			Metadata:   make(map[string]interface{}),
			Status:     StatusPending,
			RetryCount: 0,
		}

		// 解析函数名和参数
		parts := strings.SplitN(line, " = ", 2)
		if len(parts) == 2 {
			funcCall := strings.TrimSpace(parts[1])
			if strings.Contains(funcCall, "(") && strings.Contains(funcCall, ")") {
				// 检查是否有元数据
				if strings.Contains(funcCall, "){") && strings.Contains(funcCall, "}") {
					// 分离函数调用和元数据
					braceIndex := strings.Index(funcCall, "){")
					funcPart := funcCall[:braceIndex+1] // 包含右括号
					metadataPart := funcCall[braceIndex+1:]

					// 解析元数据
					stmt.Metadata = p.parseMetadata(metadataPart)
					funcCall = funcPart
				}

				// 提取函数名
				funcStart := strings.Index(funcCall, "(")
				funcName := strings.TrimSpace(funcCall[:funcStart])
				stmt.Function = funcName

				// 提取参数
				paramStart := funcStart + 1
				paramEnd := strings.LastIndex(funcCall, ")")
				if paramEnd > paramStart {
					paramStr := strings.TrimSpace(funcCall[paramStart:paramEnd])
					if paramStr != "" {
						stmt.Args = p.parseArguments(paramStr, funcName, result)
					}
				}
			}
		}

		return stmt
	}

	if strings.Contains(line, " = ") {
		// 变量赋值
		return &SimpleStatement{
			Type:       "assign",
			Content:    line,
			LineNumber: lineNumber,
			Status:     StatusPending,
			RetryCount: 0,
		}
	}

	if strings.Contains(line, "(") && strings.Contains(line, ")") {
		// 纯函数调用（无赋值）
		stmt := &SimpleStatement{
			Type:       "function-call",
			Content:    line,
			LineNumber: lineNumber,
			Metadata:   make(map[string]interface{}),
			Status:     StatusPending,
			RetryCount: 0,
		}

		// 检查是否有元数据
		funcCall := line
		if strings.Contains(line, "){") && strings.Contains(line, "}") {
			// 分离函数调用和元数据
			braceIndex := strings.Index(line, "){")
			funcPart := line[:braceIndex+1] // 包含右括号
			metadataPart := line[braceIndex+1:]

			// 解析元数据
			stmt.Metadata = p.parseMetadata(metadataPart)
			funcCall = funcPart
		}

		// 解析函数名和参数
		funcStart := strings.Index(funcCall, "(")
		funcName := strings.TrimSpace(funcCall[:funcStart])
		stmt.Function = funcName

		// 提取参数
		paramStart := funcStart + 1
		paramEnd := strings.LastIndex(funcCall, ")")
		if paramEnd > paramStart {
			paramStr := strings.TrimSpace(funcCall[paramStart:paramEnd])
			if paramStr != "" {
				stmt.Args = p.parseArguments(paramStr, funcName, result)
			}
		}

		return stmt
	}

	// 其他语句
	return &SimpleStatement{
		Type:       "other",
		Content:    line,
		LineNumber: lineNumber,
		Status:     StatusPending,
		RetryCount: 0,
		Desc:       "",
	}
}

// 打印解析结果
func (r *SimpleParseResult) Print() {
	fmt.Println("=== 简单解析结果 ===")
	if !r.Success {
		fmt.Printf("❌ 解析失败: %s\n", r.Error)
		return
	}

	fmt.Println("✅ 解析成功")
	fmt.Printf("输入变量数量: %d\n", len(r.InputVars))
	fmt.Printf("工作流步骤数量: %d\n", len(r.Steps))
	if r.MainFunc != nil {
		fmt.Printf("主函数语句数量: %d\n", len(r.MainFunc.Statements))
	}

	if len(r.InputVars) > 0 {
		fmt.Println("\n🔧 输入变量:")
		for key, value := range r.InputVars {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	if len(r.Steps) > 0 {
		fmt.Println("\n⚙️ 工作流步骤:")
		for i, step := range r.Steps {
			fmt.Printf("  %d. %s\n", i+1, step.Name)
			fmt.Printf("     函数: %s\n", step.Function)
			if step.Desc != "" {
				fmt.Printf("     描述: %s\n", step.Desc)
			}
			if step.IsStatic {
				fmt.Printf("     类型: 静态工作流\n")
				fmt.Printf("     用例ID: %s\n", step.CaseID)
			} else {
				fmt.Printf("     类型: 动态工作流\n")
			}

			if len(step.InputParams) > 0 {
				fmt.Printf("     输入参数: ")
				for j, input := range step.InputParams {
					if j > 0 {
						fmt.Printf(", ")
					}
					fmt.Printf("%s %s (%s)", input.Type, input.Name, input.Desc)
				}
				fmt.Println()
			}

			if len(step.OutputParams) > 0 {
				fmt.Printf("     输出参数: ")
				for j, output := range step.OutputParams {
					if j > 0 {
						fmt.Printf(", ")
					}
					fmt.Printf("%s %s (%s)", output.Type, output.Name, output.Desc)
				}
				fmt.Println()
			}
		}
	}

	if r.MainFunc != nil && len(r.MainFunc.Statements) > 0 {
		fmt.Println("\n🎯 主函数语句:")
		r.printStatements(r.MainFunc.Statements, 0)
	}
}

// 递归打印语句（支持嵌套）
func (r *SimpleParseResult) printStatements(statements []*SimpleStatement, depth int) {
	indent := strings.Repeat("  ", depth)

	for i, stmt := range statements {
		// 打印语句信息
		fmt.Printf("%s%d. [%s] 第%d行: %s\n", indent, i+1, stmt.Type, stmt.LineNumber, stmt.Content)

		// 打印额外信息
		if stmt.Type == "function-call" && stmt.Function != "" {
			fmt.Printf("%s   函数: %s\n", indent, stmt.Function)
			if len(stmt.Args) > 0 {
				fmt.Printf("%s   输入参数:\n", indent)
				for j, arg := range stmt.Args {
					fmt.Printf("%s     %d. %s (类型: %s, 变量: %v, 字面量: %v, 输入: %v)\n",
						indent, j+1, arg.Value, arg.Type, arg.IsVariable, arg.IsLiteral, arg.IsInput)
					if arg.Source != "" && arg.Source != arg.Value {
						fmt.Printf("%s        来源: %s\n", indent, arg.Source)
					}
				}
			}
			if len(stmt.Returns) > 0 {
				fmt.Printf("%s   输出参数:\n", indent)
				for j, ret := range stmt.Returns {
					fmt.Printf("%s     %d. %s (类型: %s, 来源: %s)\n",
						indent, j+1, ret.Value, ret.Type, ret.Source)
				}
			}
			if len(stmt.Metadata) > 0 {
				fmt.Printf("%s   元数据:\n", indent)
				for key, value := range stmt.Metadata {
					fmt.Printf("%s     %s: %v\n", indent, key, value)
				}
			}
		}
		// 处理if语句
		if stmt.Type == "if" && stmt.Condition != "" {
			fmt.Printf("%s   条件: %s\n", indent, stmt.Condition)
		}

		// 递归打印子语句
		if len(stmt.Children) > 0 {
			fmt.Printf("%s   子语句:\n", indent)
			r.printStatements(stmt.Children, depth+2)
		}
	}
}

// 解析元数据配置
func (p *SimpleParser) parseMetadata(metadataStr string) map[string]interface{} {
	metadata := make(map[string]interface{})

	// 移除大括号
	metadataStr = strings.TrimSpace(metadataStr)
	if strings.HasPrefix(metadataStr, "{") && strings.HasSuffix(metadataStr, "}") {
		metadataStr = metadataStr[1 : len(metadataStr)-1]
	}

	// 按逗号分割键值对
	pairs := strings.Split(metadataStr, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		// 分割键值对
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 类型推断
		var parsedValue interface{}
		if value == "true" {
			parsedValue = true
		} else if value == "false" {
			parsedValue = false
		} else if num, err := strconv.Atoi(value); err == nil {
			parsedValue = num
		} else if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			// 字符串字面量
			parsedValue = value[1 : len(value)-1]
		} else {
			// 默认为字符串
			parsedValue = value
		}

		// 特殊处理 err_continue 参数
		if key == "err_continue" {
			if value == "true" {
				parsedValue = true
			} else {
				parsedValue = false
			}
		}

		metadata[key] = parsedValue
	}

	return metadata
}

// 解析参数定义（用于步骤定义中的参数）
func (p *SimpleParser) parseParameterDefinitions(paramStr string) []ParameterInfo {
	params := make([]ParameterInfo, 0)

	// 智能分割参数，考虑复杂类型中的逗号
	paramList := p.splitParameters(paramStr)

	for _, param := range paramList {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}

		// 解析格式：username: string "用户名"
		paramInfo := p.parseSingleParameterDefinition(param)
		if paramInfo != nil {
			params = append(params, *paramInfo)
		}
	}

	return params
}

// 智能分割参数，考虑复杂类型中的逗号
func (p *SimpleParser) splitParameters(paramStr string) []string {
	var params []string
	var current strings.Builder
	bracketCount := 0
	parenCount := 0
	inQuotes := false

	for _, char := range paramStr {
		switch char {
		case '"':
			inQuotes = !inQuotes
			current.WriteRune(char)
		case '{', '[':
			if !inQuotes {
				bracketCount++
			}
			current.WriteRune(char)
		case '}', ']':
			if !inQuotes {
				bracketCount--
			}
			current.WriteRune(char)
		case '(':
			if !inQuotes {
				parenCount++
			}
			current.WriteRune(char)
		case ')':
			if !inQuotes {
				parenCount--
			}
			current.WriteRune(char)
		case ',':
			if !inQuotes && bracketCount == 0 && parenCount == 0 {
				// 找到真正的参数分隔符
				param := strings.TrimSpace(current.String())
				if param != "" {
					params = append(params, param)
				}
				current.Reset()
			} else {
				// 在复杂类型内部，保留逗号
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	// 添加最后一个参数
	param := strings.TrimSpace(current.String())
	if param != "" {
		params = append(params, param)
	}

	return params
}

// 解析单个参数定义
func (p *SimpleParser) parseSingleParameterDefinition(paramStr string) *ParameterInfo {
	// 格式：username: string "用户名"
	// 支持复杂类型如 map[string]interface{}, []CustomStruct 等

	// 1. 查找描述部分（以引号包围的部分）
	descStart := strings.Index(paramStr, "\"")
	if descStart == -1 {
		// 没有描述，按冒号分割
		parts := strings.SplitN(paramStr, ":", 2)
		if len(parts) != 2 {
			return nil
		}
		name := strings.TrimSpace(parts[0])
		paramType := strings.TrimSpace(parts[1])
		return &ParameterInfo{
			Name: name,
			Type: paramType,
			Desc: "",
		}
	}

	// 2. 从描述开始位置往前找，找到最后一个冒号
	// 需要跳过复杂类型中的冒号，如 map[string]interface{} 中的冒号
	lastColonIndex := -1
	bracketCount := 0
	parenCount := 0

	for i := descStart - 1; i >= 0; i-- {
		char := paramStr[i]
		if char == '}' || char == ']' || char == ')' {
			if char == '}' {
				bracketCount++
			} else if char == ']' {
				bracketCount++
			} else if char == ')' {
				parenCount++
			}
		} else if char == '{' || char == '[' || char == '(' {
			if char == '{' {
				bracketCount--
			} else if char == '[' {
				bracketCount--
			} else if char == '(' {
				parenCount--
			}
		} else if char == ':' && bracketCount == 0 && parenCount == 0 {
			// 找到了不在括号内的冒号
			lastColonIndex = i
			break
		}
	}

	if lastColonIndex == -1 {
		return nil
	}

	name := strings.TrimSpace(paramStr[:lastColonIndex])
	typeAndDesc := strings.TrimSpace(paramStr[lastColonIndex+1:])

	// 3. 从typeAndDesc中分离类型和描述
	descStart = strings.Index(typeAndDesc, "\"")
	if descStart == -1 {
		// 没有描述，整个都是类型
		paramType := strings.TrimSpace(typeAndDesc)
		return &ParameterInfo{
			Name: name,
			Type: paramType,
			Desc: "",
		}
	}

	// 找到描述开始位置，类型是引号之前的部分
	paramType := strings.TrimSpace(typeAndDesc[:descStart])
	desc := strings.TrimSpace(typeAndDesc[descStart:])

	// 4. 去掉描述中的引号
	desc = strings.Trim(desc, "\"")

	return &ParameterInfo{
		Name: name,
		Type: paramType,
		Desc: desc,
	}
}

// 解析函数参数为ArgumentInfo结构体
func (p *SimpleParser) parseArguments(paramStr string, funcName string, result *SimpleParseResult) []*ArgumentInfo {
	var args []*ArgumentInfo

	// 从步骤定义中获取输入参数信息
	var inputParams []ParameterInfo
	for _, step := range result.Steps {
		if step.Name == funcName {
			inputParams = step.InputParams
			break
		}
	}

	// 按逗号分割参数
	params := strings.Split(paramStr, ",")
	for i, param := range params {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}

		arg := &ArgumentInfo{
			Value: param,
		}

		// 从输入参数定义中获取描述信息
		if i < len(inputParams) {
			arg.Desc = inputParams[i].Desc
		}

		// 判断参数类型
		if strings.HasPrefix(param, "input[") && strings.HasSuffix(param, "]") {
			// 输入参数：input["用户名"]
			arg.IsInput = true
			arg.IsVariable = true
			arg.Type = "input"
			arg.Source = "input"
		} else if strings.Contains(param, "\"") || strings.Contains(param, "'") {
			// 字符串字面量
			arg.IsLiteral = true
			arg.Type = "string"
		} else if _, err := strconv.Atoi(param); err == nil {
			// 数字字面量
			arg.IsLiteral = true
			arg.Type = "int"
		} else {
			// 变量引用
			arg.IsVariable = true
			arg.Type = "variable"
			arg.Source = param

			// 从变量映射表中获取类型信息
			if varInfo, exists := result.Variables[param]; exists {
				arg.Type = varInfo.Type
				arg.Source = varInfo.Source
			}
		}

		args = append(args, arg)
	}

	return args
}

// 解析返回变量并建立映射
func (p *SimpleParser) parseReturnVariables(varStr, funcName string, lineNumber int, result *SimpleParseResult) []*ArgumentInfo {
	// 从步骤定义中获取输出类型
	var outputParams []ParameterInfo
	for _, step := range result.Steps {
		// 匹配步骤名称（如step1）而不是函数名
		if step.Name == funcName {
			outputParams = step.OutputParams
			break
		}
	}

	var returns []*ArgumentInfo

	// 分割变量名
	vars := strings.Split(varStr, ",")
	for i, varName := range vars {
		varName = strings.TrimSpace(varName)
		if varName == "" {
			continue
		}

		// 从输出参数中获取对应的类型和描述
		varType := "unknown"
		varDesc := ""
		if i < len(outputParams) {
			varType = outputParams[i].Type
			varDesc = outputParams[i].Desc
		}

		// 检查变量名是否重复，如果重复则重命名
		originalName := varName
		if varName == "err" {
			varName = funcName + "Err"
		}

		// 创建返回参数信息
		returnArg := &ArgumentInfo{
			Value:      varName,
			Type:       varType,
			Desc:       varDesc,
			IsVariable: true,
			IsLiteral:  false,
			IsInput:    false,
			Source:     funcName,
			LineNum:    lineNumber,
		}
		returns = append(returns, returnArg)

		// 建立变量映射
		result.Variables[varName] = VariableInfo{
			Name:    varName,
			Type:    varType,
			Source:  funcName,
			LineNum: lineNumber,
			IsInput: false,
		}

		// 如果重命名了，记录原始名称
		if originalName != varName {
			result.Variables[originalName] = VariableInfo{
				Name:    varName,
				Type:    varType,
				Source:  funcName,
				LineNum: lineNumber,
				IsInput: false,
			}
		}
	}

	return returns
}

// 设置语句状态
func (s *SimpleStatement) SetStatus(status StatementStatus) {
	s.Status = status
}

// 增加重试次数
func (s *SimpleStatement) IncrementRetry() {
	s.RetryCount++
}

// 重置重试次数
func (s *SimpleStatement) ResetRetry() {
	s.RetryCount = 0
}

// 获取步骤名称（用于日志记录）
func (s *SimpleStatement) GetStepName() string {
	if s.Type == "function-call" && s.Function != "" {
		return s.Function
	}
	return "unknown"
}

// 开始执行计时
func (s *SimpleStatement) StartExecution() {
	now := time.Now()
	s.StartTime = &now
	s.Status = StatusRunning
}

// 结束执行计时
func (s *SimpleStatement) EndExecution() {
	now := time.Now()
	s.EndTime = &now
	if s.StartTime != nil {
		s.Duration = s.EndTime.Sub(*s.StartTime)
	}
}

// 获取合并后的元数据（默认值 + 语句元数据）
func (s *SimpleStatement) GetMergedMetadata() map[string]interface{} {
	defaultMeta := GetDefaultMetadata()
	merged := make(map[string]interface{})

	// 设置默认值
	if defaultMeta.Timeout != nil {
		merged["timeout"] = defaultMeta.Timeout.Milliseconds()
	} else {
		merged["timeout"] = nil // 无超时限制
	}
	merged["retry_count"] = defaultMeta.RetryCount
	merged["async"] = defaultMeta.Async
	merged["priority"] = defaultMeta.Priority
	merged["debug"] = defaultMeta.Debug
	merged["log_level"] = defaultMeta.LogLevel
	merged["ai_model"] = defaultMeta.AIModel

	// 覆盖语句中的元数据
	for key, value := range s.Metadata {
		merged[key] = value
	}

	return merged
}

// 获取超时时间，返回nil表示无超时限制
func (s *SimpleStatement) GetTimeout() *time.Duration {
	meta := s.GetMergedMetadata()
	if timeout, exists := meta["timeout"]; exists {
		if timeout == nil {
			return nil // 无超时限制
		}
		if timeoutMs, ok := timeout.(int); ok {
			duration := time.Duration(timeoutMs) * time.Millisecond
			return &duration
		}
	}
	return GetDefaultMetadata().Timeout
}

// 获取重试次数
func (s *SimpleStatement) GetRetryCount() int {
	meta := s.GetMergedMetadata()
	// 支持 retry 和 retry_count 两种键名
	if retry, exists := meta["retry"]; exists {
		if retryCount, ok := retry.(int); ok {
			return retryCount
		}
	}
	if retry, exists := meta["retry_count"]; exists {
		if retryCount, ok := retry.(int); ok {
			return retryCount
		}
	}
	return GetDefaultMetadata().RetryCount
}

// 是否异步执行
func (s *SimpleStatement) IsAsync() bool {
	meta := s.GetMergedMetadata()
	if async, exists := meta["async"]; exists {
		if asyncVal, ok := async.(bool); ok {
			return asyncVal
		}
	}
	return GetDefaultMetadata().Async
}

// 是否调试模式
func (s *SimpleStatement) IsDebug() bool {
	meta := s.GetMergedMetadata()
	if debug, exists := meta["debug"]; exists {
		if debugVal, ok := debug.(bool); ok {
			return debugVal
		}
	}
	return GetDefaultMetadata().Debug
}

// 获取优先级
func (s *SimpleStatement) GetPriority() int {
	meta := s.GetMergedMetadata()
	if priority, exists := meta["priority"]; exists {
		if priorityVal, ok := priority.(int); ok {
			return priorityVal
		}
		// 支持字符串优先级
		if priorityStr, ok := priority.(string); ok {
			switch priorityStr {
			case "high":
				return 1
			case "medium":
				return 0
			case "low":
				return -1
			}
		}
	}
	return GetDefaultMetadata().Priority
}

// 获取日志级别
func (s *SimpleStatement) GetLogLevel() string {
	meta := s.GetMergedMetadata()
	if logLevel, exists := meta["log_level"]; exists {
		if logLevelVal, ok := logLevel.(string); ok {
			return logLevelVal
		}
	}
	return GetDefaultMetadata().LogLevel
}

// 获取AI模型
func (s *SimpleStatement) GetAIModel() string {
	meta := s.GetMergedMetadata()
	if aiModel, exists := meta["ai_model"]; exists {
		if aiModelVal, ok := aiModel.(string); ok {
			return aiModelVal
		}
	}
	return GetDefaultMetadata().AIModel
}
