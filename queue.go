/*
queue.go
sam
2022-04-25
*/

package esme

// TodoQueue interface
type TodoQueue interface {

	// Add add an element to the queue
	Add(task *Task)

	// AddTasks Add multiple elements to the queue at once
	AddTasks(list []*Task)

	// Pop pop and return an element
	Pop() *Task

	// Clear empty the queue
	Clear() bool

	// Size queue length
	Size() int

	// IsEmpty is empty
	IsEmpty() bool

	// Print print
	Print()
}
