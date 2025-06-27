package dev

// Input组件设计文档
// 文本输入相关的组件，支持多种输入模式

// ===== Input组件定义 =====

// InputWidget 文本输入组件定义
type InputWidget struct {
	Mode        string `json:"mode"`        // line_text(默认) / text_area / password
	Placeholder string `json:"placeholder"` // 占位符文本
}

// 使用示例：
// Username string `runner:"code:username;name:用户名;widget:input;mode:line_text;placeholder:请输入用户名" validate:"required,min=3,max=20"`
// Content  string `runner:"code:content;name:内容;widget:input;mode:text_area;placeholder:请输入内容" validate:"required,max=1000"`
// Password string `runner:"code:password;name:密码;widget:input;mode:password;placeholder:请输入密码" validate:"required,min=6"`

// ===== Input组件设计 =====

/*
Input组件用于处理文本输入，支持多种输入模式和验证规则。

设计目标：
1. 支持单行文本、多行文本、密码输入
2. 提供长度限制和格式验证
3. 支持占位符和提示信息
4. 与现有验证规则无缝集成
5. 提供良好的用户体验

使用场景：
- 单行文本：用户名、邮箱、标题等
- 多行文本：描述、备注、内容等
- 密码输入：登录密码、确认密码等
*/

// ===== 组件示例 =====

// InputExample 文本输入组件示例
type InputExample struct {
	// 基础单行文本输入
	Username string `runner:"code:username;name:用户名;widget:input;mode:line_text;placeholder:请输入用户名" validate:"required,min=3,max=20,alphanum" json:"username"`
	// 注释：单行文本输入，限制长度3-20，只允许字母数字

	// 邮箱输入
	Email string `runner:"code:email;name:邮箱;widget:input;mode:line_text;placeholder:请输入邮箱地址" validate:"required,email,max=100" json:"email"`
	// 注释：邮箱格式验证，最大长度100

	// 多行文本输入
	Description string `runner:"code:description;name:描述;widget:input;mode:text_area;placeholder:请输入详细描述" validate:"max=1000" json:"description"`
	// 注释：多行文本域，最大长度1000，可选字段

	// 密码输入
	Password string `runner:"code:password;name:密码;widget:input;mode:password;placeholder:请输入密码" validate:"required,min=6,max=50" json:"password"`
	// 注释：密码输入，隐藏显示，长度6-50

	// 确认密码
	ConfirmPassword string `runner:"code:confirm_password;name:确认密码;widget:input;mode:password;placeholder:请再次输入密码" validate:"required,eqfield=Password" json:"confirm_password"`
	// 注释：确认密码，必须与Password字段相等

	// 手机号输入
	Phone string `runner:"code:phone;name:手机号;widget:input;mode:line_text;placeholder:请输入11位手机号" validate:"required,len=11,numeric" json:"phone"`
	// 注释：手机号输入，固定11位数字

	// 网址输入
	Website string `runner:"code:website;name:网站;widget:input;mode:line_text;placeholder:请输入网站地址" validate:"omitempty,url" json:"website"`
	// 注释：网址输入，可选字段，URL格式验证

	// 标题输入
	Title string `runner:"code:title;name:标题;widget:input;mode:line_text;placeholder:请输入标题" validate:"required,max=200" json:"title"`
	// 注释：标题输入，必填，最大长度200
}

// ===== 标签配置详解 =====

/*
Input组件支持的标签：

核心标签：
- code: 字段代码（必需）
- name: 显示名称（必需）
- widget: input（必需）

输入模式：
- mode: 输入模式
  - line_text: 单行文本（默认）
  - text_area: 多行文本
  - password: 密码输入

显示配置：
- placeholder: 占位符文本
- readonly: 是否只读，默认false
- disabled: 是否禁用，默认false

显示控制：
- show: 显示场景控制
- hidden: 隐藏场景控制

常用验证规则：
- required: 必填验证
- min/max: 长度验证
- email: 邮箱格式验证
- url: 网址格式验证
- numeric: 数字验证
- alphanum: 字母数字验证
- eqfield: 字段相等验证
*/

// ===== 实现要点 =====

/*
前端实现要点：

1. 组件识别：
   - 根据mode选择对应的输入组件
   - 支持line_text、text_area、password模式

2. 输入处理：
   - 实时长度验证
   - 格式验证提示
   - 密码强度指示

3. 用户体验：
   - 清晰的占位符提示
   - 实时字符计数
   - 错误状态显示

4. 样式设计：
   - 统一的输入框样式
   - 焦点状态指示
   - 错误状态样式

后端实现要点：

1. 数据验证：
   - 长度验证
   - 格式验证
   - 安全性检查

2. 数据处理：
   - 字符串清理
   - XSS防护
   - 敏感信息处理

3. 存储处理：
   - 文本编码处理
   - 索引优化
   - 搜索支持
*/

// ===== 验证规则集成 =====

/*
支持的验证规则：

1. 基础验证：
   - required: 必填验证
   - omitempty: 可选字段

2. 长度验证：
   - min: 最小长度
   - max: 最大长度
   - len: 固定长度

3. 格式验证：
   - email: 邮箱格式
   - url: 网址格式
   - numeric: 纯数字
   - alpha: 纯字母
   - alphanum: 字母数字

4. 字段比较：
   - eqfield: 与其他字段相等
   - nefield: 与其他字段不相等

示例验证规则：
```go
// 用户名验证
Username string `validate:"required,min=3,max=20,alphanum"`

// 邮箱验证
Email string `validate:"required,email"`

// 密码确认
ConfirmPassword string `validate:"required,eqfield=Password"`
```
*/

// ===== 使用最佳实践 =====

/*
最佳实践建议：

1. 模式选择：
   - 简短文本使用line_text
   - 长文本使用text_area
   - 敏感信息使用password

2. 长度限制：
   - 设置合理的长度限制
   - 考虑数据库字段长度
   - 提供清晰的长度提示

3. 占位符设计：
   - 提供有用的输入提示
   - 避免过长的占位符
   - 考虑国际化需求

4. 验证规则：
   - 前后端保持一致
   - 提供友好的错误提示
   - 考虑用户输入习惯

5. 安全考虑：
   - 输入内容过滤
   - XSS防护
   - 敏感信息保护

示例：
```go
// 推荐的配置
type GoodExample struct {
    // 用户名输入
    Username string `runner:"code:username;name:用户名;widget:input;mode:line_text;placeholder:请输入用户名;max_length:20"`

    // 描述输入
    Description string `runner:"code:description;name:描述;widget:input;mode:text_area;placeholder:请输入详细描述;max_length:1000"`

    // 密码输入
    Password string `runner:"code:password;name:密码;widget:input;mode:password;placeholder:请输入密码;min_length:6"`
}

// 不推荐的配置
type BadExample struct {
    // 缺少长度限制
    Text1 string `runner:"code:text1;name:文本;widget:input"` // 应该设置max_length

    // 占位符过长
    Text2 string `runner:"code:text2;name:文本;widget:input;placeholder:请输入一个非常详细和完整的描述信息"` // 占位符太长
}
```
*/

// ===== 实现优先级 =====

/*
实现步骤：

第一步：基础文本输入
- line_text模式实现
- 基础验证支持
- 长度限制

第二步：多行文本
- text_area模式实现
- 自适应高度
- 字符计数

第三步：密码输入
- password模式实现
- 密码强度指示
- 显示/隐藏切换

第四步：高级功能
- 输入格式化
- 自动完成
- 实时验证

预计总工期：2-3天
- 第一步：1天
- 第二步：0.5天
- 第三步：0.5天
- 第四步：1天
*/
 