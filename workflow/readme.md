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
- ✅ **元数据支持**: 支持函数调用元数据配置 `{retry:1, timeout:2000}`
- ✅ **注释忽略**: 自动忽略 `//` 开头的注释行
- ✅ **中文支持**: 完整支持中文变量名和注释
- ✅ **状态管理**: 支持语句执行状态跟踪（pending/running/completed/failed/skipped）
- ✅ **重试机制**: 支持语句重试次数统计和管理
- ✅ **步骤日志**: 支持步骤级别日志记录 `step1.Printf("消息")`
- ✅ **步骤描述**: 支持 `//desc:` 注释为每个步骤添加详细说明
- ✅ **执行引擎**: 完整的工作流执行引擎，支持回调机制
- ✅ **状态持久化**: 支持SQLite持久化，实现断点续传
- ✅ **并发执行**: 支持多工作流实例并发执行
- ✅ **Context支持**: 所有回调函数支持context.Context，支持取消和超时

## 项目结构

```
workflow/
├── simple_parser.go          # 简单解析器核心实现
├── workflow_executor.go      # 工作流执行引擎
├── persistence_test.go       # 状态持久化测试
├── parameter_mapping_test.go # 参数映射测试
├── executor_test.go          # 执行器测试
├── simple_parser_test.go     # 完整测试套件
├── examples/                 # 示例项目
│   ├── user_registration/    # 用户注册工作流示例
│   ├── error_handling/       # 错误处理示例
│   └── concurrent_execution/ # 并发执行示例
├── readme.md                 # 设计文档
├── 用户使用指南.md            # 用户使用指南
└── 开发者集成指南.md          # 开发者集成指南
```

## 文档导航

- **[用户使用指南](用户使用指南.md)** - 面向工作流编写者，介绍语言语法和特性，包含丰富的完整示例
- **[开发者集成指南](开发者集成指南.md)** - 面向服务端开发者，介绍解析器API和集成方法
- **[设计文档](readme.md)** - 技术设计文档，介绍解析器架构和实现细节

## 快速开始

### 1. 运行测试

```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestSimpleParser_DynamicWorkflow

# 运行持久化测试
go test -v -run TestSQLitePersistence
```

### 2. 运行示例项目

```bash
# 用户注册工作流示例
cd examples/user_registration
go run main.go

# 错误处理示例
cd examples/error_handling
go run main.go

# 并发执行示例
cd examples/concurrent_execution
go run main.go
```

### 3. 在代码中使用

```go
package main

import (
    "context"
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
    //desc: 开始执行发布流程
    fmt.Println("开始执行发布流程...")
    
    //desc: 推送代码到远程仓库
    err := step1()
    
    //desc: 检查代码推送是否成功
    if err != nil {
        //desc: 推送失败，记录错误并退出
        step1.Printf("推送代码失败: %v", err)
        return
    }
    
    //desc: 推送成功，记录成功日志
    step1.Printf("✅ 代码推送成功")
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
    
    // 创建执行器
    executor := workflow.NewExecutor()
    
    // 设置回调函数
    executor.OnFunctionCall = func(ctx context.Context, step workflow.SimpleStep, in *workflow.ExecutorIn) (*workflow.ExecutorOut, error) {
        fmt.Printf("执行步骤: %s - %s\n", step.Name, in.StepDesc)
        fmt.Printf("输入参数: %+v\n", in.RealInput)
        
        // 模拟业务逻辑
        return &workflow.ExecutorOut{
            Success: true,
            WantOutput: map[string]interface{}{
                "err": nil,
            },
            Error: "",
            Logs: []string{"步骤执行成功"},
        }, nil
    }
    
    executor.OnWorkFlowUpdate = func(ctx context.Context, current *workflow.SimpleParseResult) error {
        fmt.Printf("工作流状态更新: FlowID=%s\n", current.FlowID)
        return nil
    }
    
    executor.OnWorkFlowExit = func(ctx context.Context, current *workflow.SimpleParseResult) error {
        fmt.Println("工作流正常结束")
        return nil
    }
    
    // 执行工作流
    ctx := context.Background()
    if err := executor.Start(ctx, result); err != nil {
        fmt.Printf("执行失败: %v\n", err)
        return
    }
    
    fmt.Println("工作流执行完成！")
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
    //desc: 推送代码到远程仓库
    err := step1()  // 参数从用例001获取
    
    //desc: 检查代码推送是否成功
    if err != nil {
        //desc: 推送失败，记录错误并退出
        step1.Printf("推送代码失败: %v", err)
        return
    }
    
    //desc: 推送成功，记录成功日志
    step1.Printf("✅ 代码推送成功")
    
    //desc: 部署到测试环境
    err = step2()   // 参数从用例002获取
    
    //desc: 检查测试环境部署是否成功
    if err != nil {
        //desc: 部署失败，记录错误并退出
        step2.Printf("发布测试环境失败: %v", err)
        return
    }
    
    //desc: 部署成功，记录成功日志
    step2.Printf("✅ 测试环境发布成功")
}
```

### 2. 动态工作流

```go
var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(
    username: string "用户名",
    phone: int "手机号"
) -> (
    workId: string "工号",
    username: string "用户名", 
    err: error "是否失败"
);

step2 = beiluo.test1.crm.crm_interview_schedule(
    username: string "用户名"
) -> (
    interviewTime: string "面试时间",
    interviewer: string "面试官名称", 
    err: error "是否失败"
);

func main() {
    //desc: 开始用户注册和面试安排流程
    fmt.Println("开始用户注册和面试安排流程...")
    
    //desc: 创建用户账号，获取工号
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    
    //desc: 检查用户创建是否成功
    if step1Err != nil {
        //desc: 用户创建失败，记录错误并退出
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    
    //desc: 用户创建成功，记录成功日志
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
    
    //desc: 安排面试时间，联系面试官
    面试时间, 面试官名称, step2Err := step2(用户名)  // 使用step1的输出
    
    //desc: 检查面试安排是否成功
    if step2Err != nil {
        //desc: 面试安排失败，记录错误并退出
        step2.Printf("安排面试失败: %v", step2Err)
        return
    }
    
    //desc: 面试安排成功，记录详细信息
    step2.Printf("✅ 面试安排成功，时间: %s", 面试时间)
}
```

### 3. 步骤描述

```go
func main() {
    //desc: 开始订单处理流程
    fmt.Println("开始订单处理流程...")
    
    //desc: 验证订单信息，检查订单是否有效
    验证结果, step1Err := step1(input["订单号"], input["金额"]){retry:2, timeout:3000}
    
    //desc: 检查订单验证结果
    if step1Err != nil {
        //desc: 订单验证失败，记录错误并退出
        step1.Printf("订单验证失败: %v", step1Err)
        return
    }
    
    //desc: 根据验证结果决定后续流程
    if 验证结果 {
        //desc: 订单验证通过，开始处理支付
        fmt.Println("订单验证通过，开始处理支付...")
        
        //desc: 处理支付流程，调用支付接口
        支付流水号, step2Err := step2(input["订单号"], input["金额"]){retry:3, timeout:5000, priority:"high"}
        
        //desc: 检查支付是否成功
        if step2Err != nil {
            //desc: 支付失败，记录错误并退出
            step2.Printf("支付失败: %v", step2Err)
            return
        }
        
        //desc: 支付成功，记录流水号
        fmt.Printf("订单处理完成，支付流水号: %s\n", 支付流水号)
    } else {
        //desc: 订单验证失败，流程结束
        step1.Printf("订单验证失败")
        fmt.Println("订单验证失败，流程结束")
    }
}
```

**描述语法说明：**
- 使用 `//desc: 描述内容` 格式为下一个语句添加描述
- 描述信息会存储在 `SimpleStatement.Desc` 字段中
- 支持所有类型的语句：`function-call`、`if`、`print`、`var`、`return` 等
- 支持嵌套语句的描述
- 描述信息可以用于界面渲染，让用户清楚了解每个步骤的作用

### 4. 模板字符串

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

### 5. 元数据配置

```go
// 函数调用带元数据配置
工号, 用户名, step1Err := step1(input["用户名"], input["手机号"]){retry:3, timeout:5000, priority:"high"}

// 纯函数调用带元数据
step2(用户名){retry:1, timeout:2000, async:true}

// 支持的元数据类型
step3(){retry:5, timeout:10000, debug:true, mode:"production"}
```

### 6. 条件判断

```go
if step1Err != nil {
    step1.Printf("创建用户失败: %v", step1Err)
    return
} else {
    面试时间, 面试官名称, step2Err := step2(用户名)
    if step2Err != nil {
        step2.Printf("安排面试失败: %v", step2Err)
        return
    }
    step2.Printf("✅ 面试安排成功，时间: %s", 面试时间)
}
```

### 7. 步骤级别日志记录

```go
// 普通日志 - 全局日志
fmt.Printf("✅ 用户创建成功，工号: %s\n", 工号)

// 步骤级别日志 - 可以清楚地知道日志来自哪个步骤
step1.Printf("✅ 用户创建成功，工号: %s", 工号)
step2.Printf("❌ 安排面试失败: %v", step2Err)
step3.Printf("✅ 通知发送成功")

// 支持格式化字符串
step1.Printf("开始执行步骤，参数: %v", args)
step2.Printf("步骤执行完成，耗时: %dms", duration)
```

**步骤日志的优势**：
- **清晰标识**: 每个日志都明确标识来源步骤
- **便于调试**: 可以快速定位问题所在的步骤
- **日志分类**: 支持按步骤分类和管理日志
- **执行追踪**: 可以追踪每个步骤的执行过程

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

### 元数据配置系统

解析器支持在函数调用后添加元数据配置，提供执行时的控制参数：

```go
type SimpleStatement struct {
    // ... 其他字段
    Metadata map[string]interface{} // 元数据配置
}

// 元数据解析示例
// step1(input["用户名"], input["手机号"]){retry:3, timeout:5000, priority:"high"}
// 解析结果:
// Metadata = {
//     "retry": 3,           // 数字类型
//     "timeout": 5000,      // 数字类型  
//     "priority": "high"    // 字符串类型
// }
```

**支持的元数据类型**：
- **数字**: `retry:3`, `timeout:5000`
- **布尔值**: `async:true`, `debug:false`
- **字符串**: `priority:"high"`, `mode:"production"`
- **混合配置**: `{retry:1, timeout:2000, async:true, priority:"high"}`

### 状态管理系统

解析器支持完整的语句执行状态跟踪：

```go
type StatementStatus string

const (
    StatusPending   StatementStatus = "pending"   // 待执行
    StatusRunning   StatementStatus = "running"   // 正在执行
    StatusCompleted StatementStatus = "completed" // 执行完成
    StatusFailed    StatementStatus = "failed"    // 执行失败
    StatusSkipped   StatementStatus = "skipped"   // 跳过执行
)

type SimpleStatement struct {
    // ... 其他字段
    Status     StatementStatus `json:"status"`      // 执行状态
    RetryCount int            `json:"retry_count"` // 重试次数
    Logs       []*StatementLog `json:"logs"`        // 步骤日志
}
```

**状态转换流程**：
1. **pending** → **running**: 语句开始执行
2. **running** → **completed**: 语句执行成功
3. **running** → **failed**: 语句执行失败
4. **running** → **skipped**: 语句被跳过（如条件不满足）

### 重试机制

解析器支持语句级别的重试次数管理：

```go
// 重试次数管理方法
func (s *SimpleStatement) IncrementRetry()    // 增加重试次数
func (s *SimpleStatement) ResetRetry()        // 重置重试次数

// 使用示例
stmt.IncrementRetry()  // 重试次数 +1
if stmt.RetryCount > maxRetries {
    stmt.SetStatus(StatusFailed)
}
```

### 步骤日志系统

解析器支持步骤级别的日志记录和管理：

```go
type StatementLog struct {
    Timestamp time.Time `json:"timestamp"` // 日志时间
    Level     string    `json:"level"`     // 日志级别 (info, warn, error)
    Message   string    `json:"message"`   // 日志内容
    StepName  string    `json:"step_name"` // 所属步骤名称
}

// 添加日志
func (s *SimpleStatement) AddLog(level, message, stepName string)

// 使用示例
stmt.AddLog("info", "开始执行步骤", "step1")
stmt.AddLog("error", "步骤执行失败", "step1")
```

**日志级别**：
- **info**: 一般信息日志
- **warn**: 警告日志
- **error**: 错误日志

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
- `SimpleStatement`: 语句（支持嵌套、元数据、状态管理和日志）
- `VariableInfo`: 变量信息
- `ArgumentInfo`: 参数详细信息
- `StatementStatus`: 语句执行状态
- `StatementLog`: 步骤日志记录

## 测试

```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestSimpleParser_DynamicWorkflow
go test -v -run TestSimpleParser_StaticWorkflow
go test -v -run TestSimpleParser_MixedWorkflow
go test -v -run TestSimpleParser_MetadataSupport
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
- ✅ **元数据配置**: 支持函数调用元数据配置
- ✅ **状态管理**: 完整的语句执行状态跟踪
- ✅ **重试机制**: 支持语句级别的重试管理
- ✅ **步骤日志**: 支持步骤级别的日志记录和分类

## 贡献

欢迎提交Issue和Pull Request来改进这个解析器！

## 许可证

MIT License