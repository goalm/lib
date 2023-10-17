package pipeline

import (
	"context"
	"log"
	"runtime"
	"sync"
)

type StageFunc func(s any)

type IProcessor interface {
	Process(ctx context.Context, wg *sync.WaitGroup, dataChan <-chan any, errChan chan error) <-chan any
}

// Single Prog
type SingleProcessor struct {
	Proc StageFunc
}

func (p *SingleProcessor) Process(ctx context.Context, wg *sync.WaitGroup, dataChan <-chan any, errChan chan error) <-chan any {
	outChan := make(chan any)

	go func() {
		defer wg.Done()
		defer close(outChan)

		for {
			select {
			case s, ok := <-dataChan:
				if !ok {
					log.Println("Single data channel closed!")
					return
				}
				p.Proc(s)
				outChan <- s
			}
		}
	}()

	return outChan
}

// Concurrent Prog
type MultiProcessor struct {
	Proc StageFunc
}

func (p *MultiProcessor) Process(ctx context.Context, wg *sync.WaitGroup, dataChan <-chan any, errChan chan error) <-chan any {
	outChan := make(chan any)
	maxCnt := 10

	go func() {
		defer wg.Done()
		defer close(outChan)

		wg2 := sync.WaitGroup{}
		wg2.Add(maxCnt)
		for i := 0; i < maxCnt; i++ {
			go func() {
				defer wg2.Done()
				for {
					select {
					case s, ok := <-dataChan:
						if !ok {
							log.Println("Parallel processing channel closed! gouroutine", runtime.NumGoroutine())
							return
						}
						p.Proc(s)
						outChan <- s
					}
				}
			}()
		}

		wg2.Wait()

	}()

	return outChan
}
