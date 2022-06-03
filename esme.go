/*
esme.go
sam
2022-04-25
*/

package esme

import (
	"sync"

	"github.com/zituocn/esme/logx"
)

// Job job struct
type Job struct {

	// task name
	name string

	// num number of goroutines
	num int

	// queue queue of tasks waiting to be executed
	queue TodoQueue

	// jobOptions job options
	jobOptions JobOptions
}

// JobOptions 任务参数
type JobOptions struct {

	// StartFunc Callback at start
	StartFunc CallbackFunc

	// SucceedFunc  Callback after success
	SucceedFunc CallbackFunc

	// RetryFunc retry callback
	RetryFunc CallbackFunc

	// FailedFunc callback after failure
	FailedFunc CallbackFunc

	// CompleteFunc Callback for request completion
	CompleteFunc CallbackFunc

	// ProxyIP proxy ip
	ProxyIP string

	// ProxyLib proxy ip library
	ProxyLib *ProxyLib

	// SheepTime Sleep time for http request execution
	// millisecond
	SheepTime int

	// TimeOut http request timeout
	// millisecond
	TimeOut int

	// 是否打印调试
	IsDebug bool
}

// NewJob returns a  *Job
func NewJob(name string, num int, queue TodoQueue, options JobOptions) *Job {
	if num < 1 {
		num = 1
	}
	return &Job{
		name:       name,
		num:        num,
		queue:      queue,
		jobOptions: options,
	}
}

// Do start the job
func (j *Job) Do() {

	logx.Infof("[%s] start job -> Goroutines : %d ", j.name, j.num)

	var wg sync.WaitGroup
	for n := 0; n < j.num; n++ {
		wg.Add(1)
		go func(i int) {
			logx.Infof("start task %d", i+1)
			defer wg.Done()
			for {
				if j.queue.IsEmpty() {
					break
				}
				task := j.queue.Pop()
				if task != nil {
					ctx := DoRequest(task.Url, task.Method, task.Header, task.FormData, task.Payload, task)

					ctx.SetStartFunc(j.jobOptions.StartFunc).
						SetSucceedFunc(j.jobOptions.SucceedFunc).
						SetRetryFunc(j.jobOptions.RetryFunc).
						SetFailedFunc(j.jobOptions.FailedFunc).
						SetCompleteFunc(j.jobOptions.CompleteFunc).
						SetIsDebug(j.jobOptions.IsDebug).
						SetTimeOut(j.jobOptions.TimeOut).
						SetSleepTime(j.jobOptions.SheepTime).
						SetProxy(j.jobOptions.ProxyIP).
						SetProxyLib(j.jobOptions.ProxyLib)

					// execute request
					ctx.Do()
				}
			}

		}(n)
	}
	wg.Wait()

	logx.Info("job done")
}
