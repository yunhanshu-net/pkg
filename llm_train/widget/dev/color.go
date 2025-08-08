package dev

// Color组件设计文档
// 颜色选择相关的组件，支持多种颜色选择格式

// ===== 重要说明 =====
/*
以下参数是在外层通用结构中定义的，组件内部不需要重复定义：
- Type: string (固定类型，对应Go类型系统)
- Widget: string (组件类型，固定为"color")
- Code: string (字段标识)
- CnName: string (显示名称)
- Required: bool (是否必填)

ColorWidget结构体只需要定义组件特有的配置参数。
*/

// ===== Color组件定义 =====

// ColorWidget 颜色选择器组件定义
type ColorWidget struct {
	Format string `json:"format"` // 颜色格式：hex/rgb/rgba/hsl/hsla

	// 显示配置
	ShowAlpha bool `json:"show_alpha"` // 是否显示透明度

	// 默认值和预设
	DefaultValue string `json:"default_value"` // 默认颜色值
	Predefine    string `json:"predefine"`     // 预定义颜色，逗号分隔

	// 交互配置
	ShowSwatches bool `json:"show_swatches"` // 是否显示色板
	AllowEmpty   bool `json:"allow_empty"`   // 是否允许清空
}

// 颜色格式映射表（系统内部使用）
var ColorFormatMapping = map[string]string{
	"hex":  "#FFFFFF", // 十六进制格式：#FFFFFF
	"rgb":  "rgb()",   // RGB格式：rgb(255, 255, 255)
	"rgba": "rgba()",  // RGBA格式：rgba(255, 255, 255, 1)
	"hsl":  "hsl()",   // HSL格式：hsl(0, 0%, 100%)
	"hsla": "hsla()",  // HSLA格式：hsla(0, 0%, 100%, 1)
}

// 使用示例：
// ThemeColor string `runner:"code:theme_color;name:主题色;type:string;widget:color;format:hex;default_value:#409EFF" validate:"required"`
// BackgroundColor string `runner:"code:bg_color;name:背景色;type:string;widget:color;format:rgba;show_alpha:true;default_value:rgba(64,158,255,0.8)"`

// ===== Color组件设计 =====

/*
Color组件用于处理颜色选择，现在使用type:string + widget:color。

设计目标：
1. type固定为"string"，对应Go的string类型
2. 使用widget:color标识颜色选择器组件
3. 使用format字段区分具体的颜色格式
4. 支持多种颜色格式输出
5. 提供预定义颜色和透明度支持

使用场景：
- 主题颜色选择：网站主题色、品牌色等 (format:hex)
- 背景颜色：支持透明度的背景色 (format:rgba)
- 文字颜色：简单的文字颜色选择 (format:hex)
- 图表颜色：数据可视化的颜色配置 (format:rgb)

格式说明：
- hex: 十六进制颜色，如 #FF0000
- rgb: RGB颜色，如 rgb(255, 0, 0)
- rgba: RGBA颜色，支持透明度，如 rgba(255, 0, 0, 0.8)
- hsl: HSL颜色，如 hsl(0, 100%, 50%)
- hsla: HSLA颜色，支持透明度，如 hsla(0, 100%, 50%, 0.8)
*/

// ===== 组件示例 =====

// ColorExample 颜色选择器组件示例
type ColorExample struct {
	// 主题颜色选择 - type:string, widget:color, format:hex
	ThemeColor string `runner:"code:theme_color;name:主题颜色;type:string;widget:color;format:hex;default_value:#409EFF;placeholder:请选择主题颜色" validate:"required" json:"theme_color"`

	// 背景颜色选择（支持透明度） - type:string, widget:color, format:rgba
	BackgroundColor string `runner:"code:background_color;name:背景颜色;type:string;widget:color;format:rgba;show_alpha:true;default_value:rgba(64,158,255,0.1)" json:"background_color"`

	// 文字颜色选择 - type:string, widget:color, format:hex
	TextColor string `runner:"code:text_color;name:文字颜色;type:string;widget:color;format:hex;default_value:#333333;predefine:#000000,#333333,#666666,#999999,#CCCCCC,#FFFFFF" json:"text_color"`

	// 边框颜色选择 - type:string, widget:color, format:rgb
	BorderColor string `runner:"code:border_color;name:边框颜色;type:string;widget:color;format:rgb;default_value:rgb(220,223,230)" json:"border_color"`

	// 高亮颜色选择（支持透明度） - type:string, widget:color, format:hsla
	HighlightColor string `runner:"code:highlight_color;name:高亮颜色;type:string;widget:color;format:hsla;show_alpha:true;default_value:hsla(210,100%,60%,0.3)" json:"highlight_color"`

	// 警告颜色选择 - type:string, widget:color, format:hex，带预定义颜色
	WarningColor string `runner:"code:warning_color;name:警告颜色;type:string;widget:color;format:hex;default_value:#E6A23C;predefine:#F56C6C,#E6A23C,#67C23A,#409EFF,#909399" json:"warning_color"`

	// 品牌色选择 - type:string, widget:color, format:hex，带色板
	BrandColor string `runner:"code:brand_color;name:品牌色;type:string;widget:color;format:hex;default_value:#409EFF;show_swatches:true" validate:"required" json:"brand_color"`

	// 可选颜色（允许清空） - type:string, widget:color, format:hex
	OptionalColor string `runner:"code:optional_color;name:可选颜色;type:string;widget:color;format:hex;allow_empty:true" validate:"omitempty" json:"optional_color"`
}

// ===== 标签配置详解 =====

/*
Color组件支持的标签配置：

核心标签：
- code: 字段代码（必需）
- name: 显示名称（必需）
- type: 固定为"string"，对应Go类型（必需）
- widget: 固定为"color"，标识颜色选择器（必需）
- format: 颜色格式（必需）- hex/rgb/rgba/hsl/hsla

默认值配置：
- default_value: 默认颜色值，必须符合对应format的格式
- predefine: 预定义颜色列表，多个颜色用逗号分隔

透明度配置：
- show_alpha: 是否显示透明度滑块（仅对rgba/hsla有效）

交互配置：
- show_swatches: 是否显示色板
- allow_empty: 是否允许清空选择

占位符配置：
- placeholder: 占位符文本

颜色格式规范：
- hex格式: #RRGGBB 或 #RGB，如 #FF0000、#F00
- rgb格式: rgb(r, g, b)，如 rgb(255, 0, 0)
- rgba格式: rgba(r, g, b, a)，如 rgba(255, 0, 0, 0.8)
- hsl格式: hsl(h, s%, l%)，如 hsl(0, 100%, 50%)
- hsla格式: hsla(h, s%, l%, a)，如 hsla(0, 100%, 50%, 0.8)

前端默认支持的功能（无需配置）：
- 拾色器：标准的颜色选择器界面
- 输入框：支持手动输入颜色值
- 预览：实时颜色预览
- 验证：自动验证颜色格式
*/

// ===== 实现要点 =====

/*
前端实现要点（基于Element Plus el-color-picker）：

1. 组件识别：
   - 根据type="string" + widget="color"识别为颜色选择器组件
   - 根据format设置颜色格式和输出格式

2. 格式处理：
   - format配置自动设置color-format属性
   - 支持hex、rgb、rgba、hsl、hsla格式
   - 自动格式验证和转换

3. 透明度处理：
   - show_alpha配置控制是否显示Alpha通道
   - 仅对rgba和hsla格式有效

4. 预定义颜色：
   - predefine配置转换为predefine数组
   - 支持常用颜色快速选择

5. 交互优化：
   - show_swatches控制色板显示
   - allow_empty控制是否可清空
   - 支持键盘操作和无障碍访问

后端实现要点：

1. 数据类型：
   - 对应Go的string类型，直接存储颜色字符串
   - 无需特殊的类型转换

2. 数据验证：
   - 根据format验证颜色值格式
   - 支持各种颜色格式的有效性检查
   - 透明度值范围验证（0-1）

3. 格式转换：
   - 支持不同颜色格式间的转换
   - 统一存储格式处理
   - 输出格式适配前端需求
*/

// ===== 验证规则集成 =====

/*
支持的验证规则：

1. 基础验证：
   - required: 必填验证
   - omitempty: 可选字段

2. 格式验证：
   - 自动根据format进行颜色格式验证
   - 支持正则表达式自定义验证

3. 自定义验证：
   - 可以添加颜色值范围验证
   - 透明度值验证

示例验证规则：
```go
// 必填颜色字段
ThemeColor string `runner:"type:string;widget:color;format:hex" validate:"required"`

// 十六进制颜色验证
HexColor string `runner:"type:string;widget:color;format:hex" validate:"required,regexp=^#[0-9A-Fa-f]{6}$"`

// 可选颜色字段
OptionalColor string `runner:"type:string;widget:color;format:hex;allow_empty:true" validate:"omitempty"`
```
*/

// ===== 使用最佳实践 =====

/*
最佳实践建议：

1. 格式选择：
   - hex: 最常用，适合主题色、品牌色等
   - rgba: 需要透明度的背景色、遮罩色等
   - rgb: 简单的RGB颜色，无透明度需求
   - hsl/hsla: 更直观的色彩控制，设计师友好

2. 默认值设置：
   - 为常用场景设置合理的默认颜色
   - default_value必须符合指定的format格式
   - 考虑品牌色和设计规范

3. 预定义颜色：
   - predefine提供常用颜色快速选择
   - 符合设计规范的色彩搭配
   - 提升用户选择效率

推荐的配置示例：
```go
// 主题色选择
ThemeColor string `runner:"code:theme_color;name:主题色;type:string;widget:color;format:hex;default_value:#409EFF"`

// 背景色选择（支持透明度）
BackgroundColor string `runner:"code:bg_color;name:背景色;type:string;widget:color;format:rgba;show_alpha:true"`

// 文字颜色（带预定义颜色）
TextColor string `runner:"code:text_color;name:文字颜色;type:string;widget:color;format:hex;predefine:#000000,#333333,#666666"`
```
*/

// ===== 实现优先级 =====

/*
实现步骤：

第一步：基础颜色选择器 (0.5天)
- 建立type:string + widget:color的组件映射
- 实现hex、rgb基础格式支持
- 默认值和格式验证

第二步：透明度和高级格式 (0.5天)
- rgba、hsl、hsla格式支持
- show_alpha透明度控制
- 格式转换和验证

第三步：预定义颜色和交互优化 (0.5天)
- predefine预定义颜色支持
- show_swatches色板显示
- allow_empty清空功能

第四步：验证和优化 (0.5天)
- 颜色格式验证集成
- 错误处理和用户提示
- 跨浏览器兼容性测试

预计总工期：2天
核心优势：类型系统一致性，widget语义清晰，颜色格式完整支持
*/
