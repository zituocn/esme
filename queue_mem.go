package esme

import (
	"fmt"
	"sync"
)

// MemQueue 内存中的队列
type MemQueue struct {
	mux  *sync.Mutex
	list []*Task
}

// NewMemQueue return a memory queue obj
func NewMemQueue() TodoQueue {
	return &MemQueue{
		list: make([]*Task, 0),
		mux:  &sync.Mutex{},
	}
}

// Add 添加一个任务
func (q *MemQueue) Add(task *Task) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.list = append(q.list, task)
}

// AddTasks 一次添加多个任务
func (q *MemQueue) AddTasks(list []*Task) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.list = append(q.list, list...)
}

// Pop 弹出第并获取第一条
func (q *MemQueue) Pop() *Task {
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
func (q *MemQueue) Clear() bool {
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
func (q *MemQueue) IsEmpty() bool {
	if len(q.list) == 0 {
		return true
	}
	return false
}

// Size 返回长度
func (q *MemQueue) Size() int {
	return len(q.list)
}

// Print 打印输出
func (q *MemQueue) Print() {
	fmt.Println(q.list)
}
