package dev

// Checkbox组件设计文档
// 固定选项的多选组件，适用于预定义选项的多选场景

// ===== 重要说明 =====
/*
以下参数是在外层通用结构中定义的，组件内部不需要重复定义：
- Type: string (固定类型，对应Go类型系统，但实际是[]string)
- Widget: string (组件类型，固定为"checkbox")
- Required: bool (是否必选至少一项)
- Code: string (字段标识)
- CnName: string (显示名称)

CheckboxWidget结构体只需要定义组件特有的配置参数。
*/

// ===== Checkbox组件定义 =====

// CheckboxWidget 复选框组件定义
type CheckboxWidget struct {
	// 选项配置
	Options      string `json:"options"`       // 选项列表：value1(label1),value2(label2)
	DefaultValue string `json:"default_value"` // 默认选中值，多个用逗号分隔

	// 多选配置
	MultipleLimit int `json:"multiple_limit"` // 最多选择数量，0为不限制

	// 布局配置
	Direction string `json:"direction"` // 排列方向：horizontal/vertical，默认vertical
	Columns   int    `json:"columns"`   // 列数（grid布局），仅当direction为vertical时有效

	// 交互配置
	ShowSelectAll bool `json:"show_select_all"` // 是否显示全选/反选按钮
	Disabled      bool `json:"disabled"`        // 是否禁用整个组件
}

// 使用示例：
// Permissions []string `runner:"code:permissions;name:权限;type:string;widget:checkbox;options:read(读取),write(写入),delete(删除),admin(管理);default_value:read,write;direction:vertical" validate:"required,min=1"`
// Workdays []string `runner:"code:workdays;name:工作日;type:string;widget:checkbox;options:monday(周一),tuesday(周二),wednesday(周三),thursday(周四),friday(周五);direction:horizontal" validate:"required,min=1"`

// ===== Checkbox组件设计 =====

/*
Checkbox组件用于处理固定选项的多选场景，选项预定义且数量有限。

设计目标：
1. type固定为"string"，对应Go的string类型（实际使用时为[]string）
2. 使用widget:checkbox标识复选框组件
3. 通过options配置预定义选项
4. 支持灵活的布局控制
5. 专注于固定选项，与动态multiselect分离

使用场景：
- 权限选择：用户权限、角色权限等固定权限项
- 功能开关：系统功能的开启/关闭控制
- 偏好设置：用户偏好选项配置
- 工作日选择：固定的星期选择
- 通知类型：邮件、短信、推送等固定通知方式
- 兴趣标签：预定义的兴趣爱好标签

选项说明：
- 选项数量：通常2-15个，超过15个建议使用multiselect
- 选项内容：固定不变的预定义选项
- 显示方式：所有选项同时可见，无需搜索
*/

// ===== 组件示例 =====

// CheckboxExample 复选框组件示例
type CheckboxExample struct {
	// 权限多选（垂直排列） - type:string, widget:checkbox
	Permissions []string `runner:"code:permissions;name:权限;type:string;widget:checkbox;options:read(读取),write(写入),delete(删除),admin(管理);default_value:read,write;direction:vertical;show_select_all:true" validate:"required,min=1" json:"permissions"`

	// 工作日选择（水平排列） - type:string, widget:checkbox
	Workdays []string `runner:"code:workdays;name:工作日;type:string;widget:checkbox;options:monday(周一),tuesday(周二),wednesday(周三),thursday(周四),friday(周五),saturday(周六),sunday(周日);default_value:monday,tuesday,wednesday,thursday,friday;direction:horizontal" validate:"required,min=1" json:"workdays"`

	// 通知类型（网格布局） - type:string, widget:checkbox
	NotificationTypes []string `runner:"code:notification_types;name:通知类型;type:string;widget:checkbox;options:email(邮件),sms(短信),push(推送),wechat(微信),dingtalk(钉钉),slack(Slack);default_value:email,push;direction:vertical;columns:2" json:"notification_types"`

	// 兴趣爱好（垂直排列，可选） - type:string, widget:checkbox
	Hobbies []string `runner:"code:hobbies;name:兴趣爱好;type:string;widget:checkbox;options:reading(阅读),music(音乐),sports(运动),travel(旅行),photography(摄影),coding(编程),gaming(游戏);direction:vertical;columns:2" validate:"omitempty" json:"hobbies"`

	// 功能开关（水平排列，数量限制） - type:string, widget:checkbox
	Features []string `runner:"code:features;name:启用功能;type:string;widget:checkbox;options:analytics(数据分析),reporting(报表),integration(集成),automation(自动化);default_value:analytics;direction:horizontal;multiple_limit:2" validate:"max=2" json:"features"`

	// 语言偏好（垂直排列） - type:string, widget:checkbox
	Languages []string `runner:"code:languages;name:语言偏好;type:string;widget:checkbox;options:zh-cn(简体中文),zh-tw(繁体中文),en(English),ja(日本語),ko(한국어);default_value:zh-cn;direction:vertical" validate:"required,min=1" json:"languages"`

	// 数据导出格式（水平排列） - type:string, widget:checkbox
	ExportFormats []string `runner:"code:export_formats;name:导出格式;type:string;widget:checkbox;options:excel(Excel),csv(CSV),pdf(PDF),json(JSON);default_value:excel,csv;direction:horizontal;show_select_all:true" json:"export_formats"`

	// 订阅内容（垂直排列，全选支持） - type:string, widget:checkbox
	Subscriptions []string `runner:"code:subscriptions;name:订阅内容;type:string;widget:checkbox;options:news(新闻资讯),tech(技术文章),product(产品更新),marketing(营销活动),system(系统通知);direction:vertical;show_select_all:true" json:"subscriptions"`

	// 安全选项（垂直排列，必选） - type:string, widget:checkbox
	SecurityOptions []string `runner:"code:security_options;name:安全选项;type:string;widget:checkbox;options:two_factor(双因素认证),login_alert(登录提醒),session_timeout(会话超时),ip_whitelist(IP白名单);default_value:two_factor,login_alert;direction:vertical;required:true" validate:"required,min=1" json:"security_options"`
}

// ===== 标签配置详解 =====

/*
Checkbox组件支持的标签配置：

核心标签：
- code: 字段代码（必需）
- name: 显示名称（必需）
- type: 固定为"string"（必需）
- widget: 固定为"checkbox"（必需）

选项配置：
- options: 选项列表（必需），格式：value1(label1),value2(label2)
- default_value: 默认选中值，多个用逗号分隔，必须在options中

数量限制：
- multiple_limit: 最多选择数量，0为不限制
- required: 是否必选至少一项

布局配置：
- direction: 排列方向，horizontal（水平）/vertical（垂直），默认vertical
- columns: 列数，用于网格布局，仅vertical时有效

交互配置：
- show_select_all: 是否显示全选/反选按钮
- disabled: 是否禁用整个组件

数据格式：
- 选择值：使用options中的value部分
- 显示文本：使用options中的label部分
- 返回值：[]string（value值的数组）

验证规则：
- required: 必填验证（至少选择一项）
- min/max: 选择数量限制
- omitempty: 可选字段
*/

// ===== 实现要点 =====

/*
前端实现要点：

1. 组件识别：
   - 根据type="string" + widget="checkbox"识别
   - 渲染为el-checkbox-group组件
   - 解析options为复选框列表

2. 选项渲染：
   - 解析options配置为选项列表
   - 支持value(label)格式解析
   - 处理默认选中状态

3. 布局控制：
   - direction控制水平/垂直排列
   - columns实现网格布局（CSS Grid）
   - 响应式适配不同屏幕尺寸

4. 交互功能：
   - show_select_all启用全选/反选按钮
   - multiple_limit控制选择数量
   - disabled状态处理

5. 状态管理：
   - 多选状态管理
   - 数量限制提示
   - 验证状态显示

后端实现要点：

1. 数据验证：
   - 验证所有选中值都在options中
   - 数组长度验证（min/max）
   - 必填验证处理

2. 数据处理：
   - []string类型处理
   - 去重和排序
   - 默认值设置

3. 选项管理：
   - options配置解析
   - 选项有效性验证
   - 国际化支持
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

3. 枚举验证：
   - dive: 深入验证数组中的每个元素
   - oneof: 每个元素必须在options中

示例验证规则：
```go
// 必须至少选择一项
Permissions []string `validate:"required,min=1"`

// 选择数量限制
Features []string `validate:"min=1,max=3"`

// 每个选中值必须在指定列表中
Workdays []string `validate:"required,min=1,dive,oneof=monday tuesday wednesday thursday friday saturday sunday"`
```
*/

// ===== MultiSelect vs Checkbox 对比 =====

/*
Checkbox vs MultiSelect 选择指南：

Checkbox（固定多选）：
✅ 适用场景：
- 选项数量：2-15个
- 选项特点：固定预定义
- 用户体验：所有选项可见
- 典型场景：权限选择、功能开关、偏好设置

✅ 优势：
- 简单直观，无需搜索
- 所有选项一目了然
- 适合快速多选操作
- 配置简单，无需回调

❌ 不适用：
- 选项数量过多（>15个）
- 选项内容动态变化
- 需要搜索功能
- 数据来源于其他系统

MultiSelect（动态多选）：
✅ 适用场景：
- 选项数量：>10个或动态
- 选项特点：通过API查询
- 用户体验：支持搜索过滤
- 典型场景：用户选择、部门选择

选择建议：
- 固定权限选择 → Checkbox
- 用户/部门选择 → MultiSelect
- 功能开关配置 → Checkbox
- 动态标签选择 → MultiSelect
- 工作日选择 → Checkbox
- 项目关联选择 → MultiSelect
*/

// ===== 使用最佳实践 =====

/*
最佳实践建议：

1. 选项设计：
   - 选项数量控制在2-15个之间
   - 选项名称简洁明了
   - 选项分组合理（如按功能分类）
   - 提供合理的默认选择

2. 布局选择：
   - 2-4个选项：推荐horizontal水平排列
   - 5-8个选项：推荐vertical垂直排列
   - 8个以上选项：使用columns网格布局
   - 选项名称较长：使用vertical垂直排列

3. 交互设计：
   - 选项较多时启用show_select_all
   - 设置合理的multiple_limit
   - 提供清晰的选择状态提示
   - 考虑disabled状态的视觉反馈

4. 验证配置：
   - 根据业务需求设置required
   - 合理设置min/max数量限制
   - 提供友好的错误提示信息

推荐配置示例：
```go
// 权限选择（垂直排列，全选支持）
Permissions []string `runner:"code:permissions;name:权限;type:string;widget:checkbox;options:read(读取),write(写入),delete(删除);direction:vertical;show_select_all:true"`

// 工作日选择（水平排列）
Workdays []string `runner:"code:workdays;name:工作日;type:string;widget:checkbox;options:mon(周一),tue(周二),wed(周三),thu(周四),fri(周五);direction:horizontal"`

// 通知类型（网格布局）
Notifications []string `runner:"code:notifications;name:通知类型;type:string;widget:checkbox;options:email(邮件),sms(短信),push(推送),wechat(微信);direction:vertical;columns:2"`
```
*/

// ===== 实现优先级 =====

/*
实现步骤：

第一步：基础复选框功能 (0.5天)
- 建立type:string + widget:checkbox的组件映射
- 实现基础的复选框渲染
- options解析和默认值处理

第二步：布局控制 (0.5天)
- direction水平/垂直排列
- columns网格布局实现
- 响应式适配

第三步：交互功能 (0.5天)
- 全选/反选功能
- 数量限制控制
- 禁用状态处理

第四步：验证和优化 (0.5天)
- 验证规则集成
- 错误提示优化
- 用户体验调优

预计总工期：2天
核心优势：固定选项，布局灵活，交互简单，与动态multiselect分离
*/
