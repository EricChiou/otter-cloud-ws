package jobqueue

import (
	"github.com/EricChiou/jobqueue"
)

type codemap struct {
	add jobqueue.Queue
}

func (u *codemap) NewAddJob(run func() interface{}) error {
	wait := make(chan interface{})
	u.add.Add(worker{run: run, wait: &wait})

	result := <-wait
	switch result.(type) {
	case error:
		return result.(error)
	default:
		return nil
	}
}

var Codemap codemap = codemap{
	add: jobqueue.New(1024),
}
