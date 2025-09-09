package main

import (
	"fmt"
	"strings"
)

// 参数定义信息
type ParameterInfo struct {
	Name string `json:"name"` // 英文参数名
	Type string `json:"type"` // 参数类型
	Desc string `json:"desc"` // 中文描述
}

// 解析单个参数定义
func parseSingleParameterDefinition(paramStr string) *ParameterInfo {
	// 格式：username: string "用户名"
	// 1. 按冒号分割
	parts := strings.SplitN(paramStr, ":", 2)
	if len(parts) != 2 {
		return nil
	}

	name := strings.TrimSpace(parts[0])
	typeAndDesc := strings.TrimSpace(parts[1])

	// 2. 按空格分割类型和描述
	typeDescParts := strings.SplitN(typeAndDesc, " ", 2)
	if len(typeDescParts) != 2 {
		return nil
	}

	paramType := strings.TrimSpace(typeDescParts[0])
	desc := strings.TrimSpace(typeDescParts[1])

	// 3. 去掉描述中的引号
	desc = strings.Trim(desc, "\"")

	return &ParameterInfo{
		Name: name,
		Type: paramType,
		Desc: desc,
	}
}

// 解析参数定义列表
func parseParameterDefinitions(paramStr string) []ParameterInfo {
	params := make([]ParameterInfo, 0)

	// 按逗号分割参数
	paramList := strings.Split(paramStr, ",")

	for _, param := range paramList {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}

		// 解析格式：username: string "用户名"
		paramInfo := parseSingleParameterDefinition(param)
		if paramInfo != nil {
			params = append(params, *paramInfo)
		}
	}

	return params
}

func main() {
	fmt.Println("=== 测试新的参数定义格式 ===")

	// 测试单个参数定义
	testCases := []string{
		`username: string "用户名"`,
		`phone: int "手机号"`,
		`email: string "邮箱"`,
		`workId: string "工号"`,
		`err: error "是否失败"`,
	}

	fmt.Println("\n--- 单个参数定义测试 ---")
	for _, testCase := range testCases {
		fmt.Printf("输入: %s\n", testCase)
		result := parseSingleParameterDefinition(testCase)
		if result != nil {
			fmt.Printf("结果: Name=%s, Type=%s, Desc=%s\n", result.Name, result.Type, result.Desc)
		} else {
			fmt.Printf("解析失败\n")
		}
		fmt.Println()
	}

	// 测试参数列表
	fmt.Println("--- 参数列表测试 ---")
	paramListStr := `username: string "用户名", phone: int "手机号", email: string "邮箱"`
	fmt.Printf("输入: %s\n", paramListStr)

	params := parseParameterDefinitions(paramListStr)
	fmt.Printf("解析结果: %d 个参数\n", len(params))
	for i, param := range params {
		fmt.Printf("  %d. Name=%s, Type=%s, Desc=%s\n", i+1, param.Name, param.Type, param.Desc)
	}

	// 测试复杂的工作流步骤定义
	fmt.Println("\n--- 工作流步骤定义测试 ---")
	stepDefinition := `beiluo.test1.devops.devops_script_create(
    username: string "用户名",
    phone: int "手机号", 
    email: string "邮箱"
) -> (
    workId: string "工号",
    username: string "用户名", 
    err: error "是否失败"
)`

	fmt.Printf("步骤定义: %s\n", stepDefinition)

	// 解析输入参数部分
	inputPart := `username: string "用户名",
    phone: int "手机号", 
    email: string "邮箱"`

	inputParams := parseParameterDefinitions(inputPart)
	fmt.Printf("输入参数: %d 个\n", len(inputParams))
	for i, param := range inputParams {
		fmt.Printf("  %d. %s (%s) - %s\n", i+1, param.Name, param.Type, param.Desc)
	}

	// 解析输出参数部分
	outputPart := `workId: string "工号",
    username: string "用户名", 
    err: error "是否失败"`

	outputParams := parseParameterDefinitions(outputPart)
	fmt.Printf("输出参数: %d 个\n", len(outputParams))
	for i, param := range outputParams {
		fmt.Printf("  %d. %s (%s) - %s\n", i+1, param.Name, param.Type, param.Desc)
	}
}
