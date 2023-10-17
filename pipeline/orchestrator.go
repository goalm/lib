package pipeline

import (
	"context"
	"sync"
)

type ProcessorManager struct {
	source  ISource
	sink    ISink
	err     IError
	ps      []IProcessor
	errChan chan error
}

func NewProcessorManager() *ProcessorManager {
	return &ProcessorManager{errChan: make(chan error, 1)}
}

func (m *ProcessorManager) AddProcessor(processor IProcessor) {
	m.ps = append(m.ps, processor)
}

func (m *ProcessorManager) AddSource(source ISource) {
	m.source = source
}

func (m *ProcessorManager) AddSink(sink ISink) {
	m.sink = sink
}

func (m *ProcessorManager) AddError(err IError) {
	m.err = err
}

func (m *ProcessorManager) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg = sync.WaitGroup{}

	// 组装pipeline
	wg.Add(1)
	dataChan := m.source.Process(ctx, &wg, m.errChan)

	for _, v := range m.ps {
		wg.Add(1)
		dataChan = v.Process(ctx, &wg, dataChan, m.errChan)
	}

	wg.Add(1)
	m.sink.Process(ctx, &wg, dataChan, m.errChan)
	//wg.Add(len(mw))
	wg.Add(1)
	m.sink.Output(ctx, &wg, m.errChan)

	go func() {
		wg.Wait()
		close(m.errChan)
	}()

	m.err.Process(ctx, &wg, m.errChan, cancel)
}
