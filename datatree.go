package datatree

import (
	"encoding/json"
	"log"
	"reflect"
)

// Local types for yaml/json/mxj style interface{} trees
type Map map[string]interface{}
type List []interface{}

// ----------------------------------------------------------------------
// Public functions

// Gimme a new Map
func NewMap(data map[string]interface{}) Map {
	return Map(data).convert()
}

// Function to simplify building test stacks
func TestMap(entries ...interface{}) Map {
	m := make(Map)
	var k string
	for i := 0; i < len(entries); i++ {
		if i%2 == 0 {
			k = entries[i].(string)
		} else {
			m[k] = entries[i]
		}
	}
	return m
}

// ----------------------------------------------------------------------
// Interface functions

// Merge new Map into current Map.
func (m Map) Merge(mm Map) Map {
	return m.merge(mm)
}

// Deep copy of the Map
func (m Map) Copy() Map { return m.deepcopy() }

// Compare 2 stacks
func (m Map) Compare(mm Map) bool {
	return reflect.DeepEqual(m, mm)
}

// Clear/Delete all data from Map
func (m Map) Clear() {
	for k, _ := range m {
		delete(m, k)
	}
}

// Test for key equality in map (used in testing)
func (m Map) KeysEqual(mm Map) bool {
	for k, _ := range m {
		if _, ok := mm[k]; !ok {
			return false
		}
	}
	for k, _ := range mm {
		if _, ok := m[k]; !ok {
			return false
		}
	}
	return true
}

// Lookup a value from deep within Map/List objects
// Returns an interface containing the value and boolean about lookup sucess
func (m Map) Lookup(path ...interface{}) (interface{}, bool) {
	return m.lookup(path...)
}

// Map as JSON
func (m Map) Json() ([]byte, error) {
	return json.Marshal(m)
}

func (m Map) PrettyJson() ([]byte, error) {
	return json.MarshalIndent(m, "", "    ")
}

// Map as a human readable string
func (m Map) String() string {
	j, err := m.Json()
	if err != nil {
		log.Println(err)
	}
	return string(j)
}

// ----------------------------------------------------------------------
// Private functions

// recursive Map tree type converter
func (m Map) convert() Map {
	for k, v := range m {
		switch val := v.(type) {
		case []interface{}:
			m[k] = List(val).convert()
		case map[string]interface{}: // json unmarshal
			m[k] = Map(val).convert()
		case map[interface{}]interface{}: // yaml unmarshal
			m[k] = Map(mii2msi(val)).convert()
		}
	}
	return m
}

// recursive List type converter
func (l List) convert() List {
	for k, v := range l {
		switch val := v.(type) {
		case []interface{}:
			l[k] = List(val).convert()
		case map[string]interface{}:
			l[k] = Map(val).convert()
		}
	}
	return l
}

// fix yaml's map[interface{}]interface{} default conversion
func mii2msi(m map[interface{}]interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	for k, v := range m {
		r[k.(string)] = v
	}
	return r
}

// merge 2 stacks
func (m Map) merge(mm Map) Map {
	for k, v := range mm {
		switch val := v.(type) {
		default:
			m[k] = val
		case List:
			if v, ok := m[k]; ok {
				m[k] = val.deepcopyto(v.(List))
			} else {
				m[k] = val.deepcopy()
			}
		case Map:
			if v, ok := m[k]; ok {
				m[k] = val.deepcopyto(v.(Map))
			} else {
				m[k] = val.deepcopy()
			}
		}
	}
	return m
}

// deep copy of Map (json/yaml data tree)
func (m Map) deepcopy() Map {
	dst := make(Map)
	return m.deepcopyto(dst)
}

// deep copy into passed in Map
func (m Map) deepcopyto(dst Map) Map {
	for k, v := range m {
		switch val := v.(type) {
		default:
			dst[k] = val
		case []interface{}:
			dst[k] = List(val).deepcopy()
		case map[string]interface{}:
			dst[k] = Map(val).deepcopy()
		}
	}
	return dst
}

// deep copy of Lists (slice)
func (l List) deepcopy() List {
	dst := make(List, len(l))
	return l.deepcopyto(dst)
}

// deep copy of passed in List
func (l List) deepcopyto(dst List) List {
	for k, v := range l {
		switch val := v.(type) {
		default:
			dst[k] = val
		case []interface{}:
			dst[k] = List(val).deepcopy()
		case map[string]interface{}:
			dst[k] = Map(val).deepcopy()
		}
	}
	return dst
}

// Lookup data from (deep) in a map/list structure
func (m Map) lookup(path ...interface{}) (interface{}, bool) {
	for i, k := range path {
		ks, ok := k.(string)
		var v interface{}
		if ok {
			v, ok = m[ks]
		}
		if ok {
			switch tv := v.(type) {
			case string, bool, float64:
				return tv, true
			case Map:
				return tv.lookup(path[i+1:]...)
			case List:
				return tv.lookup(path[i+1:]...)
			}
		}
	}
	return nil, false
}

func (l List) lookup(path ...interface{}) (interface{}, bool) {
	for i, k := range path {
		ki, ok := k.(int)
		if ok && len(l) > ki {
			v := l[ki]
			switch tv := v.(type) {
			case string, bool, float64:
				return tv, true
			case Map:
				return tv.lookup(path[i+1:]...)
			case List:
				return tv.lookup(path[i+1:]...)
			}
		}
	}
	return nil, false
}
