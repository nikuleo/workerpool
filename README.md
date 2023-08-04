## goroutine 线程池
This is a tiny goroutine workerpool.
> reference:
> bigwhite
> [fasthttp](https://github.com/valyala/fasthttp/blob/master/workerpool.go)


## workerpool 主要功能
- pool 创建与销毁
- goroutine 管理复用
- 任务调动

## option 可选参数 （functional option）
> https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html

使用函数式的方式包装可选参数，返回闭包函数
```go
func WithBlock(block bool) Option { // 调用是否阻塞
	return func(p *Pool) {
		p.block = block
	}
}
```