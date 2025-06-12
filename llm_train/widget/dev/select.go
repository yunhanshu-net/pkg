package dev

// Select组件设计文档
// 下拉选择相关的组件，支持单选和多选

// ===== Select组件定义 =====

// SelectWidget 下拉选择组件定义
type SelectWidget struct {
	Options      string `json:"options"`       // 选项列表：value1(label1),value2(label2)
	DefaultValue string `json:"default_value"` // 默认值，必须在options中
	Multiple     bool   `json:"multiple"`      // 是否多选（计划功能）
	Searchable   bool   `json:"searchable"`    // 是否可搜索（计划功能）
}

// 使用示例：
// Status string `runner:"code:status;name:状态;widget:select;options:active(启用),inactive(禁用);default_value:active(启用)" validate:"required,oneof=active(启用) inactive(禁用)"`

// ===== Select组件设计 =====

/*
Select组件用于处理下拉选择，支持单选和多选模式。

设计目标：
1. 支持单选和多选模式
2. 提供搜索功能
3. 支持分组选项
4. 与现有验证规则无缝集成
5. 提供良好的用户体验

使用场景：
- 状态选择：启用/禁用、激活/停用等
- 分类选择：类型、优先级、部门等
- 多选场景：标签、权限、功能等
*/

// ===== 组件示例 =====

// SelectExample 下拉选择组件示例
type SelectExample struct {
	// 基础单选下拉
	Status string `runner:"code:status;name:状态;widget:select;options:active(启用),inactive(禁用);default_value:active(启用)" validate:"required,oneof=active(启用) inactive(禁用)" json:"status"`
	// 注释：基础单选，两个选项，默认启用，必填验证

	// 优先级选择
	Priority string `runner:"code:priority;name:优先级;widget:select;options:low(低),medium(中),high(高);default_value:medium(中)" validate:"required,oneof=low(低) medium(中) high(高)" json:"priority"`
	// 注释：三级优先级选择，默认中等优先级

	// 用户类型选择
	UserType string `runner:"code:user_type;name:用户类型;widget:select;options:normal(普通用户),vip(VIP用户),admin(管理员);default_value:normal(普通用户)" validate:"required,oneof=normal(普通用户) vip(VIP用户) admin(管理员)" json:"user_type"`
	// 注释：用户类型选择，默认普通用户

	// 分类选择
	Category string `runner:"code:category;name:分类;widget:select;options:tech(技术),product(产品),design(设计),marketing(市场)" validate:"required,oneof=tech(技术) product(产品) design(设计) marketing(市场)" json:"category"`
	// 注释：分类选择，无默认值，必填

	// 可选的选择器
	Region string `runner:"code:region;name:地区;widget:select;options:north(北方),south(南方),east(东部),west(西部)" validate:"omitempty,oneof=north(北方) south(南方) east(东部) west(西部)" json:"region"`
	// 注释：地区选择，可选字段

	// 多选示例（计划功能）
	Tags string `runner:"code:tags;name:标签;widget:select;options:tag1(标签1),tag2(标签2),tag3(标签3);multiple:true" json:"tags"`
	// 注释：多选标签，返回逗号分隔的值

	// 可搜索选择器（计划功能）
	City string `runner:"code:city;name:城市;widget:select;options:beijing(北京),shanghai(上海),guangzhou(广州),shenzhen(深圳);searchable:true" validate:"required" json:"city"`
	// 注释：可搜索的城市选择器
}

// ===== 标签配置详解 =====

/*
Select组件支持的标签：

核心标签：
- code: 字段代码（必需）
- name: 显示名称（必需）
- widget: select（必需）

选项配置：
- options: 选项列表，格式：value1(label1),value2(label2)
- default_value: 默认值，必须在options中存在

功能配置：
- multiple: 是否多选，默认false
- searchable: 是否可搜索，默认false
- placeholder: 占位符文本

显示控制：
- show: 显示场景控制
- hidden: 隐藏场景控制

选项格式说明：
- 基本格式：value(label)
- 多个选项用逗号分隔
- value是实际存储的值
- label是显示给用户的文本
- 如果value和label相同，可以简写为：value

验证规则：
- 单选：使用oneof验证，值必须在options的value中
- 多选：返回逗号分隔的value字符串
*/

// ===== 实现要点 =====

/*
前端实现要点：

1. 选项解析：
   - 解析options配置为选项列表
   - 支持value(label)格式
   - 处理默认值设置

2. 单选模式：
   - 下拉选择器实现
   - 选中状态管理
   - 值变化处理

3. 多选模式（计划）：
   - 多选下拉实现
   - 选中项显示
   - 值的组合和分离

4. 搜索功能（计划）：
   - 选项过滤
   - 高亮匹配文本
   - 键盘导航

5. 用户体验：
   - 清晰的选项显示
   - 选中状态指示
   - 加载状态处理

后端实现要点：

1. 选项验证：
   - 验证选中值在options中
   - 多选时验证每个值
   - 提供友好的错误信息

2. 数据处理：
   - 单选值处理
   - 多选值的分割和组合
   - 默认值设置

3. 存储优化：
   - 枚举值索引
   - 查询优化
   - 统计支持
*/

// ===== 验证规则集成 =====

/*
支持的验证规则：

1. 基础验证：
   - required: 必填验证
   - omitempty: 可选字段

2. 枚举验证：
   - oneof: 值必须在指定列表中
   - 自动根据options生成验证规则

3. 多选验证（计划）：
   - 验证每个选中值
   - 最少/最多选择数量限制

示例验证规则：
```go
// 单选验证
Status string `validate:"required,oneof=active inactive"`

// 可选单选
Region string `validate:"omitempty,oneof=north south east west"`

// 多选验证（计划）
Tags string `validate:"required,min=1"`
```
*/

// ===== 使用最佳实践 =====

/*
最佳实践建议：

1. 选项设计：
   - 选项数量适中（建议不超过20个）
   - 选项文本简洁明确
   - 考虑选项的逻辑分组

2. 默认值设置：
   - 为常用场景设置合理默认值
   - 默认值必须在options中存在
   - 考虑业务逻辑的默认状态

3. 验证规则：
   - 必须使用oneof验证
   - 验证值与options的value保持一致
   - 考虑大小写敏感性

4. 用户体验：
   - 选项过多时考虑搜索功能
   - 提供清晰的选项分组
   - 支持键盘操作

5. 数据一致性：
   - 前后端选项保持同步
   - 避免硬编码选项值
   - 考虑选项的动态更新

示例：
```go
// 推荐的配置
type GoodExample struct {
    // 状态选择
    Status string `runner:"code:status;name:状态;widget:select;options:active(启用),inactive(禁用);default_value:active(启用)"`

    // 优先级选择
    Priority string `runner:"code:priority;name:优先级;widget:select;options:low(低),medium(中),high(高);default_value:medium(中)"`

    // 分类选择
    Category string `runner:"code:category;name:分类;widget:select;options:tech(技术),product(产品),design(设计)"`
}

// 不推荐的配置
type BadExample struct {
    // 选项过多
    Country string `runner:"code:country;name:国家;widget:select;options:cn(中国),us(美国),...50个选项"` // 应该使用搜索功能

    // 缺少验证
    Status string `runner:"code:status;name:状态;widget:select;options:active,inactive"` // 应该添加oneof验证

    // 默认值不在选项中
    Type string `runner:"code:type;name:类型;widget:select;options:a,b,c;default_value:d"` // 默认值必须在options中
}
```
*/

// ===== 实现优先级 =====

/*
实现步骤：

第一步：基础单选功能
- 选项解析和渲染
- 单选值处理
- 默认值设置

第二步：验证集成
- oneof验证支持
- 错误提示优化
- 边界情况处理

第三步：用户体验优化
- 样式美化
- 交互优化
- 加载状态

第四步：高级功能（计划）
- 多选支持
- 搜索功能
- 分组选项

预计总工期：2-3天
- 第一步：1天
- 第二步：0.5天
- 第三步：0.5天
- 第四步：1天
*/
