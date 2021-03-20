package errgroup

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Pool interface {
	NewGroup() Group

	Go(func())

	RunningCount() int
}

type Group interface {
	Go(func() error)
	Wait() []error
}

type pool struct {
	coreSize int
	maxSize  int
	overTime time.Duration

	tokenCh           chan struct{} // token 池
	atomicCurrentJobs int32         // 当前运行任务数

	defaultGroup Group
}

type group struct {
	p    *pool
	done sync.WaitGroup

	errs    []error
	errLock sync.Locker
}

func NewPool(coreSize, maxSize int, overTime time.Duration) Pool {
	p := pool{
		coreSize:          coreSize,
		maxSize:           maxSize,
		overTime:          overTime,
		tokenCh:           make(chan struct{}, coreSize),
		atomicCurrentJobs: 0,
	}
	p.defaultGroup = p.NewGroup()
	return &p
}

func (p *pool) NewGroup() Group {
	return &group{p: p, errLock: &sync.Mutex{}}
}

func (p *pool) RunningCount() int {
	return int(atomic.LoadInt32(&p.atomicCurrentJobs))
}

func (p *pool) getToken()     { p.tokenCh <- struct{}{} }
func (p *pool) releaseToken() { <-p.tokenCh }
func (p *pool) addJobCount()  { atomic.AddInt32(&p.atomicCurrentJobs, 1) }
func (p *pool) subJobCount()  { atomic.AddInt32(&p.atomicCurrentJobs, -1) }

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

func (p *pool) Go(fn func()) {
	p.defaultGroup.Go(func() error {
		fn()
		return nil
	})
}

func (g *group) Go(fn func() error) {
	g.done.Add(1)
	go func() {
		g.p.getToken()
		g.p.addJobCount()
		once := sync.Once{}
		done := make(chan struct{}, 1)
		go g.p.releaseWhenOvertime(done, &once)
		defer func() {
			done <- struct{}{}
			g.p.subJobCount()
			g.done.Done()
			once.Do(g.p.releaseToken)
			if err := recover(); err != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				g.errLock.Lock()
				g.errs = append(g.errs, fmt.Errorf("errgroup: panic recovered: %s\n%s", err, buf))
				g.errLock.Unlock()
			}
		}()
		if err := fn(); err != nil {
			g.errLock.Lock()
			g.errs = append(g.errs, err)
			g.errLock.Unlock()
		}
	}()
}

func (g *group) Wait() []error {
	g.done.Wait()
	g.errLock.Lock()
	defer g.errLock.Unlock()
	return g.errs
}
