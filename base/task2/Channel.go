package main

import (
	"fmt"
	"sync"
)

// task 1
// 生产者协程：生成1-10的整数发送到通道
func producer(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done() // 确保协程结束时通知WaitGroups

	for i := 1; i <= 10; i++ {
		ch <- i // 发送数据到通道（阻塞直到消费者接收）
		fmt.Printf("Sent: %d\n", i)
	}
	close(ch) // 发送完成后关闭通道
}

// 消费者协程：从通道接收并打印数据
func consumer(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done() // 确保协程结束时通知WaitGroup

	for num := range ch {
		fmt.Printf("Received: %d\n", num)
	}
}

// task 2
// 生产者协程：发送1-100的整数到通道
func producerByBuffer(ch chan<- int) {
	for i := 1; i <= 100; i++ {
		ch <- i // 数据写入缓冲通道[1,3](@ref)
	}
	close(ch) // 发送完成后关闭通道（防止消费者死锁）[4,5](@ref)
}

// 消费者协程：从通道接收并打印数据
func consumerByBuffer(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch { // 自动检测通道关闭[2,6](@ref)
		fmt.Printf("Received: %d\n", num)
	}
}
