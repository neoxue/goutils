package commons

import "strconv"

func InArray(item interface{}, arr []interface{}) bool {
	for _, v := range arr {
		if item == v {
			return true
		}
	}
	return false
}

func NumberInArray(item float64, arr []interface{}) bool {
	for _, v := range arr {
		switch v.(type) {
		case int:
			if item == float64(v.(int)) {
				return true
			}
		case float32:
			if item == float64(v.(float32)) {
				return true
			}
		case string:
			if vfloat, err := strconv.ParseFloat(v.(string), 64); err == nil && vfloat == item {
				return true
			}
		case float64:
			if item == v.(float64) {
				return true
			}
		}
	}
	return false
}
