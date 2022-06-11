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
