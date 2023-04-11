package workerpool

import (
	"fmt"
	"github.com/linqcod/cbr-currency-app/cmd/internal/currency/model"
	"log"
	"sync"
)

type Worker struct {
	ID       int
	taskChan chan *Task
	resChan  chan *model.ValResults
}

func NewWorker(tc chan *Task, rc chan *model.ValResults, ID int) *Worker {
	return &Worker{
		ID:       ID,
		taskChan: tc,
		resChan:  rc,
	}
}

func (w *Worker) Start(wg *sync.WaitGroup) {
	fmt.Printf("Starting worker %d\n", w.ID)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range w.taskChan {
			valCurs, err := task.process(w.ID)
			if err != nil {
				log.Printf("error while processing task by worker %d: %v", w.ID, err)
			}
			if len(valCurs.Records) != 0 {
				res := &model.ValResults{
					ValCurs:    valCurs,
					CurrencyId: task.currencyId,
				}

				w.resChan <- res
			}
		}
		fmt.Printf("Worker %d ended jobs! \n", w.ID)
	}()
}
