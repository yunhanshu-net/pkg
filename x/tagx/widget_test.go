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
