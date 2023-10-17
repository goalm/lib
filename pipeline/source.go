package pipeline

import (
	"context"
	"sync"
)

type ISource interface {
	Process(ctx context.Context, wg *sync.WaitGroup, errChan chan error) <-chan any
}

// 生成器 输入数据依次放入输出通道
type Source struct {
	Nums []any
}

func NewSource(nums []any) *Source {
	return &Source{Nums: nums}
}

func (t *Source) Process(ctx context.Context, wg *sync.WaitGroup, errChan chan error) <-chan any {
	defer wg.Done()

	outChan := make(chan any)
	go func() {
		defer close(outChan)
		for _, s := range t.Nums {
			//time.Sleep(time.Microsecond * 10)
			select {
			case <-ctx.Done():
				return
			case outChan <- s:
			}
		}
	}()

	return outChan
}
