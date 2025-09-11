package workflow

import (
	"fmt"
)

func main() {
	parser := NewSimpleParser()

	// 测试复杂类型解析
	code := `var input = {
    "testData": "test"
}

//desc: 合并解析结果
step1 = code.service.merge_results(basicResult: map[string]interface{} "基础解析结果", advancedResult: map[string]interface{} "高级解析结果", isAdvanced: bool "是否高级解析") -> (result: map[string]interface{} "合并结果", err: error "是否失败")

func main() {
    fmt.Println("开始执行测试...")
    
    //desc: 执行合并
    合并结果, step1Err := step1(input["testData"], input["advancedData"], true)
    
    if step1Err != nil {
        step1.Printf("合并失败: %v", step1Err)
        return
    }
    
    fmt.Println("合并完成！")
}`

	result := parser.ParseWorkflow(code)
	if !result.Success {
		fmt.Printf("解析失败: %s\n", result.Error)
		return
	}

	fmt.Println("解析成功！")
	fmt.Printf("步骤数: %d\n", len(result.Steps))
	if len(result.Steps) > 0 {
		step := result.Steps[0]
		fmt.Printf("步骤名: %s\n", step.Name)
		fmt.Printf("函数名: %s\n", step.Function)
		fmt.Printf("是否静态工作流: %v\n", step.IsStatic)
		fmt.Printf("输入参数数: %d\n", len(step.InputParams))
		for i, param := range step.InputParams {
			fmt.Printf("  参数%d: %s %s \"%s\"\n", i+1, param.Name, param.Type, param.Desc)
		}
	}
}
