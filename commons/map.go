package commons

/*
CopyMapFirstLevel returns a new map which copys the first level keys
*/
func CopyMapFirstLevel(a map[string]interface{}) map[string]interface{} {
	b := make(map[string]interface{})
	for key := range a {
		b[key] = a[key]
	}
	return b
}

func MapKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
