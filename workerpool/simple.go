package workerpool

import (
	"runtime"
	"sync"
)

// NewSimplePool returns a worker pool with the same amount of workers as the runtime available CPUs.
func NewSimplePool() WorkerPool {
	workersCnt := runtime.NumCPU()
	return &simplePool{
		workChan:   make(chan Job, workersCnt),
		stopChan:   make(chan bool, workersCnt),
		numWorkers: workersCnt,
	}
}

type simplePool struct {
	workChan   chan Job
	stopChan   chan bool
	numWorkers int
}

func (s *simplePool) Add(t Job) {
	s.workChan <- t
}

func (s *simplePool) Run(w WorkerFunc) <-chan JobResult {
	c := make(chan JobResult)

	var wg sync.WaitGroup

	for i := 0; i < s.numWorkers; i++ {
		wg.Add(1)

		go func(wg *sync.WaitGroup, in <-chan Job, out chan<- JobResult, stop <-chan bool) {
			defer wg.Done()

			select {
			case <-stop:
				return
			default:
				for t := range in {
					out <- w(t)
				}
			}

		}(&wg, s.workChan, c, s.stopChan)
	}

	// when the work is done -> close the res chan.
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(c)
	}(&wg)

	return c
}

func (s *simplePool) Close() error {
	_, ok := <-s.workChan
	if ok {
		close(s.workChan)
	}
	return nil
}

func (s *simplePool) Stop() {
	for i := 0; i < s.numWorkers; i++ {
		s.stopChan <- true
	}
}
