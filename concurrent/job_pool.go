package concurrent

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Job func() error

type JobPool interface {
	Submit(Job)

	Wait() []error

	RunningCount() int
}

type pool struct {
	coreSize int
	maxSize  int
	overTime time.Duration

	jobs              chan Job       // 任务池
	tokenCh           chan struct{}  // token 池子
	atomicCurrentJobs int32          // 当前运行任务数
	done              sync.WaitGroup // 任务结束 wg

	errs    []error
	errLock sync.Mutex
}

func NewJobPool(coreSize, maxSize int, overTime time.Duration) JobPool {
	p := pool{
		coreSize:          coreSize,
		maxSize:           maxSize,
		overTime:          overTime,
		jobs:              make(chan Job, maxSize),
		tokenCh:           make(chan struct{}, coreSize),
		atomicCurrentJobs: 0,
	}
	go p.run() // todo 会有一个协程一直运行，想办法解决 (close)
	return &p
}

func (p *pool) Submit(job Job) {
	p.done.Add(1)
	go func() { p.jobs <- job }()
}

func (p *pool) Wait() []error {
	p.done.Wait()
	p.errLock.Lock()
	defer p.errLock.Unlock()
	return p.errs
}

func (p *pool) RunningCount() int {
	return int(atomic.LoadInt32(&p.atomicCurrentJobs))
}

func (p *pool) getToken()     { p.tokenCh <- struct{}{} }
func (p *pool) releaseToken() { <-p.tokenCh }
func (p *pool) addJobCount()  { atomic.AddInt32(&p.atomicCurrentJobs, 1) }
func (p *pool) subJobCount()  { atomic.AddInt32(&p.atomicCurrentJobs, -1) }

func (p *pool) run() {
	for job := range p.jobs {
		p.runJob(job)
	}
}

func (p *pool) runJob(job Job) {
	once := sync.Once{}
	done := make(chan struct{}, 1)
	go func() {
		p.getToken()
		p.addJobCount()
		go p.releaseWhenOvertime(done, &once)
		defer func() {
			done <- struct{}{}
			p.subJobCount()
			p.done.Done()
			once.Do(p.releaseToken)
			if err := recover(); err != nil {
				p.errLock.Lock()
				p.errs = append(p.errs, fmt.Errorf("panic: %v", err))
				p.errLock.Unlock()
			}
		}()
		// run job.
		if err := job(); err != nil {
			p.errLock.Lock()
			p.errs = append(p.errs, err)
			p.errLock.Unlock()
		}
	}()
}

func (p *pool) releaseWhenOvertime(done chan struct{}, once *sync.Once) {
	if p.overTime <= 0 {
		return
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), p.overTime)
	select {
	case <-done:
		cancelFunc()
	case <-ctx.Done():
		// 超时, 启动超发任务
		if p.RunningCount() < p.maxSize {
			once.Do(p.releaseToken)
		}
		cancelFunc()
	}
}
