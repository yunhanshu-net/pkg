package usercall

const (
	// CallbackTypeOnPageLoad 页面事件
	CallbackTypeOnPageLoad     = "OnPageLoad"     // 页面加载时
	CallbackTypeOnCreateTables = "OnCreateTables" // 创建表

	// UserCallTypeOnApiCreated API 生命周期
	UserCallTypeOnApiCreated    = "OnApiCreated"    // API创建完成时
	CallbackTypeOnApiUpdated    = "OnApiUpdated"    // API更新时
	CallbackTypeBeforeApiDelete = "BeforeApiDelete" // API删除前
	CallbackTypeAfterApiDeleted = "AfterApiDeleted" // API删除后

	// CallbackTypeBeforeRunnerClose 运行器(Runner)生命周期
	CallbackTypeBeforeRunnerClose = "BeforeRunnerClose" // 运行器关闭前
	CallbackTypeAfterRunnerClose  = "AfterRunnerClose"  // 运行器关闭后

	// CallbackTypeOnVersionChange 版本控制
	CallbackTypeOnVersionChange = "OnVersionChange" // 版本变更时

	// CallbackTypeOnInputFuzzy 输入交互
	CallbackTypeOnInputFuzzy    = "OnInputFuzzy"    // 输入模糊匹配
	CallbackTypeOnInputValidate = "OnInputValidate" // 输入校验

	// CallbackTypeOnTableDeleteRows 表格操作
	CallbackTypeOnTableDeleteRows = "OnTableDeleteRows" // 删除表格行
	CallbackTypeOnTableAddRows    = "OnTableAddRows"    // 添加表格行
	CallbackTypeOnTableUpdateRows = "OnTableUpdateRows" // 更新表格行
	CallbackTypeOnTableSearch     = "OnTableSearch"     // 表格搜索

	// CallbackTypeOnDryRun 危险操作预览
	CallbackTypeOnDryRun = "OnDryRun" // DryRun 预览
)
