package workflow

type Param struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
type ExecutorIn struct {
	RealInput map[string]interface{} `json:"real_input"`
}
type ExecutorOut struct {

	//面试时间, 面试官名称, step2Err := step2(用户名)
	WantOutput map[string]interface{} `json:"want_output"`
	//这里的Output会输出
}

type OnFunctionCall func(in *ExecutorIn) (*ExecutorOut, error)

type Executor struct {
	OnFunctionCall
}

type ExecutorResp struct {
}

func (e *Executor) Start(code string) error {
	parser := NewSimpleParser()
	workflow := parser.ParseWorkflow(code)

}
