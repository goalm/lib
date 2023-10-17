package pipeline

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"sync"

	"github.com/goalm/app/dcs/config"
)

type SinkFunc func(s any) string

type ISink interface {
	Process(ctx context.Context, wg *sync.WaitGroup, dataChan <-chan any, errChan chan error)
	Output(ctx context.Context, wg *sync.WaitGroup, errChan chan error)
}

// 输出到命令行
type ConsoleSink struct {
}

func NewConsoleSink() *ConsoleSink {
	return &ConsoleSink{}
}

func (s *ConsoleSink) Process(ctx context.Context, wg *sync.WaitGroup, dataChan <-chan any, errChan chan error) {
	go func() {
		defer wg.Done()
		for {
			select {
			case val, ok := <-dataChan:
				if ok {
					fmt.Printf("sink value: %v\n", cacheData(val))
				} else {
					log.Println("sink data channel closed!")
					return
				}
			}
		}
	}()
}

type ModelSink struct {
	Fn SinkFunc
}

func NewModelSink(fn SinkFunc) *ModelSink {
	return &ModelSink{fn}
}

func (s *ModelSink) Process(ctx context.Context, wg *sync.WaitGroup, dataChan <-chan any, errChan chan error) {
	go func() {
		defer wg.Done()
		for {
			select {
			case val, ok := <-dataChan:
				if ok {
					// dispatch to different files
					pName := s.Fn(val)
					config.Writers[pName].Chn <- cacheData(val)

				} else {
					log.Println("sink data channel closed!")
					// close all channels
					for _, v := range config.Writers {
						close(v.Chn)
					}
					return
				}
			}
		}
	}()
}

func (s *ModelSink) Output(ctx context.Context, wg *sync.WaitGroup, errChan chan error) {
	go func() {
		defer wg.Done()
		wg2 := sync.WaitGroup{}
		output := func(w *config.Writer) {
			for v := range w.Chn {
				_, err := w.File.WriteString(v + "\r\n")
				if err != nil {
					panic(err)
				}
				//fmt.Println(w.Name, w.File, w.Chn)
			}
			wg2.Done()
		}

		wg2.Add(len(config.Writers))
		for _, w := range config.Writers {
			go output(w)
		}
		go func() {
			wg2.Wait()
		}()

	}()
}

// shared function
func cacheData(a any) string {
	val := reflect.ValueOf(a)
	typ := reflect.TypeOf(a)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	fields := "*"
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		res := f.String()

		switch f.Type().String() {
		case "int":
			res = strconv.Itoa(f.Interface().(int))
			if res == "" {
				res = "0"
			}
		case "float64":
			res = strconv.FormatFloat(f.Interface().(float64), 'f', 2, 64)
			if res == "" {
				res = "0.0"
			}

		case "string":
			if res == "" {
				res = "-"
			}
			res = `"` + res + `"`
			//fmt.Println(res)
		}

		fields = fields + "," + res
	}
	return fields
}
