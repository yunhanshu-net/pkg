package jsonx

import (
	"encoding/json"
	"reflect"
)

func Convert(el interface{}, resetEl interface{}) error {
	marshal, err := json.Marshal(el)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, resetEl)
	if err != nil {
		return err
	}
	return nil
}

func EQ(s1, s2 string) bool {
	var mp1 interface{}
	var mp2 interface{}
	err := json.Unmarshal([]byte(s1), &mp1)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(s2), &mp2)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(mp1, mp2)
}

func EQRawMessage(s1, s2 json.RawMessage) bool {
	var mp1 interface{}
	var mp2 interface{}
	err := json.Unmarshal(s1, &mp1)
	if err != nil {
		return false
	}
	err = json.Unmarshal(s2, &mp2)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(mp1, mp2)
}
