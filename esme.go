package esme

import (
	"sync"

	"github.com/zituocn/esme/logx"
)

// Job 任务
type Job struct {

	// 任务name 名称
	name string

	// num 协程数量
	num int

	// queue 等待执行的任务队列
	queue TodoQueue

	// jobOptions 任务选项
	jobOptions JobOptions
}

// JobOptions 任务参数
type JobOptions struct {
	StartFunc CallbackFunc

	// SucceedFunc 成功后的回调
	SucceedFunc CallbackFunc

	// RetryFunc 重试的回调
	RetryFunc CallbackFunc

	// FailedFunc 失败后的回调
	FailedFunc CallbackFunc

	// CompleteFunc 请求完成的回调
	CompleteFunc CallbackFunc

	// ProxyIP 代理IP
	ProxyIP string

	// ProxyLib 代理IP库
	ProxyLib *ProxyLib

	// SheepTime http 请求 执行的休眠时间
	SheepTime int

	// TimeOut http 请求超时时间
	// 毫秒
	TimeOut int

	// 是否打印调试
	IsDebug bool
}

// NewJob 返回一个 *Job
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

// Do 开始执行任务
func (j *Job) Do() {

	logx.Infof("[%s] 开始执行 -> 协程数: %d ", j.name, j.num)

	var wg sync.WaitGroup
	for n := 0; n < j.num; n++ {
		wg.Add(1)
		go func(i int) {
			logx.Infof("启动第 %d 个任务", i+1)
			defer wg.Done()
			for {
				if j.queue.IsEmpty() {
					break
				}
				task := j.queue.Pop()
				if task != nil {
					ctx := DoRequest(task.Url, task.Method, task.Header, task.FormData, task.Playload, task)

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

					// 执行请求
					ctx.Do()
				}
			}

		}(n)
	}
	wg.Wait()

	logx.Infof("任务执行结束")
}
