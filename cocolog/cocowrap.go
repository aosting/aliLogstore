package cocolog

import (
	"github.com/aliyun/aliyun-log-go-sdk"
	"sync"
)

/**********************************************
** @Des: cocowrap  线程安全日志包，默认100条组成一个包
** @Author: zhangxueyuan 
** @Date:   2019-01-22 11:51:27
** @Last Modified by:   zhangxueyuan 
** @Last Modified time: 2019-01-22 11:51:27
***********************************************/

const basecapcity = 100

type wrap struct {
	put    []*sls.Log //初始化注意设置长度，防止频繁扩容
	length int
	mutex  sync.RWMutex //指针copy相当于无效
}

func (cs *wrap) appends(puts []*sls.Log) {
	cs.mutex.Lock()
	if cs == nil {
		cs = initwrap(basecapcity)
	}
	cs.length = cs.length + len(puts)
	cs.put = append(cs.put, puts...)
	cs.mutex.Unlock()

}

// Appends an item to the concurrent slice
func (cs *wrap) append(put *sls.Log) {
	cs.mutex.Lock()
	if cs == nil {
		cs = initwrap(basecapcity)
	}
	cs.length = cs.length + 1
	cs.put = append(cs.put, put)
	cs.mutex.Unlock()
}

//注意线程安装
func initwrap(capcity int) *wrap {
	return &wrap{
		put:   make([]*sls.Log, 0, capcity),
		mutex: sync.RWMutex{},
	}
}
func (cs *wrap) clear(capcity int) {
	cs.length = 0
	cs.put = make([]*sls.Log, 0, capcity)
}

func (cs *wrap) set(put []*sls.Log) {
	cs.mutex.Lock()
	if cs == nil {
		cs = initwrap(basecapcity)
	}
	cs.put = put
	cs.length = len(put)
	cs.mutex.Unlock()
}

func (cs *wrap) size() int {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	return cs.length
}

func (cs *wrap) copy(capcity int) *wrap {
	cs.mutex.Lock()
	copycs := &wrap{
		put:    cs.put,
		length: cs.length,
		mutex:  sync.RWMutex{},
	}
	cs.clear(capcity)
	cs.mutex.Unlock()
	return copycs
}
func (cs *wrap) sub() (wrap, wrap) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	half := cs.length / 2
	halfwrap := wrap{
		put:    cs.put[:half],
		length: half,
		mutex:  sync.RWMutex{},
	}
	halfwrap2 := wrap{
		put:    cs.put[half:],
		length: cs.length - half,
		mutex:  sync.RWMutex{},
	}
	return halfwrap, halfwrap2

}

//func (cs *wrap) getSizeBefore(size int) []*sls.Log {
//	cs.Lock()
//	defer cs.Unlock()
//	if cs == nil {
//		cs = &wrap{}
//		cs.put = make([]*sls.Log, 0)
//	}
//	out := cs.put[:size]
//	cs.put = cs.put[size:]
//	return out
//}

//func (cs *wrap) getSizeAfter(size int) {
//	cs.Lock()
//	defer cs.Unlock()
//	if cs == nil {
//		cs = &wrap{}
//		cs.put = make([]*sls.Log, 0)
//	}
//	cs.put = cs.put[size:]
//}

//func (cs *wrap) getAll() []*sls.Log {
//	cs.Lock()
//	defer cs.Unlock()
//	if cs == nil {
//		cs = &wrap{}
//		cs.put = make([]*sls.Log, 0)
//	}
//	tmp := cs.put
//	cs.put = make([]*sls.Log, 0)
//	return tmp
//}
