// Package table 这个文件是用来展示整个table函数的全部功能和能力，就像element-plus这种ui组件的组件一览表一样，这里没有太多逻辑，
//只是展示一下怎么用，函数的签名和组件的参数，标签，回调等等，目前项目处于开发中，只展示已经实现的功能，不展示未实现的功能

package table

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/response" //这个是渲染返回的数据结构的包，必须引用这个包
	"github.com/yunhanshu-net/function-go/pkg/dto/usercall" //这是回调的请求和响应model，用到回调的请引用这个包
	"github.com/yunhanshu-net/function-go/runner"           //这个是function-go的核心包，需要引用这个才能启动
	"github.com/yunhanshu-net/pkg/logger"                   //这个是日志包，需要引用这个才能使用日志
	"github.com/yunhanshu-net/pkg/query"                    //这个是查询的包，table函数的请求可以用，query.PageInfoReq 这个通用结构，里面有分页相关信息
)

// ===== 数据模型 =====

// ExampleModel 展示table函数支持的所有字段类型和标签配置
type ExampleModel struct {
	// 基础字段
	ID          int    `gorm:"primaryKey;autoIncrement" runner:"code:id;name:ID;example:1001"`                                                                                // 主键字段：GORM自动递增主键，无需验证，示例值展示ID格式
	Title       string `json:"title" form:"title" runner:"code:title;name:标题;type:string;placeholder:请输入标题;example:产品发布会" validate:"required,max=200"`                        // 标题验证：必填 + 最大长度200，示例值帮助理解字段用途
	Description string `json:"description" form:"description" runner:"code:description;name:描述;type:string;placeholder:请输入描述信息;example:这是一个关于新产品发布的详细描述" validate:"max=1000"` // 描述验证：可选字段 + 最大长度1000，示例值展示描述格式

	// 常用验证规则示例
	Email    string `json:"email" form:"email" runner:"code:email;name:邮箱;type:string;placeholder:请输入邮箱地址;example:user@example.com" validate:"required,email"`                      // 邮箱格式验证：必填 + 邮箱格式校验，示例值展示邮箱格式
	Phone    string `json:"phone" form:"phone" runner:"code:phone;name:手机号;type:string;placeholder:请输入11位手机号;example:13812345678" validate:"required,len=11,numeric"`               // 手机号验证：必填 + 固定11位长度 + 纯数字，示例值展示手机号格式
	Username string `json:"username" form:"username" runner:"code:username;name:用户名;type:string;placeholder:请输入用户名;example:user123" validate:"required,min=3,max=20,alphanum"`      // 用户名验证：必填 + 长度3-20 + 字母数字组合，示例值展示用户名格式
	Password string `json:"password" form:"password" runner:"code:password;name:密码;type:string;placeholder:请输入密码;hidden:list;example:password123" validate:"required,min=6,max=50"` // 密码验证：必填 + 长度6-50，列表中隐藏显示，示例值展示密码格式
	Website  string `json:"website" form:"website" runner:"code:website;name:网站;type:string;placeholder:请输入网站地址;example:https://www.example.com" validate:"omitempty,url"`          // 网站验证：可选字段 + URL格式校验，示例值展示URL格式

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
	Status   string `json:"status" form:"status" runner:"code:status;name:状态;widget:select;default_value:active(启用);options:active(启用),inactive(禁用);type:string;example:active(启用)" validate:"required,oneof=active(启用) inactive(禁用)"`                         // 状态验证：必填 + 二选一枚举（启用/禁用），示例值展示启用状态
	Category string `json:"category" form:"category" runner:"code:category;name:分类;widget:select;default_value:type1(类型1);options:type1(类型1),type2(类型2),type3(类型3);type:string;example:type1(类型1)" validate:"required,oneof=type1(类型1) type2(类型2) type3(类型3)"` // 分类验证：必填 + 三选一枚举，default_value必须在options中，示例值与默认值一致
	Priority string `json:"priority" form:"priority" runner:"code:priority;name:优先级;widget:select;default_value:medium(中);options:low(低),medium(中),high(高);type:string;example:medium(中)" validate:"required,oneof=low(低) medium(中) high(高)"`                  // 优先级验证：必填 + 三级优先级枚举，示例值展示中等优先级

	// 数值字段
	Amount float64 `json:"amount" form:"amount" runner:"code:amount;name:金额;type:number;placeholder:请输入金额;example:1299.99" validate:"min=0"`                 // 金额验证：最小值0（非负数），可选字段，示例值展示金额格式
	Count  int     `json:"count" form:"count" runner:"code:count;name:数量;type:number;default_value:0;placeholder:请输入数量;example:5" validate:"min=0,max=9999"` // 数量验证：数值范围0-9999，有默认值，示例值展示合理数量

	// 布尔字段
	IsActive string `json:"is_active" form:"is_active" runner:"code:is_active;name:是否启用;widget:select;default_value:true(启用);options:true(启用),false(禁用);example:true(启用)" validate:"required,oneof=true(启用) false(禁用)"` // 布尔选择器：用select实现布尔选择，必填，示例值展示启用状态
	IsPublic bool   `json:"is_public" form:"is_public" runner:"code:is_public;name:是否公开;widget:switch;true_label:公开;false_label:私有;type:boolean;default_value:false;example:false"`                                     // 原生布尔类型：直接使用bool类型，有默认值，示例值与默认值一致,前端会展示label和switch组件

	// 字段显示控制示例
	CreatedBy string `json:"created_by" gorm:"column:created_by" runner:"code:created_by;name:创建人;show:list;example:admin"`                         // 显示控制：只在列表页面显示，创建和更新时隐藏，示例值展示创建人
	UpdatedAt string `json:"updated_at" gorm:"autoUpdateTime" runner:"code:updated_at;name:更新时间;type:string;show:list;example:2024-01-15 14:30:25"` // 时间字段：GORM自动更新时间，只在列表显示，示例值展示时间格式
	CreatedAt string `json:"created_at" gorm:"autoCreateTime" runner:"code:created_at;name:创建时间;type:string;show:list;example:2024-01-15 14:25:10"` // 时间字段：GORM自动创建时间，只在列表显示，示例值展示时间格式

	// 隐藏字段示例（只参与后端逻辑）
	UserId   int    `json:"user_id" gorm:"column:user_id" runner:"code:user_id;name:用户ID;hidden:list,create,update;example:1001"`             // 完全隐藏：在所有场景都隐藏，只用于后端逻辑，示例值展示用户ID
	TenantId string `json:"tenant_id" gorm:"column:tenant_id" runner:"code:tenant_id;name:租户ID;hidden:list,create,update;example:tenant_001"` // 多租户字段：完全隐藏，用于数据隔离，示例值展示租户ID格式
}

func (em *ExampleModel) TableName() string {
	return "example_table"
}

// ===== 请求参数结构体 =====

// ExampleListReq 展示table函数支持的所有查询参数类型
type ExampleListReq struct {
	query.PageInfoReq // 必须嵌入分页参数

	// 精确匹配查询
	Category string `json:"category" form:"category" runner:"code:category;name:分类;widget:select;options:type1(类型1),type2(类型2),type3(类型3);type:string;example:type1(类型1)"`
	Status   string `json:"status" form:"status" runner:"code:status;name:状态;widget:select;options:active(启用),inactive(禁用);type:string;example:active(启用)"`
	Priority string `json:"priority" form:"priority" runner:"code:priority;name:优先级;widget:select;options:low(低),medium(中),high(高);type:string;example:medium(中)"`

	// 模糊搜索
	Keyword string `json:"keyword" form:"keyword" runner:"code:keyword;name:关键词;type:string;placeholder:搜索标题或描述;example:产品发布"`

	// 范围查询
	MinAmount float64 `json:"min_amount" form:"min_amount" runner:"code:min_amount;name:最小金额;type:number;placeholder:最小金额;example:100.00"`
	MaxAmount float64 `json:"max_amount" form:"max_amount" runner:"code:max_amount;name:最大金额;type:number;placeholder:最大金额;example:5000.00"`

	// 布尔查询
	IsActive string `json:"is_active" form:"is_active" runner:"code:is_active;name:是否启用;widget:select;options:true(启用),false(禁用);default:true(启用);example:true(启用)"`
}

// ===== 函数配置 =====

var ExampleTableConfig = &runner.FunctionInfo{
	// 基础信息
	Tags:        []string{"示例", "表格管理", "数据展示"},
	EnglishName: "example_table",
	ChineseName: "示例表格管理",
	ApiDesc:     "展示table函数的完整功能，包括数据展示、筛选、分页、增删改查等操作",

	// 请求响应结构
	Request:  &ExampleListReq{},
	Response: query.PaginatedTable[[]ExampleModel]{}, // table函数必须用query.PaginatedTable[]来包裹

	// 渲染类型
	RenderType: response.RenderTypeTable, // 标识返回的数据是table格式

	// 数据库表管理
	CreateTables:  []interface{}{&ExampleModel{}}, // 注册创建表的变量
	AutoCrudTable: &ExampleModel{},                // 自动创建增删改查接口，AutoCrudTable如果注册来了表的话，OperateTables会自动生成表的crud操作，此时OperateTables可以为nil,等价下面的注射
	//OperateTables: map[interface{}][]runner.OperateTableType{ // AutoCrudTable: &ExampleModel{}
	//	&ExampleModel{}: {runner.OperateTableTypeGet, runner.OperateTableTypeUpdate, runner.OperateTableTypeDelete, runner.OperateTableTypeAdd},
	//},

	// 页面加载回调
	OnPageLoad: func(ctx *runner.Context, resp response.Response) (initData *usercall.OnPageLoadResp, err error) {
		// 页面加载时执行的函数，可以做初始化操作
		db := ctx.MustGetOrInitDB()
		var results []ExampleModel
		filter := &query.PageInfoReq{PageSize: 10}
		isActive := "true(启用)"

		// 默认查询启用状态的数据，按ID降序
		err = resp.Table(&results).AutoPaginated(
			db.Where("is_active = ?", isActive).Order("id DESC"),
			&ExampleModel{},
			filter,
		).Build()

		if err != nil {
			logger.Errorf(ctx, "示例表格数据加载失败: %v", err)
			return
		}

		// 返回初始化数据，设置默认筛选条件
		initData = &usercall.OnPageLoadResp{
			Request: &ExampleListReq{PageInfoReq: *filter, IsActive: isActive},
			AutoRun: false, // 是否在页面加载时自动运行一次函数
		}
		return
	},

	// API创建回调
	OnApiCreated: func(ctx *runner.Context, req *usercall.OnApiCreatedReq) error {
		logger.Infof(ctx, "示例表格管理API创建成功: %+v", req)
		return nil
	},
}

// ===== 路由注册 =====

func init() {
	// 路由注册：/package名/函数英文名
	runner.Get("/table/example_table", ExampleTable, ExampleTableConfig)
}

// ===== 主要业务逻辑 =====

func ExampleTable(ctx *runner.Context, req *ExampleListReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB()
	var results []ExampleModel

	// 构建查询条件
	query := db.Model(&ExampleModel{})

	// 精确匹配查询
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Priority != "" {
		query = query.Where("priority = ?", req.Priority)
	}
	if req.IsActive != "" {
		query = query.Where("is_active = ?", req.IsActive)
	}

	// 模糊搜索
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", keyword, keyword)
	}

	// 范围查询
	if req.MinAmount > 0 {
		query = query.Where("amount >= ?", req.MinAmount)
	}
	if req.MaxAmount > 0 {
		query = query.Where("amount <= ?", req.MaxAmount)
	}

	// 默认排序：按ID降序
	query = query.Order("id DESC")

	// table函数必须用resp.Table来返回，AutoPaginated会自动分页处理
	return resp.Table(&results).AutoPaginated(query, &ExampleModel{}, &req.PageInfoReq).Build()
}

// ===== 功能说明 =====

/*
Table函数完整功能说明：

1. 数据模型 (Model)：
   - 支持所有Go基础类型：string, int, float64, bool
   - 支持GORM标签：主键、自动递增、列名等
   - 支持runner标签：字段配置、UI组件、验证规则

2. Runner标签说明（已实现功能）：
   - code: 字段代码，用于前端识别
   - name: 字段显示名称
   - type: 字段类型 (string, number, boolean)
   - widget: UI组件 (select, input)
   - placeholder: 输入提示
   - default_value: 默认值（必须与options中的某个值完全匹配）
   - options: 选择器选项 (格式：value(label),value2(label2))
   - validate: 验证规则 (required, max, min, oneof等)

3. 字段显示控制（show/hidden）：
   - 默认行为：不填show和hidden = 全部场景显示
   - 正向控制：show:list,create,update = 只在指定场景显示
   - 反向控制：hidden:list,create,update = 在指定场景隐藏
   - 支持场景：list(列表), create(创建), update(更新)
   - 互斥规则：show和hidden只能选其一
   - 常用示例：
     * 创建时间/更新时间：show:list（只在列表显示）
     * 用户ID/租户ID：hidden:list,create,update（完全隐藏）
     * 密码字段：hidden:list（列表时隐藏）

4. 请求参数 (Request)：
   - 必须嵌入 query.PageInfoReq 用于分页
   - 支持各种查询类型：精确匹配、模糊搜索、范围查询
   - 参数也需要runner标签配置UI组件

5. 函数配置 (FunctionInfo)：
   - Tags: 功能标签，用于分类
   - EnglishName: 英文名称，用于路由
   - ChineseName: 中文名称，用于显示
   - ApiDesc: 功能描述
   - Request/Response: 请求响应结构
   - RenderType: 必须是 response.RenderTypeTable
   - CreateTables: 需要创建的数据库表
   - AutoCrudTable: 自动生成增删改查接口
   - OnPageLoad: 页面加载回调
   - OnApiCreated: API创建回调

6. 业务逻辑 (Function)：
   - 参数：ctx *runner.Context, req *请求结构, resp response.Response
   - 返回：error
   - 必须使用 resp.Table(&results).AutoPaginated() 返回数据
   - 支持复杂的查询条件构建
   - 支持分页、排序、筛选

7. 路由注册：
   - 使用 runner.Get() 注册GET路由
   - 路径格式：/package名/函数英文名
   - 必须在init()函数中注册

8. 最佳实践：
   - default_value必须与options中的某个值完全匹配
   - 业务逻辑判断要使用完整的"value(label)"格式
   - 合理使用show/hidden控制字段显示
   - 创建时间/更新时间建议使用show:list
   - 敏感字段建议使用hidden完全隐藏
   - 输入验证要充分
   - 错误处理要完善
   - 日志记录要详细
   - 合理使用索引优化查询性能
   - 支持事务处理数据一致性

9. 验证规则 (validate标签)：
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
   - eqfield=FieldName: 必须等于指定字段
   - gtfield=FieldName: 必须大于指定字段

   条件验证：
   - required_if=Field value: 当指定字段等于指定值时必填
   - required_unless=Field value: 除非指定字段等于指定值，否则必填

   组合示例：
   - validate:"required,email": 必填邮箱
   - validate:"required,min=6,max=50": 必填，长度6-50
   - validate:"omitempty,url": 可选URL
   - validate:"required,oneof=active(启用) inactive(禁用)": 必填枚举值
*/
