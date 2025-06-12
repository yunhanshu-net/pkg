package dev

// Radio组件设计文档
// 单选框相关的组件，适用于互斥选择场景

// ===== Radio组件定义 =====

// RadioWidget 单选框组件定义
type RadioWidget struct {
	Options      string `json:"options"`       // 选项列表：value1(label1),value2(label2)
	DefaultValue string `json:"default_value"` // 默认值
	Direction    string `json:"direction"`     // 排列方向：horizontal(水平) / vertical(垂直)
}

// 使用示例：
// Gender string `runner:"code:gender;name:性别;widget:radio;options:male(男),female(女);direction:horizontal" validate:"required,oneof=male female"`
// Priority string `runner:"code:priority;name:优先级;widget:radio;options:low(低),medium(中),high(高);direction:vertical;default_value:medium" validate:"required"`

// ===== Radio组件设计 =====

/*
Radio组件用于处理单选场景，适用于互斥选择的情况。

设计目标：
1. 支持单选功能
2. 提供水平和垂直排列
3. 支持默认选中项
4. 与现有验证规则无缝集成
5. 提供良好的用户体验

使用场景：
- 性别选择：男/女
- 优先级选择：高/中/低
- 状态选择：启用/禁用
- 类型选择：个人/企业
- 评价选择：满意/一般/不满意

与Select组件的区别：
- Radio：所有选项同时可见，适合选项较少的场景
- Select：下拉选择，适合选项较多的场景
*/

// ===== 组件示例 =====

// RadioExample 单选框组件示例
type RadioExample struct {
	// 性别选择（水平排列）
	Gender string `runner:"code:gender;name:性别;widget:radio;options:male(男),female(女);direction:horizontal" validate:"required,oneof=male female" json:"gender"`
	// 注释：性别单选，水平排列，必填验证

	// 优先级选择（垂直排列）
	Priority string `runner:"code:priority;name:优先级;widget:radio;options:low(低),medium(中),high(高);direction:vertical;default_value:medium" validate:"required,oneof=low medium high" json:"priority"`
	// 注释：优先级单选，垂直排列，默认中等优先级

	// 用户类型选择
	UserType string `runner:"code:user_type;name:用户类型;widget:radio;options:personal(个人),enterprise(企业);direction:horizontal;default_value:personal" validate:"required,oneof=personal enterprise" json:"user_type"`
	// 注释：用户类型选择，默认个人用户

	// 满意度评价
	Satisfaction string `runner:"code:satisfaction;name:满意度;widget:radio;options:satisfied(满意),neutral(一般),dissatisfied(不满意);direction:vertical" validate:"required,oneof=satisfied neutral dissatisfied" json:"satisfaction"`
	// 注释：满意度评价，垂直排列，必填

	// 支付方式
	PaymentMethod string `runner:"code:payment_method;name:支付方式;widget:radio;options:alipay(支付宝),wechat(微信),bank(银行卡);direction:horizontal" validate:"required,oneof=alipay wechat bank" json:"payment_method"`
	// 注释：支付方式选择，水平排列

	// 通知频率
	NotificationFreq string `runner:"code:notification_freq;name:通知频率;widget:radio;options:realtime(实时),daily(每日),weekly(每周),never(从不);direction:vertical;default_value:daily" validate:"required,oneof=realtime daily weekly never" json:"notification_freq"`
	// 注释：通知频率选择，默认每日

	// 账户状态
	AccountStatus string `runner:"code:account_status;name:账户状态;widget:radio;options:active(激活),inactive(未激活),suspended(暂停);direction:vertical;default_value:active" validate:"required,oneof=active inactive suspended" json:"account_status"`
	// 注释：账户状态选择，默认激活状态

	// 可选的单选（无默认值）
	Theme string `runner:"code:theme;name:主题;widget:radio;options:light(浅色),dark(深色),auto(自动);direction:horizontal" validate:"omitempty,oneof=light dark auto" json:"theme"`
	// 注释：主题选择，可选字段，无默认值
}

// ===== 标签配置详解 =====

/*
Radio组件支持的标签：

核心标签：
- code: 字段代码（必需）
- name: 显示名称（必需）
- widget: radio（必需）

选项配置：
- options: 选项列表，格式：value1(label1),value2(label2)
- default_value: 默认值，必须在options中存在

布局配置：
- direction: 排列方向
  - horizontal: 水平排列（适合选项较少）
  - vertical: 垂直排列（默认，适合选项较多）

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
- 返回值：string（单个选中值）
- 空值处理：未选择时返回空字符串
*/

// ===== 实现要点 =====

/*
前端实现要点：

1. 选项渲染：
   - 解析options配置为单选框列表
   - 支持value(label)格式
   - 处理默认选中状态

2. 单选逻辑：
   - 互斥选择实现
   - 选中状态管理
   - 值变化处理

3. 布局控制：
   - 水平/垂直排列
   - 响应式设计
   - 间距控制

4. 用户体验：
   - 清晰的选中状态
   - 点击区域优化
   - 键盘操作支持

5. 样式设计：
   - 统一的单选框样式
   - 选中状态指示
   - 禁用状态样式

后端实现要点：

1. 数据验证：
   - 验证选中值在options中
   - 单值验证
   - 数据类型验证

2. 数据处理：
   - 字符串值处理
   - 默认值设置
   - 空值处理

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

示例验证规则：
```go
// 必填单选验证
Gender string `validate:"required,oneof=male female"`

// 可选单选验证
Theme string `validate:"omitempty,oneof=light dark auto"`

// 带默认值的验证
Priority string `validate:"required,oneof=low medium high"`
```
*/

// ===== 使用最佳实践 =====

/*
最佳实践建议：

1. 选项设计：
   - 选项数量适中（建议2-7个）
   - 选项文本简洁明确
   - 选项之间互斥且完整

2. 排列方向：
   - 2-3个选项：推荐水平排列
   - 4个以上选项：推荐垂直排列
   - 考虑移动端适配

3. 默认值设置：
   - 为常用场景设置合理默认值
   - 默认值必须在options中存在
   - 考虑业务逻辑的默认状态

4. 验证规则：
   - 必须使用oneof验证
   - 验证值与options的value保持一致
   - 考虑必填和可选场景

5. 与Select组件选择：
   - 选项少且需要同时可见：使用Radio
   - 选项多或需要节省空间：使用Select
   - 考虑用户操作习惯

示例：
```go
// 推荐的配置
type GoodExample struct {
    // 性别选择（水平）
    Gender string `runner:"code:gender;name:性别;widget:radio;options:male(男),female(女);direction:horizontal"`

    // 优先级选择（垂直）
    Priority string `runner:"code:priority;name:优先级;widget:radio;options:low(低),medium(中),high(高);direction:vertical;default_value:medium"`

    // 用户类型（水平）
    UserType string `runner:"code:user_type;name:用户类型;widget:radio;options:personal(个人),enterprise(企业);direction:horizontal"`
}

// 不推荐的配置
type BadExample struct {
    // 选项过多
    Country string `runner:"code:country;name:国家;widget:radio;options:cn(中国),us(美国),...20个选项"` // 应该使用Select

    // 缺少验证
    Gender string `runner:"code:gender;name:性别;widget:radio;options:male,female"` // 应该添加oneof验证

    // 默认值不在选项中
    Status string `runner:"code:status;name:状态;widget:radio;options:active,inactive;default_value:pending"` // 默认值必须在options中
}
```
*/

// ===== 实现优先级 =====

/*
实现步骤：

第一步：基础单选功能
- 单选框渲染
- 单选状态管理
- 默认值设置

第二步：布局控制
- 水平/垂直排列
- 响应式适配
- 样式优化

第三步：验证集成
- oneof验证支持
- 错误提示优化
- 边界情况处理

第四步：用户体验优化
- 交互优化
- 键盘操作
- 无障碍支持

预计总工期：1-2天
- 第一步：0.5天
- 第二步：0.5天
- 第三步：0.5天
- 第四步：0.5天
*/

// ===== 设计对比 =====

/*
Radio vs Select vs Checkbox 对比：

1. Radio（单选框）：
   - 用途：互斥单选
   - 可见性：所有选项同时可见
   - 选项数量：2-7个
   - 适用场景：性别、优先级、状态等

2. Select（下拉选择）：
   - 用途：单选或多选
   - 可见性：点击后显示选项
   - 选项数量：不限
   - 适用场景：城市、分类、大量选项等

3. Checkbox（复选框）：
   - 用途：多选
   - 可见性：所有选项同时可见
   - 选项数量：不限（建议不超过15个）
   - 适用场景：权限、标签、功能开关等

选择建议：
- 2-3个互斥选项：Radio（水平）
- 4-7个互斥选项：Radio（垂直）或Select
- 8个以上互斥选项：Select
- 多选场景：Checkbox或Select（多选模式）
*/

// ===== 数据格式说明 =====

/*
数据格式处理：

1. 前端到后端：
   - 选中的单个值：string类型
   - 未选择：空字符串""

2. 后端到前端：
   - 单个字符串值
   - 空值处理

3. 存储格式：
   - 数据库：VARCHAR字段
   - 推荐使用枚举类型

4. 验证处理：
   - 使用oneof验证规则
   - 确保值在选项列表中

示例代码：
```go
// 数据结构定义
type UserProfile struct {
    Gender   string `json:"gender" validate:"required,oneof=male female"`
    Priority string `json:"priority" validate:"required,oneof=low medium high"`
    Theme    string `json:"theme" validate:"omitempty,oneof=light dark auto"`
}

// 数据处理函数
func validateRadioValue(value string, options []string) bool {
    for _, option := range options {
        if value == option {
            return true
        }
    }
    return false
}

func getDefaultRadioValue(options []string, defaultValue string) string {
    if validateRadioValue(defaultValue, options) {
        return defaultValue
    }
    return "" // 或返回第一个选项
}
```
*/
