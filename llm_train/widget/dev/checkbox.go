package dev

// Checkbox组件设计文档
// 复选框相关的组件，支持多选场景

// ===== Checkbox组件定义 =====

// CheckboxWidget 复选框组件定义
type CheckboxWidget struct {
	Options      string `json:"options"`       // 选项列表：value1(label1),value2(label2)
	DefaultValue string `json:"default_value"` // 默认选中的值，多个用逗号分隔
	MinSelect    int    `json:"min_select"`    // 最少选择数量
	MaxSelect    int    `json:"max_select"`    // 最多选择数量
}

// 使用示例：
// Permissions []string `runner:"code:permissions;name:权限;widget:checkbox;options:read(读取),write(写入),delete(删除);default_value:read,write" validate:"required,min=1"`
// Tags string `runner:"code:tags;name:标签;widget:checkbox;options:tag1(标签1),tag2(标签2),tag3(标签3);min_select:1;max_select:2"`

// ===== Checkbox组件设计 =====

/*
Checkbox组件用于处理多选场景，允许用户选择多个选项。

设计目标：
1. 支持多选功能
2. 提供选择数量限制
3. 支持默认选中项
4. 与现有验证规则无缝集成
5. 提供良好的用户体验

使用场景：
- 权限选择：读取、写入、删除等
- 标签选择：多个标签组合
- 功能开关：多个功能的启用/禁用
- 兴趣爱好：多项兴趣选择
*/

// ===== 组件示例 =====

// CheckboxExample 复选框组件示例
type CheckboxExample struct {
	// 基础多选复选框
	Permissions []string `runner:"code:permissions;name:权限;widget:checkbox;options:read(读取),write(写入),delete(删除);default_value:read,write" validate:"required,min=1" json:"permissions"`
	// 注释：权限多选，默认选中读取和写入，至少选择1项

	// 兴趣爱好选择
	Hobbies []string `runner:"code:hobbies;name:兴趣爱好;widget:checkbox;options:reading(阅读),music(音乐),sports(运动),travel(旅行),cooking(烹饪)" validate:"omitempty" json:"hobbies"`
	// 注释：兴趣爱好多选，可选字段，无默认值

	// 功能开关
	Features []string `runner:"code:features;name:功能;widget:checkbox;options:notification(通知),backup(备份),sync(同步),analytics(分析);default_value:notification,backup" json:"features"`
	// 注释：功能开关，默认开启通知和备份

	// 有数量限制的选择
	Skills []string `runner:"code:skills;name:技能;widget:checkbox;options:java(Java),python(Python),go(Go),javascript(JavaScript),react(React);min_select:1;max_select:3" validate:"required,min=1,max=3" json:"skills"`
	// 注释：技能选择，最少1项，最多3项

	// 标签选择
	Tags []string `runner:"code:tags;name:标签;widget:checkbox;options:urgent(紧急),important(重要),bug(缺陷),feature(功能),improvement(改进)" json:"tags"`
	// 注释：标签选择，可选，无数量限制

	// 工作日选择
	Workdays []string `runner:"code:workdays;name:工作日;widget:checkbox;options:monday(周一),tuesday(周二),wednesday(周三),thursday(周四),friday(周五);default_value:monday,tuesday,wednesday,thursday,friday" validate:"required,min=1" json:"workdays"`
	// 注释：工作日选择，默认全选工作日

	// 通知类型
	NotificationTypes []string `runner:"code:notification_types;name:通知类型;widget:checkbox;options:email(邮件),sms(短信),push(推送),wechat(微信);default_value:email,push" json:"notification_types"`
	// 注释：通知类型选择，默认邮件和推送
}

// ===== 标签配置详解 =====

/*
Checkbox组件支持的标签：

核心标签：
- code: 字段代码（必需）
- name: 显示名称（必需）
- widget: checkbox（必需）

选项配置：
- options: 选项列表，格式：value1(label1),value2(label2)
- default_value: 默认选中值，多个用逗号分隔

数量限制：
- min_select: 最少选择数量
- max_select: 最多选择数量

显示配置：
- direction: 排列方向，horizontal(水平)/vertical(垂直)，默认vertical
- columns: 列数，用于网格布局

显示控制：
- show: 显示场景控制
- hidden: 隐藏场景控制

选项格式说明：
- 基本格式：value(label)
- 多个选项用逗号分隔
- value是实际存储的值
- label是显示给用户的文本
- 如果value和label相同，可以简写为：value

数据格式：
- 返回值：[]string 或 string（逗号分隔）
- 空值处理：未选择时返回空数组或空字符串
*/

// ===== 实现要点 =====

/*
前端实现要点：

1. 选项渲染：
   - 解析options配置为复选框列表
   - 支持value(label)格式
   - 处理默认选中状态

2. 多选逻辑：
   - 选中状态管理
   - 全选/反选功能
   - 选择数量统计

3. 数量限制：
   - 最少选择验证
   - 最多选择限制
   - 实时提示信息

4. 布局控制：
   - 垂直/水平排列
   - 网格布局支持
   - 响应式设计

5. 用户体验：
   - 清晰的选中状态
   - 禁用状态处理
   - 键盘操作支持

后端实现要点：

1. 数据验证：
   - 验证选中值在options中
   - 数量限制验证
   - 数据类型验证

2. 数据处理：
   - 数组和字符串格式转换
   - 去重处理
   - 排序处理

3. 存储优化：
   - 多值字段存储
   - 索引优化
   - 查询支持
*/

// ===== 验证规则集成 =====

/*
支持的验证规则：

1. 基础验证：
   - required: 必填验证（至少选择一项）
   - omitempty: 可选字段

2. 数量验证：
   - min: 最少选择数量
   - max: 最多选择数量
   - len: 固定选择数量

3. 内容验证：
   - dive: 验证数组中每个元素
   - oneof: 每个值必须在指定列表中

示例验证规则：
```go
// 基础多选验证
Permissions []string `validate:"required,min=1,dive,oneof=read write delete"`

// 数量限制验证
Skills []string `validate:"required,min=1,max=3,dive,oneof=java python go"`

// 可选多选
Tags []string `validate:"omitempty,dive,oneof=urgent important bug"`
```
*/

// ===== 使用最佳实践 =====

/*
最佳实践建议：

1. 选项设计：
   - 选项数量适中（建议不超过15个）
   - 选项文本简洁明确
   - 考虑选项的逻辑分组

2. 默认值设置：
   - 为常用场景设置合理默认值
   - 默认值必须在options中存在
   - 考虑用户习惯和业务逻辑

3. 数量限制：
   - 根据业务需求设置合理限制
   - 提供清晰的数量提示
   - 考虑极端情况处理

4. 布局设计：
   - 选项较少时使用水平排列
   - 选项较多时使用垂直排列或网格
   - 考虑移动端适配

5. 数据处理：
   - 统一数据格式（推荐使用[]string）
   - 处理空值情况
   - 考虑数据迁移兼容性

示例：
```go
// 推荐的配置
type GoodExample struct {
    // 权限选择
    Permissions []string `runner:"code:permissions;name:权限;widget:checkbox;options:read(读取),write(写入),delete(删除);default_value:read"`

    // 技能选择（有限制）
    Skills []string `runner:"code:skills;name:技能;widget:checkbox;options:java(Java),python(Python),go(Go);min_select:1;max_select:2"`

    // 标签选择（可选）
    Tags []string `runner:"code:tags;name:标签;widget:checkbox;options:urgent(紧急),important(重要)"`
}

// 不推荐的配置
type BadExample struct {
    // 选项过多
    Countries []string `runner:"code:countries;name:国家;widget:checkbox;options:cn(中国),us(美国),...50个选项"` // 应该使用其他组件

    // 缺少验证
    Permissions []string `runner:"code:permissions;name:权限;widget:checkbox;options:read,write,delete"` // 应该添加验证规则

    // 不合理的限制
    Skills []string `runner:"code:skills;name:技能;widget:checkbox;options:java,python;min_select:3"` // 最少选择数不能超过选项总数
}
```
*/

// ===== 实现优先级 =====

/*
实现步骤：

第一步：基础多选功能
- 复选框渲染
- 多选状态管理
- 默认值设置

第二步：数量限制
- 最少/最多选择限制
- 实时验证提示
- 禁用状态处理

第三步：布局优化
- 排列方向控制
- 网格布局支持
- 样式美化

第四步：高级功能
- 全选/反选
- 搜索过滤
- 分组显示

预计总工期：2-3天
- 第一步：1天
- 第二步：0.5天
- 第三步：0.5天
- 第四步：1天
*/

// ===== 数据格式说明 =====

/*
数据格式处理：

1. 前端到后端：
   - 选中的值组成数组：["read", "write"]
   - 或逗号分隔字符串："read,write"

2. 后端到前端：
   - 数组格式：["read", "write"]
   - 字符串格式需要分割处理

3. 存储格式：
   - 数据库：JSON数组或逗号分隔字符串
   - 推荐使用JSON数组格式

4. 验证处理：
   - 数组类型：使用dive验证每个元素
   - 字符串类型：先分割再验证

示例代码：
```go
// 数据结构定义
type UserPreferences struct {
    Permissions []string `json:"permissions" validate:"required,min=1,dive,oneof=read write delete"`
    Features    []string `json:"features" validate:"omitempty,dive,oneof=notification backup sync"`
}

// 数据处理函数
func processCheckboxData(value string) []string {
    if value == "" {
        return []string{}
    }
    return strings.Split(value, ",")
}

func formatCheckboxData(values []string) string {
    return strings.Join(values, ",")
}
```
*/
