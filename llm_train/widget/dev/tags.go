package dev

// 这个文件定义了标签系统的设计规范
// 包括所有支持的标签和它们的使用规则

// ===== 标签系统设计 =====

/*
标签格式：runner:"key1:value1;key2:value2;key3:value3"

基础标签（必需）：
- code: 字段代码，用于API传输
- name: 字段显示名称，用于前端展示

组件标签：
- widget: 组件类型
- 其他组件特定配置

显示控制标签：
- show: 控制字段在哪些场景显示
- hidden: 控制字段在哪些场景隐藏

业务标签：
- default: 默认值
- placeholder: 占位符文本
- example: 示例值，帮助大模型理解字段用途和格式
*/

// ===== 核心标签定义 =====

// CoreTags 核心标签（必需）
type CoreTags struct {
	Code string `json:"code"` // 字段代码，必需，用于API传输
	Name string `json:"name"` // 显示名称，必需，用于前端展示
}

// 使用示例：
// Username string `runner:"code:username;name:用户名;example:user123"`

// ===== 组件标签定义 =====

// WidgetTags 组件相关标签
type WidgetTags struct {
	Widget string `json:"widget"` // 组件类型：input/select/switch/checkbox/radio/number/date/file
}

// 支持的组件类型：
const (
	WidgetInput    = "input"    // 文本输入框
	WidgetSelect   = "select"   // 下拉选择
	WidgetSwitch   = "switch"   // 开关
	WidgetCheckbox = "checkbox" // 复选框
	WidgetRadio    = "radio"    // 单选框
	WidgetNumber   = "number"   // 数字输入
	WidgetDate     = "date"     // 日期选择
	WidgetFile     = "file"     // 文件上传
)

// ===== 显示控制标签 =====

// DisplayTags 显示控制标签
type DisplayTags struct {
	Show   string `json:"show"`   // 显示场景：list,create,update,detail
	Hidden string `json:"hidden"` // 隐藏场景：list,create,update,detail
}

// 显示场景说明：
const (
	SceneList   = "list"   // 列表页面
	SceneCreate = "create" // 创建页面
	SceneUpdate = "update" // 更新页面
	SceneDetail = "detail" // 详情页面
)

// 使用示例：
// ID uint `runner:"code:id;name:ID;show:list,detail" json:"id"`                    // 只在列表和详情显示
// Password string `runner:"code:password;name:密码;hidden:list,detail"`            // 在列表和详情隐藏
// CreatedAt time.Time `runner:"code:created_at;name:创建时间;show:list,detail"`    // 只在列表和详情显示

// ===== 输入组件标签 =====

// InputTags input组件专用标签
type InputTags struct {
	Mode        string `json:"mode"`        // 输入模式：line_text/text_area/password
	Placeholder string `json:"placeholder"` // 占位符文本
	MaxLength   int    `json:"max_length"`  // 最大长度
	MinLength   int    `json:"min_length"`  // 最小长度
}

// 使用示例：
// Username string `runner:"code:username;name:用户名;widget:input;mode:line_text;placeholder:请输入用户名;example:user123"`
// Content string `runner:"code:content;name:内容;widget:input;mode:text_area;placeholder:请输入内容;example:这是一段示例内容"`
// Password string `runner:"code:password;name:密码;widget:input;mode:password;example:password123"`

// ===== 选择组件标签 =====

// SelectTags select组件专用标签
type SelectTags struct {
	Options      string `json:"options"`       // 选项：value1(label1),value2(label2)
	DefaultValue string `json:"default_value"` // 默认值
	Multiple     bool   `json:"multiple"`      // 多选（计划功能）
	Searchable   bool   `json:"searchable"`    // 可搜索（计划功能）
}

// 使用示例：
// Status string `runner:"code:status;name:状态;widget:select;options:active(启用),inactive(禁用);default_value:active(启用);example:active(启用)" validate:"required,oneof=active(启用) inactive(禁用)"`

// SwitchTags switch组件专用标签
type SwitchTags struct {
	OnText   string `json:"on_text"`   // 开启文本
	OffText  string `json:"off_text"`  // 关闭文本
	OnValue  string `json:"on_value"`  // 开启值
	OffValue string `json:"off_value"` // 关闭值
}

// 使用示例：
// IsEnabled bool `runner:"code:is_enabled;name:是否启用;widget:switch;on_text:启用;off_text:禁用"`
// Status string `runner:"code:status;name:状态;widget:switch;on_value:active;off_value:inactive"`

// CheckboxTags checkbox组件专用标签
type CheckboxTags struct {
	Options      string `json:"options"`       // 选项列表
	DefaultValue string `json:"default_value"` // 默认选中值
	MinSelect    int    `json:"min_select"`    // 最少选择数
	MaxSelect    int    `json:"max_select"`    // 最多选择数
}

// RadioTags radio组件专用标签
type RadioTags struct {
	Options      string `json:"options"`       // 选项列表
	DefaultValue string `json:"default_value"` // 默认值
	Direction    string `json:"direction"`     // 排列方向：horizontal/vertical
}

// ===== 数值组件标签 =====

// NumberTags number组件专用标签
type NumberTags struct {
	Min         float64 `json:"min"`         // 最小值
	Max         float64 `json:"max"`         // 最大值
	Step        float64 `json:"step"`        // 步长
	Precision   int     `json:"precision"`   // 小数位数
	Placeholder string  `json:"placeholder"` // 占位符
	Unit        string  `json:"unit"`        // 单位
}

// 使用示例：
// Age int `runner:"code:age;name:年龄;widget:number;min:1;max:150;step:1;unit:岁;example:25"`
// Price float64 `runner:"code:price;name:价格;widget:number;min:0;step:0.01;precision:2;unit:元;example:1299.99"`

// ===== 日期组件标签 =====

// DateTags date组件专用标签
type DateTags struct {
	Format      string `json:"format"`      // 日期格式
	Placeholder string `json:"placeholder"` // 占位符
	MinDate     string `json:"min_date"`    // 最小日期
	MaxDate     string `json:"max_date"`    // 最大日期
	ShowTime    bool   `json:"show_time"`   // 显示时间
}

// 使用示例：
// Birthday string `runner:"code:birthday;name:生日;widget:date;format:YYYY-MM-DD;placeholder:请选择生日"`
// EventTime string `runner:"code:event_time;name:活动时间;widget:date;show_time:true"`

// ===== 文件组件标签 =====

// FileTags file组件专用标签
type FileTags struct {
	Accept     string `json:"accept"`      // 文件类型
	MaxSize    int64  `json:"max_size"`    // 最大大小
	Multiple   bool   `json:"multiple"`    // 多文件
	UploadText string `json:"upload_text"` // 上传文本
}

// 使用示例：
// Avatar string `runner:"code:avatar;name:头像;widget:file;accept:.jpg,.png;max_size:2097152"`
// Documents []string `runner:"code:documents;name:文档;widget:file;multiple:true"`

// ===== 标签解析规则 =====

/*
标签解析规则：

1. 格式规范：
   - 使用分号(;)分隔不同的标签
   - 使用冒号(:)分隔标签名和值
   - 值中不能包含分号和冒号
   - 布尔值使用true/false

2. 必需标签：
   - code: 必需，字段代码
   - name: 必需，显示名称
   - widget: 组件类型（如果不是默认input）

3. 可选标签：
   - 显示控制：show, hidden
   - 组件配置：根据widget类型确定
   - 业务配置：default, placeholder等

4. 标签优先级：
   - 核心标签 > 组件标签 > 显示标签 > 业务标签

5. 冲突处理：
   - show和hidden冲突时，hidden优先
   - 相同标签重复时，后面的覆盖前面的

6. 默认值：
   - widget默认为input
   - mode默认为line_text
   - show默认为list,create,update,detail

示例完整标签：
```go
type User struct {
    ID       uint   `runner:"code:id;name:ID;show:list,detail;example:1001" json:"id"`
    Username string `runner:"code:username;name:用户名;widget:input;mode:line_text;placeholder:请输入用户名;example:user123" validate:"required,min=3,max=20"`
    Password string `runner:"code:password;name:密码;widget:input;mode:password;hidden:list,detail;example:password123" validate:"required,min=6"`
    Status   string `runner:"code:status;name:状态;widget:select;options:active(启用),inactive(禁用);default_value:active(启用);example:active(启用)" validate:"required,oneof=active(启用) inactive(禁用)"`
    IsVIP    bool   `runner:"code:is_vip;name:VIP用户;widget:switch;on_text:是;off_text:否;example:false"`
    Age      int    `runner:"code:age;name:年龄;widget:number;min:1;max:150;unit:岁;example:25" validate:"min=1,max=150"`
}
```
*/

// ===== 标签验证规则 =====

// TagValidationRules 标签验证规则
type TagValidationRules struct {
	RequiredTags []string            `json:"required_tags"` // 必需标签
	WidgetTags   map[string][]string `json:"widget_tags"`   // 各组件支持的标签
	ConflictTags map[string][]string `json:"conflict_tags"` // 冲突标签组
}

// 标签验证规则定义
var ValidationRules = TagValidationRules{
	RequiredTags: []string{"code", "name"},
	WidgetTags: map[string][]string{
		"input":    {"mode", "placeholder", "max_length", "min_length", "example"},
		"select":   {"options", "default_value", "multiple", "searchable", "example"},
		"switch":   {"on_text", "off_text", "on_value", "off_value", "example"},
		"checkbox": {"options", "default_value", "min_select", "max_select", "example"},
		"radio":    {"options", "default_value", "direction", "example"},
		"number":   {"min", "max", "step", "precision", "placeholder", "unit", "example"},
		"date":     {"format", "placeholder", "min_date", "max_date", "show_time", "example"},
		"file":     {"accept", "max_size", "multiple", "upload_text", "example"},
	},
	ConflictTags: map[string][]string{
		"display": {"show", "hidden"}, // show和hidden不能同时使用
	},
}
