# Query库使用示例

这是一个通用的GORM查询库，支持动态搜索条件和分页功能。可以轻松集成到任何Go项目中。

## 核心特性

- 🔍 **丰富的搜索操作符**：支持 eq, like, in, gt, gte, lt, lte, not_eq, not_like, not_in
- 🛡️ **安全防护**：内置SQL注入防护和字段验证
- ⚡ **高性能**：优化的查询构建和执行
- 🔧 **易于集成**：简单的API设计，易于在现有项目中使用
- 📦 **零依赖**：只依赖GORM，无其他外部依赖

## 快速开始

### 1. 基本使用

```go
package main

import (
    "github.com/yunhanshu-net/pkg/query"
    "gorm.io/gorm"
)

type Product struct {
    ID       uint    `json:"id" gorm:"primaryKey"`
    Name     string  `json:"name"`
    Category string  `json:"category"`
    Price    float64 `json:"price"`
    Status   string  `json:"status"`
}

func GetProducts(db *gorm.DB, pageInfo *query.PageInfoReq) ([]Product, error) {
    var products []Product
    
    // 方法1：使用SimplePaginate（推荐）
    result, err := query.SimplePaginate(db, &Product{}, &products, pageInfo)
    if err != nil {
        return nil, err
    }
    
    // result包含分页信息和数据
    fmt.Printf("总数: %d, 当前页: %d\n", result.TotalCount, result.CurrentPage)
    
    return products, nil
}
```

### 2. 高级使用：自定义查询

```go
func GetProductsAdvanced(db *gorm.DB, pageInfo *query.PageInfoReq) ([]Product, error) {
    var products []Product
    
    // 方法2：使用ApplySearchConditions
    // 先构建基础查询
    baseQuery := db.Model(&Product{}).Where("deleted_at IS NULL")
    
    // 应用搜索条件
    queryWithConditions, err := query.ApplySearchConditions(baseQuery, pageInfo)
    if err != nil {
        return nil, err
    }
    
    // 执行查询
    err = queryWithConditions.Find(&products).Error
    return products, err
}
```

### 3. 带权限控制的查询

```go
func GetProductsWithPermission(db *gorm.DB, pageInfo *query.PageInfoReq, userID string) ([]Product, error) {
    var products []Product
    
    // 创建查询配置
    config := query.NewQueryConfig()
    
    // 允许的搜索字段和操作符
    config.AllowField("name", "like", "not_like")
    config.AllowField("category", "eq", "in", "not_eq", "not_in")
    config.AllowField("price", "gte", "lte", "eq", "gt", "lt")
    config.AllowField("status", "eq", "not_eq")
    
    // 禁止搜索敏感字段
    config.DenyField("created_by")
    config.DenyField("internal_notes")
    
    // 添加用户权限过滤
    baseQuery := db.Model(&Product{}).Where("created_by = ? OR status = 'public'", userID)
    
    // 应用搜索条件（带权限控制）
    queryWithConditions, err := query.ApplySearchConditions(baseQuery, pageInfo, config)
    if err != nil {
        return nil, err
    }
    
    err = queryWithConditions.Find(&products).Error
    return products, err
}
```

## HTTP API集成示例

### 1. Gin框架集成

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/yunhanshu-net/pkg/query"
)

func ProductListHandler(c *gin.Context) {
    // 绑定查询参数
    var pageInfo query.PageInfoReq
    if err := c.ShouldBindQuery(&pageInfo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 查询数据
    var products []Product
    result, err := query.SimplePaginate(db, &Product{}, &products, &pageInfo)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // 返回结果
    c.JSON(http.StatusOK, result)
}
```

### 2. HTTP请求示例

```bash
# 基础分页
GET /products?page=1&page_size=10

# 精确匹配
GET /products?eq=category:手机&eq=status:启用

# 模糊搜索
GET /products?like=name:iPhone

# 范围查询
GET /products?gte=price:1000&lte=price:5000

# 包含查询
GET /products?in=category:手机&in=category:平板

# 否定查询
GET /products?not_eq=status:禁用&not_like=name:测试

# 排序
GET /products?sorts=price:DESC,created_at:ASC

# 复合查询
GET /products?like=name:iPhone&eq=status:启用&gte=price:5000&sorts=price:DESC
```

## 搜索操作符详解

| 操作符 | 说明 | 示例 | SQL等价 | 适用字段类型 |
|--------|------|------|---------|-------------|
| `eq` | 精确匹配 | `eq=status:启用` | `status = '启用'` | 字符串、数字、布尔值 |
| `like` | 模糊匹配 | `like=name:iPhone` | `name LIKE '%iPhone%'` | 字符串 |
| `in` | 包含查询 | `in=category:手机&in=category:平板` | `category IN ('手机', '平板')` | 字符串、数字 |
| `gt` | 大于 | `gt=price:1000` | `price > 1000` | 数字、时间戳 |
| `gte` | 大于等于 | `gte=price:1000` | `price >= 1000` | 数字、时间戳 |
| `lt` | 小于 | `lt=price:5000` | `price < 5000` | 数字、时间戳 |
| `lte` | 小于等于 | `lte=price:5000` | `price <= 5000` | 数字、时间戳 |
| `not_eq` | 不等于 | `not_eq=status:禁用` | `status != '禁用'` | 字符串、数字、布尔值 |
| `not_like` | 否定模糊匹配 | `not_like=name:测试` | `name NOT LIKE '%测试%'` | 字符串 |
| `not_in` | 否定包含查询 | `not_in=category:其他` | `category NOT IN ('其他')` | 字符串、数字 |

## 搜索操作符使用指南

### 字符串字段搜索
```bash
# 精确匹配
GET /products?eq=status:启用

# 模糊搜索
GET /products?like=name:iPhone

# 包含查询
GET /products?in=category:手机&in=category:平板

# 否定查询
GET /products?not_eq=status:禁用
GET /products?not_like=name:测试
GET /products?not_in=category:其他
```

### 数值字段搜索
```bash
# 精确匹配
GET /products?eq=price:1000

# 范围查询
GET /products?gt=price:1000      # 大于1000
GET /products?gte=price:1000     # 大于等于1000
GET /products?lt=price:5000      # 小于5000
GET /products?lte=price:5000     # 小于等于5000

# 否定查询
GET /products?not_eq=price:1000  # 不等于1000
```

### 时间戳字段搜索
```bash
# 时间范围查询（毫秒时间戳）
GET /products?gte=created_at:1704067200000  # 2024-01-01 00:00:00 之后
GET /products?lte=created_at:1704153600000  # 2024-01-02 00:00:00 之前

# 精确时间查询
GET /products?eq=created_at:1704067200000   # 2024-01-01 00:00:00
```

## 排序功能

```bash
# 单字段排序
GET /products?sorts=price:DESC

# 多字段排序
GET /products?sorts=price:DESC,created_at:ASC

# 支持的排序方向：ASC（升序）、DESC（降序）
```

## 安全特性

### 1. SQL注入防护
所有字段名都会通过 `SafeColumn()` 函数验证，只允许字母、数字和下划线。

### 2. 字段权限控制
```go
config := query.NewQueryConfig()
config.AllowField("name", "like")      // 只允许name字段进行like查询
config.DenyField("password")           // 禁止查询password字段
```

### 3. 查询配置示例
```go
// 创建严格的查询配置
config := query.NewQueryConfig()

// 产品名称：只允许模糊搜索
config.AllowField("name", "like", "not_like")

// 分类：允许精确匹配和包含查询
config.AllowField("category", "eq", "in", "not_eq", "not_in")

// 价格：允许范围查询
config.AllowField("price", "gte", "lte", "eq")

// 禁止查询敏感字段
config.DenyField("created_by")
config.DenyField("internal_notes")
config.DenyField("cost_price")
```

## 错误处理

```go
result, err := query.SimplePaginate(db, &Product{}, &products, pageInfo)
if err != nil {
    // 处理常见错误
    switch {
    case strings.Contains(err.Error(), "无效的字段名"):
        // 字段名不安全
        return fmt.Errorf("查询参数包含非法字段")
    case strings.Contains(err.Error(), "不允许查询"):
        // 字段被禁止查询
        return fmt.Errorf("无权限查询该字段")
    case strings.Contains(err.Error(), "不支持的操作符"):
        // 操作符不被允许
        return fmt.Errorf("该字段不支持此查询操作")
    default:
        // 其他数据库错误
        return fmt.Errorf("查询失败: %w", err)
    }
}
```

## 性能优化建议

### 1. 数据库索引
确保经常搜索的字段有适当的索引：

```sql
-- 为常用搜索字段创建索引
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_products_created_at ON products(created_at);

-- 复合索引用于多字段查询
CREATE INDEX idx_products_category_status ON products(category, status);
```

### 2. 分页大小限制
```go
func ValidatePageSize(pageInfo *query.PageInfoReq) {
    if pageInfo.PageSize > 100 {
        pageInfo.PageSize = 100  // 限制最大分页大小
    }
    if pageInfo.PageSize <= 0 {
        pageInfo.PageSize = 20   // 设置默认分页大小
    }
}
```

### 3. 查询缓存
对于频繁查询的数据，可以考虑添加缓存：

```go
func GetProductsWithCache(db *gorm.DB, pageInfo *query.PageInfoReq) ([]Product, error) {
    // 生成缓存键
    cacheKey := fmt.Sprintf("products:%s", generateCacheKey(pageInfo))
    
    // 尝试从缓存获取
    if cached := getFromCache(cacheKey); cached != nil {
        return cached.([]Product), nil
    }
    
    // 查询数据库
    var products []Product
    result, err := query.SimplePaginate(db, &Product{}, &products, pageInfo)
    if err != nil {
        return nil, err
    }
    
    // 存入缓存
    setCache(cacheKey, products, 5*time.Minute)
    
    return products, nil
}
```

## 集成到现有项目

### 1. 最小化改动
如果你已有分页功能，只需要替换查询构建部分：

```go
// 原有代码
// query := db.Model(&Product{})
// if name != "" {
//     query = query.Where("name LIKE ?", "%"+name+"%")
// }

// 新代码
query, err := query.ApplySearchConditions(db.Model(&Product{}), pageInfo)
if err != nil {
    return err
}
```

### 2. 渐进式迁移
可以先在新功能中使用，然后逐步迁移旧功能：

```go
func GetProducts(db *gorm.DB, useNewQuery bool, pageInfo *query.PageInfoReq) ([]Product, error) {
    var products []Product
    
    if useNewQuery {
        // 使用新的query库
        result, err := query.SimplePaginate(db, &Product{}, &products, pageInfo)
        return products, err
    } else {
        // 保持原有逻辑
        return getProductsLegacy(db, pageInfo)
    }
}
```

## 总结

这个query库提供了：

1. **简单易用**的API设计
2. **安全可靠**的查询构建
3. **灵活强大**的搜索功能
4. **高性能**的分页查询
5. **易于集成**的架构设计

无论是新项目还是现有项目，都可以轻松集成这个库来实现强大的搜索和分页功能。