package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// sync  task1

type Counter struct {
	mu    sync.Mutex
	count int
}

func (c *Counter) Increment() {
	c.mu.Lock()         // 获取锁
	defer c.mu.Unlock() // 确保解锁
	c.count++
}

func SyncTask1() {
	var wg sync.WaitGroup
	counter := Counter{}

	// 启动10个协程（网页8的协程数量控制参考）
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 每个协程执行1000次递
			for j := 0; j < 1000; j++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait() // 等待所有协程完成
	fmt.Printf("Final counter value: %d ", counter.count)
}

// sync task2

func SyncTask2() {
	var counter int64     // 声明原子计数器
	var wg sync.WaitGroup // 同步协程组

	wg.Add(10) // 初始化10个协程任务

	// 启动10个协程进行并发递增
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			// 每个协程执行1000次原子递增
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1) // 原子操作无需锁
			}
		}()
	}
	wg.Wait() // 等待所有协程完成
	fmt.Printf("Final Counter: %d\n", atomic.LoadInt64(&counter))
}
