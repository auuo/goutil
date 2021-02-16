package lru

import "testing"

func TestLRU(t *testing.T) {
	lru := New(5)
	input := []struct {
		cmdType int // 0: add, 1: get
		key     string
		value   interface{}
		expect  []interface{}
	}{
		{key: "0", value: 0, expect: []interface{}{0}},
		{key: "1", value: 1, expect: []interface{}{0, 1}},
		{key: "3", value: 3, expect: []interface{}{0, 1, 3}},
		{cmdType: 1, key: "1", expect: []interface{}{0, 3, 1}},
		{key: "3", value: 3, expect: []interface{}{0, 1, 3}},
		{cmdType: 1, key: "0", expect: []interface{}{1, 3, 0}},
		{key: "4", value: 4, expect: []interface{}{1, 3, 0, 4}},
		{key: "5", value: 5, expect: []interface{}{1, 3, 0, 4, 5}},
		{key: "6", value: 6, expect: []interface{}{3, 0, 4, 5, 6}},
		{cmdType: 1, key: "3", expect: []interface{}{0, 4, 5, 6, 3}},
		{key: "7", value: 7, expect: []interface{}{4, 5, 6, 3, 7}},
		{cmdType: 1, key: "120", expect: []interface{}{4, 5, 6, 3, 7}},
	}
	for _, it := range input {
		if it.cmdType == 0 {
			lru.Add(it.key, it.value)
		} else {
			lru.Get(it.key)
		}

		var res []interface{}
		length := lru.(*cache).list.Len()
		ele := lru.(*cache).list.Front()
		for i := 0; i < length; i++ {
			res = append(res, ele.Value.(*item).value)
			ele = ele.Next()
		}
		if len(it.expect) != len(res) {
			t.Errorf("expect (%v) but (%v)", it.expect, res)
		}
		for idx := range it.expect {
			if it.expect[idx] != res[idx] {
				t.Errorf("expect (%v) but (%v)", it.expect, res)
			}
		}
	}
}
