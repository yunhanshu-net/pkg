# DryRun设计文档

## 设计背景

### 问题描述

在函数系统中，经常需要执行一些**不可逆的危险操作**：

1. **数据删除操作**
   - 删除用户账户
   - 清空数据库表
   - 删除文件或目录

2. **外部API调用**
   - 调用第三方支付接口
   - 发送短信或邮件
   - 调用云服务API

3. **系统配置修改**
   - 修改数据库配置
   - 更新系统参数
   - 重启服务

4. **批量数据处理**
   - 批量删除数据
   - 批量发送通知
   - 批量更新状态

这些操作一旦执行就无法撤销，如果参数错误可能导致：
- 数据丢失
- 系统故障
- 经济损失
- 安全风险

### 解决方案

DryRun机制提供了一种**预执行预览**的方式：

1. **参数预览**: 在实际执行前展示所有操作参数
2. **用户确认**: 让用户确认参数无误后再执行
3. **风险控制**: 避免因参数错误导致的意外操作

## 设计原则

### 1. 统一API设计

```go
// 实际执行和DryRun使用完全相同的API调用链
req := httpx.Delete("https://api.example.com/users/123").
    Header("Authorization", "Bearer token").
    Body(map[string]interface{}{
        "user_id": "123",
        "reason":  "用户主动删除",
    })

// 实际执行
result, err := req.Do(nil)

// DryRun预览
dryRunCase := req.DryRun()
```

**优势**:
- 代码一致性高
- 维护成本低
- 学习成本低
- 减少重复代码

### 2. 类型安全

```go
// 保持Go的类型安全特性
type DeleteUserReq struct {
    UserID string `json:"user_id"`
    Reason string `json:"reason"`
}

// 编译时检查类型
req := httpx.Delete(fmt.Sprintf("https://api.example.com/users/%s", req.UserID)).
    Body(DeleteUserReq{
        UserID: "123",
        Reason: "用户主动删除",
    })
```

### 3. 清晰的信息展示

```go
// DryRunCase提供完整的操作信息
type DryRunCase interface {
    Type() string
    Map() map[string]interface{}
    Metadata() map[string]interface{}
}
```

## 数据结构

### OnDryRunReq 请求结构体

```go
type OnDryRunReq struct {
    Body interface{} `json:"body"` // 原始请求体，支持任何类型
}

// DecodeBody 解码请求体到指定类型
func (r *OnDryRunReq) DecodeBody(el interface{}) error {
    return jsonx.Convert(r.Body, el)
}
```

**特点**:
- `Body` 字段类型为 `interface{}`，可以接受任何类型的数据
- `DecodeBody` 方法使用 `jsonx.Convert` 进行类型转换，支持复杂的数据结构
- 支持从 map、struct、JSON 字符串等多种格式解码

**使用示例**:
```go
// 解码到结构体
var deleteReq DeleteUserReq
if err := req.DecodeBody(&deleteReq); err != nil {
    return &usercall.OnDryRunResp{
        Valid:   false,
        Message: "参数解析失败: " + err.Error(),
    }, nil
}

// 解码到 map
var params map[string]interface{}
if err := req.DecodeBody(&params); err != nil {
    return &usercall.OnDryRunResp{
        Valid:   false,
        Message: "参数解析失败: " + err.Error(),
    }, nil
}
```

### OnDryRunResp 响应结构体

```go
type OnDryRunResp struct {
    Valid   bool           `json:"valid"`   // 是否有效
    Cases   []DryRunCase   `json:"cases"`   // DryRun 案例列表
    Message string         `json:"message"` // 提示信息
}
```

## 使用场景

### 场景1: 函数系统中的危险操作

```go
// 删除用户函数
func DeleteUser(ctx *runner.Context, req *DeleteUserReq, resp response.Response) error {
    // 实际执行删除操作
    result, err := httpx.Delete(fmt.Sprintf("https://api.company.com/users/%s", req.UserID)).
        Header("Authorization", "Bearer [TOKEN]").
        Body(map[string]interface{}{
            "user_id": req.UserID,
            "reason":  req.Reason,
        }).
        Do(nil)
    
    if err != nil {
        return err
    }
    
    return resp.Form(map[string]interface{}{
        "success": true,
        "message": "用户删除成功",
    }).Build()
}

// DryRun实现
func OnDryRun(ctx *runner.Context, req *usercall.OnDryRunReq) (*usercall.OnDryRunResp, error) {
    var deleteReq DeleteUserReq
    if err := req.DecodeBody(&deleteReq); err != nil {
        return &usercall.OnDryRunResp{
            Valid:   false,
            Message: "参数解析失败: " + err.Error(),
        }, nil
    }
    
    // 获取DryRun案例
    dryRunCase := httpx.Delete(fmt.Sprintf("https://api.company.com/users/%s", deleteReq.UserID)).
        Header("Authorization", "Bearer [TOKEN]").
        Body(map[string]interface{}{
            "user_id": deleteReq.UserID,
            "reason":  deleteReq.Reason,
        }).
        DryRun()
    
    return &usercall.OnDryRunResp{
        Valid:   true,
        Cases:   []usercall.DryRunCase{dryRunCase},
        Message: "预览删除用户操作",
    }, nil
}
```

### 场景2: 前端界面预览

前端可以调用DryRun接口，在用户点击"确认执行"前显示操作预览：

```javascript
// 前端调用DryRun
const preview = await fetch('/api/dryrun', {
    method: 'POST',
    body: JSON.stringify({
        user_id: "123",
        reason: "用户主动删除"
    })
});

// 显示预览信息
console.log("即将执行的操作:");
console.log("方法:", preview.method);
console.log("URL:", preview.url);
console.log("请求体:", preview.body);

// 用户确认后执行
if (confirm("确认执行此操作？")) {
    await fetch('/api/execute', {
        method: 'POST',
        body: JSON.stringify({
            user_id: "123",
            reason: "用户主动删除"
        })
    });
}
```

### 场景3: 批量操作预览

```go
// 批量删除用户
func BatchDeleteUsers(ctx *runner.Context, req *BatchDeleteReq, resp response.Response) error {
    var cases []usercall.DryRunCase
    
    for _, userID := range req.UserIDs {
        dryRunCase := httpx.Delete(fmt.Sprintf("https://api.company.com/users/%s", userID)).
            Header("Authorization", "Bearer [TOKEN]").
            Body(map[string]interface{}{
                "user_id": userID,
                "reason":  req.Reason,
            }).
            DryRun()
        
        cases = append(cases, dryRunCase)
    }
    
    // 返回预览信息
    return resp.Form(map[string]interface{}{
        "operations": cases,
        "total":      len(cases),
        "warning":    fmt.Sprintf("即将删除 %d 个用户", len(cases)),
    }).Build()
}
```

### 场景4: 复杂操作链预览

```go
// 复杂的多步骤操作
func ComplexOperation(ctx *runner.Context, req *ComplexReq, resp response.Response) error {
    var cases []DryRunCase
    
    // 步骤1: 备份数据
    backupCase := httpx.Post("https://api.company.com/backup").
        Header("Authorization", "Bearer [TOKEN]").
        Body(map[string]interface{}{
            "table": "users",
            "date":  time.Now().Format("2006-01-02"),
        }).
        DryRun()
    cases = append(cases, *backupCase)
    
    // 步骤2: 删除数据
    deleteCase := httpx.Delete(fmt.Sprintf("https://api.company.com/users/%s", req.UserID)).
        Header("Authorization", "Bearer [TOKEN]").
        Body(map[string]interface{}{
            "user_id": req.UserID,
            "reason":  req.Reason,
        }).
        DryRun()
    cases = append(cases, *deleteCase)
    
    // 步骤3: 发送通知
    notifyCase := httpx.Post("https://api.company.com/notify").
        Header("Authorization", "Bearer [TOKEN]").
        Body(map[string]interface{}{
            "user_id": req.UserID,
            "type":    "user_deleted",
            "message": "用户已被删除",
        }).
        DryRun()
    cases = append(cases, *notifyCase)
    
    return resp.Form(map[string]interface{}{
        "operations": cases,
        "steps":      len(cases),
        "warning":    "此操作包含多个步骤，请谨慎确认",
    }).Build()
}
```

## 最佳实践

### 1. 错误处理

```go
// 在DryRun中也要处理错误
func OnDryRun(ctx *runner.Context, req *usercall.OnDryRunReq) (*usercall.OnDryRunResp, error) {
    var deleteReq DeleteUserReq
    if err := req.DecodeBody(&deleteReq); err != nil {
        return &usercall.OnDryRunResp{
            Valid:   false,
            Message: "参数解析失败: " + err.Error(),
        }, nil
    }
    
    // 参数验证
    if deleteReq.UserID == "" {
        return &usercall.OnDryRunResp{
            Valid:   false,
            Message: "用户ID不能为空",
        }, nil
    }
    
    // 获取DryRun案例
    dryRunCase := httpx.Delete(fmt.Sprintf("https://api.company.com/users/%s", deleteReq.UserID)).
        Header("Authorization", "Bearer [TOKEN]").
        Body(map[string]interface{}{
            "user_id": deleteReq.UserID,
            "reason":  deleteReq.Reason,
        }).
        DryRun()
    
    return &usercall.OnDryRunResp{
        Valid:   true,
        Cases:   []usercall.DryRunCase{*dryRunCase},
        Message: "预览删除用户操作",
    }, nil
}
```

### 2. 参数验证

```go
// 在DryRun中进行参数验证
func validateDeleteUser(req *DeleteUserReq) error {
    if req.UserID == "" {
        return fmt.Errorf("用户ID不能为空")
    }
    
    if req.Reason == "" {
        return fmt.Errorf("删除原因不能为空")
    }
    
    // 检查用户是否存在
    if !userExists(req.UserID) {
        return fmt.Errorf("用户不存在: %s", req.UserID)
    }
    
    return nil
}
```

### 3. 安全考虑

```go
// 敏感信息处理
func sanitizeDryRunCase(case *DryRunCase) *DryRunCase {
    sanitized := *case
    
    // 隐藏敏感信息
    if token, exists := sanitized.Headers["Authorization"]; exists {
        if len(token) > 10 {
            sanitized.Headers["Authorization"] = token[:10] + "..."
        }
    }
    
    return &sanitized
}
```

### 4. 日志记录

```go
// 记录DryRun操作
func logDryRunOperation(ctx *runner.Context, req *DryRunReq, resp *DryRunResp) {
    logger.Info("DryRun操作",
        "user_id", ctx.UserID,
        "operation", "delete_user",
        "valid", resp.Valid,
        "cases_count", len(resp.Cases),
    )
}
```

## 扩展性

### 1. 支持多种操作类型

```go
// 可以扩展到其他类型的操作
type DryRunCase struct {
    Type        string            `json:"type"`        // 操作类型: http, sql, file
    Method      string            `json:"method"`      // HTTP方法或SQL类型
    Url         string            `json:"url"`         // URL或SQL语句
    Headers     map[string]string `json:"headers"`     // 请求头或SQL参数
    Body        string            `json:"body"`        // 请求体
    Description string            `json:"description"` // 操作描述
}
```

### 2. 支持混合操作

```go
// 支持HTTP + SQL的混合操作
func MixedOperation(ctx *runner.Context, req *MixedReq, resp response.Response) error {
    var cases []DryRunCase
    
    // HTTP操作
    httpCase := httpx.Delete("https://api.company.com/users/123").
        Header("Authorization", "Bearer token").
        DryRun()
    cases = append(cases, *httpCase)
    
    // SQL操作
    sqlCase := &DryRunCase{
        Type:        "sql",
        Method:      "DELETE",
        Url:         "DELETE FROM users WHERE id = ?",
        Body:        "123",
        Description: "删除用户记录",
    }
    cases = append(cases, *sqlCase)
    
    return resp.Form(map[string]interface{}{
        "operations": cases,
    }).Build()
}
```

## 总结

DryRun机制通过以下方式解决了危险操作的安全问题：

1. **预执行预览**: 在实际执行前展示所有操作参数
2. **用户确认**: 让用户确认参数无误后再执行
3. **统一API**: 实际执行和DryRun使用相同的API调用链
4. **类型安全**: 保持Go的类型安全特性
5. **易于集成**: 与现有函数系统无缝集成

这种设计既保证了操作的安全性，又保持了代码的简洁性和一致性。

## 连通性检查功能

### 设计目标

在DryRun基础上增加**网络连通性检查**功能，帮助用户在实际执行前验证：

1. **API可达性**: 检查目标API是否可访问
2. **认证有效性**: 验证Token/API Key是否有效  
3. **接口可用性**: 确认API接口是否正常工作
4. **网络质量**: 检测网络延迟和响应时间

### 核心功能

#### ConnectivityCheck 方法
```go
// 执行连通性检查，返回HttpRequest支持链式调用
func (h *HttpRequest) ConnectivityCheck() *HttpRequest

// 获取连通性检查结果
func (h *HttpRequest) GetConnectivity() *ConnectivityResult
```

#### ConnectivityResult 结构体
```go
type ConnectivityResult struct {
    // 基础网络状态
    Reachable     bool                   `json:"reachable"`     // 是否可达
    StatusCode    int                    `json:"status_code"`   // HTTP状态码
    Error         string                 `json:"error,omitempty"` // 错误信息
    
    // 性能指标
    Latency       time.Duration          `json:"latency"`       // 网络延迟
    ResponseTime  time.Duration          `json:"response_time"` // 响应时间
    
    // 网络诊断
    Timeout       bool                   `json:"timeout"`       // 是否超时
    DNSResolved   bool                   `json:"dns_resolved"`  // DNS是否解析成功
    
    // 服务器响应头信息
    Server        string                 `json:"server"`        // 服务器标识
    ContentType   string                 `json:"content_type"`  // 内容类型
    ContentLength int64                  `json:"content_length"` // 内容长度
    
    // SSL信息
    SSLValid      bool                   `json:"ssl_valid"`     // SSL证书是否有效（需要实际SSL验证时设置）
    
    // 扩展数据
    Metadata      map[string]interface{} `json:"metadata"`      // 扩展信息
}
```

### 使用示例

#### 基础连通性检查
```go
// 执行连通性检查
req := httpx.Get("https://httpbin.org/get").
    Header("User-Agent", "test-agent").
    ConnectivityCheck()

// 获取检查结果
connectivity := req.GetConnectivity()
if connectivity != nil {
    fmt.Printf("可达性: %t\n", connectivity.Reachable)
    fmt.Printf("状态码: %d\n", connectivity.StatusCode)
    fmt.Printf("响应时间: %v\n", connectivity.ResponseTime)
    fmt.Printf("服务器: %s\n", connectivity.Server)
}
```

#### DryRun包含连通性检查
```go
// 创建HTTP请求并执行连通性检查
dryRunCase := httpx.Delete("https://api.company.com/users/123").
    Header("Authorization", "Bearer [TOKEN]").
    Body(map[string]interface{}{
        "user_id": "123",
        "reason":  "用户主动删除",
    }).
    ConnectivityCheck().
    DryRun()

// 连通性检查结果已包含在DryRun案例中
if dryRunCase.Connectivity != nil {
    conn := dryRunCase.Connectivity
    fmt.Printf("API可达性: %t\n", conn.Reachable)
    fmt.Printf("状态码: %d\n", conn.StatusCode)
    fmt.Printf("响应时间: %v\n", conn.ResponseTime)
    
    // 根据结果给出建议
    if !conn.Reachable {
        fmt.Printf("⚠️  警告: API不可达，请检查网络连接\n")
    } else if conn.StatusCode == 401 {
        fmt.Printf("⚠️  警告: 认证失败，请检查Token\n")
    } else if conn.StatusCode == 200 {
        fmt.Printf("✅ API正常，可以执行操作\n")
    }
}
```

### 技术实现

#### HEAD请求优化
- 使用HEAD请求进行连通性检查，避免下载响应体
- 设置合理的超时时间（默认5秒）
- 支持自定义请求头传递

#### 错误分析
- 自动识别超时错误
- 区分DNS解析失败和网络连接失败
- 提供详细的错误信息

#### 性能指标
- 测量网络延迟
- 记录响应时间
- 支持SSL证书验证（需要实际SSL验证时设置SSLValid字段）

#### SSLValid字段说明
SSLValid字段通过实际的SSL证书验证来设置，支持以下场景：

```go
// 1. HTTPS请求自动验证SSL证书
httpsResult := httpx.Get("https://api.example.com").ConnectivityCheck()
if httpsResult.SSLValid {
    // SSL证书有效
} else {
    // SSL证书无效或自签名
}

// 2. 同时测试HTTP和HTTPS时
httpResult := httpx.Get("http://api.example.com").ConnectivityCheck()
httpsResult := httpx.Get("https://api.example.com").ConnectivityCheck()

// 比较SSL有效性
if httpResult.Reachable && httpsResult.Reachable {
    if httpsResult.SSLValid {
        // HTTPS的SSL证书有效
    }
}

// 3. HTTP请求的SSLValid始终为false
httpResult := httpx.Get("http://api.example.com").ConnectivityCheck()
// httpResult.SSLValid 始终为 false
```

**SSL验证规则**:
- HTTPS请求：自动验证SSL证书有效性
- HTTP请求：SSLValid始终为false
- 自签名证书：检测为无效
- 过期证书：检测为无效
- 域名不匹配：检测为无效

### 优势

1. **预执行验证**: 在实际操作前验证API可用性
2. **风险控制**: 避免因网络问题导致的失败操作
3. **用户体验**: 提供详细的状态信息和操作建议
4. **链式调用**: 与现有API完美集成，支持链式调用
5. **标准化**: 返回标准化的数据结构，便于前端处理 