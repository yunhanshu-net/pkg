# Query库使用文档

## 📖 概述

Query库是一个强大的数据库查询和分页工具，提供了安全、灵活的查询条件构建和分页功能。支持多种查询操作符、排序、分页，并具备SQL注入防护机制。

## 🚀 核心特性

- ✅ **Search标签驱动**：通过`search`标签声明查询能力，自动构建安全配置
- ✅ **零配置查询**：使用`AutoSearchPaginated`实现开箱即用的安全分页查询
- ✅ **类型安全的分页查询**：基于泛型的类型安全分页
- ✅ **多种查询操作符**：支持eq、like、in、gt、gte、lt、lte、not_eq、not_like、not_in等
- ✅ **灵活的排序**：支持单字段和多字段排序
- ✅ **权限控制**：结合`permission`标签自动控制字段查询权限
- ✅ **安全防护**：内置SQL注入防护机制和字段白名单验证
- ✅ **向下兼容**：仍支持手动QueryConfig配置方式

## 📋 基础用法

### 1. 推荐用法：基于Search标签的自动配置

```go
import (
    "context"
    "github.com/yunhanshu-net/pkg/query"
)

// 定义数据模型 - 通过search标签声明查询能力
type Product struct {
    ID       int     `gorm:"primaryKey" json:"id" runner:"code:id;name:产品ID"`
    
    // 支持模糊搜索和精确匹配
    Name     string  `json:"name" runner:"code:name;name:产品名称" search:"like,eq"`
    
    // 支持精确匹配和多选
    Category string  `json:"category" runner:"code:category;name:产品分类" search:"eq,in"`
    
    // 支持范围查询
    Price    float64 `json:"price" runner:"code:price;name:产品价格" search:"gte,lte,eq"`
    
    // 支持精确匹配
    Status   string  `json:"status" runner:"code:status;name:产品状态" search:"eq"`
    
    // 敏感字段不添加search标签，自动禁止查询
    CostPrice float64 `json:"cost_price" runner:"code:cost_price;name:成本价"`
}

// 零配置自动搜索分页查询
func ProductList(ctx context.Context, db *gorm.DB, pageInfo *query.PageInfoReq) (*query.PaginatedTable[[]Product], error) {
    var products []Product
    // 使用AutoSearchPaginated，自动根据search标签构建查询配置
    return query.AutoSearchPaginated(db, &Product{}, products, pageInfo)
}
```

### 2. 在HTTP接口中使用

```go
// 请求结构体
type ProductListReq struct {
    query.PageInfoReq `form:",inline"` // 内嵌分页参数
    // 可以添加其他业务参数
}

// HTTP处理函数
func ProductListHandler(c *gin.Context) {
    var req ProductListReq
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    db := getDB() // 获取数据库连接
    var products []Product
    
    // 零配置查询，自动根据Product模型的search标签进行安全验证
    result, err := query.AutoSearchPaginated(db, &Product{}, products, &req.PageInfoReq)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, result)
}
```

## 🏷️ Search标签语法

### 基础语法
```go
search:"操作符1,操作符2,操作符3"
```

### 支持的操作符
- `eq` - 精确匹配
- `like` - 模糊搜索  
- `in` - 包含查询（多选）
- `gte` - 大于等于
- `lte` - 小于等于
- `gt` - 大于
- `lt` - 小于
- `not_eq` - 不等于
- `not_like` - 不模糊匹配
- `not_in` - 不包含

### 完整的模型示例

```go
type Product struct {
    // 基础字段 - 不可搜索（没有search标签）
    ID int `json:"id" gorm:"primaryKey" runner:"code:id;name:产品ID"`

    // 文本搜索字段
    Name string `json:"name" gorm:"column:name;comment:产品名称" 
                 runner:"code:name;name:产品名称" 
                 search:"like,eq"` // 支持模糊搜索和精确匹配

    // 分类选择字段
    Category string `json:"category" gorm:"column:category;comment:产品分类" 
                     runner:"code:category;name:产品分类" 
                     search:"eq,in"` // 支持精确匹配和多选

    // 价格范围搜索
    Price float64 `json:"price" gorm:"column:price;comment:产品价格" 
                   runner:"code:price;name:产品价格" 
                   search:"gte,lte,eq"` // 支持范围查询和精确匹配

    // 库存数量搜索
    Stock int `json:"stock" gorm:"column:stock;comment:库存数量" 
               runner:"code:stock;name:库存数量" 
               search:"gte,lte,eq,gt,lt"` // 支持所有数值比较

    // 状态搜索
    Status string `json:"status" gorm:"column:status;comment:产品状态" 
                   runner:"code:status;name:产品状态" 
                   search:"eq,not_eq"` // 支持正向和否定匹配

    // 标签搜索
    Tags string `json:"tags" gorm:"column:tags;comment:产品标签" 
                 runner:"code:tags;name:产品标签" 
                 search:"like,in,not_like,not_in"` // 支持模糊和否定查询

    // 敏感字段 - 使用permission标签限制
    CostPrice float64 `json:"cost_price" gorm:"column:cost_price;comment:成本价" 
                       runner:"code:cost_price;name:成本价"
                       search:"eq"
                       permission:"write"` // 只写权限，查询时自动禁止

    // 只读字段
    CreatedAt string `json:"created_at" gorm:"autoCreateTime" 
                      runner:"code:created_at;name:创建时间" 
                      search:"gte,lte"
                      permission:"read"` // 只读权限，支持时间范围查询
}
```

## 🔍 查询操作符详解

### 1. 等于查询 (eq)

```bash
# 单个条件
GET /api/products?eq=category:手机

# 多个条件（AND关系）
GET /api/products?eq=category:手机&eq=status:启用
```

```go
pageInfo := &query.PageInfoReq{
    Eq: []string{"category:手机", "status:启用"},
}
```

### 2. 模糊查询 (like)

```bash
# 产品名称包含"苹果"
GET /api/products?like=name:苹果

# 多个模糊条件
GET /api/products?like=name:苹果&like=tags:智能
```

```go
pageInfo := &query.PageInfoReq{
    Like: []string{"name:苹果", "tags:智能"},
}
```

### 3. 包含查询 (in)

```bash
# 分类为手机或平板
GET /api/products?in=category:手机&in=category:平板

# 状态为启用或禁用
GET /api/products?in=status:启用&in=status:禁用
```

```go
pageInfo := &query.PageInfoReq{
    In: []string{"category:手机", "category:平板"},
}
```

### 4. 数值比较查询

```bash
# 价格大于1000
GET /api/products?gt=price:1000

# 价格大于等于1000
GET /api/products?gte=price:1000

# 价格小于5000
GET /api/products?lt=price:5000

# 价格小于等于5000
GET /api/products?lte=price:5000

# 价格区间查询（1000-5000）
GET /api/products?gte=price:1000&lte=price:5000
```

```go
pageInfo := &query.PageInfoReq{
    Gte: []string{"price:1000"},
    Lte: []string{"price:5000"},
}
```

### 5. 否定查询操作符

#### 5.1 不等于查询 (not_eq)

```bash
# 分类不是手机
GET /api/products?not_eq=category:手机

# 状态不是禁用，且分类不是配件
GET /api/products?not_eq=status:禁用&not_eq=category:配件
```

```go
pageInfo := &query.PageInfoReq{
    NotEq: []string{"category:手机", "status:禁用"},
}
```

#### 5.2 不模糊匹配查询 (not_like)

```bash
# 产品名称不包含"测试"
GET /api/products?not_like=name:测试

# 描述不包含"临时"或"废弃"
GET /api/products?not_like=description:临时&not_like=description:废弃
```

```go
pageInfo := &query.PageInfoReq{
    NotLike: []string{"name:测试", "description:临时"},
}
```

#### 5.3 不包含查询 (not_in)

```bash
# 分类不是手机也不是平板
GET /api/products?not_in=category:手机&not_in=category:平板

# 状态不是禁用、删除、草稿
GET /api/products?not_in=status:禁用&not_in=status:删除&not_in=status:草稿
```

```go
pageInfo := &query.PageInfoReq{
    NotIn: []string{"category:手机", "category:平板", "status:禁用"},
}
```

## 📊 分页参数

### 基础分页

```bash
# 第1页，每页10条
GET /api/products?page=1&page_size=10

# 第2页，每页20条
GET /api/products?page=2&page_size=20
```

### 默认值

```go
// 如果不传分页参数，默认值为：
// page: 1
// page_size: 20
```

### 分页响应格式

```json
{
  "items": [...],           // 当前页数据
  "current_page": 1,        // 当前页码
  "total_count": 100,       // 总记录数
  "total_pages": 10,        // 总页数
  "page_size": 10          // 每页数量
}
```

## 🔄 排序功能

### 单字段排序

```bash
# 按价格升序
GET /api/products?sorts=price:asc

# 按价格降序
GET /api/products?sorts=price:desc

# 按创建时间降序
GET /api/products?sorts=created_at:desc
```

### 多字段排序

```bash
# 先按分类升序，再按价格降序
GET /api/products?sorts=category:asc,price:desc

# 先按状态升序，再按创建时间降序，最后按价格升序
GET /api/products?sorts=status:asc,created_at:desc,price:asc
```

```go
pageInfo := &query.PageInfoReq{
    Sorts: "category:asc,price:desc",
}
```

## 🛡️ 安全配置

### 1. 查询配置 (QueryConfig)

```go
// 创建查询配置
config := query.NewQueryConfig()

// 白名单：只允许指定字段的指定操作
config.AllowField("name", "like", "eq", "not_like", "not_eq")     // name字段支持正向和否定查询
config.AllowField("category", "eq", "in", "not_eq", "not_in")     // category字段支持精确匹配和否定查询
config.AllowField("price", "gte", "lte", "eq", "not_eq")          // price字段支持范围查询和否定查询
config.AllowField("tags", "like", "not_like", "in", "not_in")     // tags字段支持模糊匹配和否定查询
config.AllowField("status", "eq", "not_eq")                      // status字段支持精确匹配和否定查询

// 黑名单：禁止查询指定字段
config.DenyField("password")                      // 禁止查询password字段
config.DenyField("secret_key")                    // 禁止查询secret_key字段

// 使用配置
result, err := query.AutoPaginateTable(ctx, db, &Product{}, &products, pageInfo, config)
```

### 2. 多配置合并

```go
// 基础配置
baseConfig := query.NewQueryConfig()
baseConfig.AllowField("name", "like")
baseConfig.AllowField("category", "eq")

// 扩展配置
extConfig := query.NewQueryConfig()
extConfig.AllowField("price", "gte", "lte")
extConfig.DenyField("internal_code")

// 自动合并配置
result, err := query.AutoPaginateTable(ctx, db, &Product{}, &products, pageInfo, baseConfig, extConfig)
```

## 💡 实际应用示例

### 1. 电商产品列表（推荐用法）

```go
// 产品模型 - 通过search标签声明查询能力
type Product struct {
    ID        int     `json:"id" gorm:"primaryKey" runner:"code:id;name:产品ID"`
    Name      string  `json:"name" runner:"code:name;name:产品名称" search:"like,eq"`
    Category  string  `json:"category" runner:"code:category;name:产品分类" search:"eq,in"`
    Price     float64 `json:"price" runner:"code:price;name:产品价格" search:"gte,lte,eq"`
    Stock     int     `json:"stock" runner:"code:stock;name:库存数量" search:"gte,lte"`
    Status    string  `json:"status" runner:"code:status;name:产品状态" search:"eq"`
    CostPrice float64 `json:"cost_price" runner:"code:cost_price;name:成本价"` // 无search标签，自动禁止查询
    CreatedAt string  `json:"created_at" gorm:"autoCreateTime" runner:"code:created_at;name:创建时间" search:"gte,lte" permission:"read"`
}

func ProductList(ctx *gin.Context) {
    // 绑定查询参数
    var req struct {
        query.PageInfoReq `form:",inline"`
        // 可以添加额外的业务过滤条件
        MinPrice float64 `form:"min_price"`
        MaxPrice float64 `form:"max_price"`
    }
    
    if err := ctx.ShouldBindQuery(&req); err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 构建基础查询
    db := getDB().Model(&Product{})
    
    // 添加业务过滤条件
    if req.MinPrice > 0 {
        db = db.Where("price >= ?", req.MinPrice)
    }
    if req.MaxPrice > 0 {
        db = db.Where("price <= ?", req.MaxPrice)
    }
    
    // 零配置查询：自动根据Product模型的search标签进行安全验证和查询
    var products []Product
    result, err := query.AutoSearchPaginated(db, &Product{}, products, &req.PageInfoReq)
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(200, result)
}
```

### 2. 用户管理列表（推荐用法）

```go
// 用户模型 - 通过search标签声明查询能力
type User struct {
    ID        int    `json:"id" gorm:"primaryKey" runner:"code:id;name:用户ID"`
    Username  string `json:"username" runner:"code:username;name:用户名" search:"like,eq"`
    Email     string `json:"email" runner:"code:email;name:邮箱" search:"like,eq"`
    Role      string `json:"role" runner:"code:role;name:用户角色" search:"eq,in"`
    Status    string `json:"status" runner:"code:status;name:用户状态" search:"eq"`
    Age       int    `json:"age" runner:"code:age;name:年龄" search:"gte,lte,eq"`
    CreatedAt string `json:"created_at" gorm:"autoCreateTime" runner:"code:created_at;name:注册时间" search:"gte,lte" permission:"read"`
    Password  string `json:"-" gorm:"column:password" runner:"code:password;name:密码" permission:"write"` // 无search标签且仅写权限
}

func UserList(ctx *gin.Context) {
    var req struct {
        query.PageInfoReq `form:",inline"`
    }
    
    if err := ctx.ShouldBindQuery(&req); err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    db := getDB().Model(&User{})
    var users []User
    
    // 零配置查询：自动根据User模型的search标签进行安全验证和查询
    result, err := query.AutoSearchPaginated(db, &User{}, users, &req.PageInfoReq)
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(200, result)
}
```

## 🌐 HTTP请求示例

### 1. 基础分页

```bash
# 获取第一页，每页10条
curl "http://localhost:8080/api/products?page=1&page_size=10"
```

### 2. 复杂查询

```bash
# 查询手机分类、价格在1000-5000之间、名称包含"苹果"的产品
curl "http://localhost:8080/api/products?eq=category:手机&gte=price:1000&lte=price:5000&like=name:苹果&page=1&page_size=20&sorts=price:desc"
```

### 3. 否定查询示例

```bash
# 查询分类不是"测试"、名称不包含"废弃"、状态不是"禁用"的产品
curl "http://localhost:8080/api/products?not_eq=category:测试&not_like=name:废弃&not_eq=status:禁用"

# 查询分类不是手机和平板、标签不包含临时和测试的产品
curl "http://localhost:8080/api/products?not_in=category:手机&not_in=category:平板&not_in=tags:临时&not_in=tags:测试"

# 混合正向和否定查询：启用状态、价格大于100、分类不是配件、名称不包含测试
curl "http://localhost:8080/api/products?eq=status:启用&gte=price:100&not_eq=category:配件&not_like=name:测试"
```

### 4. 多条件查询

```bash
# 查询多个分类的启用状态产品，按分类和价格排序
curl "http://localhost:8080/api/products?in=category:手机&in=category:平板&eq=status:启用&sorts=category:asc,price:desc"
```

## 🔧 高级用法

### 1. 手动配置QueryConfig（不推荐）

如果你需要更精细的控制，仍然可以手动配置QueryConfig：

```go
func ProductListWithManualConfig(ctx *gin.Context) {
    var req struct {
        query.PageInfoReq `form:",inline"`
    }
    
    if err := ctx.ShouldBindQuery(&req); err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 手动创建安全配置
    config := query.NewQueryConfig()
    config.AllowField("name", "like")                    // 产品名称模糊搜索
    config.AllowField("category", "eq", "in")           // 分类精确匹配和多选
    config.AllowField("status", "eq")                   // 状态精确匹配
    config.AllowField("price", "gte", "lte", "eq")      // 价格范围查询
    config.AllowField("stock", "gte", "lte")            // 库存范围查询
    config.DenyField("cost_price")                      // 禁止查询成本价
    
    db := getDB().Model(&Product{})
    var products []Product
    
    // 使用手动配置
    result, err := query.AutoPaginateTable(ctx, db, &Product{}, &products, &req.PageInfoReq, config)
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(200, result)
}
```

### 2. 混合使用：基于模型自动配置 + 手动补充

```go
func buildQueryConfig(model interface{}, userRole string) *query.QueryConfig {
    // 首先从模型自动构建基础配置
    config, err := query.BuildQueryConfigFromModel(model)
    if err != nil {
        return query.NewQueryConfig()
    }
    
    // 根据用户角色动态调整权限
    switch userRole {
    case "admin":
        // 管理员可以查询所有字段（包括敏感字段）
        config.AllowField("cost_price", "gte", "lte", "eq")
        config.AllowField("created_by", "eq")
    case "manager":
        // 经理可以查询价格但不能查询成本
        config.DenyField("cost_price")
    default:
        // 普通用户不能查询价格相关信息
        config.DenyField("price")
        config.DenyField("cost_price")
    }
    
    return config
}
```

### 3. 自定义分页大小限制

```go
// 设置默认分页大小
pageInfo := &query.PageInfoReq{
    Page: 1,
    // PageSize 不设置，使用默认值
}

// 获取分页大小（如果未设置，使用默认值20）
limit := pageInfo.GetLimit()

// 获取分页大小（如果未设置，使用指定默认值）
limit := pageInfo.GetLimit(50) // 默认50条
```

### 4. 错误处理

```go
result, err := query.AutoSearchPaginated(db, &Product{}, products, pageInfo)
if err != nil {
    // 处理不同类型的错误
    if strings.Contains(err.Error(), "无效的字段名") {
        return fmt.Errorf("查询字段不合法: %w", err)
    }
    if strings.Contains(err.Error(), "字段不允许查询") {
        return fmt.Errorf("没有权限查询该字段: %w", err)
    }
    if strings.Contains(err.Error(), "不支持的操作符") {
        return fmt.Errorf("查询操作不支持: %w", err)
    }
    
    return fmt.Errorf("查询失败: %w", err)
}
```

### 5. 生成搜索表单配置

```go
// 自动生成前端搜索表单配置
func GetSearchFormConfig(ctx *gin.Context) {
    config, err := query.GenerateSearchFormConfig(&Product{})
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(200, config)
}
```

## ⚠️ 注意事项

### 1. 安全性

- ✅ **使用QueryConfig**：生产环境建议使用QueryConfig限制查询字段
- ✅ **字段验证**：所有字段名都会进行SQL注入检查
- ✅ **操作符限制**：通过白名单限制允许的查询操作符

### 2. 性能优化

- ✅ **索引优化**：确保查询字段有适当的数据库索引
- ✅ **分页大小**：建议限制分页大小，避免一次查询过多数据
- ✅ **复杂查询**：复杂查询条件可能影响性能，建议监控

### 3. 使用建议

- ✅ **优先使用search标签**：推荐使用`AutoSearchPaginated`结合`search`标签，零配置实现安全查询
- ✅ **模型驱动设计**：通过标签在模型层声明查询能力，避免重复配置
- ✅ **权限控制**：使用`permission`标签控制字段的查询权限
- ✅ **参数验证**：在业务层对查询参数进行额外验证
- ✅ **错误处理**：提供友好的错误信息给前端
- ✅ **日志记录**：记录查询日志便于调试和监控

## 📚 完整示例

查看 `pkg/query/query_test.go` 文件获取更多完整的使用示例和测试用例。

## 🔗 相关链接

- [GORM文档](https://gorm.io/docs/)
- [Gin框架文档](https://gin-gonic.com/docs/)
- [项目仓库](https://github.com/yunhanshu-net/pkg) 