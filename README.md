# sched 并发任务调度库

简洁, 高效, 并发限制, 复用 goroutine

*forked from [polyred/internal/sched at main · polyred/polyred (github.com)](https://github.com/polyred/polyred/tree/main/internal/sched)*

## 变动

- 取消 go1.18+ 使用限制
- 取消随机指定 goroutine 运行机制, 改为 goroutine 饥饿模式, 更高效
- 取消每 Worker 固定任务缓冲数设定, 改为总缓冲数配置, 待执行任务数超过该值会触发阻塞
- 默认并发数为: `runtime.NumCPU()`
- 增加方法: `p.IsRunning()` `p.WaitAndRelease()` `Queues(limit int)`

## 使用

```go
package sched // import "github.com/fufuok/sched"

type Option func(w *Pool)
    func Queues(limit int) Option
    func Workers(limit int) Option
type Pool struct{ ... }
    func New(opts ...Option) *Pool
```

示例: [examples](examples)

```go
package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/fufuok/sched"
)

func main() {
	// 默认并发数: runtime.NumCPU()
	s := sched.New()
	sum := uint32(0)
	for i := 0; i < 10; i++ {
		s.Add(1)
		s.RunWithArgs(func(n interface{}) {
			atomic.AddUint32(&sum, uint32(n.(int)))
		}, i)
	}
	s.Wait()
	// sum: 45
	fmt.Println("sum:", atomic.LoadUint32(&sum))

	// 继续下一批任务
	fn := func() {
		// is running: true 2
		fmt.Println("is running:", s.IsRunning(), s.Running())
	}
	s.Add(2)
	s.Run(fn, fn)
	s.Wait()
	s.Release()

	// is running: false
	fmt.Println("is running:", s.IsRunning())

	// 指定并发数 * 指定队列缓冲数 < 总任务数时, 会产生阻塞排队
	s = sched.New(sched.Workers(2), sched.Queues(1))
	fmt.Println("start:", time.Now().Format(time.RFC3339Nano))
	for i := 0; i < 5; i++ {
		s.Add(1)
		s.Run(func() {
			fmt.Println(i, time.Now().Format("04:05"))
			time.Sleep(time.Second)
		})
	}
	fmt.Println("not-blocking:", time.Now().Format(time.RFC3339Nano))
	s.WaitAndRelease()
	fmt.Println("done:", time.Now().Format(time.RFC3339Nano))

	// is running: false
	fmt.Println("is running:", s.IsRunning())
}
```











*ff*

