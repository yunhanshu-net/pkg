package llm_train

// ValidationExamples 展示go-playground/validator库的常用验证规则
// 这些验证规则可以直接用在runner标签的validate字段中
type ValidationExamples struct {
	// ===== 必填验证 =====
	RequiredField string `validate:"required" runner:"code:required_field;name:必填字段;type:string" json:"required_field"` // 最基础的验证：字段不能为空

	// ===== 字符串长度验证 =====
	MinLength   string `validate:"min=3" runner:"code:min_length;name:最小长度3;type:string" json:"min_length"`
	MaxLength   string `validate:"max=100" runner:"code:max_length;name:最大长度100;type:string" json:"max_length"`
	RangeLength string `validate:"min=3,max=20" runner:"code:range_length;name:长度3-20;type:string" json:"range_length"`
	ExactLength string `validate:"len=6" runner:"code:exact_length;name:固定长度6;type:string" json:"exact_length"`

	// ===== 数值验证 =====
	MinValue    int     `validate:"min=0" runner:"code:min_value;name:最小值0;type:number" json:"min_value"`
	MaxValue    int     `validate:"max=100" runner:"code:max_value;name:最大值100;type:number" json:"max_value"`
	RangeValue  int     `validate:"min=1,max=150" runner:"code:range_value;name:数值1-150;type:number" json:"range_value"`
	PositiveNum float64 `validate:"gt=0" runner:"code:positive_num;name:大于0;type:number" json:"positive_num"`
	NonNegative float64 `validate:"gte=0" runner:"code:non_negative;name:大于等于0;type:number" json:"non_negative"`
	LessThan    float64 `validate:"lt=100" runner:"code:less_than;name:小于100;type:number" json:"less_than"`
	LessOrEqual float64 `validate:"lte=100" runner:"code:less_or_equal;name:小于等于100;type:number" json:"less_or_equal"`

	// ===== 格式验证 =====
	Email string `validate:"email" runner:"code:email;name:邮箱格式;type:string" json:"email"`
	URL   string `validate:"url" runner:"code:url;name:URL格式;type:string" json:"url"`
	UUID  string `validate:"uuid" runner:"code:uuid;name:UUID格式;type:string" json:"uuid"`
	IPv4  string `validate:"ipv4" runner:"code:ipv4;name:IPv4格式;type:string" json:"ipv4"`
	IPv6  string `validate:"ipv6" runner:"code:ipv6;name:IPv6格式;type:string" json:"ipv6"`
	MAC   string `validate:"mac" runner:"code:mac;name:MAC地址格式;type:string" json:"mac"`

	// ===== 字符类型验证 =====
	AlphaOnly     string `validate:"alpha" runner:"code:alpha_only;name:只能包含字母;type:string" json:"alpha_only"`
	AlphaNum      string `validate:"alphanum" runner:"code:alpha_num;name:只能包含字母数字;type:string" json:"alpha_num"`
	NumericOnly   string `validate:"numeric" runner:"code:numeric_only;name:只能包含数字;type:string" json:"numeric_only"`
	AlphaSpace    string `validate:"alphaunicode" runner:"code:alpha_space;name:字母和Unicode字符;type:string" json:"alpha_space"`
	AlphaNumSpace string `validate:"alphanumunicode" runner:"code:alpha_num_space;name:字母数字和Unicode字符;type:string" json:"alpha_num_space"`

	// ===== 枚举值验证 =====
	OneOfString string `validate:"oneof=red green blue" runner:"code:one_of_string;name:枚举值;widget:select;options:red(红色),green(绿色),blue(蓝色);type:string" json:"one_of_string"` // 字符串枚举：注意oneof的值要与options的value匹配
	OneOfNumber int    `validate:"oneof=1 2 3 4 5" runner:"code:one_of_number;name:枚举数字;widget:select;options:1(一),2(二),3(三),4(四),5(五);type:number" json:"one_of_number"`       // 数字枚举：数字类型的枚举验证

	// ===== 正则表达式验证 =====
	ChinesePhone  string `validate:"regexp=^1[3-9]\\d{9}$" runner:"code:chinese_phone;name:中国手机号;type:string;placeholder:请输入11位手机号" json:"chinese_phone"`
	ChineseIdCard string `validate:"regexp=^[1-9]\\d{5}(18|19|20)\\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$" runner:"code:chinese_id_card;name:中国身份证号;type:string;placeholder:请输入18位身份证号" json:"chinese_id_card"`
	PostalCode    string `validate:"regexp=^\\d{6}$" runner:"code:postal_code;name:邮政编码;type:string;placeholder:请输入6位邮政编码" json:"postal_code"`

	// ===== 字段比较验证 =====
	Password        string `validate:"required,min=6" runner:"code:password;name:密码;type:string" json:"password"`                              // 密码字段：用于字段比较验证的基准字段
	ConfirmPassword string `validate:"required,eqfield=Password" runner:"code:confirm_password;name:确认密码;type:string" json:"confirm_password"` // 确认密码：必须与Password字段值相等
	StartDate       string `validate:"required" runner:"code:start_date;name:开始日期;type:string" json:"start_date"`                              // 开始日期：用于日期比较验证的基准字段
	EndDate         string `validate:"required,gtfield=StartDate" runner:"code:end_date;name:结束日期;type:string" json:"end_date"`                // 结束日期：必须大于StartDate字段值

	// ===== 条件验证 =====
	OptionalEmail  string `validate:"omitempty,email" runner:"code:optional_email;name:可选邮箱;type:string" json:"optional_email"`
	RequiredIfTrue string `validate:"required_if=IsVip true" runner:"code:required_if_true;name:VIP时必填;type:string" json:"required_if_true"`
	RequiredUnless string `validate:"required_unless=UserType guest" runner:"code:required_unless;name:非访客必填;type:string" json:"required_unless"`
	IsVip          bool   `runner:"code:is_vip;name:是否VIP;type:boolean" json:"is_vip"`
	UserType       string `runner:"code:user_type;name:用户类型;type:string" json:"user_type"`

	// ===== 数组/切片验证 =====
	Tags    []string `validate:"min=1,max=5,dive,required,min=2,max=20" runner:"code:tags;name:标签列表;type:string" json:"tags"`
	Numbers []int    `validate:"min=1,max=10,dive,min=1,max=100" runner:"code:numbers;name:数字列表;type:string" json:"numbers"`

	// ===== 自定义错误消息示例 =====
	CustomMessage string `validate:"required,min=3,max=50" runner:"code:custom_message;name:自定义消息;type:string" json:"custom_message"`
}

// CommonValidationRules 常用验证规则组合示例
var CommonValidationRules = map[string]string{
	// 用户信息相关
	"username":  "required,min=3,max=20,alphanum",
	"password":  "required,min=6,max=50",
	"email":     "required,email",
	"phone":     "required,len=11,numeric",
	"real_name": "required,min=2,max=20",
	"id_card":   "omitempty,len=18",

	// 数值相关
	"age":        "required,min=1,max=150",
	"score":      "required,min=0,max=100",
	"price":      "required,min=0",
	"quantity":   "required,min=1,max=9999",
	"percentage": "required,min=0,max=100",

	// 字符串相关
	"title":       "required,max=200",
	"description": "max=1000",
	"content":     "required",
	"code":        "required,len=6,alphanum",
	"url":         "omitempty,url",

	// 选择器相关
	"status":   "required,oneof=active inactive",
	"priority": "required,oneof=low medium high",
	"category": "required,oneof=type1 type2 type3",
	"gender":   "required,oneof=male female other",

	// 时间相关
	"date":     "required",
	"datetime": "required",

	// 文件相关
	"file_size": "required,min=1,max=10485760", // 最大10MB
	"file_type": "required,oneof=jpg jpeg png gif pdf doc docx",
}

// ValidationRuleExplanations 验证规则说明
var ValidationRuleExplanations = map[string]string{
	// 基础验证
	"required":  "必填字段，不能为空",
	"omitempty": "如果为空则跳过后续验证",

	// 长度验证
	"min": "最小长度/值",
	"max": "最大长度/值",
	"len": "固定长度",

	// 数值比较
	"gt":  "大于指定值",
	"gte": "大于等于指定值",
	"lt":  "小于指定值",
	"lte": "小于等于指定值",

	// 格式验证
	"email": "邮箱格式验证",
	"url":   "URL格式验证",
	"uuid":  "UUID格式验证",
	"ipv4":  "IPv4地址格式",
	"ipv6":  "IPv6地址格式",
	"mac":   "MAC地址格式",

	// 字符类型
	"alpha":           "只能包含字母",
	"alphanum":        "只能包含字母和数字",
	"numeric":         "只能包含数字",
	"alphaunicode":    "字母和Unicode字符",
	"alphanumunicode": "字母数字和Unicode字符",

	// 枚举验证
	"oneof": "必须是指定值中的一个",

	// 正则表达式
	"regexp": "正则表达式验证",

	// 字段比较
	"eqfield":  "必须等于指定字段的值",
	"nefield":  "必须不等于指定字段的值",
	"gtfield":  "必须大于指定字段的值",
	"gtefield": "必须大于等于指定字段的值",
	"ltfield":  "必须小于指定字段的值",
	"ltefield": "必须小于等于指定字段的值",

	// 条件验证
	"required_if":      "当指定字段等于指定值时必填",
	"required_unless":  "除非指定字段等于指定值，否则必填",
	"required_with":    "当指定字段有值时必填",
	"required_without": "当指定字段无值时必填",

	// 数组验证
	"dive":   "深入验证数组/切片中的每个元素",
	"unique": "数组/切片中的元素必须唯一",
}

// 使用示例说明
/*
验证规则使用说明：

1. 基础用法：
   validate:"required"                    // 必填
   validate:"required,email"              // 必填且邮箱格式
   validate:"omitempty,url"               // 可选，但如果填写必须是URL格式

2. 长度验证：
   validate:"min=3"                       // 最小长度3
   validate:"max=100"                     // 最大长度100
   validate:"min=3,max=20"                // 长度在3-20之间
   validate:"len=6"                       // 固定长度6

3. 数值验证：
   validate:"min=0"                       // 最小值0
   validate:"max=100"                     // 最大值100
   validate:"min=1,max=150"               // 值在1-150之间
   validate:"gt=0"                        // 大于0
   validate:"gte=0"                       // 大于等于0

4. 格式验证：
   validate:"email"                       // 邮箱格式
   validate:"url"                         // URL格式
   validate:"numeric"                     // 纯数字
   validate:"alphanum"                    // 字母数字

5. 枚举验证：
   validate:"oneof=red green blue"        // 必须是red、green或blue之一
   validate:"oneof=1 2 3"                 // 必须是1、2或3之一

6. 正则表达式：
   validate:"regexp=^1[3-9]\\d{9}$"       // 中国手机号格式

7. 字段比较：
   validate:"eqfield=Password"            // 必须等于Password字段的值
   validate:"gtfield=StartDate"           // 必须大于StartDate字段的值

8. 条件验证：
   validate:"required_if=IsVip true"      // 当IsVip为true时必填
   validate:"required_unless=UserType guest" // 除非UserType为guest，否则必填

9. 组合使用：
   validate:"required,min=6,max=50"       // 必填，长度6-50
   validate:"omitempty,email"             // 可选，但如果填写必须是邮箱格式
   validate:"required,oneof=active inactive" // 必填，且必须是active或inactive

注意事项：
- 多个验证规则用逗号分隔
- 对于select组件，oneof的值必须与options中的value完全匹配
- 使用omitempty可以让字段变为可选
- 正则表达式中的反斜杠需要双重转义
*/
