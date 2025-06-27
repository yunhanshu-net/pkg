package dev

// Number组件设计文档
// 数值输入相关的组件，支持多种数值输入方式

// ===== Number组件定义 =====

// NumberWidget 数字输入组件定义
type NumberWidget struct {
	Min         float64 `json:"min"`         // 最小值
	Max         float64 `json:"max"`         // 最大值
	Step        float64 `json:"step"`        // 步长
	Precision   int     `json:"precision"`   // 小数位数
	Placeholder string  `json:"placeholder"` // 占位符
	Unit        string  `json:"unit"`        // 单位显示
}

// 使用示例：
// Age int `runner:"code:age;name:年龄;widget:number;min:1;max:150;step:1;unit:岁" validate:"required,min=1,max=150"`
// Price float64 `runner:"code:price;name:价格;widget:number;min:0;step:0.01;precision:2;unit:元;placeholder:请输入价格" validate:"required,min=0"`
// Score float64 `runner:"code:score;name:分数;widget:number;min:0;max:100;step:0.1;precision:1" validate:"required,min=0,max=100"`

// ===== Number组件设计 =====

/*
Number组件用于处理数值输入，支持多种数值类型和输入方式。

设计目标：
1. 支持整数和浮点数输入
2. 提供数值范围限制
3. 支持步长控制
4. 提供滑块、评分等特殊数值组件
5. 与现有验证规则无缝集成

使用场景：
- 基础数值：年龄、数量、价格等
- 范围选择：评分、进度、百分比等
- 计数器：库存、人数等
- 滑块选择：音量、亮度等
*/

// ===== 组件示例 =====

// NumberExample 数值输入组件示例
type NumberExample struct {
	// 基础整数输入
	Age int `runner:"code:age;name:年龄;widget:number;min:1;max:150;step:1;placeholder:请输入年龄" validate:"required,min=1,max=150" json:"age"`
	// 注释：基础整数输入，范围1-150，步长为1，必填验证

	// 浮点数输入 - 价格
	Price float64 `runner:"code:price;name:价格;widget:number;min:0;max:999999.99;step:0.01;placeholder:请输入价格;suffix:元" validate:"required,min=0" json:"price"`
	// 注释：浮点数输入，最小值0，步长0.01，带单位后缀

	// 百分比输入
	Percentage float64 `runner:"code:percentage;name:百分比;widget:number;min:0;max:100;step:0.1;placeholder:请输入百分比;suffix:%" validate:"required,min=0,max=100" json:"percentage"`
	// 注释：百分比输入，范围0-100，步长0.1，带百分号后缀

	// 滑块组件 - 评分
	Rating int `runner:"code:rating;name:评分;widget:slider;min:1;max:10;step:1;default_value:5;show_marks:true" validate:"required,min=1,max=10" json:"rating"`
	// 注释：滑块评分，范围1-10，显示刻度标记，默认值5

	// 滑块组件 - 音量
	Volume int `runner:"code:volume;name:音量;widget:slider;min:0;max:100;step:5;default_value:50;show_tooltip:true;suffix:%" json:"volume"`
	// 注释：音量滑块，范围0-100，步长5，显示提示框

	// 评分组件 - 星级评分
	StarRating int `runner:"code:star_rating;name:星级评分;widget:rate;max:5;allow_half:true;default_value:0" validate:"required,min=1,max=5" json:"star_rating"`
	// 注释：星级评分组件，最高5星，支持半星，必填验证

	// 计数器组件 - 数量
	Quantity int `runner:"code:quantity;name:数量;widget:counter;min:1;max:999;step:1;default_value:1" validate:"required,min=1" json:"quantity"`
	// 注释：计数器组件，带加减按钮，最小值1，默认值1

	// 范围滑块 - 价格区间
	PriceRange string `runner:"code:price_range;name:价格区间;widget:range;min:0;max:10000;step:100;default_value:1000,5000;separator:-" json:"price_range"`
	// 注释：范围滑块，选择价格区间，用"-"分隔最小值和最大值

	// 带前缀的数值输入
	Salary float64 `runner:"code:salary;name:薪资;widget:number;min:0;step:100;placeholder:请输入薪资;prefix:￥" validate:"min=0" json:"salary"`
	// 注释：薪资输入，带人民币符号前缀，步长100

	// 只读数值显示
	TotalAmount float64 `runner:"code:total_amount;name:总金额;widget:number;readonly:true;suffix:元" json:"total_amount"`
	// 注释：只读数值显示，用于计算结果展示
}

// ===== 标签配置详解 =====

/*
Number组件支持的标签：

核心标签：
- code: 字段代码（必需）
- name: 显示名称（必需）
- widget: number/slider/rate/counter/range（必需）

数值配置：
- min: 最小值
- max: 最大值
- step: 步长，默认为1
- precision: 小数位数，默认自动判断

默认值配置：
- default_value: 默认值
- placeholder: 占位符文本

显示配置：
- prefix: 前缀文本，如"￥"
- suffix: 后缀文本，如"元"、"%"
- readonly: 是否只读，默认false

滑块特有配置：
- show_marks: 是否显示刻度标记
- show_tooltip: 是否显示数值提示框
- vertical: 是否垂直显示

评分特有配置：
- allow_half: 是否允许半星评分
- allow_clear: 是否允许清除评分
- character: 自定义评分字符

计数器特有配置：
- controls: 是否显示加减按钮，默认true
- controls_position: 按钮位置，right/left

范围滑块特有配置：
- separator: 范围值分隔符，默认","
- tooltip_format: 提示框格式化

显示控制：
- show: 显示场景控制
- hidden: 隐藏场景控制
*/

// ===== 实现要点 =====

/*
前端实现要点：

1. 组件识别：
   - 根据widget类型选择对应的数值组件
   - 支持number、slider、rate、counter、range等类型

2. 数值处理：
   - 正确处理整数和浮点数
   - 实现步长控制
   - 范围限制和验证

3. 用户交互：
   - 键盘输入支持
   - 鼠标滚轮支持
   - 触摸手势支持

4. 视觉反馈：
   - 数值变化动画
   - 状态指示
   - 错误提示

5. 格式化显示：
   - 前缀后缀显示
   - 千分位分隔符
   - 小数位控制

后端实现要点：

1. 数据验证：
   - 数值类型验证
   - 范围验证
   - 精度验证

2. 数据转换：
   - 字符串到数值转换
   - 精度处理
   - 范围值解析

3. 存储处理：
   - 数值类型存储
   - 精度保持
   - 索引优化
*/

// ===== 验证规则集成 =====

/*
支持的验证规则：

1. 基础验证：
   - required: 必填验证
   - omitempty: 可选字段

2. 数值验证：
   - min: 最小值验证
   - max: 最大值验证
   - gt/gte: 大于/大于等于验证
   - lt/lte: 小于/小于等于验证

3. 自定义验证：
   - 可以添加自定义的数值验证规则
   - 支持业务逻辑验证

示例验证规则：
```go
// 基础数值验证
Age int `validate:"required,min=1,max=150"`

// 浮点数验证
Price float64 `validate:"required,min=0,max=999999.99"`

// 范围验证
Score int `validate:"gte=0,lte=100"`
```
*/

// ===== 使用最佳实践 =====

/*
最佳实践建议：

1. 组件选择：
   - 简单数值输入使用number组件
   - 范围选择使用slider组件
   - 评分使用rate组件
   - 计数使用counter组件

2. 数值范围：
   - 设置合理的最小值和最大值
   - 考虑业务场景的实际需求
   - 提供清晰的范围说明

3. 步长设置：
   - 根据数值精度设置合适的步长
   - 整数使用1，小数使用0.01或0.1
   - 考虑用户输入习惯

4. 用户体验：
   - 提供默认值减少用户输入
   - 显示单位和格式说明
   - 支持多种输入方式

5. 数据处理：
   - 正确处理浮点数精度
   - 统一数值格式
   - 考虑国际化需求

示例：
```go
// 推荐的配置
type GoodExample struct {
    // 基础数值输入
    Age int `runner:"code:age;name:年龄;widget:number;min:1;max:150;step:1;placeholder:请输入年龄"`

    // 价格输入
    Price float64 `runner:"code:price;name:价格;widget:number;min:0;step:0.01;prefix:￥;placeholder:请输入价格"`

    // 评分滑块
    Rating int `runner:"code:rating;name:评分;widget:slider;min:1;max:10;step:1;show_marks:true;default_value:5"`

    // 星级评分
    Stars int `runner:"code:stars;name:星级;widget:rate;max:5;allow_half:true"`
}

// 不推荐的配置
type BadExample struct {
    // 缺少范围限制
    Number1 int `runner:"code:number1;name:数值;widget:number"` // 应该设置min/max

    // 步长不合理
    Price float64 `runner:"code:price;name:价格;widget:number;step:1"` // 价格应该用0.01

    // 缺少默认值
    Rating int `runner:"code:rating;name:评分;widget:slider;min:1;max:10"` // 应该设置默认值
}
```
*/

// ===== 实现优先级 =====

/*
实现步骤：

第一步：基础数值组件
- number组件的基本功能
- 整数和浮点数支持
- 基础验证

第二步：滑块组件
- slider组件实现
- 范围控制
- 视觉反馈

第三步：评分组件
- rate组件实现
- 星级显示
- 交互处理

第四步：计数器组件
- counter组件实现
- 加减按钮
- 键盘支持

第五步：高级功能
- 范围滑块
- 格式化显示
- 自定义样式

预计总工期：3-4天
- 第一步：1天
- 第二步：1天
- 第三步：0.5天
- 第四步：0.5天
- 第五步：1天
*/
 