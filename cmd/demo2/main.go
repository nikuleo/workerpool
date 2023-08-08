package main

import (
	"fmt"
	"sync"
	"time"
	workerpool "workerpool/pool"
)

func demoFunc() {
	time.Sleep(time.Second * 3)
}

func main() {
	p := workerpool.New(5, workerpool.WithPreAllocWorkers(false), workerpool.WithBlock(true))
	defer p.Free()

	var wg sync.WaitGroup

	time.Sleep(time.Second * 3)
	for i := 0; i < 15; i++ {
		wg.Add(1)
		task := func() {
			demoFunc()
			wg.Done()
		}
		err := p.Schedule(task)
		if err != nil {
			fmt.Printf("task[%d]: error: %s\n", i, err.Error())
		}
	}

}
