package dev

// DateTime组件设计文档
// 日期时间相关的输入组件，支持多种日期时间格式

// ===== 重要说明 =====
/*
以下参数是在外层通用结构中定义的，组件内部不需要重复定义：
- Type: string (固定类型，对应Go类型系统)
- Widget: string (组件类型，固定为"datetime")
- Code: string (字段标识)
- CnName: string (显示名称)
- Required: bool (是否必填)

DateTimeWidget结构体只需要定义组件特有的配置参数。
*/

// ===== DateTime组件定义 =====

// DateTimeWidget 日期时间组件定义（统一类型版）
type DateTimeWidget struct {
	Format string `json:"format"` // 具体格式：date/datetime/time/daterange/datetimerange/month/year/week

	// 占位符配置
	Placeholder      string `json:"placeholder"`       // 占位符文本
	StartPlaceholder string `json:"start_placeholder"` // 范围选择时开始日期占位符
	EndPlaceholder   string `json:"end_placeholder"`   // 范围选择时结束日期占位符

	// 默认值和限制（业务核心参数）
	DefaultValue string `json:"default_value"` // 默认值，支持today、now等特殊值
	DefaultTime  string `json:"default_time"`  // 选中日期后的默认具体时刻
	MinDate      string `json:"min_date"`      // 最小可选日期
	MaxDate      string `json:"max_date"`      // 最大可选日期

	// 范围选择配置
	Separator string `json:"separator"` // 日期范围分隔符，默认"至"

	// 快捷选项（高价值功能）
	Shortcuts string `json:"shortcuts"` // 快捷选项配置，JSON字符串
}

// 格式映射表（系统内部使用，前端根据format渲染对应组件）
var FormatComponentMapping = map[string]string{
	"date":          "YYYY-MM-DD",          // 2025-06-13
	"datetime":      "YYYY-MM-DD HH:mm:ss", // 2025-06-13 12:56:48
	"time":          "HH:mm:ss",            // 12:56:48
	"daterange":     "YYYY-MM-DD",          // 2025-06-13 ~ 2025-06-20
	"datetimerange": "YYYY-MM-DD HH:mm:ss", // 2025-06-13 12:56:48 ~ 2025-06-20 18:30:00
	"month":         "YYYY-MM",             // 2025-06
	"year":          "YYYY",                // 2025
	"week":          "YYYY-[W]WW",          // 2025-W24
}

// 使用示例：
// Birthday string `runner:"code:birthday;name:生日;type:string;widget:datetime;format:date;placeholder:请选择生日" validate:"required"`
// EventTime string `runner:"code:event_time;name:活动时间;type:string;widget:datetime;format:datetime;default_time:12:00:00" validate:"required"`
// ValidPeriod string `runner:"code:valid_period;name:有效期;type:string;widget:datetime;format:daterange;separator:至;shortcuts:快捷选项" validate:"required"`
// MeetingTime string `runner:"code:meeting_time;name:会议时间;type:string;widget:datetime;format:time"`

// ===== DateTime组件设计 =====

/*
DateTime组件用于处理日期和时间的输入，现在使用type:string + widget:datetime。

设计目标：
1. type固定为"string"，对应Go的string类型
2. 使用widget:datetime标识日期时间选择器组件
3. 使用format字段区分具体的日期时间组件类型
4. 提供日期范围选择功能
5. 与现有验证规则无缝集成

使用场景：
- 日期选择：生日、创建日期等 (format:date)
- 时间选择：预约时间、提醒时间等 (format:time)
- 日期时间选择：活动时间、截止时间等 (format:datetime)
- 日期范围选择：查询时间段、有效期等 (format:daterange)

格式说明：
- date: 日期选择器，单个日期
- datetime: 日期时间选择器，包含日期和时间
- time: 时间选择器，只选择时间
- daterange: 日期范围选择器，选择日期区间
- datetimerange: 日期时间范围选择器，选择日期时间区间
- month: 月份选择器，选择年月
- year: 年份选择器，选择年份
- week: 周选择器，选择周
*/

// ===== 组件示例 =====

// DateTimeExample 日期时间组件示例（统一类型版）
type DateTimeExample struct {
	// 基础日期选择器 - type:string, widget:datetime, format:date
	BirthDate string `runner:"code:birth_date;name:出生日期;type:string;widget:datetime;format:date;placeholder:请选择出生日期" validate:"required" json:"birth_date"`

	// 时间选择器 - type:string, widget:datetime, format:time
	MeetingTime string `runner:"code:meeting_time;name:会议时间;type:string;widget:datetime;format:time;placeholder:请选择会议时间" json:"meeting_time"`

	// 日期时间选择器 - type:string, widget:datetime, format:datetime
	EventDateTime string `runner:"code:event_datetime;name:活动时间;type:string;widget:datetime;format:datetime;default_time:12:00:00;placeholder:请选择活动时间" validate:"required" json:"event_datetime"`

	// 日期范围选择器 - type:string, widget:datetime, format:daterange
	ValidPeriod string `runner:"code:valid_period;name:有效期;type:string;widget:datetime;format:daterange;separator:至;start_placeholder:开始日期;end_placeholder:结束日期;shortcuts:[{\"text\":\"最近一周\",\"value\":\"week\"},{\"text\":\"最近一个月\",\"value\":\"month\"}]" json:"valid_period"`

	// 日期时间范围选择器 - type:string, widget:datetime, format:datetimerange
	EventPeriod string `runner:"code:event_period;name:活动时间段;type:string;widget:datetime;format:datetimerange;separator:至" json:"event_period"`

	// 带默认值的日期选择 - type:string, widget:datetime, format:date
	CreateDate string `runner:"code:create_date;name:创建日期;type:string;widget:datetime;format:date;default_value:today;placeholder:请选择创建日期" json:"create_date"`

	// 过去日期限制 - type:string, widget:datetime, format:date
	HistoryDate string `runner:"code:history_date;name:历史日期;type:string;widget:datetime;format:date;max_date:today;placeholder:请选择历史日期" validate:"required" json:"history_date"`

	// 未来日期限制 - type:string, widget:datetime, format:date
	FutureDate string `runner:"code:future_date;name:预约日期;type:string;widget:datetime;format:date;min_date:today;placeholder:请选择预约日期" validate:"required" json:"future_date"`

	// 年月选择器 - type:string, widget:datetime, format:month
	ExpireMonth string `runner:"code:expire_month;name:到期月份;type:string;widget:datetime;format:month;placeholder:请选择到期月份" json:"expire_month"`

	// 年份选择器 - type:string, widget:datetime, format:year
	GraduateYear string `runner:"code:graduate_year;name:毕业年份;type:string;widget:datetime;format:year;placeholder:请选择毕业年份" json:"graduate_year"`

	// 周选择器 - type:string, widget:datetime, format:week
	WorkWeek string `runner:"code:work_week;name:工作周;type:string;widget:datetime;format:week;placeholder:请选择工作周" json:"work_week"`
}

// ===== 时间运算示例 =====

/*
使用typex.Time进行时间运算的示例：

1. 日期比较：
   if user.BirthDate.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
       // 2000年前出生
   }

2. 时间差计算：
   age := time.Since(time.Time(user.BirthDate))
   years := int(age.Hours() / 24 / 365)

3. 日期加减：
   nextWeek := typex.Time(time.Time(user.CreateDate).AddDate(0, 0, 7))

4. Unix时间戳：
   timestamp := user.EventDateTime.GetUnix()

5. 格式化输出：
   formatted := time.Time(user.BirthDate).Kind("2006年01月02日")

6. 数据库查询：
   // 直接用于数据库查询，无需格式转换
   db.Where("create_date > ?", user.CreateDate)
*/

// ===== 标签配置详解 =====

/*
DateTime组件支持的标签配置：

核心标签：
- code: 字段代码（必需）
- name: 显示名称（必需）
- type: 统一类型，固定为"datetime"（必需）
- format: 具体格式（必需）- date/datetime/time/daterange/datetimerange/month/year/week

占位符配置：
- placeholder: 单个选择器的占位符文本
- start_placeholder: 范围选择器开始日期占位符
- end_placeholder: 范围选择器结束日期占位符

默认值配置：
- default_value: 默认值，支持today、now等特殊值
- default_time: 选中日期后的默认具体时刻（如12:00:00）

日期限制：
- min_date: 最小可选日期，支持today、具体日期等
- max_date: 最大可选日期，支持today、具体日期等

范围选择配置：
- separator: 日期范围分隔符，默认为"至"

快捷选项：
- shortcuts: 快捷选项配置，JSON格式字符串

统一格式规范：
- 后端统一格式: 2025-06-13 12:56:48 (对应Go的"2006-01-02 15:04:05")
- 前端显示格式: 根据format自动映射标准显示样式
- 数据传输格式: 始终使用"2025-06-13 12:56:48"格式
- 不同精度处理:
  * 年份: "2025" -> "2025-01-01 00:00:00"
  * 年月: "2025-06" -> "2025-06-01 00:00:00"
  * 日期: "2025-06-13" -> "2025-06-13 00:00:00"
  * 完整: "2025-06-13 12:56:48" -> 保持不变

前端默认支持的功能（无需配置）：
- 可清空：所有日期时间组件默认支持清空
- 可手动输入：默认支持键盘输入和日历选择
- 统一尺寸：由渲染引擎统一控制组件大小
- 标准图标：使用系统默认的日期时间图标
*/

// ===== 实现要点 =====

/*
前端实现要点（基于Element Plus组件）：

1. 组件识别：
   - 根据type="datetime"识别为时间日期组件
   - 根据format选择对应的Element Plus组件:
     * date -> el-date-picker (type="date")
     * datetime -> el-date-picker (type="datetime")
     * time -> el-time-picker
     * daterange -> el-date-picker (type="daterange")
     * datetimerange -> el-date-picker (type="datetimerange")
     * month -> el-date-picker (type="month")
     * year -> el-date-picker (type="year")
     * week -> el-date-picker (type="week")

2. 格式处理：
   - 根据format自动设置Element Plus的format和value-format
   - 支持Element Plus的所有日期格式
   - 自动格式验证

3. 默认值处理：
   - default_value: 支持today、now等动态值
   - default_time: 日期时间选择器的默认时间
   - 正确处理时区和本地化

4. 日期限制：
   - picker-options中配置disabledDate函数
   - min_date、max_date转换为日期限制

5. 快捷选项：
   - shortcuts配置转换为Element Plus的shortcuts格式
   - 支持预定义快捷选项

6. 用户体验：
   - placeholder配置: 单个和范围选择的占位符
   - separator: 自定义分隔符
   - 默认支持清空和手动输入功能

后端实现要点：

1. 数据验证：
   - 根据format验证输入格式
   - 验证日期范围和时间有效性
   - 处理时区转换和标准化

2. 存储处理：
   - 统一使用ISO 8601格式存储
   - 时区信息保存
   - 数据库索引优化

3. API处理：
   - 输入格式自动识别和转换
   - 输出格式根据前端需求格式化
   - 错误信息本地化
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

示例验证规则：
```go
// 必填的日期字段
BirthDate typex.Time `runner:"type:datetime;format:date" validate:"required"`

// 日期范围验证
StartDate typex.Time `runner:"type:datetime;format:date;min_date:today"`
EndDate   typex.Time `runner:"type:datetime;format:date;max_date:2025-12-31"`
```
*/

// ===== 使用最佳实践 =====

/*
最佳实践建议（统一type版）：

1. type统一性：
   - 所有日期时间相关组件统一使用 type:datetime
   - 通过format字段区分具体的组件类型
   - 保持与string、number等基础类型的一致性

2. format选择：
   - date: 日期选择，如生日、创建日期
   - datetime: 完整日期时间，如活动时间、截止时间
   - time: 仅时间选择，如提醒时间、会议时间
   - daterange: 日期范围，如查询时间段、有效期
   - datetimerange: 日期时间范围，如活动时间段
   - month: 月份选择，如到期月份
   - year: 年份选择，如毕业年份
   - week: 周选择，如工作周

3. 参数精简原则：
   - 只保留业务必需的核心参数
   - format和type是必需参数
   - placeholder、min_date、max_date是高价值参数
   - shortcuts适用于范围选择场景

推荐的配置示例：
```go
    // 标准日期选择
EventDate typex.Time `runner:"code:event_date;name:活动日期;type:datetime;format:date;default_value:today"`

    // 日期时间选择
StartTime typex.Time `runner:"code:start_time;name:开始时间;type:datetime;format:datetime;min_date:today"`

    // 日期范围选择
ValidPeriod string `runner:"code:valid_period;name:有效期;type:datetime;format:daterange;separator:至"`

// 时间选择
MeetingTime typex.Time `runner:"code:meeting_time;name:会议时间;type:datetime;format:time"`
```
*/

// ===== 实现优先级 =====

/*
实现步骤（基于Element Plus）：

第一步：基础组件映射 (0.5天)
- 建立type:datetime + format的组件映射关系
- 实现date、time、datetime基础组件
- 配置自动格式映射

第二步：属性配置系统 (0.5天)
- 实现placeholder、min_date、max_date等配置
- default_value、default_time处理
- 基础验证集成

第三步：范围选择组件 (0.5天)
- daterange、datetimerange组件实现
- start_placeholder、end_placeholder配置
- separator分隔符处理

第四步：高级功能和优化 (0.5天)
- shortcuts快捷选项配置
- month、year、week特殊组件
- 错误处理和用户体验优化

预计总工期：2天
核心优势：type统一性，format灵活性，配置精简化
*/
