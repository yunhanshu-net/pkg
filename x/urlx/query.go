package urlx

import "net/url"

func QueryToMap(query string) map[string]interface{} {
	values, err := url.ParseQuery(query)
	if err != nil {
		return make(map[string]interface{})
	}
	params := make(map[string]interface{})
	for key, value := range values {
		params[key] = value
	}
	return params
}
