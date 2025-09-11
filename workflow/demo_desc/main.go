package main

import (
	"github.com/yunhanshu-net/pkg/workflow"
)

func main() {
	// 测试代码，包含带描述的步骤定义
	code := `
//desc: 用户验证步骤
step1 = validateUser(username: string "用户名", password: string "密码") -> (isValid: bool "验证结果", err: string "错误信息");

//desc: 发送欢迎邮件
step2 = sendWelcomeEmail(email: string "邮箱地址") -> (success: bool "发送结果");

func main() {
	//desc: 执行用户验证
	valid, err := step1("john", "password123")
	
	//desc: 如果验证成功，发送欢迎邮件
	if valid {
		//desc: 发送邮件
		success := step2("john@example.com")
	}
}
`

	parser := workflow.NewSimpleParser()
	result := parser.ParseWorkflow(code)
	result.Print()
}
