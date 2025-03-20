package utils

import (
	"crypto/rand"

	"github.com/rs/zerolog/log"

	"sync"
	"worker-service/pkg/constant"
)

type WorkerPool struct {
	NumberOfWorkers  int
	NumberOfWorkings int
	Mux              *sync.Mutex
}

func (p *WorkerPool) Start(handler constant.WorkerHandleFunc, numberOfRetry int, cmd any) *Worker {
	if p.IsWorkerPoolFull() {
		log.Error().Msg("Worker pool is full")
		return nil
	}
	p.increaseWorking()
	worker := &Worker{
		name:          "worker_" + rand.Text(),
		handleFunc:    handler,
		numberOfRetry: numberOfRetry,
	}
	worker.Execute(cmd)
	p.decreaseWorking()
	return worker
}

func (p *WorkerPool) IsWorkerPoolFull() bool {
	if !p.Mux.TryLock() {
		return false
	}
	defer p.Mux.Unlock()
	return p.NumberOfWorkings >= p.NumberOfWorkers
}

func (p *WorkerPool) increaseWorking() {
	p.Mux.Lock()
	defer p.Mux.Unlock()
	p.NumberOfWorkings++
}

func (p *WorkerPool) decreaseWorking() {
	p.Mux.Lock()
	defer p.Mux.Unlock()
	p.NumberOfWorkings--
}
