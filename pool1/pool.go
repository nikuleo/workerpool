package workerpool

import (
	"errors"
	"fmt"
	"sync"
)

const (
	defaultCapacity = 100
	maxCapacity     = 10000
)

var (
	ErrNoIdleWorkerInPool = errors.New("no idle worker in pool")
	ErrWorkerPoolFreed    = errors.New("wokerpool freed")
)

type Task func()

type Pool struct {
	capacity int
	active   chan struct{}
	tasks    chan Task
	wg       sync.WaitGroup // 销毁时等待所有 worker 退出
	quit     chan struct{}  // 通知各个 worker 退出的信号
}

func New(capacity int) *Pool {
	if capacity <= 0 { // 防御性校验，当传入参数不合理是主动纠错
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}

	p := &Pool{
		capacity: capacity,
		tasks:    make(chan Task),
		quit:     make(chan struct{}),
		active:   make(chan struct{}, capacity),
	}

	fmt.Println("workerpool start")
	go p.run()
	return p
}

func (p *Pool) run() {
	idx := 0 // 区分 worker 的编号
	for {    // 不断循环监听 Pool 内的两个 channel。
		select {
		case <-p.quit: // 接收到 quit channel 信号时退出。
			return
		case p.active <- struct{}{}: // 当 active 可写时，创建一个新的 worker goroutine
			idx++
			p.newWorker(idx)
		}
	}
}

func (p *Pool) newWorker(i int) {
	p.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exit\n", i, err)
				<-p.active
			}
			p.wg.Done()
		}()
		fmt.Printf("worker[%03d]: start\n", i)
		for {
			select {
			case <-p.quit: // 监听 quit
				fmt.Printf("worker[%03d]: exit\n", i)
				<-p.active
				return
			case t := <-p.tasks:
				fmt.Printf("worker[%03d]: receive a task\n", i)
				t()
			}
		}
	}()
}

func (p *Pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- t:
		return nil
	}
}

func (p *Pool) Free() {
	close(p.quit)
	p.wg.Wait()
	fmt.Printf("workerpool freed\n")
}
