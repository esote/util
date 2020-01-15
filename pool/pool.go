// Package pool provides a simple, versatile worker pool implementation.
package pool

import (
	"sync"
	"sync/atomic"
)

// Pool is a worker pool.
type Pool struct {
	done uint64
	jobs chan job
	wg   sync.WaitGroup
}

type job struct {
	f    func(args ...interface{})
	args []interface{}
}

// New constructs a new worker pool. Panics if workers <= 0.
func New(workers, backlog int) *Pool {
	if workers <= 0 {
		panic("pool: not enough workers")
	}

	p := &Pool{
		jobs: make(chan job, backlog),
	}

	p.wg.Add(workers)

	for i := 0; i < workers; i++ {
		go p.worker()
	}

	return p
}

// Enlist the worker pool with a new job. The arguments will be passed to f. If
// blocking, this function waits until a worker has accepted the job. This
// function returns whether a worker was available to accept the job.
func (p *Pool) Enlist(block bool, f func(args ...interface{}), args ...interface{}) bool {
	if block {
		p.jobs <- job{f, args}
		return true
	}

	select {
	case p.jobs <- job{f, args}:
		return true
	default:
		return false
	}
}

// Close the worker pool. Waits for the job backlog to be consumed. Further use
// of the pool is undefined behavior. When flushing, the jobs in the job backlog
// are ignored and simply consumed rather than called.
func (p *Pool) Close(flush bool) {
	var n uint64
	if flush {
		n = 1
	} else {
		defer atomic.SwapUint64(&p.done, 1)
	}
	if !atomic.CompareAndSwapUint64(&p.done, 0, n) {
		panic("pool: close on closed pool")
	}
	close(p.jobs)
	p.wg.Wait()
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for job := range p.jobs {
		if atomic.CompareAndSwapUint64(&p.done, 1, 1) {
			continue
		}
		job.f(job.args...)
	}
}
