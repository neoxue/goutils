package goutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCopyMapFirstLevel(t *testing.T) {
	avalue := map[string]interface{}{"testkey": "testvalue"}
	a := map[string]interface{}{"key1": "value1", "key2": nil, "key3": 111, "key4": avalue}
	b := CopyMapFirstLevel(a)
	assert.Equal(t, avalue, b["key4"])
	delete(b, "key1")
	assert.NotEqual(t, nil, a["key1"])
	assert.Equal(t, nil, b["key1"])
}

func TestMapKeys(t *testing.T) {
	avalue := map[string]interface{}{"testkey": "testvalue", "testkey2": "simple"}
	mapkeys := MapKeys(avalue)
	assert.Equal(t, []string{"testkey", "testkey2"}, mapkeys)
	assert.Equal(t, mapkeys, []string{"testkey", "testkey2"})
}
