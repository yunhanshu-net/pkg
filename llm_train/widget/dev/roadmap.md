# Widget 实现路线图

## 总体规划

基于现有的 `input` 和 `select` 组件，我们计划按优先级逐步实现其他组件。

## 实现阶段

### 阶段一：核心组件 (P0 - 高优先级)

#### 1. Switch 开关组件
- **预计工期**: 2-3天
- **技术难点**: 布尔值和字符串值的处理
- **实现要点**:
  - 支持 `bool` 和 `string` 两种数据类型
  - 自定义开启/关闭文本和值
  - 与现有验证规则集成

```go
// 目标API
IsEnabled bool `runner:"code:is_enabled;name:是否启用;widget:switch;on_text:启用;off_text:禁用"`
Status string `runner:"code:status;name:状态;widget:switch;on_value:active;off_value:inactive"`
```

#### 2. Number 数字输入组件
- **预计工期**: 3-4天
- **技术难点**: 精度控制、步长验证
- **实现要点**:
  - 支持整数和浮点数
  - 最小值、最大值、步长控制
  - 小数位数精度控制
  - 单位显示

```go
// 目标API
Age int `runner:"code:age;name:年龄;widget:number;min:1;max:150;step:1;unit:岁"`
Price float64 `runner:"code:price;name:价格;widget:number;min:0;step:0.01;precision:2;unit:元"`
```

#### 3. Checkbox 复选框组件
- **预计工期**: 4-5天
- **技术难点**: 多选值的处理和验证
- **实现要点**:
  - 支持 `[]string` 和 `string` (逗号分隔) 两种数据格式
  - 最少/最多选择数量控制
  - 默认选中值处理

```go
// 目标API
Permissions []string `runner:"code:permissions;name:权限;widget:checkbox;options:read(读取),write(写入);min_select:1"`
```

#### 4. Radio 单选框组件
- **预计工期**: 2-3天
- **技术难点**: 与select组件的差异化
- **实现要点**:
  - 水平/垂直排列
  - 默认值设置
  - 样式区分

```go
// 目标API
Gender string `runner:"code:gender;name:性别;widget:radio;options:male(男),female(女);direction:horizontal"`
```

### 阶段二：扩展组件 (P1 - 中优先级)

#### 5. Date 日期选择组件
- **预计工期**: 5-6天
- **技术难点**: 日期格式处理、时区问题
- **实现要点**:
  - 多种日期格式支持
  - 日期范围限制
  - 时间选择支持
  - 与Go time.Time类型集成

```go
// 目标API
Birthday string `runner:"code:birthday;name:生日;widget:date;format:YYYY-MM-DD"`
EventTime string `runner:"code:event_time;name:活动时间;widget:date;show_time:true"`
```

#### 6. File 文件上传组件
- **预计工期**: 6-8天
- **技术难点**: 文件上传处理、类型验证
- **实现要点**:
  - 文件类型限制
  - 文件大小限制
  - 多文件上传支持
  - 上传进度显示
  - 文件预览功能

```go
// 目标API
Avatar string `runner:"code:avatar;name:头像;widget:file;accept:.jpg,.png;max_size:2097152"`
Documents []string `runner:"code:documents;name:文档;widget:file;multiple:true"`
```

### 阶段三：高级组件 (P2 - 低优先级)

#### 7. Cascader 级联选择组件
- **预计工期**: 8-10天
- **技术难点**: 级联数据结构、异步加载
- **实现要点**:
  - 多级数据结构支持
  - 异步数据加载
  - 自定义分隔符
  - 搜索功能

#### 8. TreeSelect 树形选择组件
- **预计工期**: 10-12天
- **技术难点**: 树形数据处理、节点选择逻辑
- **实现要点**:
  - 树形数据展示
  - 单选/多选支持
  - 复选框模式
  - 节点搜索

#### 9. RichText 富文本编辑器
- **预计工期**: 15-20天
- **技术难点**: 富文本编辑器集成、内容安全
- **实现要点**:
  - 富文本编辑功能
  - 工具栏自定义
  - 内容安全过滤
  - 图片上传集成

## 实现策略

### 1. 渐进式开发
- 每个组件独立开发和测试
- 完成一个组件后立即集成到 overview.go
- 保持向后兼容性

### 2. 测试驱动
- 每个组件都要有完整的测试用例
- 包括单元测试和集成测试
- 验证规则测试

### 3. 文档同步
- 组件实现的同时更新文档
- 提供使用示例和最佳实践
- 更新 llm_train 训练数据

### 4. 性能考虑
- 组件懒加载
- 大数据量优化
- 移动端适配

## 技术架构

### 前端架构
```
components/
├── base/           # 基础组件
├── form/           # 表单组件
│   ├── Input.vue
│   ├── Select.vue
│   ├── Switch.vue   # 新增
│   ├── Number.vue   # 新增
│   └── ...
└── advanced/       # 高级组件
```

### 后端架构
```
pkg/
├── widget/         # 组件定义
├── parser/         # 标签解析
├── validator/      # 验证器
└── generator/      # 代码生成
```

## 里程碑

### Milestone 1: 核心组件完成 (4周)
- ✅ Input (已完成)
- ✅ Select (已完成)
- 🚧 Switch
- 🚧 Number
- 🚧 Checkbox
- 🚧 Radio

### Milestone 2: 扩展组件完成 (6周)
- 📋 Date
- 📋 File

### Milestone 3: 高级组件完成 (12周)
- 📋 Cascader
- 📋 TreeSelect
- 📋 RichText

## 风险评估

### 高风险项
1. **File组件**: 文件上传涉及安全性和性能问题
2. **RichText组件**: 复杂度高，第三方依赖多
3. **TreeSelect组件**: 大数据量性能问题

### 中风险项
1. **Date组件**: 时区和格式兼容性问题
2. **Cascader组件**: 异步数据加载复杂度

### 低风险项
1. **Switch组件**: 相对简单，风险较低
2. **Number组件**: 基于现有input组件扩展
3. **Checkbox/Radio组件**: 逻辑相对简单

## 资源需求

### 开发资源
- 前端开发: 1人 × 12周
- 后端开发: 1人 × 8周
- 测试: 1人 × 4周

### 技术资源
- Vue3 + TypeScript
- Go 1.19+
- 测试框架和工具
- CI/CD 环境

## 成功标准

### 功能标准
- 所有组件功能完整实现
- 通过完整的测试用例
- 性能满足要求

### 质量标准
- 代码覆盖率 > 80%
- 无严重安全漏洞
- 用户体验良好

### 文档标准
- 完整的API文档
- 使用示例和最佳实践
- 故障排除指南

## 后续规划

### 版本规划
- v1.0: 核心组件 (Switch, Number, Checkbox, Radio)
- v1.1: 扩展组件 (Date, File)
- v2.0: 高级组件 (Cascader, TreeSelect, RichText)

### 持续改进
- 用户反馈收集
- 性能优化
- 新组件需求评估
- 技术债务清理 