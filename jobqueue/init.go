package jobqueue

import (
	"errors"
	"otter-cloud-ws/constants/api"

	"github.com/EricChiou/jobqueue"
)

type worker struct {
	run  func() interface{}
	wait *chan interface{}
}

func Init() {
	// user job queues
	run(&User.signUp)

	// codemap job queues
	run(&Codemap.add)
}

func Wait() {
	// user job queues
	User.signUp.Wait()

	// codemap job queues
	Codemap.add.Wait()
}

func run(queue *jobqueue.Queue) {
	queue.SetWorker(func(w interface{}) {
		if w, ok := w.(worker); ok {
			*w.wait <- w.run()
		} else {
			*w.wait <- errors.New(string(api.ServerError))
		}
	})
	queue.Run()
}
