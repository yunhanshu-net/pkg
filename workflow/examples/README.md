# 工作流示例项目

这个目录包含了使用工作流引擎的各种示例项目，展示了不同的使用场景和功能。

## 示例列表

### 1. 用户注册工作流 (user_registration)

**功能**：演示完整的用户注册流程，包括用户创建、部门分配、邮件发送和档案创建。

**特点**：
- 多步骤工作流
- 参数传递链
- 成功场景演示
- 详细的执行日志

**运行方式**：
```bash
cd user_registration
go run main.go
```

**预期输出**：
```
🚀 用户注册工作流演示
========================

📋 执行步骤: step1 - 创建用户账号，获取工号
📥 输入参数: map[email:zhangsan@example.com phone:13800138000 username:张三]
✅ 用户创建成功，用户ID: USER_1234567890

📋 执行步骤: step2 - 分配部门
📥 输入参数: map[department:技术部 position:高级工程师 userId:USER_1234567890]
✅ 部门分配成功: 已分配到 技术部 部门

📋 执行步骤: step3 - 发送欢迎邮件
📥 输入参数: map[department:技术部 email:zhangsan@example.com userId:USER_1234567890 username:张三]
✅ 欢迎邮件发送成功

📋 执行步骤: step4 - 创建用户档案
📥 输入参数: map[department:技术部 position:高级工程师 userId:USER_1234567890 username:张三]
✅ 用户档案创建成功，档案ID: PROFILE_1234567890

🎉 用户注册流程完成！用户: 张三, ID: USER_1234567890, 档案: PROFILE_1234567890
```

### 2. 错误处理演示 (error_handling)

**功能**：演示工作流中的错误处理机制，包括验证失败、业务逻辑错误等场景。

**特点**：
- 错误场景模拟
- 条件判断逻辑
- 错误传播机制
- 工作流中断处理

**运行方式**：
```bash
cd error_handling
go run main.go
```

**预期输出**：
```
🚨 错误处理工作流演示
========================

📋 执行步骤: step1 - 验证用户信息
📥 输入参数: map[email:lisi@example.com phone:13900139000 username:李四]
❌ 用户信息无效: 用户名已存在，请选择其他用户名

❌ 工作流因错误中断
```

### 3. 并发执行演示 (concurrent_execution)

**功能**：演示多个工作流实例的并发执行，展示系统的并发处理能力。

**特点**：
- 并发执行多个工作流
- 时间戳显示执行顺序
- 性能统计
- 资源管理

**运行方式**：
```bash
cd concurrent_execution
go run main.go
```

**预期输出**：
```
🚀 并发工作流执行演示
========================

🚀 开始并发执行 5 个工作流实例...

[15:04:05] 📋 执行步骤: step1 - 创建用户
[15:04:05] 📋 执行步骤: step1 - 创建用户
[15:04:05] 📋 执行步骤: step1 - 创建用户
[15:04:05] 📋 执行步骤: step1 - 创建用户
[15:04:05] 📋 执行步骤: step1 - 创建用户
[15:04:05] ✅ 工作流正常结束
[15:04:05] ✅ 工作流正常结束
[15:04:05] ✅ 工作流正常结束
[15:04:05] ✅ 工作流正常结束
[15:04:05] ✅ 工作流正常结束

⏱️  总执行时间: 200ms
📊 成功执行: 5/5 个工作流实例
🚀 平均每个工作流执行时间: 40ms
```

## 技术特性展示

### 1. 参数映射
- **形参名 vs 实例名**：步骤定义中的形参名正确映射到调用时的实例名
- **类型推断**：从步骤定义中自动获取参数类型信息
- **值传递**：支持变量值的正确传递和存储

### 2. 错误处理
- **条件判断**：支持 `if` 语句进行条件判断
- **错误传播**：错误信息正确传播到上层
- **工作流中断**：支持 `return` 语句中断工作流

### 3. 并发支持
- **多实例执行**：支持同时执行多个工作流实例
- **上下文管理**：每个实例有独立的执行上下文
- **资源隔离**：实例间相互独立，不会相互影响

### 4. 日志系统
- **步骤级日志**：每个步骤可以记录独立的日志
- **全局日志**：支持全局日志记录
- **时间戳**：日志包含时间戳信息

### 5. 回调机制
- **函数调用回调**：`OnFunctionCall` 处理步骤执行
- **状态更新回调**：`OnWorkFlowUpdate` 处理状态更新
- **结束回调**：`OnWorkFlowExit` 和 `OnWorkFlowReturn` 处理结束逻辑

## 自定义开发

### 1. 创建新的示例项目

```bash
mkdir my_demo
cd my_demo
go mod init my-demo
```

### 2. 基本代码结构

```go
package main

import (
    "context"
    "github.com/yunhanshu-net/pkg/workflow"
)

func main() {
    // 1. 定义工作流代码
    workflowCode := `...`
    
    // 2. 解析工作流
    parser := workflow.NewSimpleParser()
    parseResult := parser.ParseWorkflow(workflowCode)
    
    // 3. 创建执行器
    executor := workflow.NewExecutor()
    
    // 4. 设置回调函数
    executor.OnFunctionCall = func(ctx context.Context, step workflow.SimpleStep, in *workflow.ExecutorIn) (*workflow.ExecutorOut, error) {
        // 实现业务逻辑
    }
    
    // 5. 执行工作流
    ctx := context.Background()
    executor.Start(ctx, parseResult)
}
```

### 3. 工作流代码语法

```go
// 输入变量定义
var input = map[string]interface{}{
    "参数名": "参数值",
}

// 步骤定义
step1 = 函数名(
    参数名: 类型 "描述"
) -> (
    返回值名: 类型 "描述"
);

// 主函数
func main() {
    // 步骤调用
    返回值, 错误 := step1(参数){元数据}
    
    // 条件判断
    if 错误 != nil {
        return
    }
    
    // 打印日志
    step1.Printf("日志信息")
}
```

## 注意事项

1. **依赖管理**：确保正确设置 `replace` 指令指向工作流包
2. **错误处理**：在回调函数中正确处理错误情况
3. **资源管理**：长时间运行的工作流需要考虑资源清理
4. **并发安全**：多实例执行时注意共享资源的安全性
5. **性能优化**：大量并发时考虑使用连接池等优化手段
