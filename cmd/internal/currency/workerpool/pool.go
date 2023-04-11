package workerpool

import (
	"github.com/linqcod/cbr-currency-app/cmd/internal/currency/model"
	"sync"
)

type Pool struct {
	Tasks []*Task

	workersCount int
	collector    chan *Task
	resultChan   chan *model.ValResults
	wg           sync.WaitGroup
}

func NewPool(tasks []*Task, resultChan chan *model.ValResults, workersCount int) *Pool {
	return &Pool{
		Tasks:        tasks,
		workersCount: workersCount,
		collector:    make(chan *Task, len(tasks)),
		resultChan:   resultChan,
	}
}

func (p *Pool) Run() {
	for i := 1; i <= p.workersCount; i++ {
		worker := NewWorker(p.collector, p.resultChan, i)
		worker.Start(&p.wg)
	}

	for i := range p.Tasks {
		p.collector <- p.Tasks[i]
	}

	close(p.collector)

	p.wg.Wait()
}
