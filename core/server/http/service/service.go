package service

import (
	"onesite/core/dao"
	"onesite/core/worker"
)

type Service struct {
	Dao    *dao.Dao
	Worker *worker.Worker
}

func NewService(d *dao.Dao, w *worker.Worker) *Service {
	return &Service{
		d,
		w,
	}
}
