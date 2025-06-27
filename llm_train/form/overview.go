// Package form 这个文件是用来展示整个form函数的全部功能和能力，就像element-plus这种ui组件的组件一览表一样，这里没有太多逻辑，
//只是展示一下怎么用，函数的签名和组件的参数，标签，回调等等，目前项目处于开发中，只展示已经实现的功能，不展示未实现的功能

package form

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yunhanshu-net/function-go/pkg/dto/response" //这个是渲染返回的数据结构的包，必须引用这个包
	"github.com/yunhanshu-net/function-go/pkg/dto/usercall" //这是回调的请求和响应model，用到回调的请引用这个包
	"github.com/yunhanshu-net/function-go/runner"           //这个是function-go的核心包，需要引用这个才能启动
	"github.com/yunhanshu-net/pkg/logger"                   //这个是日志包，需要引用这个才能使用日志
)

// ===== 请求参数结构体 =====

// ExampleFormReq 展示form函数支持的所有输入参数类型,(没有标注widget的都默认是input)
type ExampleFormReq struct {
	// 基础字符串字段
	Title       string `json:"title" form:"title" runner:"code:title;name:标题;type:string;placeholder:请输入标题;example:产品发布会" validate:"required,max=200"`                                       // 标题验证：必填 + 最大长度200，示例值帮助理解字段用途
	Description string `json:"description" form:"description" runner:"code:description;mode:text_area;name:描述;type:string;placeholder:请输入详细描述;example:这是一个关于新产品发布的详细描述" validate:"max=1000"` // 描述验证：可选字段 + 最大长度1000，示例值展示描述格式,这里text_area会渲染成多行文本框
	Content     string `json:"content" form:"content" runner:"code:content;name:内容;mode:text_area;type:string;placeholder:请输入内容主体;example:产品介绍：本产品具有创新性功能..." validate:"required"`           // 内容验证：必填字段，无长度限制，示例值展示内容格式，这里text_area会渲染成多行文本框

	// 常用验证规则示例
	Email       string `json:"email" form:"email" runner:"code:email;name:邮箱;type:string;placeholder:请输入邮箱地址;example:user@example.com" validate:"required,email"`                              // 邮箱验证：必填 + 邮箱格式校验，示例值展示邮箱格式
	Phone       string `json:"phone" form:"phone" runner:"code:phone;name:手机号;type:string;placeholder:请输入11位手机号;example:13812345678" validate:"required,len=11,numeric"`                       // 手机号验证：必填 + 固定11位长度 + 纯数字，示例值展示手机号格式
	Username    string `json:"username" form:"username" runner:"code:username;name:用户名;type:string;placeholder:请输入用户名;example:user123" validate:"required,min=3,max=20,alphanum"`              // 用户名验证：必填 + 长度3-20 + 字母数字组合，示例值展示用户名格式
	Password    string `json:"password" form:"password" runner:"code:password;name:密码;mode:password;type:string;placeholder:请输入密码;example:password123" validate:"required,min=6,max=50"`       // 密码验证：必填 + 长度6-50，示例值展示密码格式,这里password会渲染成密码框
	ConfirmPass string `json:"confirm_pass" form:"confirm_pass" runner:"code:confirm_pass;name:确认密码;type:string;placeholder:请再次输入密码;example:password123" validate:"required,eqfield=Password"` // 确认密码验证：必填 + 必须与Password字段相等（字段比较验证），示例值与密码一致
	Website     string `json:"website" form:"website" runner:"code:website;name:网站;type:string;placeholder:请输入网站地址;example:https://www.example.com" validate:"omitempty,url"`                  // 网站验证：可选字段 + URL格式校验，示例值展示URL格式

	// 数值验证规则示例
	Age    int     `json:"age" form:"age" runner:"code:age;name:年龄;type:number;placeholder:请输入年龄;example:25" validate:"required,min=1,max=150"`                // 年龄验证：必填 + 数值范围1-150，示例值展示合理年龄
	Score  float64 `json:"score" form:"score" runner:"code:score;name:分数;type:number;placeholder:请输入分数(0-100);example:85.5" validate:"required,min=0,max=100"` // 分数验证：必填 + 数值范围0-100，示例值展示小数分数
	Salary float64 `json:"salary" form:"salary" runner:"code:salary;name:薪资;type:number;placeholder:请输入薪资;example:8000.00" validate:"omitempty,min=0"`         // 薪资验证：可选字段 + 最小值0（非负数），示例值展示薪资格式

	// 字符串长度和格式验证
	Code       string `json:"code" form:"code" runner:"code:code;name:编码;type:string;placeholder:请输入6位编码;example:ABC123" validate:"required,len=6,alphanum"`                                       // 编码验证：必填 + 固定6位长度 + 字母数字组合，示例值展示编码格式
	IdCard     string `json:"id_card" form:"id_card" runner:"code:id_card;name:身份证号;type:string;placeholder:请输入18位身份证号;example:110101199001011234" validate:"omitempty,len=18"`                    // 身份证验证：可选字段 + 固定18位长度，示例值展示身份证格式
	BankCard   string `json:"bank_card" form:"bank_card" runner:"code:bank_card;name:银行卡号;type:string;placeholder:请输入银行卡号;example:6222021234567890123" validate:"omitempty,min=16,max=19,numeric"` // 银行卡验证：可选字段 + 长度16-19位 + 纯数字，示例值展示银行卡格式
	PostalCode string `json:"postal_code" form:"postal_code" runner:"code:postal_code;name:邮政编码;type:string;placeholder:请输入6位邮政编码;example:100000" validate:"omitempty,len=6,numeric"`              // 邮编验证：可选字段 + 固定6位长度 + 纯数字，示例值展示邮编格式

	// 枚举值验证
	Gender   string `json:"gender" form:"gender" runner:"code:gender;name:性别;widget:select;default_value:male(男);options:male(男),female(女),other(其他);type:string;example:male(男)" validate:"required,oneof=male(男) female(女) other(其他)"`                                    // 性别验证：必填 + 枚举值校验（oneof值必须与options完全匹配），示例值展示选择格式
	UserType string `json:"user_type" form:"user_type" runner:"code:user_type;name:用户类型;widget:select;default_value:normal(普通用户);options:normal(普通用户),vip(VIP用户),admin(管理员);type:string;example:normal(普通用户)" validate:"required,oneof=normal(普通用户) vip(VIP用户) admin(管理员)"` // 用户类型验证：必填 + 枚举值校验，展示select组件的正确用法，示例值与默认值一致

	// 选择器字段
	Category string `json:"category" form:"category" runner:"code:category;name:分类;widget:select;default_value:type1(类型1);options:type1(类型1),type2(类型2),type3(类型3);type:string;example:type1(类型1)" validate:"required,oneof=type1(类型1) type2(类型2) type3(类型3)"`     // 分类验证：必填 + 三选一枚举，default_value必须在options中，示例值与默认值一致
	Priority string `json:"priority" form:"priority" runner:"code:priority;name:优先级;widget:select;default_value:medium(中);options:low(低),medium(中),high(高);type:string;example:medium(中)" validate:"required,oneof=low(低) medium(中) high(高)"`                      // 优先级验证：必填 + 三级优先级枚举，示例值展示中等优先级
	Status   string `json:"status" form:"status" runner:"code:status;name:状态;widget:select;default_value:draft(草稿);options:draft(草稿),published(已发布),archived(已归档);type:string;example:draft(草稿)" validate:"required,oneof=draft(草稿) published(已发布) archived(已归档)"` // 状态验证：必填 + 三状态枚举（草稿/发布/归档），示例值展示草稿状态

	// 数值字段
	Amount     float64 `json:"amount" form:"amount" runner:"code:amount;name:金额;type:number;placeholder:请输入金额;default_value:0;example:1299.99" validate:"min=0"`                           // 金额验证：最小值0（非负数），有默认值，示例值展示金额格式
	Quantity   int     `json:"quantity" form:"quantity" runner:"code:quantity;name:数量;type:number;placeholder:请输入数量;default_value:1;example:5" validate:"min=1,max=9999"`                  // 数量验证：数值范围1-9999，有默认值，示例值展示合理数量
	Percentage float64 `json:"percentage" form:"percentage" runner:"code:percentage;name:百分比;type:number;placeholder:请输入百分比(0-100);default_value:0;example:75.5" validate:"min=0,max=100"` // 百分比验证：数值范围0-100，有默认值，示例值展示百分比格式

	// 布尔字段
	IsPublic  bool   `json:"is_public" form:"is_public" runner:"code:is_public;name:是否公开;widget:switch;true_label:公开;false_label:私有;type:boolean;default_value:false;example:false"`                                     // 原生布尔类型：直接使用bool类型，有默认值，示例值与默认值一致
	IsActive  string `json:"is_active" form:"is_active" runner:"code:is_active;name:是否启用;widget:select;default_value:true(启用);options:true(启用),false(禁用);example:true(启用)" validate:"required,oneof=true(启用) false(禁用)"` // 布尔选择器：用select实现布尔选择，必填，示例值展示启用状态
	SendEmail bool   `json:"send_email" form:"send_email" runner:"code:send_email;name:发送邮件通知;widget:switch;true_label:发送;false_label:不发送;type:boolean;default_value:true;example:true"`                                 // 布尔开关：邮件通知开关，有默认值，示例值展示开启状态

	// 多选字段
	Tags     string `json:"tags" form:"tags" runner:"code:tags;name:标签;type:string;placeholder:多个标签用逗号分隔;example:技术,产品,创新"`                     // 标签字段：多个标签用逗号分隔，无验证规则（可选），示例值展示标签格式
	Keywords string `json:"keywords" form:"keywords" runner:"code:keywords;name:关键词;type:string;placeholder:多个关键词用空格分隔;example:人工智能 机器学习 深度学习"` // 关键词字段：多个关键词用空格分隔，无验证规则（可选），示例值展示关键词格式
}

// ===== 响应结果结构体 =====

// ExampleFormResp 展示form函数支持的所有输出结果类型
type ExampleFormResp struct {
	// 基础信息字段
	Message     string `json:"message" runner:"code:message;name:处理结果;type:string;example:表单处理成功"`                        // 处理结果消息：展示操作结果信息，示例值展示成功消息
	Status      string `json:"status" runner:"code:status;name:状态;type:string;example:success"`                           // 处理状态：成功/失败/处理中等状态，示例值展示成功状态
	ProcessTime string `json:"process_time" runner:"code:process_time;name:处理时间;type:string;example:2024-01-15 14:30:25"` // 处理时间：记录操作完成的时间，示例值展示时间格式

	// 数值结果字段
	TotalCount  int     `json:"total_count" runner:"code:total_count;name:总数量;type:number;example:150"`        // 总数量：处理的记录总数，示例值展示合理数量
	SuccessRate float64 `json:"success_rate" runner:"code:success_rate;name:成功率;type:number;example:95.5"`     // 成功率：操作成功的百分比，示例值展示高成功率
	TotalAmount float64 `json:"total_amount" runner:"code:total_amount;name:总金额;type:number;example:12580.50"` // 总金额：计算得出的金额总和，示例值展示金额格式

	// 复杂数据字段
	ResultData string `json:"result_data" runner:"code:result_data;name:结果数据;type:string;example:{\"processed\":150,\"failed\":7}"` // 结果数据：JSON格式的处理结果，示例值展示JSON格式
	ErrorInfo  string `json:"error_info" runner:"code:error_info;name:错误信息;type:string;example:部分数据格式不正确"`                          // 错误信息：记录处理过程中的错误，示例值展示错误描述
	LogInfo    string `json:"log_info" runner:"code:log_info;name:日志信息;type:string;example:开始处理->数据验证->业务处理->完成"`                   // 日志信息：详细的处理日志，示例值展示日志格式

	// 文件结果字段
	ReportUrl string `json:"report_url" runner:"code:report_url;name:报告文件;type:string;example:https://example.com/reports/20240115_report.pdf"` // 报告文件：生成的报告文件下载链接，示例值展示文件URL格式
	ExportUrl string `json:"export_url" runner:"code:export_url;name:导出文件;type:string;example:https://example.com/exports/20240115_data.xlsx"`  // 导出文件：生成的导出文件下载链接，示例值展示导出文件URL

	// 布尔结果字段
	IsSuccess  bool `json:"is_success" runner:"code:is_success;name:是否成功;true_label:成功;false_label:失败;type:boolean;example:true"` // 是否成功：操作是否成功完成，示例值展示成功状态
	HasWarning bool `json:"has_warning" runner:"code:has_warning;name:有警告;true_label:有;false_label:无;type:boolean;example:true"`  // 有警告：是否存在警告信息，示例值展示有警告情况

	// 时间结果字段
	CreatedAt   string `json:"created_at" runner:"code:created_at;name:创建时间;type:string;example:2024-01-15 14:25:10"`     // 创建时间：记录创建的时间，示例值展示创建时间格式
	CompletedAt string `json:"completed_at" runner:"code:completed_at;name:完成时间;type:string;example:2024-01-15 14:30:25"` // 完成时间：记录完成的时间，示例值展示完成时间格式
}

// ===== 函数配置 =====

var ExampleFormConfig = &runner.FunctionOptions{
	// 基础信息
	Tags:        []string{"示例", "表单处理", "数据处理"},
	EnglishName: "example_form",
	ChineseName: "示例表单处理",
	ApiDesc:     "展示form函数的完整功能，包括各种输入组件、数据处理、结果展示等操作",

	// 请求响应结构
	Request:  &ExampleFormReq{},
	Response: &ExampleFormResp{},

	// 渲染类型
	RenderType: response.RenderTypeForm, // 标识返回的数据是form格式

	// 数据库表管理（可选，form函数不一定需要数据库操作）
	CreateTables: []interface{}{
			//&ExampleForm{},
	},             // 如果需要创建表，在这里添加
	OperateTables: map[interface{}][]runner.OperateTableType{}, // 操作表的映射

	// 页面加载回调
	OnPageLoad: func(ctx *runner.Context, resp response.Response) (initData *usercall.OnPageLoadResp, err error) {
		// 页面加载时执行的函数，可以做初始化操作

		// 返回初始化的表单数据和结果
		resp.Form(&ExampleFormResp{
			Message:     "请填写表单信息",
			Status:      "等待处理",
			ProcessTime: time.Now().Format("2006-01-02 15:04:05"),
			TotalCount:  0,
			SuccessRate: 0.0,
			IsSuccess:   false,
			HasWarning:  false,
		}).Build()

		// 返回初始化数据，设置默认表单值
		initData = &usercall.OnPageLoadResp{
			Request: ExampleFormReq{
				Category:   "type1(类型1)",
				Priority:   "medium(中)",
				Status:     "draft(草稿)",
				IsActive:   "true(启用)",
				IsPublic:   false,
				SendEmail:  true,
				Quantity:   1,
				Amount:     0.0,
				Percentage: 0.0,
			},
			AutoRun: false, // 是否在页面加载时自动运行一次函数
		}
		return
	},

	// API创建回调
	OnApiCreated: func(ctx *runner.Context, req *usercall.OnApiCreatedReq) error {
		logger.Infof(ctx, "示例表单处理API创建成功: %+v", req)
		return nil
	},
}

// ===== 路由注册 =====

func init() {
	// 路由注册：/package名/函数英文名
	runner.Post("/form/example_form", ExampleForm, ExampleFormConfig)
}

// ===== 主要业务逻辑 =====

func ExampleForm(ctx *runner.Context, req *ExampleFormReq, resp response.Response) (err error) {
	// 记录处理开始时间
	startTime := time.Now()

	// 延迟函数，记录处理日志
	defer func() {
		logger.InfoContextf(ctx, "ExampleForm req:%+v err:%v", req, err)
	}()

	// 输入验证
	if req.Title == "" {
		return resp.Form(&ExampleFormResp{
			Message:     "标题不能为空",
			Status:      "验证失败",
			ProcessTime: time.Now().Format("2006-01-02 15:04:05"),
			IsSuccess:   false,
			HasWarning:  false,
			ErrorInfo:   "请填写标题字段",
		}).Build()
	}

	// 模拟业务处理逻辑
	var resultData map[string]interface{}
	var logMessages []string
	var hasWarning bool

	// 1. 数据处理
	logMessages = append(logMessages, "开始处理表单数据...")

	// 处理分类逻辑 - 修复后的判断逻辑
	switch req.Category {
	case "type1(类型1)":
		resultData = map[string]interface{}{
			"type":        "类型1处理",
			"description": "执行类型1的特殊处理逻辑",
		}
	case "type2(类型2)":
		resultData = map[string]interface{}{
			"type":        "类型2处理",
			"description": "执行类型2的特殊处理逻辑",
		}
		hasWarning = true
		logMessages = append(logMessages, "警告：类型2需要额外验证")
	case "type3(类型3)":
		resultData = map[string]interface{}{
			"type":        "类型3处理",
			"description": "执行类型3的特殊处理逻辑",
		}
	default:
		resultData = map[string]interface{}{
			"type":        "默认处理",
			"description": "执行默认处理逻辑",
		}
	}

	// 处理优先级逻辑 - 修复后的判断逻辑
	switch req.Priority {
	case "low(低)":
		logMessages = append(logMessages, "低优先级处理")
	case "medium(中)":
		logMessages = append(logMessages, "中等优先级处理")
	case "high(高)":
		logMessages = append(logMessages, "高优先级处理")
		hasWarning = true
		logMessages = append(logMessages, "警告：高优先级需要特别关注")
	}

	// 处理状态逻辑 - 修复后的判断逻辑
	switch req.Status {
	case "draft(草稿)":
		logMessages = append(logMessages, "草稿状态，暂不发布")
	case "published(已发布)":
		logMessages = append(logMessages, "已发布状态，执行发布流程")
	case "archived(已归档)":
		logMessages = append(logMessages, "已归档状态，只读模式")
	}

	// 2. 数值计算
	totalAmount := req.Amount * float64(req.Quantity)
	successRate := req.Percentage / 100.0

	logMessages = append(logMessages, fmt.Sprintf("计算总金额: %.2f", totalAmount))
	logMessages = append(logMessages, fmt.Sprintf("成功率: %.2f%%", req.Percentage))

	// 3. 组装结果数据
	resultData["tags"] = req.Tags
	resultData["keywords"] = req.Keywords
	resultData["settings"] = map[string]interface{}{
		"is_public":  req.IsPublic,
		"is_active":  req.IsActive,
		"send_email": req.SendEmail,
	}

	// 转换为JSON字符串
	resultJson, _ := json.MarshalIndent(resultData, "", "  ")

	// 计算处理时间
	processDuration := time.Since(startTime)
	logMessages = append(logMessages, fmt.Sprintf("处理完成，耗时: %v", processDuration))

	// 返回处理结果
	return resp.Form(&ExampleFormResp{
		Message:     "表单处理成功",
		Status:      "处理完成",
		ProcessTime: time.Now().Format("2006-01-02 15:04:05"),
		TotalCount:  req.Quantity,
		SuccessRate: successRate,
		TotalAmount: totalAmount,
		ResultData:  string(resultJson),
		ErrorInfo:   "",
		LogInfo:     fmt.Sprintf("处理日志:\n%s", joinStrings(logMessages, "\n")),
		ReportUrl:   "/reports/example_" + time.Now().Format("20060102150405") + ".pdf",
		ExportUrl:   "/exports/example_" + time.Now().Format("20060102150405") + ".xlsx",
		IsSuccess:   true,
		HasWarning:  hasWarning,
		CreatedAt:   startTime.Format("2006-01-02 15:04:05"),
		CompletedAt: time.Now().Format("2006-01-02 15:04:05"),
	}).Build()
}

// ===== 辅助函数 =====

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// ===== 功能说明 =====

/*
Form函数完整功能说明：

1. 请求参数 (Request)：
   - 支持所有Go基础类型：string, int, float64, bool
   - 支持各种UI组件：input, select
   - 支持验证规则：required, max, min, 自定义验证

2. Runner标签说明（已实现功能）：
   - code: 字段代码，用于前端识别
   - name: 字段显示名称
   - type: 字段类型 (string, number, boolean)
   - widget: UI组件 (select, input)
   - placeholder: 输入提示
   - default_value: 默认值 (必须与options中的某个值完全匹配)
   - options: 选择器选项 (格式：value(label),value2(label2))

   - validate: 验证规则 (required, max, min等)

3. 响应结果 (Response)：
   - 支持各种结果展示：文本、数值、布尔、时间、文件链接
   - 支持结果展示：用于展示处理结果和状态信息
   - 支持复杂数据：JSON格式的结果数据
   - 支持状态信息：成功/失败状态、警告信息、日志信息

4. 函数配置 (FunctionOptions)：
   - Tags: 功能标签，用于分类
   - EnglishName: 英文名称，用于路由
   - ChineseName: 中文名称，用于显示
   - ApiDesc: 功能描述
   - Request/Response: 请求响应结构
   - RenderType: 必须是 response.RenderTypeForm
   - CreateTables: 需要创建的数据库表（可选）
   - OnPageLoad: 页面加载回调
   - OnApiCreated: API创建回调

5. 业务逻辑 (Function)：
   - 参数：ctx *runner.Context, req *请求结构, resp response.Response
   - 返回：error
   - 必须使用 resp.Form(&result).Build() 返回数据
   - 支持复杂的业务处理逻辑
   - 支持数据计算、外部API调用等
   - 业务逻辑中的判断必须与options中的完整值匹配

6. 路由注册：
   - 使用 runner.Post() 注册POST路由
   - 路径格式：/package名/函数英文名
   - 必须在init()函数中注册

7. 最佳实践：
   - default_value必须与options中的某个值完全匹配
   - 业务逻辑判断要使用完整的"value(label)"格式
   - 输入验证要充分
   - 错误处理要完善
   - 日志记录要详细
   - 返回信息要明确
   - 支持异步处理长时间任务
   - 合理使用事务处理数据库操作

8. 验证规则 (validate标签)：
   基于go-playground/validator库，常用规则：

   基础验证：
   - required: 必填字段
   - omitempty: 可选字段（为空时跳过验证）

   长度验证：
   - min=3: 最小长度3
   - max=100: 最大长度100
   - len=6: 固定长度6
   - min=3,max=20: 长度范围3-20

   数值验证：
   - min=0: 最小值0
   - max=100: 最大值100
   - gt=0: 大于0
   - gte=0: 大于等于0

   格式验证：
   - email: 邮箱格式
   - url: URL格式
   - numeric: 纯数字
   - alphanum: 字母数字
   - alpha: 纯字母

   枚举验证：
   - oneof=value1 value2: 必须是指定值之一
   - 注意：对于select组件，oneof的值必须与options中的value完全匹配

   正则表达式：
   - regexp=^pattern$: 自定义正则验证

   字段比较：
   - eqfield=FieldName: 必须等于指定字段（如确认密码）
   - gtfield=FieldName: 必须大于指定字段（如结束时间大于开始时间）

   条件验证：
   - required_if=Field value: 当指定字段等于指定值时必填
   - required_unless=Field value: 除非指定字段等于指定值，否则必填

   组合示例：
   - validate:"required,email": 必填邮箱
   - validate:"required,min=6,max=50": 必填，长度6-50
   - validate:"omitempty,url": 可选URL
   - validate:"required,oneof=active(启用) inactive(禁用)": 必填枚举值
   - validate:"required,eqfield=Password": 确认密码必须与密码字段相同
*/
