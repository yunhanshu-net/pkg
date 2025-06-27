package form

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/pkg/dto/usercall"
	"github.com/yunhanshu-net/function-go/runner"
	"github.com/yunhanshu-net/pkg/logger"
)

// EnhancedFormReq 展示新标签系统的使用示例
type EnhancedFormReq struct {
	// 基础字符串字段 - 新标签系统示例
	Title string `json:"title" form:"title" 
			  runner:"code:title;name:标题" 
			  widget:"type:input;placeholder:请输入标题" 
			  data:"type:string;example:产品发布会;default_value:" 
			  validate:"required,max=200"`

	// 多行文本字段
	Description string `json:"description" form:"description"
					runner:"code:description;name:描述"
					widget:"type:input;mode:text_area;placeholder:请输入详细描述"
					data:"type:string;example:这是一个关于新产品发布的详细描述"
					validate:"max=1000"`

	// 选择器字段 - 简化options格式
	Category string `json:"category" form:"category"
				runner:"code:category;name:分类"
				widget:"type:select;placeholder:请选择分类;options:产品,营销,技术,运营"
				data:"type:string;default_value:产品"
				validate:"required"`

	// 多选字段 - 支持[]string类型
	Tags []string `json:"tags" form:"tags"
			 runner:"code:tags;name:标签"
			 widget:"type:multiselect;placeholder:选择相关标签;options:重要,紧急,公开,内部"
			 data:"type:[]string;example:重要,公开"
			 validate:"min=1"`

	// 数字字段
	Priority int `json:"priority" form:"priority"
			 runner:"code:priority;name:优先级"
			 widget:"type:slider;min:1;max:10"
			 data:"type:number;default_value:5;example:8"
			 validate:"min=1,max=10"`

	// 布尔字段
	IsPublic bool `json:"is_public" form:"is_public"
			  runner:"code:is_public;name:是否公开"
			  widget:"type:switch;true_label:公开;false_label:私有"
			  data:"type:boolean;default_value:false"`

	// 带搜索功能的用户选择 - 字段级别回调
	AssignedUser string `json:"assigned_user" form:"assigned_user"
				  runner:"code:assigned_user;name:负责人"
				  widget:"type:select;placeholder:搜索并选择负责人"
				  data:"type:string;source:api://users"
				  callback:"OnInputFuzzy(delay:300,min:2)"
				  validate:"required"`

	// 只读字段 - 新权限系统
	CreatedAt string `json:"created_at" form:"created_at"
			   runner:"code:created_at;name:创建时间"
			   widget:"type:input"
			   data:"type:string"
			   permission:"read"`

	// 仅创建时可编辑的字段
	ProjectType string `json:"project_type" form:"project_type"
				 runner:"code:project_type;name:项目类型"
				 widget:"type:select;options:内部项目,外部项目,合作项目"
				 data:"type:string;default_value:内部项目"
				 permission:"read,create"
				 validate:"required"`

	// 带多个回调的复杂字段
	Budget float64 `json:"budget" form:"budget"
			runner:"code:budget;name:预算"
			widget:"type:input;placeholder:请输入预算金额"
			data:"type:number;example:10000.50"
			callback:"OnInputChange();OnBlur(validate:true,format:currency)"
			validate:"min=0"`

	// 文件上传字段
	Attachments []string `json:"attachments" form:"attachments"
				 runner:"code:attachments;name:附件"
				 widget:"type:file_upload;accept:.pdf,.doc,.docx,.jpg,.png"
				 data:"type:[]string"
				 callback:"OnFileChange(max_size:10MB,max_count:5)"`
}

// EnhancedFormResp 响应结构体
type EnhancedFormResp struct {
	Message     string  `json:"message" runner:"code:message;name:处理结果;type:string" data:"example:表单处理成功"`
	Status      string  `json:"status" runner:"code:status;name:状态;type:string" data:"example:success"`
	ProcessTime string  `json:"process_time" runner:"code:process_time;name:处理时间;type:string" data:"example:2024-01-15 14:30:25"`
	TotalAmount float64 `json:"total_amount" runner:"code:total_amount;name:总金额;type:number" data:"example:12580.50"`
	IsSuccess   bool    `json:"is_success" runner:"code:is_success;name:是否成功;type:boolean" widget:"true_label:成功;false_label:失败" data:"example:true"`
}

// 函数配置
var EnhancedFormConfig = &runner.FunctionOptions{
	Tags:        []string{"示例", "新标签系统", "表单处理"},
	EnglishName: "enhanced_form",
	ChineseName: "增强表单示例",
	ApiDesc:     "展示新标签系统的完整功能，包括分离的标签配置、权限控制、回调函数等",

	Request:  &EnhancedFormReq{},
	Response: &EnhancedFormResp{},

	RenderType: response.RenderTypeForm,

	OnPageLoad: func(ctx *runner.Context, resp response.Response) (initData *usercall.OnPageLoadResp, err error) {
		// 返回初始化数据
		initData = &usercall.OnPageLoadResp{
			Request: EnhancedFormReq{
				Category: "类型1",
			},
			AutoRun: false,
		}
		return
	},

	OnApiCreated: func(ctx *runner.Context, req *usercall.OnApiCreatedReq) error {
		logger.Infof(ctx, "增强表单示例API创建成功: %+v", req)
		return nil
	},
}

func init() {
	runner.Post("/form/enhanced_form", EnhancedForm, EnhancedFormConfig)
}

func EnhancedForm(ctx *runner.Context, req *EnhancedFormReq, resp response.Response) error {
	// 简单的业务逻辑处理
	result := &EnhancedFormResp{
		Message:     "表单处理成功",
		Status:      "success",
		ProcessTime: "2024-01-15 14:30:25",
		IsSuccess:   true,
	}

	return resp.Form(result).Build()
}

// ReadOnlyFormReq 只读表单示例 - 展示权限控制
type ReadOnlyFormReq struct {
	// 系统生成的只读字段
	ID string `json:"id" form:"id"
		 runner:"code:id;name:任务ID"
		 widget:"type:input"
		 data:"type:string"
		 permission:"read"`

	// 普通可编辑字段
	TaskName string `json:"task_name" form:"task_name"
			  runner:"code:task_name;name:任务名称"
			  widget:"type:input;placeholder:请输入任务名称"
			  data:"type:string"
			  permission:"read,update,create"
			  validate:"required,max=100"`

	// 仅更新时可编辑
	Status string `json:"status" form:"status"
			runner:"code:status;name:状态"
			widget:"type:select;options:待开始,进行中,已完成,已取消"
			data:"type:string;default_value:待开始"
			permission:"read,update"`
}

// SearchFormReq 搜索表单示例 - 展示动态数据源
type SearchFormReq struct {
	// 用户搜索字段 - 动态数据源
	UserSearch string `json:"user_search" form:"user_search"
				runner:"code:user_search;name:用户搜索"
				widget:"type:select;placeholder:搜索用户名、邮箱或部门"
				data:"type:string;source:api://users/search"
				callback:"OnInputFuzzy(delay:300,min:2)"`

	// 部门搜索 - 级联选择
	Department string `json:"department" form:"department"
				runner:"code:department;name:部门"
				widget:"type:select;placeholder:选择部门"
				data:"type:string;source:api://departments"
				callback:"OnChange(cascade:user_search)"`

	// 日期范围 - 复合字段
	DateRange []string `json:"date_range" form:"date_range"
			   runner:"code:date_range;name:日期范围"
			   widget:"type:datetime;mode:range"
			   data:"type:[]string;format:YYYY-MM-DD"
			   validate:"required"`
}
