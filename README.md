# sched 并发任务调度库

简洁, 高效, 并发限制, 复用 goroutine

*forked from [polyred/internal/sched at main · polyred/polyred (github.com)](https://github.com/polyred/polyred/tree/main/internal/sched)*

## 变动

- 取消 go1.18+ 限制
- 默认并发数为: `runtime.NumCPU()`
- 增加方法: `IsRunning()` `WaitAndRelease()`

## 使用

```go
package sched // import "github.com/fufuok/utils/sched"

type Option func(w *Pool)
    func Randomizer(f func(min, max int) int) Option
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
	s.Add(1)
	s.Run(func() {
		// is running: true 1
		fmt.Println("is running:", s.IsRunning(), s.Running())
	})
	s.Wait()
	s.Release()
	// is running: false
	fmt.Println("is running:", s.IsRunning())

	// 指定并发数
	s = sched.New(sched.Workers(2))
	s.Add(5)
	for i := 0; i < 5; i++ {
		s.Run(func() {
			fmt.Println(time.Now())
			time.Sleep(time.Second)
		})
	}
	s.WaitAndRelease()
	// is running: false
	fmt.Println("is running:", s.IsRunning())
}
```











*ff*

