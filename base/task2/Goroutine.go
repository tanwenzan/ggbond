package main

import (
	"fmt"
	"sync"
	"time"
)

func printEven() {
	for i := 0; i <= 10; i += 2 {
		fmt.Println(i)
	}
}

func printOdd() {
	for i := 1; i <= 10; i += 2 {
		fmt.Println(i)
	}
}

type Task struct {
	ID   string
	Func func() error
}

type Scheduler struct {
	tasks       []Task
	results     map[string]time.Duration
	concurrency int
	mu          sync.Mutex
	wg          sync.WaitGroup
}

func NewScheduler(concurrency int) *Scheduler {
	return &Scheduler{
		results:     make(map[string]time.Duration),
		concurrency: concurrency,
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.tasks = append(s.tasks, task)
}

func (s *Scheduler) worker(taskChan <-chan Task) {
	for task := range taskChan {
		// 执行任务并统计时间
		start := time.Now()
		if err := task.Func(); err != nil {
			fmt.Printf("Task %s failed: %v\n", task.ID, err)
		}
		elapsed := time.Since(start)

		// 线程安全写入结果
		s.mu.Lock()
		s.results[task.ID] = elapsed
		s.mu.Unlock()

		s.wg.Done()
	}
}

func (s *Scheduler) Run() {
	taskChan := make(chan Task, len(s.tasks))

	// 启动协程池
	for i := 0; i < s.concurrency; i++ {
		go s.worker(taskChan)
	}

	// 分发任务
	s.wg.Add(len(s.tasks))
	for _, task := range s.tasks {
		taskChan <- task
	}
	close(taskChan)

	// 等待所有任务完成
	s.wg.Wait()
}

func (s *Scheduler) PrintResults() {
	fmt.Println("\nTask Execution Times:")
	for id, duration := range s.results {
		fmt.Printf("• %-10s → %v\n", id, duration)
	}
}
