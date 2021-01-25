package utils

import "encoding/json"

// strut -> json
func Struct2Json(strut interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(strut)
	if err != nil {
		return nil, err
	}

	r := make(map[string]interface{})

	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}