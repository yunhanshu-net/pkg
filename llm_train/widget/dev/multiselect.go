package dev

// MultiSelect组件设计文档
// 多选择器组件，支持固定选项和动态搜索两种模式

// ===== 重要说明 =====
/*
以下参数是在外层通用结构中定义的，组件内部不需要重复定义：
- Type: string (固定类型，对应Go类型系统，但实际是[]string)
- Widget: string (组件类型，固定为"multiselect")
- Required: bool (是否必选至少一项)
- Code: string (字段标识)
- Name: string (显示名称)

MultiSelectWidget结构体只需要定义组件特有的配置参数。
*/

// ===== MultiSelect组件定义 =====

// MultiSelectWidget 多选择器组件定义
type MultiSelectWidget struct {
	// 选项配置
	Options      string `json:"options"`       // 选项列表：value1(label1),value2(label2)
	DefaultValue string `json:"default_value"` // 默认选中值，多个用逗号分隔，必须与options中的格式完全一致，格式：value1(label1),value2(label2)

	// 多选配置
	MultipleLimit int `json:"multiple_limit"` // 最多选择数量，0为不限制

	// 显示配置
	Placeholder  string `json:"placeholder"`   // 占位符文本
	CollapseTags bool   `json:"collapse_tags"` // 是否折叠显示已选择的标签

	// 创建配置
	AllowCreate bool `json:"allow_create"` // 是否允许创建新条目
}

// 使用示例：
// Tags []string `runner:"code:tags;name:标签;type:[]string;widget:multiselect;options:tech(技术),design(设计),product(产品);default_value:tech(技术),design(设计)" validate:"required,min=1"`
// Users []string `runner:"code:users;name:用户;type:[]string;widget:multiselect;placeholder:搜索用户" callback:"OnInputFuzzy(delay:300,min:2)" validate:"required,min=1"`

// ===== MultiSelect组件设计 =====

/*
MultiSelect组件是一个多选择器，支持两种工作模式：

1. 固定选项模式：
   - 通过options配置预定义选项
   - 适用于选项较少且固定的场景
   - 类似select组件，但支持多选

2. 动态搜索模式：
   - 包含callback:OnInputFuzzy(...)时启用
   - 前端会调用回调获取搜索结果
   - 适用于选项较多或动态变化的场景

回调设计：
- callback标签独立于runner标签，职责分离
- 使用callback="CallbackName(param1:value1,param2:value2)"格式
- 多个回调用分号分隔：callback="OnInputFuzzy(...);OnInputValidate(...)"
- 参数用逗号分隔，回调用分号分隔，层次清晰
- 支持任意数量的参数和回调

设计目标：
1. 专注于多选择器的核心功能
2. 支持固定选项和动态搜索两种模式
3. 回调配置表达能力强且无解析冲突
4. 与Select组件保持一致性
*/

// ===== 组件示例 =====

// MultiSelectExample 多选择器组件示例
type MultiSelectExample struct {
	// 固定选项模式 - 技能标签
	Skills []string `runner:"code:skills;name:技能;type:[]string;widget:multiselect;options:javascript(JavaScript),python(Python),golang(Go),java(Java);default_value:javascript(JavaScript),python(Python);multiple_limit:3;placeholder:选择技能" validate:"required,min=1,max=3" json:"skills"`

	// 固定选项模式 - 权限选择
	Permissions []string `runner:"code:permissions;name:权限;type:[]string;widget:multiselect;options:read(读取),write(写入),delete(删除),admin(管理);default_value:read(读取),write(写入);collapse_tags:true" validate:"required,min=1" json:"permissions"`

	// 动态搜索模式 - 用户选择（基础回调）
	AssignedUsers []string `runner:"code:assigned_users;name:指派用户;type:[]string;widget:multiselect;placeholder:搜索并选择用户;multiple_limit:5" callback:"OnInputFuzzy(delay:300,min:2)" validate:"required,min=1,max=5" json:"assigned_users"`

	// 动态搜索模式 - 部门选择（高级回调）
	Departments []string `runner:"code:departments;name:部门;type:[]string;widget:multiselect;placeholder:搜索部门;collapse_tags:true" callback:"OnInputFuzzy(delay:500,min:3,max_results:20)" validate:"required,min=1" json:"departments"`

	// 动态搜索模式 - 标签选择（支持创建+实时验证）
	Tags []string `runner:"code:tags;name:标签;type:[]string;widget:multiselect;allow_create:true;placeholder:搜索或创建标签" callback:"OnInputFuzzy(delay:300,min:1);OnInputValidate(delay:600,min:2)" json:"tags"`

	// 混合模式 - 预设选项+动态搜索
	Categories []string `runner:"code:categories;name:分类;type:[]string;widget:multiselect;options:tech(技术),business(商业),design(设计);placeholder:选择或搜索分类" callback:"OnInputFuzzy(delay:300,min:2,cache:60)" json:"categories"`

	// 简单多选 - 通知类型
	NotificationTypes []string `runner:"code:notification_types;name:通知类型;type:[]string;widget:multiselect;options:email(邮件),sms(短信),push(推送),wechat(微信);default_value:email(邮件),push(推送)" json:"notification_types"`

	// 复杂回调 - 项目选择（实时验证）
	RelatedProjects []string `runner:"code:related_projects;name:关联项目;type:[]string;widget:multiselect;multiple_limit:2;placeholder:最多选择2个项目" callback:"OnInputFuzzy(delay:300,min:2,max_results:10);OnInputValidate(delay:800,min:1,show_loading:true);OnChange(debounce:200)" validate:"max=2" json:"related_projects"`

	// 高级搜索 - 用户选择（完整回调参数）
	AdvancedUsers []string `runner:"code:advanced_users;name:高级用户选择;type:[]string;widget:multiselect;placeholder:搜索用户" callback:"OnInputFuzzy(delay:500,min:3,max_results:20,cache:120,filter_mode:fuzzy)" json:"advanced_users"`

	// 实时验证 - 邮箱选择（后端验证）
	EmailList []string `runner:"code:email_list;name:邮箱列表;type:[]string;widget:multiselect;placeholder:输入邮箱地址;allow_create:true" callback:"OnInputFuzzy(delay:200,min:1);OnInputValidate(delay:800,min:3,cache:120)" validate:"email" json:"email_list"`

	// 级联选择 - 地区选择
	Regions []string `runner:"code:regions;name:地区;type:[]string;widget:multiselect;placeholder:选择地区" callback:"OnInputFuzzy(delay:300,min:1,cascade:true);OnChange(debounce:200)" json:"regions"`

	// ===== OnInputValidate 典型场景 =====

	// 基础验证 - 快速响应
	BasicValidate []string `runner:"code:basic_validate;name:基础验证;type:[]string;widget:multiselect;allow_create:true;placeholder:输入内容" callback:"OnInputValidate(delay:300,min:2)" validate:"required,min=2,max=20" json:"basic_validate"`

	// 慢速验证 - 复杂检查
	SlowValidate []string `runner:"code:slow_validate;name:慢速验证;type:[]string;widget:multiselect;allow_create:true;placeholder:输入内容" callback:"OnInputValidate(delay:1000,min:3,show_loading:true)" validate:"required" json:"slow_validate"`

	// 缓存验证 - 减少请求
	CachedValidate []string `runner:"code:cached_validate;name:缓存验证;type:[]string;widget:multiselect;allow_create:true;placeholder:输入内容" callback:"OnInputValidate(delay:500,min:2,cache:300)" json:"cached_validate"`

	// 重试验证 - 网络不稳定
	RetryValidate []string `runner:"code:retry_validate;name:重试验证;type:[]string;widget:multiselect;allow_create:true;placeholder:输入内容" callback:"OnInputValidate(delay:600,min:1,retry:3,show_loading:true)" json:"retry_validate"`

	// 实时验证 - 即时反馈
	RealtimeValidate []string `runner:"code:realtime_validate;name:实时验证;type:[]string;widget:multiselect;allow_create:true;placeholder:输入内容" callback:"OnInputValidate(delay:200,min:1,real_time:true)" json:"realtime_validate"`

	// 静默验证 - 无加载状态
	SilentValidate []string `runner:"code:silent_validate;name:静默验证;type:[]string;widget:multiselect;allow_create:true;placeholder:输入内容" callback:"OnInputValidate(delay:400,min:2,show_loading:false)" json:"silent_validate"`

	// 长缓存验证 - 稳定数据
	LongCacheValidate []string `runner:"code:long_cache_validate;name:长缓存验证;type:[]string;widget:multiselect;allow_create:true;placeholder:输入内容" callback:"OnInputValidate(delay:500,min:2,cache:1800)" json:"long_cache_validate"`

	// ===== OnInputFuzzy 典型场景 =====

	// 快速搜索 - 低延迟
	FastSearch []string `runner:"code:fast_search;name:快速搜索;type:[]string;widget:multiselect;placeholder:快速搜索" callback:"OnInputFuzzy(delay:200,min:1,max_results:10)" json:"fast_search"`

	// 精确搜索 - 高门槛
	PreciseSearch []string `runner:"code:precise_search;name:精确搜索;type:[]string;widget:multiselect;placeholder:至少输入3个字符" callback:"OnInputFuzzy(delay:400,min:3,max_results:20,filter_mode:fuzzy)" json:"precise_search"`

	// 大量结果 - 高限制
	MassiveSearch []string `runner:"code:massive_search;name:大量搜索;type:[]string;widget:multiselect;placeholder:搜索内容" callback:"OnInputFuzzy(delay:300,min:2,max_results:100,cache:600)" json:"massive_search"`

	// 缓存搜索 - 减少请求
	CachedSearch []string `runner:"code:cached_search;name:缓存搜索;type:[]string;widget:multiselect;placeholder:搜索内容" callback:"OnInputFuzzy(delay:300,min:2,max_results:20,cache:1800)" json:"cached_search"`

	// 多字段搜索 - 复合查询
	MultiFieldSearch []string `runner:"code:multi_field_search;name:多字段搜索;type:[]string;widget:multiselect;placeholder:搜索多个字段" callback:"OnInputFuzzy(delay:350,min:2,max_results:25,fields:name,code,description)" json:"multi_field_search"`

	// 全文搜索 - 内容检索
	FullTextSearch []string `runner:"code:full_text_search;name:全文搜索;type:[]string;widget:multiselect;placeholder:全文搜索" callback:"OnInputFuzzy(delay:500,min:3,max_results:15,filter_mode:fulltext,highlight:true)" json:"full_text_search"`

	// 层级搜索 - 树形结构
	HierarchySearch []string `runner:"code:hierarchy_search;name:层级搜索;type:[]string;widget:multiselect;placeholder:层级搜索" callback:"OnInputFuzzy(delay:400,min:1,max_results:20,hierarchy:true)" json:"hierarchy_search"`

	// 级联搜索 - 联动查询
	CascadeSearch []string `runner:"code:cascade_search;name:级联搜索;type:[]string;widget:multiselect;placeholder:级联搜索" callback:"OnInputFuzzy(delay:300,min:1,max_results:15,cascade:true)" json:"cascade_search"`

	// 排序搜索 - 结果排序
	SortedSearch []string `runner:"code:sorted_search;name:排序搜索;type:[]string;widget:multiselect;placeholder:排序搜索" callback:"OnInputFuzzy(delay:300,min:2,max_results:20,sort_by:relevance)" json:"sorted_search"`

	// 高亮搜索 - 匹配高亮
	HighlightSearch []string `runner:"code:highlight_search;name:高亮搜索;type:[]string;widget:multiselect;placeholder:高亮搜索" callback:"OnInputFuzzy(delay:350,min:2,max_results:20,highlight:true)" json:"highlight_search"`
}

// ===== 标签配置详解 =====

/*
MultiSelect组件支持的标签配置：

核心标签：
- code: 字段代码（必需）
- name: 显示名称（必需）
- type: 固定为"[]string"（必需）
- widget: 固定为"multiselect"（必需）

选项配置：
- options: 选项列表，格式：value1(label1),value2(label2)
- default_value: 默认选中值，多个用逗号分隔，必须与options中的格式完全一致，格式：value1(label1),value2(label2)

多选配置：
- multiple_limit: 最多选择数量，0为不限制

显示配置：
- placeholder: 占位符文本
- collapse_tags: 是否折叠已选择标签

交互配置：
- allow_create: 允许创建新条目

回调配置（独立callback标签）：
- callback: 独立标签，格式：callback="CallbackName(param1:value1,param2:value2);CallbackName2(...)"

支持的回调类型：
1. OnInputFuzzy: 模糊搜索回调（查询数据源）
   - delay: 搜索延迟(ms)，默认300
   - min: 最小搜索长度，默认1
   - max_results: 最大返回结果数，默认50
   - cache: 缓存时间(秒)，默认0（不缓存）
   - filter_mode: 过滤模式（contains/startswith/fuzzy/fulltext），默认contains
   - fields: 搜索字段，多个用逗号分隔（name,code,description）
   - sort_by: 排序方式（relevance/popularity/time/name），默认relevance
   - hierarchy: 是否层级搜索（true/false），默认false
   - cascade: 是否级联搜索（true/false），默认false
   - highlight: 是否高亮匹配（true/false），默认false
   - scope: 搜索范围（table,column/user,department等）

2. OnInputValidate: 实时验证回调（需要后端交互验证）
   - delay: 验证延迟(ms)，默认500
   - min: 最小验证长度，默认1
   - real_time: 实时验证（true/false），默认true
   - show_loading: 显示加载状态（true/false），默认true
   - cache: 缓存验证结果时间(秒)，默认60
   - retry: 失败重试次数，默认1

3. OnChange: 值变化回调
   - debounce: 防抖延迟(ms)
   - trigger: 触发器名称
   - cascade: 是否级联（true/false）

工作模式：
1. 仅配置options：固定选项模式
2. 仅配置callback标签：纯搜索模式
3. 同时配置options和callback标签：混合模式（预设+搜索）

数据格式：
- 固定选项：使用options中定义的格式
- 搜索结果：回调返回["value(label)", "value2(label2)"]格式
- 存储格式：[]string（解析后的value数组）
- default_value必须与options中的格式完全一致

回调解析示例：
```go
// 单个回调
callback="OnInputFuzzy(delay:300,min:2,max_results:20)"

// 多个回调
callback="OnInputFuzzy(delay:300,min:2);OnInputValidate(real_time:true,cache:60)"

// 复杂回调
callback="OnInputFuzzy(delay:300,min:2,cache:60,filter_mode:fuzzy);OnInputValidate(delay:500,show_loading:true);OnChange(debounce:200)"
```

解析算法：
1. 读取callback标签：field.Tag.Get("callback")
2. 按分号分割回调：split(";")
3. 提取回调名和参数：regex `(\w+)\(([^)]*)\)`
4. 按逗号分割参数：split(",")
5. 解析键值对：split(":")
*/

// ===== 回调解析实现 =====

/*
回调解析伪代码：

```go
// 解析回调字符串
func parseCallbacks(callbackStr string) map[string]map[string]string {
    result := make(map[string]map[string]string)

         // 1. 按分号分割回调
     callbacks := strings.Split(callbackStr, ";")

    for _, callback := range callbacks {
        callback = strings.TrimSpace(callback)

        // 2. 使用正则提取回调名和参数
        re := regexp.MustCompile(`(\w+)\(([^)]*)\)`)
        matches := re.FindStringSubmatch(callback)

        if len(matches) == 3 {
            callbackName := matches[1]
            paramsStr := matches[2]

                         // 3. 解析参数
             params := make(map[string]string)
             if paramsStr != "" {
                 paramPairs := strings.Split(paramsStr, ",")
                 for _, pair := range paramPairs {
                     kv := strings.Split(pair, ":")
                     if len(kv) == 2 {
                         params[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
                     }
                 }
             }

            result[callbackName] = params
        }
    }

    return result
}

// 使用示例
callbackTag := field.Tag.Get("callback") // "OnInputFuzzy(delay:300,min:2);OnInputValidate(real_time:true,cache:60)"
callbacks := parseCallbacks(callbackTag)

// 结果：
// {
//   "OnInputFuzzy": {"delay": "300", "min": "2"},
//   "OnInputValidate": {"real_time": "true", "cache": "60"}
// }
```

解析优势：
1. 无冲突：分号不会与逗号冲突
2. 表达力强：支持任意参数和回调组合
3. 易于扩展：新增回调类型和参数无需修改解析逻辑
4. 向后兼容：可以逐步增加新的回调类型
5. 错误处理：解析失败时可以给出明确的错误信息

测试用例：
```go
// 基础测试
"OnInputFuzzy(delay:300)" → {"OnInputFuzzy": {"delay": "300"}}

// 多参数测试
"OnInputFuzzy(delay:300,min:2,max_results:20)" → {"OnInputFuzzy": {"delay": "300", "min": "2", "max_results": "20"}}

// 多回调测试
"OnInputFuzzy(delay:300);OnInputValidate(real_time:true)" → {"OnInputFuzzy": {"delay": "300"}, "OnInputValidate": {"real_time": "true"}}

// 空参数测试
"OnChange()" → {"OnChange": {}}

// 复杂测试
"OnInputFuzzy(delay:300,min:2,cache:60);OnInputValidate(delay:500,show_loading:true);OnChange(debounce:200)"
```
*/

// ===== 实现要点 =====

/*
前端实现要点：

1. 模式识别：
   - 有callback_type配置：启用搜索功能
   - 有options配置：显示预设选项
   - 同时有两者：混合模式

2. 选项处理：
   - 解析options为固定选项列表
   - 解析搜索返回的"value(label)"格式
   - 合并固定选项和搜索结果（混合模式）

3. 搜索功能：
   - 检测到callback_type:OnInputFuzzy时启用搜索
   - 根据callback_delay设置防抖延迟
   - 根据callback_min_length设置最小搜索长度
   - 根据callback_max_results限制返回结果

4. 用户体验：
   - multiple_limit控制选择数量
   - collapse_tags控制标签显示
   - allow_create支持创建新条目

后端实现要点：

1. 配置解析：
   - 解析options为选项列表
   - 识别callback_type配置启用搜索
   - 验证default_value与options的一致性

2. 数据验证：
   - 验证选中值的有效性
   - 数量限制验证
   - 选项范围验证
   - default_value必须在options范围内

3. 模式支持：
   - 固定模式：验证值在options中
   - 搜索模式：通过回调验证
   - 混合模式：支持两种验证方式

4. 回调参数处理：
   - callback_type: 确定回调类型
   - callback_delay: 前端防抖延迟
   - callback_min_length: 搜索触发条件
   - callback_max_results: 结果数量限制
*/

// ===== 与其他组件的关系 =====

/*
组件对比：

Select（单选下拉）：
- 单选：只能选择一个值
- 配置：options + callback_type
- 存储：string
- type: string

MultiSelect（多选下拉）：
- 多选：可以选择多个值
- 配置：options + callback_type + multiple_limit
- 存储：[]string
- type: []string

Checkbox（多选复选框）：
- 多选：可以选择多个值
- 配置：仅options，固定选项
- 存储：[]string
- type: []string
- 显示：复选框组，所有选项可见

选择建议：
- 单个选择 → Select (type: string)
- 多选且选项<10个 → Checkbox (type: []string)
- 多选且选项>10个或动态 → MultiSelect (type: []string)
- 需要搜索功能 → MultiSelect (type: []string)
*/

// ===== 使用最佳实践 =====

/*
最佳实践建议：

1. 模式选择：
   - 选项固定且较少：使用options配置
   - 选项动态或较多：使用callback_type配置
   - 需要预设+搜索：使用options+callback_type

2. 配置建议：
   - 设置合理的multiple_limit
   - 提供清晰的placeholder
   - 考虑使用collapse_tags避免界面拥挤
   - 确保default_value与options完全对齐

3. 性能优化：
   - 大量选项时使用callback_type模式
   - 合理设置callback_delay防抖延迟
   - 使用callback_min_length避免无效搜索
   - 使用callback_max_results控制结果数量

4. 数据一致性：
   - default_value必须与options中的格式完全一致
   - 验证选中值的有效性
   - 处理特殊字符和编码问题

推荐配置示例：
```go
// 固定选项模式
Skills []string `runner:"code:skills;name:技能;type:[]string;widget:multiselect;options:js(JavaScript),py(Python),go(Go);default_value:js(JavaScript),py(Python)"`

// 搜索模式
Users []string `runner:"code:users;name:用户;type:[]string;widget:multiselect;placeholder:搜索用户;callback_type:OnInputFuzzy;callback_delay:300"`

// 混合模式
Categories []string `runner:"code:categories;name:分类;type:[]string;widget:multiselect;options:tech(技术),business(商业);callback_type:OnInputFuzzy;callback_delay:300"`

// 高级搜索配置
AdvancedSearch []string `runner:"code:advanced;name:高级搜索;type:[]string;widget:multiselect;callback_type:OnInputFuzzy;callback_delay:500;callback_min_length:3;callback_max_results:20"`
```
*/

// ===== 回调参数扩展设计 =====

/*
回调参数扩展方案：

当前支持的回调参数：
- callback_type: OnInputFuzzy（回调类型）
- callback_delay: 300（搜索延迟ms）
- callback_min_length: 2（最小搜索长度）
- callback_max_results: 50（最大返回结果数）

未来可扩展的回调参数：
- callback_cache_time: 缓存时间（秒）
- callback_debounce_mode: 防抖模式（leading/trailing）
- callback_filter_mode: 过滤模式（contains/startswith/fuzzy）
- callback_sort_field: 排序字段
- callback_sort_order: 排序方向（asc/desc）

扩展示例：
```go
// 完整回调配置
Users []string `runner:"code:users;name:用户;type:[]string;widget:multiselect;callback_type:OnInputFuzzy;callback_delay:300;callback_min_length:2;callback_max_results:20;callback_cache_time:60;callback_filter_mode:fuzzy"`
```

解析优势：
1. 每个参数独立，易于解析
2. 参数可选，向后兼容
3. 类型明确，验证简单
4. 扩展性强，不会冲突

对比原方案callback:OnInputFuzzy(delay:300,min:3)的问题：
1. 括号嵌套解析复杂
2. 参数顺序敏感
3. 扩展时容易冲突
4. 错误处理困难
*/

// ===== 实现优先级 =====

/*
实现步骤：

第一步：基础多选功能 (0.5天)
- 固定选项模式实现
- options解析和多选逻辑
- 基础验证和数量限制
- default_value与options对齐验证

第二步：搜索功能集成 (0.5天)
- callback_type配置识别
- 回调参数解析（delay、min_length等）
- 搜索结果解析和显示
- 搜索与固定选项的结合

第三步：用户体验优化 (0.5天)
- 标签折叠和显示优化
- 创建新条目支持
- 错误处理和提示
- 加载状态和空状态

第四步：测试和优化 (0.5天)
- 各种模式的测试
- 参数解析测试
- 性能优化
- 边界情况处理

预计总工期：2天
核心优势：配置清晰，参数独立，易于扩展，类型正确
*/
