package main

import (
	"fmt"
	"runtime"
	"sync"
)


var count int = 0

func counter(lock *sync.Mutex) {
	lock.Lock()
	count++
	fmt.Println(count)
	lock.Unlock()
}

// 通过golang中的 goroutine 与sync.Mutex进行 并发同步
func test_01() {
	lock := &sync.Mutex{}
	for i := 0; i < 10; i++ {
		//传递指针是为了防止 函数内的锁和 调用锁不一致
		go counter(lock)
	}
	for {
		lock.Lock()
		c := count
		lock.Unlock()
		///把时间片给别的goroutine  未来某个时刻运行该routine
		runtime.Gosched()
		if c >= 10 {
			fmt.Println("goroutine end")
			break
		}
	}
	fmt.Println("test_01  end")
}


// goroutine之间通过 channel进行通信,channel是和类型相关的 可以理解为  是一种类型安全的管道。 简单的channel 使用
func test_02() {
	chs := make([]chan int, 10)
	for i := 0; i < 10; i++ {
		chs[i] = make(chan int)
		go Count(chs[i])
		fmt.Println("Count",i)
	}
	for i, ch := range chs {
		<-ch
		fmt.Println("Counting",i)
	}
}

func Count(ch chan int) {
	ch <- 1
	fmt.Println("Counting  func Count")
}


// 这个示例程序展示如何创建 goroutine 以及调度器的行为
func listing01(){
	// 调度器根据逻辑处理器的数量分配同时执行的 goroutine 任务限制
	// 分配逻辑处理器给调度器使用
	runtime.GOMAXPROCS(2)

	// wg 用来等待程序完成
	// 计数加 2，表示要等待两个 goroutine
	var wg sync.WaitGroup
	wg.Add(2)

	// 声明一个匿名函数，并创建一个 goroutine
	go func(){
		// 在函数退出时调用 Done 来通知 main 函数工作已经完成
		defer wg.Done()

		for count :=0; count < 9; count++{
			for char:= 'a'; char < 'a' + 26; char++{
				fmt.Printf("%c ", char)
			}
		}
	}()

	go func(){
		// 在函数退出时调用 Done 来通知 main 函数工作已经完成
		defer wg.Done()

		for count :=0; count < 9; count++{
			for char:= 'A'; char < 'A' + 26; char++{
				fmt.Printf("%c ", char)
			}
		}
	}()

	fmt.Println("Waiting To Finish")
	// 一旦两个匿名函数创建 goroutine 来执行，main 中的代码会继续运行
	// 等待 goroutine 结束
	wg.Wait()
	fmt.Println("\nTerminating Program")
}