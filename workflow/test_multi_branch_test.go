package workflow

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	parser := NewSimpleParser()

	// 测试多分支if语句
	code := `var input = {
    "解析级别": 2,
    "代码内容": "test code"
}

//desc: 基础解析
step1 = code.service.basic_parsing(codeContent: string "代码内容") -> (result: string "基础结果", err: error "是否失败")
//desc: 高级解析
step2 = code.service.advanced_parsing(codeContent: string "代码内容", level: int "解析级别") -> (result: string "高级结果", err: error "是否失败")
//desc: 专家解析
step3 = code.service.expert_parsing(codeContent: string "代码内容", level: int "解析级别") -> (result: string "专家结果", err: error "是否失败")

func main() {
    sys.Println("开始解析工作流...")
    
    //desc: 根据解析级别选择解析方式
    if input["解析级别"].(int) == 1 {
        //desc: 执行基础解析
        基础结果, step1Err := step1(input["代码内容"])
        if step1Err != nil {
            step1.Printf("基础解析失败: %v", step1Err)
            return
        }
        step1.Printf("基础解析完成: %s", 基础结果)
    } else if input["解析级别"].(int) == 2 {
        //desc: 执行高级解析
        高级结果, step2Err := step2(input["代码内容"], input["解析级别"])
        if step2Err != nil {
            step2.Printf("高级解析失败: %v", step2Err)
            return
        }
        step2.Printf("高级解析完成: %s", 高级结果)
    } else if input["解析级别"].(int) == 3 {
        //desc: 执行专家解析
        专家结果, step3Err := step3(input["代码内容"], input["解析级别"])
        if step3Err != nil {
            step3.Printf("专家解析失败: %v", step3Err)
            return
        }
        step3.Printf("专家解析完成: %s", 专家结果)
    } else {
        //desc: 不支持的解析级别
        sys.Println("不支持的解析级别，使用默认解析")
        基础结果, step1Err := step1(input["代码内容"])
        if step1Err != nil {
            step1.Printf("默认解析失败: %v", step1Err)
            return
        }
        step1.Printf("默认解析完成: %s", 基础结果)
    }
    
    sys.Println("解析工作流完成！")
}`

	result := parser.ParseWorkflow(code)
	if !result.Success {
		fmt.Printf("解析失败: %s\n", result.Error)
		return
	}

	fmt.Println("解析成功！")
	fmt.Printf("步骤数: %d\n", len(result.Steps))

	if len(result.MainFunc.Statements) > 0 {
		fmt.Println("\n主函数语句:")
		for i, stmt := range result.MainFunc.Statements {
			fmt.Printf("语句 %d: [%s] %s\n", i+1, stmt.Type, stmt.Content)
			if stmt.Type == "if" && len(stmt.Children) > 0 {
				fmt.Printf("  分支数: %d\n", len(stmt.Children))
				for j, child := range stmt.Children {
					if child.Type == "branch" {
						fmt.Printf("    分支 %d: 条件=%s, 描述=%s, 语句数=%d\n",
							j+1, child.Condition, child.Desc, len(child.Children))
					}
				}
			}
		}
	}
}
func tes() {

}
