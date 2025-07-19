// Package jsonx ...
package jsonx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/pkg/logger"
)

// MustJSON ..
func MustJSON(el interface{}) string {
	marshal, err := json.Marshal(el)
	if err != nil {
		panic(err)
	}
	return string(marshal)
}

// MustPrintJSON ...
func MustPrintJSON(el interface{}) {
	marshal, err := json.Marshal(el)
	if err != nil {
		fmt.Println(fmt.Sprintf("[jsonx] err:%s el:%+v", err.Error(), el))
		return
	}
	fmt.Println(string(marshal))
}

// JSONString ...
func JSONString(el interface{}) string {
	marshal, err := json.Marshal(el)
	if err != nil {
		return ""
	}
	return string(marshal)
}

// String ...
func String(el interface{}) string {
	marshal, err := json.Marshal(el)
	if err != nil {
		logger.Errorf(context.Background(), "[jsonx] err:%s el:%+v", err.Error(), el)
		return ""
	}
	return string(marshal)
}
