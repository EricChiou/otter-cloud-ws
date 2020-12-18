package jobqueue

import (
	"github.com/EricChiou/jobqueue"
)

type user struct {
	signUp jobqueue.Queue
}

func (u *user) NewSignUpJob(run func() interface{}) error {
	wait := make(chan interface{})
	u.signUp.Add(worker{run: run, wait: &wait})

	result := <-wait
	switch result.(type) {
	case error:
		return result.(error)
	default:
		return nil
	}
}

var User user = user{
	signUp: jobqueue.New(1024),
}
