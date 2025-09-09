# AI工作流编排语言解析器

这是一个用Go语言实现的AI工作流编排语言解析器，能够解析基于Go语法的工作流代码并生成抽象语法树(AST)。

## 功能特性

- ✅ **简单解析器**: 基于行解析和符号分隔的高效解析器
- ✅ **静态工作流**: 支持用例ID引用的静态工作流
- ✅ **动态工作流**: 支持参数动态传递的工作流
- ✅ **模板系统**: 支持 `{{变量名}}` 和 `{{步骤名.字段名}}` 模板变量
- ✅ **条件判断**: 支持 if-else 条件分支和嵌套结构
- ✅ **变量映射**: 自动建立变量到来源函数的映射关系
- ✅ **变量防重复**: 自动重命名重复变量名（如err → step1Err）
- ✅ **变量赋值**: 支持 `var` 类型变量赋值解析
- ✅ **参数结构**: 使用 `ArgumentInfo` 结构体提供详细参数信息
- ✅ **注释忽略**: 自动忽略 `//` 开头的注释行
- ✅ **中文支持**: 完整支持中文变量名和注释

## 项目结构

```
workflow/
├── simple_parser.go      # 简单解析器核心实现
├── simple_parser_test.go # 完整测试套件
├── readme.md            # 设计文档
└── README.md            # 说明文档
```

## 快速开始

### 1. 运行测试

```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestSimpleParser_DynamicWorkflow
go test -v -run TestSimpleParser_StaticWorkflow
```

### 2. 在代码中使用

```go
package main

import (
    "fmt"
    "github.com/yunhanshu-net/pkg/workflow"
)

func main() {
    // 工作流代码
    code := `
var input = map[string]interface{}{
    "项目名称": "my-project",
}

step1 = beiluo.test1.devops.git_push[用例001] -> (err 是否失败);

func main() {
    fmt.Println("开始执行发布流程...")
    err := step1()
    if err != nil {
        fmt.Printf("推送代码失败: %v\n", err)
        return
    }
    fmt.Println("✅ 代码推送成功")
}
`

    // 创建解析器
    parser := workflow.NewSimpleParser()
    
    // 解析程序
    result := parser.ParseWorkflow(code)
    
    // 检查错误
    if !result.Success {
        fmt.Printf("解析失败: %s\n", result.Error)
        return
    }
    
    // 使用解析结果
    result.Print()
    
    // 访问变量映射表
    for varName, info := range result.Variables {
        fmt.Printf("变量: %s, 类型: %s, 来源: %s\n", varName, info.Type, info.Source)
    }
    
    // 访问主函数语句
    for i, stmt := range result.MainFunc.Statements {
        fmt.Printf("语句 %d: [%s] %s\n", i+1, stmt.Type, stmt.Content)
        if stmt.Type == "function-call" && len(stmt.Args) > 0 {
            fmt.Printf("  参数: ")
            for j, arg := range stmt.Args {
                if j > 0 {
                    fmt.Printf(", ")
                }
                fmt.Printf("%s(%s)", arg.Value, arg.Type)
            }
            fmt.Println()
        }
    }
}
```

## 语法支持

### 1. 静态工作流

```go
var (
    step1 = beiluo.test1.devops.git_push[用例001] -> (err 是否失败);
    step2 = beiluo.test1.devops.deploy_test[用例002] -> (int cost, err 是否失败);
)

func main() {
    err := step1()  // 参数从用例001获取
    err = step2()   // 参数从用例002获取
}
```

### 2. 动态工作流

```go
var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(string 用户名, int 手机号) -> (string 工号, string 用户名, err 是否失败);
step2 = beiluo.test1.crm.crm_interview_schedule(string 用户名) -> (string 面试时间, string 面试官名称, err 是否失败);

func main() {
    // 变量自动重命名：err → step1Err, step2Err
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    面试时间, 面试官名称, step2Err := step2(用户名)  // 使用step1的输出
}
```

### 3. 模板字符串

```go
// 支持 {{变量名}} 模板语法
通知信息 := `你收到了:{{用户名}},时间：{{面试时间}}的面试安排，请关注`
```

### 4. 变量赋值

```go
// 变量赋值解析为 var 类型
通知信息 := `你收到了:{{用户名}},时间：{{面试时间}}的面试安排，请关注`
状态 := "已完成"
```

### 5. 条件判断

```go
if step1Err != nil {
    fmt.Printf("创建用户失败: %v\n", step1Err)
    return
} else {
    面试时间, 面试官名称, step2Err := step2(用户名)
    if step2Err != nil {
        fmt.Printf("安排面试失败: %v\n", step2Err)
        return
    }
}
```

## 核心特性详解

### 变量防重复机制

解析器会自动检测并重命名重复的变量名，确保变量名唯一性：

```go
// 原始代码
工号, 用户名, err := step1(input["用户名"], input["手机号"])
面试时间, 面试官名称, err := step2(用户名)  // err重复了

// 解析器自动重命名为：
工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
面试时间, 面试官名称, step2Err := step2(用户名)
```

### 变量映射表

解析器会建立完整的变量映射关系：

```go
type VariableInfo struct {
    Name     string // 变量名
    Type     string // 变量类型
    Source   string // 来源函数名
    LineNum  int    // 定义行号
    IsInput  bool   // 是否来自input
}

// 示例映射
result.Variables = map[string]VariableInfo{
    "工号": {Name: "工号", Source: "step1", LineNum: 6, IsInput: false},
    "step1Err": {Name: "step1Err", Source: "step1", LineNum: 6, IsInput: false},
    "step2Err": {Name: "step2Err", Source: "step2", LineNum: 14, IsInput: false},
}
```

### 参数信息结构

解析器使用 `ArgumentInfo` 结构体提供详细的参数信息：

```go
type ArgumentInfo struct {
    Value      string `json:"value"`       // 参数值
    Type       string `json:"type"`        // 参数类型
    IsVariable bool   `json:"is_variable"` // 是否为变量引用
    IsLiteral  bool   `json:"is_literal"`  // 是否为字面量
    IsInput    bool   `json:"is_input"`    // 是否为输入参数
    Source     string `json:"source"`      // 来源（变量名或函数名）
    LineNum    int    `json:"line_num"`    // 定义行号
}

// 示例参数解析
// step1(input["用户名"], input["手机号"])
// 参数1: {Value: "input[\"用户名\"]", Type: "input", IsInput: true, Source: "input"}
// 参数2: {Value: "input[\"手机号\"]", Type: "input", IsInput: true, Source: "input"}
```

## API文档

### SimpleParser

```go
// 创建简单解析器
func NewSimpleParser() *SimpleParser

// 解析工作流代码
func (p *SimpleParser) ParseWorkflow(code string) *SimpleParseResult
```

### SimpleParseResult

```go
type SimpleParseResult struct {
    Success   bool                    // 解析是否成功
    InputVars map[string]interface{} // 输入变量
    Steps     []SimpleStep           // 工作流步骤
    MainFunc  *SimpleMainFunc        // 主函数
    Variables map[string]VariableInfo // 变量映射表
    Error     string                 // 错误信息
}
```

### 核心数据结构

- `SimpleStep`: 工作流步骤定义
- `SimpleMainFunc`: 主函数
- `SimpleStatement`: 语句（支持嵌套）
- `VariableInfo`: 变量信息
- `ArgumentInfo`: 参数详细信息

## 测试

```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestSimpleParser_DynamicWorkflow
go test -v -run TestSimpleParser_StaticWorkflow
go test -v -run TestSimpleParser_MixedWorkflow
```

## 性能特点

- **简单高效**: 基于行解析和符号分隔，性能优异
- **内存友好**: 低内存占用，适合大规模工作流
- **解析快速**: 直接字符串操作，无复杂词法分析开销
- **易于维护**: 代码简洁，逻辑清晰

## 设计理念

### 简单就是美

这个解析器遵循"简单就是美"的设计理念：

1. **按行解析**: 直接按行读取，避免复杂的词法分析
2. **符号分隔**: 使用 `=`, `->`, `:=` 等符号进行分隔
3. **字符串操作**: 基于字符串操作，简单直观
4. **功能完整**: 在简单的基础上实现所有必需功能

### 核心优势

- ✅ **零学习成本**: 基于Go语法，开发者无需学习新语法
- ✅ **中文友好**: 完整支持中文变量名和注释
- ✅ **智能重命名**: 自动处理变量名冲突
- ✅ **变量追踪**: 完整的变量来源映射
- ✅ **模板支持**: 支持 `{{变量名}}` 模板语法
- ✅ **参数分析**: 详细的参数类型和来源分析
- ✅ **变量赋值**: 智能识别变量赋值语句

## 贡献

欢迎提交Issue和Pull Request来改进这个解析器！

## 许可证

MIT License