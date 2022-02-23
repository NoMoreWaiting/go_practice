package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// 使用go完成协程池pool
/*

ceo(main)给你的领导(dispatcher)分配任务,你的领导(dispatcher)再把任务分配给你(worker),你再去执行具体的任务(playload)

编程来源于现实啊
分布式: 众多人员的分工合作
分库分表: 分部门分项目
多层架构: 现实中的公司组织架构

线程池: 众多码农
任务池: 一堆子coding任务
线程池管理者: 主管, 分配任务
任务池管理者: 需求开发
总调度者: 老板

*/



// 声明成顾客
type Payload struct {
	name string
}

// 顾客就餐
func (p *Payload) Play(waiter string) {
	fmt.Printf("waiter: %s 服务 %s 就餐中... 当前任务完成\n", waiter, p.name)
}

// 任务
type Job struct {
	Payload Payload
}

// 服务员
type Worker struct {
	name       string        // 服务员的名字
	WorkerPool chan chan Job // 线程池
	JobChannel chan Job      // 任务池
	quit       chan bool     // 停止退出
}

// 新建一个服务员
func NewWorker(workerPool chan chan Job, name string) Worker {
	fmt.Printf("创建了一个服务员, 他的名字是: %s \n", name)
	return Worker{
		name:       name,            // 服务员的名字
		WorkerPool: workerPool,      // 服务员在哪个线程池里面工作, 可以理解为部门
		JobChannel: make(chan Job),  // 服务员的任务池, 这里是无缓存的, 可以变为有缓存的, 注意有缓存时的处理
		quit:       make(chan bool), // 退出消息. 这里需要逗号, 格式需要
	}
}

// 服务员开始工作
func (w *Worker) Start() {
	// 开一个协程
	go func() {
		for {
			// 注册到线程池中
			w.WorkerPool <- w.JobChannel
			fmt.Printf("%s 准备就绪, 等待任务 \n", w.name)
			select {
			// 接收到了新任务
			case job := <-w.JobChannel:
				fmt.Printf("%s 接到新任务, 当前任务长度是 %d\n", w.name, len(w.WorkerPool))
				job.Payload.Play(w.name)
				time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
				// 任务退出
			case <-w.quit:
				return
				// 注意, 这里的default要慎用, 此协程一直在阻塞等待channel触发, 如果有default, 就会进入下一次循环. select的超时机制
				//default:
				//	fmt.Printf("%s 进行默认任务 ----- \n", w.name)
			}
		}
	}()
}

// 服务员结束工作
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

// 大堂经理
type Dispatcher struct {
	name       string        // 大堂经理的名字
	maxWorkers int           // 获取调度服务员的个数
	// 通道的通道，还是通道
	WorkerPool chan chan Job // 注册和服务员一样的通道
}

// 新建一个大堂经理
func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		WorkerPool: pool,       // 将服务员放到一个池中, 可以理解成一个部门
		name:       "李经理",      // 大堂经理的名字
		maxWorkers: maxWorkers, // 这个大堂经理管理多少服务员
	}
}

// 调度者开始工作
func (d *Dispatcher) Run() {
	// 开始运行
	for i := 0; i < d.maxWorkers; i++ {
		// 注意循环变量的陷阱, 闭包时, 引用循环变量是引用的变量的地址,  最后所有的引用都是循环变量最后的值.
		// 解决办法: 不能直接引用循环变量.
		// 1. 使用函数参数传入, 函数参数在传入的时候会进行计算
		// 2. 重新赋值计算出一个新变量
		worker := NewWorker(d.WorkerPool, fmt.Sprintf("服务员-%s", strconv.Itoa(i)))
		// 工人开始工作
		worker.Start()
	}
	// 监控
	go d.dispatch()
}

// 大堂经理分配任务
func (d *Dispatcher) dispatch() {
	for {
		fmt.Printf("%s, 等待中...\n", d.name)

		select {
		case job := <-JobQueue:
			fmt.Printf("%s, 接收到一个工作任务\n", d.name)
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			// 调度者接收到一个工作任务
			go func(job Job) {
				// 从现有的线程池中拿出一个
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		// 有default, select会出来
		//default:
		//	fmt.Printf("%s, 等待中...\n", d.name)
		}
		fmt.Printf("%s, 分配任务完成...\n", d.name)
	}
}

// 任务队列
var JobQueue chan Job

// 初始化线程池
func Initialize() {
	maxWorkers := 2
	maxQueue := 4
	//初始化一个调度者, 并制定可以调度的服务员数
	dispatcher := NewDispatcher(maxWorkers)
	JobQueue = make(chan Job, maxQueue) // 指定任务的队列
	dispatcher.Run()                    // 一直运行
}

// 运行函数
func RoutinuePoolRun() {
	// 初始化线程池
	Initialize()
	for i := 0; i < 10; i++ {
		p := Payload{
			fmt.Sprintf("顾客-%s", strconv.Itoa(i)),
		}
		JobQueue <- Job{
			Payload: p,
		}
		time.Sleep(time.Second)
	}
	close(JobQueue)
}


