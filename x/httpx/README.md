# httpx - 链式调用HTTP库

httpx是一个支持链式调用的HTTP库，提供了简洁的API、DryRun功能和响应体绑定。

## 特性

- 链式调用API
- 支持DryRun预览
- 响应体绑定
- 类型安全
- 易于使用

## DryRun设计背景

### 为什么需要DryRun？

在函数系统中，经常需要执行一些**危险操作**，比如：
- 删除用户数据
- 调用外部API
- 修改系统配置
- 批量数据处理

这些操作一旦执行就无法撤销，如果参数错误可能导致严重后果。DryRun机制可以在**实际执行前**预览操作参数，让用户确认无误后再执行。

### 设计原则

1. **统一API**: 实际执行和DryRun使用完全相同的API调用链
2. **参数预览**: 清晰展示所有请求参数（URL、方法、请求头、请求体）
3. **类型安全**: 保持Go的类型安全特性
4. **易于集成**: 与现有函数系统无缝集成

### 使用场景

#### 1. 函数系统中的危险操作预览

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

// DryRun实现 - 预览删除操作
func OnDryRun(ctx *runner.Context, req *usercall.OnDryRunReq) (*usercall.OnDryRunResp, error) {
    var deleteReq DeleteUserReq
    if err := req.DecodeBody(&deleteReq); err != nil {
        return &usercall.OnDryRunResp{
            Valid:   false,
            Message: "参数解析失败: " + err.Error(),
        }, nil
    }
    
    // 获取DryRun案例 - 不实际执行，只预览参数
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

#### 2. 前端界面预览

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
```

#### 3. 批量操作预览

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
    }).Build()
}
```

## 基本用法

### 实际执行HTTP请求

```go
package main

import (
    "fmt"
    "github.com/your-project/pkg/x/httpx"
)

func main() {
    // DELETE请求
    result, err := httpx.Delete("https://api.example.com/users/123").
        Header("Authorization", "Bearer token123").
        Body(map[string]interface{}{
            "user_id": "123",
            "reason":  "用户主动删除",
        }).
        Do(nil)
    
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
        return
    }
    
    fmt.Printf("请求成功，状态码: %d\n", result.Code)
    fmt.Printf("响应内容: %s\n", result.ResBodyString)
}
```

### 响应体绑定

```go
package main

import (
    "fmt"
    "github.com/your-project/pkg/x/httpx"
)

// 用户结构体
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// API响应结构体
type ApiResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func main() {
    // 绑定到结构体
    var user User
    result, err := httpx.Get("https://api.example.com/users/123").
        Header("Authorization", "Bearer token123").
        Do(&user)
    
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
        return
    }
    
    if result.OK() {
        fmt.Printf("用户信息: %+v\n", user)
    }
    
    // 绑定到API响应结构体
    var apiResp ApiResponse
    result, err = httpx.Post("https://api.example.com/users").
        Header("Authorization", "Bearer token123").
        Body(map[string]interface{}{
            "name":  "张三",
            "email": "zhangsan@example.com",
        }).
        Do(&apiResp)
    
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
        return
    }
    
    if result.OK() {
        fmt.Printf("API响应: %+v\n", apiResp)
    }
    
    // 绑定到map
    var responseMap map[string]interface{}
    result, err = httpx.Get("https://api.example.com/users").
        Header("Authorization", "Bearer token123").
        Do(&responseMap)
    
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
        return
    }
    
    if result.OK() {
        fmt.Printf("响应数据: %+v\n", responseMap)
    }
}
```

### DryRun预览

```go
package main

import (
    "fmt"
    "github.com/your-project/pkg/x/httpx"
)

func main() {
    // 获取DryRun案例，不实际执行
    dryRunCase := httpx.Delete("https://api.example.com/users/123").
        Header("Authorization", "Bearer token123").
        Body(map[string]interface{}{
            "user_id": "123",
            "reason":  "用户主动删除",
        }).
        DryRun()
    
    // 推荐：直接访问结构体字段
    fmt.Printf("HTTP方法: %s\n", dryRunCase.Method)
    fmt.Printf("URL: %s\n", dryRunCase.Url)
    fmt.Printf("请求头: %v\n", dryRunCase.Headers)
    fmt.Printf("请求体: %s\n", dryRunCase.Body)
    fmt.Printf("描述: %s\n", dryRunCase.Description)
    fmt.Printf("元数据: %v\n", dryRunCase.Meta)
    
    // 也可以通过接口方法访问
    fmt.Printf("HTTP方法: %s\n", dryRunCase.Map()["method"])
    fmt.Printf("URL: %s\n", dryRunCase.Map()["url"])
    fmt.Printf("请求头: %v\n", dryRunCase.Map()["headers"])
    fmt.Printf("请求体: %s\n", dryRunCase.Map()["body"])
    fmt.Printf("描述: %s\n", dryRunCase.Map()["description"])
}
```

## API参考

### 创建请求

```go
// 创建不同类型的请求
deleteReq := httpx.Delete("https://api.example.com/users/123")
getReq := httpx.Get("https://api.example.com/users")
postReq := httpx.Post("https://api.example.com/users")
putReq := httpx.Put("https://api.example.com/users/123")
```

### 设置请求头

```go
// 单个请求头
req := httpx.Post("https://api.example.com/users").
    Header("Authorization", "Bearer token123").
    Header("Content-Type", "application/json")

// 批量设置请求头
headers := map[string]string{
    "Authorization": "Bearer token123",
    "Content-Type":  "application/json",
    "User-Agent":    "MyApp/1.0",
}
req := httpx.Post("https://api.example.com/users").
    Headers(headers)
```

### 设置请求体

```go
// 结构体
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

user := User{Name: "张三", Email: "zhangsan@example.com"}
req := httpx.Post("https://api.example.com/users").
    Body(user)

// Map
req := httpx.Post("https://api.example.com/users").
    Body(map[string]interface{}{
        "name":  "张三",
        "email": "zhangsan@example.com",
    })

// 字符串
req := httpx.Post("https://api.example.com/users").
    Body(`{"name":"张三","email":"zhangsan@example.com"}`)
```

### 设置超时

```go
req := httpx.Get("https://api.example.com/users").
    Timeout(30 * time.Second)
```

### 执行请求

```go
// 执行请求并绑定响应到结构体
var user User
result, err := req.Do(&user)

// 执行请求并绑定响应到map
var response map[string]interface{}
result, err := req.Do(&response)

// 执行请求不绑定响应
result, err := req.Do(nil)

// 执行请求并返回字符串响应
result, err := req.DoString()
```

### DryRun预览

```go
// 获取DryRun案例
dryRunCase := req.DryRun()

// 输出预览信息
fmt.Printf("方法: %s\n", dryRunCase.Map()["method"])
fmt.Printf("URL: %s\n", dryRunCase.Map()["url"])
fmt.Printf("请求头: %v\n", dryRunCase.Map()["headers"])
fmt.Printf("请求体: %s\n", dryRunCase.Map()["body"])
fmt.Printf("描述: %s\n", dryRunCase.Map()["description"])
```

## 完整示例

### 用户管理API

```go
package main

import (
    "fmt"
    "github.com/your-project/pkg/x/httpx"
)

// 用户结构体
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// API响应结构体
type ApiResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

// 创建用户
func CreateUser(name, email string) (*User, error) {
    var apiResp ApiResponse
    
    result, err := httpx.Post("https://api.example.com/users").
        Header("Authorization", "Bearer token123").
        Body(map[string]interface{}{
            "name":  name,
            "email": email,
        }).
        Do(&apiResp)
    
    if err != nil {
        return nil, fmt.Errorf("创建用户失败: %v", err)
    }
    
    if !result.OK() {
        return nil, fmt.Errorf("创建用户失败，状态码: %d", result.Code)
    }
    
    // 解析用户数据
    if userData, ok := apiResp.Data.(map[string]interface{}); ok {
        user := &User{
            ID:    userData["id"].(string),
            Name:  userData["name"].(string),
            Email: userData["email"].(string),
        }
        return user, nil
    }
    
    return nil, fmt.Errorf("解析用户数据失败")
}

// 删除用户
func DeleteUser(userID string) error {
    var apiResp ApiResponse
    
    result, err := httpx.Delete(fmt.Sprintf("https://api.example.com/users/%s", userID)).
        Header("Authorization", "Bearer token123").
        Body(map[string]interface{}{
            "user_id": userID,
            "reason":  "用户主动删除",
        }).
        Do(&apiResp)
    
    if err != nil {
        return fmt.Errorf("删除用户失败: %v", err)
    }
    
    if !result.OK() {
        return fmt.Errorf("删除用户失败，状态码: %d", result.Code)
    }
    
    fmt.Printf("用户删除成功: %s\n", apiResp.Message)
    return nil
}

// 获取用户列表
func GetUsers() ([]User, error) {
    var apiResp ApiResponse
    
    result, err := httpx.Get("https://api.example.com/users").
        Header("Authorization", "Bearer token123").
        Do(&apiResp)
    
    if err != nil {
        return nil, fmt.Errorf("获取用户列表失败: %v", err)
    }
    
    if !result.OK() {
        return nil, fmt.Errorf("获取用户列表失败，状态码: %d", result.Code)
    }
    
    // 解析用户列表
    var users []User
    if usersData, ok := apiResp.Data.([]interface{}); ok {
        for _, userData := range usersData {
            if userMap, ok := userData.(map[string]interface{}); ok {
                user := User{
                    ID:    userMap["id"].(string),
                    Name:  userMap["name"].(string),
                    Email: userMap["email"].(string),
                }
                users = append(users, user)
            }
        }
    }
    
    return users, nil
}

// 获取单个用户
func GetUser(userID string) (*User, error) {
    var user User
    
    result, err := httpx.Get(fmt.Sprintf("https://api.example.com/users/%s", userID)).
        Header("Authorization", "Bearer token123").
        Do(&user)
    
    if err != nil {
        return nil, fmt.Errorf("获取用户失败: %v", err)
    }
    
    if !result.OK() {
        return nil, fmt.Errorf("获取用户失败，状态码: %d", result.Code)
    }
    
    return &user, nil
}

// DryRun预览
func PreviewUserOperations() {
    fmt.Println("=== 预览用户操作 ===")
    
    // 预览创建用户
    createCase := httpx.Post("https://api.example.com/users").
        Header("Authorization", "Bearer token123").
        Body(map[string]interface{}{
            "name":  "张三",
            "email": "zhangsan@example.com",
        }).
        DryRun()
    
    fmt.Printf("创建用户预览:\n")
    fmt.Printf("  方法: %s\n", createCase.Method)
    fmt.Printf("  URL: %s\n", createCase.Url)
    fmt.Printf("  请求体: %s\n", createCase.Body)
    
    // 预览删除用户
    deleteCase := httpx.Delete("https://api.example.com/users/123").
        Header("Authorization", "Bearer token123").
        Body(map[string]interface{}{
            "user_id": "123",
            "reason":  "用户主动删除",
        }).
        DryRun()
    
    fmt.Printf("删除用户预览:\n")
    fmt.Printf("  方法: %s\n", deleteCase.Method)
    fmt.Printf("  URL: %s\n", deleteCase.Url)
    fmt.Printf("  请求体: %s\n", deleteCase.Body)
}

func main() {
    // 预览操作
    PreviewUserOperations()
    
    // 实际执行
    user, err := CreateUser("张三", "zhangsan@example.com")
    if err != nil {
        fmt.Printf("错误: %v\n", err)
    } else {
        fmt.Printf("创建用户成功: %+v\n", user)
    }
    
    users, err := GetUsers()
    if err != nil {
        fmt.Printf("错误: %v\n", err)
    } else {
        fmt.Printf("用户列表: %+v\n", users)
    }
    
    if err := DeleteUser("123"); err != nil {
        fmt.Printf("错误: %v\n", err)
    }
}
```

## 最佳实践

1. **统一API**: 实际执行和DryRun使用完全相同的API调用链
2. **响应绑定**: 使用Do()方法绑定响应到结构体
3. **错误处理**: 始终检查返回的错误和状态码
4. **超时设置**: 为长时间运行的请求设置合适的超时时间
5. **请求头**: 使用Headers()方法批量设置多个请求头
6. **DryRun**: 在函数系统中使用DryRun预览API调用

## 注意事项

- 默认超时时间为60秒
- 默认Content-Type为application/json
- DryRun不会实际执行HTTP请求
- 所有请求都支持链式调用
- Do()方法会自动解析JSON响应到结构体 