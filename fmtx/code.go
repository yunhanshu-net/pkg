package fmtx

import (
	"fmt"
	"regexp"
	"strings"
)

// DeleteVar 删除代码中重复定义的变量
// variable 指的是变量名称，code是代码内容
// 主要用于删除框架自动定义的变量（如RouterGroup）的重复定义
func DeleteVar(variable string, code string) string {
	if variable == "" || code == "" {
		return code
	}

	// 匹配变量定义的多种模式
	patterns := []string{
		// var VariableName = ... (支持换行)
		`var\s+` + regexp.QuoteMeta(variable) + `\s*=[^;\n]*;?\s*\n`,
		// var VariableName\n    = ... (换行情况)
		`var\s+` + regexp.QuoteMeta(variable) + `\s*\n\s*=\s*[^;\n]*;?\s*\n`,
		// VariableName := ... (支持换行)
		regexp.QuoteMeta(variable) + `\s*:=\s*[^;\n]*;?\s*\n`,
		// VariableName\n    := ... (换行情况)
		regexp.QuoteMeta(variable) + `\s*\n\s*:=\s*[^;\n]*;?\s*\n`,
		// const VariableName = ... (支持换行)
		`const\s+` + regexp.QuoteMeta(variable) + `\s*=[^;\n]*;?\s*\n`,
		// const VariableName\n    = ... (换行情况)
		`const\s+` + regexp.QuoteMeta(variable) + `\s*\n\s*=\s*[^;\n]*;?\s*\n`,
	}

	result := code
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		result = re.ReplaceAllString(result, "")
	}

	return result
}

// DecodeFileName 从生成的代码中解析出文件名称
// 支持两种解析方式：
// 1. 从函数组的EnName字段提取
// 2. 从路由注册中提取函数名
func DecodeFileName(code string) string {
	if code == "" {
		return ""
	}

	// 方式1：从文件名注释中提取（首选方案）
	// 匹配模式：文件名：xxx.go 或 文件：xxx.go
	commentPattern := `(?:文件名|文件)[：:]\s*([a-zA-Z0-9_]+\.go)`
	commentRe := regexp.MustCompile(commentPattern)
	commentMatches := commentRe.FindStringSubmatch(code)
	if len(commentMatches) > 1 {
		return commentMatches[1]
	}

	// 方式2：从函数组中提取EnName
	// 匹配模式：var XxxGroup = &runner.FunctionGroup{...EnName: "xxx"...} (支持多行)
	groupPattern := `var\s+\w+Group\s*=\s*&runner\.FunctionGroup\s*\{[\s\S]*?EnName:\s*"([^"]+)"`
	groupRe := regexp.MustCompile(groupPattern)
	groupMatches := groupRe.FindStringSubmatch(code)
	if len(groupMatches) > 1 {
		return groupMatches[1] + ".go"
	}

	// 方式3：从路由注册中提取函数名
	// 匹配模式：runner.Post/Get(RouterGroup+"/function_name", ...)
	routePattern := `runner\.(?:Post|Get)\(RouterGroup\+"/([^"]+)"`
	routeRe := regexp.MustCompile(routePattern)
	routeMatches := routeRe.FindStringSubmatch(code)
	if len(routeMatches) > 1 {
		// 提取第一个匹配的路由作为文件名
		routeName := routeMatches[1]
		// 如果路由名包含下划线，直接使用；否则可能需要进一步处理
		return routeName + ".go"
	}

	return ""
}

// ConvertToRouterGroup 将硬编码的路由路径转换为使用RouterGroup变量的标准格式
// 例如：runner.Post("/任意路径/function_name", ...) -> runner.Post(RouterGroup+"/function_name", ...)
func ConvertToRouterGroup(code string) string {
	if code == "" {
		return ""
	}

	// 匹配模式：runner.Post/Get("/任意路径", ...)
	// 提取路径作为函数名，替换为RouterGroup+"/路径"
	// 支持单路径段和多路径段
	pattern := `runner\.(Post|Get)\(["']/[^"']+["']`
	re := regexp.MustCompile(pattern)

	result := re.ReplaceAllStringFunc(code, func(match string) string {
		// 手动解析匹配的字符串来提取HTTP方法和函数名
		// match 格式：runner.Post("/path/to/function_name"
		// 或者：runner.Get("/path/to/function_name"

		// 提取HTTP方法
		methodMatch := regexp.MustCompile(`runner\.(Post|Get)\(`)
		methodMatches := methodMatch.FindStringSubmatch(match)
		if len(methodMatches) < 2 {
			return match
		}
		method := methodMatches[1]

		// 提取路径部分
		pathMatch := regexp.MustCompile(`["']([^"']+)["']`)
		pathMatches := pathMatch.FindStringSubmatch(match)
		if len(pathMatches) < 2 {
			return match
		}
		path := pathMatches[1]

		// 提取路径作为函数名（去掉开头的斜杠）
		// 对于单路径段：/single_path -> single_path
		// 对于多路径段：/path/to/function -> function（取最后一个路径段）
		var functionName string
		if strings.Count(path, "/") == 1 {
			// 单路径段：/single_path -> single_path
			functionName = path[1:] // 去掉开头的斜杠
		} else {
			// 多路径段：取最后一个路径段
			lastSlashIndex := strings.LastIndex(path, "/")
			if lastSlashIndex == -1 {
				return match
			}
			functionName = path[lastSlashIndex+1:]
		}

		return fmt.Sprintf("runner.%s(RouterGroup+\"/%s\"", method, functionName)
	})

	return result
}
