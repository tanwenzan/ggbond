package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// pointer
	num := 1
	i := &num
	inc10(i)
	fmt.Println(num)
	//
	arr := []int{1, 2, 3}
	double(&arr)
	fmt.Println(arr)

	// goroutine task1
	go printEven()
	go printOdd()
	time.Sleep(3 * time.Second)

	// goroutine task2
	// 初始化调度器（最大并发数=3）
	scheduler := NewScheduler(3)
	// 添加示例任务
	scheduler.AddTask(Task{
		ID: "Task1",
		Func: func() error {
			time.Sleep(1 * time.Second)
			return nil
		},
	})
	scheduler.AddTask(Task{
		ID: "Task2",
		Func: func() error {
			time.Sleep(2 * time.Second)
			return nil
		},
	})
	// 运行调度器
	scheduler.Run()
	// 输出统计结果
	scheduler.PrintResults()

	// OOP task 1
	rect := Rectangle{Width: 5, Height: 3}
	circle := Circle{Radius: 4}
	printShapeInfo(rect)
	printShapeInfo(circle)

	// OOP task2
	emp := Employee{
		Person: Person{
			Name: "Alice",
			Age:  28,
		},
		EmployeeID: 1001,
	}
	emp.PrintInfo()

	// channel task1
	// 创建无缓冲通道（同步通信）
	ch := make(chan int)

	// 使用WaitGroup等待协程完成
	var wg sync.WaitGroup
	wg.Add(2) // 等待生产者和消费者两个协程

	// 启动生产者协程
	go producer(ch, &wg)

	// 启动消费者协程
	go consumer(ch, &wg)

	// 等待所有协程完成
	wg.Wait()
	fmt.Println("All tasks completed")

	// channel task2
	const bufferSize = 10 // 缓冲区大小可调整
	bufferCh := make(chan int, bufferSize)
	var bufferWg sync.WaitGroup

	bufferWg.Add(1)
	go producerByBuffer(bufferCh)            // 启动生产者
	go consumerByBuffer(bufferCh, &bufferWg) // 启动消费者

	bufferWg.Wait() // 等待消费者完成
	fmt.Println("All done!")

	// sync task1
	SyncTask1()

	// sync task2
	SyncTask2()
}
