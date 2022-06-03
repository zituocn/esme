/*
queue_mem.go
内存队列实现
sam
*/

package esme

import (
	"fmt"
	"sync"
)

// MemQueue in-memory queue
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

// Add add a task
func (q *MemQueue) Add(task *Task) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.list = append(q.list, task)
}

// AddTasks add multiple tasks at once
func (q *MemQueue) AddTasks(list []*Task) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.list = append(q.list, list...)
}

// Pop get the first task
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

// Clear clear queue
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

// IsEmpty is empty
//	return bool
func (q *MemQueue) IsEmpty() bool {
	if len(q.list) == 0 {
		return true
	}
	return false
}

// Size returns queue length
func (q *MemQueue) Size() int {
	return len(q.list)
}

// Print print
func (q *MemQueue) Print() {
	fmt.Println(q.list)
}
