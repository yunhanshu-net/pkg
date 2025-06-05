package jsonx

import "encoding/json"

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
