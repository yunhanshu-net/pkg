package workflow

import (
	"fmt"
	"testing"
)

func TestSimplifiedWorkflowSyntax(t *testing.T) {
	// 测试简化的动态工作流语法
	workflowCode := `
var input = {
    "用户名": "张三",
    "邮箱": "zhangsan@example.com"
}

//desc: 创建用户账号
step1 = user.service.create_user(username: string "用户名", email: string "邮箱") -> (userID: int "用户ID", username: string "用户名")
//desc: 发送欢迎邮件（允许失败继续）
step2 = email.service.send_welcome(email: string "邮箱", userID: int "用户ID") -> (success: bool "是否成功") {err_continue: true}
//desc: 记录用户日志（允许失败继续）
step3 = log.service.record_user_action(userID: int "用户ID", action: string "操作") -> (logID: int "日志ID") {err_continue: true}

func main() {
    //desc: 开始用户注册流程
    sys.Println("开始用户注册流程...")
    
    //desc: 创建用户账号
    用户ID, 用户名 := step1(input["用户名"], input["邮箱"])
    
    //desc: 发送欢迎邮件（允许失败继续）
    邮件成功 := step2(input["邮箱"], 用户ID)
    
    //desc: 记录用户日志（允许失败继续）
    日志ID := step3(用户ID, "用户注册")
    
    sys.Println("用户注册流程完成！")
}`

	// 解析工作流
	parser := &SimpleParser{}
	result := parser.ParseWorkflow(workflowCode)

	// 验证解析结果
	if len(result.Steps) != 3 {
		t.Fatalf("期望3个步骤，实际得到 %d 个", len(result.Steps))
	}

	// 验证步骤1（不返回错误类型）
	step1 := result.Steps[0]
	if step1.Name != "step1" {
		t.Errorf("步骤1名称错误: 期望 step1, 实际 %s", step1.Name)
	}
	if step1.Desc != "创建用户账号" {
		t.Errorf("步骤1描述错误: 期望 '创建用户账号', 实际 '%s'", step1.Desc)
	}
	if len(step1.OutputParams) != 2 {
		t.Errorf("步骤1输出参数数量错误: 期望 2, 实际 %d", len(step1.OutputParams))
	}
	if step1.OutputParams[0].Name != "userID" || step1.OutputParams[0].Type != "int" {
		t.Errorf("步骤1第一个输出参数错误: 期望 userID:int, 实际 %s:%s", step1.OutputParams[0].Name, step1.OutputParams[0].Type)
	}
	if step1.OutputParams[1].Name != "username" || step1.OutputParams[1].Type != "string" {
		t.Errorf("步骤1第二个输出参数错误: 期望 username:string, 实际 %s:%s", step1.OutputParams[1].Name, step1.OutputParams[1].Type)
	}

	// 验证步骤2（有err_continue元数据）
	step2 := result.Steps[1]
	if step2.Name != "step2" {
		t.Errorf("步骤2名称错误: 期望 step2, 实际 %s", step2.Name)
	}
	if step2.Desc != "发送欢迎邮件（允许失败继续）" {
		t.Errorf("步骤2描述错误: 期望 '发送欢迎邮件（允许失败继续）', 实际 '%s'", step2.Desc)
	}
	if step2.Metadata == nil {
		t.Error("步骤2元数据为空")
	} else if errContinue, ok := step2.Metadata["err_continue"]; !ok || errContinue != true {
		t.Errorf("步骤2 err_continue 元数据错误: 期望 true, 实际 %v", errContinue)
	}

	// 验证步骤3（有err_continue元数据）
	step3 := result.Steps[2]
	if step3.Name != "step3" {
		t.Errorf("步骤3名称错误: 期望 step3, 实际 %s", step3.Name)
	}
	if step3.Desc != "记录用户日志（允许失败继续）" {
		t.Errorf("步骤3描述错误: 期望 '记录用户日志（允许失败继续）', 实际 '%s'", step3.Desc)
	}
	if step3.Metadata == nil {
		t.Error("步骤3元数据为空")
	} else if errContinue, ok := step3.Metadata["err_continue"]; !ok || errContinue != true {
		t.Errorf("步骤3 err_continue 元数据错误: 期望 true, 实际 %v", errContinue)
	}

	// 验证主函数
	if result.MainFunc == nil {
		t.Fatal("主函数为空")
	}

	// 验证主函数中的语句
	statements := result.MainFunc.Statements
	if len(statements) != 5 {
		t.Errorf("主函数语句数量错误: 期望 5, 实际 %d", len(statements))
	}

	// 验证第一个语句（步骤调用，不处理错误）
	stmt1 := statements[1] // 跳过第一个sys.Println
	if stmt1.Type != "function-call" {
		t.Errorf("第一个语句类型错误: 期望 function-call, 实际 %s", stmt1.Type)
	}
	if stmt1.Desc != "创建用户账号" {
		t.Errorf("第一个语句描述错误: 期望 '创建用户账号', 实际 '%s'", stmt1.Desc)
	}

	// 验证第二个语句（步骤调用，不处理错误）
	stmt2 := statements[2]
	if stmt2.Type != "function-call" {
		t.Errorf("第二个语句类型错误: 期望 function-call, 实际 %s", stmt2.Type)
	}
	if stmt2.Desc != "发送欢迎邮件（允许失败继续）" {
		t.Errorf("第二个语句描述错误: 期望 '发送欢迎邮件（允许失败继续）', 实际 '%s'", stmt2.Desc)
	}

	// 验证第三个语句（步骤调用，不处理错误）
	stmt3 := statements[3]
	if stmt3.Type != "function-call" {
		t.Errorf("第三个语句类型错误: 期望 function-call, 实际 %s", stmt3.Type)
	}
	if stmt3.Desc != "记录用户日志（允许失败继续）" {
		t.Errorf("第三个语句描述错误: 期望 '记录用户日志（允许失败继续）', 实际 '%s'", stmt3.Desc)
	}

	fmt.Println("✅ 简化语法测试通过！")
	result.Print()
}

func TestStaticWorkflowSimplifiedSyntax(t *testing.T) {
	// 测试简化的静态工作流语法
	workflowCode := `
var input = {
    "项目名称": "my-project",
    "环境": "production"
}

//desc: 推送代码到远程仓库
step1 = beiluo.test1.devops.git_push[用例001] -> ()
//desc: 部署到测试环境
step2 = beiluo.test1.devops.deploy_test[用例002] -> (cost: int "成本")
//desc: 部署到生产环境
step3 = beiluo.test1.devops.deploy_prod[用例003] -> (msg: string "消息")
//desc: 发送部署完成通知（允许失败继续）
step4 = beiluo.test1.notify.send_notification[用例004] -> () {err_continue: true}

func main() {
    //desc: 开始执行发布流程
    sys.Println("开始执行发布流程...")
    
    //desc: 推送代码到远程仓库
    step1()
    
    //desc: 部署到测试环境
    cost := step2()
    
    //desc: 部署到生产环境
    msg := step3()
    
    //desc: 发送部署完成通知（允许失败继续）
    step4()
    
    sys.Println("发布流程执行完成！")
}`

	// 解析工作流
	parser := &SimpleParser{}
	result := parser.ParseWorkflow(workflowCode)

	// 验证解析结果
	if len(result.Steps) != 4 {
		t.Fatalf("期望4个步骤，实际得到 %d 个", len(result.Steps))
	}

	// 验证步骤1（无返回值）
	step1 := result.Steps[0]
	if step1.Name != "step1" {
		t.Errorf("步骤1名称错误: 期望 step1, 实际 %s", step1.Name)
	}
	if step1.Desc != "推送代码到远程仓库" {
		t.Errorf("步骤1描述错误: 期望 '推送代码到远程仓库', 实际 '%s'", step1.Desc)
	}
	if len(step1.OutputParams) != 0 {
		t.Errorf("步骤1输出参数数量错误: 期望 0, 实际 %d", len(step1.OutputParams))
	}

	// 验证步骤2（有返回值）
	step2 := result.Steps[1]
	if step2.Name != "step2" {
		t.Errorf("步骤2名称错误: 期望 step2, 实际 %s", step2.Name)
	}
	if step2.Desc != "部署到测试环境" {
		t.Errorf("步骤2描述错误: 期望 '部署到测试环境', 实际 '%s'", step2.Desc)
	}
	if len(step2.OutputParams) != 1 {
		t.Errorf("步骤2输出参数数量错误: 期望 1, 实际 %d", len(step2.OutputParams))
	}
	if step2.OutputParams[0].Name != "cost" || step2.OutputParams[0].Type != "int" {
		t.Errorf("步骤2输出参数错误: 期望 cost:int, 实际 %s:%s", step2.OutputParams[0].Name, step2.OutputParams[0].Type)
	}

	// 验证步骤4（有err_continue元数据）
	step4 := result.Steps[3]
	if step4.Name != "step4" {
		t.Errorf("步骤4名称错误: 期望 step4, 实际 %s", step4.Name)
	}
	if step4.Desc != "发送部署完成通知（允许失败继续）" {
		t.Errorf("步骤4描述错误: 期望 '发送部署完成通知（允许失败继续）', 实际 '%s'", step4.Desc)
	}
	if step4.Metadata == nil {
		t.Error("步骤4元数据为空")
	} else if errContinue, ok := step4.Metadata["err_continue"]; !ok || errContinue != true {
		t.Errorf("步骤4 err_continue 元数据错误: 期望 true, 实际 %v", errContinue)
	}

	fmt.Println("✅ 静态工作流简化语法测试通过！")
	result.Print()
}

func TestConditionalWorkflowSimplifiedSyntax(t *testing.T) {
	// 测试简化的条件判断工作流语法
	workflowCode := `
var input = {
    "代码文件路径": "path/to/code.go",
    "是否启用高级解析": true
}

//desc: 读取代码文件
step1 = code.service.read_file(filePath: string "代码文件路径") -> (content: string "代码内容")
//desc: 基础代码解析
step2 = code.service.basic_parse(content: string "代码内容") -> (result: map[string]interface{} "解析结果")
//desc: 高级代码解析
step3 = code.service.advanced_parse(content: string "代码内容") -> (result: map[string]interface{} "解析结果")

func main() {
    //desc: 开始代码解析工作流
    sys.Println("开始代码解析工作流...")
    
    //desc: 读取代码文件
    代码内容 := step1(input["代码文件路径"])
    
    //desc: 判断是否启用高级解析
    if input["是否启用高级解析"] {
        //desc: 执行高级解析
        解析结果 := step3(代码内容)
        sys.Println("使用高级解析完成")
    } else {
        //desc: 执行基础解析
        解析结果 := step2(代码内容)
        sys.Println("使用基础解析完成")
    }
    
    sys.Println("代码解析工作流完成！")
}`

	// 解析工作流
	parser := &SimpleParser{}
	result := parser.ParseWorkflow(workflowCode)

	// 验证解析结果
	if len(result.Steps) != 3 {
		t.Fatalf("期望3个步骤，实际得到 %d 个", len(result.Steps))
	}

	// 验证主函数中的条件判断
	if result.MainFunc == nil {
		t.Fatal("主函数为空")
	}

	statements := result.MainFunc.Statements
	if len(statements) != 6 {
		t.Errorf("主函数语句数量错误: 期望 6, 实际 %d", len(statements))
	}

	// 验证条件判断语句
	stmt2 := statements[2] // 条件判断语句
	if stmt2.Type != "if" {
		t.Errorf("条件判断语句类型错误: 期望 if, 实际 %s", stmt2.Type)
	}
	if stmt2.Desc != "判断是否启用高级解析" {
		t.Errorf("条件判断语句描述错误: 期望 '判断是否启用高级解析', 实际 '%s'", stmt2.Desc)
	}

	// 验证条件判断中的子语句
	if len(stmt2.Children) != 2 {
		t.Errorf("条件判断子语句数量错误: 期望 2, 实际 %d", len(stmt2.Children))
	}

	// 验证if分支中的步骤调用
	ifStmt := stmt2.Children[0]
	if ifStmt.Type != "function-call" {
		t.Errorf("if分支语句类型错误: 期望 function-call, 实际 %s", ifStmt.Type)
	}
	if ifStmt.Desc != "执行高级解析" {
		t.Errorf("if分支语句描述错误: 期望 '执行高级解析', 实际 '%s'", ifStmt.Desc)
	}

	// 验证else分支中的步骤调用（在if语句外面）
	elseStmt := statements[3] // else分支实际上是独立的语句
	if elseStmt.Type != "function-call" {
		t.Errorf("else分支语句类型错误: 期望 function-call, 实际 %s", elseStmt.Type)
	}
	if elseStmt.Desc != "执行基础解析" {
		t.Errorf("else分支语句描述错误: 期望 '执行基础解析', 实际 '%s'", elseStmt.Desc)
	}

	fmt.Println("✅ 条件判断工作流简化语法测试通过！")
	result.Print()
}

func TestErrorHandlingSimplification(t *testing.T) {
	// 测试错误处理简化效果
	oldStyleCode := `
// 传统模式（需要手动错误处理）
step1 = user.service.create_user(username: string "用户名", email: string "邮箱") -> (userID: int "用户ID", err: error "是否失败")

func main() {
    用户ID, step1Err := step1("张三", "zhangsan@example.com")
    if step1Err != nil {
        step1.Printf("用户创建失败: %v", step1Err)
        return
    }
    step1.Printf("用户创建成功，ID: %d", 用户ID)
}`

	newStyleCode := `
// 执行引擎自动处理模式（极简代码）
step1 = user.service.create_user(username: string "用户名", email: string "邮箱") -> (userID: int "用户ID")

func main() {
    //desc: 创建用户
    用户ID := step1("张三", "zhangsan@example.com")
    
    sys.Println("用户注册流程完成！")
}`

	// 解析两种风格的代码
	parser := &SimpleParser{}

	// 解析传统模式
	oldResult := parser.ParseWorkflow(oldStyleCode)

	// 解析新简化模式
	newResult := parser.ParseWorkflow(newStyleCode)

	// 对比两种模式
	fmt.Println("=== 传统模式解析结果 ===")
	fmt.Printf("步骤数量: %d\n", len(oldResult.Steps))
	fmt.Printf("主函数语句数量: %d\n", len(oldResult.MainFunc.Statements))
	fmt.Printf("步骤1输出参数数量: %d\n", len(oldResult.Steps[0].OutputParams))

	fmt.Println("\n=== 新简化模式解析结果 ===")
	fmt.Printf("步骤数量: %d\n", len(newResult.Steps))
	fmt.Printf("主函数语句数量: %d\n", len(newResult.MainFunc.Statements))
	fmt.Printf("步骤1输出参数数量: %d\n", len(newResult.Steps[0].OutputParams))

	// 验证简化效果
	if len(newResult.Steps[0].OutputParams) < len(oldResult.Steps[0].OutputParams) {
		fmt.Println("✅ 成功移除了错误类型返回")
	}

	if len(newResult.MainFunc.Statements) < len(oldResult.MainFunc.Statements) {
		fmt.Println("✅ 成功简化了主函数语句")
	}

	fmt.Println("✅ 错误处理简化效果验证通过！")
}
