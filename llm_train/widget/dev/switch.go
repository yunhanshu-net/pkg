package dev

// Switch组件设计文档
// 这是第一个要实现的新组件，作为其他组件实现的参考模板

// ===== Switch组件设计 =====

/*
Switch组件是一个开关控件，用于表示两种状态的切换。

设计目标：
1. 支持bool和string两种数据类型
2. 自定义开启/关闭的显示文本
3. 自定义开启/关闭的实际值
4. 与现有验证规则无缝集成
5. 提供良好的用户体验

使用场景：
- 功能开关：启用/禁用某个功能
- 状态切换：激活/停用某个状态
- 布尔选择：是/否的选择
*/

// SwitchWidget 开关组件定义
type SwitchWidget struct {
	TrueLabel    string `json:"true_label"`    // true状态的显示文本，默认"开启"
	FalseLabel   string `json:"false_label"`   // false状态的显示文本，默认"关闭"
	DefaultValue bool   `json:"default_value"` // 默认值，支持true/false，默认为false
}

// 使用示例：
// IsEnabled bool `runner:"code:is_enabled;name:是否启用;widget:switch;true_label:启用;false_label:禁用;default_value:false" validate:"required"`
// SendEmail bool `runner:"code:send_email;name:发送邮件;widget:switch;true_label:发送;false_label:不发送;default_value:true"`

// ===== 数据类型支持 =====

// SwitchExample 开关组件示例（统一使用bool类型）
type SwitchExample struct {
	// 基础布尔开关（使用默认文本，默认关闭）
	IsEnabled bool `runner:"code:is_enabled;name:是否启用;widget:switch;default_value:false" json:"is_enabled"`

	// 自定义文本的布尔开关（默认关闭）
	IsVIP bool `runner:"code:is_vip;name:VIP用户;widget:switch;true_label:是;false_label:否;default_value:false" json:"is_vip"`

	// 带验证的布尔开关（默认私有）
	IsPublic bool `runner:"code:is_public;name:公开状态;widget:switch;true_label:公开;false_label:私有;default_value:false" validate:"required" json:"is_public"`

	// 功能开关示例（默认发送邮件）
	SendEmail bool `runner:"code:send_email;name:发送邮件;widget:switch;true_label:发送;false_label:不发送;default_value:true" json:"send_email"`

	// 状态开关示例（默认激活）
	IsActive bool `runner:"code:is_active;name:是否激活;widget:switch;true_label:激活;false_label:停用;default_value:true" json:"is_active"`

	// 权限开关示例（默认拒绝）
	HasPermission bool `runner:"code:has_permission;name:权限状态;widget:switch;true_label:允许;false_label:拒绝;default_value:false" json:"has_permission"`

	// 模式开关示例（默认自动模式）
	AutoMode bool `runner:"code:auto_mode;name:自动模式;widget:switch;true_label:自动;false_label:手动;default_value:true" json:"auto_mode"`

	// 通知开关示例（默认开启）
	EnableNotification bool `runner:"code:enable_notification;name:启用通知;widget:switch;true_label:开启;false_label:关闭;default_value:true" json:"enable_notification"`

	// 调试模式示例（默认关闭）
	DebugMode bool `runner:"code:debug_mode;name:调试模式;widget:switch;true_label:开启;false_label:关闭;default_value:false" json:"debug_mode"`
}

// ===== 标签配置详解 =====

/*
Switch组件支持的标签：

核心标签：
- code: 字段代码（必需）
- name: 显示名称（必需）
- widget: switch（必需）

显示配置：
- true_label: true状态的显示文本，默认"开启"
- false_label: false状态的显示文本，默认"关闭"
- default_value: 默认值，支持true/false，默认为false

显示控制：
- show: 显示场景控制
- hidden: 隐藏场景控制

设计原则：
1. 统一使用bool类型，不支持string类型
2. 值固定为true/false，简化设计
3. 通过true_label/false_label控制显示文本
4. 通过default_value设置初始状态
5. 请求参数用switch渲染，响应参数用label文本展示
6. 前后端保持一致的数据类型

默认值使用场景：
- 通知功能：通常默认开启（default_value:true）
- 自动模式：根据业务需求默认开启
- 调试模式：通常默认关闭（default_value:false）
- 权限控制：通常默认拒绝（default_value:false）
- 公开状态：通常默认私有（default_value:false）
*/

// ===== 实现要点 =====

/*
前端实现要点：

1. 组件识别：
   - 根据widget:switch识别为开关组件
   - 统一处理bool类型字段

2. 值处理：
   - 统一使用true/false布尔值
   - 简化数据处理逻辑

3. 显示处理：
   - 请求参数：渲染为switch开关组件
   - 响应参数：根据true_label/false_label显示对应文本
   - 提供默认文本："开启"/"关闭"

4. 验证集成：
   - 支持required验证（如果字段必填）
   - 不需要复杂的枚举验证

5. 样式设计：
   - 现代化的开关样式
   - 清晰的状态指示
   - 平滑的切换动画
   - 支持禁用状态

后端实现要点：

1. 标签解析：
   - 解析true_label/false_label配置
   - 提供默认标签文本
   - 验证配置合理性

2. 代码生成：
   - 生成bool类型的字段定义
   - 生成对应的前端switch组件
   - 生成响应展示的label逻辑

3. 验证处理：
   - 支持required验证
   - 类型安全的bool验证
*/

// ===== 测试用例设计 =====

/*
测试用例覆盖：

1. 基础功能测试：
   - bool类型开关的基本切换
   - true/false值的正确处理
   - 默认值处理

2. 配置测试：
   - 自定义true_label/false_label显示
   - 默认标签文本处理
   - 标签配置验证

3. 验证测试：
   - required验证
   - bool类型验证
   - 错误处理

4. 边界测试：
   - 空值处理
   - 无效配置处理
   - 标签文本为空的处理

5. 集成测试：
   - 与表单系统集成
   - 与验证系统集成
   - 请求/响应渲染集成
*/

// ===== 使用最佳实践 =====

/*
最佳实践建议：

1. 数据类型选择：
   - 统一使用bool类型，简化设计
   - 避免使用string或数字类型
   - 保持前后端数据类型一致

2. 文本配置：
   - 使用简洁明确的标签文本
   - 保持true_label和false_label长度相近
   - 考虑国际化需求
   - 语义化标签，如"启用/禁用"、"公开/私有"

3. 验证规则：
   - 根据业务需求添加required验证
   - 不需要复杂的枚举验证
   - 利用bool类型的天然类型安全

4. 命名规范：
   - 使用Is/Has/Can/Enable等前缀
   - 语义化命名，见名知意
   - 保持命名的一致性

5. 渲染规则：
   - 请求参数：渲染为switch开关组件
   - 响应参数：显示对应的label文本
   - 列表展示：优先显示label文本而非true/false

示例：
```go
// 推荐的命名和配置
type GoodExample struct {
    // 基础开关（默认关闭）
    IsEnabled     bool `runner:"code:is_enabled;name:是否启用;widget:switch;true_label:启用;false_label:禁用;default_value:false"`

    // 公开状态（默认私有）
    IsPublic      bool `runner:"code:is_public;name:是否公开;widget:switch;true_label:公开;false_label:私有;default_value:false"`

    // 权限控制（默认拒绝）
    HasPermission bool `runner:"code:has_permission;name:权限状态;widget:switch;true_label:允许;false_label:拒绝;default_value:false"`

    // 自动模式（默认开启）
    AutoMode      bool `runner:"code:auto_mode;name:自动模式;widget:switch;true_label:自动;false_label:手动;default_value:true"`

    // 通知功能（默认开启）
    SendNotification bool `runner:"code:send_notification;name:发送通知;widget:switch;true_label:开启;false_label:关闭;default_value:true"`

    // 调试模式（默认关闭）
    DebugMode     bool `runner:"code:debug_mode;name:调试模式;widget:switch;true_label:开启;false_label:关闭;default_value:false"`
}

// 不推荐的配置
type BadExample struct {
    Enable int    `runner:"code:enable;name:启用;widget:switch"` // 不要使用int类型
    Status string `runner:"code:status;name:状态;widget:switch"` // 不要使用string类型
    IsActive bool `runner:"code:is_active;name:是否激活;widget:switch"` // 缺少default_value，建议明确指定
}
```
*/

// ===== 实现优先级 =====

/*
实现步骤：

第一步：基础功能实现
- bool类型开关的基本功能
- 默认标签文本显示
- 基础样式

第二步：标签配置扩展
- true_label/false_label自定义配置
- 标签解析和验证
- 默认值处理

第三步：渲染逻辑实现
- 请求参数：switch组件渲染
- 响应参数：label文本展示
- 列表页面：优先显示label文本

第四步：验证和优化
- 与现有验证系统集成
- 错误处理和边界情况
- 样式优化和动画效果

预计总工期：1.5-2天
- 第一步：0.5天
- 第二步：0.5天
- 第三步：0.5天
- 第四步：0.5天

优势：
- 设计更简洁统一
- 减少了string类型的复杂性
- 前后端数据类型一致
- 降低了实现和维护成本
*/
