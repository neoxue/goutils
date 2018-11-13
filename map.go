package goutils

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
