package tagx

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAllWidgetParsersAvailable(t *testing.T) {
	// 测试所有组件解析器是否都能正常创建
	expectedWidgets := []string{
		"file_upload", "color", "datetime", "slider", "multiselect",
		"select", "input", "switch", "radio", "checkbox",
		"file_display", "table",
	}

	for _, widgetType := range expectedWidgets {
		parser, exists := widgetParsers[widgetType]
		if !exists {
			t.Errorf("Widget parser for %s not found", widgetType)
			continue
		}

		// 测试每个解析器的基本功能
		keys := parser.GetSupportedKeys()
		if len(keys) == 0 {
			t.Errorf("Widget parser for %s has no supported keys", widgetType)
		}

		fmt.Printf("✓ %s: %v\n", widgetType, keys)
	}

	fmt.Printf("\n总共支持 %d 种组件类型\n", len(widgetParsers))
}

func TestWidgetConfigParsing(t *testing.T) {
	parser := NewMultiTagParser()

	// 测试基本的widget标签解析
	testCases := []struct {
		tag      string
		expected string
	}{
		{"type:input;placeholder:请输入", "input"},
		{"type:select;options:a,b,c", "select"},
		{"type:switch;true_label:是;false_label:否", "switch"},
		{"type:slider;min:0;max:100", "slider"},
		{"type:multiselect;options:x,y,z", "multiselect"},
		{"type:radio;options:male,female", "radio"},
		{"type:checkbox;options:a,b,c;show_select_all:true", "checkbox"},
		{"type:file_upload;accept:.jpg,.png;max_size:2MB", "file_upload"},
		{"type:file_display;display_mode:card;preview:true", "file_display"},
		{"type:color;format:hex;show_alpha:true", "color"},
		{"type:datetime;format:date;placeholder:选择日期", "datetime"},
		{"type:table;pagination:true;sortable:true", "table"},
	}

	for _, tc := range testCases {
		config := parser.parseWidgetTag(tc.tag)
		if config.Type != tc.expected {
			t.Errorf("Expected widget type %s, got %s for tag: %s", tc.expected, config.Type, tc.tag)
		}
	}
}

func TestPermissionTagParsing(t *testing.T) {
	parser := NewMultiTagParser()

	// 测试permission标签解析
	testCases := []struct {
		tag      string
		expected *PermissionConfig
	}{
		{
			tag:      "",
			expected: &PermissionConfig{Read: true, Update: true, Create: true}, // 空标签默认全部权限
		},
		{
			tag:      "read",
			expected: &PermissionConfig{Read: true, Update: false, Create: false},
		},
		{
			tag:      "read,update",
			expected: &PermissionConfig{Read: true, Update: true, Create: false},
		},
		{
			tag:      "read,update,create",
			expected: &PermissionConfig{Read: true, Update: true, Create: true},
		},
		{
			tag:      "create",
			expected: &PermissionConfig{Read: false, Update: false, Create: true},
		},
		{
			tag:      "update,create",
			expected: &PermissionConfig{Read: false, Update: true, Create: true},
		},
	}

	for _, tc := range testCases {
		result := parser.parsePermissionTag(tc.tag)
		if result.Read != tc.expected.Read || result.Update != tc.expected.Update || result.Create != tc.expected.Create {
			t.Errorf("Permission parsing failed for tag '%s'. Expected: %+v, Got: %+v",
				tc.tag, tc.expected, result)
		}
	}
}

func TestCompleteStructParsing(t *testing.T) {
	parser := NewMultiTagParser()

	// 定义测试结构体
	type TestStruct struct {
		NormalField     string `json:"normal_field" form:"normal_field" runner:"code:normal_field;name:普通字段" widget:"type:input" data:"type:string"`
		ReadOnlyField   string `json:"readonly_field" form:"readonly_field" runner:"code:readonly_field;name:只读字段" widget:"type:input" data:"type:string" permission:"read"`
		CreateOnlyField string `json:"create_only_field" form:"create_only_field" runner:"code:create_only_field;name:仅创建字段" widget:"type:input" data:"type:string" permission:"read,create"`
		UpdateOnlyField string `json:"update_only_field" form:"update_only_field" runner:"code:update_only_field;name:仅更新字段" widget:"type:input" data:"type:string" permission:"read,update"`
	}

	// 解析结构体
	fields, err := parser.ParseStruct(reflect.TypeOf(TestStruct{}))
	if err != nil {
		t.Fatalf("Failed to parse struct: %v", err)
	}

	// 验证字段数量
	if len(fields) != 4 {
		t.Fatalf("Expected 4 fields, got %d", len(fields))
	}

	// 验证每个字段的permission解析
	testCases := []struct {
		fieldName string
		expected  *PermissionConfig
	}{
		{
			fieldName: "NormalField",
			expected:  nil, // 没有permission标签，应该为nil
		},
		{
			fieldName: "ReadOnlyField",
			expected:  &PermissionConfig{Read: true, Update: false, Create: false},
		},
		{
			fieldName: "CreateOnlyField",
			expected:  &PermissionConfig{Read: true, Update: false, Create: true},
		},
		{
			fieldName: "UpdateOnlyField",
			expected:  &PermissionConfig{Read: true, Update: true, Create: false},
		},
	}

	for _, tc := range testCases {
		// 找到对应的字段
		var field *FieldConfig
		for _, f := range fields {
			if f.FieldName == tc.fieldName {
				field = f
				break
			}
		}

		if field == nil {
			t.Errorf("Field %s not found", tc.fieldName)
			continue
		}

		// 验证permission配置
		if tc.expected == nil {
			if field.Permission != nil {
				t.Errorf("Field %s should have nil permission, got %+v", tc.fieldName, field.Permission)
			}
		} else {
			if field.Permission == nil {
				t.Errorf("Field %s should have permission config, got nil", tc.fieldName)
				continue
			}
			if field.Permission.Read != tc.expected.Read ||
				field.Permission.Update != tc.expected.Update ||
				field.Permission.Create != tc.expected.Create {
				t.Errorf("Field %s permission mismatch. Expected: %+v, Got: %+v",
					tc.fieldName, tc.expected, field.Permission)
			}
		}
	}
}

type WorkflowDemoReq struct {
	// 时间安排 - 演示不同时间格式的使用场景
	StartDate    int64 `json:"start_date" form:"start_date" runner:"code:start_date;name:开始日期" widget:"type:datetime;kind:date;placeholder:请选择开始日期;min_date:$today;format:yyyy-MM-dd" data:"type:number;default_value:$today;example:1705292200000" validate:"required"`  // 业务日期：format指定只显示日期
	DueDate      int64 `json:"due_date" form:"due_date" runner:"code:due_date;name:截止日期" widget:"type:datetime;kind:date;placeholder:请选择截止日期;min_date:$today;format:yyyy-MM-dd" data:"type:number;default_value:$7_days_later;example:1705897600000" validate:"required"` // 业务日期：format指定只显示日期
	Birthday     int64 `json:"birthday" form:"birthday" runner:"code:birthday;name:团队成员生日" widget:"type:datetime;kind:date;placeholder:请选择生日;format:yyyy-MM-dd" data:"type:number;example:694224000000"`                                                                  // 生日：只显示日期，如1992-01-15
	ClassEndTime int64 `json:"class_end_time" form:"class_end_time" runner:"code:class_end_time;name:培训结束时间" widget:"type:datetime;kind:time;placeholder:请选择结束时间;format:HH:mm" data:"type:number;example:1705309800000"`                                                  // 时间：只显示时间，如18:30
	MeetingTime  int64 `json:"meeting_time" form:"meeting_time" runner:"code:meeting_time;name:会议时间" widget:"type:datetime;kind:datetime;placeholder:请选择会议时间" data:"type:number;example:1705467600000"`                                                                   // 默认格式：无format标签，显示完整日期时间
}

func TestName(t *testing.T) {
	of := reflect.TypeOf(WorkflowDemoReq{})
	fields, err := NewMultiTagParser().ParseStruct(of)
	if err != nil {
		t.Fatalf("Failed to parse struct: %v", err)
	}
	fmt.Println(fields)

}
