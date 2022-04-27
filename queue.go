/*
队列
queue.go
*/

package esme

// TodoQueue 实现接口
type TodoQueue interface {

	// Add 向队列中添加元素
	Add(task *Task)

	// AddTasks 一次向队列中添加多个元素
	AddTasks(list []*Task)

	// Pop 弹出并返回一个元素
	Pop() *Task

	// Clear 清空队列
	Clear() bool

	// Size 队列长度
	Size() int

	// IsEmpty 判断是否为空
	IsEmpty() bool

	// Print 打印
	Print()
}
