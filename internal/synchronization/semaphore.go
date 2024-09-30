package synchronization

import "sync"

type ISemaphore interface {
	Acquire()
	Release()
}

type Semaphore struct {
	cond  *sync.Cond
	max   int
	count int
}

func NewSemaphore(max int) ISemaphore {
	if max == 0 {
		return &NoOpSemaphore{}
	}

	m := &sync.Mutex{}
	return &Semaphore{
		cond: sync.NewCond(m),
		max:  max,
	}
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for s.count >= s.max {
		s.cond.Wait()
	}

	s.count++
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.count--
	s.cond.Signal()
}
