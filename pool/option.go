package workerpool

type Option func(*Pool)

func WithBlock(block bool) Option { // 调用是否阻塞
	return func(p *Pool) {
		p.block = block
	}
}

func WithPreAllocWorkers(preAlloc bool) Option { // 是否预创建 worker
	return func(p *Pool) {
		p.preAlloc = preAlloc
	}
}
