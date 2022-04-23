package esme

import (
	"fmt"
	"sync"
)

// HttpQueue http请求的队列
type HttpQueue struct {
	mux  *sync.Mutex
	list []*Task
}

// NewHttpQueue return a http queue obj
func NewHttpQueue() TodoQueue {
	return &HttpQueue{
		list: make([]*Task, 0),
		mux:  &sync.Mutex{},
	}
}

// Add 添加一个任务
func (q *HttpQueue) Add(task *Task) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.list = append(q.list, task)
}

// Pop 弹出第并获取第一条
func (q *HttpQueue) Pop() *Task {
	q.mux.Lock()
	defer q.mux.Unlock()

	if q.IsEmpty() {
		return nil
	}

	first := q.list[0]
	q.list = q.list[1:]
	return first
}

// Clear 清理所有
func (q *HttpQueue) Clear() bool {
	if q.IsEmpty() {
		return false
	}

	for i := 0; i < q.Size(); i++ {
		q.list[i].Url = ""
	}
	q.list = nil
	return true
}

// IsEmpty 是否为空
//	return bool
func (q *HttpQueue) IsEmpty() bool {
	if len(q.list) == 0 {
		return true
	}
	return false
}

// Size 返回长度
func (q *HttpQueue) Size() int {
	return len(q.list)
}

// Print 打印输出
func (q *HttpQueue) Print() {
	fmt.Println(q.list)
}
