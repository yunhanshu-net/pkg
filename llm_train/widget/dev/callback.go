package dev

// Callback标签设计文档
// 独立的回调配置标签，专门处理动态交互逻辑

// ===== 重要说明 =====
/*
Callback标签是完全独立的标签，与runner标签分离：

外层通用参数（在runner标签中定义）：
- Type: string (对应Go类型系统，如string、[]string等)
- Widget: string (组件类型，如multiselect、select等)
- Required: bool (是否必选)
- Code: string (字段标识)
- CnName: string (显示名称)

Callback标签职责：
- 专门处理动态交互逻辑
- 不包含任何通用参数
- 纯技术属性配置
- 与业务逻辑解耦

标签分工：
- runner标签：组件基础配置
- callback标签：交互行为配置
- validate标签：静态验证规则
*/

// ===== Callback标签概述 =====

/*
Callback标签是一个独立的标签，专门用于配置组件的动态交互行为。
它与runner标签分离，职责清晰，专注于处理需要后端交互的逻辑。

设计理念：
1. 职责分离：callback专注交互逻辑，runner专注组件配置
2. 技术属性：关注如何交互，而不是交互什么
3. 无业务耦合：参数都是技术层面的配置
4. 高度可扩展：支持任意数量的回调和参数
*/

// ===== 语法规则 =====

/*
基本语法：
callback="CallbackName(param1:value1,param2:value2)"

多回调语法：
callback="CallbackName1(param1:value1,param2:value2);CallbackName2(param3:value3)"

分隔符规则：
- 参数分隔符：逗号(,)
- 回调分隔符：分号(;)
- 键值分隔符：冒号(:)

语法特点：
1. 独立标签：不与runner标签冲突
2. 层次清晰：参数用逗号，回调用分号
3. 易于解析：正则表达式友好
4. 扩展性强：支持任意参数组合
*/

// ===== 支持的回调类型 =====

// CallbackType 回调类型定义
type CallbackType string

const (
	OnInputFuzzy    CallbackType = "OnInputFuzzy"    // 模糊搜索回调
	OnInputValidate CallbackType = "OnInputValidate" // 实时验证回调
	OnChange        CallbackType = "OnChange"        // 值变化回调
	OnFocus         CallbackType = "OnFocus"         // 获得焦点回调
	OnBlur          CallbackType = "OnBlur"          // 失去焦点回调
	OnSelect        CallbackType = "OnSelect"        // 选择项回调
	OnRemove        CallbackType = "OnRemove"        // 移除项回调
	OnClear         CallbackType = "OnClear"         // 清空回调
)

// ===== OnInputFuzzy 模糊搜索回调 =====

/*
OnInputFuzzy: 模糊搜索回调（查询数据源）

核心参数：
- delay: 搜索延迟(ms)，默认300
- min: 最小搜索长度，默认1
- max_results: 最大返回结果数，默认50

性能参数：
- cache: 缓存时间(秒)，默认0（不缓存）
- timeout: 请求超时(ms)，默认5000

搜索参数：
- filter_mode: 过滤模式（contains/startswith/fuzzy/fulltext），默认contains
- fields: 搜索字段，多个用逗号分隔（name,code,description）
- case_sensitive: 是否大小写敏感（true/false），默认false

结果参数：
- sort_by: 排序方式（relevance/popularity/time/name），默认relevance
- sort_order: 排序方向（asc/desc），默认asc
- highlight: 是否高亮匹配（true/false），默认false

结构参数：
- hierarchy: 是否层级搜索（true/false），默认false
- cascade: 是否级联搜索（true/false），默认false
- group_by: 分组字段，用于结果分组

使用示例：
```go
// 基础搜索
callback="OnInputFuzzy(delay:300,min:2)"

// 高性能搜索
callback="OnInputFuzzy(delay:200,min:1,max_results:20,cache:600)"

// 多字段搜索
callback="OnInputFuzzy(delay:300,min:2,fields:name,code,description,filter_mode:fuzzy)"

// 层级搜索
callback="OnInputFuzzy(delay:400,min:1,hierarchy:true,sort_by:name)"

// 全文搜索
callback="OnInputFuzzy(delay:500,min:3,filter_mode:fulltext,highlight:true,max_results:15)"
```
*/

// ===== OnInputValidate 实时验证回调 =====

/*
OnInputValidate: 实时验证回调（需要后端交互验证）

核心参数：
- delay: 验证延迟(ms)，默认500
- min: 最小验证长度，默认1
- timeout: 验证超时(ms)，默认3000

交互参数：
- real_time: 实时验证（true/false），默认true
- show_loading: 显示加载状态（true/false），默认true
- show_success: 显示成功状态（true/false），默认false

性能参数：
- cache: 缓存验证结果时间(秒)，默认60
- retry: 失败重试次数，默认1
- retry_delay: 重试延迟(ms)，默认1000

错误处理：
- silent_error: 静默错误（true/false），默认false
- error_position: 错误显示位置（inline/tooltip/bottom），默认inline

使用示例：
```go
// 基础验证
callback="OnInputValidate(delay:500,min:2)"

// 快速验证
callback="OnInputValidate(delay:300,min:1,show_loading:true)"

// 缓存验证
callback="OnInputValidate(delay:600,min:2,cache:300,retry:2)"

// 静默验证
callback="OnInputValidate(delay:400,min:2,show_loading:false,silent_error:true)"

// 实时验证
callback="OnInputValidate(delay:200,min:1,real_time:true,show_success:true)"
```
*/

// ===== OnChange 值变化回调 =====

/*
OnChange: 值变化回调（值发生变化时触发）

核心参数：
- debounce: 防抖延迟(ms)，默认300
- immediate: 是否立即触发（true/false），默认false

触发条件：
- trigger_on_add: 添加项时触发（true/false），默认true
- trigger_on_remove: 移除项时触发（true/false），默认true
- trigger_on_clear: 清空时触发（true/false），默认true
- trigger_on_sort: 排序时触发（true/false），默认false

数据传递：
- include_old_value: 包含旧值（true/false），默认false
- include_diff: 包含差异信息（true/false），默认false

级联控制：
- cascade: 是否级联触发（true/false），默认false
- cascade_delay: 级联延迟(ms)，默认100

使用示例：
```go
// 基础变化监听
callback="OnChange(debounce:300)"

// 立即触发
callback="OnChange(debounce:0,immediate:true)"

// 详细监听
callback="OnChange(debounce:500,include_old_value:true,include_diff:true)"

// 级联触发
callback="OnChange(debounce:200,cascade:true,cascade_delay:150)"

// 选择性触发
callback="OnChange(debounce:300,trigger_on_add:true,trigger_on_remove:false)"
```
*/

// ===== 其他回调类型 =====

/*
OnFocus: 获得焦点回调
- delay: 触发延迟(ms)，默认0
- once: 是否只触发一次（true/false），默认false

OnBlur: 失去焦点回调
- delay: 触发延迟(ms)，默认0
- validate: 是否触发验证（true/false），默认false

OnSelect: 选择项回调
- delay: 触发延迟(ms)，默认0
- include_item: 包含选择项信息（true/false），默认true

OnRemove: 移除项回调
- delay: 触发延迟(ms)，默认0
- confirm: 是否需要确认（true/false），默认false

OnClear: 清空回调
- delay: 触发延迟(ms)，默认0
- confirm: 是否需要确认（true/false），默认false
*/

// ===== 多回调组合示例 =====

/*
多回调组合可以实现复杂的交互逻辑：

// 搜索 + 验证
callback="OnInputFuzzy(delay:300,min:2);OnInputValidate(delay:600,min:2)"

// 搜索 + 验证 + 变化监听
callback="OnInputFuzzy(delay:300,min:2,cache:600);OnInputValidate(delay:500,show_loading:true);OnChange(debounce:200)"

// 完整交互链
callback="OnInputFuzzy(delay:300,min:2,max_results:20);OnInputValidate(delay:500,cache:120);OnChange(debounce:300,cascade:true);OnSelect(include_item:true)"

// 焦点管理 + 搜索
callback="OnFocus(delay:100);OnInputFuzzy(delay:300,min:1);OnBlur(validate:true)"

// 高级验证链
callback="OnInputValidate(delay:500,retry:2,cache:300);OnChange(debounce:400,include_diff:true);OnRemove(confirm:true)"
*/

// ===== 解析实现 =====

/*
解析算法伪代码：

```go
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
```

解析结果示例：
```go
// 输入
callback="OnInputFuzzy(delay:300,min:2);OnInputValidate(delay:500,cache:60)"

// 输出
{
    "OnInputFuzzy": {
        "delay": "300",
        "min": "2"
    },
    "OnInputValidate": {
        "delay": "500",
        "cache": "60"
    }
}
```
*/

// ===== 最佳实践 =====

/*
1. 参数选择：
   - 根据实际需求选择合适的参数
   - 避免过度配置，保持简洁
   - 优先使用默认值，只配置必要参数

2. 性能优化：
   - 合理设置delay避免频繁请求
   - 使用cache减少重复请求
   - 设置合适的timeout避免长时间等待

3. 用户体验：
   - 设置合适的min避免无效搜索
   - 使用show_loading提供视觉反馈
   - 合理设置max_results控制结果数量

4. 错误处理：
   - 设置retry处理网络问题
   - 使用silent_error控制错误显示
   - 合理设置timeout避免卡死

5. 组合使用：
   - 搜索场景：OnInputFuzzy + OnChange
   - 验证场景：OnInputValidate + OnBlur
   - 复杂场景：多回调组合使用

推荐配置：
```go
// 用户搜索（常用）
callback="OnInputFuzzy(delay:300,min:2,max_results:20,cache:600)"

// 实时验证（常用）
callback="OnInputValidate(delay:500,min:2,cache:120,show_loading:true)"

// 复杂交互（高级）
callback="OnInputFuzzy(delay:300,min:2,cache:600);OnInputValidate(delay:500,show_loading:true);OnChange(debounce:300)"
```
*/

// ===== 扩展设计 =====

/*
未来可扩展的回调类型：

1. OnLoad: 组件加载回调
   - 用于初始化数据加载
   - 参数：delay, cache, auto_load

2. OnScroll: 滚动回调
   - 用于无限滚动加载
   - 参数：threshold, delay, batch_size

3. OnError: 错误回调
   - 用于错误处理和上报
   - 参数：retry, report, fallback

4. OnSuccess: 成功回调
   - 用于成功状态处理
   - 参数：delay, message, auto_hide

5. OnTimeout: 超时回调
   - 用于超时处理
   - 参数：timeout, retry, fallback

扩展示例：
```go
// 完整生命周期
callback="OnLoad(delay:100);OnInputFuzzy(delay:300,min:2);OnError(retry:2);OnSuccess(message:true)"

// 无限滚动
callback="OnInputFuzzy(delay:300,min:2);OnScroll(threshold:80,batch_size:20)"
```

设计原则：
1. 保持技术属性导向
2. 避免业务逻辑耦合
3. 确保向后兼容
4. 维护解析简洁性
*/

// ===== 总结 =====

/*
Callback标签的核心价值：

1. 职责分离：专注交互逻辑，与组件配置分离
2. 技术导向：关注如何交互，而不是交互什么
3. 高度灵活：支持任意回调和参数组合
4. 易于扩展：新增回调类型无需修改解析逻辑
5. 解析简单：清晰的分隔符规则，正则友好

使用建议：
- 简单场景：使用单一回调
- 复杂场景：组合多个回调
- 性能敏感：合理配置缓存和延迟
- 用户体验：提供适当的视觉反馈

Callback标签是组件动态交互的核心，通过合理配置可以实现丰富的用户交互体验。
*/
