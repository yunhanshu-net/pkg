package dev

// DateTime组件设计文档
// 日期时间相关的输入组件，支持多种日期时间格式

// ===== DateTime组件定义 =====

// DateTimeWidget 日期时间组件定义
type DateTimeWidget struct {
	Format       string `json:"format"`        // 日期时间格式：YYYY-MM-DD / YYYY-MM-DD HH:mm:ss 等
	Placeholder  string `json:"placeholder"`   // 占位符
	MinDate      string `json:"min_date"`      // 最小日期
	MaxDate      string `json:"max_date"`      // 最大日期
	DefaultValue string `json:"default_value"` // 默认值，支持today、now等
	Separator    string `json:"separator"`     // 日期范围分隔符，默认"~"
}

// 使用示例：
// Birthday string `runner:"code:birthday;name:生日;widget:date;format:YYYY-MM-DD;placeholder:请选择生日" validate:"required"`
// StartDate string `runner:"code:start_date;name:开始日期;widget:date;format:YYYY-MM-DD;min_date:today" validate:"required"`
// EventTime string `runner:"code:event_time;name:活动时间;widget:datetime;format:YYYY-MM-DD HH:mm;show_time:true" validate:"required"`

// ===== DateTime组件设计 =====

/*
DateTime组件用于处理日期和时间的输入，支持多种格式和场景。

设计目标：
1. 支持多种日期时间格式
2. 提供日期范围选择功能
3. 支持时区处理
4. 与现有验证规则无缝集成
5. 提供良好的用户体验

使用场景：
- 日期选择：生日、创建日期等
- 时间选择：预约时间、提醒时间等
- 日期时间选择：活动时间、截止时间等
- 日期范围选择：查询时间段、有效期等
*/

// ===== 组件示例 =====

// DateTimeExample 日期时间组件示例
type DateTimeExample struct {
	// 基础日期选择器
	BirthDate string `runner:"code:birth_date;name:出生日期;widget:date;format:YYYY-MM-DD;placeholder:请选择出生日期" validate:"required" json:"birth_date"`
	// 注释：基础日期选择，使用标准日期格式，必填验证

	// 时间选择器
	MeetingTime string `runner:"code:meeting_time;name:会议时间;widget:time;format:HH:mm;placeholder:请选择会议时间" json:"meeting_time"`
	// 注释：时间选择器，24小时格式，可选字段

	// 日期时间选择器
	EventDateTime string `runner:"code:event_datetime;name:活动时间;widget:datetime;format:YYYY-MM-DD HH:mm:ss;placeholder:请选择活动时间" validate:"required" json:"event_datetime"`
	// 注释：完整的日期时间选择，包含秒，必填验证

	// 日期范围选择器
	ValidPeriod string `runner:"code:valid_period;name:有效期;widget:daterange;format:YYYY-MM-DD;separator:至;placeholder:请选择有效期范围" json:"valid_period"`
	// 注释：日期范围选择，用"至"分隔开始和结束日期

	// 带默认值的日期选择
	CreateDate string `runner:"code:create_date;name:创建日期;widget:date;format:YYYY-MM-DD;default_value:today;placeholder:请选择创建日期" json:"create_date"`
	// 注释：默认值为今天，适用于创建记录的场景

	// 过去日期限制
	HistoryDate string `runner:"code:history_date;name:历史日期;widget:date;format:YYYY-MM-DD;max_date:today;placeholder:请选择历史日期" validate:"required" json:"history_date"`
	// 注释：限制只能选择今天及以前的日期，适用于历史记录

	// 未来日期限制
	FutureDate string `runner:"code:future_date;name:预约日期;widget:date;format:YYYY-MM-DD;min_date:today;placeholder:请选择预约日期" validate:"required" json:"future_date"`
	// 注释：限制只能选择今天及以后的日期，适用于预约场景

	// 年月选择器
	ExpireMonth string `runner:"code:expire_month;name:到期月份;widget:month;format:YYYY-MM;placeholder:请选择到期月份" json:"expire_month"`
	// 注释：只选择年月，适用于信用卡到期日等场景

	// 年份选择器
	GraduateYear string `runner:"code:graduate_year;name:毕业年份;widget:year;format:YYYY;placeholder:请选择毕业年份" json:"graduate_year"`
	// 注释：只选择年份，适用于毕业年份、入职年份等场景
}

// ===== 标签配置详解 =====

/*
DateTime组件支持的标签：

核心标签：
- code: 字段代码（必需）
- name: 显示名称（必需）
- widget: date/time/datetime/daterange/month/year（必需）

格式配置：
- format: 日期时间格式，如YYYY-MM-DD、HH:mm:ss等
- separator: 日期范围分隔符，默认为"~"

默认值配置：
- default_value: 默认值，支持today、now等特殊值
- placeholder: 占位符文本

日期限制：
- min_date: 最小日期，支持today、具体日期等
- max_date: 最大日期，支持today、具体日期等
- disabled_dates: 禁用的日期列表

显示控制：
- show: 显示场景控制
- hidden: 隐藏场景控制

常用格式：
- YYYY-MM-DD: 标准日期格式
- YYYY-MM-DD HH:mm:ss: 完整日期时间格式
- HH:mm: 24小时时间格式
- YYYY-MM: 年月格式
- YYYY: 年份格式
*/

// ===== 实现要点 =====

/*
前端实现要点：

1. 组件识别：
   - 根据widget类型选择对应的日期时间组件
   - 支持date、time、datetime、daterange等类型

2. 格式处理：
   - 根据format配置显示和解析日期时间
   - 支持多种常用格式
   - 提供格式验证

3. 默认值处理：
   - 支持today、now等动态默认值
   - 支持具体日期时间的默认值
   - 正确处理时区

4. 日期限制：
   - 实现min_date、max_date限制
   - 支持禁用特定日期
   - 提供友好的错误提示

5. 用户体验：
   - 提供日历弹窗选择
   - 支持键盘输入
   - 快捷选择功能（今天、昨天等）

后端实现要点：

1. 数据验证：
   - 验证日期时间格式
   - 验证日期范围
   - 处理时区转换

2. 存储处理：
   - 统一存储格式
   - 时区标准化
   - 索引优化

3. API处理：
   - 输入格式验证
   - 输出格式化
   - 错误处理
*/

// ===== 验证规则集成 =====

/*
支持的验证规则：

1. 基础验证：
   - required: 必填验证
   - omitempty: 可选字段

2. 格式验证：
   - 自动根据format进行格式验证
   - 支持自定义日期时间格式

3. 范围验证：
   - 通过min_date、max_date进行范围限制
   - 支持相对日期（如today）

4. 自定义验证：
   - 可以添加自定义的日期验证规则
   - 支持业务逻辑验证

示例验证规则：
```go
// 必填的日期字段
BirthDate string `validate:"required"`

// 日期格式验证（自动根据format验证）
EventDate string `runner:"format:YYYY-MM-DD"`

// 日期范围验证
StartDate string `runner:"min_date:today"`
EndDate   string `runner:"max_date:2025-12-31"`
```
*/

// ===== 使用最佳实践 =====

/*
最佳实践建议：

1. 格式选择：
   - 使用标准的日期时间格式
   - 考虑用户习惯和地区差异
   - 保持格式的一致性

2. 默认值设置：
   - 合理设置默认值，减少用户输入
   - 使用today、now等动态默认值
   - 考虑业务场景的常用值

3. 日期限制：
   - 根据业务逻辑设置合理的日期范围
   - 提供清晰的限制说明
   - 避免过于严格的限制

4. 用户体验：
   - 提供多种输入方式
   - 显示清晰的格式提示
   - 支持快捷操作

5. 数据处理：
   - 统一后端存储格式
   - 正确处理时区
   - 考虑国际化需求

示例：
```go
// 推荐的配置
type GoodExample struct {
    // 标准日期选择
    EventDate string `runner:"code:event_date;name:活动日期;widget:date;format:YYYY-MM-DD;default_value:today"`

    // 日期时间选择
    StartTime string `runner:"code:start_time;name:开始时间;widget:datetime;format:YYYY-MM-DD HH:mm;min_date:today"`

    // 日期范围选择
    ValidPeriod string `runner:"code:valid_period;name:有效期;widget:daterange;format:YYYY-MM-DD;separator:至"`
}

// 不推荐的配置
type BadExample struct {
    // 缺少格式配置
    Date1 string `runner:"code:date1;name:日期;widget:date"` // 应该指定format

    // 格式不标准
    Date2 string `runner:"code:date2;name:日期;widget:date;format:MM/DD/YYYY"` // 建议使用标准格式
}
```
*/

// ===== 实现优先级 =====

/*
实现步骤：

第一步：基础日期组件
- date组件的基本功能
- 标准格式支持
- 基础验证

第二步：时间组件
- time组件实现
- 时间格式处理
- 时间验证

第三步：日期时间组件
- datetime组件实现
- 复合格式处理
- 时区处理

第四步：范围选择组件
- daterange组件实现
- 范围验证
- 用户体验优化

第五步：高级功能
- 年月、年份选择器
- 快捷选择功能
- 自定义验证规则

预计总工期：3-4天
- 第一步：1天
- 第二步：0.5天
- 第三步：1天
- 第四步：1天
- 第五步：0.5-1天
*/
