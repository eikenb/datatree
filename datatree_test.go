package datatree

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var jsonstr = `{"one":{"hi":"bye"},"two":["foo",["bar"]]}`

// Make sure NewMap takes nil arg
func TestNilMap(t *testing.T) {
	s := NewMap(nil)
	assert.Nil(t, s)
}

// Check that deep copy is full copy w/o references to original
func TestDeepCopy(t *testing.T) {
	m := make(Map)
	err := json.Unmarshal([]byte(jsonstr), &m)
	assert.Nil(t, err)

	// fmt.Println(m)
	m1 := m.deepcopy()
	m["one"].(map[string]interface{})["hi"] = "world"
	m["two"].([]interface{})[1].([]interface{})[0] = "zed"
	// fmt.Println(m)
	// fmt.Println(m1)

	assert.NotEqual(t, m, m1)
	assert.Equal(t, m1["one"].(Map)["hi"], "bye")
	assert.Equal(t, m1["two"].(List)[1].(List)[0], "bar")
}

var jsonstr_1 = `{"one":{"hi":"mom"},"two":["zoo",["bar"]]}`
var jsonstr_2 = `{"one":{"hi":"mom"},"two":["zoo",["bar"]],"three":"zed"}`

func newjsonMap(j []byte) Map {
	s := make(Map)
	err := json.Unmarshal([]byte(j), &s)
	if err != nil {
		log.Fatal(err)
	}
	return s.convert()
}

// Test merging of 2 Map
func TestMerge(t *testing.T) {
	m := newjsonMap([]byte(jsonstr))
	n := newjsonMap([]byte(jsonstr_1))
	m1 := m.deepcopy()
	m1.merge(n)
	// fmt.Println(m1)
	assert.Equal(t, m, m1)
	assert.Equal(t, m1["one"].(Map)["hi"], "mom")
	assert.Equal(t, m1["two"].(List)[0], "zoo")

	n = newjsonMap([]byte(jsonstr_2))
	m1.merge(n)
	assert.NotEqual(t, m, m1)
}

// Test comparing of 2 Map
func TestCompare(t *testing.T) {
	m := newjsonMap([]byte(jsonstr))
	n := newjsonMap([]byte(jsonstr))
	o := newjsonMap([]byte(jsonstr_1))

	assert.True(t, m.Compare(n))
	assert.False(t, n.Compare(o))
}

var jsonstr_3 = `{"one":{"hi":"mom"},"two":["zoo",[true]],"three":"zed"}`

// Test Lookup method
func TestLookup(t *testing.T) {
	m := make(Map)
	err := json.Unmarshal([]byte(jsonstr_3), &m)
	if err != nil {
		log.Fatal(err)
	}
	m = m.convert()
	v, ok := m.Lookup("two", 1, 0)
	assert.True(t, ok)
	assert.True(t, v.(bool))

	v, ok = m.Lookup("one", "hi")
	assert.True(t, ok)
	assert.Equal(t, v.(string), "mom")
}
