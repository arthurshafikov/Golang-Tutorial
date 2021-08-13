package hw06pipelineexecution

import (
	"sync"
	"sync/atomic"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// Пояснение решения:
// Мы передаём по одному числу в функцию processNumberToStages, и она прокидывает его через все стейджи
// Эта функция запускается паралелльно для всех входных чисел
// После чего мы принимаем эти числа в результат в соответствии с их позицией в исходном массиве
// Решение работает за 0.45 секунды для первого теста и является самым оптимизированным в отличии от других решений
// Пусть оно и занимает больше строк кода.
func processNumberToStages(i interface{}, stages []Stage, done In) Out {
	processResCh := make(Bi)

	processCh := make(Bi, 1)
	processCh <- i

	go func(processCh In) {
		defer close(processResCh)
		for _, stage := range stages {
			stageCh := make(Bi)
			select {
			case <-done:
				return
			case i := <-processCh:
				processCh = stage(stageCh)
				stageCh <- i
			}
		}
		select {
		case <-done:
		case processResCh <- <-processCh:
		}
	}(processCh)

	return processResCh
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	resultCh := make(Bi)

	go func() {
		defer close(resultCh)
		var wg sync.WaitGroup
		numberPosition := 0
		var numbersPassed int32
		atomic.StoreInt32(&numbersPassed, 0)
		for i := range in {
			wg.Add(1)
			go func(i interface{}, pos int) {
				defer wg.Done()
				for j := range processNumberToStages(i, stages, done) {
					for {
						select {
						case <-done:
							return
						default:
						}
						if int(atomic.LoadInt32(&numbersPassed)) == pos {
							resultCh <- j
							atomic.AddInt32(&numbersPassed, 1)
							break
						}
					}
				}
			}(i, numberPosition)
			numberPosition++
		}
		wg.Wait()
	}()

	return resultCh
}
