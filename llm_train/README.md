# LLM训练数据 - 验证规则示例

本目录包含了function-go框架中验证规则的完整示例，用于大模型学习和参考。

## 目录结构

```
pkg/llm_train/
├── table/
│   └── overview.go              # Table函数完整示例
├── form/
│   └── overview.go              # Form函数完整示例
├── validation_examples.go       # 通用验证规则示例
├── chinese_validation_examples.go # 中国本土化验证规则
└── README.md                    # 本文件
```

## 验证规则概述

基于 `go-playground/validator` 库，支持前后端一致性验证。

### 基础验证规则

| 规则 | 说明 | 示例 |
|------|------|------|
| `required` | 必填字段 | `validate:"required"` |
| `omitempty` | 可选字段（为空时跳过验证） | `validate:"omitempty,email"` |

### 长度验证

| 规则 | 说明 | 示例 |
|------|------|------|
| `min=3` | 最小长度3 | `validate:"min=3"` |
| `max=100` | 最大长度100 | `validate:"max=100"` |
| `len=6` | 固定长度6 | `validate:"len=6"` |
| `min=3,max=20` | 长度范围3-20 | `validate:"min=3,max=20"` |

### 数值验证

| 规则 | 说明 | 示例 |
|------|------|------|
| `min=0` | 最小值0 | `validate:"min=0"` |
| `max=100` | 最大值100 | `validate:"max=100"` |
| `gt=0` | 大于0 | `validate:"gt=0"` |
| `gte=0` | 大于等于0 | `validate:"gte=0"` |
| `lt=100` | 小于100 | `validate:"lt=100"` |
| `lte=100` | 小于等于100 | `validate:"lte=100"` |

### 格式验证

| 规则 | 说明 | 示例 |
|------|------|------|
| `email` | 邮箱格式 | `validate:"email"` |
| `url` | URL格式 | `validate:"url"` |
| `numeric` | 纯数字 | `validate:"numeric"` |
| `alphanum` | 字母数字 | `validate:"alphanum"` |
| `alpha` | 纯字母 | `validate:"alpha"` |
| `uuid` | UUID格式 | `validate:"uuid"` |
| `ipv4` | IPv4地址 | `validate:"ipv4"` |
| `ipv6` | IPv6地址 | `validate:"ipv6"` |

### 枚举验证

| 规则 | 说明 | 示例 |
|------|------|------|
| `oneof=value1 value2` | 必须是指定值之一 | `validate:"oneof=red green blue"` |

**注意**：对于select组件，oneof的值必须与options中的value完全匹配。

### 正则表达式验证

| 规则 | 说明 | 示例 |
|------|------|------|
| `regexp=^pattern$` | 自定义正则验证 | `validate:"regexp=^1[3-9]\\d{9}$"` |

### 字段比较验证

| 规则 | 说明 | 示例 |
|------|------|------|
| `eqfield=FieldName` | 必须等于指定字段 | `validate:"eqfield=Password"` |
| `nefield=FieldName` | 必须不等于指定字段 | `validate:"nefield=Username"` |
| `gtfield=FieldName` | 必须大于指定字段 | `validate:"gtfield=StartDate"` |
| `ltfield=FieldName` | 必须小于指定字段 | `validate:"ltfield=EndDate"` |

### 条件验证

| 规则 | 说明 | 示例 |
|------|------|------|
| `required_if=Field value` | 当指定字段等于指定值时必填 | `validate:"required_if=IsVip true"` |
| `required_unless=Field value` | 除非指定字段等于指定值，否则必填 | `validate:"required_unless=UserType guest"` |
| `required_with=Field` | 当指定字段有值时必填 | `validate:"required_with=Phone"` |
| `required_without=Field` | 当指定字段无值时必填 | `validate:"required_without=Email"` |

## 常用验证规则组合

### 用户信息
```go
Username string `validate:"required,min=3,max=20,alphanum"`
Password string `validate:"required,min=6,max=50"`
Email    string `validate:"required,email"`
Phone    string `validate:"required,len=11,numeric"`
Age      int    `validate:"required,min=1,max=150"`
```

### 业务数据
```go
Title       string  `validate:"required,max=200"`
Description string  `validate:"max=1000"`
Amount      float64 `validate:"min=0"`
Score       float64 `validate:"required,min=0,max=100"`
Status      string  `validate:"required,oneof=active inactive"`
```

### 可选字段
```go
Website     string `validate:"omitempty,url"`
OptionalAge int    `validate:"omitempty,min=1,max=150"`
```

## 中国本土化验证规则

### 个人信息
```go
// 手机号：11位，1开头，第二位3-9
MobilePhone string `validate:"required,regexp=^1[3-9]\\d{9}$"`

// 身份证号：18位，包含校验位
IdCard string `validate:"required,regexp=^[1-9]\\d{5}(18|19|20)\\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$"`

// 中文姓名：2-20个中文字符
ChineseName string `validate:"required,regexp=^[\\p{Han}]{2,20}$"`

// 邮政编码：6位数字
PostalCode string `validate:"required,regexp=^\\d{6}$"`
```

### 银行卡相关
```go
// 银行卡号：16-19位数字
BankCard string `validate:"required,regexp=^\\d{16,19}$"`
```

### 企业信息
```go
// 统一社会信用代码：18位字母数字组合
UnifiedSocialCreditCode string `validate:"required,regexp=^[0-9A-HJ-NPQRTUWXY]{2}\\d{6}[0-9A-HJ-NPQRTUWXY]{10}$"`

// 营业执照注册号：15位数字
BusinessLicenseNumber string `validate:"required,regexp=^\\d{15}$"`
```

### 车辆相关
```go
// 车牌号：支持新能源和传统车牌
LicensePlate string `validate:"required,regexp=^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领][A-Z][A-HJ-NP-Z0-9]{4,5}[A-HJ-NP-Z0-9挂学警港澳]$"`

// 车架号：17位字母数字，不含I、O、Q
VinCode string `validate:"required,regexp=^[A-HJ-NPR-Z0-9]{17}$"`
```

## 在Runner标签中使用

### 基本语法
```go
type User struct {
    Email    string `validate:"required,email" runner:"code:email;name:邮箱;type:string;placeholder:请输入邮箱地址"`
    Phone    string `validate:"required,len=11,numeric" runner:"code:phone;name:手机号;type:string;placeholder:请输入11位手机号"`
    Age      int    `validate:"required,min=1,max=150" runner:"code:age;name:年龄;type:number;placeholder:请输入年龄"`
    Status   string `validate:"required,oneof=active(启用) inactive(禁用)" runner:"code:status;name:状态;widget:select;options:active(启用),inactive(禁用);type:string"`
}
```

### 与Select组件配合
```go
// 正确：oneof的值与options完全匹配
Category string `validate:"required,oneof=type1(类型1) type2(类型2) type3(类型3)" runner:"code:category;name:分类;widget:select;default_value:type1(类型1);options:type1(类型1),type2(类型2),type3(类型3);type:string"`

// 错误：oneof的值与options不匹配
Category string `validate:"required,oneof=type1 type2 type3" runner:"code:category;name:分类;widget:select;options:type1(类型1),type2(类型2),type3(类型3);type:string"`
```

## 注意事项

1. **正则表达式转义**：在Go字符串中，反斜杠需要双重转义
   ```go
   // 正确
   validate:"regexp=^1[3-9]\\d{9}$"
   // 错误
   validate:"regexp=^1[3-9]\d{9}$"
   ```

2. **Select组件枚举值匹配**：oneof的值必须与options中的value完全匹配
   ```go
   // 正确
   validate:"oneof=active(启用) inactive(禁用)"
   runner:"options:active(启用),inactive(禁用)"
   
   // 错误
   validate:"oneof=active inactive"
   runner:"options:active(启用),inactive(禁用)"
   ```

3. **中文字符匹配**：使用`\p{Han}`匹配中文字符
   ```go
   validate:"regexp=^[\\p{Han}]{2,20}$"
   ```

4. **可选字段**：使用`omitempty`让字段变为可选
   ```go
   validate:"omitempty,email"  // 可选邮箱
   validate:"omitempty,url"    // 可选URL
   ```

5. **组合验证**：多个规则用逗号分隔
   ```go
   validate:"required,min=6,max=50"
   validate:"omitempty,email"
   validate:"required,oneof=red green blue"
   ```

## 最佳实践

1. **前后端一致性**：使用相同的验证规则确保前后端验证一致
2. **用户体验**：提供清晰的错误提示和placeholder
3. **安全性**：对敏感信息进行额外的安全验证
4. **本土化**：根据业务需求使用合适的本土化验证规则
5. **性能考虑**：复杂正则表达式可能影响性能，需要权衡
6. **测试覆盖**：为验证规则编写充分的测试用例

## 参考文档

- [go-playground/validator官方文档](https://github.com/go-playground/validator)
- [function-go框架文档](../README.md)
- [Runner标签规范](../runner/README.md) 