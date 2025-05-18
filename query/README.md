# 通用查询参数使用文档

本文档详细说明了通用查询参数的使用方法，包括分页、排序和各种查询条件。

## 基础参数

### 分页参数
- `page`: 当前页码，从1开始
- `page_size`: 每页数量，默认20

示例：
```http
GET /api/v1/xxx?page=1&page_size=10
```

### 排序参数
- `sorts`: 排序字段和方向，格式为 `字段名:方向,字段名:方向`
- 方向可选值：`asc`（升序）或 `desc`（降序）

示例：
```http
GET /api/v1/xxx?sorts=created_at:desc,cost:asc
```

## 查询条件

### 等于查询 (eq)
格式：`eq=字段名:值`

示例：
```http
GET /api/v1/xxx?eq=status:1
GET /api/v1/xxx?eq=status:1&eq=type:test
```

### 模糊查询 (like)
格式：`like=字段名:值`

示例：
```http
GET /api/v1/xxx?like=name:张
GET /api/v1/xxx?like=name:张&like=email:example
```

### IN查询 (in)
格式：`in=字段名:值,字段名:值`

示例：
```http
# 单个字段多个值
GET /api/v1/xxx?in=status:1,status:2,status:3

# 多个字段多个值
GET /api/v1/xxx?in=status:1,status:2,role:admin,role:user
```

### 大于查询 (gt)
格式：`gt=字段名:值`

示例：
```http
GET /api/v1/xxx?gt=cost:100
GET /api/v1/xxx?gt=cost:100&gt=age:18
```

### 大于等于查询 (gte)
格式：`gte=字段名:值`

示例：
```http
GET /api/v1/xxx?gte=cost:100
GET /api/v1/xxx?gte=cost:100&gte=age:18
```

### 小于查询 (lt)
格式：`lt=字段名:值`

示例：
```http
GET /api/v1/xxx?lt=cost:100
GET /api/v1/xxx?lt=cost:100&lt=age:18
```

### 小于等于查询 (lte)
格式：`lte=字段名:值`

示例：
```http
GET /api/v1/xxx?lte=cost:100
GET /api/v1/xxx?lte=cost:100&lte=age:18
```

## 组合查询示例

### 基础分页和排序
```http
GET /api/v1/xxx?page=1&page_size=10&sorts=created_at:desc
```

### 多条件组合查询
```http
GET /api/v1/xxx?page=1&page_size=10&sorts=created_at:desc&eq=status:1&like=name:张&in=role:admin,role:user&gt=cost:100
```

### 复杂条件查询
```http
GET /api/v1/xxx?page=1&page_size=10&sorts=created_at:desc,cost:asc&eq=status:1&eq=type:test&like=name:张&like=email:example&in=role:admin,role:user&gt=cost:100&lt=age:30
```

## 注意事项

1. 所有字段名必须符合数据库字段命名规范（字母、数字、下划线）
2. 查询条件中的值会自动进行类型转换
3. 多个相同类型的条件会以 AND 关系组合
4. 不同类型的条件也会以 AND 关系组合
5. 如果使用了未配置的字段或操作符，会返回错误提示

## 错误码说明

- 400: 参数格式错误
- 403: 字段不允许查询
- 500: 服务器内部错误

## 返回格式

```json
{
    "items": [], // 查询结果列表
    "current_page": 1, // 当前页码
    "total_count": 100, // 总数据量
    "total_pages": 10, // 总页数
    "page_size": 10 // 每页数量
}
``` 