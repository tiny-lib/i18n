package pool

import "sync"

var PayloadPool = sync.Pool{New: func() interface{} {
	return make(map[string]string, 2)
}}
