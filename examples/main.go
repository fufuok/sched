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
