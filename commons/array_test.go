package commons

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInArray(t *testing.T) {
	item1 := 1
	item2 := "2"
	item3 := "c"
	item4 := 'c'
	arr := []interface{}{1, "2", "c", 'c', 'd'}
	assert.Equal(t, true, InArray(item1, arr))
	assert.Equal(t, true, InArray(item2, arr))
	assert.Equal(t, true, InArray(item3, arr))
	assert.Equal(t, true, InArray(item4, arr))

	assert.Equal(t, false, InArray("test", arr))

	arr2 := []interface{}{"cTime"}
	var a interface{}
	a = "cTime"
	assert.Equal(t, true, InArray(a.(string), arr2))
}
