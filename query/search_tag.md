# Search标签设计文档

## 🎯 设计目标

基于现有的标签体系（`runner`、`widget`、`data`、`validate`、`permission`），新增`search`标签来定义字段的查询能力，实现：

1. **声明式查询配置**：通过标签声明字段支持的查询操作符
2. **自动查询验证**：根据search标签自动验证查询请求
3. **智能UI生成**：根据search标签和widget标签自动生成搜索界面
4. **零配置使用**：开发者只需要添加标签，无需手写查询逻辑

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

## 📝 使用示例

### 1. 完整的产品模型示例

```go
type Product struct {
    // 基础字段 - 不可搜索
    ID int `json:"id" gorm:"primaryKey" 
           runner:"code:id;name:产品ID" 
           data:"type:number"`

    // 文本搜索字段
    Name string `json:"name" gorm:"column:name;comment:产品名称" 
                 runner:"code:name;name:产品名称" 
                 widget:"type:input;placeholder:请输入产品名称" 
                 data:"type:string" 
                 search:"like,eq"
                 validate:"required"`

    // 分类选择字段 - 支持精确匹配和多选
    Category string `json:"category" gorm:"column:category;comment:产品分类" 
                     runner:"code:category;name:产品分类" 
                     widget:"type:select;options:手机,笔记本,平板,耳机,其他" 
                     data:"type:string" 
                     search:"eq,in"
                     validate:"required"`

    // 价格范围搜索
    Price float64 `json:"price" gorm:"column:price;comment:产品价格" 
                   runner:"code:price;name:产品价格" 
                   widget:"type:input;prefix:¥;precision:2" 
                   data:"type:float" 
                   search:"gte,lte,eq"
                   validate:"required,min=0"`

    // 库存数量搜索
    Stock int `json:"stock" gorm:"column:stock;comment:库存数量" 
               runner:"code:stock;name:库存数量" 
               widget:"type:input;suffix:件" 
               data:"type:number" 
               search:"gte,lte,eq,gt,lt"
               validate:"required,min=0"`

    // 状态开关搜索
    Status string `json:"status" gorm:"column:status;comment:产品状态" 
                   runner:"code:status;name:产品状态" 
                   widget:"type:switch;true_value:启用;false_value:禁用" 
                   data:"type:string;default_value:启用" 
                   search:"eq"
                   validate:"required"`

    // 标签模糊搜索
    Tags string `json:"tags" gorm:"column:tags;comment:产品标签" 
                 runner:"code:tags;name:产品标签" 
                 widget:"type:tag;separator:,;max_tags:5" 
                 data:"type:string" 
                 search:"like,in"`

    // 只读字段 - 不可搜索
    CreatedAt typex.Time `json:"created_at" gorm:"autoCreateTime" 
                         runner:"code:created_at;name:创建时间" 
                         widget:"type:datetime;format:datetime" 
                         data:"type:string" 
                         permission:"read"`
}
```

### 2. 用户模型示例

```go
type User struct {
    ID int `json:"id" gorm:"primaryKey" 
           runner:"code:id;name:用户ID" 
           data:"type:number"`

    // 用户名模糊搜索
    Username string `json:"username" gorm:"column:username;comment:用户名" 
                     runner:"code:username;name:用户名" 
                     widget:"type:input;placeholder:请输入用户名" 
                     data:"type:string" 
                     search:"like,eq"
                     validate:"required"`

    // 邮箱模糊搜索
    Email string `json:"email" gorm:"column:email;comment:邮箱" 
                  runner:"code:email;name:邮箱" 
                  widget:"type:input;placeholder:请输入邮箱" 
                  data:"type:string" 
                  search:"like,eq"
                  validate:"required,email"`

    // 角色多选搜索
    Role string `json:"role" gorm:"column:role;comment:用户角色" 
                 runner:"code:role;name:用户角色" 
                 widget:"type:select;options:admin,user,guest" 
                 data:"type:string;default_value:user" 
                 search:"eq,in"
                 validate:"required"`

    // 年龄范围搜索
    Age int `json:"age" gorm:"column:age;comment:年龄" 
            runner:"code:age;name:年龄" 
            widget:"type:input;suffix:岁" 
            data:"type:number" 
            search:"gte,lte,eq"`

    // 注册时间范围搜索
    CreatedAt typex.Time `json:"created_at" gorm:"autoCreateTime" 
                         runner:"code:created_at;name:注册时间" 
                         widget:"type:datetime;format:datetime" 
                         data:"type:string" 
                         search:"gte,lte"
                         permission:"read"`

    // 密码字段 - 不可搜索
    Password string `json:"-" gorm:"column:password;comment:密码" 
                     runner:"code:password;name:密码" 
                     widget:"type:input;mode:password" 
                     data:"type:string" 
                     permission:"write"
                     validate:"required"`
}
```

## 🔄 智能UI组件映射

基于search标签和widget标签的组合，自动生成对应的搜索组件：

### 1. 文本搜索组件

```go
// 配置：search:"like" + widget:"type:input"
// 生成：文本输入框搜索
Name string `search:"like" widget:"type:input;placeholder:搜索产品名称"`
```

**前端渲染：**
```html
<input type="text" placeholder="搜索产品名称" name="like" />
```

### 2. 下拉选择搜索

```go
// 配置：search:"eq" + widget:"type:select"
// 生成：下拉选择搜索
Category string `search:"eq" widget:"type:select;options:手机,笔记本,平板"`
```

**前端渲染：**
```html
<select name="eq">
  <option value="">请选择分类</option>
  <option value="手机">手机</option>
  <option value="笔记本">笔记本</option>
  <option value="平板">平板</option>
</select>
```

### 3. 多选搜索

```go
// 配置：search:"in" + widget:"type:select"
// 生成：多选下拉搜索
Category string `search:"in" widget:"type:select;options:手机,笔记本,平板"`
```

**前端渲染：**
```html
<select name="in" multiple>
  <option value="手机">手机</option>
  <option value="笔记本">笔记本</option>
  <option value="平板">平板</option>
</select>
```

### 4. 数值范围搜索

```go
// 配置：search:"gte,lte" + widget:"type:input"
// 生成：范围输入组件
Price float64 `search:"gte,lte" widget:"type:input;prefix:¥"`
```

**前端渲染：**
```html
<div class="range-input">
  <input type="number" placeholder="最低价格" name="gte" />
  <span>-</span>
  <input type="number" placeholder="最高价格" name="lte" />
</div>
```

### 5. 开关搜索

```go
// 配置：search:"eq" + widget:"type:switch"
// 生成：开关选择搜索
Status string `search:"eq" widget:"type:switch;true_value:启用;false_value:禁用"`
```

**前端渲染：**
```html
<select name="eq">
  <option value="">全部状态</option>
  <option value="启用">启用</option>
  <option value="禁用">禁用</option>
</select>
```

### 6. 日期范围搜索

```go
// 配置：search:"gte,lte" + widget:"type:datetime"
// 生成：日期范围选择器
CreatedAt typex.Time `search:"gte,lte" widget:"type:datetime;format:date"`
```

**前端渲染：**
```html
<div class="date-range">
  <input type="date" placeholder="开始日期" name="gte" />
  <span>-</span>
  <input type="date" placeholder="结束日期" name="lte" />
</div>
```

## 🔧 自动化功能实现

### 1. 自动生成QueryConfig

```go
// 根据search标签自动生成查询配置
func BuildQueryConfigFromModel(model interface{}) (*query.QueryConfig, error) {
    config := query.NewQueryConfig()
    
    // 解析模型的search标签
    // 自动构建白名单配置
    // 返回安全的查询配置
    
    return config, nil
}
```

### 2. 自动验证查询请求

```go
// 根据search标签自动验证查询参数
func ValidateSearchRequest(model interface{}, pageInfo *query.PageInfoReq) error {
    // 检查查询字段是否有search标签
    // 检查操作符是否在允许列表中
    // 返回验证结果
    
    return nil
}
```

### 3. 自动生成搜索表单配置

```go
// 根据search标签和widget标签生成前端搜索表单配置
func GenerateSearchFormConfig(model interface{}) (*SearchFormConfig, error) {
    // 解析search标签和widget标签
    // 生成前端表单配置JSON
    // 支持各种组件类型的智能映射
    
    return config, nil
}
```

## 📊 使用流程

### 1. 模型定义阶段

```go
// 开发者只需要在模型字段上添加search标签
type Product struct {
    Name     string `search:"like,eq" widget:"type:input"`
    Category string `search:"eq,in" widget:"type:select;options:手机,笔记本"`
    Price    float64 `search:"gte,lte" widget:"type:input"`
}
```

### 2. API实现阶段

```go
// 使用增强的查询函数，自动处理search标签
func ProductList(ctx *runner.Context, req *ProductListReq, resp response.Response) error {
    db := ctx.MustGetOrInitDB()
    var results []Product
    
    // 自动验证和查询，基于search标签
    return resp.Table(&results).AutoSearchPaginated(
        db.Model(&Product{}),
        &Product{},
        &req.PageInfoReq,
    ).Build()
}
```

### 3. 前端使用阶段

```javascript
// 获取搜索表单配置
const searchConfig = await api.get('/api/product/search-config');

// 根据配置自动渲染搜索表单
renderSearchForm(searchConfig);

// 提交搜索请求
const searchParams = {
    like: 'name:苹果',
    eq: 'category:手机',
    gte: 'price:1000',
    lte: 'price:5000'
};
const results = await api.get('/api/product/list', searchParams);
```

## 🎯 预期效果

### 开发效率提升
- **模型定义**：只需添加search标签，无需手写查询逻辑
- **API开发**：一行代码实现复杂搜索功能
- **前端开发**：自动生成搜索表单，无需手写UI

### 代码质量提升
- **类型安全**：编译时检查search标签配置
- **安全防护**：自动SQL注入防护和权限控制
- **一致性**：统一的搜索体验和API格式

### 维护成本降低
- **配置集中**：搜索配置与模型定义在一起
- **自动同步**：模型变更自动同步到搜索功能
- **文档自生成**：根据标签自动生成API文档

## 🔄 与现有标签的协同

Search标签与现有标签体系完美配合：

- **runner标签**：提供字段基础信息（code、name）
- **widget标签**：提供UI组件配置，用于生成搜索界面
- **data标签**：提供数据类型信息，用于查询参数验证
- **validate标签**：提供验证规则，确保查询参数合法
- **permission标签**：提供权限控制，限制查询字段访问
- **search标签**：定义查询操作符，实现声明式查询配置

这样的设计既保持了标签体系的一致性，又提供了强大的搜索功能。 