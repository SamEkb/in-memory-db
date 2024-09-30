package synchronization

type NoOpSemaphore struct{}

func (s *NoOpSemaphore) Acquire() {
}

func (s *NoOpSemaphore) Release() {

}
