package utils

import (
	"fmt"
	"worker-service/pkg/constant"

	"github.com/rs/zerolog/log"
)

type Worker struct {
	name          string
	handleFunc    constant.WorkerHandleFunc
	numberOfRetry int
}

func (w *Worker) GetName() string {
	return w.name
}

func (w *Worker) Execute(cmd any) {
	err := w.handleFunc(cmd)
	if err == nil || w.numberOfRetry == 0 {
		return
	}
	log.Info().Msg("Retry worker: " + w.name + " with error: " + err.Error() + " and retry time: " + fmt.Sprint(w.numberOfRetry))
	w.numberOfRetry--
	w.Execute(cmd)
}
